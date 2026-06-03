package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sssnukke/url-shortener/internal/model"
)

type URLRepository interface {
	Create(ctx context.Context, url *model.URL) error
	GetByCode(ctx context.Context, shortCode string) (*model.URL, error)
	IncrementClicks(ctx context.Context, shortCode string) error
}

type postgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRepo(db *pgxpool.Pool) URLRepository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) Create(ctx context.Context, url *model.URL) error {
	query := `
		INSERT INTO url (short_code, original_url)
		VALUES ($1, $2)
		RETURNING id, created_at
	`

	return r.db.QueryRow(ctx, query, url.ShortCode, url.OriginalURl).Scan(&url.ID, &url.CreatedAt)
}

func (r *postgresRepo) GetByCode(ctx context.Context, shortCode string) (*model.URL, error) {
	query := `
		SELECT id, short_code, original_url, clicks, created_at
		FROM urls WHERE short_code = $1
	`

	url := &model.URL{}

	err := r.db.QueryRow(ctx, query, shortCode).Scan(&url.ID, &url.ShortCode, &url.OriginalURl, &url.Clicks, &url.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("url not found: %w", err)
	}

	return url, nil
}

func (r *postgresRepo) IncrementClicks(ctx context.Context, shortCode string) error {
	query := `
		UPDATE urls SET clicks = clicks + 1
		WHERE short_code = $1
	`
	
	_, err := r.db.Exec(ctx, query, shortCode)
	return err
}
