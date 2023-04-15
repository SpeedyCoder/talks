package main

import (
	"net/http"
	"net/http/httputil" // HL
	"net/url"
	"os"
)

func run() error {
	oldServiceURL, err := url.Parse(os.Getenv("OLD_SERVICE_URL"))
	if err != nil {
		return err
	}
	server := http.Server{
		Addr:    ":3000",
		Handler: httputil.NewSingleHostReverseProxy(oldServiceURL), // HL
	}
	return server.ListenAndServe()
}
