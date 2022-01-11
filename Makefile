.PHONY: default run fmt help

default: help

run: ## Run the program
	go run server.go

build: Dockerfile ## Build container image
	@docker build . -t aeolyus/gull

release: Dockerfile ## Build and push the container image for all platforms
	@docker buildx build . \
		-t aeolyus/gull \
		-t ghcr.io/aeolyus/gull \
		--platform linux/arm,linux/arm64,linux/amd64 \
		--push

fmt: ## Format the code
	go fmt ./...

test: ## Run tests
	go test -cover ./...

help: Makefile ## Print this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; \
		{printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
