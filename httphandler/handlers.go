package httphandler

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"
)

func HandleError(err error, msg string, statusCode int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err != nil {
			log.Printf("%s %s %v", r.Method, r.URL.Path, err)
		}
		w.WriteHeader(statusCode)
		_ = json.NewEncoder(w).Encode(msg)
	}
}

func HandleInternalServerError(err error) http.HandlerFunc {
	return HandleError(err, "something bad happened on our side", http.StatusInternalServerError)
}

func HandleBadRequest() http.HandlerFunc {
	return HandleError(nil, "bad request", http.StatusBadRequest)
}

func HandleUnauthorized() http.HandlerFunc {
	return HandleError(nil, "unauthorized", http.StatusUnauthorized)
}

func HandleForbidden() http.HandlerFunc {
	return HandleError(nil, "forbidden", http.StatusForbidden)
}

func HandleNotFound() http.HandlerFunc {
	return HandleError(nil, "not found", http.StatusNotFound)
}

func HandleMethod(h http.Handler, methods ...string) http.HandlerFunc {
	isMethodAvailable := func(method string) bool {
		for _, m := range methods {
			if method == m {
				return true
			}
		}
		return false
	}

	handleOptions := func(w http.ResponseWriter) {
		slice := sort.StringSlice(methods)
		slice.Sort()
		allowedMethods := strings.Join(methods, ", ")
		w.Header().Set("Allow", allowedMethods)
		w.WriteHeader(http.StatusNoContent)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			handleOptions(w)
			return
		}

		if !isMethodAvailable(r.Method) {
			HandleError(nil, "method not allowed", http.StatusMethodNotAllowed)(w, r)
			return
		}

		h.ServeHTTP(w, r)
	}
}
