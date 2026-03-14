
# Name Generator Terraform Provider

This directory contains the Go source code for a custom Terraform provider that generates names using a configurable API and supports multiple name generation types.

## 🚀 Quick Start

### Local Go Build (Requires Go 1.21+)

```bash
# From this directory
make install
```

### Example Usage

```bash
# (Optional) Change to your Terraform example directory
cd ../../examples/name-generator-provider-example/
terraform init
terraform plan
terraform apply
```

## 📁 Provider Structure

```
name-generator-provider/
├── main.go              # Provider entry point
├── provider.go          # Provider schema and configuration
├── resource_name.go     # namegen_name resource implementation
├── go.mod               # Go module definition
├── Makefile             # Build automation
└── README.md            # This file
```


## 🔧 Development Workflow

### Initial Setup
```bash
make init
```

### Build and Test
```bash
make build      # Build provider binary
make install    # Install locally for testing
make clean      # Clean build artifacts
```

### GitHub Actions Workflow

This repo includes a workflow at `.github/workflows/build-provider.yml` that builds the provider for Linux and Windows and uploads the binaries as artifacts on every push.

## Terraform Usage Example

```hcl
terraform {
  required_providers {
    namegen = {
      source  = "local/namegen"  # For local development
      version = "~> 1.0"
    }
  }
}

provider "namegen" {
  api_base_url = "https://your-api-endpoint.com" # Required
  # Optionally set provider-level defaults for fields below
  default_cloudregion    = "gfr"
  default_platform_code  = "CC"
  default_environment    = "d"
}

# --- Option 1: type = "host" ---
resource "namegen_name" "host_example" {
  type          = "host"         # Required: "host"
  api_url       = "https://your-api-endpoint.com" # Required
  hostname_type = "app"          # Required for host
  stack_id      = "stack01"      # Required for host
}

# --- Option 2: type = "DB" ---
resource "namegen_name" "db_example" {
  type          = "DB"           # Required: "DB"
  api_url       = "https://your-api-endpoint.com" # Required
  hostname_type = "db"           # Required for DB
  stack_id      = "stack02"      # Required for DB
}

# --- Option 3: type = "gcpname" ---
resource "namegen_name" "gcp_example" {
  type          = "gcpname"            # Required: "gcpname"
  api_url       = "https://your-api-endpoint.com" # Required
  resource_type = "gcp_cloudstorage"   # Example: GCP Cloud Storage bucket
  cloudregion   = "us-central1"        # Required for gcpname
  platform_code = "CC"                 # Required for gcpname
  environment   = "d"                  # Required for gcpname
  assettag      = "100001"             # Required for gcpname
  name_context  = "Web Server"         # Required for gcpname
}

output "name" {
  value = namegen_name.gcp_example.name
}
```


## 🎯 Provider & Resource Schema

### Provider Configuration
- `api_base_url` (**Required**) – Name generation API endpoint
- `default_cloudregion` (Optional) – Default region for all resources
- `default_platform_code` (Optional) – Default platform code
- `default_environment` (Optional) – Default environment

### Resource: `namegen_name`

**Required Arguments:**
- `type` (**Required**): One of `host`, `DB`, or `gcpname`
- `api_url` (**Required**): API endpoint URL for name generation

**Parameters by Type:**

- If `type = "host"`:
  - `hostname_type` (**Required**)
  - `stack_id` (**Required**)

- If `type = "DB"`:
  - `hostname_type` (**Required**)
  - `stack_id` (**Required**)

- If `type = "gcpname"`:
  - `resource_type` (**Required**)
  - `cloudregion` (**Required**)
  - `platform_code` (**Required**)
  - `environment` (**Required**)
  - `assettag` (**Required**)
  - `name_context` (**Required**)

**Computed Attributes:**
- `id` – Terraform resource ID
- `name` – Generated resource name
- `cache_key` – Unique cache identifier
- `cached` – Whether result was cached
- `last_updated` – Timestamp of last update


## 🔄 Local Installation

When you run `make install`, the provider is installed to:

```
~/.terraform.d/plugins/local/namegen/1.0.0/<os>_<arch>/terraform-provider-namegen
```

Terraform finds it when you specify:

```hcl
source = "local/namegen"
```

## 🚀 Publishing

1. Tag and push a release:
   ```bash
   git tag v1.0.0
   git push --tags
   ```
2. GitHub Actions will build and upload binaries as artifacts.
3. For Terraform Registry, follow the [Provider Publishing Guide](https://www.terraform.io/docs/registry/providers/publishing.html).


## 🛠️ Customization & Extensibility

- Add new resource types by creating new files and updating `provider.go`.
- Enhance API integration (authentication, retries, caching, validation).
- Add data sources for read-only operations.


## 📚 Further Reading

- [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework)
- [Provider Development Guide](https://developer.hashicorp.com/terraform/plugin/best-practices)
- [Go HTTP Client Tutorial](https://gobyexample.com/http-clients)


