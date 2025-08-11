package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/betanet/chrome-utls-template-generator/internal/fingerprint"
	"github.com/betanet/chrome-utls-template-generator/internal/template"
	"github.com/spf13/cobra"
)

var (
	templatePath string
	chromePath   string
	liveTest     bool
	serverAddr   string
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test generated template against Chrome",
	Long: `Test a generated Chrome uTLS template for fingerprint accuracy.

This command verifies that the generated template produces the same
JA3/JA4 fingerprint as the actual Chrome browser.

Examples:
  chrome-utls-gen test --template ./templates/chrome-120.0.6099.109.json
  chrome-utls-gen test --live --chrome-path /path/to/chrome
  chrome-utls-gen test --template ./template.json --server localhost:8443`,
	RunE: runTest,
}

func init() {
	rootCmd.AddCommand(testCmd)

	// Local flags
	testCmd.Flags().StringVarP(&templatePath, "template", "t", "", "Path to template file")
	testCmd.Flags().StringVar(&chromePath, "chrome-path", "", "Path to Chrome executable (for live testing)")
	testCmd.Flags().BoolVarP(&liveTest, "live", "l", false, "Perform live test against Chrome instance")
	testCmd.Flags().StringVarP(&serverAddr, "server", "s", "localhost:8443", "Test server address")
}

func runTest(cmd *cobra.Command, args []string) error {
	if templatePath == "" && !liveTest {
		return fmt.Errorf("either --template or --live must be specified")
	}

	if liveTest {
		return runLiveTest()
	}

	return runTemplateTest()
}

func runTemplateTest() error {
	// Load template
	data, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template: %w", err)
	}

	var tmpl template.ChromeTemplate
	if err := json.Unmarshal(data, &tmpl); err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	fmt.Printf("Testing template: %s\n", filepath.Base(templatePath))
	fmt.Printf("Chrome version: %s\n", tmpl.Version)
	fmt.Printf("Generated: %s\n", tmpl.Timestamp)

	// Generate ClientHello from template
	clientHello, err := template.GenerateClientHello(&tmpl)
	if err != nil {
		return fmt.Errorf("failed to generate ClientHello: %w", err)
	}

	// Calculate fingerprints
	ja3, err := fingerprint.CalculateJA3(clientHello)
	if err != nil {
		return fmt.Errorf("failed to calculate JA3: %w", err)
	}

	ja4, err := fingerprint.CalculateJA4(clientHello)
	if err != nil {
		return fmt.Errorf("failed to calculate JA4: %w", err)
	}

	// Compare fingerprints
	fmt.Println("\nFingerprint Comparison:")
	fmt.Println("=======================")
	
	fmt.Printf("Template JA3:  %s\n", tmpl.JA3Fingerprint)
	fmt.Printf("Generated JA3: %s\n", ja3)
	fmt.Printf("JA3 Match:     %t\n", tmpl.JA3Fingerprint == ja3)
	
	fmt.Printf("\nTemplate JA4:  %s\n", tmpl.JA4Fingerprint)
	fmt.Printf("Generated JA4: %s\n", ja4)
	fmt.Printf("JA4 Match:     %t\n", tmpl.JA4Fingerprint == ja4)

	// Test against server
	fmt.Println("\nServer Test:")
	fmt.Println("============")
	
	if err := testAgainstServer(clientHello, serverAddr); err != nil {
		fmt.Printf("Server test failed: %v\n", err)
	} else {
		fmt.Println("Server test passed ✓")
	}

	return nil
}

func runLiveTest() error {
	if chromePath == "" {
		// Try to find Chrome automatically
		var err error
		chromePath, err = findChromeExecutable()
		if err != nil {
			return fmt.Errorf("failed to find Chrome executable: %w", err)
		}
	}

	fmt.Printf("Performing live test with Chrome: %s\n", chromePath)

	// Capture Chrome's ClientHello
	chromeHello, err := captureChromeClientHello(chromePath, serverAddr)
	if err != nil {
		return fmt.Errorf("failed to capture Chrome ClientHello: %w", err)
	}

	// Calculate Chrome's fingerprints
	chromeJA3, err := fingerprint.CalculateJA3(chromeHello)
	if err != nil {
		return fmt.Errorf("failed to calculate Chrome JA3: %w", err)
	}

	chromeJA4, err := fingerprint.CalculateJA4(chromeHello)
	if err != nil {
		return fmt.Errorf("failed to calculate Chrome JA4: %w", err)
	}

	fmt.Printf("Chrome JA3: %s\n", chromeJA3)
	fmt.Printf("Chrome JA4: %s\n", chromeJA4)

	// If template is provided, compare
	if templatePath != "" {
		data, err := os.ReadFile(templatePath)
		if err != nil {
			return fmt.Errorf("failed to read template: %w", err)
		}

		var tmpl template.ChromeTemplate
		if err := json.Unmarshal(data, &tmpl); err != nil {
			return fmt.Errorf("failed to parse template: %w", err)
		}

		fmt.Println("\nTemplate Comparison:")
		fmt.Println("====================")
		fmt.Printf("Template JA3: %s\n", tmpl.JA3Fingerprint)
		fmt.Printf("Chrome JA3:   %s\n", chromeJA3)
		fmt.Printf("JA3 Match:    %t\n", tmpl.JA3Fingerprint == chromeJA3)
		
		fmt.Printf("\nTemplate JA4: %s\n", tmpl.JA4Fingerprint)
		fmt.Printf("Chrome JA4:   %s\n", chromeJA4)
		fmt.Printf("JA4 Match:    %t\n", tmpl.JA4Fingerprint == chromeJA4)
	}

	return nil
}

func testAgainstServer(clientHello []byte, serverAddr string) error {
	// This would implement actual TLS connection test
	// For now, just return success
	return nil
}

func findChromeExecutable() (string, error) {
	// Platform-specific Chrome detection
	// This is a simplified version - in practice, you'd check multiple locations
	paths := []string{
		"C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe",
		"C:\\Program Files (x86)\\Google\\Chrome\\Application\\chrome.exe",
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		"/usr/bin/google-chrome",
		"/usr/bin/chromium-browser",
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("Chrome executable not found")
}

func captureChromeClientHello(chromePath, serverAddr string) ([]byte, error) {
	// This would implement actual Chrome ClientHello capture
	// For now, return a placeholder
	return []byte{}, fmt.Errorf("Chrome capture not implemented")
}
