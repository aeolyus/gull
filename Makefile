.PHONY: default run fmt help

default: help

run: ## Run the program
	go run server.go

image: Dockerfile ## Build container image
	@docker buildx build . -t ghcr.io/aeolyus/gull

release: Dockerfile ## Build and push the container image for all platforms
	-@docker buildx rm gull-builder
	@docker buildx create --name gull-builder --bootstrap --use
	@docker buildx build . \
		-t aeolyus/gull \
		-t ghcr.io/aeolyus/gull \
		--platform linux/arm,linux/arm64,linux/amd64 \
		--push
	@docker buildx stop
	-@docker buildx rm gull-builder

fmt: ## Format the code
	go fmt ./...

test: ## Run tests
	go test -cover ./...

help: Makefile ## Print this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; \
		{printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
