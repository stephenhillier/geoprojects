package server

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/stephenhillier/geoprojects/earthworks/db"
)

func TestHealthCheck(t *testing.T) {
	mockDB, _, err := sqlmock.New()

	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	store, err := db.NewDatastore(sqlxDB)
	if err != nil {
		t.Fatalf("Could not create mock datastore")
	}

	svc, err := NewEarthworksService(store, &Config{})
	if err != nil {
		t.Fatal("Could not create earthworks service")
	}

	router := svc.appRoutes(chi.NewRouter())

	ts := httptest.NewServer(router)
	defer ts.Close()

	// check for HTTP 200
	if _, code, _ := testRequest(t, ts, "GET", "/api/v1/health", nil); code != 200 {
		t.Fatalf("Response was %v, expected 200", code)
	}

}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, int, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)
		return nil, 0, ""
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, 0, ""
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return nil, 0, ""
	}
	defer resp.Body.Close()

	return resp, resp.StatusCode, string(respBody)
}
