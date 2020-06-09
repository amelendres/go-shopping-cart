
build: ## up docker containers
	@docker-compose up --build -d

start: ## run cart server
	@docker-compose up -d

test: ## run tests
	@go test

coverage:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out

