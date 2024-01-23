GOBIN ?= $$(go env GOPATH)/bin
BIN := "./bin/gomigrator"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -tags "postgres" -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/gomigrator

run: build
	$(BIN) -config ./configs/config.yml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.55.2

integration_test:
	DSN=postgresql://postgres:postgres@localhost:5432/gomigrator DIR=./migrations go test -tags integration

lint: install-lint-deps
	golangci-lint run ./...

install-go-test-coverage:
	go install github.com/vladopajic/go-test-coverage/v2@latest

check-coverage: install-go-test-coverage
	go test ./internal/... ./pkg/... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...
	${GOBIN}/go-test-coverage --config=./.testcoverage.yml

.PHONY: build run build-img run-img version test lint install-lint-deps integration_test install-go-test-coverage check-coverage
