package auth

import (
	"encoding/json"
	"hash"
	"net/http"
	"simpleblog/httphandler"
	"time"
)

func HandleRegister(s UserStorage, hash hash.Hash, generatePassword func() (string, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := struct {
			Username string
			Role     string
		}{}

		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httphandler.HandleError(err, "bad request body", http.StatusBadRequest)(w, r)
			return
		}

		password, err := generatePassword()
		if err != nil {
			httphandler.HandleInternalServerError(err)(w, r)
			return
		}

		timestamp := time.Now().UTC().Truncate(time.Millisecond)

		u := StoredUser{
			Username:  req.Username,
			Role:      req.Role,
			Password:  GetPasswordHash(password, hash),
			Timestamp: timestamp,
		}

		if !u.isValid() {
			httphandler.HandleBadRequest()(w, r)
			return
		}

		if err := s.AddUser(u); err != nil {
			httphandler.HandleInternalServerError(err)(w, r)
			return
		}

		_ = json.NewEncoder(w).Encode(map[string]any{
			"username":  req.Username,
			"role":      req.Role,
			"password":  password,
			"timestamp": timestamp,
		})
	}
}
