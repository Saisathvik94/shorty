package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/url"

	"github.com/Saisathvik94/shorty/apps/api/internal/repository"
	"github.com/jackc/pgx/v5/pgconn"
)

type URLService struct {
	repo *repository.URLRepository
}

func generateShortCode() (string, error) {
	length := 6
	bytes := make([]byte, length)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

func (s *URLService) CreateURL(ctx context.Context, originalURL string) (string, error) {
	u, err := url.Parse(originalURL)
	if err != nil {
		return "", err
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return "", errors.New("invalid URL scheme")
	}

	if u.Hostname() == "" {
		return "", errors.New("URL must have a hostname")
	}

	for {
		shortCode, err := generateShortCode()
		if err != nil {
			return "", err
		}

		err = s.repo.CreateURL(ctx, shortCode, originalURL)
		if err == nil {
			return shortCode, nil
		}

		if isUniqueViolation(err) {
			continue
		}

		return "", err
	}
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}

	return false
}
