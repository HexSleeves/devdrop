# Plan 009: Add cross-process locking around config/state/manifest read-modify-write

> **Executor instructions**: Follow this plan step by step. Run every
> verification command and confirm the expected result before moving on. On any
> STOP condition, stop and report. When done, update this plan's status row in
> `plans/README.md` — unless a reviewer told you they maintain the index.
>
> **Drift check (run first)**: `git diff --stat 595d158..HEAD -- internal/devspace/commands.go internal/devspace/watch.go internal/devspace/paths.go go.mod`
> On drift, re-verify excerpts; on mismatch, STOP.

## Status

- **Priority**: P2
- **Effort**: M
- **Risk**: MED (lock scope must not deadlock; enumeration of entry points must be complete)
- **Depends on**: none
- **Category**: bug
- **Planned at**: commit `595d158`, 2026-07-02

## Why this matters

Every mutating operation is a multi-file read-modify-write with no
inter-process coordination: `ScanWorkspace` does `LoadManifest → mutate →
SaveManifest` then `LoadState → mutate → SaveState`; `AddProject`, the `env`
mutations, and hosted-sync bookkeeping do the same. `atomicWriteFile` makes
each single write atomic, but not the read-then-write span. Two concurrent
invocations — the documented `devspace watch` running while the user runs
`devspace project add`, or two watch processes — silently lose whichever
process's save lands first (a just-added project, fresh sync state), with no
error to either process.

## Current state

- `internal/devspace/config.go` — `LoadConfig/SaveConfig/LoadState/SaveState`,
  all independent read/write calls.
- `internal/devspace/workspace.go:20-103` — `ScanWorkspace`: load manifest +
  state up front, save both at the end.
- `internal/devspace/watch.go:105-116` — the watch loop's `runRefresh` calls
  `RefreshWorkspaceForWatch` (which calls `ScanWorkspace` then push).
- `internal/devspace/commands.go` — cobra wiring; every command's `RunE`
  closure is the outermost entry point of one CLI invocation.
- `internal/devspace/paths.go:28-40` — `appHome()` resolves the app-home dir;
  the lock file belongs there.
