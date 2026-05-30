package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/craftaholic/ocp-cli/internal/config"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <profile>",
	Short: "Create a new profile",
	Long:  `Create a new empty profile and open it in $EDITOR.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profileName := args[0]

		exists, err := config.ProfileExists(profileName)
		if err != nil {
			return fmt.Errorf("failed to check profile: %w", err)
		}
		if exists {
			return fmt.Errorf("profile '%s' already exists", profileName)
		}

		profile := &config.Profile{
			Name: profileName,
			Vars: make(map[string]string),
		}

		if err := config.SaveProfile(profile); err != nil {
			return fmt.Errorf("failed to create profile: %w", err)
		}

		fmt.Printf("Created profile '%s'\n", profileName)

		profileDir, err := config.GetProfileDir(profileName)
		if err != nil {
			return fmt.Errorf("failed to get profile directory: %w", err)
		}

		profilePath := filepath.Join(profileDir, "profile.json")

		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vi"
		}

		editorCmd := exec.Command(editor, profilePath)
		editorCmd.Stdin = os.Stdin
		editorCmd.Stdout = os.Stdout
		editorCmd.Stderr = os.Stderr

		if err := editorCmd.Run(); err != nil {
			return fmt.Errorf("failed to run editor: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
