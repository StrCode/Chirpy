package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/StrCode/Chirpy/internal/auth"
	"github.com/StrCode/Chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, req *http.Request) {
	type requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	defer req.Body.Close()

	data, err := io.ReadAll(req.Body)
	if err != nil {
		return
	}

	var requestVals requestBody

	if err := json.Unmarshal(data, &requestVals); err != nil {
		respondWithError(w, 500, "unable to marshal the data", err)
	}

	hashedPwd, err := auth.HashPassword(requestVals.Password)
	if err != nil {
		respondWithError(w, 500, "could not hash", err)
	}

	newUser, err := cfg.dbQueries.CreateUser(context.Background(), database.CreateUserParams{
		Email:          requestVals.Email,
		HashedPassword: hashedPwd,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	createdUser := User{
		newUser.ID,
		newUser.CreatedAt,
		newUser.UpdatedAt,
		newUser.Email,
	}

	respondWithJSON(w, 201, createdUser)
}
