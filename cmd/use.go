package cmd

import (
	"fmt"

	"github.com/craftaholic/ocp-cli/internal/config"
	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use <profile>",
	Short: "Set the active profile",
	Long:  `Set the active profile. This profile will be used by default for subsequent operations.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profileName := args[0]

		exists, err := config.ProfileExists(profileName)
		if err != nil {
			return fmt.Errorf("failed to check profile: %w", err)
		}
		if !exists {
			return fmt.Errorf("profile '%s' does not exist", profileName)
		}

		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		cfg.Active = profileName

		if err := config.SaveConfig(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		if err := config.UpdateSymlink(profileName); err != nil {
			return fmt.Errorf("failed to update symlink: %w", err)
		}

		fmt.Printf("Switched to profile '%s'\n", profileName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(useCmd)
}
