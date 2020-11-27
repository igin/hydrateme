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
			GetHydrateRoute,
		},
		{
			"Hydrations",
			"POST",
			"/hydrate",
			CreateHydrateRoute,
		},
		{
			"Default Route",
			"GET",
			"/",
			DefaultRoute,
		},
	}
}
