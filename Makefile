GO ?= go
BINARY_NAME ?= bin/terminal
CMD_DIR ?= ./cmd/terminal
# GLFW picks its display backend at compile time; empty TAGS builds the X11 backend.
TAGS ?= wayland
GOFLAGS_TAGS := $(if $(TAGS),-tags $(TAGS),)

.PHONY: all build run test clean help

all: build

build: ## Build the terminal application binary
	@mkdir -p bin
	$(GO) build $(GOFLAGS_TAGS) -o $(BINARY_NAME) $(CMD_DIR)

run: ## Run the terminal application directly
	$(GO) run $(GOFLAGS_TAGS) $(CMD_DIR)

test: ## Run unit tests across all packages
	$(GO) test -v ./...

clean: ## Remove built binaries and output artifacts
	rm -rf bin/

help: ## Display available make targets
	@echo "Usage:"
	@echo "  make <target>"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}'
