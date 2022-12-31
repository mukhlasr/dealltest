package auth

import "time"

const (
	refreshTokenTimeoutDuration = 24 * time.Hour
)

type RefreshTokenContent struct {
	Token    string
	Username string
	Role     string
}

type RefreshTokenStorage interface {
	StoreRefreshToken(token string, content RefreshTokenContent, exp time.Duration) error
	GetRefreshToken(token string) (RefreshTokenContent, error)
	RevokeRefreshToken(token string) error
}
