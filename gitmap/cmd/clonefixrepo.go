// Package cmd — clonefixrepo.go: entry points for `gitmap clone-fix-repo`
// (alias `cfr`) and `gitmap clone-fix-repo-pub` (alias `cfrp`).
//
// These are convenience pipelines that chain three existing commands
// in one shot:
//
//	cfr  : clone <url>  →  cd <folder>  →  fix-repo --all
//	cfrp : clone <url>  →  cd <folder>  →  fix-repo --all  →  make-public --yes
//
// Implementation strategy: the chained commands (runFixRepo,
// runMakePublic) all call os.Exit at the end, which would terminate
// our parent process before the next step runs. To stay decoupled
// and side-effect-clean, we shell out to our own binary (resolved
// via os.Executable) for the fix-repo and make-public steps after
// invoking executeDirectClone in-process. This also keeps each
// step's exit code, stdout, and stderr semantics intact.
package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v19/gitmap/clonenext"
	"github.com/alimtvnetwork/gitmap-v19/gitmap/constants"
)

// runCloneFixRepo implements `gitmap clone-fix-repo` (alias cfr).
func runCloneFixRepo(args []string) {
	checkHelp(constants.CmdCloneFixRepo, args)
	runCloneFixRepoPipeline(args, false)
}

// runCloneFixRepoPub implements `gitmap clone-fix-repo-pub` (alias cfrp).
func runCloneFixRepoPub(args []string) {
	checkHelp(constants.CmdCloneFixRepoPub, args)
	runCloneFixRepoPipeline(args, true)
}

// runCloneFixRepoPipeline is the shared core. `makePublic` controls
// whether the optional 3rd step (visibility flip) runs.
func runCloneFixRepoPipeline(args []string, makePublic bool) {
	url, folderName, noVSCodeSync, requireVersion := parseCloneFixRepoArgs(args)
	if len(url) == 0 {
		fmt.Fprint(os.Stderr, constants.ErrCloneFixRepoUsage)
		os.Exit(constants.ExitCloneFixRepoBadFlag)
	}

	absPath := resolveCloneTargetFolder(url, folderName)
	requireOnline()
	executeDirectClone(url, folderName, true, false, "", noVSCodeSync)

	if err := os.Chdir(absPath); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCloneFixRepoChdirFmt, absPath, err)
		os.Exit(constants.ExitCloneFixRepoChdir)
	}

	maybeRunFixRepoStep(absPath, requireVersion)
	if makePublic {
		runChainedGitmapStep([]string{constants.CmdMakePublic, "--" + constants.FlagVisYes})
	}
	fmt.Printf(constants.MsgCloneFixRepoDone, absPath)
}

// maybeRunFixRepoStep runs `fix-repo --all` only when the repo folder
// name carries a `-vN` suffix. Repos without a version suffix have
// nothing the rewriter can target, so we skip with a one-line notice.
// `--require-version` restores the strict (exit-4) failure mode for
// CI pipelines that want the old contract.
func maybeRunFixRepoStep(absPath string, requireVersion bool) {
	parsed := clonenext.ParseRepoName(filepath.Base(absPath))
	if parsed.HasVersion {
		runChainedGitmapStep([]string{constants.CmdFixRepo, "--" + constants.FixRepoFlagAll})

		return
	}
	if requireVersion {
		fmt.Fprintf(os.Stderr, constants.ErrCloneFixRepoNeedVersion, parsed.BaseName)
		os.Exit(constants.ExitCloneFixRepoChainFailed)
	}
	fmt.Printf(constants.MsgCloneFixRepoSkipNoVer, parsed.BaseName)
}

// parseCloneFixRepoArgs returns (url, folderName, noVSCodeSync, requireVersion).
// First non-flag arg is the URL; second non-flag is the destination folder.
// Recognized flags: --no-vscode-sync, --require-version.
func parseCloneFixRepoArgs(args []string) (string, string, bool, bool) {
	positional := make([]string, 0, len(args))
	noVSCodeSync := false
	requireVersion := false
	syncFlag := "--" + constants.FlagNoVSCodeSync
	reqFlag := "--" + constants.FlagRequireVersion
	for _, a := range args {
		if a == syncFlag {
			noVSCodeSync = true

			continue
		}
		if a == reqFlag {
			requireVersion = true

			continue
		}
		if len(a) > 0 && a[0] != '-' {
			positional = append(positional, a)
		}
	}
	url := ""
	folder := ""
	if len(positional) > 0 {
		url = positional[0]
	}
	if len(positional) > 1 {
		folder = positional[1]
	}

	return url, folder, noVSCodeSync, requireVersion
}

// resolveCloneTargetFolder mirrors the folder-naming logic in
// executeDirectClone so we know which directory to cd into after
// the clone step finishes. Versioned URLs auto-flatten to BaseName.
func resolveCloneTargetFolder(url, folderName string) string {
	if len(folderName) == 0 {
		repoName := repoNameFromURL(url)
		parsed := clonenext.ParseRepoName(repoName)
		if parsed.HasVersion {
			folderName = parsed.BaseName
		} else {
			folderName = repoName
		}
	}
	abs, err := filepath.Abs(folderName)
	if err != nil {
		return folderName
	}

	return abs
}

// runChainedGitmapStep re-execs the current gitmap binary with the
// given args, streaming stdin/stdout/stderr through. Any non-zero
// exit propagates immediately so the pipeline halts on first failure.
func runChainedGitmapStep(args []string) {
	bin, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCloneFixRepoExecFmt, err)
		os.Exit(constants.ExitCloneFixRepoChainFailed)
	}
	cmd := exec.Command(bin, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if runErr := cmd.Run(); runErr != nil {
		var exitErr *exec.ExitError
		if errors.As(runErr, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(os.Stderr, constants.ErrCloneFixRepoExecFmt, runErr)
		os.Exit(constants.ExitCloneFixRepoChainFailed)
	}
}
