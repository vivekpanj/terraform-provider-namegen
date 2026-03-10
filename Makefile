# Makefile for building and installing the namegen provider locally

.PHONY: build install clean

# Build the provider binary
build:
	go build -o terraform-provider-namegen

# Install the provider locally for development
install: build
	@echo "Installing provider locally..."
	@mkdir -p ~/.terraform.d/plugins/local/namegen/1.0.0/windows_amd64/
	@cp terraform-provider-namegen ~/.terraform.d/plugins/local/namegen/1.0.0/windows_amd64/
	@echo "Provider installed to ~/.terraform.d/plugins/local/namegen/1.0.0/windows_amd64/"

# Install for Linux (if needed)
install-linux: build
	@echo "Installing provider locally for Linux..."
	@mkdir -p ~/.terraform.d/plugins/local/namegen/1.0.0/linux_amd64/
	@cp terraform-provider-namegen ~/.terraform.d/plugins/local/namegen/1.0.0/linux_amd64/

# Clean build artifacts
clean:
	@rm -f terraform-provider-namegen
	@rm -f terraform-provider-namegen.exe

# Initialize Go modules
init:
	go mod init github.com/your-org/terraform-provider-namegen
	go mod tidy

# Run tests (when you add them)
test:
	go test -v ./...

# Help
help:
	@echo "Available commands:"
	@echo "  build        - Build the provider binary"
	@echo "  install      - Build and install provider locally"
	@echo "  clean        - Remove build artifacts"
	@echo "  init         - Initialize Go modules"
	@echo "  test         - Run tests"