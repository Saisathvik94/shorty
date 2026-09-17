package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type URLRepository struct {
	db *pgxpool.Pool
}

func NewURLRepository(db *pgxpool.Pool) *URLRepository {
	return &URLRepository{
		db: db,
	}
}

func (r *URLRepository) CreateURL(
	ctx context.Context,
	shortCode string,
	originalURL string,
) error {
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
