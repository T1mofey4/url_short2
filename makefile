include .env

LOCAL_BIN:=$(CURDIR)/bin

run:
	go run ./cmd/shortener

test:
	go test ./cmd/shortener

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0

migration-status:
	${BIN}/goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} status -v

migration-up:
	${BIN}/goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} up -v

migration-down:
	${BIN}/goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} down -v
