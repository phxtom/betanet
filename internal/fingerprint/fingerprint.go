package fingerprint

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

// CalculateJA3 calculates the JA3 fingerprint from ClientHello bytes
func CalculateJA3(clientHello []byte) (string, error) {
	if len(clientHello) < 5 {
		return "", fmt.Errorf("ClientHello too short")
	}

	// Parse TLS record header
	if clientHello[0] != 0x16 {
		return "", fmt.Errorf("not a handshake record")
	}

	// Parse handshake header
	if clientHello[5] != 0x01 {
		return "", fmt.Errorf("not a ClientHello")
	}

	// Extract components
	components, err := extractJA3Components(clientHello)
	if err != nil {
		return "", fmt.Errorf("failed to extract components: %w", err)
	}

	// Build JA3 string
	ja3Parts := []string{
		strconv.Itoa(components.TLSVersion),
		strings.Join(components.CipherSuites, ","),
		strings.Join(components.Extensions, ","),
		strings.Join(components.EllipticCurves, ","),
		strings.Join(components.EllipticCurvePointFormats, ","),
	}

	ja3String := strings.Join(ja3Parts, ",")

	// Calculate MD5 hash
	hash := md5.Sum([]byte(ja3String))
	return fmt.Sprintf("%x", hash), nil
}

// CalculateJA4 calculates the JA4 fingerprint from ClientHello bytes
func CalculateJA4(clientHello []byte) (string, error) {
	if len(clientHello) < 5 {
		return "", fmt.Errorf("ClientHello too short")
	}

	// Parse TLS record header
	if clientHello[0] != 0x16 {
		return "", fmt.Errorf("not a handshake record")
	}

	// Parse handshake header
	if clientHello[5] != 0x01 {
		return "", fmt.Errorf("not a ClientHello")
	}

	// Extract components
	components, err := extractJA4Components(clientHello)
	if err != nil {
		return "", fmt.Errorf("failed to extract components: %w", err)
	}

	// Build JA4 string
	ja4Parts := []string{
		"t" + strconv.Itoa(components.TLSVersion),
		"d" + strconv.Itoa(components.DTLSVersion),
		"h" + strconv.Itoa(components.HandshakeVersion),
		"c" + strings.Join(components.CipherSuites, "-"),
		"e" + strings.Join(components.Extensions, "-"),
		"g" + strings.Join(components.EllipticCurves, "-"),
		"p" + strings.Join(components.EllipticCurvePointFormats, "-"),
		"a" + strings.Join(components.SignatureAlgorithms, "-"),
		"r" + strings.Join(components.RenegotiationInfo, "-"),
		"u" + strings.Join(components.UnknownExtensions, "-"),
	}

	ja4String := strings.Join(ja4Parts, "_")

	// Calculate MD5 hash
	hash := md5.Sum([]byte(ja4String))
	return fmt.Sprintf("%s_%x", ja4String, hash), nil
}

// JA3Components represents the components used in JA3 calculation
type JA3Components struct {
	TLSVersion                int
	CipherSuites             []string
	Extensions               []string
	EllipticCurves           []string
	EllipticCurvePointFormats []string
}

// JA4Components represents the components used in JA4 calculation
type JA4Components struct {
	TLSVersion                int
	DTLSVersion              int
	HandshakeVersion         int
	CipherSuites             []string
	Extensions               []string
	EllipticCurves           []string
	EllipticCurvePointFormats []string
	SignatureAlgorithms      []string
	RenegotiationInfo        []string
	UnknownExtensions        []string
}

// extractJA3Components extracts components from ClientHello for JA3 calculation
func extractJA3Components(clientHello []byte) (*JA3Components, error) {
	components := &JA3Components{}

	// Skip TLS record header (5 bytes)
	// Skip handshake header (4 bytes)
	pos := 9

	// Extract TLS version
	if pos+2 > len(clientHello) {
		return nil, fmt.Errorf("insufficient data for TLS version")
	}
	components.TLSVersion = int(binary.BigEndian.Uint16(clientHello[pos:pos+2]))
	pos += 2

	// Skip random (32 bytes)
	pos += 32

	// Skip session ID
	if pos >= len(clientHello) {
		return nil, fmt.Errorf("insufficient data for session ID length")
	}
	sessionIDLen := int(clientHello[pos])
	pos++
	pos += sessionIDLen

	// Extract cipher suites
	if pos+2 > len(clientHello) {
		return nil, fmt.Errorf("insufficient data for cipher suites length")
	}
	cipherSuitesLen := int(binary.BigEndian.Uint16(clientHello[pos:pos+2]))
	pos += 2

	if pos+cipherSuitesLen > len(clientHello) {
		return nil, fmt.Errorf("insufficient data for cipher suites")
	}

	for i := 0; i < cipherSuitesLen; i += 2 {
		if pos+i+2 > len(clientHello) {
			break
		}
		suite := binary.BigEndian.Uint16(clientHello[pos+i : pos+i+2])
		components.CipherSuites = append(components.CipherSuites, strconv.Itoa(int(suite)))
	}
	pos += cipherSuitesLen

	// Skip compression methods
	if pos >= len(clientHello) {
		return nil, fmt.Errorf("insufficient data for compression methods length")
	}
	compressionMethodsLen := int(clientHello[pos])
	pos++
	pos += compressionMethodsLen

	// Extract extensions
	if pos+2 > len(clientHello) {
		return nil, fmt.Errorf("insufficient data for extensions length")
	}
	extensionsLen := int(binary.BigEndian.Uint16(clientHello[pos:pos+2]))
	pos += 2

	extensionsEnd := pos + extensionsLen
	for pos < extensionsEnd && pos+4 <= len(clientHello) {
		extType := binary.BigEndian.Uint16(clientHello[pos:pos+2])
		extLen := int(binary.BigEndian.Uint16(clientHello[pos+2:pos+4]))
		pos += 4

		components.Extensions = append(components.Extensions, strconv.Itoa(int(extType)))

		// Parse specific extensions
		switch extType {
		case 0x000a: // supported_groups
			if pos+2 <= len(clientHello) && pos+2 <= pos+extLen {
				groupsLen := int(binary.BigEndian.Uint16(clientHello[pos:pos+2]))
				pos += 2
				for i := 0; i < groupsLen && pos+i+2 <= len(clientHello) && pos+i+2 <= pos+extLen; i += 2 {
					group := binary.BigEndian.Uint16(clientHello[pos+i : pos+i+2])
					components.EllipticCurves = append(components.EllipticCurves, strconv.Itoa(int(group)))
				}
			}
		case 0x000b: // ec_point_formats
			if pos+1 <= len(clientHello) && pos+1 <= pos+extLen {
				formatsLen := int(clientHello[pos])
				pos++
				for i := 0; i < formatsLen && pos+i < len(clientHello) && pos+i < pos+extLen; i++ {
					format := int(clientHello[pos+i])
					components.EllipticCurvePointFormats = append(components.EllipticCurvePointFormats, strconv.Itoa(format))
				}
			}
		}

		pos += extLen
	}

	return components, nil
}

