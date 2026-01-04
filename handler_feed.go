package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Kunal0007/RSS_Aggregator/internal/database"
	"github.com/google/uuid"
)

// Create User handler
func (t *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {

	// Body values
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body) // decoder for request body

	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	// create Feed
	feed, err := t.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	})

	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldn't create a feed: %v", err))
		return
	}

	// passed created user
	responseWithJSON(w, 201, databaseFeedToFeed(feed))
}
