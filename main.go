package main

import (
	"net/http"

	"github.com/git-cst/bootdev_chirpy/internal/configs"
	"github.com/git-cst/bootdev_chirpy/internal/handlers"
)

func main() {
	apiConfig := &configs.ApiConfig{}

	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fileServer := http.FileServer(http.Dir("./static"))

	mux.Handle("/app/", apiConfig.MiddlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	mux.HandleFunc("GET /api/healthz", handlers.HandlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiConfig.HandlerShowMetrics)
	mux.HandleFunc("POST /admin/reset", apiConfig.HandlerResetMetrics)

	server.ListenAndServe()
}
