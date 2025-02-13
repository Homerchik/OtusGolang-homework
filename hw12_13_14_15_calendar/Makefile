CALENDAR_BIN := "./bin/calendar"
SCHEDULER_BIN := "./bin/scheduler"
SENDER_BIN := "./bin/sender"

DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(CALENDAR_BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(SCHEDULER_BIN) -ldflags "$(LDFLAGS)" ./cmd/scheduler
	go build -v -o $(SENDER_BIN) -ldflags "$(LDFLAGS)" ./cmd/sender

run: build
	$(CALENDAR_BIN) -config ./configs/config.yaml

run-scheduler: build
	$(SCHEDULER_BIN) -config ./configs/config.yaml

run-sender: build
	$(SENDER_BIN) -config ./configs/config.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(CALENDAR_BIN) version

test:
	go test -race ./internal/...

migrate-up:
	goose postgres "user=${DB_USERNAME} password=${DB_PASSWORD} dbname=${DB_NAME} host=${DB_HOST} port=${DB_PORT}" -dir ${DB_MIGRATIONS} up

migrate-down:
	goose postgres "user=${DB_USERNAME} password=${DB_PASSWORD} dbname=${DB_NAME} host=${DB_HOST} port=${DB_PORT}" -dir ${DB_MIGRATIONS} down

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.61.0

lint: install-lint-deps
	golangci-lint run ./... -v

generate:
	rm -rf ./internal/server/protobuf/pb
	mkdir -p ./internal/server/protobuf/pb

	protoc \
		--go_out=./internal/server/protobuf \
		--go-grpc_out=./internal/server/protobuf \
		./api/*.proto

.PHONY: build run-calendar build-img run-img version test lint generate run-scheduler run-sender
