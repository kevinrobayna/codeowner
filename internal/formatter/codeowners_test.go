package formatter_test

import (
	"path/filepath"
	"testing"

	"github.com/kevin-robayna/codeowner/internal/formatter"
	"github.com/kevin-robayna/codeowner/internal/owner"
)

func TestCodeOwners(t *testing.T) {
	mappings := []owner.Mapping{
		// Intentionally unordered to verify sorting.
		{Path: filepath.Join("src", "cmd", "main.go"), Owners: []string{"@backend"}},
		{Path: "README.md", Owners: []string{"@docs"}},
		{Path: filepath.Join(".github", "workflows", "ci.yml"), Owners: []string{"@devops"}},
		{Path: filepath.Join("src", "cmd", "helper.go"), Owners: []string{"@backend"}},
		{Path: filepath.Join("src", "lib", "utils.go"), Owners: []string{"@platform"}},
		{Path: "Makefile", Owners: []string{"@infra"}},
		{Path: filepath.Join(".github", "workflows", "deploy.yml"), Owners: []string{"@devops"}},
		{Path: filepath.Join("web", "index.html"), Owners: []string{"@frontend"}},
	}

	got := formatter.CodeOwners(mappings)
	want := "Makefile @infra\n" +
		"README.md @docs\n" +
		"\n" +
		filepath.Join(".github", "workflows", "ci.yml") + " @devops\n" +
		filepath.Join(".github", "workflows", "deploy.yml") + " @devops\n" +
		"\n" +
		filepath.Join("src", "cmd", "helper.go") + " @backend\n" +
		filepath.Join("src", "cmd", "main.go") + " @backend\n" +
		"\n" +
		filepath.Join("src", "lib", "utils.go") + " @platform\n" +
		"\n" +
		filepath.Join("web", "index.html") + " @frontend\n"

	if got != want {
		t.Errorf("CodeOwners:\ngot:\n%s\nwant:\n%s", got, want)
	}
}

func TestCodeOwners_RootFilesGroupedTogether(t *testing.T) {
	mappings := []owner.Mapping{
		{Path: "README.md", Owners: []string{"@docs"}},
		{Path: "LICENSE.txt", Owners: []string{"@legal"}},
		{Path: "docs.md", Owners: []string{"@docs"}},
	}

	got := formatter.CodeOwners(mappings)
	want := "LICENSE.txt @legal\n" +
		"README.md @docs\n" +
		"docs.md @docs\n"

	if got != want {
		t.Errorf("CodeOwners:\ngot:\n%s\nwant:\n%s", got, want)
	}
}

func TestCodeOwners_SameDirectoryNoBlankLines(t *testing.T) {
	mappings := []owner.Mapping{
		{Path: filepath.Join("testdata", "example.go"), Owners: []string{"@go_owner"}},
		{Path: filepath.Join("testdata", "example.py"), Owners: []string{"@python_owner"}},
		{Path: filepath.Join("testdata", "example.rs"), Owners: []string{"@rust_owner"}},
	}

	got := formatter.CodeOwners(mappings)
	want := filepath.Join("testdata", "example.go") + " @go_owner\n" +
		filepath.Join("testdata", "example.py") + " @python_owner\n" +
		filepath.Join("testdata", "example.rs") + " @rust_owner\n"

	if got != want {
		t.Errorf("CodeOwners:\ngot:\n%s\nwant:\n%s", got, want)
	}
}
