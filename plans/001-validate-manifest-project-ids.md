# Plan 001: Reject unsafe project IDs so a synced manifest cannot traverse the secrets directory

> **Executor instructions**: Follow this plan step by step. Run every
> verification command and confirm the expected result before moving to the
> next step. If anything in the "STOP conditions" section occurs, stop and
> report — do not improvise. When done, update the status row for this plan
> in `plans/README.md` — unless a reviewer dispatched you and told you they
> maintain the index.
>
> **Drift check (run first)**: `git diff --stat 595d158..HEAD -- internal/devspace/manifest.go internal/devspace/secrets.go internal/devspace/devspace_test.go`
> If any in-scope file changed since this plan was written, compare the
> "Current state" excerpts against the live code before proceeding; on a
> mismatch, treat it as a STOP condition.

## Status

- **Priority**: P1
- **Effort**: S
- **Risk**: LOW
- **Depends on**: none
- **Category**: security
- **Planned at**: commit `595d158`, 2026-07-02

## Why this matters

The workspace manifest (`.devspace/manifest.json`) is a *synced document*: it
arrives from teammates via `devspace workspace pull` or from the hosted server
via `devspace hosted pull`, so its contents are untrusted input. `ValidateManifest`
runs every project **path** through `safeWorkspacePath`, but a project **ID** is
only checked for non-empty/unique — and the ID is later joined raw into a
filesystem path by `secretPath`. A manifest whose project `id` contains `../`
sequences makes an everyday command like `devspace env set <name> KEY` write the
encrypted secret file outside the intended secrets directory (arbitrary file
create/overwrite, gated only by OS permissions). `findProject` matches by name
or path too, so the user never has to type the forged ID.

## Current state

Relevant files:

- `internal/devspace/manifest.go` — manifest schema validation; `ValidateManifest` (line 39), `projectID` generator (line 147)
- `internal/devspace/secrets.go` — `secretPath` (line 666) joins the ID raw; `validSecretName` (line 670) already guards the *profile* segment the same way this plan must guard the ID segment
- `internal/devspace/workspace.go` — `findProject` (line 397) resolves a user-typed name/path to a `Project` carrying the forged ID
- `internal/devspace/devspace_test.go` — has `TestValidateManifestRejectsUnsafeProjects` (starts line 60) to extend

`ValidateManifest` today checks IDs only for presence and uniqueness
(`manifest.go:51-61`):

```go
if p.ID == "" || p.Name == "" || p.Path == "" {
    return fmt.Errorf("project id, name, and path are required")
}
if _, _, err := safeWorkspacePath(m.WorkspaceRoot, p.Path); err != nil {
    return fmt.Errorf("project %s has invalid relative path %q: %w", p.Name, p.Path, err)
}
if ids[p.ID] {
    return fmt.Errorf("duplicate project id %q", p.ID)
}
```

`secretPath` joins the ID with no check (`secrets.go:666-668`):

```go
func secretPath(cfg Config, projectID, profile string) string {
    return filepath.Join(workspaceDevdrop(cfg.WorkspaceRoot), "secrets", projectID, profile+".age")
}
```

The profile segment right next to it IS validated — this is the convention to
mirror (`secrets.go:670-675`):

```go
// validSecretName rejects profile names that could escape the per-project
// secrets directory. Only a plain single path segment is allowed.
func validSecretName(name string) error {
    if name == "" || name == "." || name == ".." {
        return fmt.Errorf("invalid profile name %q", name)
    }
```

Legitimate IDs are always machine-generated (`manifest.go:147-150`):

```go
func projectID(rel string) string {
    h := sha1.Sum([]byte(filepath.ToSlash(rel)))
    return "project_" + hex.EncodeToString(h[:])[:12]
}
```

Repo conventions: single package `internal/devspace`; tests isolate state with
`t.Setenv(envHome, t.TempDir())` and a temp workspace — see
`TestInitWorkspaceIsIdempotent` at `devspace_test.go:16` for the pattern.
Errors are `fmt.Errorf` with the offending value quoted.

## Commands you will need

| Purpose | Command | Expected on success |
|---------|---------|---------------------|
| Full gate | `make verify` | exit 0 (test + vet + build) |
| One test | `go test ./internal/devspace -run TestValidateManifest -v` | PASS |
| Secrets tests | `go test ./internal/devspace -run TestEncryptedEnv -v` | PASS |

## Scope

**In scope** (the only files you should modify):

- `internal/devspace/manifest.go`
- `internal/devspace/secrets.go`
- `internal/devspace/devspace_test.go`

**Out of scope** (do NOT touch, even though they look related):

- `internal/devspace/paths.go` — `safeWorkspacePath` is correct; do not "reuse" it for IDs (IDs are single segments, not workspace-relative paths).
- `internal/devspace/hosted_sync.go` — its `validateHostedWorkspaceID` is a different concern (server workspace IDs, already validated).
- The `projectID()` generator — do not change how IDs are produced; old manifests must keep validating.

## Git workflow

