# Go uygulaması için Makefile

.PHONY: run build test lint docker-up docker-down migrate

run:
	go run ./cmd/main.go

build:
	go build -o hospital-app ./cmd/main.go

test:
	go test ./...

lint:
	golangci-lint run || true

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down

migrate:
	go run ./cmd/main.go migrate 