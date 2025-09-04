package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/DXICIDE/remote_server_go/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerWebhook(w http.ResponseWriter, r *http.Request) {
	type data struct {
		UserID string `json:"user_id"`
	}

	type parameters struct {
		Event string `json:"event"`
		Data  data   `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		log.Printf("Error getting Token: %s", err)
		w.WriteHeader(401)
		return
	}

	if apiKey != cfg.polka_key {
		respondWithError(w, 401, "Wrong Api key")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, "")
		return
	}

	id, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chirp id")
		return
	}

	err = cfg.db.UpgradeChirpyRed(r.Context(), id)
	if err != nil {
		respondWithError(w, 404, "Couldnt upgrade user")
		return
	}

	respondWithJSON(w, 204, "")
}
