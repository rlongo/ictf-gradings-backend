package app

import (
	"net/http"

	"github.com/rlongo/ictf-gradings-backend/api"
	"github.com/rlongo/ictf-gradings-backend/app/handler"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc func(api.StorageService) http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{"BeltTest Indicies", "GET", "/tests", handler.GetBeltTests},
	Route{"BeltTest Index", "GET", "/test/{id}", handler.GetBeltTest},
	Route{"BeltTest Create", "POST", "/test", handler.CreateBeltTest},
}
