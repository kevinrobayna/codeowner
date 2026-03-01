# Codeowner

A language-agnostic CLI tool that scans source files for `CodeOwner:` annotations and generates GitHub CODEOWNERS files.

## Architecture

- CLI built with `github.com/spf13/cobra`
- Minimal dependencies by design — only Cobra for CLI; all core logic is in-house

## Key Directories

- `cmd/codeowner/` — CLI entry point (`main.go`)
- `internal/cmd/` — Cobra command definitions (root command + version subcommand)
- `internal/scanning/` — File parsing and directory walking (`ParseFile`, `ParseDir`, `ParseCodeOwnerFile`, `ParseProtect`)
- `internal/formatter/` — CODEOWNERS output formatting (`CodeOwners`)
- `internal/appinfo/` — Build-time version info injected via ldflags
- `testdata/` — 47 test fixtures covering 30+ languages, nested dirs, edge cases

## Tech Stack

- Language: Go 1.25
- Build: Make + GoReleaser
- Linting: golangci-lint (28 linters, see `.golangci.yml`)
- CI: GitHub Actions (test matrix, lint, CodeQL, release with cosign signing)

## Development

### Building

```bash
make build
```

### Running Tests

```bash
make test
```

Runs `go test -race -cover ./...` across all packages.

### Test Parallelization

Always use `t.Parallel()` in tests. Every top-level test function and subtest should call `t.Parallel()` unless it modifies process-global state (e.g., `os.Chdir()`, `os.Setenv()`).

```go
func TestFeature(t *testing.T) {
    t.Parallel()
    // ...
}

func TestFeature_Subtests(t *testing.T) {
    t.Parallel()
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            t.Parallel()
            // ...
        })
    }
}
```

### Linting and Formatting

```bash
make lint        # Auto-fix lint issues
make lint-check  # Check only (used in CI)
```

### Before Every Commit

```bash
make lint && make test
```

### Cleaning

```bash
make clean
```

## CLI Usage

```bash
codeowner [path]
```

**Flags:**

- `--prefix string` — Annotation prefix (default: `CodeOwner:`)
- `--dirowner string` — Directory ownership filename (default: `.codeowner`)
- `--protect string` — Owners for the CODEOWNERS file itself (e.g., `@admin @platform`)

**Subcommands:**

- `version` — Print version, commit, and build date

## Code Patterns

### Separation of Concerns

- **Scanning** reads files and extracts annotations — no formatting logic
- **Formatting** takes parsed mappings and produces output — no file I/O
- **Commands** wire scanning and formatting together via CLI flags

### Owner Validation

- Owners must start with `@`
- Valid characters: alphanumeric, `-`, `_`, `/` (for `@org/team` format)
- `CodeOwner:` prefix must be at line start or preceded by a space
- A space is required after the prefix

### Output Formatting

The formatter sorts entries into sections:

1. CODEOWNERS protect rule (if `--protect` used), separated by blank line
2. Root-level files (no subdirectory)
3. Hidden directories (`.github/`, etc.)
4. Regular directories, grouped by top-level path

Groups are separated by blank lines. Within groups, entries are sorted alphabetically.

## Go Code Style

- Write lint-compliant Go code. Reference `.golangci.yml` for enabled linters.
- Follow standard Go idioms: proper error handling, no unused variables/imports, `gofmt` formatting.
- Handle all errors explicitly.
