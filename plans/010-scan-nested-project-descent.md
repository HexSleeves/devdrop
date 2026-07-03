# Plan 010: Stop scan from creating one project per nested package manifest in monorepos

> **Executor instructions**: Follow this plan step by step. Run every
> verification command and confirm the expected result before moving on. On any
> STOP condition, stop and report. When done, update this plan's status row in
> `plans/README.md` — unless a reviewer told you they maintain the index.
>
> **Drift check (run first)**: `git diff --stat 595d158..HEAD -- internal/devspace/workspace.go internal/devspace/hardening_test.go internal/devspace/devspace_test.go`
> On drift, re-verify excerpts; on mismatch, STOP.

## Status

- **Priority**: P2
- **Effort**: M
- **Risk**: MED (changes which projects a scan discovers; needs both-direction regression tests)
- **Depends on**: 006 (merge semantics — land first so rescan tests aren't fighting two bugs)
- **Category**: bug
- **Planned at**: commit `595d158`, 2026-07-02

## Why this matters

`ScanWorkspace` prunes the walk (`filepath.SkipDir`) only for **git** projects.
A non-git directory tracked via a dependency marker (`package.json`, etc.) or a
`.env` file falls through and its subdirectories are re-tested with the same
marker check. A pnpm/turborepo-style monorepo — root `package.json` plus
`packages/*/package.json` — therefore produces one manifest `Project` per
nested package in addition to the root, with overlapping paths. The manifest
bloats, syncs that bloat to every machine, and `plan`/`apply` enumerate
redundant nested entries. The asymmetry with git repos (which are correctly
skipped) reads as unintended.

## Current state

`internal/devspace/workspace.go:39-84`, the WalkDir callback:

```go
info := gitInfo(path)
hasMarker := info.IsRepo || hasDependencyMarker(path) || exists(filepath.Join(path, ".env"))
if !hasMarker {
    summary.UntrackedFolders++
    return nil
}
p := projectFromPath(clean, path, info)
upsertProject(&m, p)
...
if info.IsRepo {
    return filepath.SkipDir
}
return nil            // ← non-git project: walk keeps descending
```

Existing regression net: `TestHardeningScanIgnoresDependencyFoldersAndNestedRepos`
(`hardening_test.go:93-114`) asserts nested **git** repos inside a tracked repo
are not tracked; there is no nested-non-git case.

`ignoredName` + `DefaultIgnores` (`types.go:23-35`) already prune
`node_modules` etc., so the duplication hits real source directories like
`packages/`, `apps/`.

### Chosen semantics (decision made — implement exactly this)

Once a directory is tracked as a **local** (non-git) project, its descendants
are only eligible to become projects if they are **git repositories**. Marker
files (`package.json`, `.env`, …) below a tracked local project do NOT create
additional projects. Nested git repos remain tracked (that behavior is
deliberate and covered by the walk continuing).

Rationale: a monorepo is one project; an embedded git checkout inside a
non-git tree is genuinely separate.

## Commands you will need

| Purpose | Command | Expected on success |
|---------|---------|---------------------|
| Full gate | `make verify` | exit 0 |
| Targeted | `go test ./internal/devspace -run 'Scan|Hardening' -v` | PASS |

## Scope

**In scope**:
- `internal/devspace/workspace.go` (the WalkDir callback only)
- `internal/devspace/devspace_test.go` (new regression test)

**Out of scope**:
- `hasDependencyMarker`, `projectFromPath`, `upsertProject` — unchanged.
- Removing already-duplicated entries from existing manifests — that is
  Plan 012's `project remove` territory (manual cleanup), not auto-pruning
  (non-destructive invariant).
- `hardening_test.go` — extend `devspace_test.go` instead; hardening contract
  stays untouched.

## Git workflow

- Branch: `advisor/010-scan-monorepo-descent`
- Conventional commit, e.g. `fix(scan): do not track nested package manifests under an already-tracked local project`

## Steps

### Step 1: Track the active local-project root in the walk

`filepath.WalkDir` visits depth-first in lexical order, so a simple "current
local project root" string suffices:

```go
var localRoot string // abs path of the innermost tracked non-git project being descended
```

In the callback, before the marker check:

```go
underLocal := localRoot != "" && strings.HasPrefix(path+string(filepath.Separator), localRoot+string(filepath.Separator))
if !underLocal {
    localRoot = "" // walked out of the previous local project's subtree
}
```

Then adjust eligibility: when `underLocal`, only `info.IsRepo` qualifies:

```go
hasMarker := info.IsRepo || (!underLocal && (hasDependencyMarker(path) || exists(filepath.Join(path, ".env"))))
```

And when a non-git project is tracked (`!info.IsRepo` after the upsert), set
`localRoot = path` if `localRoot == ""`.

(Exact code shape is the executor's; the semantics table is not negotiable.)

**Verify**: `go build ./...` → exit 0

### Step 2: Regression tests (both directions)

In `devspace_test.go`, add `TestScanTreatsMonorepoAsOneProject`:

- Build a temp workspace: `mono/package.json`, `mono/packages/a/package.json`,
  `mono/packages/b/package.json`, and a nested git repo `mono/vendor-tool/`
  (init a real repo — see how existing tests shell out to git, e.g. in
  `workspace_sync_test.go`, or use `exec.Command("git", "init")` as
  `devspace_test.go` already imports `os/exec`).
- Run `ScanWorkspace`; load manifest; assert exactly **2** projects: `mono`
  (local) and `mono/vendor-tool` (git). Assert `packages/a`, `packages/b` are
  absent.
- Second direction — standalone siblings still tracked: `apps/x/package.json`
  and `apps/y/package.json` where `apps/` itself has **no** marker → both `x`
  and `y` tracked as before (the fix must not collapse mere siblings).

**Verify**: `go test ./internal/devspace -run TestScanTreatsMonorepo -v` → PASS

### Step 3: Full suite — especially existing scan expectations

**Verify**: `make verify` → exit 0. If any existing test expected nested-marker
tracking, see STOP conditions.

## Test plan

Step 2's two-direction test plus the untouched hardening scan tests. Note:
`summary.UntrackedFolders` counts will shift for descendants of local
projects; assert on project sets, not on untracked counts.

## Done criteria

- [ ] `make verify` exits 0
- [ ] `TestScanTreatsMonorepoAsOneProject` passes (both directions)
- [ ] `TestHardeningScanIgnoresDependencyFoldersAndNestedRepos` passes unmodified
- [ ] No files outside scope modified (`git status`)
- [ ] `plans/README.md` status row updated

## STOP conditions

- An existing test asserts nested non-git markers ARE tracked — the current
  behavior may be deliberate after all; report the test name and stop.
- The lexical-order assumption of WalkDir (children visited immediately after
  parent, before siblings) fails on some input — report; do not switch to a
  path-set approach without review.
- Scan results become order-dependent between platforms (case-insensitive FS) —
  report.

## Maintenance notes

- Existing manifests that already contain nested duplicates are NOT cleaned by
  this fix (scan never deletes). Cleanup path: Plan 012's `project remove`.
- If a future "workspaces-aware" feature wants per-package tracking inside
  monorepos, it should be explicit opt-in metadata, not walk-descent behavior.
- Reviewer: the `underLocal` prefix check must use the separator-suffixed
  form (as shown) to avoid `mono-tools` matching `mono`.
