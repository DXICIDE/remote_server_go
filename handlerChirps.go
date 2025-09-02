package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/DXICIDE/remote_server_go/internal/auth"
	"github.com/DXICIDE/remote_server_go/internal/database"
	"github.com/google/uuid"
)

type Chirps struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	User_id   uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting Token: %s", err)
		w.WriteHeader(500)
		return
	}

	id, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		log.Printf("Error Validating: %s", err)
		w.WriteHeader(401)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	bodysplit := strings.Split(params.Body, " ")

	for i := 0; i < len(bodysplit); i++ {
		body := strings.ToLower(bodysplit[i])
		if body == "kerfuffle" || body == "sharbert" || body == "fornax" {
			bodysplit[i] = "****"
		}
	}
	body := strings.Join(bodysplit, " ")
	cleanBody := database.CreateChirpParams{}
	cleanBody.Body = body
	cleanBody.UserID = id

	chirp, err := cfg.db.CreateChirp(r.Context(), cleanBody)
	if err != nil {
		log.Printf("Error creating Chirp: %s", err)
		w.WriteHeader(500)
		return
	}

	chirpJson := Chirps{}
	chirpJson.ID = chirp.ID
	chirpJson.Body = chirp.Body
	chirpJson.CreatedAt = chirp.CreatedAt
	chirpJson.UpdatedAt = chirp.UpdatedAt
	chirpJson.User_id = chirp.UserID
	respondWithJSON(w, 201, chirpJson)
}