// extractJA4Components extracts components from ClientHello for JA4 calculation
func extractJA4Components(clientHello []byte) (*JA4Components, error) {
	components := &JA4Components{}

	// Extract basic components (similar to JA3)
	ja3Components, err := extractJA3Components(clientHello)
	if err != nil {
		return nil, err
	}

	components.TLSVersion = ja3Components.TLSVersion
	components.CipherSuites = ja3Components.CipherSuites
	components.Extensions = ja3Components.Extensions
	components.EllipticCurves = ja3Components.EllipticCurves
	components.EllipticCurvePointFormats = ja3Components.EllipticCurvePointFormats

	// JA4-specific processing
	// For now, we'll use the same data as JA3
	// In practice, JA4 would have more detailed parsing

	return components, nil
}

// ParseJA3 parses a JA3 fingerprint string
func ParseJA3(ja3 string) (*JA3Components, error) {
	parts := strings.Split(ja3, ",")
	if len(parts) < 5 {
		return nil, fmt.Errorf("invalid JA3 format")
	}

	components := &JA3Components{}

	// Parse TLS version
	tlsVersion, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid TLS version: %w", err)
	}
	components.TLSVersion = tlsVersion

	// Parse cipher suites
	if parts[1] != "" {
		components.CipherSuites = strings.Split(parts[1], ",")
	}

	// Parse extensions
	if parts[2] != "" {
		components.Extensions = strings.Split(parts[2], ",")
	}

	// Parse elliptic curves
	if parts[3] != "" {
		components.EllipticCurves = strings.Split(parts[3], ",")
	}

	// Parse elliptic curve point formats
	if parts[4] != "" {
		components.EllipticCurvePointFormats = strings.Split(parts[4], ",")
	}

	return components, nil
}

// ParseJA4 parses a JA4 fingerprint string
func ParseJA4(ja4 string) (*JA4Components, error) {
	// Remove hash suffix if present
	parts := strings.Split(ja4, "_")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid JA4 format")
	}

	// Remove hash from last part
	lastPart := parts[len(parts)-1]
	if len(lastPart) == 32 && isHex(lastPart) {
		parts = parts[:len(parts)-1]
	}

	components := &JA4Components{}

	for _, part := range parts {
		if len(part) < 2 {
			continue
		}

		prefix := part[0]
		value := part[1:]

		switch prefix {
		case 't':
			if v, err := strconv.Atoi(value); err == nil {
				components.TLSVersion = v
			}
		case 'd':
			if v, err := strconv.Atoi(value); err == nil {
				components.DTLSVersion = v
			}
		case 'h':
			if v, err := strconv.Atoi(value); err == nil {
				components.HandshakeVersion = v
			}
		case 'c':
			if value != "" {
				components.CipherSuites = strings.Split(value, "-")
			}
		case 'e':
			if value != "" {
				components.Extensions = strings.Split(value, "-")
			}
		case 'g':
			if value != "" {
				components.EllipticCurves = strings.Split(value, "-")
			}
		case 'p':
			if value != "" {
				components.EllipticCurvePointFormats = strings.Split(value, "-")
			}
		case 'a':
			if value != "" {
				components.SignatureAlgorithms = strings.Split(value, "-")
			}
		case 'r':
			if value != "" {
				components.RenegotiationInfo = strings.Split(value, "-")
			}
		case 'u':
			if value != "" {
				components.UnknownExtensions = strings.Split(value, "-")
			}
		}
	}

	return components, nil
}

// isHex checks if a string is a valid hexadecimal string
func isHex(s string) bool {
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}
