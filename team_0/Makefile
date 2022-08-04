SERVER_OUT := transmitter
CLIENT_OUT := receiver
API_PATH	= ./internal/transport
PKG := "gitlab.com/pantomath-io/demo-grpc"

.PHONY: all api build_server build_client

all: build_server build_client

api: internal/transport/api.pb.go internal/transport/api_grpc.pb.go ## Auto-generate grpc go sources

internal/transport/api.pb.go: internal/transport/api.proto
	@protoc --go_out=. $(API_PATH)/api.proto

internal/transport/api_grpc.pb.go: internal/transport/api.proto
	@protoc --go-grpc_out=. $(API_PATH)/api.proto

dep: ## Get the dependencies
	go get -v -d ./...

build_server: dep api ## Build the binary file for server
	go build -o $(SERVER_OUT) ./cmd/$(SERVER_OUT)/main.go

build_client: dep api ## Build the binary file for client
	go build -o $(CLIENT_OUT) ./cmd/$(CLIENT_OUT)/main.go

clean: ## Remove previous builds
	@rm $(SERVER_OUT) $(CLIENT_OUT) $(API_OUT)

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
