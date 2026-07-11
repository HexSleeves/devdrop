# Plan 027: Serialize diff-cache and mount hydration mutations

> **Executor instructions**: Follow the plan in order. `withAppLock` is
> non-reentrant; stop rather than adding nested acquisition. Update
> `plans/README.md` when done.
>
> **Drift check (run first)**: `git diff --stat 7b521c3..HEAD -- internal/devspace/commands.go internal/devspace/ui_actions.go internal/devspace/mount.go internal/devspace/commands_test.go internal/devspace/ui_test.go internal/devspace/mount_test.go internal/devspace/mount_fuse_test.go`

## Status

- **Priority**: P1
- **Effort**: M
- **Risk**: MED
- **Depends on**: none
- **Category**: bug / concurrency
- **Planned at**: commit `7b521c3`, 2026-07-10

## Why this matters

`sync diff` and dashboard status appear read-only but fetch and fast-forward the
shared manifest cache. FUSE lookup can also clone a project and rewrite shared
state. These paths bypass the cross-process application lock, allowing Git-cache
collisions, duplicate clones, failed lookups, and racing state writes.

## Current state

- `commands.go:477-486` invokes `DiffWorkspaceManifest` without `withAppLock`.
- `ui_actions.go:159-190` locks config/state reads, then lines 193-203 run the
  cache-mutating diff outside the lock.
- `mount.go:322-334` checks `shouldHydrate`, then calls `HydrateProject` directly
  from a potentially concurrent FUSE lookup.
- `lock.go:19-24` states the rule: acquire only at outermost command/action
  boundaries; domain functions do not lock themselves.
- Existing command tests use `executeCommand`; FUSE integration tests are tagged
  and skip cleanly without `/dev/fuse`.

## Commands you will need

| Purpose | Command | Expected on success |
|---|---|---|
| Focused tests | `go test ./internal/devspace -run 'Test(SyncDiffUsesAppLock|DashboardSyncStatus|MountHydration)' -count=1` | pass |
| Race tests | `go test ./internal/devspace -race -run 'Test(SyncDiffUsesAppLock|DashboardSyncStatus|MountHydration)' -count=1` | pass |
| FUSE tests | `go test ./internal/devspace -tags fusetest -run 'TestFuseMountHydratesOnLookup' -count=1` | pass or documented environment skip |
| Full gate | `make verify` | exit 0 |

## Scope

**In scope**:

- `internal/devspace/commands.go`
- `internal/devspace/ui_actions.go`
- `internal/devspace/mount.go`
- `internal/devspace/commands_test.go`
- `internal/devspace/ui_test.go`
- `internal/devspace/mount_test.go`
- `internal/devspace/mount_fuse_test.go` only for tagged concurrency coverage

**Out of scope**:

- Adding locks inside `DiffWorkspaceManifest` or `HydrateProject`.
- Holding the app lock for the mount lifetime.
- Replacing the global lock with a new lock architecture.
- Parallel project update behavior.

## Git workflow

- Branch: `advisor/027-lock-hidden-mutations`
- Commit: `fix: serialize diff cache and mount hydration`
- Do not push or open a PR unless instructed.

## Steps

### Step 1: Lock CLI diff at the command boundary

Wrap the entire `DiffWorkspaceManifest` call in the `sync diff` RunE with
`withAppLock`, while rendering output after the result is captured or inside the
same callback. Match the existing push/pull command boundary and do not lock the
domain function.

Add a command test that holds the application lock from another goroutine or
process, shortens the test lock timeout using existing package variables, and
proves `sync diff` returns the standard lock error before touching its remote.

**Verify**: focused command test passes.

### Step 2: Keep dashboard diff inside its existing lock

Move `DiffWorkspaceManifest` into the `runLocked` closure in
`dashboardSyncStatusCmd`. Preserve status construction and unavailable-reason
behavior. Do not acquire a second lock around it. The cache may still avoid
repeat work; the lock protects a real miss.

Add/adjust `dashboardSyncStatusCmd` tests to prove lock contention becomes a
status-unavailable result and that ordinary status remains correct.

**Verify**: focused dashboard tests pass under `-race`.

### Step 3: Lock and re-check FUSE hydration

Extract the smallest helper needed so a project lookup can:

1. acquire `withAppLock` for one hydration attempt;
2. re-run `shouldHydrate` after acquiring the lock;
3. call `HydrateProject` only if still needed.

The second check is required: two lookups may both decide hydration is needed
before either obtains the lock. Never hold the lock for the mount lifetime.

Add an untagged helper-level concurrency test using a local bare remote. If FUSE
is available, extend the tagged lookup test to issue concurrent reads and prove
one valid checkout with no hydration-failure diagnostic.

**Verify**: focused and FUSE commands pass or the FUSE command reports only its
existing environment skip.

### Step 4: Run full verification

**Verify**: race command and `make verify` pass.

## Test plan

Follow lock timing/control patterns in `lock_test.go`, dashboard fixtures in
`ui_test.go:130-230`, and the local-remote FUSE test in
`mount_fuse_test.go:66-106`. Cover lock contention, normal diff, cache status,
same-project concurrent hydration, and post-lock recheck.

## Done criteria

- [ ] CLI diff mutates the cache only while holding the app lock.
- [ ] Dashboard cache misses run diff inside one non-reentrant lock boundary.
- [ ] Concurrent mount lookups perform at most one hydration.
- [ ] No lock is held for the mount lifetime.
- [ ] Focused, race, and full gates pass.
- [ ] Tagged FUSE test passes or skips only because FUSE is unavailable.

## STOP conditions

- Any path already holds `withAppLock` when entering the proposed boundary.
- Correctness requires making `withAppLock` reentrant.
- FUSE callbacks cannot safely wait for the existing lock timeout.
- Tests require a network remote rather than a local bare repository.

## Maintenance notes

Commands that fetch, pull, clone, rename, or write state are mutating even when
their user-facing name says "diff," "status," or "lookup." Keep locks at outer
boundaries and re-check filesystem predicates after waiting.
