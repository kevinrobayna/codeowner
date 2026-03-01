package scanning_test

import (
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"testing"

	"github.com/kevin-robayna/codeowner/internal/scanning"
)

func testdataDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata")
}

func TestParseFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		file   string
		owners []string
	}{
		{"example.py", []string{"@python_owner"}},
		{"example.rb", []string{"@ruby_owner"}},
		{"example.sh", []string{"@shell_owner"}},
		{"example.pl", []string{"@perl_owner"}},
		{"Dockerfile", []string{"@docker_owner"}},
		{"example.r", []string{"@r_owner"}},
		{"example.ex", []string{"@elixir_owner"}},
		{"example.yaml", []string{"@yaml_owner"}},
		{"example.sql", []string{"@sql_owner"}},
		{"example.lua", []string{"@lua_owner"}},
		{"example.hs", []string{"@haskell_owner"}},
		{"example.c", []string{"@c_owner"}},
		{"example.cpp", []string{"@cpp_owner"}},
		{"Example.java", []string{"@java_owner"}},
		{"example.js", []string{"@js_owner"}},
		{"example.ts", []string{"@ts_owner"}},
		{"example.go", []string{"@go_owner"}},
		{"example.rs", []string{"@rust_owner"}},
		{"example.swift", []string{"@swift_owner"}},
		{"Example.kt", []string{"@kotlin_owner"}},
		{"Example.cs", []string{"@csharp_owner"}},
		{"example.scala", []string{"@scala_owner"}},
		{"example.php", []string{"@php_owner"}},
		{"example.clj", []string{"@clojure_owner"}},
		{"example.el", []string{"@elisp_owner"}},
		{"example.html", []string{"@html_owner"}},
		{"example.xml", []string{"@xml_owner"}},
		{"example.tex", []string{"@latex_owner"}},
		{"example.erl", []string{"@erlang_owner"}},
		{"example.css", []string{"@css_owner"}},
	}

	dir := testdataDir()

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			t.Parallel()

			path := filepath.Join(dir, tt.file)
			got := scanning.ParseFile(path, scanning.DefaultPrefix)
			if len(got) == 0 {
				t.Fatalf("expected to find CodeOwner in %s, got nothing", tt.file)
			}
			if !slices.Equal(got, tt.owners) {
				t.Errorf("ParseFile(%s) = %v, want %v", tt.file, got, tt.owners)
			}
		})
	}
}

