package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/DXICIDE/remote_server_go/internal/auth"
	"github.com/DXICIDE/remote_server_go/internal/database"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
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

	updateToken := database.UpdateTokenParams{
		RevokedAt: sql.NullTime{Time: time.Now(), Valid: true},
		Token:     refreshToken.Token,
		UpdatedAt: time.Now(),
	}

	err = cfg.db.UpdateToken(r.Context(), updateToken)
	if refreshToken.RevokedAt.Valid {
		log.Printf("Token couldnt be updated: %v", err)
		w.WriteHeader(401)
		return
	}

	respondWithJSON(w, 204, "")
}
