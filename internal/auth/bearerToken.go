package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	header := headers.Values("Authorization")
	if len(header) == 0 {
		return "", fmt.Errorf("no authorization header")
	}
	const prefix = "Bearer "
	token := strings.Replace(header[0], prefix, "", -1)
	token = strings.TrimSpace(token)
	if token == "" {
		return "", fmt.Errorf("empty bearer token")
	}
	return token, nil
}
