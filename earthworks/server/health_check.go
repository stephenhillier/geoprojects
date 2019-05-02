package server

import "net/http"

// health is a simple health check handler that returns HTTP 200 OK.
func health(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
}
