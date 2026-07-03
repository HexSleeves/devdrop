# Plan 007: Bound the hosted server's per-workspace mutex map

> **Executor instructions**: Follow this plan step by step. Run every
> verification command and confirm the expected result before moving on. On any
> STOP condition, stop and report. When done, update this plan's status row in
> `plans/README.md` — unless a reviewer told you they maintain the index.
>
> **Drift check (run first)**: `git diff --stat 595d158..HEAD -- internal/devspace/hosted_sync.go internal/devspace/hosted_sync_test.go internal/devspace/hardening_test.go`
> On drift, re-verify excerpts; on mismatch, STOP.

## Status

- **Priority**: P2
- **Effort**: S
- **Risk**: LOW
- **Depends on**: none
- **Category**: bug (resource growth on a public-facing server)
- **Planned at**: commit `595d158`, 2026-07-02

## Why this matters

The hosted sync server was deliberately hardened for public exposure — its
per-IP rate-limiter map is explicitly capped and evicted
(`defaultHostedSyncMaxLimiters = 4096`, idle TTL, dedicated test). But the
per-workspace mutex map right next to it grows without bound: any
authenticated client PUTting to many distinct workspace IDs (IDs are
client-chosen, up to 120 chars, alnum/`-`/`_`/`.`) permanently retains one
`*sync.Mutex` per ID for the life of the process. Same server, same growth
class, no cap.

## Current state

- `internal/devspace/hosted_sync.go:497-508`:

```go
// workspaceMutex lazily creates and returns the mutex used to serialize
// read-check-write access to a single workspace's stored manifest.
func (s *hostedSyncServer) workspaceMutex(workspace string) *sync.Mutex {
    s.mu.Lock()
    defer s.mu.Unlock()
    m, ok := s.workspaceMus[workspace]
    if !ok {
        m = &sync.Mutex{}
        s.workspaceMus[workspace] = m
    }
    return m
}
```

- Used at `hosted_sync.go:616`: `wsMutex := s.workspaceMutex(workspace)` —
  held for the duration of a single request's read-check-write.
- The struct (`hosted_sync.go:477-491`) holds `mu sync.Mutex` +
  `workspaceMus map[string]*sync.Mutex`.
- The contrast pattern — bounded limiter map with eviction — lives at
  `hosted_sync.go:415-423` (constants) and `:516-546` (`allowRequest`,
  `evictLimitersLocked`), with test `TestHostedServerRateLimiterMapIsBounded`
  (`hosted_sync_test.go:43`).

### Chosen approach (decision made — implement this, not eviction)

Evicting mutexes is subtly wrong: you can't safely delete a mutex another
goroutine may still hold. Replace the map with **striped locks** — a fixed
array of mutexes indexed by a hash of the workspace ID:

- Bounded by construction (no eviction, no growth).
- Two different workspaces occasionally sharing a stripe merely serialize —
  correctness is unaffected (the lock only guards read-check-write of that
  workspace's files; over-locking is safe, under-locking is not).

## Commands you will need

| Purpose | Command | Expected on success |
|---------|---------|---------------------|
| Full gate | `make verify` | exit 0 |
| Server tests | `go test ./internal/devspace -run 'Hosted|Hardening' -v` | PASS |

## Scope

**In scope**:
- `internal/devspace/hosted_sync.go`
- `internal/devspace/hosted_sync_test.go`

**Out of scope**:
- The limiter map and its eviction logic — already correct; don't touch.
- `hardening_test.go` — its contract must pass **unmodified**.

## Git workflow

- Branch: `advisor/007-bound-workspace-mutexes`
- Conventional commit, e.g. `fix(hosted): replace unbounded per-workspace mutex map with striped locks`

## Steps

### Step 1: Replace the map with stripes

In `hostedSyncServer`, replace

```go
mu           sync.Mutex
workspaceMus map[string]*sync.Mutex
```

with

```go
workspaceMus [256]sync.Mutex
```

Rewrite `workspaceMutex`:

```go
// workspaceMutex returns the stripe lock that serializes read-check-write
// access to a workspace's stored manifest. Stripes are a fixed array so
// client-chosen workspace IDs cannot grow server memory; distinct workspaces
// that hash to the same stripe merely serialize, which is safe.
func (s *hostedSyncServer) workspaceMutex(workspace string) *sync.Mutex {
    h := fnv.New32a()
    _, _ = h.Write([]byte(workspace))
    return &s.workspaceMus[h.Sum32()%uint32(len(s.workspaceMus))]
}
```

Remove the `workspaceMus: map[string]*sync.Mutex{}` initializer in
`NewHostedSyncServer` (line 468) — a fixed array needs none. Remove the now
unused `s.mu` field if nothing else uses it (`grep -n "s.mu" internal/devspace/hosted_sync.go`
first; at planning time it was only used by `workspaceMutex`). Add
`hash/fnv` to imports.

**Verify**: `go vet ./...` → exit 0

### Step 2: Tests

In `hosted_sync_test.go`:

1. `TestHostedServerWorkspaceLocksAreBounded`: create the server, call
   `workspaceMutex` with 10,000 distinct IDs, and assert no per-call
   allocation growth — with stripes this reduces to asserting the function
   returns a pointer into the fixed array (e.g. same pointer for the same ID
   twice, and total distinct pointers ≤ 256 across the 10,000 IDs).
2. Keep/confirm the concurrency contract: run the existing PUT-serialization
   tests (whatever in `hosted_sync_test.go`/`hardening_test.go` exercises
   concurrent PUTs to the same workspace) and confirm they pass.

**Verify**: `go test ./internal/devspace -run 'TestHostedServerWorkspaceLocks|Hosted' -v` → PASS

### Step 3: Full gate

**Verify**: `make verify` → exit 0 (hardening suite green, unmodified)

## Test plan

Step 2. Model the new test's structure on
`TestHostedServerRateLimiterMapIsBounded` (`hosted_sync_test.go:43`).

## Done criteria

- [ ] `make verify` exits 0
- [ ] `grep -n "map\[string\]\*sync.Mutex" internal/devspace/hosted_sync.go` → no matches
- [ ] New boundedness test passes; hardening tests pass unmodified
- [ ] No files outside scope modified (`git status`)
- [ ] `plans/README.md` status row updated

## STOP conditions

- `s.mu` turns out to guard anything besides the mutex map — report before
  removing it.
- Any hardening test fails — the striping changed observable behavior; report
  which test and how.

## Maintenance notes

- If per-workspace state beyond file access ever attaches to these locks
  (e.g. in-memory caches), striping stops being appropriate — revisit then.
- Reviewer: confirm stripe count (256) vs. expected concurrent workspaces; it
  only needs to exceed realistic concurrent request diversity, not total
  workspace count.
