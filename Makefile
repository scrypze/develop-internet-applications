PHONY: up
up:
	go run cmd/exocalc/main.go

PHONY: migrate
migrate:
	go run cmd/migrate/main.go
