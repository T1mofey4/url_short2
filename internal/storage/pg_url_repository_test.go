package storage

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"url_shortener/internal/app"
)

func TestURLRepository(t *testing.T) {
	ctx := context.Background()

	// // to start container PG
	// req := testcontainers.ContainerRequest{
	// 	Image:        "postgres:14-alpine3.17",
	// 	ExposedPorts: []string{"5432/tcp"},
	// 	Env: map[string]string{
	// 		"PG_USER":     "testuser",
	// 		"PG_PASSWORD": "testpass",
	// 		"PG_DB":       "testdb",
	// 	},
	// 	WaitingFor: wait.ForListeningPort("5432/tcp"),
	// },

	// container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
	// 	ContainerRequest: req,
	// 	Started:          true,
	// })
	// if err != nil {
	// 	t.Fatalf("failed to start container: %v", err)
	// }
	// defer container.Terminate(ctx)

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:14-alpine3.17",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_USER":     "testuser",
				"POSTGRES_PASSWORD": "password",
				"POSTGRES_DB":       "testdb",
			},
			WaitingFor: wait.ForListeningPort("5432/tcp"),
		},
		Started: true,
	})
	defer container.Terminate(ctx)

	// Получить информацию о контейнере
	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container host: %v", err)
	}
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("failed to get container port: %v", err)
	}

	// Составить DSN для подключения к БД
	dsn := fmt.Sprintf("postgres://testuser:password@%s:%s/testdb?sslmode=disable", host, port.Port())

	// Создать пул подключений к БД
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Создать репозиторий
	repo := NewURLRepository(pool)

	// Создать URL
	url := app.URL{
		Slug:      "test-slug",
		LongURL:   "https://example.com",
		TTL:       time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}

	// Сохранить URL
	slug, err := repo.Save(ctx, url)
	if err != nil {
		t.Fatalf("failed to save URL: %v", err)
	}
	if slug != url.Slug {
		t.Errorf("expected slug %q, got %q", url.Slug, slug)
	}
}
