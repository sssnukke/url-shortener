package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/sssnukke/url-shortener/internal/service"
	"net/http"
)

type RedirectHandler struct {
	svc service.ShortenerService
}

func NewRedirectHandler(svc service.ShortenerService) *RedirectHandler {
	return &RedirectHandler{
		svc: svc,
	}
}

func (h *RedirectHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "code")

	originalURL, err := h.svc.Redirect(r.Context(), shortCode)
	if err != nil {
		writeError(w, http.StatusNotFound, "not found url")
		return
	}

	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}
