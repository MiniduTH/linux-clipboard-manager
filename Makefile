# Makefile for Clipboard Manager

.PHONY: build test clean install help

# Default target
all: build

# Build the clipboard manager
build:
	@echo "🔨 Building Clipboard Manager..."
	@mkdir -p build
	go build -o clipboard-manager
	@cp clipboard-manager build/
	@echo "✅ Build completed!"
	@echo "   • Main binary: ./clipboard-manager"
	@echo "   • Build copy: ./build/clipboard-manager"

# Build different variants
build-all:
	@echo "🔨 Building all variants..."
	@mkdir -p build
	@echo "Building standard version..."
	go build -o build/clipboard-manager
	@echo "Building with debug info..."
	go build -gcflags="all=-N -l" -o build/clipboard-manager-debug
	@echo "Building optimized version..."
	go build -ldflags="-s -w" -o build/clipboard-manager-optimized
	@echo "✅ All builds completed!"
	@ls -la build/clipboard-manager*

# Create release package
release: build-all
	@echo "📦 Creating release package..."
	@mkdir -p build/release
	@cp build/clipboard-manager build/release/
	@cp README.md build/release/
	@cp LICENSE build/release/
	@cp -r scripts build/release/
	@cp -r docs build/release/
	@echo "✅ Release package created in build/release/"
	@echo "Contents:"
	@ls -la build/release/

# Run all tests
test:
	@echo "🧪 Running tests..."
	@./tests/run_tests.sh

# Run tests with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	go test -v -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -f clipboard-manager
	rm -f clipboard-manager-test
	rm -f tests/clipboard-manager-test
	rm -f build/clipboard-manager*
	rm -f coverage.out
	rm -f coverage.html
	go clean
	@echo "✅ Clean completed!"

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	go mod tidy
	go mod download
	@echo "✅ Dependencies installed!"

# Run the application
run:
	@echo "🚀 Starting Clipboard Manager..."
	./clipboard-manager

# Run in daemon mode
daemon:
	@echo "🔧 Starting Clipboard Manager daemon..."
	./clipboard-manager daemon

# Show GUI
show:
	@echo "🖥️  Opening Clipboard Manager GUI..."
	./clipboard-manager show

# Install system-wide (requires sudo)
install: build
	@echo "📥 Installing system-wide..."
	sudo cp clipboard-manager /usr/local/bin/
	@echo "✅ Installed to /usr/local/bin/clipboard-manager"

# Uninstall system-wide (requires sudo)
uninstall:
	@echo "🗑️  Uninstalling system-wide..."
	sudo rm -f /usr/local/bin/clipboard-manager
	@echo "✅ Uninstalled from /usr/local/bin/"

# Complete uninstall (removes everything)
uninstall-complete:
	@echo "🗑️  Running complete uninstall..."
	@./scripts/uninstall.sh

# Help
help:
	@echo "Clipboard Manager - Available targets:"
	@echo ""
	@echo "  build          - Build the clipboard manager binary"
	@echo "  build-all      - Build multiple variants (standard, debug, optimized)
	@echo "  release        - Create release package with documentation"
	@echo "  test           - Run all tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  clean          - Clean build artifacts"
	@echo "  deps           - Install/update dependencies"
	@echo "  run            - Run the application"
	@echo "  daemon         - Run in daemon mode"
	@echo "  show           - Show GUI"
	@echo "  install        - Install system-wide (requires sudo)"
	@echo "  uninstall      - Uninstall system-wide (requires sudo)"
	@echo "  uninstall-complete - Complete uninstall (removes everything)"
	@echo "  help           - Show this help"
	@echo ""
	@echo "Examples:"
	@echo "  make build     # Build the application"
	@echo "  make build-all # Build all variants"
	@echo "  make release   # Create release package"
	@echo "  make test      # Run tests"
	@echo "  make clean     # Clean up"