.PHONY: proto proto-install help ingestion-grpc ingestion-http

# Root makefile for YouTube Analytics project

# Proto generation commands
PROTO_DIR = proto
PROTO_OUT_DIR = services/pkg/pb

# Help command
help:
	@echo "YouTube Analytics - Available commands:"
	@echo ""
	@echo "Proto generation:"
	@echo "  make proto-install  - Install required tools for proto generation"
	@echo "  make proto         - Generate Go code from proto files"
	@echo ""
	@echo "Ingestion Service:"
	@echo "  make ingestion-grpc - Run ingestion service gRPC server"
	@echo "  make ingestion-http - Run ingestion service HTTP server"
	@echo ""
	@echo "Testing:"
	@echo "  make test          - Run all tests"
	@echo "  make lint          - Run linter on all services"

# Install required tools for proto generation
proto-install:
	@echo "Installing protoc-gen-go and protoc-gen-go-grpc..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Generate Go code from proto files
proto:
	@echo "Generating Go code from proto files..."
	@mkdir -p $(PROTO_OUT_DIR)/ingestion/v1
	protoc \
		--go_out=$(PROTO_OUT_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT_DIR) \
		--go-grpc_opt=paths=source_relative \
		-I $(PROTO_DIR) \
		$(PROTO_DIR)/ingestion/v1/ingestion.proto

# Run ingestion service gRPC server
ingestion-grpc:
	cd services/ingestion-service && go run cmd/grpc/main.go

# Run ingestion service HTTP server
ingestion-http:
	cd services/ingestion-service && go run cmd/http/main.go

# Run tests
test:
	@echo "Running tests..."
	cd services/ingestion-service && go test -v ./...

# Run linter
lint:
	@echo "Running linter..."
	cd services/ingestion-service && golangci-lint run