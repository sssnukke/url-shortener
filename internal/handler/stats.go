package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/sssnukke/url-shortener/internal/service"
	"net/http"
)

type StatsHandler struct {
	svc service.ShortenerService
}

func NewStatsHandler(svc service.ShortenerService) *StatsHandler {
	return &StatsHandler{
		svc: svc,
	}
}

func (h *StatsHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "code")

	stats, err := h.svc.GetStats(r.Context(), shortCode)
	if err != nil {
		writeError(w, http.StatusNotFound, "not found url")
		return
	}

	writeJSON(w, http.StatusOK, stats)
}
