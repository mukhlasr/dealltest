package posts_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"simpleblog/posts"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreatePost(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"title":"mamat pergi ke desa", "content":"konten khusus"}`))
		rec := httptest.NewRecorder()

		posts.HandleCreatePost(&stubPostCreator{
			CreateFunc: func(ctx context.Context, p posts.StoredPost) (id int64, err error) {
				return 1, nil
			},
		}).ServeHTTP(rec, req)

		if statusCode := rec.Result().StatusCode; statusCode != http.StatusOK {
			t.Fatal("expecting 200 OK but got:", statusCode)
		}

		type Respond struct {
			ID        int64
			Title     string
			Content   string
			Timestamp time.Time
		}

		var res, expected Respond

		if err := json.NewDecoder(rec.Result().Body).Decode(&res); err != nil {
			t.Fatal("malformed json returned from the server:", err)
		}

		expected = Respond{
			ID:      1,
			Title:   "mamat pergi ke desa",
			Content: "konten khusus",
		}

		assert.Equal(t, expected.ID, res.ID)
		assert.Equal(t, expected.Title, res.Title)
		assert.Equal(t, expected.Content, res.Content)
	})
}

type stubPostCreator struct {
	CreateFunc func(ctx context.Context, p posts.StoredPost) (id int64, err error)
}

func (s *stubPostCreator) Create(ctx context.Context, p posts.StoredPost) (int64, error) {
	return s.CreateFunc(ctx, p)
}
