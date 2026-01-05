package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Kunal0007/RSS_Aggregator/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

// Create FeedFollow handler
func (t *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {

	// Body values
	type parameters struct {
		FeedId uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body) // decoder for request body

	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	// create Feed
	feedFollow, err := t.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedId,
	})

	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldn't create a feedfollow: %v", err))
		return
	}

	responseWithJSON(w, 201, databaseFeedFollowToFeedFollow(feedFollow))
}

// Get FeedFollows handler
func (t *apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {

	// get feed_follows
	feedFollows, err := t.DB.GetFeedFollows(r.Context(), user.ID)

	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldn't get feedfollows: %v", err))
		return
	}

	responseWithJSON(w, 200, databaseFeedFollowsToFeedFollows(feedFollows))
}

// Delete FeedFollow handler
func (t *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIDStr := chi.URLParam(r, "feedFollowID")

	feedFollowID, err := uuid.Parse(feedFollowIDStr)

	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	err = t.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})

	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldn't delete the feedFollow: %v", err))
		return
	}

	responseWithJSON(w, 200, struct{}{})
}
