package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

const (
	sessionCookieName = "ggdrive_session"
	sessionMaxAge     = 24 * time.Hour
)

type Session struct {
	Token      *oauth2.Token
	ExpiresAt  time.Time
}

type Store struct {
	mu      sync.RWMutex
	sessions map[string]*Session
	secret  []byte
}

func NewStore(secret string) *Store {
	gob.Register(&oauth2.Token{})
	return &Store{
		sessions: make(map[string]*Session),
		secret:   []byte(secret),
	}
}

func (s *Store) Set(w http.ResponseWriter, token *oauth2.Token) {
	id := generateSessionID()
	s.mu.Lock()
	s.sessions[id] = &Session{
		Token:     token,
		ExpiresAt: time.Now().Add(sessionMaxAge),
	}
	s.mu.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    id,
		Path:     "/",
		MaxAge:   int(sessionMaxAge.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
	})
}

func (s *Store) Get(r *http.Request) (*oauth2.Token, bool) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return nil, false
	}

	s.mu.RLock()
	sess, ok := s.sessions[cookie.Value]
	s.mu.RUnlock()

	if !ok || sess == nil || time.Now().After(sess.ExpiresAt) {
		return nil, false
	}

	return sess.Token, true
}

func (s *Store) Delete(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
}

func generateSessionID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return uuid.New().String()
	}
	return base64.URLEncoding.EncodeToString(b)
}

func encodeTokenForCookie(token *oauth2.Token) (string, error) {
	data, err := json.Marshal(token)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(data), nil
}
