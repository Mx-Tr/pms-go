include .env

.PHONY: migrate-up migrate-down migrate-create

migrate-up:
	migrate -path ./migrations -database "$(DB_URL)?sslmode=disable" up

migrate-down:
	migrate -path ./migrations -database "$(DB_URL)?sslmode=disable" down $(filter-out $@,$(MAKECMDGOALS))

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir ./migrations -seq $$name

#Для грязных миграций
migrate-force:
	@read -p "Enter version to force: " version; \
	migrate -path ./migrations -database "$(DB_URL)?sslmode=disable" force $$version

migrate-version:
	migrate -path ./migrations -database "$(DB_URL)?sslmode=disable" version

run:
	go run ./cmd/api

build:
	go run build -o bin/api ./cmd/api