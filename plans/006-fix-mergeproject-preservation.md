# Plan 006: Make rescans actually preserve user-set hydrateMode and ignore lists

> **Executor instructions**: Follow this plan step by step. Run every
> verification command and confirm the expected result before moving on. On any
> STOP condition, stop and report. When done, update this plan's status row in
> `plans/README.md` — unless a reviewer told you they maintain the index.
>
> **Drift check (run first)**: `git diff --stat 595d158..HEAD -- internal/devspace/manifest.go internal/devspace/workspace.go internal/devspace/devspace_test.go`
> On drift, re-verify excerpts; on mismatch, STOP.

## Status

- **Priority**: P1
- **Effort**: S
- **Risk**: LOW
- **Depends on**: none
- **Category**: bug
- **Planned at**: commit `595d158`, 2026-07-02

## Why this matters

`mergeProject` exists to preserve a known project's customizations across
rescans, but two of its three guards are dead code: they trigger on
empty/zero values that the scanner *never produces*. Consequence: any user- or
teammate-set `hydrateMode` (e.g. `immediate`) or custom `ignore` list is
silently reset to auto-detected defaults on the very next `devspace scan` —
and `devspace watch` rescans continuously, so the reset is near-immediate.
Only `EnvProfiles` preservation actually works today.

## Current state

- `internal/devspace/manifest.go:198-209`:

```go
func mergeProject(old, next Project) Project {
    if old.EnvProfiles != nil && next.EnvProfiles == nil {
        next.EnvProfiles = old.EnvProfiles
    }
    if len(next.Ignore) == 0 {
        next.Ignore = old.Ignore
    }
    if next.HydrateMode == "" {
        next.HydrateMode = old.HydrateMode
    }
    return next
}
```

- `internal/devspace/workspace.go:105-122` — `projectFromPath`, the ONLY
  producer of `next` (via `upsertProject`, called from `ScanWorkspace` and
  `AddProject`), always sets both fields:

```go
p := Project{
    ...
    HydrateMode: HydrateManual,
    Ignore:      append([]string{}, DefaultIgnores...),
    ...
}
if info.IsRepo {
    p.Type = ProjectTypeGit
    ...
    p.HydrateMode = HydrateOnDemand
}
```

So `len(next.Ignore) == 0` and `next.HydrateMode == ""` are never true on a
rescan.

- Valid hydrate modes (`types.go:12-15`): `immediate`, `on-demand`,
  `metadata-only`, `manual`. Scanner defaults: `manual` for local projects,
  `on-demand` for git repos.
- `upsertProject` (`manifest.go:157-169`) matches existing projects and calls
  `mergeProject(m.Projects[i], p)`.

### Required semantics (decision already made — implement exactly this)

| Field | Rule |
|-------|------|
| `EnvProfiles` | unchanged (current behavior is correct) |
| `Ignore` | **always preserve `old.Ignore`** when it is non-empty. The scanner only ever proposes `DefaultIgnores`; the old value is either that same default or a deliberate edit. |
| `HydrateMode` | preserve `old.HydrateMode`, **except**: when the project's `Type` changed between `old` and `next` AND `old.HydrateMode` equals the scanner default for `old.Type` (`manual` for local, `on-demand` for git), take `next.HydrateMode`. This keeps the existing local→git upgrade (`manual`→`on-demand`) working while never clobbering a non-default choice. |
| `Setup`, `Remote`, `DefaultBranch`, `Type`, `Name`, `Path` | keep current behavior (scanner wins) — do not add preservation for these. |

## Commands you will need

| Purpose | Command | Expected on success |
|---------|---------|---------------------|
| Full gate | `make verify` | exit 0 |
| Targeted | `go test ./internal/devspace -run 'MergeProject|Scan|Hardening' -v` | PASS |

## Scope

**In scope**:

- `internal/devspace/manifest.go` (`mergeProject` only)
- `internal/devspace/devspace_test.go` (tests)

