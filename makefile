include .env

run:
	go run ./cmd/shortener

test:
	go test ./cmd/shortener