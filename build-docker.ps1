# PowerShell script to build Terraform provider using Docker
# No Go installation required!

Write-Host "🔨 Building Terraform Provider using Docker..." -ForegroundColor Cyan

# Check if Docker is running
try {
    docker version | Out-Null
    Write-Host "✅ Docker is running" -ForegroundColor Green
} catch {
    Write-Host "❌ Docker is not running. Please start Docker Desktop." -ForegroundColor Red
    exit 1
}

# Create go.sum if it doesn't exist
if (-not (Test-Path "go.sum")) {
    Write-Host "📦 Creating go.sum..." -ForegroundColor Yellow
    # Create an empty go.sum file - Docker will populate it
    New-Item -Path "go.sum" -ItemType File -Value "" -Force | Out-Null
}

# Build Docker image using Linux base (faster and more reliable)
Write-Host "🐳 Building Docker image..." -ForegroundColor Cyan
docker build -f Dockerfile.linux -t terraform-provider-namegen-builder .

if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Docker build failed" -ForegroundColor Red
    exit 1
}

# Create temporary container and copy the binary
Write-Host "📦 Extracting provider binary..." -ForegroundColor Cyan
$containerId = docker create terraform-provider-namegen-builder

# Copy binary from container to current directory
docker cp "${containerId}:/terraform-provider-namegen.exe" "./terraform-provider-namegen.exe"

# Clean up temporary container
docker rm $containerId | Out-Null

if (-not (Test-Path "terraform-provider-namegen.exe")) {
    Write-Host "❌ Failed to extract binary" -ForegroundColor Red
    exit 1
}

Write-Host "✅ Provider binary built successfully!" -ForegroundColor Green

# Install provider locally
Write-Host "🚀 Installing provider locally..." -ForegroundColor Cyan

# Create plugin directory
$pluginDir = "$env:USERPROFILE\.terraform.d\plugins\local\namegen\1.0.0\windows_amd64"
New-Item -ItemType Directory -Path $pluginDir -Force | Out-Null

# Copy binary to plugin directory
Copy-Item "terraform-provider-namegen.exe" "$pluginDir\terraform-provider-namegen.exe" -Force

Write-Host "✅ Provider installed to: $pluginDir" -ForegroundColor Green

# Verify installation
Write-Host "🔍 Verifying installation..." -ForegroundColor Cyan
if (Test-Path "$pluginDir\terraform-provider-namegen.exe") {
    $fileInfo = Get-Item "$pluginDir\terraform-provider-namegen.exe"
    Write-Host "✅ Provider binary size: $($fileInfo.Length) bytes" -ForegroundColor Green
    Write-Host "✅ Provider ready to use!" -ForegroundColor Green
    Write-Host ""
    Write-Host "🎯 Next steps:" -ForegroundColor Yellow
    Write-Host "   cd ../../examples/name-generator-provider-example" -ForegroundColor White
    Write-Host "   terraform init" -ForegroundColor White
    Write-Host "   terraform plan" -ForegroundColor White
} else {
    Write-Host "❌ Installation verification failed" -ForegroundColor Red
}

# Clean up Docker image (optional)
Write-Host "🧹 Cleaning up..." -ForegroundColor Cyan
docker rmi terraform-provider-namegen-builder | Out-Null

Write-Host "🎉 Docker build complete!" -ForegroundColor Green