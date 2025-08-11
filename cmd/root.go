package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	verbose bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "chrome-utls-gen",
	Short: "Chrome-Stable uTLS Template Generator for Betanet",
	Long: `Chrome-Stable uTLS Template Generator

A utility that produces exact TLS handshake bytes Chrome Stable sends, 
enabling Betanet traffic to blend in with normal web browsing.

This tool generates deterministic ClientHello templates that match Chrome 
Stable's TLS fingerprint, making it impossible for deep-packet inspectors 
to distinguish Betanet traffic from legitimate Chrome browsing.

Key Features:
- Deterministic ClientHello generation matching Chrome Stable
- JA3/JA4 self-test CLI for fingerprint verification
- Auto-refresh when new Chromium stable tags appear
- Origin mirroring support for Betanet L2 requirements
- Multi-platform support (Windows, macOS, Linux)

For more information about Betanet, see the official specification.`,
	Version: "1.0.0",
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.chrome-utls-gen.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Note: Removed toggle flag to avoid conflicts

	// Bind flags to viper
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".chrome-utls-gen" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".chrome-utls-gen")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if verbose {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}
}
