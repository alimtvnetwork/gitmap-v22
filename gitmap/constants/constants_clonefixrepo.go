// Package constants — clone-fix-repo command IDs, messages, and
// exit codes for `gitmap clone-fix-repo` (cfr) and
// `gitmap clone-fix-repo-pub` (cfrp).
//
// These commands chain `clone` → `fix-repo` (and optionally
// `make-public`) into a single invocation. See cmd/clonefixrepo.go.
package constants

// gitmap:cmd top-level
// Clone-fix-repo command IDs and short aliases.
const (
	CmdCloneFixRepo         = "clone-fix-repo"
	CmdCloneFixRepoAlias    = "cfr"
	CmdCloneFixRepoPub      = "clone-fix-repo-pub"
	CmdCloneFixRepoPubAlias = "cfrp"
)

// Clone-fix-repo help-line entries surfaced by `gitmap help`.
const (
	HelpCloneFixRepo    = "  clone-fix-repo (cfr) <url> [folder]      Clone, then run fix-repo --all in the new folder"
	HelpCloneFixRepoPub = "  clone-fix-repo-pub (cfrp) <url> [folder] Clone, fix-repo --all, then make-public --yes"
)

// Clone-fix-repo user-facing messages and errors.
const (
	MsgCloneFixRepoDone        = "clone-fix-repo: pipeline completed in %s\n"
	MsgCloneFixRepoSkipNoVer   = "  fix-repo: skipped (repo %q has no -vN suffix, nothing to rewrite)\n    pass --require-version to fail instead.\n"
	ErrCloneFixRepoUsage       = "clone-fix-repo: ERROR <url> is required\n  usage: gitmap clone-fix-repo <url> [folder]\n  usage: gitmap clone-fix-repo-pub <url> [folder]\n"
	ErrCloneFixRepoChdirFmt    = "clone-fix-repo: ERROR cannot cd into %q: %v\n"
	ErrCloneFixRepoExecFmt     = "clone-fix-repo: ERROR could not run chained step: %v\n"
	ErrCloneFixRepoNeedVersion = "clone-fix-repo: ERROR --require-version set but repo %q has no -vN suffix\n"
)

// Clone-fix-repo flags.
const (
	FlagRequireVersion = "require-version"
)

// Clone-fix-repo exit codes.
const (
	ExitCloneFixRepoOk          = 0
	ExitCloneFixRepoBadFlag     = 6
	ExitCloneFixRepoChdir       = 9
	ExitCloneFixRepoChainFailed = 10
)
