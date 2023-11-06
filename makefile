.PHONY: build run test clean watch

# Binary name for our application
BINARY_NAME=nnnn

# Root directory for your Go source files
SRC_DIR=./src

# Default command to run when no arguments are passed to Make
all: test build

# Compile the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@cd $(SRC_DIR) && go build -o ../$(BINARY_NAME)

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	@./$(BINARY_NAME)

# Test the application
test:
	@echo "Testing..."
	@cd $(SRC_DIR) && go test ./...

# Clean up the binary
clean:
	@echo "Cleaning up..."
	@cd $(SRC_DIR) && go clean
	@rm -f $(BINARY_NAME)

# Watch for file changes and rebuild the binary
watch:
	@echo "Watching for changes..."
	@reflex -r '\.go$$' -s -- sh -c 'make test && make run'
