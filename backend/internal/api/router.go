package api

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const frontendURLKey contextKey = "frontend_url"

func (h *Handler) Router(frontendURL string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), frontendURLKey, frontendURL)
		r = r.WithContext(ctx)

		path := strings.TrimSuffix(r.URL.Path, "/")

		switch {
		case path == "/api/auth/login":
			h.AuthLogin(w, r)
		case path == "/api/auth/callback":
			h.AuthCallback(w, r)
		case path == "/api/auth/logout" && r.Method == http.MethodPost:
			h.AuthLogout(w, r)
		case path == "/api/auth/me":
			h.AuthMe(w, r)
		case path == "/api/folders":
			h.ListFolders(w, r)
		case strings.HasPrefix(path, "/api/folders/") && strings.HasSuffix(path, "/children"):
			h.ListFolderChildren(w, r)
		case path == "/api/courses/build" && r.Method == http.MethodPost:
			h.BuildCourse(w, r)
		case strings.HasPrefix(path, "/api/videos/") && strings.HasSuffix(path, "/stream"):
			h.StreamVideo(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}
