package handlers

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

type mockStore struct{}

func (ds mockStore) GetEntry(id string) (string, error) {
	if id == "12345678" {
		return "dummy entry", nil
	}
	return "", errors.New("not found")
}

func (ds *mockStore) InsertEntry(id string, contents string) error {
	return errors.New("not implemented")
}

func TestPasteGet(t *testing.T) {
	var pasteTests = []struct {
		id                 string
		httpStatusExpected int
		contentExpected    string
	}{
		// Valid entry
		{
			"12345678",
			http.StatusOK,
			"dummy entry",
		},
		// Non-existent entry
		{
			"missing1",
			http.StatusNotFound,
			"",
		},
	}

	ds := mockStore{}
	router := mux.NewRouter()
	s := defaultServer{
		store:  &ds,
		router: router,
	}
	s.routes()

	for _, tt := range pasteTests {
		req, err := http.NewRequest("GET", "/"+tt.id, nil)
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		if status := w.Code; status != tt.httpStatusExpected {
			t.Fatalf("for ID [%s], handler returned wrong status code: got %v want %v",
				tt.id, status, tt.httpStatusExpected)
		}
		if tt.httpStatusExpected != http.StatusOK {
			continue
		}
		bodyBytes, err := ioutil.ReadAll(w.Body)
		if err != nil {
			t.Fatalf("failed to read HTTP response body: %v", err)
		}
		if tt.contentExpected != string(bodyBytes) {
			log.Fatalf("for ID [%s], got %s, want %s", tt.id, string(bodyBytes), tt.contentExpected)
		}
	}
}