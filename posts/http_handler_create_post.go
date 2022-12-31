package posts

import (
	"encoding/json"
	"net/http"
	"simpleblog/httphandler"
	"time"
)

func HandleCreatePost(pc PostCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := struct {
			Title   string
			Content string
		}{}

		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httphandler.HandleError(err, "bad request body", http.StatusBadRequest)(w, r)
			return
		}

		timestamp := time.Now().UTC().Truncate(time.Millisecond)

		post := StoredPost{
			Title:     req.Title,
			Content:   req.Content,
			Timestamp: timestamp,
		}

		if !post.isValid() {
			httphandler.HandleBadRequest()(w, r)
			return
		}

		id, err := pc.Create(r.Context(), post)
		if err != nil {
			httphandler.HandleInternalServerError(err)(w, r)
			return
		}

		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":        id,
			"title":     req.Title,
			"content":   req.Content,
			"timestamp": timestamp,
		})
	}
}
