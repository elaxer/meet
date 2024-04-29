include .env
.PHONY: create-env build up down migrate-up dictionary-load migrate-down fixture-load tests 

create-env:
	touch .env
	cat .env.dist > .env

build:
	mkdir -p ./database
	mkdir -p ./redis/local/data
	mkdir -p ./redis/local/redis.conf

	go mod download
	sudo docker-compose build

up:
	sudo docker-compose up -d
down:
	sudo docker-compose down

migrate-up:
	migrate -path migrations -database "${DB_DRIVER_NAME}://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}" up
dictionary-load:
	go run cmd/dictionary/main.go
migrate-down:
	migrate -path migrations -database "${DB_DRIVER_NAME}://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}" down
fixture-load:
	go run cmd/fixture/main.go
tests:
	go test -timeout 30s ./internal/...