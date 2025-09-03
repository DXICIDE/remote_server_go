package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/DXICIDE/remote_server_go/internal/auth"
	"github.com/DXICIDE/remote_server_go/internal/database"
	"github.com/google/uuid"
)

type UserResponse struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

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

	token, err := cfg.TokenForLogin(user.ID)
	if err != nil {
		log.Printf("Couldnt make the token: %s", err)
		w.WriteHeader(401)
		return
	}

	refreshToken, err := cfg.TokenForRefresh(user.ID, r)
	if err != nil {
		log.Printf("Couldnt make the refresh token: %s", err)
		w.WriteHeader(401)
		return
	}

	userMap := UserResponse{}
	userMap.ID = user.ID
	userMap.Email = user.Email
	userMap.CreatedAt = user.CreatedAt
	userMap.UpdatedAt = user.UpdatedAt
	userMap.Token = token
	userMap.RefreshToken = refreshToken
	userMap.IsChirpyRed = user.IsChirpyRed
	respondWithJSON(w, 200, userMap)
}

func (cfg *apiConfig) TokenForRefresh(id uuid.UUID, r *http.Request) (string, error) {
	refreshTokenID, err := auth.MakeRefreshToken()
	if err != nil {
		return "", err
	}

	createToken := database.CreateTokensParams{
		UserID: id,
		Token:  refreshTokenID,
	}

	refreshToken, err := cfg.db.CreateTokens(r.Context(), createToken)
	return refreshToken.Token, err
}

func (cfg *apiConfig) TokenForLogin(id uuid.UUID) (string, error) {
	seconds := 3600
	token, err := auth.MakeJWT(id, cfg.secret, time.Duration(time.Duration(seconds)*time.Second))
	if err != nil {
		return "", err
	}
	return token, err
}
