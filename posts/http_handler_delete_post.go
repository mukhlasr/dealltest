package posts

import (
	"errors"
	"net/http"
	"simpleblog/httphandler"
	"strconv"
)

func HandleDeletePostByID(pm PostMutator, idParam string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(idParam, 0, 64)
		if err != nil {
			httphandler.HandleBadRequest()(w, r)
			return
		}
		err = pm.DeletePostByID(r.Context(), id)
		if errors.Is(err, ErrNotFound) {
			httphandler.HandleNotFound()(w, r)
			return
		}
		if err != nil {
			httphandler.HandleInternalServerError(err)(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
