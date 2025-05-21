package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/git-cst/bootdev_chirpy/internal/configs"
	"github.com/git-cst/bootdev_chirpy/internal/httputil"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type UserCreationData struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func HandlerUsers(w http.ResponseWriter, r *http.Request, c *configs.ApiConfig) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		httputil.RespondWithError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	user, err := c.Db.CreateUser(ctx, params.Email)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique_violation code
				httputil.RespondWithError(w, "User already exists", http.StatusBadRequest)
				return
			}
			// all other errors
			httputil.RespondWithError(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
	}

	respBody := UserCreationData{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	httputil.RespondWithJSON(w, respBody, http.StatusCreated)
}
