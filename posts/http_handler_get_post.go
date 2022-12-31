package posts

import (
	"encoding/json"
	"errors"
	"net/http"
	"simpleblog/httphandler"
	"strconv"
	"time"
)

func HandleGetAllPosts(pa PostAccessor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		posts, err := pa.GetAll(r.Context())
		if errors.Is(err, ErrNotFound) {
			_ = json.NewEncoder(w).Encode(nil)
			return
		}
		if err != nil {
			httphandler.HandleInternalServerError(err)(w, r)
			return
		}

		type resp struct {
			ID        int64     `json:"id"`
			Title     string    `json:"title"`
			Content   string    `json:"content"`
			Timestamp time.Time `json:"timestamp"`
		}

		var data []resp

		for _, r := range posts {
			data = append(data, resp(r))
		}
		_ = json.NewEncoder(w).Encode(data)
	}
}

func HandleGetPostByID(pa PostAccessor, paramID string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(paramID, 0, 64)
		if err != nil {
			httphandler.HandleNotFound()(w, r)
		}

		post, err := pa.GetByID(r.Context(), id)
		if err != nil {
			httphandler.HandleInternalServerError(err)(w, r)
			return
		}

		type resp struct {
			ID        int64     `json:"id"`
			Title     string    `json:"title"`
			Content   string    `json:"content"`
			Timestamp time.Time `json:"timestamp"`
		}
		_ = json.NewEncoder(w).Encode([]resp{
			resp(post),
		})
	}
}
