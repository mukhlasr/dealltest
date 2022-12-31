package auth

import "time"

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type StoredUser struct {
	Username  string
	Password  string
	Role      string
	Timestamp time.Time
}

func (s StoredUser) isValid() bool {
	if len(s.Username) < 3 {
		return false
	}

	if len(s.Password) < 3 {
		return false
	}

	if s.Role != RoleAdmin && s.Role != RoleUser {
		return false
	}

	return true
}

type UserStorage interface {
	AddUser(u StoredUser) error
	GetUserByUsername(phone string) (StoredUser, error)
}
