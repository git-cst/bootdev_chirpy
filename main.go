package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/git-cst/bootdev_chirpy/internal/configs"
	"github.com/git-cst/bootdev_chirpy/internal/database"
	"github.com/git-cst/bootdev_chirpy/internal/handlers"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatal(fmt.Sprintf("Could not access DB. Err: %v", err))
	}

	dbQueries := database.New(db)

	apiConfig := &configs.ApiConfig{
		Db:       dbQueries,
		Platform: platform,
	}

	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fileServer := http.FileServer(http.Dir("./static"))

	mux.Handle("/app/", apiConfig.MiddlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	mux.HandleFunc("GET /api/healthz", handlers.HandlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlers.HandlerValidate)
	mux.HandleFunc("POST /api/users", handlers.MakeHandlerWithConfig(apiConfig, handlers.HandlerUsers))

	/* admin commands */
	mux.HandleFunc("GET /admin/metrics", apiConfig.HandlerShowMetrics)
	mux.HandleFunc("POST /admin/reset", apiConfig.HandlerResetUsers)
	mux.HandleFunc("POST /admin/resetmetrics", apiConfig.HandlerResetMetrics)

	server.ListenAndServe()
}