- Branch: `advisor/001-validate-manifest-project-ids`
- Conventional commits, e.g. `fix(manifest): reject project ids that are not a plain path segment`
- Do NOT push or open a PR unless the operator instructed it.

## Steps

### Step 1: Add `validateProjectID` and call it from `ValidateManifest`

In `internal/devspace/manifest.go`, add:

```go
// validateProjectID rejects IDs that could escape the per-project secrets
// directory when joined into a filesystem path. IDs are normally generated by
// projectID() ("project_" + 12 hex chars), but manifests are synced documents,
// so treat the field as untrusted: only a plain single path segment is allowed.
func validateProjectID(id string) error {
    if id == "" || id == "." || id == ".." {
        return fmt.Errorf("invalid project id %q", id)
    }
    if len(id) > 64 {
        return fmt.Errorf("project id too long: %q", id)
    }
    for _, r := range id {
        if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '_' || r == '-' || r == '.' {
            continue
        }
        return fmt.Errorf("project id contains unsupported character %q", r)
    }
    if strings.Contains(id, "..") {
        return fmt.Errorf("project id is unsafe: %q", id)
    }
    return nil
}
```

In `ValidateManifest`, immediately after the existing
`if p.ID == "" || p.Name == "" || p.Path == ""` block, add:

```go
if err := validateProjectID(p.ID); err != nil {
    return fmt.Errorf("project %s has invalid id: %w", p.Name, err)
}
```

Note: charset intentionally allows more than the generated shape
(`project_<hex>`), so hand-crafted-but-safe IDs in existing manifests keep
working. The check bans separators (`/`, `\` are not in the allowed set) and
`..`, which is what makes traversal impossible.

**Verify**: `go build ./...` → exit 0

### Step 2: Make `secretPath` defense-in-depth by validating the ID at point of use

`ValidateManifest` runs on the pull paths, but a manifest already on disk (or a
future code path that skips validation) must not reach the filesystem join.
Change `secretPath` in `internal/devspace/secrets.go` to return an error:

```go
func secretPath(cfg Config, projectID, profile string) (string, error) {
    if err := validateProjectID(projectID); err != nil {
        return "", err
    }
    return filepath.Join(workspaceDevdrop(cfg.WorkspaceRoot), "secrets", projectID, profile+".age"), nil
}
```

Find all call sites with `grep -n "secretPath(" internal/devspace/*.go` (known:
`readSecretProfile`, `writeSecretProfile`; there may be one or two more) and
propagate the error at each — every caller already returns `error`.

**Verify**: `go vet ./...` → exit 0, and `go test ./internal/devspace -run TestEncryptedEnv -v` → PASS

### Step 3: Tests

In `internal/devspace/devspace_test.go`:

1. Extend `TestValidateManifestRejectsUnsafeProjects` (line 60) with cases where
   `Project.ID` is: a traversal sequence (contains `..` and a separator), a
   value containing `/`, a value containing `\`, and `.` — each must make
   `ValidateManifest` return an error mentioning the id.
2. Add `TestSecretPathRejectsUnsafeProjectID`: build a config against a temp
   workspace (`t.Setenv(envHome, t.TempDir())` pattern from
   `TestInitWorkspaceIsIdempotent`, line 16), call the env-mutation path (e.g.
   `writeSecretProfile` or `EnvSet` after hand-writing a manifest whose project
   has an unsafe ID) and assert it errors *and* no file was created outside
   `<workspace>/.devspace/secrets/` (walk the temp dirs and assert).

**Verify**: `go test ./internal/devspace -run 'TestValidateManifest|TestSecretPath' -v` → all PASS

### Step 4: Full gate

**Verify**: `make verify` → exit 0

## Test plan

Covered in Step 3. Model the state isolation after
`TestInitWorkspaceIsIdempotent` (`devspace_test.go:16-58`). New tests: 2
(one extended, one new), covering: traversal ID rejected at validation,
traversal ID rejected at `secretPath`, no out-of-tree file created.

## Done criteria

- [ ] `make verify` exits 0
- [ ] `go test ./internal/devspace -run 'TestValidateManifest|TestSecretPath' -v` passes with the new cases
- [ ] `grep -n "filepath.Join(workspaceDevdrop" internal/devspace/secrets.go` shows the join only happens after `validateProjectID`
- [ ] No files outside the in-scope list modified (`git status`)
- [ ] `plans/README.md` status row updated

## STOP conditions

Stop and report back (do not improvise) if:

- The excerpts above don't match the live code (drift).
- You find a call site of `secretPath` whose function does not return `error` —
  report it rather than swallowing the validation error.
- Any existing test fails because a *legitimate* fixture uses an ID the new
  validation rejects — that means the charset is too strict; report instead of
  loosening it silently.

## Maintenance notes

- If a future feature lets users choose project IDs, it must funnel through
  `validateProjectID`.
- Reviewer should scrutinize: every `secretPath` call site propagates the new
  error; the charset ban includes both separators by omission.
- Deferred: hosted-server-side manifest validation already calls
  `ValidateManifest` via `validateHostedManifest`, so the server inherits this
  fix automatically — verify in review, no extra change planned.
