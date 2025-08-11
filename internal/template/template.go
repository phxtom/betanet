package template

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
)

// ChromeTemplate represents a Chrome uTLS template
type ChromeTemplate struct {
	Version        string                 `json:"version"`
	Timestamp      string                 `json:"timestamp"`
	ClientHello    ClientHelloTemplate    `json:"client_hello"`
	JA3Fingerprint string                 `json:"ja3_fingerprint"`
	JA4Fingerprint string                 `json:"ja4_fingerprint"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// ClientHelloTemplate represents the ClientHello structure
type ClientHelloTemplate struct {
	Version            string              `json:"version"`
	Random             string              `json:"random"`
	SessionID          string              `json:"session_id"`
	CipherSuites       []string            `json:"cipher_suites"`
	CompressionMethods []int               `json:"compression_methods"`
	Extensions         []ExtensionTemplate `json:"extensions"`
}

// ExtensionTemplate represents a TLS extension
type ExtensionTemplate struct {
	Type int    `json:"type"`
	Data string `json:"data"`
}

// Chrome version-specific configurations
var ChromeConfigs = map[string]ChromeConfig{
	"120": {
		TLSVersion: 0x0303, // TLS 1.2
		CipherSuites: []uint16{
			0x1301, 0x1302, 0x1303, 0xc02f, 0xc02b, 0xc030, 0xc02c, 0xc027, 0xc028, 0xc014, 0xc013, 0xc011, 0xc012, 0x009c, 0x009d, 0x002f, 0x0035,
		},
		Extensions: []ExtensionConfig{
			{Type: 0x0000, Name: "server_name"},
			{Type: 0x0017, Name: "status_request"},
			{Type: 0xff01, Name: "renegotiation_info"},
			{Type: 0x000a, Name: "supported_groups"},
			{Type: 0x000b, Name: "ec_point_formats"},
			{Type: 0x000d, Name: "signature_algorithms"},
			{Type: 0x0010, Name: "application_layer_protocol_negotiation"},
			{Type: 0x001b, Name: "extended_master_secret"},
			{Type: 0x0018, Name: "signed_certificate_timestamp"},
			{Type: 0x0022, Name: "supported_versions"},
			{Type: 0x0029, Name: "key_share"},
		},
		CompressionMethods: []int{0},
	},
}

// ChromeConfig represents Chrome version-specific configuration
type ChromeConfig struct {
	TLSVersion         uint16
	CipherSuites       []uint16
	Extensions         []ExtensionConfig
	CompressionMethods []int
}

// ExtensionConfig represents extension configuration
type ExtensionConfig struct {
	Type int    `json:"type"`
	Name string `json:"name"`
}

// GenerateChromeTemplate generates a Chrome uTLS template for the given version
func GenerateChromeTemplate(version string) (*ChromeTemplate, error) {
	// Extract major version
	majorVersion := strings.Split(version, ".")[0]

	config, exists := ChromeConfigs[majorVersion]
	if !exists {
		return nil, fmt.Errorf("unsupported Chrome version: %s", version)
	}

	// Generate random data
	random := make([]byte, 32)
	if _, err := rand.Read(random); err != nil {
		return nil, fmt.Errorf("failed to generate random: %w", err)
	}

	sessionID := make([]byte, 32)
	if _, err := rand.Read(sessionID); err != nil {
		return nil, fmt.Errorf("failed to generate session ID: %w", err)
	}

	// Build template
	template := &ChromeTemplate{
		Version: version,
		ClientHello: ClientHelloTemplate{
			Version:            fmt.Sprintf("TLS 1.%d", config.TLSVersion&0xFF),
			Random:             base64.StdEncoding.EncodeToString(random),
			SessionID:          base64.StdEncoding.EncodeToString(sessionID),
			CipherSuites:       make([]string, len(config.CipherSuites)),
			CompressionMethods: config.CompressionMethods,
			Extensions:         make([]ExtensionTemplate, len(config.Extensions)),
		},
	}

	// Convert cipher suites to hex strings
	for i, suite := range config.CipherSuites {
		template.ClientHello.CipherSuites[i] = fmt.Sprintf("0x%04x", suite)
	}

	// Build extensions
	for i, ext := range config.Extensions {
		template.ClientHello.Extensions[i] = ExtensionTemplate{
			Type: ext.Type,
			Data: GenerateExtensionData(ext),
		}
	}

	// Set timestamp
	template.Timestamp = "2025-08-11T10:30:00Z"

	return template, nil
}

// GenerateClientHello generates the actual ClientHello bytes from a template
func GenerateClientHello(template *ChromeTemplate) ([]byte, error) {
	// This is a simplified implementation
	// In practice, this would generate the exact TLS ClientHello structure

	// Build the handshake payload first
	handshake := make([]byte, 0, 1024)

	// Client Version
	handshake = append(handshake, 0x03, 0x03) // TLS 1.2

	// Random
	random, err := base64.StdEncoding.DecodeString(template.ClientHello.Random)
	if err != nil {
		return nil, fmt.Errorf("invalid random data: %w", err)
	}
	handshake = append(handshake, random...)

	// Session ID
	sessionID, err := base64.StdEncoding.DecodeString(template.ClientHello.SessionID)
	if err != nil {
		return nil, fmt.Errorf("invalid session ID: %w", err)
	}
	handshake = append(handshake, byte(len(sessionID)))
	handshake = append(handshake, sessionID...)

	// Cipher Suites
	cipherSuites := make([]byte, 0, len(template.ClientHello.CipherSuites)*2)
	for _, suiteStr := range template.ClientHello.CipherSuites {
		var suite uint16
		if _, err := fmt.Sscanf(suiteStr, "0x%x", &suite); err != nil {
			return nil, fmt.Errorf("invalid cipher suite: %s", suiteStr)
		}
		cipherSuites = append(cipherSuites, byte(suite>>8), byte(suite))
	}
	handshake = append(handshake, byte(len(cipherSuites)>>8), byte(len(cipherSuites)))
	handshake = append(handshake, cipherSuites...)

	// Compression Methods
	handshake = append(handshake, byte(len(template.ClientHello.CompressionMethods)))
	for _, method := range template.ClientHello.CompressionMethods {
		handshake = append(handshake, byte(method))
	}

	// Extensions
	extensions := make([]byte, 0)
	for _, ext := range template.ClientHello.Extensions {
		extData, err := base64.StdEncoding.DecodeString(ext.Data)
		if err != nil {
			return nil, fmt.Errorf("invalid extension data: %w", err)
		}

		extensions = append(extensions, byte(ext.Type>>8), byte(ext.Type))
		extensions = append(extensions, byte(len(extData)>>8), byte(len(extData)))
		extensions = append(extensions, extData...)
	}

	handshake = append(handshake, byte(len(extensions)>>8), byte(len(extensions)))
	handshake = append(handshake, extensions...)

	// Now build the complete TLS record
	result := make([]byte, 0, len(handshake)+9)

	// TLS Record Header
	result = append(result, 0x16)                                          // Handshake
	result = append(result, 0x03, 0x03)                                    // TLS 1.2
	result = append(result, byte(len(handshake)>>8), byte(len(handshake))) // Record length

	// Handshake Header
	result = append(result, 0x01)                                                                                              // ClientHello
	result = append(result, byte(len(handshake)>>24), byte(len(handshake)>>16), byte(len(handshake)>>8), byte(len(handshake))) // Handshake length

	// Handshake payload
	result = append(result, handshake...)

	return result, nil
}

// generateExtensionData generates extension-specific data
func GenerateExtensionData(ext ExtensionConfig) string {
	// This is a simplified implementation
	// In practice, this would generate the actual extension data based on type

	switch ext.Type {
	case 0x0000: // server_name
		// SNI extension data
		data := []byte{0x00, 0x00, 0x00, 0x00} // Name type: hostname, length placeholder
		return base64.StdEncoding.EncodeToString(data)

	case 0x000a: // supported_groups
		// Supported groups: x25519, secp256r1, secp384r1
		data := []byte{0x00, 0x06}      // Length
		data = append(data, 0x00, 0x1d) // x25519
		data = append(data, 0x00, 0x17) // secp256r1
		data = append(data, 0x00, 0x18) // secp384r1
		return base64.StdEncoding.EncodeToString(data)

	case 0x000b: // ec_point_formats
		// EC point formats: uncompressed
		data := []byte{0x01, 0x00} // Length: 1, Format: uncompressed
		return base64.StdEncoding.EncodeToString(data)

	case 0x000d: // signature_algorithms
		// Signature algorithms
		data := []byte{0x00, 0x08}      // Length
		data = append(data, 0x08, 0x07) // rsa_pss_rsae_sha256
		data = append(data, 0x08, 0x08) // rsa_pss_rsae_sha384
		data = append(data, 0x08, 0x09) // rsa_pss_rsae_sha512
		data = append(data, 0x08, 0x04) // rsa_pkcs1_sha256
		return base64.StdEncoding.EncodeToString(data)

	case 0x0010: // application_layer_protocol_negotiation
		// ALPN: h2, http/1.1
		data := []byte{0x00, 0x05}                                                // Length
		data = append(data, 0x02, 0x68, 0x32)                                     // h2
		data = append(data, 0x08, 0x68, 0x74, 0x74, 0x70, 0x2f, 0x31, 0x2e, 0x31) // http/1.1
		return base64.StdEncoding.EncodeToString(data)

	case 0x0017: // status_request
		// OCSP status request
		data := []byte{0x01, 0x00, 0x00, 0x00, 0x00} // Status type: OCSP, responder ID list: empty, request extensions: empty
		return base64.StdEncoding.EncodeToString(data)

	case 0xff01: // renegotiation_info
		// Renegotiation info
		data := []byte{0x01, 0x00} // Length: 1, renegotiated_connection: empty
		return base64.StdEncoding.EncodeToString(data)

	default:
		// Default minimal extension data - ensure it's valid base64
		return base64.StdEncoding.EncodeToString([]byte{0x00})
	}
}
