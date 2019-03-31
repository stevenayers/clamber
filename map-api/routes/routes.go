package routes

import (
	"github.com/gorilla/mux"
	"go-clamber/map-api/handlers"
	"go-clamber/map-api/logger"
	"net/http"
)

type Routes []Route

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
	Queries     []string
}

var routes = Routes{
	Route{
		"Search",
		"GET",
		"/search",
		handlers.Search,
		[]string{
			"url", "{url}",
			"depth", "{depth}",
			"allow_external_links", "{allow_external_links}",
		},
	},
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = logger.Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler).Queries(route.Queries...)
	}

	return router
}
