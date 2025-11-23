.PHONY: build install uninstall clean test run fmt

BINARY_NAME=clash-fish
INSTALL_PATH=/usr/local/bin
BUILD_DIR=build

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/clash-fish

install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_PATH)/
	sudo chmod +x $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "✓ Installation complete"

uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "✓ Uninstallation complete"

clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	go clean
	@echo "✓ Clean complete"

test:
	@echo "Running tests..."
	go test -v ./...

fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "✓ Format complete"

run: build
	@echo "Running $(BINARY_NAME)..."
	sudo $(BUILD_DIR)/$(BINARY_NAME) start

help:
	@echo "Available targets:"
	@echo "  build      - Build the binary"
	@echo "  install    - Install to system"
	@echo "  uninstall  - Remove from system"
	@echo "  clean      - Clean build artifacts"
	@echo "  test       - Run tests"
	@echo "  fmt        - Format code"
	@echo "  run        - Build and run"
