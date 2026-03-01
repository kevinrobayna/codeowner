package formatter_test

import (
	"testing"

	"github.com/kevin-robayna/codeowner/internal/formatter"
	"github.com/kevin-robayna/codeowner/internal/scanning"
)

func TestCodeOwners(t *testing.T) {
	t.Parallel()

	mappings := []scanning.Mapping{
		// Intentionally unordered to verify sorting.
		{Path: "/src/cmd/main.go", Owners: []string{"@backend"}},
		{Path: "/README.md", Owners: []string{"@docs"}},
		{Path: "/.github/workflows/ci.yml", Owners: []string{"@devops"}},
		{Path: "/src/cmd/helper.go", Owners: []string{"@backend"}},
		{Path: "/src/lib/utils.go", Owners: []string{"@platform"}},
		{Path: "/Makefile", Owners: []string{"@infra"}},
		{Path: "/.github/workflows/deploy.yml", Owners: []string{"@devops"}},
		{Path: "/web/index.html", Owners: []string{"@frontend"}},
	}

	got := formatter.CodeOwners(mappings)
	want := "/Makefile @infra\n" +
		"/README.md @docs\n" +
		"\n" +
		"/.github/workflows/ci.yml @devops\n" +
		"/.github/workflows/deploy.yml @devops\n" +
		"\n" +
		"/src/cmd/helper.go @backend\n" +
		"/src/cmd/main.go @backend\n" +
		"\n" +
		"/src/lib/utils.go @platform\n" +
		"\n" +
		"/web/index.html @frontend\n"

	if got != want {
		t.Errorf("CodeOwners:\ngot:\n%s\nwant:\n%s", got, want)
	}
}

func TestCodeOwners_RootFilesGroupedTogether(t *testing.T) {
	t.Parallel()

	mappings := []scanning.Mapping{
		{Path: "/README.md", Owners: []string{"@docs"}},
		{Path: "/LICENSE.txt", Owners: []string{"@legal"}},
		{Path: "/docs.md", Owners: []string{"@docs"}},
	}

	got := formatter.CodeOwners(mappings)
	want := "/LICENSE.txt @legal\n" +
		"/README.md @docs\n" +
		"/docs.md @docs\n"

	if got != want {
		t.Errorf("CodeOwners:\ngot:\n%s\nwant:\n%s", got, want)
	}
}

func TestCodeOwners_DirOwnerPathsSortAndGroup(t *testing.T) {
	t.Parallel()

	mappings := []scanning.Mapping{
		{Path: "/src/cmd/main.go", Owners: []string{"@backend"}},
		{Path: "/src/cmd/", Owners: []string{"@cmd-team"}},
		{Path: "/lib/", Owners: []string{"@lib-team"}},
		{Path: "/lib/utils.go", Owners: []string{"@utils"}},
	}

	got := formatter.CodeOwners(mappings)
	want := "/lib/ @lib-team\n" +
		"/lib/utils.go @utils\n" +
		"\n" +
		"/src/cmd/ @cmd-team\n" +
		"/src/cmd/main.go @backend\n"

	if got != want {
		t.Errorf("CodeOwners:\ngot:\n%s\nwant:\n%s", got, want)
	}
}

func TestCodeOwners_ProtectMapping(t *testing.T) {
	t.Parallel()

	mappings := []scanning.Mapping{
		{Path: "/README.md", Owners: []string{"@docs"}},
		{Path: "/Makefile", Owners: []string{"@infra"}},
		{Path: "CODEOWNERS", Owners: []string{"@kevinrobayna", "@admin"}},
	}

	got := formatter.CodeOwners(mappings)
	want := "CODEOWNERS @kevinrobayna @admin\n" +
		"\n" +
		"/Makefile @infra\n" +
		"/README.md @docs\n"

	if got != want {
		t.Errorf("CodeOwners with protect:\ngot:\n%s\nwant:\n%s", got, want)
	}
}

func TestCodeOwners_RootDirOwnerGroupsWithRootFiles(t *testing.T) {
	t.Parallel()

	mappings := []scanning.Mapping{
		{Path: "/README.md", Owners: []string{"@docs"}},
		{Path: "/", Owners: []string{"@root-team"}},
		{Path: "/src/main.go", Owners: []string{"@backend"}},
		{Path: "/lib/", Owners: []string{"@lib-team"}},
	}

	got := formatter.CodeOwners(mappings)
	// Root dir owner "/" should be grouped with root-level entries.
	want := "/ @root-team\n" +
		"/README.md @docs\n" +
		"\n" +
		"/lib/ @lib-team\n" +
		"\n" +
		"/src/main.go @backend\n"

	if got != want {
		t.Errorf("CodeOwners:\ngot:\n%s\nwant:\n%s", got, want)
	}
}

func TestCodeOwners_SameDirectoryNoBlankLines(t *testing.T) {
	t.Parallel()

	mappings := []scanning.Mapping{
		{Path: "/testdata/example.go", Owners: []string{"@go_owner"}},
		{Path: "/testdata/example.py", Owners: []string{"@python_owner"}},
		{Path: "/testdata/example.rs", Owners: []string{"@rust_owner"}},
	}

	got := formatter.CodeOwners(mappings)
	want := "/testdata/example.go @go_owner\n" +
		"/testdata/example.py @python_owner\n" +
		"/testdata/example.rs @rust_owner\n"

	if got != want {
		t.Errorf("CodeOwners:\ngot:\n%s\nwant:\n%s", got, want)
	}
}
