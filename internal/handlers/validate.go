package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/git-cst/bootdev_chirpy/internal/httputil"
)

type ValidationData struct {
	CleanedBody string `json:"cleaned_body"`
}

func HandlerValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type successResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		httputil.ErrorHandleValidation(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	if len(params.Body) > 140 {
		httputil.ErrorHandleValidation(w, "Chirp is too long", http.StatusBadRequest)
		return
	}

	respBody := ValidationData{
		CleanedBody: filterProfanity(params.Body),
	}

	httputil.RespondWithJSON(w, respBody, http.StatusOK)
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
