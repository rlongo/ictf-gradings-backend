package main

import (
	"log"
	"net/http"
	"os"

	"github.com/rlongo/ictf-gradings-backend/app"
	"github.com/rlongo/ictf-gradings-backend/storage/psql"
)

func main() {

	dbUrl := os.Getenv("DATABASE_URL")
	if len(dbUrl) == 0 {
		panic("env DATABASE_URL isn't set!")
	}

	storageService, err := psql.Open(dbUrl)
	if err != nil {
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
