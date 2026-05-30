package cmd

import (
	"fmt"
	"sort"

	"github.com/craftaholic/ocp/internal/config"
	"github.com/craftaholic/ocp/internal/env"
	"github.com/spf13/cobra"
)

var (
	exportFlag   bool
	nameOnlyFlag bool
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show active profile and its environment variables",
	Long:  `Display the currently active profile and its environment variables. Sensitive values are masked.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if cfg.Active == "" {
			if !exportFlag && !nameOnlyFlag {
				fmt.Println("No active profile set. Use 'ocp use <profile>' to set one.")
			}
			return nil
		}

		if nameOnlyFlag {
			fmt.Println(cfg.Active)
			return nil
		}

		profile, err := config.LoadProfile(cfg.Active)
		if err != nil {
			return err
		}

		if exportFlag {
			var keys []string
			for k := range profile.Vars {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			for _, key := range keys {
				value := env.ExpandPath(profile.Vars[key])
				fmt.Printf("export %s=%q\n", key, value)
			}
			return nil
		}

		fmt.Printf("Active profile: %s\n", cfg.Active)
		fmt.Println("\nEnvironment variables:")

		if len(profile.Vars) == 0 {
			fmt.Println("  (none)")
			return nil
		}

		var keys []string
		for k := range profile.Vars {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, key := range keys {
			value := profile.Vars[key]
			displayValue := value

			if env.IsSensitive(key) {
				displayValue = env.MaskValue(value)
			}

			fmt.Printf("  %s=%s\n", key, displayValue)
		}

		return nil
	},
}

func init() {
	statusCmd.Flags().BoolVar(&exportFlag, "export", false, "Output in export format for shell eval")
	statusCmd.Flags().BoolVar(&nameOnlyFlag, "name-only", false, "Output only the active profile name")
	rootCmd.AddCommand(statusCmd)
}
