PROJECT_NAME=cloud-back
APPLICATION_DIR := $(shell pwd)/application

.PHONY: tidy run-cloud-api run-cloud-rpc build-cloud-api build-cloud-rpc

# Sync module dependencies.
tidy:
	@echo "updating go dependencies..."
	@go mod tidy

# Run API service with local config.
run-cloud-api:
	@echo "starting cloud-api..."
	@go run $(APPLICATION_DIR)/cloud-api/cloud.go -f $(APPLICATION_DIR)/cloud-api/etc/cloud-api.yaml

# Run RPC service with local config.
run-cloud-rpc:
	@echo "starting cloud-rpc..."
	@go run $(APPLICATION_DIR)/cloud-rpc/cloud.go -f $(APPLICATION_DIR)/cloud-rpc/etc/cloud-rpc.yaml

# Build API binary into dist/.
build-cloud-api:
	@mkdir -p dist
	@go build -o dist/cloud-api $(APPLICATION_DIR)/cloud-api/cloud.go

# Build RPC binary into dist/.
build-cloud-rpc:
	@mkdir -p dist
	@go build -o dist/cloud-rpc $(APPLICATION_DIR)/cloud-rpc/cloud.go
