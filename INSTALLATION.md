# Installation Guide

This guide will help you install and set up the Chrome-Stable uTLS Template Generator.

## Prerequisites

### Required Software

1. **Go 1.21 or later**
   - Download from: https://golang.org/dl/
   - Verify installation: `go version`

2. **Git**
   - Download from: https://git-scm.com/downloads
   - Verify installation: `git --version`

3. **Chrome/Chromium Browser** (for testing)
   - Download from: https://www.google.com/chrome/
   - Or Chromium: https://www.chromium.org/getting-involved/download-chromium

### Optional Software

- **Docker** (for containerized deployment)
- **Make** (for using the Makefile)
- **golangci-lint** (for code linting)

## Installation Methods

### Method 1: Build from Source (Recommended)

1. **Clone the repository**
   ```bash
   git clone https://github.com/phxtom/betanet/edit/chrome-stable-utls
   cd chrome-utls-template-generator
   ```

2. **Install dependencies**
   ```bash
   go mod download
   go mod tidy
   ```

3. **Build the application**
   ```bash
   go build -o chrome-utls-gen .
   ```

4. **Test the installation**
   ```bash
   ./chrome-utls-gen --help
   ```

### Method 2: Using Make (if available)

1. **Clone and build**
   ```bash
   git clone https://github.com/betanet/chrome-utls-template-generator.git
   cd chrome-utls-template-generator
   make build
   ```

2. **Run tests**
   ```bash
   make test
   ```

3. **Install locally**
   ```bash
   make install
   ```

### Method 3: Docker

1. **Build Docker image**
   ```bash
   docker build -t chrome-utls-gen .
   ```

2. **Run container**
   ```bash
   docker run --rm chrome-utls-gen --help
   ```

## Platform-Specific Instructions

### Windows

1. **Install Go**
   - Download the Windows MSI installer from https://golang.org/dl/
   - Run the installer and follow the prompts
   - Add Go to your PATH if not done automatically

2. **Install Git**
   - Download from https://git-scm.com/download/win
   - Use default settings during installation

3. **Build the application**
   ```cmd
   go build -o chrome-utls-gen.exe .
   ```

4. **Run the application**
   ```cmd
   chrome-utls-gen.exe --help
   ```

### macOS

1. **Install Go using Homebrew**
   ```bash
   brew install go
   ```

2. **Install Git**
   ```bash
   brew install git
   ```

3. **Build and run**
   ```bash
   go build -o chrome-utls-gen .
   ./chrome-utls-gen --help
   ```

### Linux (Ubuntu/Debian)

1. **Install Go**
   ```bash
   sudo apt update
   sudo apt install golang-go
   ```

2. **Install Git**
   ```bash
   sudo apt install git
   ```

3. **Build and run**
   ```bash
   go build -o chrome-utls-gen .
   ./chrome-utls-gen --help
   ```

## Configuration

### 1. Create Configuration File

Copy the example configuration:
```bash
cp config.yaml.example ~/.chrome-utls-gen.yaml
```

### 2. Edit Configuration

Edit the configuration file to match your environment:
```yaml
# Example configuration
app:
  log_level: "info"
  verbose: false

chrome:
  default_version: "120.0.6099.109"
  check_interval: "1h"

template:
  output_dir: "./templates"
```

## Quick Start

### 1. Generate Your First Template

```bash
# Generate template for latest Chrome
./chrome-utls-gen generate

# Generate for specific version
./chrome-utls-gen generate --version 120.0.6099.109
```

### 2. Test the Template

```bash
# Test generated template
./chrome-utls-gen test --template ./templates/chrome-120.0.6099.109.json

# Test with live Chrome
./chrome-utls-gen test --live
```

### 3. Monitor for Updates

```bash
# Start monitoring for new Chrome releases
./chrome-utls-gen monitor --interval 1h

# Auto-generate templates on new releases
./chrome-utls-gen monitor --auto-generate
```

## Verification

### Check Installation

1. **Verify Go installation**
   ```bash
   go version
   # Should show: go version go1.21.x windows/amd64
   ```

2. **Verify application**
   ```bash
   ./chrome-utls-gen --version
   # Should show: 1.0.0
   ```

3. **Run basic test**
   ```bash
   ./chrome-utls-gen generate --version 120.0.6099.109
   # Should create a template file
   ```

### Test Template Generation

```bash
# Generate template
./chrome-utls-gen generate --version 120.0.6099.109

# Verify template was created
ls -la templates/

# Test template
./chrome-utls-gen test --template templates/chrome-120.0.6099.109.json
```

## Troubleshooting

### Common Issues

1. **"go: command not found"**
   - Install Go from https://golang.org/dl/
   - Add Go to your PATH environment variable

2. **"git: command not found"**
   - Install Git from https://git-scm.com/downloads
   - Add Git to your PATH environment variable

3. **"Permission denied"**
   - On Unix-like systems: `chmod +x chrome-utls-gen`
   - On Windows: Run as Administrator if needed

4. **"Template generation failed"**
   - Check internet connectivity
   - Verify Chrome version exists
   - Check configuration file

5. **"Chrome not found"**
   - Install Chrome/Chromium browser
   - Specify Chrome path manually: `--chrome-path /path/to/chrome`

### Debug Mode

Enable verbose logging for troubleshooting:
```bash
./chrome-utls-gen --verbose generate
```

### Network Issues

If you're behind a proxy:
```bash
# Set proxy environment variables
export HTTP_PROXY=http://proxy.example.com:8080
export HTTPS_PROXY=http://proxy.example.com:8080

# Or configure in config file
network:
  proxy:
    enabled: true
    http_proxy: "http://proxy.example.com:8080"
```

## Development Setup

### Install Development Tools

1. **Install linter**
   ```bash
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```

2. **Install documentation generator**
   ```bash
   go install golang.org/x/tools/cmd/godoc@latest
   ```

3. **Install security scanner**
   ```bash
   go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
   ```

### Development Commands

```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Generate documentation
godoc -http=:6060
```

## Next Steps

After installation:

1. **Read the documentation**: See `docs/` directory
2. **Configure for your environment**: Edit `~/.chrome-utls-gen.yaml`
3. **Generate your first template**: `./chrome-utls-gen generate`
4. **Test with Betanet**: See `docs/BETANET_INTEGRATION.md`
5. **Set up monitoring**: `./chrome-utls-gen monitor --auto-generate`

## Support

For issues and questions:

1. Check the troubleshooting section above
2. Review the documentation in `docs/`
3. Check the configuration file
4. Enable debug mode for detailed logs
5. Open an issue on GitHub

## Uninstallation

To remove the application:

```bash
# Remove binary
rm chrome-utls-gen

# Remove configuration
rm ~/.chrome-utls-gen.yaml

# Remove templates (if desired)
rm -rf templates/
```
