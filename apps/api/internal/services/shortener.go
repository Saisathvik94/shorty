package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/url"
	"time"

	"github.com/Saisathvik94/shorty/apps/api/internal/cache"
	"github.com/Saisathvik94/shorty/apps/api/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/redis/go-redis/v9"
)

type URLService struct {
	repo  *repository.URLRepository
	cache *cache.RedisCache
}

func NewURLService(cache *cache.RedisCache, repo *repository.URLRepository) *URLService {
	return &URLService{
		cache: cache,
		repo:  repo,
	}
}

func generateShortCode() (string, error) {
	length := 6
	bytes := make([]byte, length)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

func (s *URLService) CreateURL(ctx context.Context, originalURL string, expiresAt string) (string, error) {
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

	expiresAtTime, err := time.Parse(time.RFC3339, expiresAt)

	if err != nil {
		return "", errors.New("invalid expiration time")
	}

	now := time.Now()

	if !expiresAtTime.After(now) {
		return "", errors.New("expiration time must be in the future")
	}

	maxExpiration := now.Add(30 * 24 * time.Hour)

	if expiresAtTime.After(maxExpiration) {
		return "", errors.New("expiration cannot exceed 30 days")
	}

	for {
		shortCode, err := generateShortCode()
		if err != nil {
			return "", err
		}

		err = s.repo.CreateURL(ctx, shortCode, originalURL, expiresAtTime)
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

var (
	ErrorURLNotFound = errors.New("URL is Not Found")
	ErrorURLIsactive = errors.New("URL is inactive")
	ErrorExpired     = errors.New("URL has Expired")
)

func (s *URLService) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	var record repository.URLRecord
	cacheKey := "url:" + shortCode
	cacheRecord, err := s.cache.Get(ctx, cacheKey)

	if errors.Is(err, redis.Nil) {
		dbRecord, err := s.repo.GetURLByShortCode(ctx, shortCode)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return "", ErrorURLNotFound
			}
			return "", err
		}
		record = *dbRecord
		data, err := json.Marshal(record)

		if err != nil {
			return "", err
		}

		if record.ExpiresAt == nil {
			return "", errors.New("URL has no expiration time")
		}

		ttl := time.Until(*record.ExpiresAt)

		if ttl <= 0 {
			return "", ErrorExpired
		}

		setErr := s.cache.Set(ctx, cacheKey, string(data), ttl)
		if setErr != nil {
			return "", setErr
		}
	} else if err != nil {
		return "", err
	} else {
		err := json.Unmarshal([]byte(cacheRecord), &record)

		if err != nil {
			return "", err
		}
	}

	if !record.IsActive {
		return "", ErrorURLIsactive
	}
	if record.ExpiresAt != nil {
		currentTime := time.Now()
		if currentTime.After(*record.ExpiresAt) || currentTime.Equal(*record.ExpiresAt) {
			return "", ErrorExpired
		}
	}

	return record.OriginalURL, nil
}

func (s *URLService) DeactivateURL(ctx context.Context, shortCode string) error {
	found, err := s.repo.DeactivateURL(ctx, shortCode)

	if err != nil {
		return err
	}

	if !found {
		return ErrorURLNotFound
	}

	cacheErr := s.cache.Delete(ctx, "url:"+shortCode)

	if cacheErr != nil {
		return cacheErr
	}

	return nil
}
func (s *URLService) DeleteURL(ctx context.Context, shortCode string) error {
	found, err := s.repo.DeleteURL(ctx, shortCode)

	if err != nil {
		return err
	}

	if !found {
		return ErrorURLNotFound
	}

	cacheErr := s.cache.Delete(ctx, "url:"+shortCode)

	if cacheErr != nil {
		return cacheErr
	}

	return nil
}

func (s *URLService) UpdateExpiration(ctx context.Context, expiresAt string, shortCode string) error {
	expiresAtTime, err := time.Parse(time.RFC3339, expiresAt)

	if err != nil {
		return errors.New("invalid expiration time")
	}

	now := time.Now()

	if !expiresAtTime.After(now) {
		return errors.New("expiration time must be in the future")
	}

	maxExpiration := now.Add(30 * 24 * time.Hour)

	if expiresAtTime.After(maxExpiration) {
		return errors.New("expiration cannot exceed 30 days")
	}

	found, err := s.repo.UpdateExpiration(ctx, expiresAtTime, shortCode)

	if err != nil {
		return err
	}

	if !found {
		return ErrorURLNotFound
	}

	cacheErr := s.cache.Delete(ctx, "url:"+shortCode)

	if cacheErr != nil {
		return cacheErr
	}

	return nil
}
