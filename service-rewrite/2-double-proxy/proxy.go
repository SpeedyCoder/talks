package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

// START OMIT
type Manager interface {
	UseOld(r *http.Request) bool // HL
}

func newProxyHandler(manager Manager, oldSvcURL, newSvcURL *url.URL) http.Handler { // HL
	oldServiceHandler := httputil.NewSingleHostReverseProxy(oldSvcURL)
	newServiceHandler := httputil.NewSingleHostReverseProxy(newSvcURL)

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if manager.UseOld(req) {
			oldServiceHandler.ServeHTTP(w, req) // HL
		} else {
			newServiceHandler.ServeHTTP(w, req) // HL
		}
	})
}

// END OMIT
