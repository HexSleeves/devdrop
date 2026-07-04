# Plan 004: Validate manifest-supplied git remotes before clone

> **Executor instructions**: Follow this plan step by step. Run every
> verification command and confirm the expected result before moving on. On any
> STOP condition, stop and report. When done, update this plan's status row in
> `plans/README.md` — unless a reviewer told you they maintain the index.
>
> **Drift check (run first)**: `git diff --stat 595d158..HEAD -- internal/devspace/manifest.go internal/devspace/workspace.go internal/devspace/git.go internal/devspace/devspace_test.go`
> On drift, re-verify excerpts; on mismatch, STOP.

## Status

- **Priority**: P2
- **Effort**: S
- **Risk**: LOW (normal https/ssh/local remotes unaffected)
- **Depends on**: none
- **Category**: security
- **Planned at**: commit `595d158`, 2026-07-02

## Why this matters

`Project.Remote` comes from the synced manifest — untrusted input from
teammates or the hosted server. `HydrateProject` passes it to `git clone`
with a correct `--` separator (so `-`-prefixed flag injection is already
blocked), but nothing constrains the remote's *form*: git transport-helper
syntax (`ext::…`, `fd::…`) can execute commands during clone on machines whose
git config permits those transports. Hydrate is an everyday user action; the
CLI should reject remote forms nobody legitimately uses before git ever sees
them.

## Current state

- `internal/devspace/workspace.go:322-357` — `HydrateProject` loads the
  manifest, resolves the project, and calls `cloneRepo(p.Remote, tmp)` with no
  check on `p.Remote`'s form.
- `internal/devspace/git.go:120-137` — `cloneRepo` runs
  `exec.CommandContext(ctx, "git", "clone", "--", remote, dest)`. Keep the `--`.
- `internal/devspace/manifest.go:39-76` — `ValidateManifest` checks `p.Path`,
  ID uniqueness, type, hydrate mode — never `p.Remote`.
- Legitimate remote forms in this project's domain: `https://…`, `http://…`
  (private hosts), `ssh://…`, `git://…`, scp-style `user@host:path`, and plain
  local filesystem paths (local-first tool; local remotes are a supported
  workflow — e.g. tests use local bare repos, see `workspace_sync_test.go`).

## Commands you will need

| Purpose | Command | Expected on success |
|---------|---------|---------------------|
| Full gate | `make verify` | exit 0 |
| Targeted | `go test ./internal/devspace -run 'Remote|Hydrate|ValidateManifest' -v` | PASS |

## Scope

**In scope**:

- `internal/devspace/manifest.go` (new `validateProjectRemote` + call in `ValidateManifest`)
- `internal/devspace/workspace.go` (defense-in-depth call in `HydrateProject`)
- `internal/devspace/devspace_test.go` (tests)

**Out of scope**:

- `internal/devspace/git.go` — `cloneRepo` is correct as-is; don't touch.
- `internal/devspace/workspace_sync.go` — the manifest *sync repo* remote
  (`Config.ManifestRemote`) is user-typed local config, not synced input; a
  follow-up could cover it but it is explicitly out of scope here.

## Git workflow

- Branch: `advisor/004-validate-git-remotes`
- Conventional commit, e.g. `fix(manifest): reject transport-helper git remotes from synced manifests`

## Steps

### Step 1: Add `validateProjectRemote` in `manifest.go`

```go
// validateProjectRemote rejects remote forms that can invoke git transport
// helpers or smuggle options. Manifests are synced documents, so the remote is
// untrusted input. Allowed: http(s)/ssh/git URLs, scp-style user@host:path,
// and plain local filesystem paths.
func validateProjectRemote(remote string) error {
    if remote == "" {
        return nil // optional field; absence is handled elsewhere
    }
    if strings.HasPrefix(remote, "-") {
        return fmt.Errorf("git remote must not begin with '-': %q", remote)
    }
    // Transport-helper syntax is "<helper>::<address>" (e.g. ext::, fd::).
    if idx := strings.Index(remote, "::"); idx != -1 {
        return fmt.Errorf("git remote uses unsupported transport-helper syntax: %q", remote)
    }
    if u, err := url.Parse(remote); err == nil && u.Scheme != "" && u.Host != "" {
        switch u.Scheme {
        case "http", "https", "ssh", "git":
            return nil
        default:
            return fmt.Errorf("git remote has unsupported scheme %q", u.Scheme)
        }
    }
    // scp-style (user@host:path) and local paths fall through as allowed.
    return nil
}
```

Call it inside `ValidateManifest`'s project loop:

```go
if err := validateProjectRemote(p.Remote); err != nil {
    return fmt.Errorf("project %s: %w", p.Name, err)
}
```

Note `url.Parse` on scp-style `git@github.com:org/repo.git` yields scheme
`git@github.com` only when it parses at all — the `u.Host != ""` guard makes
such strings fall through to the allowed tail. Verify this in the tests rather
than assuming.

**Verify**: `go build ./...` → exit 0

### Step 2: Defense-in-depth in `HydrateProject`

In `internal/devspace/workspace.go`, immediately before the
`cloneRepo(p.Remote, tmp)` call (line ~356), add:

```go
if err := validateProjectRemote(p.Remote); err != nil {
    return Project{}, err
}
```

(Manifests already on disk never went through the new validation; hydrate is
the action that spends the trust, so it re-checks.)

**Verify**: `go test ./internal/devspace -run Hydrate -v` → PASS

### Step 3: Tests

In `devspace_test.go`:

1. Table-test `TestValidateProjectRemote`: allowed — `https://github.com/x/y.git`,
   `ssh://git@host/x.git`, `git@github.com:org/repo.git`, `/abs/local/path`,
   empty string; rejected — a remote beginning with `-`, `ext::sh …` form,
   `fd::17` form, `foo://host/x`.
2. Extend `TestValidateManifestRejectsUnsafeProjects` with one project whose
   `Remote` is transport-helper syntax → `ValidateManifest` errors.
3. `TestHydrateRejectsUnsupportedRemote`: manifest with a git project whose
   remote is `ext::…` form; `HydrateProject` returns the validation error and
   creates nothing at the destination.

**Verify**: `go test ./internal/devspace -run 'Remote|ValidateManifest' -v` → PASS

### Step 4: Full gate

**Verify**: `make verify` → exit 0. Also confirm the existing
`workspace_sync_test.go` suite (which clones local bare repos) still passes —
local paths must remain allowed.

## Test plan

Steps 3's three tests; pattern per `TestInitWorkspaceIsIdempotent`
(`devspace_test.go:16`) for state isolation.

## Done criteria

- [ ] `make verify` exits 0
- [ ] `TestValidateProjectRemote` table covers ≥ 9 cases and passes
- [ ] `grep -n "validateProjectRemote" internal/devspace/workspace.go` shows the hydrate-time check
- [ ] No files outside scope modified (`git status`)
- [ ] `plans/README.md` status row updated

## STOP conditions

- Any existing sync/hydrate test fails because it uses a remote form the new
  validation rejects — the allowlist is then too strict for this repo's own
  supported workflows; report which form, don't loosen ad hoc.
- `url.Parse` behavior on scp-style strings differs from the note in Step 1 on
  the toolchain in use — report what you observed.

## Maintenance notes

- If a future feature adds config-level remotes (e.g. `workspace remote set`),
  route them through `validateProjectRemote` too (explicitly deferred here).
- Reviewer: confirm the `--` in `cloneRepo` was NOT removed and the new check
  is *additive*.
