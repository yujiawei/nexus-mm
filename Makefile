.PHONY: build run dev deps docker-up docker-down migrate clean

BINARY=nexus-mm
CONFIG=configs/nexus.yaml

build:
	go build -o bin/$(BINARY) ./cmd/server

run: build
	./bin/$(BINARY) -config $(CONFIG)

dev:
	go run ./cmd/server -config $(CONFIG)

deps:
	go mod tidy

docker-up:
	docker compose up -d

docker-down:
	docker compose down

migrate:
	@echo "Migrations are auto-applied by PostgreSQL init scripts in docker-compose"
	@echo "For manual apply: psql -h localhost -U nexus -d nexus_mm -f migrations/001_init.sql"

clean:
	rm -rf bin/

lint:
	golangci-lint run ./...

test:
	go test ./... -v
