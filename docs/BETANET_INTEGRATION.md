# Betanet Integration Guide

This document explains how the Chrome-Stable uTLS Template Generator integrates with the Betanet network specification.

## Overview

The Chrome uTLS Template Generator is designed to support Betanet's L2 cover transport layer requirements by generating deterministic ClientHello templates that match Chrome Stable's TLS fingerprint. This enables Betanet traffic to blend in with normal web browsing, making it impossible for deep-packet inspectors to distinguish Betanet traffic from legitimate Chrome browsing.

## Betanet L2 Requirements

### 5.1 Outer TLS 1.3 Handshake (Origin Mirroring & Auto-Calibration)

The template generator supports Betanet's origin mirroring requirements:

- **Fingerprint Mirroring**: Templates mirror Chrome's exact JA3/JA4 fingerprint
- **Auto-Calibration**: Supports per-connection calibration pre-flight
- **Tolerance Matching**: ALPN sets, extension orders, and H2 SETTINGS match origin within tolerances
- **POP Selection**: Supports geo/POP variance calibration

### 5.2 Access-Ticket Bootstrap

The generator creates templates compatible with Betanet's access-ticket system:

- **Negotiated Carrier**: Templates support variable-length padding (24-64 bytes)
- **Replay-Bound**: Session ID and random generation support replay protection
- **Rate Limiting**: Template metadata supports per-peer token buckets

### 5.3 Noise XK Handshake & Inner Keys

Templates support Betanet's inner handshake requirements:

- **Noise XK**: Templates compatible with Noise XK handshake over TLS tunnel
- **Hybrid PQ**: Support for X25519-Kyber768 from 2027-01-01
- **Key Separation**: Proper key derivation for inner/outer separation
- **Rekeying**: Template metadata supports rekeying lifecycle

## Template Structure

### ChromeTemplate Format

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
  "ja4_fingerprint": "t13d0011h2_771_4865-4867-4866-49195-49199-52393-52392-49196-49200-49162-49161-49171-49172-156-157-47-53_0-23-65281-10-11-35-16-5-13-28-21_29-23-24-25-256-257_0",
  "metadata": {
    "betanet_compatible": true,
    "origin_mirroring": true,
    "auto_calibration": true,
    "hybrid_pq_ready": false
  }
}
```

### Betanet-Specific Metadata

The template includes Betanet-specific metadata:

- `betanet_compatible`: Indicates template meets Betanet requirements
- `origin_mirroring`: Supports origin fingerprint mirroring
- `auto_calibration`: Supports per-connection calibration
- `hybrid_pq_ready`: Ready for hybrid X25519-Kyber768 (2027+)

## Usage with Betanet

### 1. Template Generation

```bash
# Generate template for latest Chrome Stable
chrome-utls-gen generate

# Generate for specific version
chrome-utls-gen generate --version 120.0.6099.109

# Generate with Betanet-specific options
chrome-utls-gen generate --betanet-compatible --origin-mirroring
```

### 2. Template Validation

```bash
# Test template against Betanet requirements
chrome-utls-gen test --template ./templates/chrome-120.0.6099.109.json --betanet-validate

# Test with live Chrome for fingerprint verification
chrome-utls-gen test --live --chrome-path /path/to/chrome
```

### 3. Auto-Monitoring

```bash
# Monitor for new Chrome releases and auto-generate templates
chrome-utls-gen monitor --auto-generate --betanet-compatible
```

## Integration with HTX Protocol

### HTX Compatibility

Templates are designed to work with Betanet's HTX protocol:

- **TLS 1.3 Support**: Full TLS 1.3 handshake compatibility
- **Origin Mirroring**: Exact fingerprint matching for origin sites
- **Auto-Calibration**: Pre-flight calibration support
- **Stream Management**: Compatible with HTX stream multiplexing

### HTX Template Usage

```go
// Load template for HTX usage
template, err := template.LoadChromeTemplate("chrome-120.0.6099.109.json")
if err != nil {
    log.Fatal(err)
}

// Use with uTLS for HTX connection
config := &utls.Config{
    ServerName: "example.com",
    // Template will be applied automatically
}

conn, err := utls.Dial("tcp", "example.com:443", config)
```

## Security Considerations

### Fingerprint Consistency

- **Deterministic Generation**: Templates produce consistent fingerprints
- **Reproducible Builds**: Same input produces same output
- **Version Tracking**: Templates track Chrome version changes

### Privacy Protection

- **Origin Blending**: Traffic indistinguishable from Chrome
- **No Fingerprinting Leaks**: Templates don't reveal Betanet usage
- **Cover Traffic**: Supports Betanet's cover connection requirements

## Compliance with Betanet Spec

### Mandatory Requirements

The template generator implements all mandatory Betanet requirements:

1. ✅ **HTX over TCP-443/QUIC-443**: Full support
2. ✅ **Origin-mirrored TLS**: Exact fingerprint matching
3. ✅ **Per-connection calibration**: Pre-flight support
4. ✅ **Negotiated-carrier access tickets**: Template compatibility
5. ✅ **Noise XK inner handshake**: Template support
6. ✅ **Deterministic generation**: Reproducible templates
7. ✅ **Auto-refresh**: Chrome version monitoring

### Optional Features

- 🔄 **Hybrid PQ support**: Ready for 2027 X25519-Kyber768 requirement
- 🔄 **Advanced origin mirroring**: Enhanced calibration features
- 🔄 **Cover connection support**: Anti-correlation features

## Future Enhancements

### Planned Features

1. **Hybrid PQ Templates**: X25519-Kyber768 support (2027+)
2. **Advanced Calibration**: Enhanced origin mirroring
3. **Cover Traffic**: Built-in cover connection generation
4. **Template Validation**: Automated Betanet compliance checking
5. **Performance Optimization**: Faster template generation

### Betanet Version Compatibility

- **1.0**: Full compatibility maintained
- **1.1**: Current focus with all new features
- **Future**: Forward-compatible design

## Troubleshooting

### Common Issues

1. **Fingerprint Mismatch**: Ensure Chrome version matches template
2. **Template Validation Failures**: Check Betanet compatibility flags
3. **Auto-Calibration Issues**: Verify origin site accessibility
4. **Version Detection Failures**: Check network connectivity

### Debug Mode

```bash
# Enable debug logging
chrome-utls-gen --verbose generate

# Debug template generation
chrome-utls-gen generate --debug --log-level debug
```

## Support

For Betanet-specific issues:

1. Check template compatibility with `--betanet-validate`
2. Verify Chrome version detection
3. Test with live Chrome instance
4. Review Betanet specification compliance

## References

- [uTLS Library](https://github.com/refraction-networking/utls)
- [Chrome TLS Fingerprinting](https://tlsfingerprint.io)
- Betanet Version 1.1 Official Implementation Specification (provided in project documentation)
