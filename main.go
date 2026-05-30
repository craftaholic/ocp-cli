package main

import (
	"os"

	"github.com/craftaholic/ocp/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
