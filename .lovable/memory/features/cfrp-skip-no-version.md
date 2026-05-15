---
name: cfrp/cfr skip fix-repo on no-vN
description: clone-fix-repo / clone-fix-repo-pub gracefully skip the fix-repo step (with a one-line notice) when the cloned repo has no -vN suffix; --require-version restores the strict E_NO_VERSION_SUFFIX exit. v4.43.0+.
type: feature
---

# Feature: cfrp/cfr skip fix-repo on no-vN suffix (v4.43.0+)

## Rule
`gitmap clone-fix-repo` (cfr) and `gitmap clone-fix-repo-pub` (cfrp)
MUST NOT fail their pipeline when the cloned repo's folder name has
no `-vN` suffix. Skip the fix-repo step with a single notice line and
continue to make-public (cfrp only).

## Behavior
- Default: skip with `fix-repo: skipped (repo "<name>" has no -vN suffix, nothing to rewrite)`. Exit 0.
- `--require-version` flag: restore strict mode → exit `ExitCloneFixRepoChainFailed` (10) with a clear "ERROR --require-version set" message.

## Why
Original report (May 2026): `gitmap cfrp https://github.com/.../img-pdf`
clone succeeded, then `fix-repo: ERROR no -vN suffix on repo name "img-pdf"`
killed the pipeline AFTER the clone + Desktop + VS Code side effects had
already run. Result: the user got the artifacts but a non-zero exit and a
scary error. The fix-repo step is meaningless on non-versioned repos —
there's nothing to rewrite — so the correct default is "skip silently
(with a notice)".

## Files
- `gitmap/cmd/clonefixrepo.go::maybeRunFixRepoStep` — gates the chained step on `clonenext.ParseRepoName(...).HasVersion`.
- `gitmap/constants/constants_clonefixrepo.go` — `FlagRequireVersion`, `MsgCloneFixRepoSkipNoVer`, `ErrCloneFixRepoNeedVersion`.

## Spec
- `spec/02-app-issues/26-cfrp-no-version-suffix-hard-error.md` (TODO follow-up)
