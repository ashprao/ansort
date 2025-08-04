# Makefile for ansort - Natural Alphanumeric Sorting Package
# Provides targets for testing, examples, and development workflow

.PHONY: all test examples clean fmt lint vet mod-tidy help

# Default target
all: test examples

# Colors for output
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
BLUE := \033[34m
RESET := \033[0m

# Help target - shows available commands
help:
	@echo "$(BLUE)Ansort Package - Available Make Targets:$(RESET)"
	@echo ""
	@echo "$(GREEN)Development:$(RESET)"
	@echo "  make test        - Run all tests with verbose output"
	@echo "  make examples    - Run all examples to verify functionality"
	@echo "  make all         - Run both tests and examples (default)"
	@echo ""
	@echo "$(GREEN)Code Quality:$(RESET)"
	@echo "  make fmt         - Format Go code"
	@echo "  make vet         - Run go vet for static analysis"
	@echo "  make lint        - Run golint (if available)"
	@echo "  make mod-tidy    - Clean up go.mod dependencies"
	@echo ""
	@echo "$(GREEN)Utilities:$(RESET)"
	@echo "  make clean       - Clean build artifacts and temporary files"
	@echo "  make verify      - Run comprehensive verification (fmt + vet + test + examples)"
	@echo "  make help        - Show this help message"
	@echo ""
	@echo "$(GREEN)Benchmarks:$(RESET)"
	@echo "  make bench       - Run all tests and benchmarks"
	@echo "  make benchmarks  - Run only benchmark tests (no regular tests)"
	@echo "  make bench-performance - Run performance benchmarks with extended timing"

# Run all tests with verbose output and coverage
test:
	@echo "$(BLUE)Running all tests...$(RESET)"
	@go test -v -race -coverprofile=coverage.out ./...
	@echo "$(GREEN)✅ All tests passed$(RESET)"
	@echo ""
	@echo "$(BLUE)Test coverage:$(RESET)"
	@go tool cover -func=coverage.out | tail -1

# Run all examples
examples:
	@echo "$(BLUE)Running all examples...$(RESET)"
	@echo ""
	@for dir in examples/*/; do \
		if [ -f "$$dir/main.go" ]; then \
			example_name=$$(basename "$$dir"); \
			echo "$(YELLOW)Testing $$example_name example...$(RESET)"; \
			cd "$$dir" && go run main.go > /dev/null 2>&1; \
			if [ $$? -eq 0 ]; then \
				echo "$(GREEN)✅ $$example_name example passed$(RESET)"; \
			else \
				echo "$(RED)❌ $$example_name example failed$(RESET)"; \
				exit 1; \
			fi; \
			cd ../..; \
		fi; \
	done
	@echo ""
	@echo "$(GREEN)✅ All examples completed successfully$(RESET)"

# Run examples with output (for demonstration)
examples-verbose:
	@echo "$(BLUE)Running all examples with output...$(RESET)"
	@echo ""
	@for dir in examples/*/; do \
		if [ -f "$$dir/main.go" ]; then \
			example_name=$$(basename "$$dir"); \
			echo "$(YELLOW)Running $$example_name example:$(RESET)"; \
			echo "----------------------------------------"; \
			cd "$$dir" && go run main.go; \
			echo ""; \
			cd ../..; \
		fi; \
	done

# Format Go code
fmt:
	@echo "$(BLUE)Formatting Go code...$(RESET)"
	@go fmt ./...
	@echo "$(GREEN)✅ Code formatted$(RESET)"

# Run go vet for static analysis
vet:
	@echo "$(BLUE)Running go vet...$(RESET)"
	@go vet ./...
	@echo "$(GREEN)✅ No issues found by go vet$(RESET)"

# Run golint if available
lint:
	@echo "$(BLUE)Running golint...$(RESET)"
	@if command -v golint >/dev/null 2>&1; then \
		golint ./...; \
		echo "$(GREEN)✅ Linting completed$(RESET)"; \
	else \
		echo "$(YELLOW)⚠️  golint not available, skipping$(RESET)"; \
	fi

# Clean up go.mod dependencies
mod-tidy:
	@echo "$(BLUE)Cleaning up go.mod dependencies...$(RESET)"
	@go mod tidy
	@echo "$(GREEN)✅ Dependencies cleaned up$(RESET)"

# Clean build artifacts and temporary files
clean:
	@echo "$(BLUE)Cleaning build artifacts...$(RESET)"
	@rm -f coverage.out
	@go clean ./...
	@echo "$(GREEN)✅ Clean completed$(RESET)"

# Comprehensive verification - runs all quality checks
verify: fmt vet mod-tidy test examples
	@echo ""
	@echo "$(GREEN)🎉 All verification steps passed!$(RESET)"
	@echo "$(BLUE)Package is ready for commit/release$(RESET)"

# Quick verification without examples (for faster feedback during development)
check: fmt vet test
	@echo ""
	@echo "$(GREEN)✅ Quick verification completed$(RESET)"

# Build verification (ensure package compiles)
build:
	@echo "$(BLUE)Verifying package builds...$(RESET)"
	@go build ./...
	@echo "$(GREEN)✅ Package builds successfully$(RESET)"

# Run benchmarks if any exist
bench:
	@echo "$(BLUE)Running benchmarks...$(RESET)"
	@go test -bench=. -benchmem ./...

# Run only benchmark tests (no regular tests)
benchmarks:
	@echo "$(BLUE)Running only benchmark tests...$(RESET)"
	@go test -run=^$$ -bench=. -benchmem ./...

# Run performance benchmarks with more iterations
bench-performance:
	@echo "$(BLUE)Running performance benchmarks (extended)...$(RESET)"
	@go test -run=^$$ -bench=. -benchmem -benchtime=5s ./...

# Show test coverage in browser
coverage: test
	@echo "$(BLUE)Opening test coverage in browser...$(RESET)"
	@go tool cover -html=coverage.out

# List all available examples
list-examples:
	@echo "$(BLUE)Available examples:$(RESET)"
	@for dir in examples/*/; do \
		if [ -f "$$dir/main.go" ]; then \
			example_name=$$(basename "$$dir"); \
			echo "  - $$example_name"; \
		fi; \
	done

# Development workflow targets
dev-setup: mod-tidy
	@echo "$(BLUE)Setting up development environment...$(RESET)"
	@go mod download
	@echo "$(GREEN)✅ Development environment ready$(RESET)"

# Pre-commit hook target
pre-commit: verify
	@echo "$(GREEN)✅ Pre-commit verification passed$(RESET)"

# CI/CD target (what continuous integration should run)
ci: build verify
	@echo "$(GREEN)✅ CI verification completed$(RESET)"
