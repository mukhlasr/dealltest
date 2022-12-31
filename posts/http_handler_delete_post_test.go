package posts_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"simpleblog/posts"
	"testing"
)

func TestDeletePostByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/?id=123", nil)
		rec := httptest.NewRecorder()

		posts.HandleDeletePostByID(&stubPostMutator{
			DeleteFunc: func(ctx context.Context, id int64) error {
				return nil
			},
		}, "123").ServeHTTP(rec, req)

		if statusCode := rec.Result().StatusCode; statusCode != http.StatusNoContent {
			t.Fatal("expecting 201 OK but got:", statusCode)
		}
	})
}
