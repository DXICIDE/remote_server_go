package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	header := headers.Values("Authorization")
	if len(header) == 0 {
		return "", fmt.Errorf("no authorization header")
	}
	const prefix = "ApiKey "
	apikey := strings.Replace(header[0], prefix, "", -1)
	apikey = strings.TrimSpace(apikey)
	if apikey == "" {
		return "", fmt.Errorf("empty apikey")
	}
	return apikey, nil
}
