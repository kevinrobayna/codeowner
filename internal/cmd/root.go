package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:          "codeowner",
		Short:        "Print the CODEOWNERS file for a repository",
		Long:         "A CLI tool that prints and works with CODEOWNERS files.",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			fmt.Println("codeowner is not yet implemented")
			return nil
		},
	}

	root.AddCommand(newVersionCmd())

	return root
}

func Execute() error {
	return NewRootCmd().Execute()
}
