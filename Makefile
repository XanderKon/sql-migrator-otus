GOBIN ?= $$(go env GOPATH)/bin
BIN := "./bin/gomigrator"
DOCKER_IMG := "gomigrator-image"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

# Docker & docker-compose stuff
D_DSN := "postgresql://postgres:postgres@db:5432/gomigrator?sslmode=disable"

build:
	go build -tags "postgres" -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/gomigrator

build-for-docker:
	CGO_ENABLED=0 go build -tags "postgres" -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/gomigrator

run: build
	$(BIN) -config ./configs/config.yml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		--build-arg=BIN="$(BIN)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

run-compose:
	docker compose -f "./build/docker-compose.yml" up -d --no-deps --build

# Just for demo
run-compose-demo:
	docker compose -f "./build/docker-compose.yml" up -d --no-deps --build
	docker compose -f "./build/docker-compose.yml" exec -it app gomigrator -dsn=${D_DSN} status && \
	docker compose -f "./build/docker-compose.yml" exec -it app gomigrator -dsn=${D_DSN} up && \
	docker compose -f "./build/docker-compose.yml" exec -it app gomigrator -dsn=${D_DSN} dbversion && \
	docker compose -f "./build/docker-compose.yml" exec -it app gomigrator -dsn=${D_DSN} status && \
	docker compose -f "./build/docker-compose.yml" exec -it app gomigrator -dsn=${D_DSN} down && \
	docker compose -f "./build/docker-compose.yml" exec -it app gomigrator -dsn=${D_DSN} status && \
	docker compose -f "./build/docker-compose.yml" exec -it app gomigrator -dsn=${D_DSN} redo && \
	docker compose -f "./build/docker-compose.yml" exec -it app gomigrator -dsn=${D_DSN} status && \
	docker compose -f "./build/docker-compose.yml" exec -it app gomigrator -dsn=${D_DSN} create test_new_migration && \
	docker compose -f "./build/docker-compose.yml" exec -it app gomigrator -dsn=${D_DSN} up && \
	docker compose -f "./build/docker-compose.yml" exec -it app gomigrator -dsn=${D_DSN} status && \
	docker compose -f "./build/docker-compose.yml" down

down-compose:
	docker compose -f "./build/docker-compose.yml" down

version: build
	$(BIN) version

test:
	go test -race ./...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.55.2

integration-tests:
	cd test && DSN=postgresql://postgres:postgres@localhost:5432/gomigrator DIR=./migrations go test -tags integration

lint: install-lint-deps
	golangci-lint run ./...

install-go-test-coverage:
	go install github.com/vladopajic/go-test-coverage/v2@latest

check-coverage: install-go-test-coverage
	go test ./internal/... ./pkg/... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...
	${GOBIN}/go-test-coverage --config=./.testcoverage.yml

.PHONY: build build-for-docker run build-img run-img run-compose run-compose-demo down-compose version test lint install-lint-deps integration_test install-go-test-coverage check-coverage
