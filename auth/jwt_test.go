package auth_test

import (
	"simpleblog/auth"
	"simpleblog/httphandler"
	"strings"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func TestLoadX509ECKey(t *testing.T) {
	f := httphandler.DefaultJWTPrivKey
	privKey, err := auth.LoadX509ECDSAKey(strings.NewReader(f))
	if err != nil {
		t.Fatal("failed to load the key", err)
	}
	if privKey == nil {
		t.Fatal("expecting not null key")
	}
}

func TestGenerateToken(t *testing.T) {
	f := httphandler.DefaultJWTPrivKey
	privKey, err := auth.LoadX509ECDSAKey(strings.NewReader(f))
	if err != nil {
		t.Fatal("failed to load the key", err)
	}

	expectedClaims := httphandler.JWTClaims{
		Username: "mamat",
		Role:     "admin",
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Unix() + 10,
		},
	}

	tokenGenerator := auth.ES256TokenGenerator(privKey)
	token, err := tokenGenerator(expectedClaims)
	if err != nil {
		t.Fatal("failed to generate token", err)
	}

	tokenParser := auth.ES256TokenParser(&privKey.PublicKey)
	claims, err := tokenParser(token)
	if err != nil {
		t.Fatal("error on parsing the output token", err)
	}

	if claims != expectedClaims {
		t.Fatal("wrong claims", claims)
	}
}
