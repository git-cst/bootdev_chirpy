package handlers

import (
	"net/http"

	"github.com/git-cst/bootdev_chirpy/internal/configs"
)

func MakeHandlerWithConfig(cfg *configs.ApiConfig, handler func(http.ResponseWriter, *http.Request, *configs.ApiConfig)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, cfg)
	}
}
