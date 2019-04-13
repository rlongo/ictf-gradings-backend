package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/urfave/negroni"

	"github.com/rlongo/ictf-gradings-backend/app"
	"github.com/rlongo/ictf-gradings-backend/storage/psql"
)

func getOSVariable(key string) string {
	if v := os.Getenv(key); len(v) >= 0 {
		return v
	}

	panic(fmt.Sprintf("env %s isn't set!", key))
}

func main() {
	dbURL := getOSVariable("DATABASE_URL")
	port := getOSVariable("PORT")
	auth0AUD := getOSVariable("AUTH0_AUD")
	auth0ISS := getOSVariable("AUTH0_ISS")

	storageService, err := psql.Open(dbURL)
	if err != nil {
		panic(err)
	}
	defer storageService.Close()

	authenticator := NewAuthMiddleware(auth0AUD, auth0ISS)
	router := app.NewRouter(storageService, authenticator.Handler)

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
