package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"

	"github.com/google/go-cmp/cmp"
)

// START DiffHandler OMIT
func newDiffHandler(oldHandler, newHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		oldHandlerReq, newHandlerReq, err := getRequestsToForward(req) // HL
		if err != nil {
			http.Error(w, "failed to read payload", http.StatusInternalServerError)
			return
		}
		oldHandlerWriter, newHandlerWriter := httptest.NewRecorder(), httptest.NewRecorder() // HL

		diffWG := &sync.WaitGroup{}
		diffWG.Add(2)
		// Asynchronously check for differences after both handlers are done.
		go diffResponses(diffWG, oldHandlerWriter, newHandlerWriter) // HL

		go func() {
			defer diffWG.Done()
			newHandler.ServeHTTP(newHandlerWriter, newHandlerReq)
		}()

		defer diffWG.Done()
		oldHandler.ServeHTTP(oldHandlerWriter, oldHandlerReq)
		copyResponse(oldHandlerWriter, w) // HL
	})
}

// END DiffHandler OMIT

// START Requests OMIT
func getRequestsToForward(req *http.Request) (*http.Request, *http.Request, error) {
	payload, err := io.ReadAll(req.Body) // HL
	if err != nil {
		return nil, nil, err
	}
	oldHandlerReq := req.Clone(req.Context())                   // HL
	oldHandlerReq.Body = io.NopCloser(bytes.NewReader(payload)) // HL

	newHandlerReq := req.Clone(context.Background())            // HL
	newHandlerReq.Body = io.NopCloser(bytes.NewReader(payload)) // HL

	return oldHandlerReq, newHandlerReq, nil
}

// END Requests OMIT

// START copyResponse OMIT
func copyResponse(recorder *httptest.ResponseRecorder, w http.ResponseWriter) {
	for name, values := range recorder.Header() {
		for _, val := range values {
			w.Header().Add(name, val) // HL
		}
	}
	w.WriteHeader(recorder.Code)          // HL
	_, _ = w.Write(recorder.Body.Bytes()) // HL
}

// END copyResponse OMIT

// START diffResponses OMIT
func diffResponses(wg *sync.WaitGroup, oldResponse, newResponse *httptest.ResponseRecorder) {
	wg.Wait() // Wait for both requests to finish. // HL

	if diff := cmp.Diff(oldResponse.Header(), newResponse.Header()); diff != "" { // HL
		fmt.Println("Header Diff:", diff) // HL
	} // HL

	var oldJSON, newJSON any

	if err := json.Unmarshal(oldResponse.Body.Bytes(), &oldJSON); err != nil {
		fmt.Printf("failed to unmarshal old json: %s\n", err)
	}
	if err := json.Unmarshal(newResponse.Body.Bytes(), &newJSON); err != nil {
		fmt.Printf("failed to unmarshal old json: %s\n", err)
	}
	if diff := cmp.Diff(oldJSON, newJSON); diff != "" { // HL
		fmt.Println("Body Diff:", diff) // HL
	} // HL
}

// END diffResponses OMIT
