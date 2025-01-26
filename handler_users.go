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
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
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
		newUser.IsChirpyRed,
	}

	respondWithJSON(w, 201, createdUser)
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	access_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "user not authorized", err)
	}

	userID, err := auth.ValidateJWT(access_token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	defer r.Body.Close()
	var requestVals requestBody

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestVals); err != nil {
		respondWithError(w, http.StatusBadRequest, "request malformed", err)
		return
	}

	hashedPwd, err := auth.HashPassword(requestVals.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "something went wrong", err)
		return
	}

	updatedUser, err := cfg.dbQueries.UpdateUser(context.Background(), database.UpdateUserParams{
		Email:          requestVals.Email,
		HashedPassword: hashedPwd,
		ID:             userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not update user", err)
		return
	}

	respondWithJSON(w, 200, User{
		updatedUser.ID,
		updatedUser.CreatedAt,
		updatedUser.UpdatedAt,
		updatedUser.Email,
		updatedUser.IsChirpyRed,
	})
}
