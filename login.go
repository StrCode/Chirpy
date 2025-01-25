package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/StrCode/Chirpy/internal/auth"
	"github.com/google/uuid"
)

type LoggedInUser struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Password           string `json:"password"`
		Email              string `json:"email"`
		Expires_in_Seconds int    `json:"expires_in_seconds"`
	}

	defer r.Body.Close()

	var requestVals requestBody
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestVals)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	foundUser, err := cfg.dbQueries.GetUserByEmail(context.Background(), requestVals.Email)
	if err != nil {
		respondWithError(w, 401, "incorrect email or password", nil)
		return
	}

	if err := auth.CheckPasswordHash(requestVals.Password, foundUser.HashedPassword); err != nil {
		respondWithError(w, 401, "incorrect email or password", nil)
		return
	}

	expiresIn, _ := time.ParseDuration("1h")
	if requestVals.Expires_in_Seconds != 0 {
		expiresIn = time.Duration(requestVals.Expires_in_Seconds) * time.Second
	}

	token, err := auth.MakeJWT(
		foundUser.ID,
		cfg.jwtSecret,
		expiresIn,
	)
	if err != nil {
		respondWithError(w, 401, "incorrect email or password", nil)
		return
	}

	user := LoggedInUser{
		ID:        foundUser.ID,
		CreatedAt: foundUser.CreatedAt,
		UpdatedAt: foundUser.UpdatedAt,
		Email:     foundUser.Email,
		Token:     token,
	}

	respondWithJSON(w, 200, user)
}
