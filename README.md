# codeowner

A CLI tool that scans source files for `CodeOwner:` annotations and generates [GitHub CODEOWNERS](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-code-owners) files.

## How it works

Add a `CodeOwner:` annotation inside any comment in your source files:

```python
# CodeOwner: @backend-team
```

```go
// CodeOwner: @backend-team
```

```html
<!-- CodeOwner: @frontend-team -->
```

Run `codeowner` and it prints a CODEOWNERS file to stdout:

```
/src/api/handler.go @backend-team
/web/index.html @frontend-team
```

The tool is **language-agnostic** — it searches for the annotation as plain text, so it works with any comment syntax.

### Multiple owners

Multiple owners on a single line:

```python
# CodeOwner: @backend-team @sre-team
```

Or across separate lines:

```python
# CodeOwner: @backend-team
# CodeOwner: @sre-team
```

Both produce:

```
/src/service.py @backend-team @sre-team
```

Duplicates are automatically deduplicated.

### Custom prefix

Use `--prefix` to search for a different annotation:

```sh
codeowner --prefix "Owner:" .
```

This matches `Owner: @my-team` instead of `CodeOwner: @my-team`.

### Annotation rules

- There **must** be a space between the prefix and the owner (`CodeOwner: @team`, not `CodeOwner:@team`)
- Owners **must** start with `@`

## Install

### Homebrew

```sh
brew install kevinrobayna/tap/codeowner
```

### Go

```sh
go install github.com/kevin-robayna/codeowner/cmd/codeowner@latest
```

### Binary releases

Download a binary from the [releases page](https://github.com/kevinrobayna/codeowner/releases).

## Usage

```sh
# Scan current directory
codeowner

# Scan a specific path
codeowner ./src

# Use a custom prefix
codeowner --prefix "Owner:" .

# Print version
codeowner version
```

## Development

```sh
# Build
make build

# Run tests
make test

# Lint
make lint

# Clean build artifacts
make clean
```

## License

[MIT](LICENSE)
