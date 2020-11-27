package main

import "net/http"

// DefaultRoute is the handler for all unknown requests
func DefaultRoute(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}
