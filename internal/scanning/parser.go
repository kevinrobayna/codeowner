package scanning

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Mapping holds a file path and its code owners.
type Mapping struct {
	Path   string
	Owners []string
}

// DefaultPrefix is the default annotation prefix to search for.
const DefaultPrefix = "CodeOwner:"

// CodeOwnerFile is the name of the directory-level ownership file.
const CodeOwnerFile = ".codeowner"

// ParseProtect parses a whitespace-separated string of owner handles and
// returns a Mapping that protects the CODEOWNERS file itself. Each token must
// start with @ and contain only valid characters.
func ParseProtect(s string) (Mapping, error) {
	var owners []string
	for _, tok := range strings.Fields(s) {
		if !strings.HasPrefix(tok, "@") {
			return Mapping{}, fmt.Errorf("invalid owner %q: must start with @", tok)
		}
		if !isValidOwner(tok) {
			return Mapping{}, fmt.Errorf("invalid owner %q: contains invalid characters", tok)
		}
		owners = append(owners, tok)
	}
	return Mapping{Path: "CODEOWNERS", Owners: owners}, nil
}

// ParseFile reads a file and returns all code owners found in annotations
// matching the given prefix. Owners can appear on one line
// (CodeOwner: @a @b) or across multiple lines.
func ParseFile(path, prefix string) []string {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	seen := make(map[string]struct{})
	var owners []string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		for _, o := range extractOwners(scanner.Text(), prefix) {
			if _, dup := seen[o]; !dup {
				seen[o] = struct{}{}
				owners = append(owners, o)
			}
		}
	}

	return owners
}

// ParseCodeOwnerFile reads a .codeowner file and returns valid owner handles.
// Each line is split into whitespace-separated tokens; tokens must start with @
// and pass validation. Duplicate owners are removed.
func ParseCodeOwnerFile(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	seen := make(map[string]struct{})
	var owners []string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		for _, token := range strings.Fields(scanner.Text()) {
			if strings.HasPrefix(token, "@") && isValidOwner(token) {
				if _, dup := seen[token]; !dup {
					seen[token] = struct{}{}
					owners = append(owners, token)
				}
			}
		}
	}

	return owners
}

// ParseDir walks a directory and returns all CodeOwner mappings.
func ParseDir(root, prefix, dirOwnerFile string) ([]Mapping, error) {
	var mappings []Mapping

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && d.Name() == ".git" {
			return filepath.SkipDir
		}
		if d.IsDir() {
			return nil
		}

		if m, ok := parseEntry(root, path, d.Name(), prefix, dirOwnerFile); ok {
			mappings = append(mappings, m)
		}
		return nil
	})

	return mappings, err
}

// parseEntry handles a single file during directory walking, returning a
// Mapping and true if ownership was found.
func parseEntry(root, path, name, prefix, dirOwnerFile string) (Mapping, bool) {
	if name == dirOwnerFile {
		return parseDirOwnerEntry(root, path)
	}

	owners := ParseFile(path, prefix)
	if len(owners) == 0 {
		return Mapping{}, false
	}
	rel, relErr := filepath.Rel(root, path)
	if relErr != nil {
		rel = path
	}
	return Mapping{Path: "/" + rel, Owners: owners}, true
}

// parseDirOwnerEntry handles a .codeowner file, returning a directory-level
// Mapping with a root-anchored trailing-slash path.
func parseDirOwnerEntry(root, path string) (Mapping, bool) {
	owners := ParseCodeOwnerFile(path)
	if len(owners) == 0 {
		return Mapping{}, false
	}
	rel, relErr := filepath.Rel(root, filepath.Dir(path))
	if relErr != nil {
		rel = filepath.Dir(path)
	}
	if rel == "." {
		return Mapping{Path: "/", Owners: owners}, true
	}
	return Mapping{Path: "/" + rel + "/", Owners: owners}, true
}

// extractOwners parses all @-prefixed tokens after the prefix on a line.
func extractOwners(line, prefix string) []string {
	idx := strings.Index(line, prefix)
	if idx < 0 {
		return nil
	}

	// The prefix must be at the start of the line or preceded by a space.
	if idx > 0 && line[idx-1] != ' ' {
		return nil
	}

	rest := line[idx+len(prefix):]

	// Require a space between the prefix and the owners.
	if rest == "" || rest[0] != ' ' {
		return nil
	}

	var owners []string
	for _, token := range strings.Fields(rest) {
		if strings.HasPrefix(token, "@") && isValidOwner(token) {
			owners = append(owners, token)
		}
	}

	return owners
}

// isValidOwner checks that an owner handle contains only valid characters:
// @, letters, digits, hyphens, underscores, and slashes (for org/team).
func isValidOwner(s string) bool {
	for _, c := range s {
		switch {
		case c >= 'a' && c <= 'z',
			c >= 'A' && c <= 'Z',
			c >= '0' && c <= '9',
			c == '@', c == '-', c == '_', c == '/':
			continue
		default:
			return false
		}
	}
	return true
}
