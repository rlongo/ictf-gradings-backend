package app

import (
	"fmt"
	"io"
	"net/http"

	"github.com/urfave/negroni"

	"github.com/gorilla/mux"
	"github.com/rlongo/ictf-gradings-backend/api"
)

// RoleParser adds an extra authentication step, only allowing
// users whoe fulfill the given role to use our endpoint
type RoleParser func(*http.Request) Role

// NewRouter exports a new router class and used Dependencu Injection to introduce
// any externally required items
func NewRouter(storage api.StorageService, middleware *negroni.Negroni, roleParser RoleParser) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	dojangRouter := router.PathPrefix("/api/v1/").Subrouter()
	for _, route := range dojangRoutes {
		var handler http.Handler
		handler = route.HandlerFunc(storage)
		handler = middleware.With(
			negroni.HandlerFunc(handlerRoleAuthentication(route.AuthRole, roleParser)),
			negroni.Wrap(handler),
		)

		dojangRouter.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func PrintRoutes(w io.Writer, router *mux.Router) {
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, _ := route.GetPathTemplate()
		fmt.Fprintf(w, "%s %s\n", route.GetName(), path)
		return nil
	})
}

func handlerRoleAuthentication(requiredRole Role, roleParser RoleParser) func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		requestRole := roleParser(r)

		if roleParser != nil && !requiredRole.Matches(requestRole) {
			rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		next(rw, r)
	}
}
