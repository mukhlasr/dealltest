package posts_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"simpleblog/posts"
	"strings"
	"testing"
)

func TestUpdatePost(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"id": 123, "title":"mamat pergi ke desa", "content":"konten khusus"}`))
		rec := httptest.NewRecorder()

		posts.HandleUpdatePost(&stubPostMutator{
			UpdateFunc: func(ctx context.Context, p posts.StoredPost) error {
				return nil
			},
		}).ServeHTTP(rec, req)

		if statusCode := rec.Result().StatusCode; statusCode != http.StatusNoContent {
			t.Fatal("expecting 201 OK but got:", statusCode)
		}
	})
}

type stubPostMutator struct {
	UpdateFunc func(ctx context.Context, p posts.StoredPost) error
	DeleteFunc func(ctx context.Context, id int64) error
}

func (s *stubPostMutator) Update(ctx context.Context, p posts.StoredPost) error {
	return s.UpdateFunc(ctx, p)
}

func (s *stubPostMutator) DeletePostByID(ctx context.Context, id int64) error {
	return s.DeleteFunc(ctx, id)
}
