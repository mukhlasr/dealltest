package auth_test

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"simpleblog/auth"
	"simpleblog/httphandler"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHandleLogin(t *testing.T) {
	defaultReq := httptest.NewRequest(http.MethodGet, "/auth-token", nil)
	defaultReq.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("1234"+":"+"admin123")))

	defaultStorageStub := &stubUserStorage{
		GetUserByUsernameFunc: func(phoneNum string) (auth.StoredUser, error) {
			return auth.StoredUser{
				Username:  "foo",
				Password:  "240be518fabd2724ddb6f04eeb1da5967448d7e831c08c8fa822809f74c720a9",
				Role:      "admin",
				Timestamp: time.Now().UTC(),
			}, nil
		},
	}

	t.Run("success", func(t *testing.T) {
		req := defaultReq
		rec := httptest.NewRecorder()

		auth.HandleLogin(
			defaultStorageStub,
			sha256.New(),
			func(j httphandler.JWTClaims) (string, error) {
				return "supertoken", nil
			},
		).ServeHTTP(rec, req)

		if statusCode := rec.Result().StatusCode; statusCode != http.StatusOK {
			t.Fatal("expecting 200 OK but got:", statusCode)
		}

		type Respond struct {
			Token string
		}

		var res, expected Respond

		if err := json.NewDecoder(rec.Result().Body).Decode(&res); err != nil {
			t.Fatal("malformed json returned from the server:", err)
		}

		expected = Respond{
			Token: "supertoken",
		}

		assert.Equal(t, expected, res)
	})

	for _, tt := range []struct {
		Name               string
		GetRequest         func() *http.Request
		StubStorage        *stubUserStorage
		TokenGenerator     func(httphandler.JWTClaims) (string, error)
		ExpectedStatusCode int
	}{
		{
			Name: "storage error",
			GetRequest: func() *http.Request {
				return defaultReq
			},
			StubStorage: &stubUserStorage{
				GetUserByUsernameFunc: func(phone string) (auth.StoredUser, error) {
					return auth.StoredUser{}, errors.New("something happened in the storage")
				},
			},
			TokenGenerator:     nil,
			ExpectedStatusCode: http.StatusInternalServerError,
		},
		{
			Name: "token generation error",
			GetRequest: func() *http.Request {
				return defaultReq
			},
			StubStorage: defaultStorageStub,
			TokenGenerator: func(j httphandler.JWTClaims) (string, error) {
				return "", errors.New("unexpected error")
			},
			ExpectedStatusCode: http.StatusInternalServerError,
		},
		{
			Name: "no auth",
			GetRequest: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/auth-token", nil)
			},
			StubStorage:        defaultStorageStub,
			TokenGenerator:     nil,
			ExpectedStatusCode: http.StatusUnauthorized,
		},
		{
			Name: "wrong password",
			GetRequest: func() *http.Request {
				r := &http.Request{}
				*r = *defaultReq
				log.Println("ini", r)
				r.SetBasicAuth("foo", "wrongpassword")
				return r
			},
			StubStorage:        defaultStorageStub,
			TokenGenerator:     nil,
			ExpectedStatusCode: http.StatusUnauthorized,
		},
	} {
		t.Run(tt.Name, func(t *testing.T) {
			req := tt.GetRequest()
			rec := httptest.NewRecorder()

			auth.HandleLogin(tt.StubStorage, sha256.New(), tt.TokenGenerator).ServeHTTP(rec, req)

			if statusCode := rec.Result().StatusCode; statusCode != tt.ExpectedStatusCode {
				t.Fatal("wrong status code:", statusCode)
			}
		})
	}
}
