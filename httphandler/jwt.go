package httphandler

import (
	_ "embed"

	"github.com/dgrijalva/jwt-go"
)

//go:embed testkeys/test_priv.ec.key
var DefaultJWTPrivKey string

type JWTClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}
