package posts_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"simpleblog/posts"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetAllPosts(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()

		posts.HandleGetAllPosts(&stubPostAccessor{
			GetAllFunc: func(ctx context.Context) ([]posts.StoredPost, error) {
				return []posts.StoredPost{
					{
						ID:        1,
						Title:     "title 1",
						Content:   "content 1",
						Timestamp: time.Time{},
					},
					{
						ID:        2,
						Title:     "title 2",
						Content:   "content 2",
						Timestamp: time.Time{},
					},
				}, nil
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

		var res []Respond

		if err := json.NewDecoder(rec.Result().Body).Decode(&res); err != nil {
			t.Fatal("malformed json returned from the server:", err)
		}

		if len(res) != 2 {
			t.Fatal("wrong result length")
		}
	})
}

func TestGetPostByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/?id=123", nil)
		rec := httptest.NewRecorder()

		posts.HandleGetPostByID(&stubPostAccessor{
			GetByIDFunc: func(ctx context.Context, id int64) (posts.StoredPost, error) {
				return posts.StoredPost{
					ID:        id,
					Title:     "title 1",
					Content:   "content 1",
					Timestamp: time.Time{},
				}, nil
			},
		}, "123").ServeHTTP(rec, req)

		if statusCode := rec.Result().StatusCode; statusCode != http.StatusOK {
			t.Fatal("expecting 200 OK but got:", statusCode)
		}

		type Respond struct {
			ID        int64
			Title     string
			Content   string
			Timestamp time.Time
		}

		var resp []Respond

		if err := json.NewDecoder(rec.Result().Body).Decode(&resp); err != nil {
			t.Fatal("malformed json returned from the server:", err)
		}

		if len(resp) != 1 {
			t.Fatal("result length should be 1")
		}

		res := resp[0]

		assert.Equal(t, res.ID, int64(123))
		assert.Equal(t, res.Title, "title 1")
		assert.Equal(t, res.Content, "content 1")
		assert.Equal(t, res.Timestamp, time.Time{})
	})
}

type stubPostAccessor struct {
	GetAllFunc  func(ctx context.Context) ([]posts.StoredPost, error)
	GetByIDFunc func(ctx context.Context, id int64) (posts.StoredPost, error)
}

func (s *stubPostAccessor) GetByID(ctx context.Context, id int64) (posts.StoredPost, error) {
	return s.GetByIDFunc(ctx, id)
}

func (s *stubPostAccessor) GetAll(ctx context.Context) ([]posts.StoredPost, error) {
	return s.GetAllFunc(ctx)
}
