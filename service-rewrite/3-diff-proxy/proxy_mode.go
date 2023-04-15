package main

import "net/http"

// START OMIT
type Manager interface {
	GetProxyMode(r *http.Request) ProxyMode // HL
}

type ProxyMode int

const (
	ProxyModeUseOld ProxyMode = iota
	ProxyModeUseNew
	ProxyModeUseOldAndDiff // HL
)

// END OMIT
