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

func (r *URLRepository) CreateURL(ctx context.Context, shortCode string, originalURL string, expiresAt time.Time) error {
	_, err := r.db.Exec(
		ctx,
		`
		INSERT INTO urls (short_code, original_url, expires_at)
		VALUES ($1, $2, $3)
		`,
		shortCode,
		originalURL,
		expiresAt,
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

func (r *URLRepository) DeactivateURL(ctx context.Context, shortCode string) (bool, error) {
	result, err := r.db.Exec(
		ctx,
		`UPDATE urls SET is_active = false WHERE short_code = $1`, shortCode,
	)

	if err != nil {
		return false, err
	}

	return result.RowsAffected() > 0, nil
}

func (r *URLRepository) DeleteURL(ctx context.Context, shortCode string) (bool, error) {
	result, err := r.db.Exec(
		ctx,
		`DELETE FROM urls WHERE short_code = $1`,
		shortCode,
	)

	if err != nil {
		return false, err
	}

	return result.RowsAffected() > 0, nil

}
func (r *URLRepository) UpdateExpiration(ctx context.Context, expiresAt time.Time, shortCode string) (bool, error) {
	result, err := r.db.Exec(
		ctx,
		`UPDATE urls SET expires_at = $1 WHERE short_code = $2`,
		expiresAt,
		shortCode,
	)

	if err != nil {
		return false, err
	}

	return result.RowsAffected() > 0, nil
}
