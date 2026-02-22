package main

import (
	"os"

	"github.com/kevin-robayna/codeowner/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
