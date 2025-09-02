package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/DXICIDE/remote_server_go/internal/auth"
	"github.com/google/uuid"
)

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := body{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	user, err := cfg.db.GetUser(r.Context(), params.Email)
	if err != nil {
		log.Printf("Error fetching user or user doesnt exist: %s", err)
		w.WriteHeader(500)
		return
	}

	log.Printf("Password from request: %s", params.Password)
	log.Printf("Hashed password from DB: %s", user.HashedPassword)
	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		log.Printf("Wrong password: %s", err)
		w.WriteHeader(401)
		return
	}
	if params.ExpiresInSeconds == 0 || params.ExpiresInSeconds > 3600 {
		params.ExpiresInSeconds = 3600
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Duration(time.Duration(params.ExpiresInSeconds)*time.Second))
	if err != nil {
		log.Printf("Couldnt make the token: %s", err)
		w.WriteHeader(401)
		return
	}

	userMap := UserResponse{}
	userMap.ID = user.ID
	userMap.Email = user.Email
	userMap.CreatedAt = user.CreatedAt
	userMap.UpdatedAt = user.UpdatedAt
	userMap.Token = token
	respondWithJSON(w, 200, userMap)
}
