package storage

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"url_shortener/internal/app"
)

func TestURLRepository(t *testing.T) {
	ctx := context.Background()

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

	// Миграции для БД
	cmd := exec.Command("../../bin/goose", "-dir", "../../migrations", "postgres", dsn, "up")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	require.NoError(t, cmd.Run())

	// Создать пул подключений к БД
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Создать репозиторий
	repo := NewURLRepository(pool)

	t.Run("Save/Get round-trip", func(t *testing.T) {
		example := app.URL{
			Slug:      "abc123",
			LongURL:   "https;//examle.com",
			TTL:       time.Now().Add(24 * time.Hour),
			CreatedAt: time.Now(),
		}

		err := repo.Save(ctx, example)
		require.NoError(t, err)

		got, err := repo.GetBySlug(ctx, example.Slug)
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, example.LongURL, got.LongURL)
	})

	t.Run("DeleteExpired", func(t *testing.T) {
		count, err := repo.DeleteExpired(ctx)
		require.NoError(t, err)
		require.Equal(t, int64(0), count)
	})
}
