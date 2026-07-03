# Plan 002: Route secret, .env, and age-identity writes through the atomic write helper

> **Executor instructions**: Follow this plan step by step. Run every
> verification command and confirm the expected result before moving to the
> next step. If anything in the "STOP conditions" section occurs, stop and
> report — do not improvise. When done, update the status row for this plan
> in `plans/README.md` — unless a reviewer dispatched you and told you they
> maintain the index.
>
> **Drift check (run first)**: `git diff --stat 595d158..HEAD -- internal/devspace/secrets.go internal/devspace/init.go internal/devspace/devspace_test.go`
> On any in-scope change since `595d158`, compare the "Current state" excerpts
> against the live code; on a mismatch, STOP.

## Status

- **Priority**: P1
- **Effort**: S
- **Risk**: LOW
- **Depends on**: none
- **Category**: security
- **Planned at**: commit `595d158`, 2026-07-02

## Why this matters

Every persisted JSON artifact in this repo (`config.json`, `state.json`,
`manifest.json`, plans, hosted envelopes) goes through `atomicWriteFile`
(temp file → fsync → rename), which survives crashes and — because rename
*replaces* a destination symlink instead of following it — cannot be redirected
by a symlink planted at the destination. Three writes bypass it:

1. `writeSecretProfile` writes the age-encrypted profile with `os.WriteFile` —
   a crash mid-write leaves truncated ciphertext that fails decryption for
   every recipient, with no backup to recover from.
2. `EnvPull` writes **decrypted plaintext secrets** to `<project>/.env` with
   `os.WriteFile`, which follows symlinks: a symlink committed into a cloned
   repo named `.env` silently redirects the plaintext anywhere.
3. `ensureAgeIdentity` writes the freshly generated **age private key** with
   `os.WriteFile` (the Lstat regular-file check only protects pre-existing
   files, not first creation).

## Current state

Relevant files:

- `internal/devspace/jsonio.go` — `atomicWriteFile(path, data, perm, backup)` (line 44), the helper to reuse
- `internal/devspace/secrets.go` — `EnvPull` `.env` write (line 293), `writeSecretProfile` (line 377)
- `internal/devspace/init.go` — `ensureAgeIdentity` (line 98), identity write (line 121)

`secrets.go:290-295` (inside `EnvPull`):

```go
if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
    return "", err
}
if err := os.WriteFile(full, []byte(b.String()), 0o600); err != nil {
    return "", err
}
```

`secrets.go:373-378` (end of `writeSecretProfile`):

```go
path := secretPath(cfg, normalized.ProjectID, normalized.Profile)
if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
    return err
}
return os.WriteFile(path, buf.Bytes(), 0o600)
```

(If Plan 001 has landed, `secretPath` returns `(string, error)` — adapt
mechanically.)

`init.go:112-121` (tail of `ensureAgeIdentity`):

```go
content := fmt.Sprintf("# devspace age identity\n%s\n", identity.String())
if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
    return err
}
return os.WriteFile(path, []byte(content), 0o600)
```

The helper signature (`jsonio.go:44`):

```go
func atomicWriteFile(path string, data []byte, perm os.FileMode, backup bool) error {
```

## Commands you will need

| Purpose | Command | Expected on success |
|---------|---------|---------------------|
| Full gate | `make verify` | exit 0 |
| Secrets tests | `go test ./internal/devspace -run 'TestEncryptedEnv|TestInit' -v` | PASS |

## Scope

**In scope**:
- `internal/devspace/secrets.go` (the two write sites only)
- `internal/devspace/init.go` (the identity write only)
- `internal/devspace/devspace_test.go` (new tests)

**Out of scope**:
- `internal/devspace/jsonio.go` — do not modify the helper.
- Any change to file modes (keep `0o600`) or file content/format.

## Git workflow

- Branch: `advisor/002-atomic-secret-writes`
- Conventional commit, e.g. `fix(secrets): write secret profiles, .env, and age identity atomically`
- Do NOT push or open a PR unless instructed.

## Steps

### Step 1: Replace the three `os.WriteFile` calls

Replace each with `atomicWriteFile(<path>, <data>, 0o600, false)`.

**`backup` must be `false` at all three sites — this is load-bearing:**

- For the `.age` profile: `EnvRevoke`/`EnvRotateRecipients` re-encrypt to a new
  recipient set. A `.bak` of the previous ciphertext would remain decryptable
  by the *revoked* recipient, silently defeating revocation.
- For `.env`: a `<project>/.env.bak` would be a second plaintext copy of
  secrets that `EnvPull` never cleans up.
- For `identity.txt`: a `.bak` would be a second copy of the private key.

**Verify**: `grep -n "os.WriteFile" internal/devspace/secrets.go internal/devspace/init.go` → no matches; `go build ./...` → exit 0

### Step 2: Existing behavior unchanged

**Verify**: `go test ./internal/devspace -run 'TestEncryptedEnv|TestInit' -v` → all PASS

### Step 3: New tests

In `internal/devspace/devspace_test.go` (state isolation pattern:
`t.Setenv(envHome, t.TempDir())`, see `TestInitWorkspaceIsIdempotent` line 16):

1. `TestEnvPullReplacesSymlinkedEnvFile`: set up a project, store a secret,
   create `target.txt` outside the project and `os.Symlink` it to
   `<project>/.env`, run `EnvPull`, then assert (a) `<project>/.env` is now a
   regular file (`os.Lstat` mode), (b) `target.txt`'s content is unchanged.
2. `TestSecretWritesLeaveNoBackupFiles`: after `EnvSet` + `EnvRevoke`, glob
   `<workspace>/.devspace/secrets/**` and assert no `*.bak` and no leftover
   `.*.tmp-*` files exist.

**Verify**: `go test ./internal/devspace -run 'TestEnvPullReplacesSymlink|TestSecretWritesLeaveNoBackup' -v` → PASS

### Step 4: Full gate

**Verify**: `make verify` → exit 0

## Test plan

Two new tests (Step 3): symlink-replacement semantics and no-backup invariant.
Skip the symlink test on Windows if the suite ever runs there
(`runtime.GOOS == "windows"` guard) — repo CI is `ubuntu-latest` only.

## Done criteria

- [ ] `make verify` exits 0
- [ ] `grep -rn "os.WriteFile" internal/devspace/*.go | grep -v _test` returns no matches
- [ ] Both new tests exist and pass
- [ ] No files outside the in-scope list modified (`git status`)
- [ ] `plans/README.md` status row updated

## STOP conditions

- Excerpts don't match live code (drift — check whether Plan 001 landed and
  only the `secretPath` signature differs; that specific drift is expected and
  fine).
- Any existing test starts failing after Step 1 — the swap must be
  behavior-preserving; report the failure rather than adapting the test.
- You find additional non-test `os.WriteFile` sites beyond the three listed —
  report them; do not expand scope silently.

## Maintenance notes

- Future file writes in this package should use `writeJSON`/`atomicWriteFile`;
  a reviewer seeing a new `os.WriteFile` should push back.
- The `backup=false` rationale (revocation semantics) belongs in the PR
  description — reviewers will otherwise suggest enabling backups "for safety".
- Deferred: `ensureAgeIdentity` could additionally use `O_EXCL` create
  semantics; rename-into-place already removes the symlink-follow risk, so this
  was left out to keep the change minimal.
