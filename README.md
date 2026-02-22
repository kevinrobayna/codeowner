# codeowner

A CLI tool that parses repositories and generates CODEOWNERS files.

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
# Generate a CODEOWNERS file
codeowner generate

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
