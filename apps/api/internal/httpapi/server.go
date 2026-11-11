// Package httpapi exposes the public API (design v1.1 §9).
package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"

	"codebaseintel/internal/config"
	"codebaseintel/internal/gemini"
	"codebaseintel/internal/github"
	"codebaseintel/internal/store"
)

type Server struct {
	Store  *store.Store
	Gemini *gemini.Client
	Cfg    config.Config
}

func New(st *store.Store, gm *gemini.Client, cfg config.Config) *Server {
	return &Server{Store: st, Gemini: gm, Cfg: cfg}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.health)
	mux.HandleFunc("POST /api/analyses", s.createAnalysis)
	mux.HandleFunc("GET /api/analyses/{id}", s.getAnalysis)
	mux.HandleFunc("POST /api/analyses/{id}/questions", s.askQuestion)
	return cors(mux)
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, code, msg string) {
	writeJSON(w, status, map[string]any{"error": map[string]any{"code": code, "message": msg}})
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "version": "0.1.0"})
}

func (s *Server) createAnalysis(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RepoURL string `json:"repo_url"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 10<<10)).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, "INVALID_INPUT", "body must be {\"repo_url\": \"https://github.com/owner/repo\"}")
		return
	}
	owner, repo, err := github.ParseRepoURL(body.RepoURL)
	if err != nil {
		writeErr(w, http.StatusUnprocessableEntity, "INVALID_URL", err.Error())
		return
	}
	a, err := s.Store.CreateAnalysis(r.Context(), owner, repo, body.RepoURL)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "INTERNAL", "could not create analysis")
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]any{
		"analysis_id": a.ID,
		"status":      a.Status,
		"repository":  a.RepoFullName,
	})
}

func (s *Server) getAnalysis(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	a, err := s.Store.GetAnalysis(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeErr(w, http.StatusNotFound, "NOT_FOUND", "analysis not found")
			return
		}
		writeErr(w, http.StatusInternalServerError, "INTERNAL", "could not load analysis")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"id":               a.ID,
		"repository":       a.RepoFullName,
		"default_branch":   a.DefaultBranch,
		"commit_sha":       a.CommitSHA,
		"status":           a.Status,
		"status_detail":    a.StatusDetail,
		"size_kb":          a.SizeKB,
		"counts": map[string]int{
			"files":         a.FileCount,
			"indexed_files": a.IndexedFileCount,
			"symbols":       a.SymbolCount,
			"chunks":        a.ChunkCount,
		},
		"parse_error_rate": a.ParseErrorRate,
	})
}
