MIGRATIONS_DIR=./migrations
DB_DSN=host=localhost port=5432 user=user password=password dbname=products sslmode=disable
LOCAL_BIN:=$(CURDIR)/bin

.PHONY: validator 
validator:
	go build -o bin/validator cmd/validator/main.go && ./bin/validator

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

.PHONY: generateMock
generateMock:
	~/go/bin/mockgen -source=./internal/pkg/repository/repository.go -destination=./internal/pkg/repository/mocks/repository.go 

.PHONY: cover
cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

.PHONY: test
test:
	go test -v count=1 ./...

.PHONY: test100
test100:
	go test -v count=100 ./...

.PHONY: race
race:
	go test -v -race -count=1 ./...