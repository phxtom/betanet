package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/betanet/chrome-utls-template-generator/internal/chrome"
	"github.com/spf13/cobra"
)

var (
	interval      string
	autoGenerate  bool
	outputDir     string
	checkInterval time.Duration
)

// monitorCmd represents the monitor command
var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitor for new Chrome releases",
	Long: `Monitor for new Chrome Stable releases and optionally auto-generate templates.

This command continuously checks for new Chrome Stable versions and can
automatically generate templates when new versions are detected.

Examples:
  chrome-utls-gen monitor --interval 1h
  chrome-utls-gen monitor --auto-generate --output ./templates/
  chrome-utls-gen monitor --interval 30m --auto-generate`,
	RunE: runMonitor,
}

func init() {
	rootCmd.AddCommand(monitorCmd)

	// Local flags
	monitorCmd.Flags().StringVarP(&interval, "interval", "i", "1h", "Check interval (e.g., 30m, 1h, 6h)")
	monitorCmd.Flags().BoolVarP(&autoGenerate, "auto-generate", "a", false, "Auto-generate templates for new versions")
	monitorCmd.Flags().StringVarP(&outputDir, "output", "o", "./templates", "Output directory for auto-generated templates")
}

func runMonitor(cmd *cobra.Command, args []string) error {
	// Parse interval
	var err error
	checkInterval, err = time.ParseDuration(interval)
	if err != nil {
		return fmt.Errorf("invalid interval format: %w", err)
	}

	if checkInterval < time.Minute {
		return fmt.Errorf("interval must be at least 1 minute")
	}

	// Create output directory if auto-generating
	if autoGenerate {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	fmt.Printf("Starting Chrome version monitor (interval: %s)\n", checkInterval)
	if autoGenerate {
		fmt.Printf("Auto-generate enabled (output: %s)\n", outputDir)
	}
	fmt.Println("Press Ctrl+C to stop")

	// Track known versions
	knownVersions := make(map[string]bool)

	// Initial check
	if err := checkForNewVersions(knownVersions, autoGenerate, outputDir); err != nil {
		fmt.Printf("Initial check failed: %v\n", err)
	}

	// Start monitoring loop
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := checkForNewVersions(knownVersions, autoGenerate, outputDir); err != nil {
				fmt.Printf("Check failed: %v\n", err)
			}
		}
	}
}

func checkForNewVersions(knownVersions map[string]bool, autoGenerate bool, outputDir string) error {
	// Get latest stable version
	latest, err := chrome.GetLatestStableVersion()
	if err != nil {
		return fmt.Errorf("failed to get latest version: %w", err)
	}

	// Check if this is a new version
	if !knownVersions[latest] {
		fmt.Printf("[%s] New Chrome version detected: %s\n", time.Now().Format("2006-01-02 15:04:05"), latest)
		knownVersions[latest] = true

		// Auto-generate template if enabled
		if autoGenerate {
			if err := generateTemplateForVersion(latest, outputDir); err != nil {
				fmt.Printf("Failed to generate template for %s: %v\n", latest, err)
			} else {
				fmt.Printf("Template generated for Chrome %s\n", latest)
			}
		}
	} else {
		fmt.Printf("[%s] No new versions (latest: %s)\n", time.Now().Format("2006-01-02 15:04:05"), latest)
	}

	return nil
}

func generateTemplateForVersion(version, outputDir string) error {
	// Import the generate command logic here
	// For now, we'll just create a placeholder
	templatePath := filepath.Join(outputDir, fmt.Sprintf("chrome-%s.json", version))
	
	// Check if template already exists
	if _, err := os.Stat(templatePath); err == nil {
		return fmt.Errorf("template already exists: %s", templatePath)
	}

	// This would call the actual template generation logic
	// For now, create a placeholder file
	placeholder := fmt.Sprintf(`{
  "version": "%s",
  "timestamp": "%s",
  "status": "placeholder - run 'chrome-utls-gen generate --version %s' to generate full template"
}`, version, time.Now().UTC().Format(time.RFC3339), version)

	return os.WriteFile(templatePath, []byte(placeholder), 0644)
}
