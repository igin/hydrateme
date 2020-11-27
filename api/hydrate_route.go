package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// StartHydrationRequest is a request to start hydration alerts
type StartHydrationRequest struct {
	User string `json:"user"`
}

// GetHydrateRoute handles creation of hydration alerts
func GetHydrateRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Let's hydrate!")
}

// CreateHydrateRoute handles creation of hydration alerts
func CreateHydrateRoute(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var startRequest StartHydrationRequest
	err := decoder.Decode(&startRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if startRequest.User == "" {
		http.Error(w, "user is required", http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, "Let's hydrate!")
}
