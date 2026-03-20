# Makefile for building and installing the namegen provider locally

.PHONY: build install clean version

# Get version from git tags or use 'dev' if no tag exists
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

# Build the provider binary with version embedded
build:
	go build -ldflags="-X 'main.version=$(VERSION)'" -o terraform-provider-namegen

# Install the provider locally for development
install: build
	@echo "Installing provider locally with version $(VERSION)..."
	@mkdir -p ~/.terraform.d/plugins/local/namegen/$(VERSION)/windows_amd64/
	@cp terraform-provider-namegen ~/.terraform.d/plugins/local/namegen/$(VERSION)/windows_amd64/
	@echo "Provider installed to ~/.terraform.d/plugins/local/namegen/$(VERSION)/windows_amd64/"

# Install for Linux (if needed)
install-linux: build
	@echo "Installing provider locally for Linux with version $(VERSION)..."
	@mkdir -p ~/.terraform.d/plugins/local/namegen/$(VERSION)/linux_amd64/
	@cp terraform-provider-namegen ~/.terraform.d/plugins/local/namegen/$(VERSION)/linux_amd64/

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

# Show current version
version:
	@echo $(VERSION)

# Help
help:
	@echo "Available commands:"
	@echo "  build        - Build the provider binary (embeds version from git tag)"
	@echo "  install      - Build and install provider locally"
	@echo "  clean        - Remove build artifacts"
	@echo "  init         - Initialize Go modules"
	@echo "  test         - Run tests"
	@echo "  version      - Show current version from git tags"