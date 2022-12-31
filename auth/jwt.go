package auth

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"simpleblog/httphandler"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func ES256TokenGenerator(privateKey *ecdsa.PrivateKey) func(c httphandler.JWTClaims) (string, error) {
	return func(c httphandler.JWTClaims) (string, error) {
		return jwt.NewWithClaims(jwt.SigningMethodES256, c).SignedString(privateKey)
	}
}

func ES256TokenParser(publicKey *ecdsa.PublicKey) func(token string) (httphandler.JWTClaims, error) {
	return func(strToken string) (httphandler.JWTClaims, error) {
		token, err := jwt.ParseWithClaims(strToken, &httphandler.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
				return nil, errors.New("unknown signing method")
			}
			return publicKey, nil
		})
		if err != nil {
			return httphandler.JWTClaims{}, fmt.Errorf("failed to parse jwt token: %w", err)
		}

		claims, ok := token.Claims.(*httphandler.JWTClaims)
		if !ok {
			return httphandler.JWTClaims{}, errors.New("the claims is not convertable to JWTClaims")
		}

		return *claims, nil

	}
}

func LoadX509ECDSAKey(r io.Reader) (*ecdsa.PrivateKey, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	blck, _ := pem.Decode(data)
	return x509.ParseECPrivateKey(blck.Bytes)
}

func generateDefaultToken(username, role string, tokenGenerator func(httphandler.JWTClaims) (string, error)) (string, error) {
	issuedAt := time.Now().UTC()
	expiredAt := issuedAt.Add(3 * time.Minute)

	token, err := tokenGenerator(httphandler.JWTClaims{
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  issuedAt.Unix(),
			ExpiresAt: expiredAt.Unix(),
		},
	})
	return token, err
}
