SHELL := /usr/bin/env bash

BINARY_NAME ?= devspace
CMD_PATH ?= ./cmd/devspace
BIN_DIR ?= bin
DIST_DIR ?= dist

.PHONY: all test vet build verify clean

all: verify

test:
	go test ./...

vet:
	go vet ./...

build:
	mkdir -p $(BIN_DIR)
	go build -trimpath -o $(BIN_DIR)/$(BINARY_NAME) $(CMD_PATH)

verify: test vet build

clean:
	rm -rf $(BIN_DIR) $(DIST_DIR)
