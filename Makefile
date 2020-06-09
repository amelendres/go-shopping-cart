
build: ## up docker containers
	@docker-compose build
	@docker-compose up -d

start: ## run cart server
	@docker-compose go run ./cmd/webserver/main.go

test: ## run tests
	@docker exec -it appto-go-shopping-cart go test

coverage:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out

