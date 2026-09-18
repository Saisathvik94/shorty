package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type URLRepository struct {
	db *pgxpool.Pool
}

type URLRecord struct {
	OriginalURL string
	IsActive    bool
	ExpiresAt   *time.Time
}

func NewURLRepository(db *pgxpool.Pool) *URLRepository {
	return &URLRepository{
		db: db,
	}
}

func (r *URLRepository) CreateURL(ctx context.Context, shortCode string, originalURL string) error {
	_, err := r.db.Exec(
		ctx,
		`
		INSERT INTO urls (short_code, original_url)
		VALUES ($1, $2)
		`,
		shortCode,
		originalURL,
	)

	return err
}

func (r *URLRepository) GetURLByShortCode(ctx context.Context, shortCode string) (*URLRecord, error) {
	var record URLRecord

	err := r.db.QueryRow(
		ctx,
		`
		SELECT original_url, is_active, expires_at FROM urls WHERE short_code = $1`,
		shortCode,
	).Scan(
		&record.OriginalURL,
		&record.IsActive,
		&record.ExpiresAt,
	)

	if err != nil {
		return nil, err
	}

	return &record, nil
}
