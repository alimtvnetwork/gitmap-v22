package cmd

import (
	"github.com/alimtvnetwork/gitmap-v19/gitmap/constants"
)

// flatCtxEntry is the macOS/Linux representation of a menu item — those
// platforms do not support arbitrary nested cascades inside Finder
// Services / Nautilus scripts, so we flatten "Category ▸ Child" into a
// single labelled entry. Single source of truth: ctxMenu().
type flatCtxEntry struct {
	Label string   // "gitmap: Release — Release next (bump minor)"
	Slug  string   // filesystem-safe id derived from label, "gitmap-release-release-next"
	Args  []string // {"release", "--bump", "minor"}
	Mode  constants.CtxMode
	Exe   string // override executable; empty => use the gitmap binary
}

// flattenCtxMenu walks ctxMenu() into a flat list. Categories with
// children become "<prefix>: <Category> — <Child>"; top-level leaves
// become "<prefix>: <Label>". Order is preserved.
func flattenCtxMenu() []flatCtxEntry {
	var out []flatCtxEntry
	for _, e := range ctxMenu() {
		if len(e.Children) > 0 {
			for _, c := range e.Children {
				out = append(out, flatEntry(e.MUIVerb, c))
			}

			continue
		}
		out = append(out, flatEntry("", e))
	}

	return out
}

// flatEntry builds one flatCtxEntry from a category + child (or empty
// category for top-level leaves).
func flatEntry(category string, e ctxEntry) flatCtxEntry {
	label := constants.CtxFlatPrefix + constants.CtxFlatSeparator
	if category != "" {
		label += category + constants.CtxFlatChildJoiner
	}
	label += e.MUIVerb

	return flatCtxEntry{
		Label: label,
		Slug:  slugifyCtx(label),
		Args:  append([]string(nil), e.Args...),
		Mode:  e.Mode,
		Exe:   e.Exe,
	}
}

// slugifyCtx returns a filesystem-safe id: lowercase alphanumerics
// joined by "-". Used as workflow folder, .desktop file, and Nautilus
// script base names.
func slugifyCtx(s string) string {
	out := make([]byte, 0, len(s))
	dashed := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= 'A' && c <= 'Z':
			out = append(out, c+32)
			dashed = false
		case (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9'):
			out = append(out, c)
			dashed = false
		default:
			if !dashed && len(out) > 0 {
				out = append(out, '-')
				dashed = true
			}
		}
	}
	for len(out) > 0 && out[len(out)-1] == '-' {
		out = out[:len(out)-1]
	}

	return string(out)
}
