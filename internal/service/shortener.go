package service

import (
	"context"
	"fmt"
	"github.com/sssnukke/url-shortener/internal/model"
	"github.com/sssnukke/url-shortener/internal/repository"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type ShortenerService interface {
	Shorten(ctx context.Context, originalURL string) (*model.ShortenResponse, error)
	Redirect(ctx context.Context, shortCode string) (string, error)
	GetStats(ctx context.Context, shortCode string) (*model.StaticResponse, error)
}

type shortenerService struct {
	urlRepo   repository.URLRepository
	cacheRepo repository.CacheRepository
	codeLen   int
}

func NewShortenerService(urlRepo repository.URLRepository, cacheRepo repository.CacheRepository) ShortenerService {
	codeLen, err := strconv.Atoi(os.Getenv("SHORT_CODE_LENGTH"))
	if err != nil && codeLen <= 0 {
		codeLen = 6
	}

	return &shortenerService{
		urlRepo:   urlRepo,
		cacheRepo: cacheRepo,
		codeLen:   codeLen,
	}
}

func (s *shortenerService) Shorten(ctx context.Context, originalURL string) (*model.ShortenResponse, error) {
	if err := validateURL(originalURL); err != nil {
		return nil, err
	}

	shortCode := s.generateCode()

	url := &model.URL{
		ShortCode:   shortCode,
		OriginalURl: originalURL,
	}

	if err := s.urlRepo.Create(ctx, url); err != nil {
		return nil, fmt.Errorf("failed to save url: %w", err)
	}

	_ = s.cacheRepo.Set(ctx, shortCode, originalURL, 24*time.Hour)

	return &model.ShortenResponse{
		ShortCode:   shortCode,
		ShortURL:    "http://localhost:8080/" + shortCode,
		OriginalURL: originalURL,
	}, nil
}

func (s *shortenerService) Redirect(ctx context.Context, shortCode string) (string, error) {
	cache, err := s.cacheRepo.Get(ctx, shortCode)
	if err == nil && cache != "" {
		go s.urlRepo.IncrementClicks(context.Background(), shortCode)
		return cache, nil
	}
	url, err := s.urlRepo.GetByCode(ctx, shortCode)
	if err != nil {
		return "", fmt.Errorf("failed to get url: %w", err)
	}

	_ = s.cacheRepo.Set(ctx, shortCode, url.OriginalURl, 24*time.Hour)

	go s.urlRepo.IncrementClicks(context.Background(), shortCode)
	return url.OriginalURl, nil
}

func (s *shortenerService) GetStats(ctx context.Context, shortCode string) (*model.StaticResponse, error) {
	url, err := s.urlRepo.GetByCode(ctx, shortCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get url: %w", err)
	}

	return &model.StaticResponse{
		ShortCode:   url.ShortCode,
		OriginalURL: url.OriginalURl,
		Clicks:      url.Clicks,
		CreatedAt:   url.CreatedAt,
	}, nil
}

func (s *shortenerService) generateCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	code := make([]byte, s.codeLen)
	for i := range code {
		code[i] = alphabet[r.Intn(len(alphabet))]
	}
	return string(code)
}

func validateURL(rawURL string) error {
	if rawURL == "" {
		return fmt.Errorf("url is empty")
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid url: %w", err)
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("url must start with http:// or https://")
	}

	if parsed.Host == "" {
		return fmt.Errorf("url must have a host")
	}

	return nil
}
