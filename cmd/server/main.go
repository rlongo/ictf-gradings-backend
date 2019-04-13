package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/urfave/negroni"

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

	router := app.NewRouter(storageService, nil)
	n := negroni.Classic()
	n.UseHandler(router)

	log.Printf("listening on IPv4 address \"0.0.0.0\", port %s", port)
	log.Printf("listening on IPv6 address \"::\", port %s", port)

	s := &http.Server{
		Addr:           ":" + port,
		Handler:        n,
		WriteTimeout:   time.Second * 15,
		ReadTimeout:    time.Second * 15,
		IdleTimeout:    time.Second * 60,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(s.ListenAndServe())
}
