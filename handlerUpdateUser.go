package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/DXICIDE/remote_server_go/internal/auth"
	"github.com/DXICIDE/remote_server_go/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
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

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting Token: %s", err)
		w.WriteHeader(401)
		return
	}

	id, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		log.Printf("Error Validating: %s", err)
		w.WriteHeader(401)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "Unsuccesful Hashing")
		return
	}

	updateUser := database.UpdateUserParams{
		Email:          params.Email,
		ID:             id,
		HashedPassword: hashedPassword,
	}

	err = cfg.db.UpdateUser(r.Context(), updateUser)
	if err != nil {
		respondWithError(w, 500, "Unsuccesful User Update")
	}

	user, err := cfg.db.GetUserByID(r.Context(), id)
	if err != nil {
		log.Printf("Error fetching user or user doesnt exist: %s", err)
		w.WriteHeader(500)
		return
	}

	userBody := User{
		UpdatedAt: user.UpdatedAt,
		CreatedAt: user.CreatedAt,
		Email:     user.Email,
		ID:        user.ID,
	}

	respondWithJSON(w, 200, userBody)

}
