package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	PathValue := r.PathValue("chirpID")
	id, err := uuid.Parse(PathValue)

	if err != nil {
		log.Printf("Error getting ID : %s", err)
		respondWithError(w, http.StatusBadRequest, "invalid chirp id")
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not found")
		return
	}
	chirpJson := Chirps{}
	chirpJson.ID = chirp.ID
	chirpJson.Body = chirp.Body
	chirpJson.CreatedAt = chirp.CreatedAt
	chirpJson.UpdatedAt = chirp.UpdatedAt
	chirpJson.User_id = chirp.UserID
	respondWithJSON(w, http.StatusOK, chirpJson)
}
