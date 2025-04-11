# Project configuration
TARGET := gmeter
EBPF_SRC := ebpf/counter.c
GENERATED_DIR := gen
CMD_DIR := cmd/gmeter

# Tools
CLANG ?= clang
GO ?= go

.PHONY: all build generate clean run

all: generate build

# Generate eBPF bindings
generate:
	@echo "Generating eBPF bindings..."
	@cd $(GENERATED_DIR) && $(GO) generate

# Build the main executable
build: generate
	@echo "Building $(TARGET)..."
	@cd $(CMD_DIR) && $(GO) build -o $(TARGET) .

# Clean all gen files (preserves gen.go)
clean:
	@echo "Cleaning up generated files..."
	@rm -f $(GENERATED_DIR)/counter_bpf*.go $(GENERATED_DIR)/counter_bpf*.o
	@cd $(CMD_DIR) && rm -f $(TARGET)
	@echo "Note: Preserved $(GENERATED_DIR)/gen.go"

# Run with default interface (eth0)
run: build
	@echo "Running $(TARGET)..."
	@sudo $(CMD_DIR)/$(TARGET) -i eth0

# Install dependencies (clang, libbpf, etc.)
deps:
	@echo "Installing dependencies..."
	@sudo apt-get update && sudo apt-get install -y clang llvm libbpf-dev

# Format code
fmt:
	@find . -name '*.go' -not -path "./$(GENERATED_DIR)/*" -exec gofmt -w {} \;
	@find . -name '*.c' -exec clang-format -i {} \;

# Help target
help:
	@echo "Available targets:"
	@echo "  all       - Generate bindings and build (default)"
	@echo "  generate  - Generate eBPF bindings"
	@echo "  build     - Build the executable"
	@echo "  clean     - Remove generated files (preserves gen.go)"
	@echo "  run       - Build and run with default interface"
	@echo "  deps      - Install build dependencies"
	@echo "  fmt       - Format source code (excludes generated files)"