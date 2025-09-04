package main

import (
	"log"
	"net/http"
	"sort"

	"github.com/DXICIDE/remote_server_go/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("author_id")
	sortQ := r.URL.Query().Get("sort")

	var chirps []database.Chirp
	var err error

	if id == "" {
		chirps, err = cfg.db.GetChirps(r.Context())
		if err != nil {
			log.Printf("Error fetching Chirps : %s", err)
			w.WriteHeader(500)
			return
		}
	} else {
		id, err := uuid.Parse(id)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid chirp id")
			return
		}

		chirps, err = cfg.db.GetManyChirps(r.Context(), id)
		if err != nil {
			log.Printf("Error fetching Chirps : %s", err)
			w.WriteHeader(500)
			return
		}
	}

	validChirp := make([]Chirps, len(chirps))
	for i := 0; i < len(chirps); i++ {
		validChirp[i].Body = chirps[i].Body
		validChirp[i].CreatedAt = chirps[i].CreatedAt
		validChirp[i].UpdatedAt = chirps[i].UpdatedAt
		validChirp[i].ID = chirps[i].ID
		validChirp[i].User_id = chirps[i].UserID
	}

	if sortQ == "desc" {
		sort.Slice(validChirp, func(i, j int) bool {
			return validChirp[i].CreatedAt.After(validChirp[j].CreatedAt)
		})
	}

	respondWithJSON(w, 200, validChirp)
}