**Out of scope**:

- `projectFromPath` / `ScanWorkspace` — do not thread "explicitly set" flags
  through the scanner; the merge-side rules above are sufficient and smaller.
- Any manifest schema change.

## Git workflow

- Branch: `advisor/006-mergeproject-preservation`
- Conventional commit, e.g. `fix(manifest): preserve user-set hydrateMode and ignore across rescans`

## Steps

### Step 1: Rewrite `mergeProject` to the semantics table

Suggested shape:

```go
func mergeProject(old, next Project) Project {
    if old.EnvProfiles != nil && next.EnvProfiles == nil {
        next.EnvProfiles = old.EnvProfiles
    }
    if len(old.Ignore) > 0 {
        next.Ignore = old.Ignore
    }
    if old.HydrateMode != "" {
        typeChanged := old.Type != next.Type
        oldWasDefault := old.HydrateMode == defaultHydrateModeForType(old.Type)
        if !(typeChanged && oldWasDefault) {
            next.HydrateMode = old.HydrateMode
        }
    }
    return next
}

func defaultHydrateModeForType(projectType string) string {
    if projectType == ProjectTypeGit {
        return HydrateOnDemand
    }
    return HydrateManual
}
```

**Verify**: `go build ./...` → exit 0

### Step 2: Unit tests for the semantics table

Add `TestMergeProjectPreservesUserOverrides` in `devspace_test.go`, direct
unit tests on `mergeProject` (it's unexported; the test file is in-package):

1. old git project with `HydrateMode: HydrateImmediate` + scanner-next
   (`on-demand`) → stays `immediate`.
2. old local project with custom `Ignore: []string{"tmp"}` + scanner-next
   (DefaultIgnores) → stays `["tmp"]`.
3. local→git transition where old mode was the local default (`manual`) →
   becomes `on-demand` (upgrade still works).
4. local→git transition where old mode was `metadata-only` → stays
   `metadata-only`.
5. `EnvProfiles` preservation still works (old non-nil, next nil).

### Step 3: Integration test through `ScanWorkspace`

Add `TestScanPreservesManualHydrateModeAndIgnore`: init a temp workspace with
one git project (pattern: `TestInitWorkspaceIsIdempotent`, `devspace_test.go:16`,
plus whatever existing scan test builds a git repo — see
`TestWorkspaceScanDetectsGitAndSetup`), run `ScanWorkspace`, load the manifest,
set the project's `HydrateMode` to `immediate` and `Ignore` to `["custom"]`,
save it, run `ScanWorkspace` again, reload, and assert both fields survived.

**Verify**: `go test ./internal/devspace -run 'MergeProject|ScanPreserves' -v` → PASS

### Step 4: Full gate — especially the hardening suite

**Verify**: `make verify` → exit 0. `hardening_test.go`'s scan/plan/apply
idempotency tests must pass unchanged.

## Test plan

Steps 2–3: six cases total, unit + integration. No existing test asserts the
buggy behavior (verified at planning time), so no test updates are expected.

## Done criteria

- [ ] `make verify` exits 0
- [ ] `TestMergeProjectPreservesUserOverrides` (5 cases) and
      `TestScanPreservesManualHydrateModeAndIgnore` pass
- [ ] No files outside scope modified (`git status`)
- [ ] `plans/README.md` status row updated

## STOP conditions

- Any existing test fails after Step 1 — a caller depends on scanner-wins
  semantics this plan didn't identify; report which test.
- You find another producer of `next` besides `projectFromPath` feeding
  `upsertProject` — the semantics table may not fit it; report.

## Maintenance notes

- If `DefaultIgnores` changes in a future release, existing projects keep
  their old stored `Ignore` (old-wins). That's intended; a migration would be
  a separate decision.
- Reviewer: check case 3 vs 4 in the unit test carefully — that pair encodes
  the whole design.
