.DEFAULT_GOAL := build

.PHONY: run-server
run-server: build-server
	./bin/server

.PHONY: run-client
run-client: build-client
	./bin/client

.PHONY: build-server
build-server:
	go build -o bin/server cmd/server/main.go cmd/server/server.go cmd/server/consumer.go

.PHONY: build-client
build-client:
	go build -o bin/client client/client.go

build: build-server build-client

.PHONY: .test
.test:
	$(info Running tests...)
	go test ./...
	
LOCAL_BIN:=$(PWD)/bin
.PHONY: .deps
.deps:
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway && \
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 && \
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go && \
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc && \
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@latest && \
	GOBIN=$(LOCAL_BIN) go install github.com/golang/mock/mockgen@v1.6.0 && \
	GOBIN=$(LOCAL_BIN) go install github.com/edbmanniwood/pgxpoolmock

MIGRATIONS_DIR:=./migrations
.PHONY: migration
migration:
	goose -dir=${MIGRATIONS_DIR} create $(NAME) sql

GOOSE_DIR:=./migrations
GOOSE_DBSTRING:=host=localhost port=6432 user=${POSTGRESQL_USER} password=${POSTGRESQL_PASS} dbname=${POSTGRESQL_DB} sslmode=disable
.PHONY: integration-tests
integration-tests:
	./bin/goose -dir $(GOOSE_DIR) postgres "$(GOOSE_DBSTRING)" up
	go test --tags=integration ./...