package auth_test

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"simpleblog/auth"
	"testing"
)

func TestGeneratePassword(t *testing.T) {
	t.Run("successfully generating password", func(t *testing.T) {
		pass, err := auth.GeneratePassword(4, dumbReader{})
		if err != nil {
			t.Fatal("unexpected error:", err)
		}
		if pass != "0000" {
			t.Fatal("wrong result: ", pass)
		}
	})

	t.Run("error on reader", func(t *testing.T) {
		_, err := auth.GeneratePassword(4, dumbReader{err: errors.New("something happened")})
		if err == nil {
			t.Fatal("expecting error but got nil")
		}
	})

	t.Run("using real rand.Reader", func(t *testing.T) {
		pass, err := auth.GeneratePassword(8, rand.Reader)
		if err != nil {
			t.Fatal("unexpected error:", err)
		}
		if len(pass) != 8 {
			t.Fatal("wrong length: ", len(pass), pass)
		}
	})
}

func TestIsPasswordMatched(t *testing.T) {
	h := sha256.New()
	for _, tt := range []struct {
		name           string
		password       string
		hashedPassword string
	}{
		{
			name:           "admin",
			password:       "admin",
			hashedPassword: "8c6976e5b5410415bde908bd4dee15dfb167a9c873fc4bb8a81f6f2ab448a918",
		},
		{
			name:           "admin123",
			password:       "admin123",
			hashedPassword: "240be518fabd2724ddb6f04eeb1da5967448d7e831c08c8fa822809f74c720a9",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if !auth.PasswordMatched(tt.password, tt.hashedPassword, h) {
				t.Error("wrong result for:", tt.password)
			}
		})
	}
}

type dumbReader struct {
	err error
}

func (r dumbReader) Read(b []byte) (int, error) {
	if r.err != nil {
		return 0, r.err
	}
	return len(b), nil
}
