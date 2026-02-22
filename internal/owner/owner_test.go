package owner_test

import (
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"testing"

	"github.com/kevin-robayna/codeowner/internal/owner"
)

func testdataDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata")
}

func TestParseFile(t *testing.T) {
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
			path := filepath.Join(dir, tt.file)
			got := owner.ParseFile(path, owner.DefaultPrefix)
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
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.py")
	if err := os.WriteFile(path, []byte("# just a regular comment\nx = 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	got := owner.ParseFile(path, owner.DefaultPrefix)
	if len(got) > 0 {
		t.Errorf("expected no CodeOwner, got %v", got)
	}
}

func TestParseFile_MultipleOwnersOneLine(t *testing.T) {
	path := filepath.Join(testdataDir(), "multi_owners_single_line.py")
	got := owner.ParseFile(path, owner.DefaultPrefix)
	want := []string{"@team-a", "@team-b", "@person-c"}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseFile_MultipleOwnersMultipleLines(t *testing.T) {
	path := filepath.Join(testdataDir(), "multi_owners_multi_line.py")
	got := owner.ParseFile(path, owner.DefaultPrefix)
	want := []string{"@team-frontend", "@team-backend"}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseFile_DeduplicatesOwners(t *testing.T) {
	path := filepath.Join(testdataDir(), "multi_owners_deduplicated.py")
	got := owner.ParseFile(path, owner.DefaultPrefix)
	want := []string{"@team-a"}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseDir(t *testing.T) {
	dir := testdataDir()

	mappings, err := owner.ParseDir(dir, owner.DefaultPrefix)
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
		"example.py":                          {"@python_owner"},
		"example.go":                          {"@go_owner"},
		"example.rs":                          {"@rust_owner"},
		"example.css":                         {"@css_owner"},
		filepath.Join("nested", "handler.go"): {"@api-team"},
		filepath.Join("nested", "deeply", "service.py"):            {"@platform-team"},
		filepath.Join("nested", "deeply", "nested", "config.yaml"): {"@infra-team"},
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
	dir := t.TempDir()
	path := filepath.Join(dir, "main.go")
	content := "// Owner: @team-backend\npackage main\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	got := owner.ParseFile(path, "Owner:")
	want := []string{"@team-backend"}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}

	got = owner.ParseFile(path, owner.DefaultPrefix)
	if len(got) > 0 {
		t.Errorf("default prefix should not match custom annotation, got %v", got)
	}
}

func TestParseFile_RejectsNoSpace(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nospace.py")
	content := "# CodeOwner:@team\nx = 1\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	got := owner.ParseFile(path, owner.DefaultPrefix)
	if len(got) > 0 {
		t.Errorf("should reject annotation without space after prefix, got %v", got)
	}
}

func TestParseFile_RejectsNoAt(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "noat.py")
	content := "# CodeOwner: team-backend\nx = 1\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	got := owner.ParseFile(path, owner.DefaultPrefix)
	if len(got) > 0 {
		t.Errorf("should reject owner without @ prefix, got %v", got)
	}
}

func TestFormatCodeOwners(t *testing.T) {
	mappings := []owner.Mapping{
		{Path: "src/main.go", Owners: []string{"@backend"}},
		{Path: "web/index.html", Owners: []string{"@frontend", "@design"}},
	}

	got := owner.FormatCodeOwners(mappings)
	want := "/src/main.go @backend\n/web/index.html @frontend @design\n"

	if got != want {
		t.Errorf("FormatCodeOwners:\ngot:  %q\nwant: %q", got, want)
	}
}
