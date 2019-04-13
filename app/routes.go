package app

import (
	"net/http"

	"github.com/rlongo/ictf-gradings-backend/api"
	"github.com/rlongo/ictf-gradings-backend/app/handler"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	ReqAuth     bool
	HandlerFunc func(api.StorageService) http.HandlerFunc
}

type routes []route

var appRoutes = routes{
	route{"BeltTest Indicies", "GET", "/tests", false, handler.GetBeltTests},
	route{"BeltTest Index", "GET", "/test/{id}", true, handler.GetBeltTest},
	route{"BeltTest Create", "POST", "/test", true, handler.CreateBeltTest},
}
