package models

import (
	"net/http"
	"time"

	"github.com/motikan2010/guardian/helpers"
)

//HTTPLog represents http log
type HTTPLog struct {
	TargetID                  string
	RequestURI                string
	StatusCode                int
	RequestRulesCheckElapsed  int64
	ResponseRulesCheckElapsed int64
	HTTPElapsed               int64
	RequestSize               int64
	ResponseSize              int64

	timer time.Time
}

//NewHTTPLog inits HTTP log
func NewHTTPLog() *HTTPLog {
	return &HTTPLog{"", "", 0, 0, 0, 0, 0, 0, time.Now()}
}

//Build Fills HTTP Log
func (h *HTTPLog) Build(target *Target, request *http.Request, response *http.Response) *HTTPLog {
	h.TargetID = target.ID
	h.RequestURI = request.RequestURI
	h.RequestSize = request.ContentLength

	if response == nil {
		return h
	}

	h.ResponseSize = response.ContentLength
	h.StatusCode = response.StatusCode

	return h
}

//NoResponse handles when no response
func (h *HTTPLog) NoResponse() *HTTPLog {
	h.StatusCode = -1
	h.HTTPElapsed = helpers.CalcTimeNow(h.timer)

	return h
}

//RequestRulesExecutionEnd Calculates the time for execution of rules
func (h *HTTPLog) RequestRulesExecutionEnd() *HTTPLog {
	h.RequestRulesCheckElapsed = helpers.CalcTimeNow(h.timer)

	return h
}

//ResponseRulesExecutionStart Response execution time measure starter
func (h *HTTPLog) ResponseRulesExecutionStart() *HTTPLog {
	h.timer = time.Now()

	return h
}

//ResponseRulesExecutionEnd Response execution time measure ender
func (h *HTTPLog) ResponseRulesExecutionEnd() *HTTPLog {
	h.ResponseRulesCheckElapsed = helpers.CalcTimeNow(h.timer)

	return h
}

//OriginRequestStart Origin request time measure starter
func (h *HTTPLog) OriginRequestStart() *HTTPLog {
	h.timer = time.Now()

	return h
}

//OriginRequestEnd Origin request execution time measure ender
func (h *HTTPLog) OriginRequestEnd() *HTTPLog {
	h.HTTPElapsed = helpers.CalcTimeNow(h.timer)

	return h
}
