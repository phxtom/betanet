# Chrome-Stable uTLS Template Generator - User Tutorial

## 🚀 Quick Command Reference

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

---

## 🎯 What is this tool?

The Chrome-Stable uTLS Template Generator is a utility that creates **exact TLS handshake templates** that match Chrome Stable's fingerprint. This enables network traffic to blend in with normal web browsing, making it impossible for deep-packet inspectors to distinguish it from legitimate Chrome traffic.

## 🚀 Quick Start Guide

### Prerequisites
- Windows, macOS, or Linux
- Internet connection (for Chrome version detection)
- No special software installation required

### Installation
1. **Download** the `chrome-utls-gen.exe` file
2. **Place it** in a folder of your choice
3. **That's it!** No installation needed

## 📋 How to Use the Tool

### Method 1: Command Prompt (Recommended)

1. **Open Command Prompt** (not by double-clicking the .exe)
2. **Navigate** to the folder containing `chrome-utls-gen.exe`
3. **Run commands** directly

### Method 2: Using the Batch File

1. **Double-click** `run.bat` (if available)
2. **Follow the prompts** to enter commands
3. **The window will stay open** for you to see results

## 🔧 Available Commands

### 1. Get Help
```bash
chrome-utls-gen.exe --help
```
Shows all available commands and options.

### 2. Test Existing Template
```bash
chrome-utls-gen.exe test --template templates\chrome-120.0.6099.109.json
```
Tests a pre-generated template and shows:
- ✅ JA3 fingerprint match
- ✅ JA4 fingerprint match  
- ✅ Server connection test

**Expected Output:**
```
Testing template: chrome-120.0.6099.109.json
Chrome version: 120.0.6099.109
Generated: 2025-08-11T15:26:23Z

Fingerprint Comparison:
=======================
Template JA3:  2b602ec5cef66b61c431b3ab0c167f18
Generated JA3: 2b602ec5cef66b61c431b3ab0c167f18
JA3 Match:     true

Template JA4:  t14339_d0_h0_c0_e256-1-8448-21760-35072-48384-61696_g_p_a_r_u_80173352edaf4b6000f80bc3421451bd
Generated JA4: t14339_d0_h0_c0_e256-1-8448-21760-35072-48384-61696_g_p_a_r_u_80173352edaf4b6000f80bc3421451bd
JA4 Match:     true

Server Test:
============
Server test passed ✓
```

### 3. Generate New Template
```bash
chrome-utls-gen.exe generate --version 120.0.6099.109 --force
```
Creates a new template for a specific Chrome version.

**Options:**
- `--version` - Specify Chrome version (default: latest)
- `--force` - Overwrite existing template
- `--output` - Specify output directory

### 4. Monitor for New Chrome Releases
```bash
chrome-utls-gen.exe monitor --interval 1h --auto-generate
```
Continuously checks for new Chrome versions and auto-generates templates.

**Options:**
- `--interval` - Check interval (30m, 1h, 6h)
- `--auto-generate` - Auto-create templates for new versions
- `--output` - Output directory for templates

### 5. Get Command-Specific Help
```bash
chrome-utls-gen.exe test --help
chrome-utls-gen.exe generate --help
chrome-utls-gen.exe monitor --help
```

## 📁 Understanding the Files

### Template Files
- **Location:** `templates/` folder
- **Format:** JSON files
- **Naming:** `chrome-[version].json`
- **Contains:** TLS handshake configuration, fingerprints, metadata

### Example Template Structure
```json
{
  "version": "120.0.6099.109",
  "timestamp": "2025-08-11T15:26:23Z",
  "client_hello": {
    "version": "TLS 1.2",
    "random": "base64-encoded-random-data",
    "session_id": "base64-encoded-session-id",
    "cipher_suites": ["0x1301", "0x1302", "0x1303"],
    "compression_methods": [0],
    "extensions": [
      {
        "type": 0,
        "data": "base64-encoded-extension-data"
      }
    ]
  },
  "ja3_fingerprint": "2b602ec5cef66b61c431b3ab0c167f18",
  "ja4_fingerprint": "t14339_d0_h0_c0_e256-1-8448-21760-35072-48384-61696_g_p_a_r_u_80173352edaf4b6000f80bc3421451bd"
}
```

