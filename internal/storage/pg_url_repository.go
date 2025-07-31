package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"

	"url_shortener/internal/app"
)

type pgURLRepository struct {
	pool *pgxpool.Pool
}

func NewURLRepository(pool *pgxpool.Pool) URLRepository {
	return &pgURLRepository{pool: pool}
}

func (p *pgURLRepository) Save(ctx context.Context, url app.URL) (string, error) {
	var slug string
	err := p.pool.QueryRow(ctx, `
		INSERT INTO urls(slug, long_url, ttl, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING slug
	`, url.Slug, url.LongURL, url.TTL, url.CreatedAt).Scan(&slug)
	if err != nil {
		return "", err
	}

	return slug, nil
}

func (p *pgURLRepository) GetBySlug(ctx context.Context, slug string) (*app.URL, error) {
	var url app.URL
	err := p.pool.QueryRow(ctx, `
		SELECT id, slug, long_url, ttl, created_at
		FROM urls
		WHERE slug = $1
	`, slug).Scan(&url.ID, &url.Slug, &url.LongURL, &url.TTL, &url.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
	}

	return &url, nil
}

func (p *pgURLRepository) DeleteExpired(ctx context.Context) (int64, error) {
	res, err := p.pool.Exec(ctx, `
		DELETE FROM urls
		WHERE ttl < NOW()
	`)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected(), nil
}
