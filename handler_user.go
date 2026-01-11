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
func (t *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	// Body values
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body) // decoder for request body

	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	// create user
	user, err := t.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})

	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldn't create a user: %v", err))
		return
	}

	// passed created user
	responseWithJSON(w, 201, databaseUserToUser(user))
}

// GetUser handler
func (t *apiConfig) handlerGetUserByAPIKey(w http.ResponseWriter, r *http.Request, user database.User) {
	responseWithJSON(w, 200, databaseUserToUser(user))
}

// GetPostsForUser handler
func (t *apiConfig) handlerGetPostsForUser(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := t.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldn't get posts: %v", err))
		return
	}

	responseWithJSON(w, 200, databasePostsToPosts(posts))
}
