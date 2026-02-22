package cmd

import (
	"fmt"
	"os"

	"github.com/kevin-robayna/codeowner/internal/formatter"
	"github.com/kevin-robayna/codeowner/internal/scanning"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	var prefix string
	var dirOwner string
	var protect string

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

			mappings, err := scanning.ParseDir(dir, prefix, dirOwner)
			if err != nil {
				return fmt.Errorf("scanning directory: %w", err)
			}

			if protect != "" {
				pm, pErr := scanning.ParseProtect(protect)
				if pErr != nil {
					return fmt.Errorf("--protect: %w", pErr)
				}
				mappings = append(mappings, pm)
			}

			if len(mappings) == 0 {
				fmt.Fprintln(os.Stderr, "no CodeOwner annotations found")
				return nil
			}

			fmt.Print(formatter.CodeOwners(mappings))
			return nil
		},
	}

	root.Flags().StringVar(&prefix, "prefix", scanning.DefaultPrefix, "annotation prefix to search for")
	root.Flags().StringVar(&dirOwner, "dirowner", scanning.CodeOwnerFile, "filename for directory-level ownership")
	root.Flags().StringVar(&protect, "protect", "", "owners for the CODEOWNERS file itself (e.g. \"@admin @team\")")
	root.AddCommand(newVersionCmd())

	return root
}

func Execute() error {
	return NewRootCmd().Execute()
}
