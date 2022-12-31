package main

import (
	"hash"
	"io"
	"net/http"
	"simpleblog/auth"
	"simpleblog/httphandler"
	"simpleblog/posts"
	"strings"
)

type RouterDeps struct {
	UserStorage       auth.UserStorage
	PostStorage       posts.PostStorage
	PasswordGenerator func() (string, error)
	PasswordHash      hash.Hash
	TokenGenerator    func(httphandler.JWTClaims) (string, error)
	TokenParser       func(string) (httphandler.JWTClaims, error)
}

func RouteHTTPHandler(deps RouterDeps) http.Handler {
	useMiddlewares := func(h http.Handler) http.HandlerFunc {
		return httphandler.MiddlewareCORS([]string{"*"})(
			httphandler.MiddlewareJSONResp(
				h,
			),
		)
	}

	mux := http.NewServeMux()

	mux.Handle("/register", useMiddlewares(httphandler.HandleMethod(
		auth.HandleRegister(deps.UserStorage, deps.PasswordHash, deps.PasswordGenerator),
		http.MethodPost),
	))

	mux.Handle("/get-token", useMiddlewares(httphandler.HandleMethod(
		auth.HandleLogin(deps.UserStorage, deps.PasswordHash, deps.TokenGenerator),
		http.MethodGet),
	))

	mux.Handle("/posts", useMiddlewares(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleMethodGet := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if id := r.URL.Query().Get("id"); id != "" {
				posts.HandleGetPostByID(deps.PostStorage, id).ServeHTTP(w, r)
				return
			}
			posts.HandleGetAllPosts(deps.PostStorage)(w, r)
		})

		handleMethodPost := posts.HandleCreatePost(deps.PostStorage)

		handleMethodPut := posts.HandleUpdatePost(deps.PostStorage)

		handleMethodDelete := posts.HandleDeletePostByID(deps.PostStorage, r.URL.Query().Get("id"))

		adminOnly := httphandler.MiddlewareAllowRoles(deps.TokenParser, auth.RoleAdmin)
		allRole := httphandler.MiddlewareAllowRoles(deps.TokenParser, auth.RoleAdmin, auth.RoleUser)

		switch r.Method {
		case http.MethodGet:
			allRole(handleMethodGet).ServeHTTP(w, r)
		case http.MethodPost:
			adminOnly(handleMethodPost).ServeHTTP(w, r)
		case http.MethodPut:
			adminOnly(handleMethodPut).ServeHTTP(w, r)
		case http.MethodDelete:
			adminOnly(handleMethodDelete).ServeHTTP(w, r)
		default:
			httphandler.HandleError(nil, "method not allowed", http.StatusMethodNotAllowed).ServeHTTP(w, r)
		}
	}),
	))

	mux.Handle("/swaggerui/",
		http.FileServer(http.FS(swaggeruifs)),
	)

	mux.Handle("/swagger.yml", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "text/yaml")
			_, _ = io.Copy(w, strings.NewReader(swaggerContent))
		},
	))

	return mux
}