func TestParseFile_NoAnnotation(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "empty.py")
	if err := os.WriteFile(path, []byte("# just a regular comment\nx = 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	got := scanning.ParseFile(path, scanning.DefaultPrefix)
	if len(got) > 0 {
		t.Errorf("expected no CodeOwner, got %v", got)
	}
}

func TestParseFile_MultipleOwnersOneLine(t *testing.T) {
	t.Parallel()

	path := filepath.Join(testdataDir(), "multi_owners_single_line.py")
	got := scanning.ParseFile(path, scanning.DefaultPrefix)
	want := []string{"@team-a", "@team-b", "@person-c"}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseFile_MultipleOwnersMultipleLines(t *testing.T) {
	t.Parallel()

	path := filepath.Join(testdataDir(), "multi_owners_multi_line.py")
	got := scanning.ParseFile(path, scanning.DefaultPrefix)
	want := []string{"@team-frontend", "@team-backend"}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseFile_DeduplicatesOwners(t *testing.T) {
	t.Parallel()

	path := filepath.Join(testdataDir(), "multi_owners_deduplicated.py")
	got := scanning.ParseFile(path, scanning.DefaultPrefix)
	want := []string{"@team-a"}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseDir(t *testing.T) {
	t.Parallel()

	dir := testdataDir()

	mappings, err := scanning.ParseDir(dir, scanning.DefaultPrefix, scanning.CodeOwnerFile)
	if err != nil {
		t.Fatalf("ParseDir(%s) error: %v", dir, err)
	}

	if len(mappings) < 30 {
		t.Errorf("expected at least 30 mappings, got %d", len(mappings))
	}

	found := make(map[string][]string)
	for _, m := range mappings {
		found[m.Path] = m.Owners
	}

	checks := map[string][]string{
		"/example.py":  {"@python_owner"},
		"/example.go":  {"@go_owner"},
		"/example.rs":  {"@rust_owner"},
		"/example.css": {"@css_owner"},
		"/" + filepath.Join("nested", "handler.go"):                      {"@api-team"},
		"/" + filepath.Join("nested", "deeply", "service.py"):            {"@platform-team"},
		"/" + filepath.Join("nested", "deeply", "nested", "config.yaml"): {"@infra-team"},
		"/" + filepath.Join(".github", "workflows", "ci.yml"):            {"@devops-team"},
		"/" + filepath.Join("nested", "deeply", ".hidden", "secret.rb"):  {"@secret-team"},
	}

	for path, wantOwners := range checks {
		if got, ok := found[path]; !ok {
			t.Errorf("missing mapping for %s", path)
		} else if !slices.Equal(got, wantOwners) {
			t.Errorf("mapping for %s = %v, want %v", path, got, wantOwners)
		}
	}
}

func TestParseFile_CustomPrefix(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "main.go")
	content := "// Owner: @team-backend\npackage main\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	got := scanning.ParseFile(path, "Owner:")
	want := []string{"@team-backend"}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}

	got = scanning.ParseFile(path, scanning.DefaultPrefix)
	if len(got) > 0 {
		t.Errorf("default prefix should not match custom annotation, got %v", got)
	}
}

func TestParseFile_RejectsNoSpace(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "nospace.py")
	content := "# CodeOwner:@team\nx = 1\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	got := scanning.ParseFile(path, scanning.DefaultPrefix)
	if len(got) > 0 {
		t.Errorf("should reject annotation without space after prefix, got %v", got)
	}
}

func TestParseFile_OrgTeamOwner(t *testing.T) {
	t.Parallel()

	path := filepath.Join(testdataDir(), "org_team_owner.py")
	got := scanning.ParseFile(path, scanning.DefaultPrefix)
	want := []string{"@myorg/backend-team"}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseFile_OrgTeamMultipleOwners(t *testing.T) {
	t.Parallel()

	path := filepath.Join(testdataDir(), "org_team_multiple_owners.js")
	got := scanning.ParseFile(path, scanning.DefaultPrefix)
	want := []string{"@myorg/frontend-team", "@myorg/design-team", "@individual-dev"}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseFile_RejectsPrefixNotPrecededBySpace(t *testing.T) {
	t.Parallel()

	path := filepath.Join(testdataDir(), "invalid_prefix_not_preceded_by_space.md")
	got := scanning.ParseFile(path, scanning.DefaultPrefix)
	if len(got) > 0 {
		t.Errorf("should reject prefix not preceded by space, got %v", got)
	}
}

func TestParseFile_TabBeforePrefix(t *testing.T) {
	t.Parallel()

	path := filepath.Join(testdataDir(), "tab_before_prefix.py")
	got := scanning.ParseFile(path, scanning.DefaultPrefix)
	want := []string{"@tab-team"}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseFile_RejectsOwnerWithSpecialChars(t *testing.T) {
	t.Parallel()

	path := filepath.Join(testdataDir(), "invalid_owner_special_chars.txt")
	got := scanning.ParseFile(path, scanning.DefaultPrefix)
	if len(got) > 0 {
		t.Errorf("should reject owners with special characters, got %v", got)
	}
}

func TestParseCodeOwnerFile_SingleOwner(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, ".codeowner")
	if err := os.WriteFile(path, []byte("@team-a\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	got := scanning.ParseCodeOwnerFile(path)
	want := []string{"@team-a"}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseCodeOwnerFile_MultipleOwnersSpaceSeparated(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, ".codeowner")
	if err := os.WriteFile(path, []byte("@team-a @team-b\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	got := scanning.ParseCodeOwnerFile(path)
	want := []string{"@team-a", "@team-b"}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseCodeOwnerFile_MultipleOwnersMultiLine(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, ".codeowner")
	if err := os.WriteFile(path, []byte("@team-a\n@team-b\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	got := scanning.ParseCodeOwnerFile(path)
	want := []string{"@team-a", "@team-b"}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseCodeOwnerFile_Deduplication(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, ".codeowner")
	if err := os.WriteFile(path, []byte("@team-a @team-a\n@team-a\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	got := scanning.ParseCodeOwnerFile(path)
	want := []string{"@team-a"}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseCodeOwnerFile_InvalidOwnersIgnored(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, ".codeowner")
	if err := os.WriteFile(path, []byte("@valid-team not-an-owner @bad!owner\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	got := scanning.ParseCodeOwnerFile(path)
	want := []string{"@valid-team"}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseDir_CodeOwnerFile(t *testing.T) {
	t.Parallel()

	dir := testdataDir()

	mappings, err := scanning.ParseDir(dir, scanning.DefaultPrefix, scanning.CodeOwnerFile)
	if err != nil {
		t.Fatalf("ParseDir(%s) error: %v", dir, err)
	}

	found := make(map[string][]string)
	for _, m := range mappings {
		found[m.Path] = m.Owners
	}

	// .codeowner file should produce a /dir/ mapping
	dirOwnerPath := "/dirowner/"
	if got, ok := found[dirOwnerPath]; !ok {
		t.Errorf("missing mapping for %s", dirOwnerPath)
	} else if !slices.Equal(got, []string{"@dir-owner"}) {
		t.Errorf("mapping for %s = %v, want %v", dirOwnerPath, got, []string{"@dir-owner"})
	}

	// Regular file in same directory should still be parsed normally
	appPath := "/" + filepath.Join("dirowner", "app.go")
	if got, ok := found[appPath]; !ok {
		t.Errorf("missing mapping for %s", appPath)
	} else if !slices.Equal(got, []string{"@app-team"}) {
		t.Errorf("mapping for %s = %v, want %v", appPath, got, []string{"@app-team"})
	}

	// Nested .codeowner should produce its own /dir/sub/ mapping
	subPath := "/" + filepath.Join("dirowner", "sub") + "/"
	if got, ok := found[subPath]; !ok {
		t.Errorf("missing mapping for %s", subPath)
	} else if !slices.Equal(got, []string{"@sub-team"}) {
		t.Errorf("mapping for %s = %v, want %v", subPath, got, []string{"@sub-team"})
	}

	// Annotated file alongside nested .codeowner
	handlerPath := "/" + filepath.Join("dirowner", "sub", "handler.go")
	if got, ok := found[handlerPath]; !ok {
		t.Errorf("missing mapping for %s", handlerPath)
	} else if !slices.Equal(got, []string{"@handler-team"}) {
		t.Errorf("mapping for %s = %v, want %v", handlerPath, got, []string{"@handler-team"})
	}

	// Deeper nested file without its own .codeowner
	utilPath := "/" + filepath.Join("dirowner", "sub", "deep", "util.go")
	if got, ok := found[utilPath]; !ok {
		t.Errorf("missing mapping for %s", utilPath)
	} else if !slices.Equal(got, []string{"@util-team"}) {
		t.Errorf("mapping for %s = %v, want %v", utilPath, got, []string{"@util-team"})
	}
}

func TestParseProtect(t *testing.T) {
	t.Parallel()

	m, err := scanning.ParseProtect("@admin @platform-team")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Path != "CODEOWNERS" {
		t.Errorf("path = %q, want %q", m.Path, "CODEOWNERS")
	}
	want := []string{"@admin", "@platform-team"}
	if !slices.Equal(m.Owners, want) {
		t.Errorf("owners = %v, want %v", m.Owners, want)
	}
}

func TestParseProtect_InvalidNoAt(t *testing.T) {
	t.Parallel()

	_, err := scanning.ParseProtect("admin")
	if err == nil {
		t.Fatal("expected error for owner without @")
	}
}

func TestParseProtect_InvalidChars(t *testing.T) {
	t.Parallel()

	_, err := scanning.ParseProtect("@bad!owner")
	if err == nil {
		t.Fatal("expected error for owner with invalid characters")
	}
}

func TestParseFile_RejectsBareAt(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "bareat.py")
	content := "# CodeOwner: @\nx = 1\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	got := scanning.ParseFile(path, scanning.DefaultPrefix)
	if len(got) > 0 {
		t.Errorf("should reject bare @ as owner, got %v", got)
	}
}

func TestParseFile_RejectsNoAt(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "noat.py")
	content := "# CodeOwner: team-backend\nx = 1\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	got := scanning.ParseFile(path, scanning.DefaultPrefix)
	if len(got) > 0 {
		t.Errorf("should reject owner without @ prefix, got %v", got)
	}
}
