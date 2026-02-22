package formatter

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kevin-robayna/codeowner/internal/owner"
)

// CodeOwners formats mappings as a GitHub CODEOWNERS file.
// Output is sorted and grouped: root files first, then hidden-directory files,
// then everything else. Within each section, entries are grouped by their
// top-2-level directory with blank lines between groups.
func CodeOwners(mappings []owner.Mapping) string {
	sorted := make([]owner.Mapping, len(mappings))
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

// groupKey returns the grouping key for a path based on its directory
// structure: "" for root files, the first directory for single-depth paths,
// or the first two directories for deeper paths.
func groupKey(path string) string {
	dir := filepath.Dir(path)
	if dir == "." {
		return ""
	}
	parts := strings.SplitN(dir, string(filepath.Separator), 3)
	if len(parts) >= 2 {
		return parts[0] + string(filepath.Separator) + parts[1]
	}
	return parts[0]
}
