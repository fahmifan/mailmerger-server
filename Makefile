test:
	go test -race ./...

run-server:
	modd -f server.modd.conf

run-dev-postgres:
	docker compose -f docker-compose.dev.yml up -d postgres

run-dev-mailhog:
	docker compose -f docker-compose.dev.yml up -d mailhog

migrate-up-postgres:
	sql-migrate up -env=postgres

migrate-up-sqlite:
	sql-migrate up -env=sqlite
