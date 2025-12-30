package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Kunal0007/RSS_Aggregator/internal/database"
	"github.com/google/uuid"
)

func (t *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	// Body values
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body) // decoder for request body

	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Error parsing JSON:%v", err))
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
		responseWithError(w, 400, fmt.Sprintf("Couldn't create a user:%v", err))
		return
	}

	// passed created user
	responseWithJSON(w, 200, databaseUserToUser(user))
}
