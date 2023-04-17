package service_rewrite

import (
	"io"
	"net/url"
)

// START OMIT

type Handler interface {
	ServeHTTP(w ResponseWriter, req *Request)
}

type ResponseWriter interface {
	Header() Header
	Write([]byte) (int, error)
	WriteHeader(statusCode int)
}

type Request struct {
	Method string
	URL    *url.URL
	Header Header
	Body   io.ReadCloser
	//...
}

type Header map[string][]string

// END OMIT
