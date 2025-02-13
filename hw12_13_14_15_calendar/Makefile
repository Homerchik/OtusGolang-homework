CALENDAR_BIN := "./bin/calendar"
SCHEDULER_BIN := "./bin/scheduler"
SENDER_BIN := "./bin/sender"

CALENDAR_IMG="calendar:develop"
SCHEDULER_IMG="scheduler:develop"
SENDER_IMG="sender:develop"
MIGRATOR_IMG="migrator:develop"



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

build-calendar-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" --build-arg SERVICE_NAME=calendar \
		-t $(CALENDAR_IMG) \
		-f build/Dockerfile .

build-scheduler-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" --build-arg=SERVICE_NAME=scheduler \
		-t $(SCHEDULER_IMG) \
		-f build/Dockerfile .

build-sender-img:
	docker build \
		--build-arg LDFLAGS="$(LDFLAGS)" --build-arg=SERVICE_NAME=sender \
		-t $(SENDER_IMG) \
		-f build/Dockerfile .

build-migrator-img:
	docker build \
		-t $(MIGRATOR_IMG) \
		-f build/Dockerfile.migrator .

build-img: build-calendar-img build-scheduler-img build-sender-img build-migrator-img

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

.PHONY: build run-calendar build-scheduler-img build-calendar-img build-sender-img build-img run-img version test lint generate run-scheduler run-sender build-migrator-img
