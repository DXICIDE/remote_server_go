package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/DXICIDE/remote_server_go/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Password string `json:"password"`
		Email    string `json:"email"`
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

	userMap := User{}
	userMap.ID = user.ID
	userMap.Email = user.Email
	userMap.CreatedAt = user.CreatedAt
	userMap.UpdatedAt = user.UpdatedAt
	respondWithJSON(w, 200, userMap)
}
