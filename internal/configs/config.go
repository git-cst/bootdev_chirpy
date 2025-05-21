package configs

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/git-cst/bootdev_chirpy/internal/database"
	"github.com/git-cst/bootdev_chirpy/internal/httputil"
)

type ApiConfig struct {
	fileserverHits atomic.Int32
	Db             *database.Queries
	Platform       string
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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	response := fmt.Sprintf(MetricsTemplate, cfg.fileserverHits.Load())
	w.Write([]byte(response))
}

func (cfg *ApiConfig) HandlerResetMetrics(w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" {
		httputil.RespondWithError(w, "Forbidden", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	cfg.fileserverHits.Store(0)
	response := fmt.Sprintf("Hits reset to %d\n", cfg.fileserverHits.Load())
	w.Write([]byte(response))
}

func (cfg *ApiConfig) HandlerResetUsers(w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" {
		httputil.RespondWithError(w, "Forbidden", http.StatusForbidden)
		return
	}

	ctx := r.Context()
	err := cfg.Db.ResetUsers(ctx)
	if err != nil {
		httputil.RespondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respBody := struct {
		Message string `json:"message"`
	}{
		Message: "Users reset",
	}

	httputil.RespondWithJSON(w, respBody, http.StatusOK)
}
