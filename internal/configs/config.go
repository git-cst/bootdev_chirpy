package configs

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type ApiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

const MetricsTemplate = `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`

func (cfg *ApiConfig) HandlerShowMetrics(w http.ResponseWriter, r *http.Request) {
	// Set Header
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Write status code
	w.WriteHeader(http.StatusOK)

	response := fmt.Sprintf(MetricsTemplate, cfg.fileserverHits.Load())

	// Write response body
	w.Write([]byte(response))
}

func (cfg *ApiConfig) HandlerResetMetrics(w http.ResponseWriter, r *http.Request) {
	// Set Header
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// Write status code
	w.WriteHeader(http.StatusOK)

	cfg.fileserverHits.Store(0)

	response := fmt.Sprintf("Hits reset to %d\n", cfg.fileserverHits.Load())

	// Write response body
	w.Write([]byte(response))
}
