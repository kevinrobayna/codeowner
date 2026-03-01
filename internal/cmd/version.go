package cmd

import (
	"fmt"

	"github.com/kevin-robayna/codeowner/internal/appinfo"
	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Run: func(cmd *cobra.Command, _ []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "codeowner %s (commit: %s, built at: %s)\n",
				appinfo.Version, appinfo.Commit, appinfo.Date)
		},
	}
}
