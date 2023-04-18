package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

// START OMIT
func newProxyHandler(manager Manager, oldSvcURL, newSvcURL *url.URL) http.Handler {
	oldServiceHandler := httputil.NewSingleHostReverseProxy(oldSvcURL)
	newServiceHandler := httputil.NewSingleHostReverseProxy(newSvcURL)

	diffHandler := newDiffHandler(oldServiceHandler, newServiceHandler) // HL

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch manager.GetProxyMode(r) {

		case ProxyModeUseOld:
			oldServiceHandler.ServeHTTP(w, r)

		case ProxyModeUseNew:
			newServiceHandler.ServeHTTP(w, r)

		case ProxyModeUseOldAndDiff: // HL
			diffHandler.ServeHTTP(w, r) // HL
		}
	})
}

// END OMIT
