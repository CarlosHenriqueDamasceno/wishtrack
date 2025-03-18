include .env

BIN_PATH=etc/bin/
MIGRATIONS_PATH = ./etc/migrations

.PHONY: build
build:
	@go build -o $(BIN_PATH)/api cmd/api/main.go

.PHONY: run
run: build
	@./$(BIN_PATH)/api

.PHONY: tests
tests:
	@go test ./test/... -v

.PHONY: migrate-create
migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_DRIVER)://$(DB_DSN) up

.PHONY: migrate-down
migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_DRIVER)://$(DB_DSN) down $(filter-out $@,$(MAKECMDGOALS))

.PHONY: gen-docs
gen-docs:
	@swag init -o ./etc/doc -g ./api/main.go -d cmd,internal,pkg && swag fmt

.PHONY: front
front:
	@cd web && bun dev --host

.PHONY: back
back:
	@air

.PHONY: dev
dev:
	@docker compose up -d && make -j back front