- No locking primitive exists anywhere in the CLI paths (only the hosted
  server's in-process mutexes). No third-party deps beyond those in `go.mod`
  (age, fsnotify, go-fuse, cobra, x/term, x/time).

### Chosen design (implement exactly this)

- **One process-wide advisory lock file**: `<appHome>/.lock`, acquired
  exclusively for the span of any mutating operation. Single coarse lock —
  no per-file locks, no reader locks. Contention is rare (a human + a watcher);
  coarse and obviously-correct beats granular and deadlock-prone.
- **Acquire at outermost entry points only** (a cobra `RunE` / the watch
  refresh callback), never inside domain functions — the lock is
  **not reentrant**, so nested acquisition would deadlock.
- Dependency: `github.com/gofrs/flock` (mature, cross-platform flock
  wrapper). Adding one dep is acceptable here; hand-rolling
  flock/LockFileEx portability is not S-effort.
- Timeout behavior: try for 10s with backoff
  (`flock.TryLockContext(ctx, 250ms)`), then fail with
  `another devspace process holds the lock (<lockpath>); retry when it finishes`.

## Commands you will need

| Purpose | Command | Expected on success |
|---------|---------|---------------------|
| Add dep | `go get github.com/gofrs/flock@latest && go mod tidy` | go.mod/go.sum updated |
| Full gate | `make verify` | exit 0 |
| Targeted | `go test ./internal/devspace -run Lock -v` | PASS |

## Scope

**In scope**:
- `internal/devspace/lock.go` (create)
- `internal/devspace/lock_test.go` (create)
- `internal/devspace/commands.go` (wrap mutating RunE bodies)
- `internal/devspace/watch.go` (wrap `runRefresh`'s call)
- `go.mod`, `go.sum`

**Out of scope**:
- `jsonio.go`, `config.go`, domain functions in `workspace.go`/`secrets.go` —
  no lock calls inside them (reentrancy hazard).
- The hosted **server** (`hosted serve`) — it has its own in-process locking
  and its own storage dir; do not wrap it.
- Any behavior change to read-only commands (`status`, `doctor`, `plan` in
  preview, `hosted config get`, …) — they stay lock-free.

## Git workflow

- Branch: `advisor/009-cross-process-locking`
- Conventional commit, e.g. `fix: serialize concurrent devspace invocations with an app-home lock`

## Steps

### Step 1: `lock.go`

```go
package devspace

// withAppLock serializes mutating CLI operations across processes. It guards
// the read-modify-write span over config.json/state.json/manifest.json, which
// atomicWriteFile alone cannot (each write is atomic; the span is not).
// NOT reentrant: acquire only at an outermost entry point (a command RunE or
// the watch refresh), never inside domain functions.
func withAppLock(fn func() error) error {
    home, err := appHome()
    if err != nil {
        return err
    }
    if err := os.MkdirAll(home, 0o700); err != nil {
        return err
    }
    l := flock.New(filepath.Join(home, ".lock"))
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    ok, err := l.TryLockContext(ctx, 250*time.Millisecond)
    if err != nil {
        return err
    }
    if !ok {
        return fmt.Errorf("another devspace process holds the lock (%s); retry when it finishes", l.Path())
    }
    defer l.Unlock()
    return fn()
}
```

**Verify**: `go build ./...` → exit 0

### Step 2: Enumerate and wrap mutating entry points

In `commands.go`, wrap the body of each mutating `RunE` in
`withAppLock(func() error { … })`. Enumerate by reading `NewRootCommand`'s
tree; at `595d158` the mutating commands are:

- `init`
- `workspace scan`, `workspace plan` (writes `last-plan.json`), `workspace apply`,
  `workspace push`, `workspace pull`, `workspace sync`, `workspace remote set/create`
- `project add`, `project hydrate`
- `env set`, `env pull`, `env recipient invite/revoke/rotate`
- `hosted config set`, `hosted push`, `hosted pull`
- `setup` (writes state) — confirm by reading its RunE; wrap if it saves
  state/config
- `migrate`-related paths run inside `init`/first command via
  `migrateLegacyHome` — covered by wrapping the commands themselves

Read-only commands (`status`, `doctor`, `env list`, `env recipients`,
`hosted config get`, `workspace diff`, `mount --preview`, `version`, …) are
NOT wrapped.

In `watch.go`, wrap the `runRefresh` body's `RefreshWorkspaceForWatch` call
(`watch.go:105-116`) in `withAppLock` — the watch loop itself must not hold
the lock between refreshes, only during one.

**Verify**: `go build ./...` → exit 0; then `go test ./internal/devspace -v` →
all existing tests PASS (tests call domain functions directly, not RunE, so
they bypass the lock and must be unaffected).

### Step 3: Tests

`lock_test.go`:

1. `TestWithAppLockSerializesWriters`: two goroutines call `withAppLock`, each
   incrementing a shared counter after a small sleep inside the critical
   section; assert no interleaving (e.g. track max concurrent = 1 with an
   atomic).
2. `TestWithAppLockTimesOut`: hold the lock in one flock handle, call
   `withAppLock` with the 10s timeout shortened for test via… the timeout is
   a constant — instead assert `TryLockContext` path by holding the lock and
   using a context deadline: extract the timeout as a package-level `var
   appLockTimeout = 10 * time.Second` so the test can shorten it with
   `t.Cleanup` restore. Assert the error mentions "another devspace process".
3. `TestConcurrentScanAndAddProjectDoNotLoseWrites` (the money test): in one
   process but through the *locked wrappers*, run `withAppLock(ScanWorkspace)`
   and `withAppLock(AddProject…)` concurrently in a loop (~20 iterations) and
   assert the added project is present in the manifest afterwards every time.

Use `t.Setenv(envHome, t.TempDir())` isolation (pattern:
`TestInitWorkspaceIsIdempotent`, `devspace_test.go:16`).

**Verify**: `go test ./internal/devspace -run Lock -race -v` → PASS

### Step 4: Full gate + race

**Verify**: `make verify` → exit 0, and `go test ./internal/devspace -race` → PASS

## Test plan

Step 3's three tests, all run with `-race`.

## Done criteria

- [ ] `make verify` exits 0; `go test ./internal/devspace -race` exits 0
- [ ] `grep -c "withAppLock" internal/devspace/commands.go` ≥ 12 (one per mutating command)
- [ ] `grep -n "withAppLock" internal/devspace/watch.go` shows the refresh wrap
- [ ] `grep -rn "withAppLock" internal/devspace/workspace.go internal/devspace/secrets.go internal/devspace/config.go` → no matches (no nested acquisition)
- [ ] `plans/README.md` status row updated

## STOP conditions

- Any RunE body already calls another wrapped command's function (nested
  lock) — report the pair instead of restructuring.
- `gofrs/flock` fails on the target platform set (check `.goreleaser.yaml`
  builds: linux/darwin/windows) — report before substituting.
- Existing tests break because a domain function newly requires the lock —
  that means a lock call leaked below the entry-point layer; fix placement,
  and if unclear, STOP.

## Maintenance notes

- New mutating commands MUST wrap with `withAppLock`; reviewers should check
  this on any commands.go PR. Consider a follow-up lint/test that asserts the
  wrap by inspection.
- The lock is app-home-scoped: two different `DEVSPACE_HOME`s don't contend
  (correct — separate state).
- Deferred: finer-grained locks (per workspace) only if the coarse lock ever
  causes real contention; watch-mode + human is the only realistic overlap
  today.
