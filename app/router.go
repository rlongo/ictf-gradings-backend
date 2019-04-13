package app

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rlongo/ictf-gradings-backend/api"
)

// NewRouter exports a new router class and used Dependencu Injection to introduce
// any externally required items
func NewRouter(storage api.StorageService, authenticator mux.MiddlewareFunc) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range appRoutes {
		var handler http.Handler
		handler = route.HandlerFunc(storage)

		if route.ReqAuth == true && authenticator != nil {
			handler = authenticator(handler)
		}

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}

	return router
}
