# Plan 024: Publish a pending manifest commit when `sync push` is retried

> **Executor instructions**: Follow this plan step by step. Run every
> verification command and confirm the expected result before moving on. Stop
> on any condition listed below. When done, update `plans/README.md`.
>
> **Drift check (run first)**: `git diff --stat 7b521c3..HEAD -- internal/devspace/workspace_sync.go internal/devspace/workspace_sync_test.go internal/devspace/base_manifest_test.go`

## Status

- **Priority**: P1
- **Effort**: S
- **Risk**: LOW
- **Depends on**: none
- **Category**: bug
- **Planned at**: commit `7b521c3`, 2026-07-10

## Why this matters

`PushWorkspaceManifest` commits before pushing. If the network push fails, the
cache is clean but ahead of its upstream. A retry sees no file change, records a
base snapshot, returns "unchanged," and never publishes the pending commit.
Retries must finish the interrupted operation and report that a remote change
was made.

## Current state

```go
// internal/devspace/workspace_sync.go:162-171
changed = changed || ignoreChanged
if !changed {
    recordBaseManifestAfterSync(normalized)
    return false, nil
}
if err := commitManifestRepo(repo, cfg); err != nil { ... }
if err := pushManifestRepo(repo); err != nil { ... }
```

`ensureManifestRepoNotBehind` at lines 418-434 rejects behind/diverged state
but deliberately permits an ahead-only cache. Tests use local bare remotes and
the helpers at `workspace_sync_test.go:918-953`. Base snapshots are written only
after a successful sync boundary; preserve that invariant.

## Commands you will need

| Purpose | Command | Expected on success |
|---|---|---|
| Focused tests | `go test ./internal/devspace -run 'TestWorkspacePush(RetriesPendingCommit|IdempotentWhenNothingChanged)|TestSyncRecordsBaseManifest' -count=1` | pass |
| Race check | `go test ./internal/devspace -race -run 'TestWorkspacePush(RetriesPendingCommit|IdempotentWhenNothingChanged)' -count=1` | pass |
| Full gate | `make verify` | exit 0 |

## Scope

**In scope**:

- `internal/devspace/workspace_sync.go`
- `internal/devspace/workspace_sync_test.go`
- `internal/devspace/base_manifest_test.go` only if a base-boundary assertion belongs there

**Out of scope**:

- Hosted sync.
- Pull/reconcile semantics.
- Retrying network operations automatically inside one invocation.
- Changing the public meaning of an ordinary no-op push.

## Git workflow

- Branch: `advisor/024-retry-pending-manifest-push`
- Commit: `fix(sync): publish pending manifest commit on retry`
- Do not push or open a PR unless instructed.

## Steps

### Step 1: Add the failing recovery test

Create `TestWorkspacePushRetriesPendingCommit` beside the idempotency test.
Initialize and push once, then make a second change to the workspace manifest.
Write the matching normalized manifest into the configured cache, commit it
there, and do not push, leaving the workspace and cached file equal while the
cache is exactly one commit ahead of `origin`. Call `PushWorkspaceManifest`
without another change. Assert:

- it returns `changed == true`;
- local cache and upstream have the same HEAD afterward;
- a fresh clone contains the pending manifest content;
- the base snapshot matches published content only after the retry succeeds.

This deterministic ahead-cache setup represents the state left by a failed
network push; do not add a flaky network failure harness.

**Verify**: focused test fails before the fix because upstream stays behind.

### Step 2: Distinguish no-op from pending publication

After fetch and behind/divergence checks, determine whether the cache is ahead
of its upstream. Preserve new workspace changes: write and commit them first as
today. If files are unchanged but a commit is pending, call the existing
`pushManifestRepo`, record the base only after it succeeds, and return `true`.
If neither files nor remote state changed, retain the existing `false` result.

Use the existing `aheadBehind`/`upstreamRef` helpers; do not introduce another
Git-state abstraction. Handle an initial repository with no upstream through
the existing normal commit-and-push path.

**Verify**: focused tests pass.

### Step 3: Run race and full gates

**Verify**: race command and `make verify` both pass.

## Test plan

Model setup after `TestWorkspacePushIdempotentWhenNothingChanged`. Cover the
ahead-only retry, unchanged/idempotent push, push failure not recording a base,
and normal first push.

## Done criteria

- [ ] Ahead-only cache commits are pushed on retry.
- [ ] Retry reports `changed == true` when it publishes pending work.
- [ ] True no-op pushes still report `false` and create no commit.
- [ ] Base state is recorded only after a successful push.
- [ ] Focused, race, and full gates pass.
- [ ] Only in-scope files and `plans/README.md` changed.

## STOP conditions

- The cache can be ahead for a documented reason other than an interrupted push.
- Correctness requires changing pull/reconcile behavior.
- The test cannot model the failure state with a local bare remote.
- Existing return-value callers require a meaning incompatible with this plan.

## Maintenance notes

Keep the cached Git repository as a recoverable queue: an ahead commit is
unsent work, not an idempotent no-op. Future sync changes must preserve the
"record base only after publication" boundary.
