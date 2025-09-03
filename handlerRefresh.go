package main

import (
	"log"
	"net/http"
	"time"

	"github.com/DXICIDE/remote_server_go/internal/auth"
)

type Response struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	tokenID, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting Token from header: %s", err)
		w.WriteHeader(401)
		return
	}
	refreshToken, err := cfg.db.GetRefreshToken(r.Context(), tokenID)
	if err != nil {
		log.Printf("Error getting Refresh Token from Database: %s", err)
		w.WriteHeader(401)
		return
	}
	expired := refreshToken.ExpiresAt.Compare(time.Now())

	if expired < 1 {
		log.Printf("Token has expired")
		w.WriteHeader(401)
		return
	}

	if refreshToken.RevokedAt.Valid {
		log.Printf("Token has been revoked")
		w.WriteHeader(401)
		return
	}

	token, err := cfg.TokenForLogin(refreshToken.UserID)
	if err != nil {
		log.Printf("Couldn't make the token: %s", err)
		w.WriteHeader(401)
		return
	}

	res := Response{
		Token: token,
	}
	respondWithJSON(w, 200, res)
}
