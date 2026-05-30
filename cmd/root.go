package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "ocp",
	Version: "1.0.0",
	Short:   "OpenCode Profile Switcher - manage multiple named profiles for opencode and claude",
	Long: `ocp is a CLI tool for managing multiple named profiles for opencode and claude code.
Each profile contains a set of environment variables (API keys, config directories, model preferences).
Switch between profiles seamlessly in your shell without manually exporting variables.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}
	return nil
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
