# gitmap open (alias: op)

Open the current repository in **both** GitHub Desktop and VS Code in
one shot. Detects the repo from the current working directory
(prefers `git rev-parse --show-toplevel`, falls back to cwd) and
re-injects on every call so freshly cloned or moved repos always
land in both tools.

## Usage

    gitmap open
    gitmap op

No arguments. Run from anywhere inside the repo folder.

## What it does

1. Resolves the repo root (git toplevel, or cwd if not a git repo).
2. Best-effort DB upsert when an `origin` remote is configured —
   same logic as `gitmap inject`. Local-only folders are accepted
   silently; the upsert step is just skipped.
3. Registers the folder with GitHub Desktop. Already-registered
   repos are silently re-added (no duplication).
4. Opens the folder in VS Code.

## Examples

    cd ~/dev/macro-ahk-v31
    gitmap open
      → adds to Desktop, opens VS Code window, upserts DB row.

    cd /tmp/scratch        # not a git repo
    gitmap op
      → DB step skipped (no origin), Desktop + VS Code still proceed.

## See also

- [inject](inject.md) — Inject any folder (positional path argument)
- [code](code.md) — VS Code only (Project Manager registration)
- [cd](cd.md) — Navigate to a tracked repo by slug
