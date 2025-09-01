package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := &jwt.RegisteredClaims{
		Issuer:   "chirpy",
		IssuedAt: &jwt.NumericDate{Time: time.Now().UTC()},
		Subject:  userID.String(),
	}
	claims.ExpiresAt = &jwt.NumericDate{Time: claims.IssuedAt.Time.Add(expiresIn)}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tokenSecret))
}
