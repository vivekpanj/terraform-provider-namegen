# Dockerfile for building Terraform provider without local Go installation
FROM golang:1.21-windowsservercore-ltsc2019 as builder

# Set working directory
WORKDIR /app

# Copy Go modules files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY *.go ./

# Build the provider for Windows
RUN go build -o terraform-provider-namegen.exe

# Use a minimal Windows image for the final stage
FROM mcr.microsoft.com/windows/nanoserver:ltsc2019

# Copy the built binary
COPY --from=builder /app/terraform-provider-namegen.exe /terraform-provider-namegen.exe

# Default command
CMD ["terraform-provider-namegen.exe"]