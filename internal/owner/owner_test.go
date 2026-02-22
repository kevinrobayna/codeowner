package owner_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/kevin-robayna/codeowner/internal/owner"
)

func testdataDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata")
}

func TestParseFile(t *testing.T) {
	tests := []struct {
		file  string
		owner string
	}{
		{"example.py", "@python_owner"},
		{"example.rb", "@ruby_owner"},
		{"example.sh", "@shell_owner"},
		{"example.pl", "@perl_owner"},
		{"Dockerfile", "@docker_owner"},
		{"example.r", "@r_owner"},
		{"example.ex", "@elixir_owner"},
		{"example.yaml", "@yaml_owner"},
		{"example.sql", "@sql_owner"},
		{"example.lua", "@lua_owner"},
		{"example.hs", "@haskell_owner"},
		{"example.c", "@c_owner"},
		{"example.cpp", "@cpp_owner"},
		{"Example.java", "@java_owner"},
		{"example.js", "@js_owner"},
		{"example.ts", "@ts_owner"},
		{"example.go", "@go_owner"},
		{"example.rs", "@rust_owner"},
		{"example.swift", "@swift_owner"},
		{"Example.kt", "@kotlin_owner"},
		{"Example.cs", "@csharp_owner"},
		{"example.scala", "@scala_owner"},
		{"example.php", "@php_owner"},
		{"example.clj", "@clojure_owner"},
		{"example.el", "@elisp_owner"},
		{"example.html", "@html_owner"},
		{"example.xml", "@xml_owner"},
		{"example.tex", "@latex_owner"},
		{"example.erl", "@erlang_owner"},
		{"example.css", "@css_owner"},
	}

	dir := testdataDir()

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			path := filepath.Join(dir, tt.file)
			got, ok := owner.ParseFile(path, owner.DefaultPrefix)
			if !ok {
				t.Fatalf("expected to find CodeOwner in %s, got nothing", tt.file)
			}
			if got != tt.owner {
				t.Errorf("ParseFile(%s) = %q, want %q", tt.file, got, tt.owner)
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

	got, ok := owner.ParseFile(path, owner.DefaultPrefix)
	if ok {
		t.Errorf("expected no CodeOwner, got %q", got)
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

	found := make(map[string]string)
	for _, m := range mappings {
		found[m.Path] = m.Owner
	}

	checks := map[string]string{
		"example.py":  "@python_owner",
		"example.go":  "@go_owner",
		"example.rs":  "@rust_owner",
		"example.css": "@css_owner",
	}

	for path, wantOwner := range checks {
		if got, ok := found[path]; !ok {
			t.Errorf("missing mapping for %s", path)
		} else if got != wantOwner {
			t.Errorf("mapping for %s = %q, want %q", path, got, wantOwner)
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

	got, ok := owner.ParseFile(path, "Owner:")
	if !ok {
		t.Fatal("expected to find owner with custom prefix")
	}
	if got != "@team-backend" {
		t.Errorf("got %q, want %q", got, "@team-backend")
	}

	_, ok = owner.ParseFile(path, owner.DefaultPrefix)
	if ok {
		t.Error("default prefix should not match custom annotation")
	}
}

func TestParseFile_RejectsNoSpace(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nospace.py")
	content := "# CodeOwner:@team\nx = 1\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	_, ok := owner.ParseFile(path, owner.DefaultPrefix)
	if ok {
		t.Error("should reject annotation without space after prefix")
	}
}

func TestParseFile_RejectsNoAt(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "noat.py")
	content := "# CodeOwner: team-backend\nx = 1\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	_, ok := owner.ParseFile(path, owner.DefaultPrefix)
	if ok {
		t.Error("should reject owner without @ prefix")
	}
}

func TestFormatCodeOwners(t *testing.T) {
	mappings := []owner.Mapping{
		{Path: "src/main.go", Owner: "@backend"},
		{Path: "web/index.html", Owner: "@frontend"},
	}

	got := owner.FormatCodeOwners(mappings)
	want := "/src/main.go @backend\n/web/index.html @frontend\n"

	if got != want {
		t.Errorf("FormatCodeOwners:\ngot:  %q\nwant: %q", got, want)
	}
}
