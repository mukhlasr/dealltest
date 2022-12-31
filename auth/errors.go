package auth

import "errors"

var (
	ErrUsernameNotFound = errors.New("username not found")
	ErrWrongAuth        = errors.New("username/password mismatched")
)
