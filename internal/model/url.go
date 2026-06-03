package model

import (
	"github.com/google/uuid"
	"time"
)

type URL struct {
	ID          uuid.UUID `db:"id"`
	ShortCode   string    `db:"short_code"`
	OriginalURl string    `db:"original_uri"`
	Clicks      int64     `db:"clicks"`
	CreatedAt   time.Time `db:"created_at"`
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortCode   string `json:"short_code"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type StaticResponse struct {
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	Clicks      int64     `json:"clicks"`
	CreatedAt   time.Time `json:"created_at"`
}
