package httputil

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

func RespondWithError(w http.ResponseWriter, msg string, statusCode int) {
	respBody := errorResponse{
		Error: msg,
	}

	dat, err := json.Marshal(respBody)
	if err != nil {
		// Even if marshalling fails, we can still send a basic error
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(dat)
}

func RespondWithJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	dat, err := json.Marshal(data)
	if err != nil {
		// If marshalling fails, fall back to the error handler
		RespondWithError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(dat)
}
