package auth

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"simpleblog/httphandler"
	"time"
)

func HandleRefreshToken(us UserStorage, ts RefreshTokenStorage, tokenGenerator func(httphandler.JWTClaims) (string, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		refreshToken, err := r.Cookie("refresh_token")
		if err == http.ErrNoCookie {
			httphandler.HandleError(err, "no refresh token", http.StatusUnauthorized)(w, r)
			return
		}
		if err != nil {
			httphandler.HandleInternalServerError(err)(w, r)
			return
		}

		content, err := ts.GetRefreshToken(refreshToken.Raw)
		if err != nil {
			httphandler.HandleInternalServerError(err)(w, r)
			return
		}

		u, err := us.GetUserByUsername(content.Username)
		if err != nil {
			httphandler.HandleInternalServerError(err)(w, r)
			return
		}

		if err := ts.RevokeRefreshToken(refreshToken.Raw); err != nil {
			httphandler.HandleInternalServerError(fmt.Errorf("failed to revoke refresh token: %w", err))(w, r)
			return
		}

		newRefreshToken, err := generateRefreshToken(ts, u)
		if err != nil {
			httphandler.HandleInternalServerError(err)(w, r)
			return
		}

		token, err := generateDefaultToken(u.Username, u.Role, tokenGenerator)
		if err != nil {
			httphandler.HandleInternalServerError(err)(w, r)
			return
		}

		setRefreshTokenCookie(w, newRefreshToken)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"token": token,
		})
	}
}

func generateRefreshToken(ts RefreshTokenStorage, u StoredUser) (string, error) {
	randBytes := make([]byte, 32)
	_, err := rand.Read(randBytes)
	if err != nil {
		return "", err
	}
	token := fmt.Sprintf("%x", randBytes)
	return token, ts.StoreRefreshToken(token, RefreshTokenContent{
		Token:    token,
		Username: u.Username,
		Role:     u.Role,
	}, refreshTokenTimeoutDuration)
}

func setRefreshTokenCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		Expires:  time.Now().UTC().Add(refreshTokenTimeoutDuration),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	})
}
