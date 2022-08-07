MIGRATIONS_DIR=./migrations
DB_DSN=host=localhost port=5432 user=user password=password dbname=products sslmode=disable
LOCAL_BIN:=$(CURDIR)/bin

.PHONY: run
run:
	cp pkg/api/product.swagger.json bin/
	go build -o bin/bot cmd/bot/main.go && ./bin/bot
build:
	go build -o bin/bot cmd/bot/main.go
proto:
	buf mod update
	buf generate api

.PHONY: addMigration
addMigration:
	goose -dir=${MIGRATIONS_DIR} create $(NAME) sql 

.PHONY: migrate
migrate:
	goose -v -dir=${MIGRATIONS_DIR} postgres "${DB_DSN}" up 

.PHONY: .deps
.deps:
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway && \
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 && \
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go && \
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc 
