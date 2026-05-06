// Package cmd — vscodepmsync_testhelper_test.go: test fixtures and
// helpers shared by vscodepmsync_test.go. Kept in a separate file so
// the primary test file stays under the 200-line code-style cap.
package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v18/gitmap/vscodepm"
)

// setupVSCodePMSyncFixture creates a temp HOME, a real repo dir, and
// a single-entry projects.json file pointing at the repo. Returns the
// repo path and a restore function that resets HOME / XDG_CONFIG_HOME.
func setupVSCodePMSyncFixture(t *testing.T) (string, func()) {
	t.Helper()
	tmp := t.TempDir()
	repoDir := filepath.Join(tmp, "demo-repo")
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		t.Fatalf("mkdir repo: %v", err)
	}

	jsonPath := vscodepmSyncFixturePath(t, tmp)
	writeVSCodePMSyncSeed(t, jsonPath, repoDir)

	prevHome, prevXDG := os.Getenv("HOME"), os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("HOME", tmp)
	os.Setenv("XDG_CONFIG_HOME", filepath.Join(tmp, ".config"))

	return repoDir, func() {
		os.Setenv("HOME", prevHome)
		os.Setenv("XDG_CONFIG_HOME", prevXDG)
	}
}

// vscodepmSyncFixturePath returns the projects.json path inside the
// faux XDG_CONFIG_HOME and ensures the parent directory exists.
func vscodepmSyncFixturePath(t *testing.T, tmp string) string {
	t.Helper()
	dir := filepath.Join(tmp, ".config", "Code", "User",
		"globalStorage", "alefragnani.project-manager")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir extdir: %v", err)
	}

	return filepath.Join(dir, "projects.json")
}

// writeVSCodePMSyncSeed writes a single-entry projects.json with a
// pre-existing "user" tag so we can prove UNION preservation.
func writeVSCodePMSyncSeed(t *testing.T, jsonPath, repoDir string) {
	t.Helper()
	seed := []vscodepm.Entry{{
		Name: "demo-repo", RootPath: repoDir,
		Paths: []string{}, Tags: []string{"user"}, Enabled: true,
	}}
	data, err := json.MarshalIndent(seed, "", "\t")
	if err != nil {
		t.Fatalf("marshal seed: %v", err)
	}
	if err := os.WriteFile(jsonPath, data, 0o644); err != nil {
		t.Fatalf("write seed: %v", err)
	}
}

// loadVSCodePMSyncProjectsJSON reads the on-disk projects.json after
// the runner finished and unmarshals it for assertion.
func loadVSCodePMSyncProjectsJSON(t *testing.T) []vscodepm.Entry {
	t.Helper()
	got, err := vscodepm.ListEntries()
	if err != nil {
		t.Fatalf("list entries: %v", err)
	}

	return got
}

// containsTag reports whether tag appears in tags.
func containsTag(tags []string, tag string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}

	return false
}
