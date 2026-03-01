package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExecute(t *testing.T) {
	t.Parallel()

	// Execute() calls NewRootCmd().Execute() with no args, scanning ".".
	// We just verify it doesn't panic or return an unexpected error.
	_ = Execute()
}

func TestRootCmd_VersionSubcommand(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	cmd := NewRootCmd()
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"version"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "codeowner") {
		t.Errorf("version output should contain 'codeowner', got: %s", got)
	}
}

func TestRootCmd_NoAnnotationsFound(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "empty.txt"), []byte("no annotations here\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	cmd := NewRootCmd()
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{dir})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stdout.Len() > 0 {
		t.Errorf("expected no stdout output, got: %s", stdout.String())
	}
	if !strings.Contains(stderr.String(), "no CodeOwner annotations found") {
		t.Errorf("expected stderr message about no annotations, got: %s", stderr.String())
	}
}

func TestRootCmd_InvalidDirectory(t *testing.T) {
	t.Parallel()

	cmd := NewRootCmd()
	cmd.SetArgs([]string{"/nonexistent/path/that/does/not/exist"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for non-existent directory")
	}
}

func TestRootCmd_TooManyArgs(t *testing.T) {
	t.Parallel()

	cmd := NewRootCmd()
	cmd.SetArgs([]string{"arg1", "arg2"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for too many arguments")
	}
}

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
