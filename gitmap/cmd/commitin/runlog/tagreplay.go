package runlog

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"sync"

	"github.com/alimtvnetwork/gitmap-v18/gitmap/constants"
)

// TagReplayFacts is the persistence-layer projection of one mirrored
// annotated tag. Mirrors spec §9.4 columns. Caller resolves the
// `RewrittenCommit` row id BEFORE invoking RecordTagReplay so this
// stays a pure write helper (no FK guessing in the runlog layer).
//
// `Outcome` MUST be a constants.TagReplayOutcome* literal — the enum
// FK is resolved here via the standard mirror-table lookup pattern.
type TagReplayFacts struct {
	SourceTagName         string
	SourceTagSha          string
	SourceCommitSha       string
	DestTagSha            string // empty for DryRun / Failed / Skipped
	DestCommitSha         string // empty for DryRun
	MirroredReleaseBranch string // empty when no branch was mirrored
	IsVersionTag          bool
	Outcome               string
}

// TagReplayLookup is the cross-run idempotency view returned by
// LookupTagReplay (spec §9.5). Every field carries the destination
// state recorded by the previous successful mirror; an empty
// DestCommitSha means the prior row was Created/AlreadyExists under
// a flow that did NOT capture the dest commit (defensive — should
// not occur for those two outcomes per §9.4).
type TagReplayLookup struct {
	DestTagSha            string
	DestCommitSha         string
	MirroredReleaseBranch string
}

// ErrTagReplayMiss is returned by LookupTagReplay when no row matches
// the (SourceTagName, SourceTagSha) pair under a Created /
// AlreadyExists outcome. Callers MUST use errors.Is to detect this
// (Core memory: zero-swallow, errors.Is everywhere).
var ErrTagReplayMiss = errors.New("runlog: tag replay lookup miss")

// RecordTagReplay persists one CommitInReplayMap row. Empty string
// fields are stored as SQL NULL where the column is nullable per
// spec §9.4 (DestTagSha, DestCommitSha, MirroredReleaseBranch).
func RecordTagReplay(db *sql.DB, runID, rewrittenID int64, f TagReplayFacts) (int64, error) {
	outcomeID, err := lookupEnumID(db, constants.TableCommitInTagOutcome,
		"TagReplayOutcomeId", f.Outcome)
	if err != nil {
		return 0, fmt.Errorf("runlog: lookup tag outcome %q: %w", f.Outcome, err)
	}
	res, err := db.Exec(constants.SQLInsertCommitInReplayMap,
		runID, rewrittenID,
		f.SourceTagName, f.SourceTagSha, f.SourceCommitSha,
		nullIfEmpty(f.DestTagSha), nullIfEmpty(f.DestCommitSha),
		nullIfEmpty(f.MirroredReleaseBranch),
		boolToInt(f.IsVersionTag), outcomeID,
	)
	if err != nil {
		return 0, fmt.Errorf("runlog: insert CommitInReplayMap: %w", err)
	}
	return res.LastInsertId()
}

// LookupTagReplay implements the §9.5 cross-run short-circuit query.
// Returns ErrTagReplayMiss when no prior Created / AlreadyExists row
// matches. Pure read; safe to call before §08 attempts `git tag`.
func LookupTagReplay(db *sql.DB, sourceTagName, sourceTagSha string) (TagReplayLookup, error) {
	var got TagReplayLookup
	var dt, dc, mb sql.NullString
	err := db.QueryRow(constants.SQLSelectCommitInReplayLookup,
		sourceTagName, sourceTagSha).Scan(&dt, &dc, &mb)
	if errors.Is(err, sql.ErrNoRows) {
		return got, ErrTagReplayMiss
	}
	if err != nil {
		return got, fmt.Errorf("runlog: lookup tag replay %q: %w", sourceTagName, err)
	}
	got.DestTagSha = dt.String
	got.DestCommitSha = dc.String
	got.MirroredReleaseBranch = mb.String
	return got, nil
}

// IsAnnotatedSemverVersionTag returns true iff `tagName` matches the
// canonical version-tag pattern (spec §08, §9.4 `IsVersionTag`). The
// regex compile is cached once via sync.Once because the pattern is a
// package constant and never mutates.
//
// NOTE: this function is name-only. It deliberately does NOT take an
// "is annotated" boolean — §08's tag walker passes ONLY annotated tags
// to this helper (lightweight tags are filtered upstream). Folding
// the annotated check in here would invert the dependency.
func IsAnnotatedSemverVersionTag(tagName string) bool {
	return versionTagRegex().MatchString(tagName)
}

var (
	versionTagRegexOnce sync.Once
	versionTagRegexVal  *regexp.Regexp
)

func versionTagRegex() *regexp.Regexp {
	versionTagRegexOnce.Do(func() {
		versionTagRegexVal = regexp.MustCompile(constants.VersionTagPattern)
	})
	return versionTagRegexVal
}

// nullIfEmpty maps "" → SQL NULL so spec §9.4's NULL-on-dry-run /
// NULL-on-failed contract is honoured without the caller juggling
// `any` typed parameters.
func nullIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
