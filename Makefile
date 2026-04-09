PROJECT_NAME=cloud-back
APPLICATION_DIR := $(shell pwd)/application

.PHONY: tidy run-portal-api run-manager-api run-console-api run-cloud-api run-cloud-rpc build-portal-api build-manager-api build-console-api build-cloud-api build-cloud-rpc test

tidy:
	@echo "updating go dependencies..."
	@go mod tidy

run-portal-api:
	@echo "starting portal-api..."
	@go run $(APPLICATION_DIR)/portal-api/portal.go -f $(APPLICATION_DIR)/portal-api/etc/portal-api.yaml

run-manager-api:
	@echo "starting manager-api..."
	@go run $(APPLICATION_DIR)/manager-api/manager.go -f $(APPLICATION_DIR)/manager-api/etc/manager-api.yaml

run-console-api:
	@echo "starting console-api..."
	@go run $(APPLICATION_DIR)/console-api/console.go -f $(APPLICATION_DIR)/console-api/etc/console-api.yaml

# backward-compatible targets
run-cloud-api:
	@echo "starting cloud-api (legacy scaffold)..."
	@go run $(APPLICATION_DIR)/cloud-api/cloud.go -f $(APPLICATION_DIR)/cloud-api/etc/cloud-api.yaml

run-cloud-rpc:
	@echo "starting cloud-rpc (legacy scaffold)..."
	@go run $(APPLICATION_DIR)/cloud-rpc/cloud.go -f $(APPLICATION_DIR)/cloud-rpc/etc/cloud-rpc.yaml

build-portal-api:
	@mkdir -p dist
	@go build -o dist/portal-api $(APPLICATION_DIR)/portal-api/portal.go

build-manager-api:
	@mkdir -p dist
	@go build -o dist/manager-api $(APPLICATION_DIR)/manager-api/manager.go

build-console-api:
	@mkdir -p dist
	@go build -o dist/console-api $(APPLICATION_DIR)/console-api/console.go

build-cloud-api:
	@mkdir -p dist
	@go build -o dist/cloud-api $(APPLICATION_DIR)/cloud-api/cloud.go

build-cloud-rpc:
	@mkdir -p dist
	@go build -o dist/cloud-rpc $(APPLICATION_DIR)/cloud-rpc/cloud.go

test:
	@go test ./...
