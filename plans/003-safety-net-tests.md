# Plan 003: Add missing safety-net tests — symlink-escape containment and secrets recipient listing

> **Executor instructions**: Follow this plan step by step. Run every
> verification command and confirm the expected result before moving on. On any
> STOP condition, stop and report. When done, update this plan's status row in
> `plans/README.md` — unless a reviewer told you they maintain the index.
>
> **Drift check (run first)**: `git diff --stat 595d158..HEAD -- internal/devspace/paths.go internal/devspace/secrets.go internal/devspace/devspace_test.go`
> On drift, re-verify the excerpts below before proceeding; on mismatch, STOP.

## Status

- **Priority**: P1
- **Effort**: S
- **Risk**: LOW (test-only)
- **Depends on**: none (independent of Plans 001/002; if they landed, their new tests coexist fine)
- **Category**: tests
- **Planned at**: commit `595d158`, 2026-07-02

## Why this matters

Two safety-relevant behaviors have zero coverage:

1. `safeWorkspacePath` contains a *second* containment check built specifically
   to catch symlink-based escapes (its comment says so), but no test in the
   repo calls `os.Symlink` at all. The flagship path-safety check could be
   "simplified" away in a refactor and nothing would fail.
2. `EnvRecipients` is the function behind `devspace env recipient list` — the
   only surface where an operator audits who can decrypt a secrets profile
   before inviting/revoking. It and `EnvRecipientExport` are at 0% coverage; a
   bug showing a revoked recipient as active would ship silently.

## Current state

- `internal/devspace/paths.go:102-148` — `safeWorkspacePath(workspace, rel)`
  returns `(full, clean, error)`. Lines 130-147 implement the symlink check:

```go
// Lexical checks above cannot catch symlink-based escapes (a directory
// inside the workspace that links outside it). Resolve symlinks on the
// existing portions of both root and candidate and re-check containment.
...
if realBack == ".." || strings.HasPrefix(realBack, "../") || filepath.IsAbs(realBack) {
    return "", "", fmt.Errorf("project path escapes workspace via symlink: %s", rel)
}
```

- `internal/devspace/secrets.go:84-107` — `EnvRecipients(projectRef, profile)`
  loads the profile and returns recipients sorted by `Name` then `ID`. Each
  `SecretRecipient` has `RevokedAt string` (`types.go:171-177`), set on
  revocation.
- `internal/devspace/secrets.go:76-82` — `EnvRecipientExport()` returns the
  local machine's `SecretRecipient` (its age public key et al.).
- Existing exemplar for the secrets test setup:
  `TestEncryptedEnvProfilesCanInviteAndRevokeRecipients` in
  `internal/devspace/devspace_test.go` (around line 265) — reuse its
  setup/flow.
- State isolation convention: `t.Setenv(envHome, t.TempDir())` + temp
  workspace, see `TestInitWorkspaceIsIdempotent` (`devspace_test.go:16`).

## Commands you will need

| Purpose | Command | Expected on success |
|---------|---------|---------------------|
| Full gate | `make verify` | exit 0 |
| Targeted | `go test ./internal/devspace -run 'Symlink|Recipient' -v` | PASS |
| Coverage spot-check | `go test ./internal/devspace -coverprofile=/tmp/c.out && go tool cover -func=/tmp/c.out \| grep -E 'EnvRecipients\|EnvRecipientExport'` | both > 0% |

## Scope

**In scope**:
- `internal/devspace/devspace_test.go` (all new tests go here)

**Out of scope**:
- ANY non-test file. This plan must not change production behavior. If a test
  you write fails against current code, that is a STOP condition (you found a
  real bug — report it), not a license to change the code.

## Git workflow

- Branch: `advisor/003-safety-net-tests`
- Conventional commit, e.g. `test: cover symlink-escape containment and env recipient listing`

## Steps

### Step 1: Symlink-escape unit test

Add `TestSafeWorkspacePathRejectsSymlinkEscape`:

- Create `workspace := t.TempDir()` and `outside := t.TempDir()`.
- `os.Symlink(outside, filepath.Join(workspace, "linked"))`.
- Call `safeWorkspacePath(workspace, "linked/project")`.
- Assert it returns an error containing `escapes workspace via symlink`.
- Also assert the happy path still works: a real subdirectory resolves without
  error.

**Verify**: `go test ./internal/devspace -run TestSafeWorkspacePathRejectsSymlinkEscape -v` → PASS

### Step 2: Symlink-escape integration test

Add `TestScanRejectsSymlinkEscapeProjectPath` (or hydrate-level, whichever is
cleaner): hand-write a manifest whose project `path` traverses through a
symlinked directory that points outside the workspace, then run the code path
that consumes it (`HydrateProject` on that ref is the most direct — it calls
`safeWorkspacePath` at `workspace.go:338`) and assert the symlink error
surfaces and nothing was created outside the workspace.

**Verify**: `go test ./internal/devspace -run TestScanRejectsSymlinkEscape -v` → PASS

### Step 3: Recipient listing tests

Extend/add alongside `TestEncryptedEnvProfilesCanInviteAndRevokeRecipients`:

1. After inviting a second recipient, call `EnvRecipients` and assert: both
   recipients present, sorted by name, `RevokedAt` empty on both.
2. After revoking one, call `EnvRecipients` again and assert the revoked
   entry's `RevokedAt` is non-empty and the active one's is still empty.
3. Add `TestEnvRecipientExportReturnsLocalIdentity`: after `InitWorkspace`,
   call `EnvRecipientExport` and assert the returned `SecretRecipient` has a
   non-empty `AgeRecipient` beginning with `age1` and a non-empty `Name`/`ID`.

**Verify**: `go test ./internal/devspace -run 'Recipient' -v` → PASS

### Step 4: Full gate + coverage spot-check

**Verify**: `make verify` → exit 0, and the coverage spot-check command from
the table shows `EnvRecipients` and `EnvRecipientExport` above 0%.

## Test plan

This plan IS the test plan: 4 new/extended tests listed in Steps 1–3.

## Done criteria

- [ ] `make verify` exits 0
- [ ] `grep -c "os.Symlink" internal/devspace/devspace_test.go` ≥ 2
- [ ] Coverage spot-check: `EnvRecipients` and `EnvRecipientExport` > 0%
- [ ] Only `devspace_test.go` modified (`git status`)
- [ ] `plans/README.md` status row updated

## STOP conditions

- Any new test FAILS against unmodified production code — you likely found a
  real bug; report the failing assertion and stop (do not fix production code
  under this plan).
- `os.Symlink` fails on the CI/dev platform — report the platform instead of
  skipping silently (a `runtime.GOOS == "windows"` skip is acceptable only if
  you confirm the repo's CI is Linux, which it is at `595d158`:
  `.github/workflows/ci.yml` uses `ubuntu-latest`).

## Maintenance notes

- These tests are the regression net for `safeWorkspacePath`'s symlink branch;
  any future refactor of `paths.go` should keep them green without edits.
- If Plan 001 landed first, its unsafe-ID tests are adjacent but disjoint —
  don't merge them.
