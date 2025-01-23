package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/StrCode/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Password string `json:"password"`
		Email    string `json:"email"`
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

	user := User{
		ID:        foundUser.ID,
		CreatedAt: foundUser.CreatedAt,
		UpdatedAt: foundUser.UpdatedAt,
		Email:     foundUser.Email,
	}

	respondWithJSON(w, 200, user)
}
