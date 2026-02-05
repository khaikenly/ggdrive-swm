package main

import (
	"context"
	"log"
	"net/http"

	"github.com/ggdrive-swm/backend/internal/api"
	"github.com/ggdrive-swm/backend/internal/auth"
	"github.com/ggdrive-swm/backend/internal/config"
	"github.com/ggdrive-swm/backend/internal/drive"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()
	if cfg.GoogleClientID == "" || cfg.GoogleClientSecret == "" {
		log.Fatal("GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET must be set")
	}

	oauthConfig := auth.NewOAuth2Config(cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.BackendURL)
	sessionStore := auth.NewStore(cfg.SessionSecret)

	driveClientFn := func(token *oauth2.Token) (*drive.Client, error) {
		svc, err := auth.DriveServiceFromToken(context.Background(), token)
		if err != nil {
			return nil, err
		}
		return drive.NewClient(svc), nil
	}

	h := api.NewHandler(oauthConfig, sessionStore, driveClientFn)
	router := h.Router(cfg.FrontendURL)

	withCORS := corsMiddleware(router, cfg.FrontendURL)

	log.Printf("Listening on %s", cfg.BackendURL)
	log.Fatal(http.ListenAndServe(":8080", withCORS))
}

func corsMiddleware(next http.Handler, frontendURL string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", frontendURL)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
