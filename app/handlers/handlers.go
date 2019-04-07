package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

// PostMaxSize is the biggest size we will allow for a post
const PostMaxSize = 2014 * 8 * 20 // 20KB

// setErrorResponse is used to reply with an error
func setErrorResponse(w http.ResponseWriter, err error) {
	log.Print(err.Error())
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(422) // unprocessable entity
	if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
		panic(err)
	}
}
