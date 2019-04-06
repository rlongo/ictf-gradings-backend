 
package main

import (
	"log"
	"net/http"
	"github.com/rlongo/itcf-gradings-backend/storage/psql"
	"github.com/rlongo/itcf-gradings-backend/app"
	"os"
)

func main() {

	storageService, err := psql.Open(os.Getenv("DATABASE_URL"))
	if err!=nil {
		panic(err)
	}
	defer storageService.Close()

	port := os.Getenv("PORT")

        if len(port) == 0 {
            panic("env PORT isn't set!")
	}

	log.Printf("listening on IPv4 address \"0.0.0.0\", port %s", port)
	log.Printf("listening on IPv6 address \"::\", port %s", port)

	router := app.NewRouter(storageService)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
