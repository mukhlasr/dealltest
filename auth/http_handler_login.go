package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"net/http"
	"simpleblog/httphandler"
)

func HandleLogin(s UserStorage, h hash.Hash, tokenGenerator func(httphandler.JWTClaims) (string, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			httphandler.HandleUnauthorized()(w, r)
			return
		}

		user, err := login(s, h, username, password)
		if errors.Is(err, ErrUsernameNotFound) || errors.Is(err, ErrWrongAuth) {
			httphandler.HandleUnauthorized()(w, r)
			return
		}

		if err != nil {
			httphandler.HandleInternalServerError(err)(w, r)
			return
		}

		token, err := generateDefaultToken(user.Username, user.Role, tokenGenerator)
		if err != nil {
			httphandler.HandleInternalServerError(err)(w, r)
			return
		}

		_ = json.NewEncoder(w).Encode(map[string]any{
			"token": token,
		})
	}
}

func login(s UserStorage, h hash.Hash, username, password string) (StoredUser, error) {
	user, err := s.GetUserByUsername(username)
	if err != nil {
		return user, fmt.Errorf("failed to get username: %w", err)
	}

	if !PasswordMatched(password, user.Password, h) {
		return user, ErrWrongAuth
	}
	return user, nil
}
