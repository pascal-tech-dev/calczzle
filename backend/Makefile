.PHONY: help fmt vet test cover cover-html check run tidy

help: ## Show available targets
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-12s %s\n", $$1, $$2}'

fmt: ## Format all Go packages
	go fmt ./...

vet: ## Run go vet on all packages
	go vet ./...

test: ## Run all tests
	go test ./...

cover: ## Run tests with coverage; print total and per-package summary
	go test ./... -coverprofile=coverage.out -covermode=atomic
	@echo
	@go tool cover -func=coverage.out | grep -v '^total:'
	@echo
	@go tool cover -func=coverage.out | grep '^total:'

cover-html: cover ## Generate HTML coverage report (coverage.html)
	go tool cover -html=coverage.out -o coverage.html
	@echo "Wrote coverage.html"

check: fmt vet test ## Format, vet, and test

run: ## Start the API server
	go run ./cmd/api

tidy: ## Tidy go.mod / go.sum
	go mod tidy
