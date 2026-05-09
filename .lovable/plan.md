## Commit-In Tag Mirroring & Auto Release Branches — Implementation Plan

Source: `spec/03-commit-in/08-tag-mirroring-and-release-branches.md` (frozen).
Memory: `mem://features/commit-in-tag-mirroring`.

The spec gates implementation on this plan being greenlit, so no code is written until you approve.

### Scope at a glance

- Mirror annotated git tags from the source repo onto the destination NewSha during `commit-in`.
- For tags matching `constants.VersionTagPattern`, also create `release/<TagName>` branches at the NewSha.
- Persist `(MirroredTagName, MirroredReleaseBranch)` per `RewrittenCommit` row.
- Three new flags: `--tags`, `--no-release-branch`, `--release-branch-prefix`. Profile JSON adds `TagsMode`, `CreateReleaseBranch`, `ReleaseBranchPrefix`.
- N-tags-per-commit: extra annotated tags become sibling `Skipped / AdditionalTagAlias` rows.

### Implementation (8 ordered chunks, each ≤ ~150 LOC across small files)

1. **Migration 006 + enum seed.** New `gitmap/store/migrations/006_commit_in_tag_mirroring.sql` (or `migrate_commitin_tagmirror.go` matching the existing migration style). Idempotent `ALTER TABLE RewrittenCommit ADD COLUMN MirroredTagName / MirroredReleaseBranch`, two indexes, `TagsMode` lookup table seeded `Annotated|All|None`, and a new `SkipReason.AdditionalTagAlias` row. Bump schema version. Idempotent-migration ERD parity test updated.

2. **Constants & enums.** Add CLI ID + flag default constants in `constants_commitin.go` / `constants_cli.go` (`FlagTags`, `FlagNoReleaseBranch`, `FlagReleaseBranchPrefix`, `DefaultTagsMode`, `DefaultReleaseBranchPrefix`). Reuse existing `constants.VersionTagPattern` and `constants.ReleaseBranchPrefix` (no siblings). Add `TagsMode` Go enum in `gitmap/cmd/commitin/enums.go`.

3. **Flag parsing + validation.** Extend `parse_flags.go` / `parse_validate.go` for the three flags, with the spec's two validation errors (`--tags=None` + sibling flag → exit 2; prefix not ending `/` → exit 2). Test in `parse_test.go` (table-driven).

4. **Profile JSON.** Extend `profile/` types + resolver to read `TagsMode`, `CreateReleaseBranch`, `ReleaseBranchPrefix` with the standard CLI > profile > default precedence. Round-trip test.

5. **gitutil helpers.** New small file `gitmap/gitutil/tagmirror.go` (≤200 LOC):
   - `AnnotatedTagsAt(repo, sha) ([]TagInfo, error)` — uses `git for-each-ref refs/tags --contains <sha>` + `git cat-file -t` filter, returns tagger ident/date/message verbatim.
   - `CreateAnnotatedTag(repo, name, sha, ident, date, message) error` — wraps `git tag -a` with `GIT_COMMITTER_DATE` / `GIT_TAGGER_DATE` env pin.
   - `CreateBranchAt(repo, name, sha) error` — `git branch <name> <sha>`, idempotent (no-op if already at sha).
   Each function ≤15 lines per the project's code-style rule. Unit tests with a temp repo fixture.

6. **Replay integration (stage 14).** In `gitmap/cmd/commitin/replay/replay.go` (or a new `tagmirror.go` sibling so replay.go stays under the 200-line limit), add inline post-`ApplyCommit` logic per spec §8.4. Single SQLite transaction with the `RewrittenCommit` insert. `--dry-run` short-circuits writes but still records "would mirror" in the run-log.

7. **Run-log + results.** Extend `runlog/tagreplay.go` to surface mirrored tag/branch in summary output and the JSON run-log shape. New result fields are additive (no breaking change to existing JSON consumers).

8. **Acceptance tests T1–T7.** New `gitmap/cmd/commitin/replay/tagmirror_test.go` plus end-to-end coverage in `gitmap/cmd/commitin/e2e/` driving each row of spec §8.7's matrix against a real temp source repo.

### Validation gates I will run before marking done

- `go test ./...` (full suite)
- `golangci-lint run ./...`
- `goldenguard` determinism pre-check on any new fixtures
- `gitmap fix-repo --strict` not applicable (no version bump in this feature)
- File-size + function-size lint enforced by repo-policy CI (already gated)

### Out of scope (intentionally deferred)

- Lightweight-tag dereferencing nuances beyond `--tags=All` literal mirror.
- Conflict policy when destination already has a tag with the same name (treated as failure → `PartiallyFailed`, no overwrite — matches §02 conformance).
- Cross-repo tag pushing — `commit-in` writes locally only; pushing remains a separate `gitmap` command.

### Risk register

| Risk | Mitigation |
|------|-----------|
| Tagger date hashing changes commit identity | We do NOT touch the commit; only the tag object. Verified against §8.1's "re-create the tag, never re-point" rule. |
| Multiple annotated tags per commit | Spec §8.5 N-tags rule: sibling `Skipped/AdditionalTagAlias` rows, NewSha shared. Tested in T7. |
| Migration on existing live DBs | `ADD COLUMN` is idempotent in SQLite; migration runner pattern (see `migrate_commitin.go`) wraps in `IF NOT EXISTS` guard. |
| Schema-version + ERD parity test drift | Update `erd_parity_test.go` in chunk 1 alongside the migration. |

### Delivery cadence

I will deliver the 8 chunks one per `next` so each lands as a reviewable, test-green increment. Chunk 1 (migration + enum seed) ships first because every later chunk depends on the new columns existing.

**Greenlight needed:** approve this plan, push back on any chunk, or pick a different memory follow-up to tackle instead.
