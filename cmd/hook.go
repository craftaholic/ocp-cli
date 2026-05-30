package cmd

import (
	"fmt"

	"github.com/craftaholic/ocp-cli/internal/hook"
	"github.com/spf13/cobra"
)

var hookCmd = &cobra.Command{
	Use:   "hook <shell>",
	Short: "Print shell hook code",
	Long: `Print shell hook code for the specified shell (zsh, bash, or fish).
Add this to your shell's RC file to enable automatic environment variable loading.

Examples:
  # For zsh, add to ~/.zshrc:
  eval "$(ocp init hook zsh)"

  # For bash, add to ~/.bashrc:
  eval "$(ocp init hook bash)"

  # For fish, add to ~/.config/fish/config.fish:
  ocp init hook fish | source`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		shell := args[0]

		var hookCode string
		switch shell {
		case "zsh":
			hookCode = hook.GetZshHook()
		case "bash":
			hookCode = hook.GetBashHook()
		case "fish":
			hookCode = hook.GetFishHook()
		default:
			return fmt.Errorf("unsupported shell: %s (supported: zsh, bash, fish)", shell)
		}

		fmt.Print(hookCode)
		return nil
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize shell integration",
	Long:  `Initialize shell integration for ocp.`,
}

func init() {
	initCmd.AddCommand(hookCmd)
	rootCmd.AddCommand(initCmd)
}