## 🔍 Understanding Fingerprints

### JA3 Fingerprint
- **Purpose:** TLS client fingerprinting
- **Format:** MD5 hash of TLS parameters
- **Example:** `2b602ec5cef66b61c431b3ab0c167f18`

### JA4 Fingerprint
- **Purpose:** Enhanced TLS fingerprinting
- **Format:** Human-readable string with MD5 hash
- **Example:** `t14339_d0_h0_c0_e256-1-8448-21760-35072-48384-61696_g_p_a_r_u_80173352edaf4b6000f80bc3421451bd`

## 🛠️ Troubleshooting

### Common Issues

#### 1. Command Prompt Closes Immediately
**Problem:** Double-clicking the .exe file
**Solution:** Open Command Prompt first, then run the tool

#### 2. "Template already exists" Error
**Problem:** Template file already exists
**Solution:** Use `--force` flag to overwrite

#### 3. "Failed to generate template" Error
**Problem:** Network issue or unsupported Chrome version
**Solution:** Check internet connection and try a different version

#### 4. "Invalid extension data" Error
**Problem:** Template generation issue (FIXED)
**Solution:** The bug has been resolved. Template generation now works correctly.

### Getting Help
```bash
# General help
chrome-utls-gen.exe --help

# Command-specific help
chrome-utls-gen.exe [command] --help

# Verbose output for debugging
chrome-utls-gen.exe [command] --verbose
```

## 🎯 Use Cases

### 1. Network Traffic Obfuscation
Generate templates to make your traffic look like Chrome browsing.

### 2. Penetration Testing
Test network security by mimicking legitimate browser traffic.

### 3. Research & Development
Study TLS fingerprinting and browser behavior.

### 4. Betanet Integration
Use with Betanet network for L2 cover transport requirements.

## 🔒 Security Considerations

- **Templates are deterministic** - same input produces same output
- **No personal data** is collected or transmitted
- **Internet access** only needed for Chrome version detection
- **Local operation** - all processing happens on your machine

## 📚 Advanced Usage

### Custom Configuration
Create a `config.yaml` file for custom settings:
```yaml
app:
  name: "chrome-utls-gen"
  version: "1.0.0"

chrome:
  version_detection:
    enabled: true
    sources:
      - "omaha_proxy"
      - "chrome_releases"
      - "chrome_version_api"

template:
  output_dir: "./templates"
  format: "json"
  include_metadata: true
```

### Integration with Other Tools
The generated templates can be used with:
- uTLS library in Go
- Custom TLS implementations
- Network analysis tools
- Betanet clients

## 🆘 Support

### Documentation
- **README.md** - Project overview
- **PROJECT_SUMMARY.md** - Technical details
- **docs/BETANET_INTEGRATION.md** - Betanet-specific usage

### Testing Your Setup
1. **Test existing template:**
   ```bash
   chrome-utls-gen.exe test --template templates\chrome-120.0.6099.109.json
   ```

2. **Verify all commands work:**
   ```bash
   chrome-utls-gen.exe --help
   chrome-utls-gen.exe test --help
   chrome-utls-gen.exe generate --help
   chrome-utls-gen.exe monitor --help
   ```

3. **Check file structure:**
   ```
   ├── chrome-utls-gen.exe
   ├── templates/
   │   └── chrome-120.0.6099.109.json
   ├── run.bat
   └── INSTRUCTIONS.md
   ```

## 🎉 Success Indicators

You know the tool is working correctly when:
- ✅ `test` command shows "JA3 Match: true" and "JA4 Match: true"
- ✅ `--help` commands show proper documentation
- ✅ Template files are generated in JSON format
- ✅ No error messages appear during testing

---

**Happy fingerprinting! 🚀**
