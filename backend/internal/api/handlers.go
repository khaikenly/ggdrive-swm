package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ggdrive-swm/backend/internal/auth"
	"github.com/ggdrive-swm/backend/internal/drive"
	"golang.org/x/oauth2"
)

const maxFolderIDs = 50

type Handler struct {
	oauthConfig  *oauth2.Config
	sessionStore *auth.Store
	driveClient  func(token *oauth2.Token) (*drive.Client, error)
}

func NewHandler(oauthConfig *oauth2.Config, sessionStore *auth.Store, driveClient func(token *oauth2.Token) (*drive.Client, error)) *Handler {
	return &Handler{
		oauthConfig:  oauthConfig,
		sessionStore: sessionStore,
		driveClient:  driveClient,
	}
}

func (h *Handler) AuthLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	state := generateState()
	url := h.oauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusFound)
}

func (h *Handler) AuthCallback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "missing code", http.StatusBadRequest)
		return
	}

	token, err := h.oauthConfig.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "exchange failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.sessionStore.Set(w, token)

	frontendURL := r.Context().Value(frontendURLKey).(string)
	http.Redirect(w, r, frontendURL+"/", http.StatusFound)
}

func (h *Handler) AuthLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.sessionStore.Delete(w)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) AuthMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token, ok := h.sessionStore.Get(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	writeJSON(w, map[string]bool{"authenticated": token != nil})
}

func (h *Handler) ListFolders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	client, err := h.authenticatedClient(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	folders, err := client.ListSharedFolders(r.Context())
	if err != nil {
		http.Error(w, "list folders: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]interface{}{"folders": folders})
}

func (h *Handler) ListFolderChildren(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	folderID := strings.TrimPrefix(r.URL.Path, "/api/folders/")
	folderID = strings.TrimSuffix(folderID, "/children")
	folderID = strings.TrimSpace(folderID)

	if folderID == "" {
		http.Error(w, "folder id required", http.StatusBadRequest)
		return
	}

	client, err := h.authenticatedClient(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	folders, videos, err := client.ListSubfoldersAndVideos(r.Context(), folderID)
	if err != nil {
		http.Error(w, "list children: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]interface{}{
		"folders": folders,
		"videos":  videos,
	})
}

type buildCourseRequest struct {
	FolderIDs []string `json:"folderIds"`
}

func (h *Handler) BuildCourse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req buildCourseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if len(req.FolderIDs) == 0 {
		http.Error(w, "folderIds required", http.StatusBadRequest)
		return
	}

	if len(req.FolderIDs) > maxFolderIDs {
		http.Error(w, "too many folders", http.StatusBadRequest)
		return
	}

	client, err := h.authenticatedClient(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	course, err := client.BuildCourse(r.Context(), req.FolderIDs)
	if err != nil {
		http.Error(w, "build course: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, course)
}

func (h *Handler) StreamVideo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fileID := strings.TrimPrefix(r.URL.Path, "/api/videos/")
	fileID = strings.TrimSuffix(fileID, "/stream")
	fileID = strings.TrimSpace(fileID)

	if fileID == "" {
		http.Error(w, "file id required", http.StatusBadRequest)
		return
	}

	client, err := h.authenticatedClient(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	resp, err := client.DownloadFile(r.Context(), fileID)
	if err != nil {
		http.Error(w, "download failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "drive error: "+resp.Status, resp.StatusCode)
		return
	}

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.Header().Set("Accept-Ranges", "bytes")
	if resp.ContentLength > 0 {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", resp.ContentLength))
	}
	io.Copy(w, resp.Body)
}

func (h *Handler) authenticatedClient(r *http.Request) (*drive.Client, error) {
	token, ok := h.sessionStore.Get(r)
	if !ok || token == nil {
		return nil, ErrUnauthorized
	}

	return h.driveClient(token)
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

var ErrUnauthorized = &authError{msg: "unauthorized"}

type authError struct {
	msg string
}

func (e *authError) Error() string {
	return e.msg
}

func generateState() string {
	return "random-state"
}
