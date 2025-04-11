# Project Configuration
TARGET := gmeter
BUILD_DIR := build
PKG_DIRS := $(wildcard pkg/*)  # Auto-detects all pkg/ subdirectories
CMD_DIR := cmd/gmeter

# Tools
CLANG ?= clang
GO ?= go

.PHONY: all build generate clean run deps fmt help

# Default target
all: generate build

# Generate eBPF bindings for ALL packages
generate:
	@echo "Generating eBPF bindings..."
	@for pkg in $(PKG_DIRS); do \
		if [ -f "$$pkg/gen.go" ]; then \
			echo "Generating $$pkg..."; \
			(cd "$$pkg" && $(GO) generate); \
		fi; \
	done

# Build main executable
build: generate
	@echo "Building $(TARGET)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/$(TARGET) ./$(CMD_DIR)

# Clean ALL generated files
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@find pkg/ -name '*_bpf*.go' -delete
	@find pkg/ -name '*.o' -delete
	@echo "Preserved: $$(find pkg/ -name gen.go)"

# Run with default interface
run: build
	@echo "Running $(TARGET)..."
	@sudo $(BUILD_DIR)/$(TARGET) -i eth0

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@sudo apt-get update && sudo apt-get install -y clang llvm libbpf-dev

# Format all source code (excluding generated files)
fmt:
	@echo "Formatting code..."
	@find . -name '*.go' -not -name '*_bpf*.go' -exec gofmt -w {} \;
	@find . -name '*.c' -exec clang-format -i {} \;

# Help
help:
	@echo "Available targets:"
	@echo "  all       - Generate bindings and build (default)"
	@echo "  generate  - Generate eBPF bindings for all packages"
	@echo "  build     - Build the executable (outputs to build/)"
	@echo "  clean     - Remove ALL generated files"
	@echo "  run       - Build and run with default interface"
	@echo "  deps      - Install build dependencies"
	@echo "  fmt       - Format source code (excludes generated files)"