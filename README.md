
# Name Generator Terraform Provider

This directory contains the Go source code for a custom Terraform provider that generates names using a configurable API and supports multiple name generation types.

## � Overview

This provider supports three name generation types:
- **host**: Generate hostnames (requires `hostname_type` and `stack_id`)
- **DB**: Generate database names (requires `hostname_type` and `stack_id`)
- **gcpname**: Generate GCP resource names (requires `resource_type`, `cloudregion`, `platform_code`, `environment`, `assettag`, `name_context`)

**Distribution Options:**
- 🏢 **Terraform Cloud Private Registry** (requires Business/Enterprise plan)
- 📦 **GitHub Releases** (manual installation, free)
- 💻 **Local Development** (using `make install`)

## �🚀 Quick Start

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

This repo includes a workflow at `.github/workflows/build-provider.yml` that:
- Triggers automatically when you push a version tag (e.g., `v1.0.0`)
- Builds the provider for **4 platforms**:
  - Linux (amd64)
  - Windows (amd64)
  - macOS Intel (darwin_amd64)
  - macOS Apple Silicon (darwin_arm64)
- Generates SHA256SUMS for security verification
- Creates a GitHub Release with all binaries attached
- Uploads artifacts for each build

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

## 🚀 Publishing to Terraform Cloud Private Registry

### Prerequisites
- GitHub repository named `terraform-provider-namegen` (must start with `terraform-provider-`)
- Terraform Cloud account with **Business or Enterprise plan** (Free/Standard plans cannot publish providers)
- Repository must have at least one version tag (e.g., `v1.0.0`)

---

### Publishing Process (5 Categories)

#### **Category 1: Build, package, and prepare provider metadata**
1. Ensure your code is pushed to GitHub:
   ```bash
   git add .
   git commit -m "Prepare provider for release"
   git push origin main
   ```

2. Create and push a version tag:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```
   - This triggers the GitHub Actions workflow
   - Workflow builds binaries for all platforms and creates a GitHub Release

---

#### **Category 2: Authenticate to Terraform Cloud and set up the private registry**
1. Create a Terraform Cloud account at https://app.terraform.io/
2. Create an organization (e.g., `your-company-name`)
3. Verify your plan supports private providers:
   - Go to **Settings** in your organization
   - Only **Business** and **Enterprise** plans can publish private providers
   - If on Free/Standard, upgrade is required

---

#### **Category 3: Upload the provider to the private registry**
1. Connect GitHub to Terraform Cloud:
   - Go to **Settings** → **VCS Providers**
   - Click **Add VCS Provider** → Choose **GitHub.com**
   - Authorize Terraform Cloud to access your repositories

2. Publish the provider:
   - Go to **Registry** → **Providers**
   - Click **Publish** → **Provider**
   - Select your GitHub VCS connection
   - Select the `terraform-provider-namegen` repository
   - Click **Publish Provider**

3. Verify the provider appears in the registry with your version (e.g., v1.0.0)

---

#### **Category 4: Update Terraform configuration and test integration**
1. Create a Terraform configuration to use your private provider:
   ```hcl
   terraform {
     required_providers {
       namegen = {
         source  = "app.terraform.io/YOUR-ORG/namegen"
         version = "1.0.0"
       }
     }
   }
   
   provider "namegen" {}
   
   resource "namegen_name" "test" {
     type          = "host"
     api_url       = "https://your-api-endpoint.com"
     hostname_type = "app"
     stack_id      = "test01"
   }
   ```

2. Authenticate Terraform CLI:
   ```bash
   terraform login
   ```

3. Test the provider:
   ```bash
   terraform init    # Downloads provider from private registry
   terraform plan    # Verify configuration
   terraform apply   # Test provider functionality
   ```

---

#### **Category 5: Document the process for the team**
- Share the provider source address with your team: `app.terraform.io/YOUR-ORG/namegen`
- Document any required API endpoints or configuration
- Create usage examples for common scenarios
- Update this README with team-specific instructions

---

### Alternative: GitHub Releases Distribution (No Private Registry Required)

If you don't have a Business/Enterprise plan, users can install the provider manually from GitHub Releases:

1. Download the appropriate zip file for your platform from the [Releases page](../../releases)
2. Extract the binary to your Terraform plugins directory:
   - Linux/Mac: `~/.terraform.d/plugins/github.com/YOUR-USERNAME/namegen/v1.0.0/OS_ARCH/`
   - Windows: `%APPDATA%\terraform.d\plugins\github.com\YOUR-USERNAME\namegen\v1.0.0\OS_ARCH\`
3. Use in Terraform:
   ```hcl
   terraform {
     required_providers {
       namegen = {
         source  = "github.com/YOUR-USERNAME/namegen"
         version = "1.0.0"
       }
     }
   }
   ```

---

### Publishing New Versions

To release a new version:
```bash
git tag v1.0.1
git push origin v1.0.1
```

The workflow will automatically:
- Build binaries for all platforms
- Generate SHA256SUMS
- Create a GitHub Release
- Terraform Cloud will detect and publish the new version

---

### Official Documentation
- [Terraform Cloud Private Registry](https://developer.hashicorp.com/terraform/cloud-docs/registry/publish-providers)
- [Terraform Cloud Pricing](https://developer.hashicorp.com/terraform/cloud/pricing)
- [Provider Development Guide](https://developer.hashicorp.com/terraform/plugin/best-practices)


## 🛠️ Customization & Extensibility

- Add new resource types by creating new files and updating `provider.go`.
- Enhance API integration (authentication, retries, caching, validation).
- Add data sources for read-only operations.


## 📚 Further Reading

- [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework)
- [Provider Development Guide](https://developer.hashicorp.com/terraform/plugin/best-practices)
- [Go HTTP Client Tutorial](https://gobyexample.com/http-clients)


