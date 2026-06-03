package handler

import (
	"encoding/json"
	"github.com/sssnukke/url-shortener/internal/model"
	"github.com/sssnukke/url-shortener/internal/service"
	"net/http"
)

type ShortenerHandler struct {
	svc service.ShortenerService
}

func NewShortenerHandler(svc service.ShortenerService) *ShortenerHandler {
	return &ShortenerHandler{
		svc: svc,
	}
}

func (h *ShortenerHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	var req model.ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}

	resp, err := h.svc.Shorten(r.Context(), req.URL)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}
