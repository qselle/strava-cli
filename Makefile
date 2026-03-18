APP_NAME := strava-cli
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -s -w -X main.version=$(VERSION)

.PHONY: build install lint test clean

build:
	go build -ldflags "$(LDFLAGS)" -o bin/$(APP_NAME) .

install:
	go install -ldflags "$(LDFLAGS)" .

lint:
	golangci-lint run

test:
	go test ./...

clean:
	rm -rf bin/ dist/

.DEFAULT_GOAL := build
