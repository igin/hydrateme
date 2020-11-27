package main

import (
	"github.com/gorilla/mux"
)

// BuildRouter builds a router containing all the handlers of this api
func BuildRouter() *mux.Router {
	r := mux.NewRouter()
	routes := BuildRoutes()
	for _, route := range routes {
		r.HandleFunc(route.Pattern, route.HandlerFunc)
	}
	return r
}
