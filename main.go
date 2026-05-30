package main

import (
	"os"

	"github.com/craftaholic/ocp-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
