package app

import (
	"io/ioutil"
	"io"
	"fmt"
	"log"
	"encoding/json"
    "strconv"
    "net/http"
    "github.com/gorilla/mux"
	"github.com/rlongo/ictf-gradings-backend/api"
)

const PostMaxSize = 2014 * 8 * 20 // 20KB

func setErrorResponse(w http.ResponseWriter, err error) {
	log.Print(err.Error())
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(422) // unprocessable entity
	if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
		panic(err)
	}
}

func HandlerGetBeltTests(storage api.StorageService) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		if tests, err := storage.AllBeltTests(); err==nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(tests)
		} else {
			setErrorResponse(w, err)
		}

	}
}

func HandlerGetBeltTest(storage api.StorageService) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		var testID = params["id"]
		if testID, err := strconv.Atoi(testID); err==nil {
			if test, err := storage.GetBeltTest(testID); err==nil {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(test)
			} else {
				setErrorResponse(w, err)
			}
		} else {
			setErrorResponse(w, fmt.Errorf("Test ID Wasn't Valid"))
		}
	}
}

func HandlerCreateBeltTest(storage api.StorageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var test api.BeltTest

		// Read out 20KB of body to avoid spurious crashes
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, PostMaxSize))
		if err != nil {
			panic(err)
		}
		if err := r.Body.Close(); err != nil {
			panic(err)
		}

		if err := json.Unmarshal(body, &test); err!=nil {
			setErrorResponse(w, err)
			return
		}
		
		if _, err := storage.CreateBeltTest(test); err==nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
		} else {
			setErrorResponse(w, err)
		}
	}
}
