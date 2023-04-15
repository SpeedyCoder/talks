package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	oldServiceURL, err := url.Parse(os.Getenv("OLD_SERVICE_URL"))
	if err != nil {
		return err
	}
	newServiceURL, err := url.Parse(os.Getenv("NEW_SERVICE_URL"))
	if err != nil {
		return err
	}
	handler := newProxyHandler(boolManager(false), oldServiceURL, newServiceURL)
	server := http.Server{
		Addr:    ":3000",
		Handler: handler,
	}
	return server.ListenAndServe()
}

type boolManager bool

func (b boolManager) UseOld(*http.Request) bool {
	return bool(b)
}
