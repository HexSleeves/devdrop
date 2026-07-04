# Plan 011: Scope watch-mode refreshes to the projects that actually changed

> **Executor instructions**: Follow this plan step by step. Run every
> verification command and confirm the expected result before moving on. On any
> STOP condition, stop and report. When done, update this plan's status row in
> `plans/README.md` — unless a reviewer told you they maintain the index.
>
> **Drift check (run first)**: `git diff --stat 595d158..HEAD -- internal/devspace/watch.go internal/devspace/workspace.go internal/devspace/watch_test.go`
> On drift, re-verify excerpts; on mismatch, STOP.

## Status

- **Priority**: P3
- **Effort**: M
- **Risk**: MED (must not change full-scan semantics for scan/plan/doctor; removed-project reconciliation must survive)
- **Depends on**: 010 (scan descent semantics must be settled first), 009 (locking — the refresh runs under the app lock)
- **Category**: perf
- **Planned at**: commit `595d158`, 2026-07-02

## Why this matters

Every debounced watch event triggers `RefreshWorkspaceForWatch`, which runs a
**full** `ScanWorkspace`: a WalkDir over the whole workspace plus `gitInfo`
per project — and `gitInfo` shells out up to ~8 git subprocesses per repo
(`rev-parse`, `branch --show-current`, `remote`, `config --get`, `remote
get-url`, `rev-parse --short HEAD`, `status --porcelain`, `symbolic-ref`).
With dozens of tracked repos, one file save costs a multi-second, fork-heavy
sweep — continuously, for the life of the watch process. The debounce
coalesces bursts but never scopes the work.

## Current state

- `internal/devspace/watch.go:105-116` — `runRefresh` closure calls
  `RefreshWorkspaceForWatch(mode)` unconditionally; the fsnotify event that
  scheduled it is discarded (`schedule()` is called with no arguments,
  `watch.go:167-169`).
- `internal/devspace/watch.go:40-73` — `RefreshWorkspaceForWatch` = full
  `ScanWorkspace()` + manifest push (git or hosted mode).
- `internal/devspace/workspace.go:36-103` — `ScanWorkspace` walks everything,
  then reconciles projects no longer on disk (lines 88-97: manifest projects
  not `seen` get their state refreshed via `gitInfo(full)`).
- `internal/devspace/git.go:34-77` — `gitInfo`'s subprocess fan-out.
- `internal/devspace/watch_test.go` exists (event relevance, debounce-level
  tests) — extend it.

### Chosen design (implement exactly this)

1. Collect the workspace-relative **top-level project directories** touched by
   events between refreshes: in the event loop, map `event.Name` to the
   tracked project whose path prefixes it (load once per refresh from the
   manifest, not per event), and accumulate into a `map[string]bool`
   (`pending`). An event that maps to **no** tracked project (new top-level
   dir, deletes, renames) sets a `fullScan` flag instead.
2. On timer fire: if `fullScan` or `pending` is empty → current behavior
   (full `RefreshWorkspaceForWatch`). Otherwise call a new
   `RefreshProjectsForWatch(mode, changed []string)` that:
   - loads config/manifest/state once,
   - for each changed project only: re-runs `gitInfo` + `stateForProject`
     (and `projectFromPath`+`upsertProject` so manifest fields refresh),
   - saves manifest + state, then does the same push step as
     `RefreshWorkspaceForWatch`.
3. Safety valve: force a full scan every 10th refresh or every 5 minutes,
   whichever first (drift catch-all). Keep both constants as package-level
   `var`s so tests can shorten them.
4. `ScanWorkspace` itself is untouched — `scan`/`plan`/`doctor` callers keep
   full-scan behavior.

## Commands you will need

| Purpose | Command | Expected on success |
|---------|---------|---------------------|
| Full gate | `make verify` | exit 0 |
| Targeted | `go test ./internal/devspace -run Watch -v` | PASS |
| Race | `go test ./internal/devspace -run Watch -race` | PASS |

## Scope

**In scope**:

- `internal/devspace/watch.go`
- `internal/devspace/workspace.go` (ONLY to add `RefreshProjectsForWatch` /
  a scoped-refresh helper; `ScanWorkspace` body unchanged)
- `internal/devspace/watch_test.go`

**Out of scope**:

