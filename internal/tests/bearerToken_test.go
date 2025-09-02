package auth

import (
	"net/http/httptest"
	"testing"

	"github.com/DXICIDE/remote_server_go/internal/auth"
)

func TestBearerToken(t *testing.T) {
	w := httptest.NewRecorder()
	w.Header().Set("Authorization", "Bearer ${jwtTokenSaul}")
	res, err := auth.GetBearerToken(w.Header())
	if err != nil {
		t.Fatalf("GetBearerToken returned error: %v", err)
	}

	if res != "${jwtTokenSaul}" {
		t.Fatalf("wrong token: got %q, want %q", res, "${jwtTokenSaul}")
	}
}

func TestBearerTokenNoHeader(t *testing.T) {
	w := httptest.NewRecorder()
	_, err := auth.GetBearerToken(w.Header())
	if err == nil {
		t.Errorf("There should've been an error")
	}
}
