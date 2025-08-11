# Chrome-Stable uTLS Template Generator

A program that produces exact TLS handshake bytes Chrome Stable sends, enabling Betanet traffic to blend in with normal web browsing.

## Quick Command Reference

### Essential Commands:
```bash
# Generate a new template
chrome-utls-gen.exe generate --version 120.0.6099.109 --force

# Test an existing template
chrome-utls-gen.exe test --template templates\chrome-120.0.6099.109.json

# Monitor for new Chrome versions
chrome-utls-gen.exe monitor --interval 1h --auto-generate

# Get help
chrome-utls-gen.exe --help
```

### Common Options:
- `--version` - Specify Chrome version (default: latest)
- `--force` - Overwrite existing template
- `--output` - Specify output directory
- `--verbose` - Show detailed output

## Overview

This tool generates deterministic ClientHello templates that match Chrome Stable's TLS fingerprint, making it impossible for deep-packet inspectors to distinguish Betanet traffic from legitimate Chrome browsing.

## Features

- **Deterministic ClientHello Generation**: Produces exact TLS handshake bytes matching Chrome Stable
- **JA3/JA4 Self-Test CLI**: Verifies generated templates match Chrome's fingerprint
- **Auto-Refresh**: Automatically updates when new Chromium stable tags are released
- **Origin Mirroring Support**: Enables Betanet's L2 cover transport requirements
- **Multi-Platform**: Works on Windows, macOS, and Linux

## Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Chrome Tag    │───▶│  Template Gen    │───▶│  uTLS Template  │
│   Monitor       │    │                  │    │                 │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌──────────────────┐
                       │   Self-Test CLI  │
                       │   (JA3/JA4)      │
                       └──────────────────┘
```

## Installation

```
# Download Zip from GitHub Repo

```

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


```

## Usage

### Generate Templates

```bash
# Generate template for latest Chrome Stable
chrome-utls-gen generate

# Generate for specific Chrome version
chrome-utls-gen generate --version 120.0.6099.109

# Generate with custom output path
chrome-utls-gen generate --output ./templates/
```

### Self-Test

```bash
# Test generated template against Chrome
chrome-utls-gen test --template ./templates/chrome-120.0.6099.109.json

# Test with live Chrome instance
chrome-utls-gen test --live --chrome-path /path/to/chrome
```

### Monitor for Updates

```bash
# Start monitoring for new Chrome releases
chrome-utls-gen monitor --interval 1h

# Auto-generate templates on new releases
chrome-utls-gen monitor --auto-generate
```

## Template Format

Generated templates are JSON files containing:

```json
{
  "version": "120.0.6099.109",
  "timestamp": "2024-01-15T10:30:00Z",
  "client_hello": {
    "version": "TLS 1.2",
    "random": "base64-encoded-random",
    "session_id": "base64-encoded-session-id",
    "cipher_suites": ["0x1301", "0x1302", ...],
    "compression_methods": [0],
    "extensions": [
      {
        "type": "server_name",
        "data": "base64-encoded-data"
      }
    ]
  },
  "ja3_fingerprint": "771,4865-4867-4866-49195-49199-52393-52392-49196-49200-49162-49161-49171-49172-156-157-47-53,0-23-65281-10-11-35-16-5-13-28-21,29-23-24-25-256-257,0",
  "ja4_fingerprint": "t13d0011h2_771_4865-4867-4866-49195-49199-52393-52392-49196-49200-49162-49161-49171-49172-156-157-47-53_0-23-65281-10-11-35-16-5-13-28-21_29-23-24-25-256-257_0"
}
```

## Betanet Integration

This tool is designed to support Betanet's L2 cover transport layer requirements:

- **Origin Mirroring**: Templates enable exact fingerprint matching
- **Auto-Calibration**: Supports per-connection calibration pre-flight
- **HTX Compatibility**: Generates templates compatible with HTX protocol
- **TLS 1.3 Support**: Full support for modern TLS versions

## Development

### Prerequisites

- Go 1.21+
- Chrome/Chromium browser for testing
- Network access for Chrome tag monitoring

### Building

```bash
# Build for current platform
go build -o chrome-utls-gen .

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o chrome-utls-gen-linux .
GOOS=windows GOARCH=amd64 go build -o chrome-utls-gen.exe .
GOOS=darwin GOARCH=amd64 go build -o chrome-utls-gen-mac .
```

### Testing

```bash
# Run unit tests
go test ./...

# Run integration tests
go test ./... -tags=integration

# Run with coverage
go test ./... -cover
```

## License

MIT License - see LICENSE file for details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## Security

This tool is designed for legitimate network analysis and Betanet development. Users are responsible for complying with applicable laws and regulations.
