package owner

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Mapping holds a file path and its code owners.
type Mapping struct {
	Path   string
	Owners []string
}

// DefaultPrefix is the default annotation prefix to search for.
const DefaultPrefix = "CodeOwner:"

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

// ParseDir walks a directory and returns all CodeOwner mappings.
func ParseDir(root, prefix string) ([]Mapping, error) {
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

		if owners := ParseFile(path, prefix); len(owners) > 0 {
			rel, relErr := filepath.Rel(root, path)
			if relErr != nil {
				rel = path
			}
			mappings = append(mappings, Mapping{Path: rel, Owners: owners})
		}
		return nil
	})

	return mappings, err
}

// FormatCodeOwners formats mappings as a GitHub CODEOWNERS file.
// Output is sorted and grouped: root files first, then hidden-directory files,
// then everything else. Within each section, entries are grouped by their
// top-2-level directory with blank lines between groups.
func FormatCodeOwners(mappings []Mapping) string {
	sorted := make([]Mapping, len(mappings))
	copy(sorted, mappings)

	sort.Slice(sorted, func(i, j int) bool {
		si, sj := pathSection(sorted[i].Path), pathSection(sorted[j].Path)
		if si != sj {
			return si < sj
		}
		gi, gj := groupKey(sorted[i].Path), groupKey(sorted[j].Path)
		if gi != gj {
			return gi < gj
		}
		return sorted[i].Path < sorted[j].Path
	})

	var b strings.Builder
	prevGroup := ""
	for i, m := range sorted {
		g := groupKey(m.Path)
		if i > 0 && g != prevGroup {
			b.WriteByte('\n')
		}
		prevGroup = g
		fmt.Fprintf(&b, "%s %s\n", m.Path, strings.Join(m.Owners, " "))
	}
	return b.String()
}

// pathSection returns the sort section for a path:
// 0 for root-level files, 1 for hidden directories, 2 for everything else.
func pathSection(path string) int {
	if !strings.Contains(path, string(filepath.Separator)) {
		return 0
	}
	first := strings.SplitN(path, string(filepath.Separator), 2)[0]
	if strings.HasPrefix(first, ".") {
		return 1
	}
	return 2
}

// groupKey returns the grouping key for a path: the first two path segments
// joined, or the first segment for single-depth paths, or "" for root files.
func groupKey(path string) string {
	parts := strings.SplitN(path, string(filepath.Separator), 3)
	if len(parts) == 1 {
		return ""
	}
	if len(parts) >= 2 {
		return parts[0] + string(filepath.Separator) + parts[1]
	}
	return parts[0]
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
