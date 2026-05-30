package cmd

import (
	"fmt"

	"github.com/craftaholic/ocp-cli/internal/config"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <profile>",
	Short: "Delete a profile",
	Long:  `Delete the specified profile. This action cannot be undone.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profileName := args[0]

		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if cfg.Active == profileName {
			cfg.Active = ""
			if err := config.SaveConfig(cfg); err != nil {
				return fmt.Errorf("failed to update config: %w", err)
			}
		}

		if err := config.DeleteProfile(profileName); err != nil {
			return err
		}

		fmt.Printf("Deleted profile '%s'\n", profileName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
