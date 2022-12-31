package posts

import (
	"encoding/json"
	"log"
	"net/http"
	"simpleblog/httphandler"
	"time"
)

func HandleUpdatePost(pm PostMutator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("masuk siin")
		req := struct {
			ID      int64
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
			ID:        req.ID,
			Title:     req.Title,
			Content:   req.Content,
			Timestamp: timestamp,
		}

		if !post.isValid() {
			httphandler.HandleBadRequest()(w, r)
			return
		}

		err := pm.Update(r.Context(), post)
		if err != nil {
			httphandler.HandleInternalServerError(err)(w, r)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
