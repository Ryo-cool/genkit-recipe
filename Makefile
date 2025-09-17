SHELL := /bin/bash

ROOT_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
BACKEND_DIR := $(ROOT_DIR)/backend
FRONTEND_DIR := $(ROOT_DIR)/frontend

.PHONY: help setup backend-run backend-test backend-dev-ui frontend-dev frontend-build frontend-lint clean

help:
	@echo "Available targets:"
	@echo "  setup            Install backend and frontend dependencies"
	@echo "  backend-run      Run the Go recipe flow server"
	@echo "  backend-test     Execute Go unit tests"
	@echo "  backend-dev-ui   Launch Genkit Dev UI with the Go flow"
	@echo "  frontend-dev     Start the Next.js development server"
	@echo "  frontend-build   Build the Next.js app for production"
	@echo "  frontend-lint    Run ESLint checks"
	@echo "  clean            Remove build artifacts (frontend/.next)"

setup:
	cd $(BACKEND_DIR) && go mod tidy
	cd $(FRONTEND_DIR) && npm install

backend-run:
	cd $(BACKEND_DIR) && go run ./cmd/recipe

backend-test:
	cd $(BACKEND_DIR) && go test ./...

backend-dev-ui:
	cd $(BACKEND_DIR) && genkit start -- go run ./cmd/recipe

frontend-dev:
	cd $(FRONTEND_DIR) && npm run dev

frontend-build:
	cd $(FRONTEND_DIR) && npm run build

frontend-lint:
	cd $(FRONTEND_DIR) && npm run lint

clean:
	rm -rf $(FRONTEND_DIR)/.next
