package owner

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Mapping holds a file path and its code owner.
type Mapping struct {
	Path  string
	Owner string
}

// DefaultPrefix is the default annotation prefix to search for.
const DefaultPrefix = "CodeOwner:"

// ParseFile reads a file and returns the code owner if an annotation matching
// the given prefix is found on any line, regardless of comment syntax.
func ParseFile(path, prefix string) (string, bool) {
	f, err := os.Open(path)
	if err != nil {
		return "", false
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if o, ok := extractOwner(scanner.Text(), prefix); ok {
			return o, true
		}
	}

	return "", false
}

// ParseDir walks a directory and returns all CodeOwner mappings.
func ParseDir(root, prefix string) ([]Mapping, error) {
	var mappings []Mapping

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && path != root && strings.HasPrefix(d.Name(), ".") {
			return filepath.SkipDir
		}
		if d.IsDir() {
			return nil
		}

		if o, ok := ParseFile(path, prefix); ok {
			rel, relErr := filepath.Rel(root, path)
			if relErr != nil {
				rel = path
			}
			mappings = append(mappings, Mapping{Path: rel, Owner: o})
		}
		return nil
	})

	return mappings, err
}

// FormatCodeOwners formats mappings as a GitHub CODEOWNERS file.
func FormatCodeOwners(mappings []Mapping) string {
	var b strings.Builder
	for _, m := range mappings {
		fmt.Fprintf(&b, "/%s %s\n", m.Path, m.Owner)
	}
	return b.String()
}

func extractOwner(line, prefix string) (string, bool) {
	idx := strings.Index(line, prefix)
	if idx < 0 {
		return "", false
	}

	rest := line[idx+len(prefix):]

	// Require a space between the prefix and the owner.
	if rest == "" || rest[0] != ' ' {
		return "", false
	}

	value := strings.TrimSpace(rest)

	// Owner must start with @.
	if !strings.HasPrefix(value, "@") {
		return "", false
	}

	if spaceIdx := strings.IndexByte(value, ' '); spaceIdx > 0 {
		value = value[:spaceIdx]
	}

	return value, true
}
