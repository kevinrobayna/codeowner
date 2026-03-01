package formatter

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kevin-robayna/codeowner/internal/scanning"
)

// CodeOwners formats mappings as a GitHub CODEOWNERS file.
// Output is sorted and grouped: root files first, then hidden-directory files,
// then everything else. Within each section, entries are grouped by their
// top-2-level directory with blank lines between groups.
func CodeOwners(mappings []scanning.Mapping) string {
	var protect *scanning.Mapping
	sorted := make([]scanning.Mapping, 0, len(mappings))
	for i := range mappings {
		if mappings[i].Path == "CODEOWNERS" {
			protect = &mappings[i]
		} else {
			sorted = append(sorted, mappings[i])
		}
	}

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
	if protect != nil {
		fmt.Fprintf(&b, "%s %s\n", protect.Path, strings.Join(protect.Owners, " "))
		if len(sorted) > 0 {
			b.WriteByte('\n')
		}
	}
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

// stripRoot removes the leading "/" root-anchor prefix from a CODEOWNERS path.
func stripRoot(path string) string {
	return strings.TrimPrefix(path, "/")
}

// pathSection returns the sort section for a path:
// 0 for root-level files, 1 for hidden directories, 2 for everything else.
func pathSection(path string) int {
	p := stripRoot(path)
	if !strings.Contains(p, "/") {
		return 0
	}
	first := strings.SplitN(p, "/", 2)[0]
	if strings.HasPrefix(first, ".") {
		return 1
	}
	return 2
}

// groupKey returns the grouping key for a path based on its directory
// structure: "" for root files, the first directory for single-depth paths,
// or the first two directories for deeper paths.
func groupKey(path string) string {
	p := stripRoot(path)
	idx := strings.LastIndex(p, "/")
	if idx < 0 {
		return ""
	}
	dir := p[:idx]
	parts := strings.SplitN(dir, "/", 3)
	if len(parts) >= 2 {
		return parts[0] + "/" + parts[1]
	}
	return parts[0]
}
