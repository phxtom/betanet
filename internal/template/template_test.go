package template

import (
	"encoding/json"
	"testing"
)

func TestGenerateChromeTemplate(t *testing.T) {
	version := "120.0.6099.109"
	
	template, err := GenerateChromeTemplate(version)
	if err != nil {
		t.Fatalf("Failed to generate template: %v", err)
	}

	// Check basic structure
	if template.Version != version {
		t.Errorf("Expected version %s, got %s", version, template.Version)
	}

	if template.ClientHello.Version == "" {
		t.Error("ClientHello version should not be empty")
	}

	if template.ClientHello.Random == "" {
		t.Error("ClientHello random should not be empty")
	}

	if template.ClientHello.SessionID == "" {
		t.Error("ClientHello session ID should not be empty")
	}

	if len(template.ClientHello.CipherSuites) == 0 {
		t.Error("ClientHello should have cipher suites")
	}

	if len(template.ClientHello.Extensions) == 0 {
		t.Error("ClientHello should have extensions")
	}

	// Check fingerprints
	if template.JA3Fingerprint == "" {
		t.Error("JA3 fingerprint should not be empty")
	}

	if template.JA4Fingerprint == "" {
		t.Error("JA4 fingerprint should not be empty")
	}

	// Test JSON marshaling
	_, err = json.Marshal(template)
	if err != nil {
		t.Errorf("Failed to marshal template to JSON: %v", err)
	}
}

func TestGenerateChromeTemplateUnsupportedVersion(t *testing.T) {
	version := "999.0.0.0"
	
	_, err := GenerateChromeTemplate(version)
	if err == nil {
		t.Error("Expected error for unsupported version")
	}
}

func TestGenerateClientHello(t *testing.T) {
	version := "120.0.6099.109"
	
	template, err := GenerateChromeTemplate(version)
	if err != nil {
		t.Fatalf("Failed to generate template: %v", err)
	}

	clientHello, err := GenerateClientHello(template)
	if err != nil {
		t.Fatalf("Failed to generate ClientHello: %v", err)
	}

	// Check minimum length
	if len(clientHello) < 50 {
		t.Errorf("ClientHello too short: %d bytes", len(clientHello))
	}

	// Check TLS record header
	if clientHello[0] != 0x16 {
		t.Error("Expected TLS handshake record type")
	}

	// Check TLS version
	if clientHello[1] != 0x03 || clientHello[2] != 0x03 {
		t.Error("Expected TLS 1.2 version")
	}

	// Check handshake type
	if clientHello[5] != 0x01 {
		t.Error("Expected ClientHello handshake type")
	}
}

func TestChromeConfigs(t *testing.T) {
	// Test that we have configurations for supported versions
	supportedVersions := []string{"120"}
	
	for _, version := range supportedVersions {
		config, exists := chromeConfigs[version]
		if !exists {
			t.Errorf("Missing configuration for Chrome version %s", version)
			continue
		}

		if config.TLSVersion == 0 {
			t.Errorf("Invalid TLS version for Chrome %s", version)
		}

		if len(config.CipherSuites) == 0 {
			t.Errorf("No cipher suites for Chrome %s", version)
		}

		if len(config.Extensions) == 0 {
			t.Errorf("No extensions for Chrome %s", version)
		}
	}
}
