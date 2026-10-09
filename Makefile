SHELL := /bin/sh
.DEFAULT_GOAL := help
.PHONY: help doctor install-system-deps setup setup-go setup-frontend setup-wails test test-go test-frontend build build-linux dev clean

GO ?= go
NPM ?= npm
GOPATH := $(shell $(GO) env GOPATH 2>/dev/null)
WAILS ?= $(GOPATH)/bin/wails
BIN_DIR ?= build/bin
APP ?= agentforge

help: ## Show available commands
	@awk 'BEGIN {FS = ":.*## "} /^[a-zA-Z0-9_-]+:.*## / {printf "  %-22s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

doctor: ## Check build tools and Linux desktop libraries
	@set -eu; \
	for cmd in $(GO) $(NPM) pkg-config; do \
	  if ! command -v "$$cmd" >/dev/null 2>&1; then echo "Missing: $$cmd"; exit 1; fi; \
	done; \
	if [ "$$(uname -s)" != Linux ]; then echo "This Makefile currently supports Linux builds only"; exit 1; fi; \
	pkg-config --exists gtk+-3.0 || { echo "Missing GTK3 development libraries"; exit 1; }; \
	pkg-config --exists webkit2gtk-4.1 || { echo "Missing WebKit2GTK 4.1 development libraries"; exit 1; }; \
	echo "Go: $$($(GO) version)"; node --version; $(NPM) --version; \
	echo "Linux desktop build prerequisites: OK"

install-system-deps: ## Install Linux dependencies (explicit sudo, Arch/CachyOS or Debian/Ubuntu)
	@set -eu; \
	if command -v pacman >/dev/null 2>&1; then \
	  sudo pacman -S --needed --noconfirm base-devel go nodejs npm pkgconf gtk3 webkit2gtk-4.1; \
	elif command -v apt-get >/dev/null 2>&1; then \
	  sudo apt-get update; \
	  sudo apt-get install -y build-essential golang-go nodejs npm pkg-config libgtk-3-dev libwebkit2gtk-4.1-dev; \
	else \
	  echo "Unsupported package manager. Install Go, Node/npm, GTK3 and WebKit2GTK 4.1 manually."; exit 1; \
	fi

setup-go: ## Download Go modules and generate go.sum
	$(GO) mod tidy

setup-frontend: ## Install frontend dependencies
	@if [ -f frontend/package-lock.json ]; then \
	  cd frontend && $(NPM) ci; \
	else \
	  cd frontend && $(NPM) install; \
	fi

setup-wails: ## Install Wails v2 CLI in GOPATH/bin
	$(GO) install github.com/wailsapp/wails/v2/cmd/wails@v2.10.2

setup: doctor setup-go setup-frontend setup-wails ## Set up all application dependencies

test-go: ## Run Go unit tests (no desktop runtime required)
	$(GO) test ./internal/...

test-frontend: ## Run Vue/TypeScript checks
	cd frontend && $(NPM) run typecheck

test: test-go test-frontend ## Run backend and frontend checks

build-linux: doctor setup-go setup-frontend ## Compile Linux AMD64/native architecture binary
	cd frontend && $(NPM) run build
	mkdir -p $(BIN_DIR)
	$(GO) build -tags webkit2_41 -trimpath -o $(BIN_DIR)/$(APP) .
	@echo "Built: $(BIN_DIR)/$(APP)"

build: build-linux ## Alias for Linux build

dev: doctor setup-go setup-frontend setup-wails ## Start Wails development mode
	$(WAILS) dev -tags webkit2_41

clean: ## Remove generated frontend and desktop artifacts
	rm -rf $(BIN_DIR) frontend/dist
