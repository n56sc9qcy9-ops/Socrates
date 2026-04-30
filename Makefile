.PHONY: all build create fmt fmt-check lint test usage run clean help

APP := socrates
CMD := ./cmd/socrates
BIN_DIR := bin
BIN := $(BIN_DIR)/$(APP)
GO_PACKAGES := ./...

all: fmt lint test build usage

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN) $(CMD)

create: build

fmt:
	gofmt -w .

fmt-check:
	@test -z "$$(gofmt -l .)" || (gofmt -l . && exit 1)

lint:
	go vet $(GO_PACKAGES)

test:
	go test $(GO_PACKAGES)

usage:
	go run $(CMD) help

run: usage

clean:
	rm -rf $(BIN_DIR)

help:
	@echo "Available targets:"
	@echo "  make          Format, lint, test, build, and show CLI usage"
	@echo "  make create   Build ./$(BIN)"
	@echo "  make fmt      Format Go files with gofmt"
	@echo "  make lint     Run go vet"
	@echo "  make test     Run all Go tests"
	@echo "  make run      Show CLI usage"
	@echo "  make clean    Remove build output"
