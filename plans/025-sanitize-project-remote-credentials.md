# Plan 025: Keep credentials out of project remotes and sync artifacts

> **Executor instructions**: Follow this plan exactly, never copy any real
> credential into source, tests, output, or plan artifacts, and stop on the
> conditions below. Update `plans/README.md` when done.
>
> **Drift check (run first)**: `git diff --stat 7b521c3..HEAD -- internal/devspace/git.go internal/devspace/manifest.go internal/devspace/workspace.go internal/devspace/devspace_test.go internal/devspace/workspace_sync_test.go internal/devspace/hosted_sync_test.go`

## Status

- **Priority**: P1
- **Effort**: M
- **Risk**: MED
- **Depends on**: none
- **Category**: security
- **Planned at**: commit `7b521c3`, 2026-07-10

## Why this matters

Git origin URLs are read verbatim into `Project.Remote`. Project remotes survive
manifest normalization, so embedded HTTPS userinfo can be written locally,
committed to the manifest repository, uploaded to hosted sync, included in saved
plan warnings, and printed by clone failures. Prevent new persistence and redact
every diagnostic. Credentials already synced must be treated as exposed and
rotated outside this code change.

## Current state

```go
// internal/devspace/git.go:55-58
remote := mustGit(ctx, path, "config", "--get", "remote.origin.url")
```

```go
// internal/devspace/workspace.go:140-143
if info.IsRepo {
    p.Type = ProjectTypeGit
    p.Remote = info.Remote
}
```

`cloneRepo` prints `remote` repeatedly on failure (`git.go:145-151`).
`BuildPlan` prints both observed and manifest remotes (`workspace.go:399-400`).
`redactRemote` and `sanitizeRemoteInText` already exist in
`workspace_sync.go:636-655`; reuse them for display. `validateProjectRemote`
is the trust boundary for manifest-supplied project URLs.

## Commands you will need

| Purpose | Command | Expected on success |
|---|---|---|
| Focused tests | `go test ./internal/devspace -run 'Test(ProjectRemote|ValidateProjectRemote|RedactRemote|WorkspacePush|Hosted)' -count=1` | pass |
| Security suite | `go test ./internal/devspace -run 'TestHardening|TestSync|TestHosted' -count=1` | pass |
| Full gate | `make verify` | exit 0 |

## Scope

**In scope**:

- `internal/devspace/git.go`
- `internal/devspace/manifest.go`
- `internal/devspace/workspace.go`
- `internal/devspace/devspace_test.go`
- `internal/devspace/workspace_sync_test.go`
- `internal/devspace/hosted_sync_test.go` only for hosted boundary coverage

**Out of scope**:

- Manifest-remote configuration; this finding concerns project remotes.
- Credential storage or helper configuration inside Git itself.
- Automatic credential rotation or remote rewriting in user repositories.
- Printing, committing, or documenting any real credential value.

## Git workflow

- Branch: `advisor/025-sanitize-project-remotes`
- Commit: `fix(security): keep credentials out of project remotes`
- Do not push or open a PR unless instructed.

## Steps

### Step 1: Add credential-shape regression tests

Use synthetic, clearly non-secret test data assembled locally in the test.
Prove:

- `validateProjectRemote` rejects HTTPS userinfo with a remediation-oriented error;
- scanning a Git repo never saves userinfo in `Project.Remote`;
- plan warnings and clone errors never contain the original userinfo;
- Git-backed and hosted normalized manifests cannot carry credentialed project remotes.

Do not write the credential-shaped input with `%v`/`%q` in failure messages.

**Verify**: focused tests fail before implementation.

### Step 2: Strip userinfo at Git discovery

Add one small URL helper in `git.go` that removes `url.User` from absolute URLs
while leaving SSH scp syntax and local paths unchanged. Apply it when populating
`GitInfo.Remote`, so scans persist a cloneable credential-free URL. Do not use
the display-only `redactRemote`, which inserts a `redacted` username.

**Verify**: scan regression test passes.

### Step 3: Reject credentialed manifest input

Extend `validateProjectRemote` to reject URL userinfo. This protects manifests
received from Git/hosted sync and hand-edited local manifests. The error must
name the project/remote problem without echoing the credential-bearing URL.
Keep supported HTTPS, SSH URL, scp-style, and local-path forms unchanged.

**Verify**: validation and sync-boundary tests pass.

### Step 4: Redact all project-remote diagnostics

Route clone stderr/free text through `sanitizeRemoteInText`, render the remote
with `redactRemote`, and redact both sides of BuildPlan mismatch warnings.
Search all user-facing render paths for raw `Project.Remote` output and fix only
diagnostics that can expose URL userinfo.

**Verify**: `rg -n 'info\.Remote, p\.Remote|Reason:.*Remote|Remote:\\n%s' internal/devspace` shows no unredacted diagnostic path; focused tests pass.

### Step 5: Run security and full gates

**Verify**: security suite and `make verify` pass.

## Test plan

Follow the table-driven remote tests at `devspace_test.go:116-163` and redaction
tests at `workspace_sync_test.go:956-968`. Include scan persistence, validation,
clone error, plan warning, Git sync, and hosted sync boundaries.

## Done criteria

- [ ] Newly scanned project remotes contain no URL userinfo.
- [ ] Credential-bearing project remotes fail manifest validation without echoing input.
- [ ] Clone and plan diagnostics are redacted.
- [ ] Git and hosted sync cannot publish a credential-bearing project remote.
- [ ] Focused, security, and full gates pass.
- [ ] No real credentials appear in the diff or test output.

## STOP conditions

- Removing userinfo changes non-URL local/scp remotes.
- Existing production manifests require an automatic migration to remain loadable.
- A fix requires changing Git credential-helper configuration.
- Any test or diagnostic emits the credential-shaped input.

## Maintenance notes

Keep two concepts separate: stored remotes must be credential-free; displayed
remotes must be redacted defensively. Review new manifest fields for the same
secret-copying risk.
