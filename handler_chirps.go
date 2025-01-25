package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/StrCode/Chirpy/internal/auth"
	"github.com/StrCode/Chirpy/internal/database"
	"github.com/google/uuid"
)

type chirp struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	User_Id   uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	data, err := cfg.dbQueries.GetAllChirps(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
		return
	}

	var chirps []chirp
	for _, ch := range data {
		chirps = append(chirps, chirp{
			ch.ID,
			ch.CreatedAt,
			ch.UpdatedAt,
			ch.Body,
			ch.UserID,
		})
	}
	respondWithJSON(w, 200, chirps)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("chirpId")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "No ChirpId provided", nil)
		return
	}

	chirpId, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "chirpId is not valid", err)
		return
	}

	selectedChirp, err := cfg.dbQueries.GetChirp(context.Background(), chirpId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, chirp{
		Id:        selectedChirp.ID,
		CreatedAt: selectedChirp.CreatedAt,
		UpdatedAt: selectedChirp.UpdatedAt,
		Body:      selectedChirp.Body,
		User_Id:   selectedChirp.UserID,
	})
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	isOk := decoder.Decode(&params)
	if isOk != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", isOk)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(params.Body, badWords)

	newChirp, err := cfg.dbQueries.CreateChirp(context.Background(), database.CreateChirpParams{
		Body:   cleaned,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	createdChirp := chirp{
		Id:        newChirp.ID,
		CreatedAt: newChirp.CreatedAt,
		UpdatedAt: newChirp.UpdatedAt,
		Body:      newChirp.Body,
		User_Id:   newChirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, createdChirp)
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}
