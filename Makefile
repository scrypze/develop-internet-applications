PHONY: up
up:
	go run cmd/exocalc/main.go

PHONY: migrate
migrate:
	go run cmd/migrate/main.go

PHONY: insert
insert-stars:
	docker compose exec -T postgres psql -U myuser -d exocalc -v ON_ERROR_STOP=1 -f - < '/Users/erikray/Developer/develop-internet-applications/insert_stars.sql'
	docker compose exec -T postgres psql -U myuser -d exocalc -v ON_ERROR_STOP=1 -f - < '/Users/erikray/Developer/develop-internet-applications/insert.sql'
