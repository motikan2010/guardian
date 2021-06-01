package main

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/motikan2010/guardian/waf/engine"

	"net"
	"net/http"

	"net/url"
	"time"

	"golang.org/x/crypto/acme/autocert"

	"github.com/motikan2010/guardian/response"

	"github.com/motikan2010/guardian/data"
	"github.com/motikan2010/guardian/models"
	"github.com/motikan2010/guardian/request"
)

var dialer = &net.Dialer{
	Timeout:   30 * time.Second,
	KeepAlive: 30 * time.Second,
	DualStack: true,
}

/*GuardianHandler Guardian HTTPS Handler is the transport handler*/
type GuardianHandler struct {
	IsHTTPPortListener bool
	CertManager        *autocert.Manager
	IPRateLimiter      *models.IPRateLimiter
}

/*NewGuardianHandler Https Guardian handler init*/
func NewGuardianHandler(isHTTPPortListener bool, certManager *autocert.Manager, ipRateLimiter *models.IPRateLimiter) *GuardianHandler {
	return &GuardianHandler{isHTTPPortListener, certManager, ipRateLimiter}
}

func (h *GuardianHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	DB := data.NewDBHelper()

	ipAddress := h.getIPAddress(r)

	if !h.IPRateLimiter.IsAllowed(ipAddress) {
		w.WriteHeader(http.StatusTooManyRequests)
		fmt.Fprintf(w, "Too Many Requests. %s", url.QueryEscape(r.URL.Path))

		go DB.LogThrottleRequest(ipAddress)

		return
	}

	target := DB.GetTarget(r.Host)

	if target == nil {
		fmt.Fprintf(w, "Your application not authorized yet! Check your implementation. %s", r.URL.Path)
		fmt.Println("Unauthorized Application requested." + r.Host + r.URL.Path)
		fmt.Println(r)

		return
	}

	if target.AutoCert && h.IsHTTPPortListener {
		fmt.Println("AutoCert in progress. " + r.Host + r.URL.Path)
		h.CertManager.HTTPHandler(nil).ServeHTTP(w, r)
		return
	}

	if target.UseHTTPS && h.IsHTTPPortListener {
		redirectToURI := "https://" + r.Host

		if r.RequestURI != "" {
			redirectToURI += r.RequestURI
		}

		http.Redirect(w, r, redirectToURI, 301)

		return
	}

	httpLog := models.NewHTTPLog()

	requestChecker := request.NewRequestChecker(w, r, target)
	requestIsNotSafe := requestChecker.Handle()

	httpLog = httpLog.RequestRulesExecutionEnd()

	if requestIsNotSafe {
		go h.logHTTPRequest(httpLog.Build(target, r, nil))

		return
	}
	httpLog.OriginRequestStart()
	uriToReq := r.Host

	if target.Proto == 0 {
		uriToReq = "http://" + uriToReq
	} else {
		uriToReq = "https://" + uriToReq
	}

	if r.RequestURI != "" {
		uriToReq += r.RequestURI
	}

	transportResponse := h.transportRequest(uriToReq, w, requestChecker.Transaction, target)

	if transportResponse == nil {
		go h.logHTTPRequest(httpLog.Build(target, r, nil).NoResponse())

		return
	}

	httpLog.OriginRequestEnd()

	httpLog.ResponseRulesExecutionStart()

	responseIsNotSafe := response.NewResponseChecker(w, requestChecker.Transaction, transportResponse, target).Handle()

	httpLog = httpLog.ResponseRulesExecutionEnd()

	if responseIsNotSafe {
		go h.logHTTPRequest(httpLog.Build(target, r, nil))

		return
	}

	h.transformResponse(w, r, transportResponse)

	go h.logHTTPRequest(httpLog.Build(target, r, transportResponse))
}

//TransportRequest Transports the incoming request
func (h *GuardianHandler) transportRequest(uriToReq string,
	incomingWriter http.ResponseWriter,
	transaction *engine.Transaction,
	target *models.Target) *http.Response {

	var response *http.Response
	var err error
	var req *http.Request

	//timeout is 45 secs for to pass to origin server.
	client := &http.Client{
		Timeout: time.Second * 45,
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				//TODO: Check better solutions for dialcontext like timeouts.
				uri, ferr := url.Parse(addr)
				if ferr != nil {
					panic(ferr)
				}

				addr = target.OriginIPAddress + ":" + uri.Opaque

				return dialer.DialContext(ctx, network, addr)
			},
		},
	}

	req, err = http.NewRequest(transaction.Request.Method, uriToReq, bytes.NewBuffer(transaction.RequestBodyProcessor.GetBodyBuffer()))
	for name, value := range transaction.Request.Header {
		//TODO: Do not pass the headers except whitelisted
		if name == "X-Forwarded-For" {
			continue
		}

		req.Header.Set(name, value[0])
	}

	fwIP := h.getIPAddress(transaction.Request)
	if fwIP != "" {
		req.Header.Set("X-Forwarded-For", fwIP)
	}

	response, err = client.Do(req)

	if err != nil {
		http.Error(incomingWriter, err.Error(), http.StatusInternalServerError)
		return nil
	}

	return response
}

func (h *GuardianHandler) transformResponse(incomingWriter http.ResponseWriter, incomingRequest *http.Request, response *http.Response) {
	for k, v := range response.Header {
		incomingWriter.Header().Set(k, v[0])
	}
	incomingWriter.WriteHeader(response.StatusCode)
	io.Copy(incomingWriter, response.Body)
	defer incomingRequest.Body.Close()
}

func (h *GuardianHandler) logHTTPRequest(log *models.HTTPLog) {
	data.NewDBHelper().LogHTTPRequest(log)
}

func (h *GuardianHandler) getIPAddress(incomingRequest *http.Request) string {

	//TODO: IP forwarding must be enabled for the target
	/*
		ipAddress := incomingRequest.Header.Get("X-Real-Ip")
		if ipAddress == "" {
			ipAddress = incomingRequest.Header.Get("X-Forwarded-For")
		}
		if ipAddress == "" {
			ipAddress = incomingRequest.RemoteAddr
		}

		return ipAddress
	*/

	return incomingRequest.RemoteAddr

}
