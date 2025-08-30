package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		log.Printf("Error fetching Chirps : %s", err)
		w.WriteHeader(500)
		return
	}
	validChirp := make([]Chirps, len(chirps))
	for i := 0; i < len(chirps); i++ {
		validChirp[i].Body = chirps[i].Body
		validChirp[i].CreatedAt = chirps[i].CreatedAt
		validChirp[i].UpdatedAt = chirps[i].UpdatedAt
		validChirp[i].ID = chirps[i].ID
		validChirp[i].User_id = chirps[i].UserID
	}

	respondWithJSON(w, 200, validChirp)
}
