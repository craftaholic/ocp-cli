package cmd

import (
	"fmt"
	"os/exec"
	"syscall"

	"github.com/craftaholic/ocp/internal/config"
	"github.com/craftaholic/ocp/internal/env"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run <profile> [-- <cmd> [args...]]",
	Short: "Run a command with profile environment variables injected",
	Long: `Run a command with the specified profile's environment variables injected.
If no command is specified, defaults to running 'opencode'.

Examples:
  ocp run work
  ocp run personal -- opencode
  ocp run work -- claude --version`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profileName := args[0]

		profile, err := config.LoadProfile(profileName)
		if err != nil {
			return err
		}

		var cmdArgs []string
		dashIndex := -1
		for i, arg := range args {
			if arg == "--" {
				dashIndex = i
				break
			}
		}

		if dashIndex >= 0 && dashIndex+1 < len(args) {
			cmdArgs = args[dashIndex+1:]
		} else {
			cmdArgs = []string{"opencode"}
		}

		if len(cmdArgs) == 0 {
			return fmt.Errorf("no command specified after --")
		}

		cmdPath, err := findCommand(cmdArgs[0])
		if err != nil {
			return fmt.Errorf("command not found: %s", cmdArgs[0])
		}

		envVars := env.InjectProfileVars(profile)

		if err := syscall.Exec(cmdPath, cmdArgs, envVars); err != nil {
			return fmt.Errorf("failed to exec command: %w", err)
		}

		return nil
	},
}

func findCommand(cmd string) (string, error) {
	return exec.LookPath(cmd)
}

func init() {
	rootCmd.AddCommand(runCmd)
}
