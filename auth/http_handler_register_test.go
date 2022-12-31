package auth_test

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"simpleblog/auth"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHandleRegister(t *testing.T) {
	passwordGenerator := func() (string, error) {
		return "dumbpassword", nil
	}

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"username":"mamat", "role":"admin"}`))
		rec := httptest.NewRecorder()

		auth.HandleRegister(
			&stubUserStorage{
				AddUserFunc: func(auth.StoredUser) error {
					return nil
				},
			},
			sha256.New(),
			passwordGenerator,
		).ServeHTTP(rec, req)

		if statusCode := rec.Result().StatusCode; statusCode != http.StatusOK {
			t.Fatal("expecting 200 OK but got:", statusCode)
		}

		type Respond struct {
			Username  string
			Phone     string
			Password  string
			Role      string
			Timestamp time.Time
		}

		var res, expected Respond

		if err := json.NewDecoder(rec.Result().Body).Decode(&res); err != nil {
			t.Fatal("malformed json returned from the server:", err)
		}

		expected = Respond{
			Username: "mamat",
			Phone:    "6280989999",
			Password: "dumbpassword",
			Role:     "admin",
		}

		assert.Equal(t, expected.Username, res.Username)
		assert.Equal(t, expected.Password, res.Password)
		assert.Equal(t, expected.Role, res.Role)
	})

	for _, tt := range []struct {
		Name               string
		Request            *http.Request
		StubStorage        *stubUserStorage
		PasswordGenerator  func() (string, error)
		ExpectedStatusCode int
	}{
		{
			Name:    "storage error",
			Request: httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"username":"mamat", "role":"admin"}`)),
			StubStorage: &stubUserStorage{
				AddUserFunc: func(auth.StoredUser) error {
					return errors.New("something happened in the storage")
				},
			},
			PasswordGenerator:  passwordGenerator,
			ExpectedStatusCode: http.StatusInternalServerError,
		},
		{
			Name:               "username less than 3 characters",
			Request:            httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"username":"ma", "role":"admin"}`)),
			StubStorage:        nil,
			PasswordGenerator:  passwordGenerator,
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name:    "password generation error",
			Request: httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"username":"mamat", "role":"admin"}`)),
			StubStorage: &stubUserStorage{
				AddUserFunc: func(auth.StoredUser) error {
					return nil
				},
			},
			PasswordGenerator: func() (string, error) {
				return "", errors.New("something happened in the password generator")
			},
			ExpectedStatusCode: http.StatusInternalServerError,
		},
		{
			Name:               "malformed json request",
			Request:            httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{`)),
			ExpectedStatusCode: http.StatusBadRequest,
		},
	} {
		t.Run(tt.Name, func(t *testing.T) {
			req := tt.Request
			rec := httptest.NewRecorder()

			auth.HandleRegister(tt.StubStorage, sha256.New(), tt.PasswordGenerator).ServeHTTP(rec, req)

			if statusCode := rec.Result().StatusCode; statusCode != tt.ExpectedStatusCode {
				t.Fatal("wrong status code:", statusCode)
			}
		})
	}
}

type stubUserStorage struct {
	GetUserByUsernameFunc func(phone string) (auth.StoredUser, error)
	AddUserFunc           func(auth.StoredUser) error
}

func (s *stubUserStorage) GetUserByUsername(username string) (auth.StoredUser, error) {
	return s.GetUserByUsernameFunc(username)
}

func (s *stubUserStorage) AddUser(u auth.StoredUser) error {
	return s.AddUserFunc(u)
}
