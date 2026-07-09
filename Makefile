SHELL := /usr/bin/env bash

BINARY_NAME ?= devspace
CMD_PATH ?= ./cmd/devspace
BIN_DIR ?= bin
DIST_DIR ?= dist
GOLANGCI_LINT ?= go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2
GOVULNCHECK ?= go run golang.org/x/vuln/cmd/govulncheck@v1.1.4

.PHONY: \
	all fmt format fmt-check test test-race test-one cover cover-html vet lint vulncheck build verify command-surface-check install run \
	precommit install-hooks clean \
	tui-install tui-verify tui-build tui-build-all tui-dev tui-install-local \
	snapshot ci \
	help

.DEFAULT_GOAL := help

##@ Go

all: verify ## Run all checks and build the binary

fmt format: ## Format the code
	gofmt -w cmd internal

fmt-check: ## Check that the code is gofmt-formatted
	test -z "$$(gofmt -l cmd internal)" || (gofmt -l cmd internal && exit 1)

test: ## Run the tests
	go test ./...

test-race: ## Run the tests with the race detector
	go test ./... -race

test-one: ## Run a single test: make test-one T=TestName
	@test -n "$(T)" || (echo "Usage: make test-one T=TestName" && exit 1)
	go test ./... -run '^$(T)$$' -v

cover: ## Run tests with a coverage profile and print a summary
	mkdir -p $(BIN_DIR)
	go test ./... -coverprofile=$(BIN_DIR)/coverage.out -covermode=atomic
	go tool cover -func=$(BIN_DIR)/coverage.out

cover-html: cover ## Generate an HTML coverage report (prints the path; does not open a browser)
	go tool cover -html=$(BIN_DIR)/coverage.out -o $(BIN_DIR)/coverage.html
	@echo "Coverage report: $(BIN_DIR)/coverage.html"

vet: ## Run go vet
	go vet ./...

lint: fmt-check ## Run golangci-lint
	$(GOLANGCI_LINT) run ./...

vulncheck: ## Run govulncheck
	$(GOVULNCHECK) ./...

build: ## Build the binary
	mkdir -p $(BIN_DIR)
	go build -trimpath -o $(BIN_DIR)/$(BINARY_NAME) $(CMD_PATH)

install: ## go install the binary (into GOBIN/GOPATH bin)
	go install -trimpath $(CMD_PATH)

run: build ## Build then run the binary, e.g. make run ARGS="status --verbose"
	./$(BIN_DIR)/$(BINARY_NAME) $(ARGS)

command-surface-check: ## Check maintained docs and demos for removed CLI paths
	scripts/check-command-surface.sh
	scripts/check-command-surface-self-test.sh

verify: command-surface-check test vet lint vulncheck build ## Run all checks and build the binary (the local CI gate)

##@ Dev

precommit: fmt lint test build ## Run the pre-commit hook checks locally

install-hooks: ## Point git at .githooks (runs precommit checks on commit)
	git config core.hooksPath .githooks
	chmod +x .githooks/pre-commit

clean: ## Remove build, dist, and coverage artifacts
	rm -rf $(BIN_DIR) $(DIST_DIR)

##@ TUI

# devspace-tui (Bun) — not part of `verify` so Go-only work never needs Bun.
tui-install: ## Install devspace-tui (Bun) dependencies
	cd tui && bun install --frozen-lockfile

tui-verify: tui-install ## Typecheck and test devspace-tui
	cd tui && bun run typecheck && bun test

tui-build: tui-install ## Build devspace-tui for the current platform
	cd tui && bun run build

tui-build-all: tui-install ## Build devspace-tui for all release platforms
	cd tui && ./build-all.sh

tui-dev: build tui-install ## Run devspace-tui in dev mode against the local Go binary
	cd tui && DEVSPACE_BIN=../$(BIN_DIR)/$(BINARY_NAME) bun run dev

tui-install-local: tui-build ## Build devspace-tui and install it into $$DEVSPACE_HOME/bin (or ~/.devspace/bin)
	@dest="$${DEVSPACE_HOME:-$$HOME/.devspace}/bin"; \
	mkdir -p "$$dest"; \
	cp tui/dist/devspace-tui "$$dest/"; \
	echo "Installed devspace-tui to $$dest"

##@ Release

snapshot: ## Build a local snapshot release via goreleaser (no publish)
	@command -v goreleaser >/dev/null 2>&1 || (echo "goreleaser not found; see https://goreleaser.com/install/" && exit 1)
	goreleaser release --snapshot --clean --skip=publish

ci: verify tui-verify ## Run everything CI runs (Go verify + TUI verify)

##@ Help

help: ## Show this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage: make \033[36m<target>\033[0m\n"} \
		/^[a-zA-Z0-9_ -]+:.*##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } \
		/^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) }' $(MAKEFILE_LIST)
