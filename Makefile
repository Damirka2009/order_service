# имя бинарника
BINARY=server

# пути
CMD_DIR=cmd/server
PROTO_DIR=proto
OUT_DIR=.

# переменные окружения
ENV_FILE=.env

.PHONY: proto build run tidy clean

# генерация gRPC кода
proto:
	protoc \
	--go_out=$(OUT_DIR) \
	--go-grpc_out=$(OUT_DIR) \
	--grpc-gateway_out= \
	$(PROTO_DIR)/*.proto

# установить зависимости
tidy:
	go mod tidy

# сборка бинарника
build:
	mkdir -p bin
	go build -o bin/$(BINARY) $(CMD_DIR)/main.go

# запуск сервера
run:
	go run $(CMD_DIR)/main.go

# помощь 
help:
	@echo "make proto   - generate gRPC code"
	@echo "make build   - build binary"
	@echo "make run     - run server"