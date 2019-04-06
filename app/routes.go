package app

import (
    "net/http"
	"github.com/rlongo/itcf-gradings-backend/api"
)

type Route struct {
	Name string
	Method string
	Pattern string
	HandlerFunc func(api.StorageService) http.HandlerFunc
}

type Routes []Route

var routes = Routes {
	Route{"TestIndicies", "GET", "/tests", HandlerGetBeltTests},
	Route{"TestIndex", "GET", "/test/{id}", HandlerGetBeltTest},
	Route{"CreateTest", "POST", "/test", HandlerCreateBeltTest},
}