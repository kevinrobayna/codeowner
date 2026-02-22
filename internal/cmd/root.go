package cmd

import (
	"fmt"
	"os"

	"github.com/kevin-robayna/codeowner/internal/owner"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	var prefix string

	root := &cobra.Command{
		Use:          "codeowner [path]",
		Short:        "Print the CODEOWNERS file for a repository",
		Long:         "A CLI tool that scans source files for CodeOwner annotations and prints a CODEOWNERS file.",
		SilenceUsage: true,
		Args:         cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			dir := "."
			if len(args) > 0 {
				dir = args[0]
			}

			mappings, err := owner.ParseDir(dir, prefix)
			if err != nil {
				return fmt.Errorf("scanning directory: %w", err)
			}

			if len(mappings) == 0 {
				fmt.Fprintln(os.Stderr, "no CodeOwner annotations found")
				return nil
			}

			fmt.Print(owner.FormatCodeOwners(mappings))
			return nil
		},
	}

	root.Flags().StringVar(&prefix, "prefix", owner.DefaultPrefix, "annotation prefix to search for")
	root.AddCommand(newVersionCmd())

	return root
}

func Execute() error {
	return NewRootCmd().Execute()
}
