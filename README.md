# Name Generator Provider - Local Development

This directory contains the Go source code for a custom Terraform provider for name generation.

## 🚀 Quick Start

### Option A: Docker Build (No Go Installation Required) ⭐
 
```powershell
# From this directory - requires Docker Desktop
.\build-docker.ps1
```

**Alternative Docker commands:**
```powershell
# Using docker-compose (if you prefer)
docker-compose up --build
```

### Option B: Local Go Build (Requires Go 1.21+)

```bash
# From this directory
make install
```

### 1. Build and Install Provider Locally

### 2. Use in Terraform Examples

```bash
cd ../../examples/name-generator-provider-example/
terraform init   # ✅ Now works - downloads from local plugins
terraform plan   # ✅ Uses the namegen_name resource
terraform apply  # ✅ Generates names via provider
```

## 📁 Provider Structure

```
name-generator-provider/
├── main.go              # Provider entry point
├── provider.go          # Provider schema and configuration
├── resource_name.go     # namegen_name resource implementation
├── go.mod              # Go module definition
├── Makefile            # Build automation
└── README.md           # This file
```

## 🔧 Development Workflow

### Initial Setup
```bash
# Initialize Go modules (first time only)
make init
```

### Build and Test
```bash
# Build provider binary
make build

# Install locally for testing
make install

# Clean build artifacts
make clean
```

### Terraform Usage
```hcl
terraform {
  required_providers {
    namegen = {
      source  = "local/namegen"  # Uses locally installed provider
      version = "~> 1.0"
    }
  }
}

provider "namegen" {
  api_base_url           = "https://your-api-endpoint.com"
  default_cloudregion    = "gfr"
  default_platform_code  = "CC"
  default_environment    = "d"
}

resource "namegen_name" "example" {
  project_id   = "my-project"
  assettag     = "100001"
  name_context = "Web Server"
  resource_type = "st"
}

output "name" {
  value = namegen_name.example.name
}
```

## 🎯 Resource Schema

### Provider Configuration
- `api_base_url` (Optional) - Name generation API endpoint
- `default_cloudregion` (Optional) - Default region for all resources
- `default_platform_code` (Optional) - Default platform code
- `default_environment` (Optional) - Default environment

### Resource: `namegen_name`

**Required Arguments:**
- `project_id` - GCP project identifier
- `assettag` - 6-digit asset tag
- `name_context` - Resource context/purpose

**Optional Arguments:**
- `resource_type` - Resource type code (inherits provider default)
- `environment` - Environment identifier (inherits provider default)
- `cloudregion` - Cloud region code (inherits provider default)
- `platform_code` - Platform code (inherits provider default)

**Computed Attributes:**
- `id` - Terraform resource ID
- `name` - Generated resource name
- `cache_key` - Unique cache identifier
- `cached` - Whether result was cached
- `last_updated` - Timestamp of last update

## 🔄 Local Installation Process

When you run `make install`, the provider is installed to:
```
~/.terraform.d/plugins/local/namegen/1.0.0/windows_amd64/terraform-provider-namegen
```

Terraform finds it when you specify:
```hcl
source = "local/namegen"
```

## 🚀 Moving to Production

### 1. Publish to GitHub Releases
```bash
# Tag and release
git tag v1.0.0
git push --tags

# GitHub Actions can build and release binaries
```

### 2. Update Terraform Configuration
```hcl
terraform {
  required_providers {
    namegen = {
      source  = "github.com/your-org/terraform-provider-namegen"
      version = "~> 1.0"
    }
  }
}
```

### 3. Publish to Terraform Registry
Follow [Terraform Registry Provider Publishing Guide](https://www.terraform.io/docs/registry/providers/publishing.html)

## 🛠️ Customization

### Add New Resource Types
1. Create new resource files (e.g., `resource_database.go`)
2. Add to `provider.go` resources list
3. Rebuild and reinstall

### Enhance API Integration
- Add authentication (API keys, OAuth)
- Add retry logic and error handling
- Add caching mechanisms
- Add validation

### Add Data Sources
Create data source implementations for read-only operations.

## 📚 Further Reading

- [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework)
- [Provider Development Guide](https://developer.hashicorp.com/terraform/plugin/best-practices)
- [Go HTTP Client Tutorial](https://gobyexample.com/http-clients)