- `git.go` — do not "optimize" `gitInfo`'s subprocess count here (separate
  concern, not planned).
- `ScanWorkspace`'s walk or reconciliation logic.
- Push/sync behavior (`PushWorkspaceManifest`/`PushHostedManifest`) — reuse
  as-is.

## Git workflow

- Branch: `advisor/011-watch-incremental-refresh`
- Conventional commit, e.g. `perf(watch): rescan only changed projects on debounced refresh`

## Steps

### Step 1: Thread changed paths from the event loop

In `WatchWorkspace`'s select loop (`watch.go:149-181`): where
`watchEventRelevant` passes, resolve `event.Name` against the tracked project
paths and either add to `pending` or set `fullScan = true`. Reset both after
each refresh. Project paths can be cached and refreshed after each refresh
(they only change when a scan runs — which is exactly a refresh).

**Verify**: `go build ./...` → exit 0

### Step 2: Implement `RefreshProjectsForWatch`

In `workspace.go`, following `ScanWorkspace`'s structure for the per-project
body (lines 61-79: `gitInfo` → `projectFromPath` → `upsertProject` →
`stateForProject`). Do NOT duplicate the walk; iterate the given relative
paths, resolving each through `safeWorkspacePath` (never join by hand — repo
invariant).

Return the same `WatchRefresh` shape; `Summary` for a scoped refresh reports
counts over the refreshed subset (document that in a comment on the field
usage — or leave Summary zeroed except found counts; pick one and test it).

**Verify**: `go test ./internal/devspace -run Watch -v` → existing tests PASS

### Step 3: Wire scoped/full decision + safety valve in `runRefresh`

Full scan when: `fullScan` set, `pending` empty, refresh counter % 10 == 0,
or 5 minutes since last full scan.

**Verify**: `go build ./...` → exit 0

### Step 4: Tests

In `watch_test.go` (its existing tests show how to drive `WatchWorkspace`
with a context and `OnRefresh` callback — follow that):

1. `TestWatchScopedRefreshOnlyTouchesChangedProject`: workspace with two git
   projects; instrument by asserting on `WatchRefresh` (e.g. scoped refresh
   reports 1 project) after touching a file in project A only.
2. `TestWatchFallsBackToFullScanOnNewTopLevelDir`: create a brand-new project
   dir during watch → next refresh is a full scan and tracks it.
3. `TestWatchPeriodicFullScan`: shorten the package-level interval vars; after
   N scoped refreshes assert a full one occurred (observable via a manifest
   change that only a full scan would pick up, e.g. a project added on disk
   while events were suppressed... if this is awkward, assert via the counter
   through `OnRefresh` metadata — add a `FullScan bool` field to
   `WatchRefresh` to make this testable and printable).

Add `FullScan bool` to `WatchRefresh` and set it accordingly (also print it in
`printWatchRefresh` — one extra word, e.g. "(full)" vs "(scoped)").

**Verify**: `go test ./internal/devspace -run Watch -race -v` → PASS

### Step 5: Full gate

**Verify**: `make verify` → exit 0

## Test plan

Step 4's three tests plus existing watch tests unchanged. All watch tests run
with `-race`.

## Done criteria

- [ ] `make verify` exits 0; watch tests pass with `-race`
- [ ] `ScanWorkspace`'s body is byte-identical to `595d158` except imports (`git diff 595d158 -- internal/devspace/workspace.go` shows only additions)
- [ ] `WatchRefresh` has `FullScan`; scoped refreshes report it false
- [ ] `plans/README.md` status row updated

## STOP conditions

- Scoped refresh would need to duplicate `ScanWorkspace`'s reconciliation
  (removed projects) to be correct for some case you find — that's the design
  breaking; report the case.
- Watch tests turn flaky (timing) — report rather than papering over with
  sleeps; the package-level interval vars exist to make timing deterministic.
- Plan 010 has not landed and its descent change conflicts — STOP and note the
  ordering violation.

## Maintenance notes

- The safety-valve full scan is the correctness backstop; anyone tuning the
  constants should understand scoped refreshes never reconcile removals.
- Future: if `gitInfo`'s ~8 subprocesses per repo get consolidated (e.g.
  `git status --porcelain=v2 --branch` covers several), scoped refresh gets
  proportionally cheaper — separate plan.
