package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
)

func TestHealthCheck(t *testing.T) {
	api := &server{}
	api.router = api.appRoutes(chi.NewRouter())

	ts := httptest.NewServer(api.router)
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
