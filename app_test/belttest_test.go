package app_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rlongo/ictf-gradings-backend/api"
	"github.com/rlongo/ictf-gradings-backend/app"
	"github.com/rlongo/ictf-gradings-backend/storage/mock"
)

type authenticator bool

func (a *authenticator) authenticate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*a = true
		h.ServeHTTP(w, r)
	})
}

func (a *authenticator) assertRequestedAuthentication(t *testing.T, requested bool) {
	if (bool)(*a) != requested {
		msg := "use"
		if !requested {
			msg = "bypass"
		}
		t.Errorf("%s: Was expecting to %s authentication", t.Name(), msg)
	}
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func assertResponseBodyTests(t *testing.T, got []byte, want api.BeltTests) {
	t.Helper()
	var results api.BeltTests

	if err := json.Unmarshal(got, &results); err != nil {
		t.Errorf("Response was invalid JSON")
	}

	if len(results) != len(want) {
		t.Errorf("Response size is wrong. Expected: %d, Got: %d", len(want), len(results))
	}

	for i := range results {
		if results[i].Name != want[i].Name {
			t.Errorf("Response mismatch at index %d. Expected: '%s', Got: '%s'",
				i, want[i].Name, results[i].Name)
		}
	}
}

func assertResponseBodyTest(t *testing.T, got []byte, want api.BeltTest) {
	t.Helper()
	var result api.BeltTest

	if err := json.Unmarshal(got, &result); err != nil {
		t.Errorf("Response was invalid JSON")
	}

	if result.Name != want.Name {
		t.Errorf("Response mismatch. Expected: '%s', Got: '%s'",
			want.Name, result.Name)
	}
}

func TestGETBeltTests(t *testing.T) {
	expected := api.BeltTests{
		&api.BeltTest{ID: 0, Name: "test1", Date: 1, Location: "", Admins: nil},
		&api.BeltTest{ID: 1, Name: "test2", Date: 2, Location: "", Admins: nil},
		&api.BeltTest{ID: 2, Name: "test3", Date: 3, Location: "", Admins: nil},
	}

	storageService := mock.MockStorageService{BeltTestsDB: expected}

	t.Run("ignores authentication", func(t *testing.T) {
		var auth authenticator
		router := app.NewRouter(&storageService, auth.authenticate)

		request, _ := http.NewRequest(http.MethodGet, "/tests", nil)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		auth.assertRequestedAuthentication(t, false)
	})

	t.Run("returns Existing Tests", func(t *testing.T) {
		router := app.NewRouter(&storageService, nil)
		request, _ := http.NewRequest(http.MethodGet, "/tests", nil)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBodyTests(t, response.Body.Bytes(), expected)
	})
}

func TestGETBeltTest(t *testing.T) {
	expected := api.BeltTests{
		&api.BeltTest{ID: 0, Name: "test1", Date: 1, Location: "", Admins: nil},
		&api.BeltTest{ID: 1, Name: "test2", Date: 2, Location: "", Admins: nil},
		&api.BeltTest{ID: 2, Name: "test3", Date: 3, Location: "", Admins: nil},
	}

	storageService := mock.MockStorageService{BeltTestsDB: expected}

	t.Run("requires authentication", func(t *testing.T) {
		var auth authenticator
		router := app.NewRouter(&storageService, auth.authenticate)

		expectedTest := expected[2]
		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/test/%d", expectedTest.ID), nil)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		auth.assertRequestedAuthentication(t, true)
	})

	t.Run("returns Existing Test", func(t *testing.T) {
		router := app.NewRouter(&storageService, nil)

		expectedTest := expected[2]
		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/test/%d", expectedTest.ID), nil)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBodyTest(t, response.Body.Bytes(), *expectedTest)
	})

	t.Run("returns 404 on Missing Test", func(t *testing.T) {
		router := app.NewRouter(&storageService, nil)

		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/test/%d", len(expected)+12), nil)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})

	t.Run("returns 404 on Invalid Test ID", func(t *testing.T) {
		router := app.NewRouter(&storageService, nil)

		request, _ := http.NewRequest(http.MethodGet, "/test/dinosaur", nil)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})
}

func TestPOSTBeltTest(t *testing.T) {

	t.Run("requires authentication", func(t *testing.T) {
		storageService := mock.MockStorageService{BeltTestsDB: nil}
		var auth authenticator
		router := app.NewRouter(&storageService, auth.authenticate)

		request, _ := http.NewRequest(http.MethodPost, "/test", nil)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		auth.assertRequestedAuthentication(t, true)
	})

	t.Run("returns 201 on Valid POST", func(t *testing.T) {
		expectedTest := api.BeltTest{ID: 0, Name: "test1", Date: 1, Location: "", Admins: nil}
		storageService := mock.MockStorageService{BeltTestsDB: nil}
		router := app.NewRouter(&storageService, nil)

		expectedTestJSON, _ := json.Marshal(expectedTest)
		b := bytes.NewBuffer(expectedTestJSON)

		request, _ := http.NewRequest(http.MethodPost, "/test", b)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusCreated)

		// Check that it's in
		if len(storageService.BeltTestsDB) != 1 {
			t.Fatal("Record wasn't inserted")
		}

		if storageService.BeltTestsDB[0].Name != expectedTest.Name {
			t.Errorf("Corrupted Post. Expected: '%s', Got: '%s'",
				expectedTest.Name, storageService.BeltTestsDB[0].Name)
		}
	})

	t.Run("returns 400 on an Invalid POST", func(t *testing.T) {
		storageService := mock.MockStorageService{BeltTestsDB: nil}
		router := app.NewRouter(&storageService, nil)

		expectedTestJSON, _ := json.Marshal("foo")
		b := bytes.NewBuffer(expectedTestJSON)

		request, _ := http.NewRequest(http.MethodPost, "/test", b)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("returns 400 on an Empty POST", func(t *testing.T) {
		storageService := mock.MockStorageService{BeltTestsDB: nil}
		router := app.NewRouter(&storageService, nil)

		request, _ := http.NewRequest(http.MethodPost, "/test", nil)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusBadRequest)
	})
}
