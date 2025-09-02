package auth

import (
	"testing"
	"time"

	"github.com/DXICIDE/remote_server_go/internal/auth"
	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	id := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	secret := "test-secret"
	tok, err := auth.MakeJWT(id, secret, time.Minute)
	if err != nil {
		t.Errorf("Auth.MakeJMT did not work: %v", err)
	}
	gotID, err := auth.ValidateJWT(tok, secret)
	if err != nil {
		t.Errorf("Auth.ValidateJmt did not work: %v", err)
	}
	if gotID != id {
		t.Errorf("ID's are not identical %v != %v ", id, gotID)
	}
}

func TestJMTExpiredate(t *testing.T) {
	id := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	secret := "test-secret"
	tok, err := auth.MakeJWT(id, secret, time.Second)
	if err != nil {
		t.Errorf("Auth.MakeJMT did not work: %v", err)
	}
	time.Sleep(time.Second * 2)
	_, err = auth.ValidateJWT(tok, secret)
	if err == nil {
		t.Errorf("Token should've expired")
	}
}

func TestJMTWrongSecret(t *testing.T) {
	id := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	secret := "test-secret"
	secret1 := "secret"
	tok, err := auth.MakeJWT(id, secret, time.Minute)
	if err != nil {
		t.Errorf("Auth.MakeJMT did not work: %v", err)
	}
	_, err = auth.ValidateJWT(tok, secret1)
	if err == nil {
		t.Errorf("Token should be invalid, wrong secter")
	}
}
