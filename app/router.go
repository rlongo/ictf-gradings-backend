package app

import (
    "net/http"
    "github.com/gorilla/mux"
	"github.com/rlongo/ictf-gradings-backend/api"
)

func NewRouter(storage api.StorageService) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
    for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc(storage)
		handler = Logger(handler, route.Name)

        router.
            Methods(route.Method).
            Path(route.Pattern).
            Name(route.Name).
            Handler(handler)
    }
 
    return router
}