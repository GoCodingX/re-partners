TOOLS_PATH := $(CURDIR)/.tools
LINTER_BINARY := $(TOOLS_PATH)/golangci-lint

# install golangci-lint linter checks
.PHONY: install-linter
install-linter:
	@if [ ! -f $(LINTER_BINARY) ]; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(TOOLS_PATH) v1.63.4; \
	fi

# runs golangci-lint linter checks
.PHONY: lint
lint: install-linter
	@$(LINTER_BINARY) run --fix

# recreate the vendors directory and tidy dependencies
.PHONY: mod-clean
mod-clean:
	@rm -rf vendor && go mod tidy && go mod vendor

# generate code: mocks, openapi, etc
.PHONY: generate
generate:
	@go generate ./...

# run the db in docker, suitable for local development, when
# a db is needed to properly run the application locally
.PHONY: run-db
run-db:
	@docker compose up -d db

dev: dev-docker

# run the app (suitable for development)
.PHONY: dev-local
dev-local:
	@go run cmd/packs-api/main.go

# run the app using docker (suitable for development)
.PHONY: dev-docker
dev-docker:
	@DATABASE_URL=postgres://user:password@db:5432/packs?sslmode=disable \
	MIGRATIONS_DIR=./internal/repository/pg/migrations \
	docker compose up app

test: test-docker

# run unit tests using docker
.PHONY: test-docker
test-docker:
	@docker compose run --rm unit-tests

# run unit tests
.PHONY: test-local
test-local:
	@go test ./...
