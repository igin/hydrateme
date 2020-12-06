package testspyserver

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

// SpyServer enables testing code that performs http calls
type SpyServer struct {
	*httptest.Server
	LastRequest *http.Request
	WasCalled   bool
	LastBody    string
}

// NewSpyServer creates and returns a new spy server
func NewSpyServer() *SpyServer {
	server := &SpyServer{
		WasCalled:   false,
		LastRequest: nil,
		LastBody:    "",
		Server:      nil,
	}

	server.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.WasCalled = true
		server.LastRequest = r.Clone(r.Context())
		body, _ := ioutil.ReadAll(r.Body)
		server.LastBody = string(body)
		defer r.Body.Close()
	}))
	return server
}

// NewSpyServerReturning creates a spy server that can react to calls
func NewSpyServerReturning(handler http.HandlerFunc) *SpyServer {
	server := &SpyServer{
		WasCalled:   false,
		LastRequest: nil,
		Server:      nil,
	}

	server.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.WasCalled = true
		server.LastRequest = r

		handler(w, r)
	}))
	return server
}
