package storage

import (
	"context"

	"url_shortener/internal/app"
)

type URLRepository interface {
	Save(ctx context.Context, url app.URL) (string, error)
	GetBySlug(ctx context.Context, slug string) (*app.URL, error)
	DeleteExpired(ctx context.Context,) (int64, error)
}