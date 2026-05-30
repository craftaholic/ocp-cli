package cmd

import (
	"fmt"

	"github.com/craftaholic/ocp-cli/internal/config"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all profiles",
	Long:  `List all available profiles. The active profile is marked with *.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		profiles, err := config.ListProfiles()
		if err != nil {
			return fmt.Errorf("failed to list profiles: %w", err)
		}

		if len(profiles) == 0 {
			fmt.Println("No profiles found. Create one with 'ocp add <profile>'")
			return nil
		}

		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		for _, profile := range profiles {
			if profile == cfg.Active {
				fmt.Printf("* %s\n", profile)
			} else {
				fmt.Printf("  %s\n", profile)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
