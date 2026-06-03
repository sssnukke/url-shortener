package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sssnukke/url-shortener/internal/service"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, ErrorResponse{Error: msg})
}

func NewRouter(svc service.ShortenerService) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RealIP)

	r.Use(middleware.Logger)

	r.Use(middleware.Recoverer)

	r.Post("/api/shorten", NewShortenerHandler(svc).Shorten)
	r.Get("/{code}", NewRedirectHandler(svc).Redirect)
	r.Get("/api/stats/{code}", NewStatsHandler(svc).GetStats)

	return r
}
