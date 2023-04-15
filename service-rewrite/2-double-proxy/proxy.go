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

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if manager.UseOld(r) {
			oldServiceHandler.ServeHTTP(w, r) // HL
		} else {
			newServiceHandler.ServeHTTP(w, r) // HL
		}
	})
}

// END OMIT
