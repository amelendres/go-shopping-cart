#.PHONY: help

help:
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

## DOCKER
CONTAINER_NAME=go-shopping-cart

build: ## build and up docker containers
	@docker-compose up --build -d

start: ## run cart server
	@docker-compose up -d

make sh:
	@docker exec -it $(CONTAINER_NAME) sh

##TEST
test: ## run tests
	@go test ./...
	#go test ./pkg -run TestCartAddProduct -v
	#go test ./pkg/server  -v

coverage:
	mkdir -p var/test_results
	@go test -coverprofile=var/test_results/coverage.out ./...
	@go tool cover -func=var/test_results/coverage.out



##gRPC
gplugins: ##grpc plugins
	go get google.golang.org/protobuf/cmd/protoc-gen-go
	go get google.golang.org/grpc/cmd/protoc-gen-go-grpc


PROTO_PATH=proto
PATH_TYPE=source_relative
PROTO_OUT=.

gproto: ##compile proto files
	protoc --go_out=$(PROTO_OUT) --go_opt=paths=$(PATH_TYPE) \
        --go-grpc_out=require_unimplemented_servers=false:$(PROTO_OUT) --go-grpc_opt=paths=$(PATH_TYPE) \
        $(PROTO_PATH)/*.proto

cproto: ##clean pb
	rm -f $(PROTO_PATH)/*.pb.go

