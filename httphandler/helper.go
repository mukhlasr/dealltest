package httphandler

import (
	"net/http"
	"strings"
)

func MiddlewareJSONResp(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	}
}

func MiddlewareCORS(origins []string) func(http.Handler) http.HandlerFunc {
	return func(h http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", strings.Join(origins, ", "))
			h.ServeHTTP(w, r)
		}
	}
}
