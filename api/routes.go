package main

import "net/http"

// Route definition
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// BuildRoutes returns all routes
func BuildRoutes() []Route {
	return []Route{
		{
			"Hydrations",
			"GET",
			"/hydrate",
			HydrateRoute,
		},
		{
			"Default Route",
			"GET",
			"/",
			DefaultRoute,
		},
	}
}
