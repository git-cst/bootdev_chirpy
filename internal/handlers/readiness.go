package handlers

import "net/http"

func HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	// Set Header
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// Write status code
	w.WriteHeader(http.StatusOK)

	// Write response body
	w.Write([]byte("OK"))
}
