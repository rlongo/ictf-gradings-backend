package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rlongo/ictf-gradings-backend/api"
)

func GetBeltTests(storage api.StorageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if tests, err := storage.AllBeltTests(); err == nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(tests)
		} else {
			setErrorResponse(w, http.StatusInternalServerError, err)
		}

	}
}

func GetBeltTest(storage api.StorageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		var testID = params["id"]
		if testID, err := strconv.ParseInt(testID, 10, 64); err == nil {
			if test, err := storage.GetBeltTest(testID); err == nil {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(test)
			} else {
				setErrorResponse(w, http.StatusNotFound, err)
			}
		} else {
			setErrorResponse(w, http.StatusNotFound, fmt.Errorf("Test ID Wasn't Valid"))
		}
	}
}

func CreateBeltTest(storage api.StorageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var test api.BeltTest

		if r.Body == nil {
			setErrorResponse(w, http.StatusBadRequest, fmt.Errorf("Can't process an empty response"))
			return
		}

		// Read out 20KB of body to avoid spurious crashes
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, PostMaxSize))
		if err != nil {
			panic(err)
		}
		if err := r.Body.Close(); err != nil {
			panic(err)
		}

		if err := json.Unmarshal(body, &test); err != nil {
			setErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		if _, err := storage.CreateBeltTest(test); err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
		} else {
			setErrorResponse(w, http.StatusBadRequest, err)
		}
	}
}
