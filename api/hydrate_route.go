package main

import (
	"fmt"
	"net/http"
)

// HydrateRoute handles creation of hydration alerts
func HydrateRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Let's hydrate!")
}
