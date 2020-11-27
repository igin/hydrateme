package main

import (
	"fmt"
	"net/http"
)

// GetHydrateRoute handles creation of hydration alerts
func GetHydrateRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Let's hydrate!")
}

// CreateHydrateRoute handles creation of hydration alerts
func CreateHydrateRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Let's hydrate!")
}
