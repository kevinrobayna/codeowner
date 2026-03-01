package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRootCmd_Protect(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("// CodeOwner: @backend\npackage main\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	cmd := NewRootCmd()
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"--protect", "@admin @platform", dir})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()

	if !strings.HasPrefix(got, "CODEOWNERS @admin @platform\n") {
		t.Errorf("expected output to start with CODEOWNERS protect line, got:\n%s", got)
	}
	if !strings.Contains(got, "/main.go @backend") {
		t.Errorf("expected output to contain /main.go mapping, got:\n%s", got)
	}
}

func TestRootCmd_ProtectInvalidOwner(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("// CodeOwner: @backend\npackage main\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := NewRootCmd()
	cmd.SetArgs([]string{"--protect", "no-at-sign", dir})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for invalid --protect value")
	}
	if !strings.Contains(err.Error(), "--protect") {
		t.Errorf("error should mention --protect, got: %v", err)
	}
}

func TestRootCmd_NoProtect(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("// CodeOwner: @backend\npackage main\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	cmd := NewRootCmd()
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{dir})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()

	if strings.Contains(got, "CODEOWNERS") {
		t.Errorf("expected no CODEOWNERS protect line without --protect, got:\n%s", got)
	}
	if !strings.Contains(got, "/main.go @backend") {
		t.Errorf("expected output to contain /main.go mapping, got:\n%s", got)
	}
}
