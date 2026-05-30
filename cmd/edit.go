package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/craftaholic/ocp/internal/config"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit <profile>",
	Short: "Edit a profile in $EDITOR",
	Long:  `Open the specified profile JSON file in your $EDITOR.`,
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
	rootCmd.AddCommand(editCmd)
}
