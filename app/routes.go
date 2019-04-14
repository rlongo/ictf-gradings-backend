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
	AuthRole    Role
	HandlerFunc func(api.StorageService) http.HandlerFunc
}

type routes []route

var dojangRoutes = routes{
	route{"BeltTest Indicies", "GET", "/dojang/tests", (RoleInstructor | RoleSupervisor), handler.GetBeltTests},
	route{"BeltTest Index", "GET", "/dojang/test/{id}", (RoleInstructor | RoleSupervisor), handler.GetBeltTest},
	route{"BeltTest Create", "POST", "/dojang/test", RoleSupervisor, handler.CreateBeltTest},
}

var adminRoutes = routes{}
