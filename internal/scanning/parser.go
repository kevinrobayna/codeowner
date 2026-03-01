package scanning

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// maxFileSize is the maximum file size (1 MB) that ParseDir will scan.
// Files larger than this are skipped to avoid wasting time on binaries or
// generated artifacts.
const maxFileSize = 1 << 20

// binarySniffSize is the number of bytes read to detect binary content.
const binarySniffSize = 512

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
	fields := strings.Fields(s)
	if len(fields) == 0 {
		return Mapping{}, fmt.Errorf("empty protect string: at least one owner is required")
	}
	owners := make([]string, 0, len(fields))
	for _, tok := range fields {
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
func ParseFile(path, prefix string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return scanOwners(f, path, prefix)
}

// scanOwners extracts code owners from a reader, used by both ParseFile and
// parseEntry (which passes an already-open file to avoid a double open).
func scanOwners(r io.Reader, path, prefix string) ([]string, error) {
	seen := make(map[string]struct{})
	var owners []string

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		for _, o := range extractOwners(scanner.Text(), prefix) {
			owners = appendUnique(seen, owners, o)
		}
	}
	if err := scanner.Err(); err != nil {
		return owners, fmt.Errorf("reading %s: %w", path, err)
	}

	return owners, nil
}

// ParseCodeOwnerFile reads a .codeowner file and returns valid owner handles.
// Each line is split into whitespace-separated tokens; tokens must start with @
// and pass validation. Duplicate owners are removed.
func ParseCodeOwnerFile(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	seen := make(map[string]struct{})
	var owners []string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		for _, token := range strings.Fields(scanner.Text()) {
			if strings.HasPrefix(token, "@") && isValidOwner(token) {
				owners = appendUnique(seen, owners, token)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return owners, fmt.Errorf("reading %s: %w", path, err)
	}

	return owners, nil
}

// ParseDir walks a directory and returns all CodeOwner mappings.
func ParseDir(root, prefix, dirOwnerFile string) ([]Mapping, error) {
	var mappings []Mapping

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Type()&fs.ModeSymlink != 0 {
			return nil
		}
		if d.IsDir() && d.Name() == ".git" {
			return filepath.SkipDir
		}
		if d.IsDir() {
			return nil
		}

		m, ok, parseErr := parseEntry(root, path, d, prefix, dirOwnerFile)
		if parseErr != nil {
			return parseErr
		}
		if ok {
			mappings = append(mappings, m)
		}
		return nil
	})

	return mappings, err
}

// isBinary reports whether buf contains a null byte, indicating binary content.
func isBinary(buf []byte) bool {
	return bytes.IndexByte(buf, 0) >= 0
}

// parseEntry handles a single file during directory walking, returning a
// Mapping and true if ownership was found. It uses the DirEntry from WalkDir
// to avoid a redundant os.Stat call, and opens the file once for both binary
// detection and annotation scanning.
func parseEntry(root, path string, d fs.DirEntry, prefix, dirOwnerFile string) (Mapping, bool, error) {
	if d.Name() == dirOwnerFile {
		return parseDirOwnerEntry(root, path)
	}

	info, err := d.Info()
	if err != nil {
		return Mapping{}, false, err
	}
	if info.Size() > maxFileSize {
		return Mapping{}, false, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return Mapping{}, false, err
	}
	defer f.Close()

	// Sniff the first bytes to detect binary content.
	buf := make([]byte, binarySniffSize)
	n, err := f.Read(buf)
	if err != nil && !errors.Is(err, io.EOF) {
		return Mapping{}, false, err
	}
	if isBinary(buf[:n]) {
		return Mapping{}, false, nil
	}

	// Seek back to the start so scanOwners reads the full file.
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return Mapping{}, false, err
	}

	owners, err := scanOwners(f, path, prefix)
	if err != nil {
		return Mapping{}, false, err
	}
	if len(owners) == 0 {
		return Mapping{}, false, nil
	}
	rel, relErr := filepath.Rel(root, path)
	if relErr != nil {
		rel = path
	}
	return Mapping{Path: "/" + filepath.ToSlash(rel), Owners: owners}, true, nil
}

// parseDirOwnerEntry handles a .codeowner file, returning a directory-level
// Mapping with a root-anchored trailing-slash path.
func parseDirOwnerEntry(root, path string) (Mapping, bool, error) {
	owners, err := ParseCodeOwnerFile(path)
	if err != nil {
		return Mapping{}, false, err
	}
	if len(owners) == 0 {
		return Mapping{}, false, nil
	}
	rel, relErr := filepath.Rel(root, filepath.Dir(path))
	if relErr != nil {
		rel = filepath.Dir(path)
	}
	if rel == "." {
		return Mapping{Path: "/", Owners: owners}, true, nil
	}
	return Mapping{Path: "/" + filepath.ToSlash(rel) + "/", Owners: owners}, true, nil
}

// extractOwners parses all @-prefixed tokens after the prefix on a line.
func extractOwners(line, prefix string) []string {
	idx := strings.Index(line, prefix)
	if idx < 0 {
		return nil
	}

	// The prefix must be at the start of the line or preceded by whitespace.
	if idx > 0 && line[idx-1] != ' ' && line[idx-1] != '\t' {
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

// appendUnique appends token to owners if it has not been seen before.
func appendUnique(seen map[string]struct{}, owners []string, token string) []string {
	if _, dup := seen[token]; !dup {
		seen[token] = struct{}{}
		owners = append(owners, token)
	}
	return owners
}

// isValidOwner checks that an owner handle contains only valid characters:
// @, letters, digits, hyphens, underscores, and slashes (for org/team).
func isValidOwner(s string) bool {
	if len(s) < 2 {
		return false
	}
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
