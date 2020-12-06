package main

import (
	"github.com/gorilla/mux"
)

// BuildRouter builds a router containing all the handlers of this api
func BuildRouter() *mux.Router {
	router := mux.NewRouter()
	routes := BuildRoutes()
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	return router
}
