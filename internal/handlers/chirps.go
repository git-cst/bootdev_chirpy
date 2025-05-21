package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/git-cst/bootdev_chirpy/internal/configs"
	"github.com/git-cst/bootdev_chirpy/internal/database"
	"github.com/git-cst/bootdev_chirpy/internal/httputil"
	"github.com/google/uuid"
)

type CreatedPost struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func HandlerPostChirp(w http.ResponseWriter, r *http.Request, c *configs.ApiConfig) {
	type PostCreationData struct {
		Body   string `json:"body"`
		UserID string `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	postData := PostCreationData{}
	err := decoder.Decode(&postData)
	if err != nil {
		httputil.RespondWithError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	postLength := len(postData.Body)
	if postLength == 0 {
		httputil.RespondWithError(w, "Chirp contains no data", http.StatusUnprocessableEntity)
		return
	} else if postLength > 140 {
		httputil.RespondWithError(w, "Chirp is too long", http.StatusBadRequest)
		return
	}

	postUser, err := uuid.Parse(postData.UserID)
	if err != nil {
		httputil.RespondWithError(w, "Invalid UUID string passed", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	_, err = c.Db.GetUser(ctx, postUser)
	if err != nil {
		// no data returned
		if errors.Is(err, sql.ErrNoRows) {
			httputil.RespondWithError(w, "User does not exist", http.StatusUnprocessableEntity)
			return
		}
		// all other errors
		httputil.RespondWithError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	params := database.CreatePostParams{
		Body:   filterProfanity(postData.Body),
		UserID: postUser,
	}

	createdPost, err := c.Db.CreatePost(ctx, params)
	if err != nil {
		errMsg := fmt.Sprintf("Error creating post. Error: %s.", err)
		httputil.RespondWithError(w, errMsg, http.StatusInternalServerError)
		return
	}

	respBody := CreatedPost{
		ID:        createdPost.ID,
		CreatedAt: createdPost.CreatedAt,
		UpdatedAt: createdPost.UpdatedAt,
		Body:      createdPost.Body,
		UserID:    createdPost.UserID,
	}

	httputil.RespondWithJSON(w, respBody, http.StatusCreated)
}

func filterProfanity(rBody string) string {
	var profanityList = [3]string{"kerfuffle", "sharbert", "fornax"}
	var filterArray []string

	for _, word := range strings.Split(rBody, " ") {
		wordFound := false

		for _, profanity := range profanityList {
			if strings.ToLower(word) == profanity {
				filterArray = append(filterArray, "****")
				wordFound = true
				break
			}
		}

		if wordFound {
			continue
		} else {
			filterArray = append(filterArray, word)
		}
	}

	filteredMessage := strings.Join(filterArray, " ")

	return filteredMessage
}

func HandlerGetChirps(w http.ResponseWriter, r *http.Request, c *configs.ApiConfig) {
	ctx := context.Background()
	posts, err := c.Db.GetPosts(ctx)
	if err != nil {
		httputil.RespondWithError(w, "Failed to get chirps", http.StatusInternalServerError)
		return
	}

	var chirps []CreatedPost
	for _, post := range posts {
		chirps = append(chirps, CreatedPost{
			ID:        post.ID,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
			Body:      post.Body,
			UserID:    post.UserID,
		})
	}

	httputil.RespondWithJSON(w, chirps, http.StatusOK)
}

func HandlerGetChirp(w http.ResponseWriter, r *http.Request, c *configs.ApiConfig) {
	chirpString := r.PathValue("chirpID")
	chirpId, err := uuid.Parse(chirpString)
	if err != nil {
		httputil.RespondWithError(w, "Invalid UUID string passed", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	post, err := c.Db.GetPost(ctx, chirpId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httputil.RespondWithError(w, "No chirp with that ID", http.StatusBadRequest)
			return
		}
		httputil.RespondWithError(w, "Something went wrong", http.StatusInternalServerError)
	}

	respPost := CreatedPost{
		ID:        post.ID,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		Body:      post.Body,
		UserID:    post.UserID,
	}

	httputil.RespondWithJSON(w, respPost, http.StatusOK)
}
