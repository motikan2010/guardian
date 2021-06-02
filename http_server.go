package main

import (
	"net/http"
	"os"
	"time"
)

/*HTTPServer The http server handler*/
type HTTPServer struct {
}

/*NewHTTPServer HTTP server initializer*/
func NewHTTPServer() *HTTPServer {
	return &HTTPServer{}
}

func (h *HTTPServer) ServeHTTP() {
	srv := &http.Server{
		ReadHeaderTimeout: 20 * time.Second,
		WriteTimeout:      2 * time.Minute,
		ReadTimeout:       1 * time.Minute,
		Handler:           NewGuardianHandler(),
		Addr:              ":" + os.Getenv("LISTEN_PORT"),
	}

	srv.ListenAndServe()
}
