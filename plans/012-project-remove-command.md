# Plan 012: Add `devspace project remove` — the missing untrack path

> **Executor instructions**: Follow this plan step by step. Run every
> verification command and confirm the expected result before moving on. On any
> STOP condition, stop and report. When done, update this plan's status row in
> `plans/README.md` — unless a reviewer told you they maintain the index.
>
> **Drift check (run first)**: `git diff --stat 595d158..HEAD -- internal/devspace/commands.go internal/devspace/workspace.go internal/devspace/manifest.go internal/devspace/devspace_test.go`
> On drift, re-verify excerpts; on mismatch, STOP.

## Status

- **Priority**: P2
- **Effort**: M
- **Risk**: LOW-MED (mutates the shared manifest; must cascade referential integrity correctly)
- **Depends on**: 009 recommended first (the new command is a mutating entry point and should be wrapped in `withAppLock` if that landed; if not, proceed without)
- **Category**: direction / feature
- **Planned at**: commit `595d158`, 2026-07-02

## Why this matters

`devspace project` has `add`, `hydrate`, `status` — and no way to untrack.
`ScanWorkspace` deliberately never deletes manifest entries (projects gone
from disk are only marked `Missing`), so a project added by mistake, renamed,
or retired stays in the manifest forever, syncs to every machine, and
`plan`/`apply` recreate its placeholder folder everywhere, indefinitely. The
only workaround is hand-editing `manifest.json`, which the rest of the CLI
treats as machine-owned. This also becomes the cleanup path for nested
duplicates created before Plan 010's fix.

## Current state

- `internal/devspace/commands.go:579-620` — `newProjectCommand` registers
  exactly `add`, `hydrate`, `status`. `add`'s shape (the pattern to mirror):

```go
cmd.AddCommand(&cobra.Command{
    Use:   "add <relative-path>",
    Short: "Track a project path",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        p, err := AddProject(args[0])
        ...
```

- `internal/devspace/workspace.go:397-404` — `findProject(m, ref)` matches by
  `ID || Name || Path`; reuse it for the ref argument.
- `internal/devspace/workspace.go:142-181` — `AddProject`: the
  load-config→load-manifest→mutate→save-manifest→save-state shape to mirror.
- `internal/devspace/manifest.go:77-129` — `validateManifest` enforces
  referential integrity: `Access` entries must reference existing project IDs
  (line 114: `if !projectIDs[access.ProjectID]`), so removal MUST cascade to
  `Access` or the resulting manifest fails its own validation.
- Secrets on disk live at `<workspace>/.devspace/secrets/<projectID>/<profile>.age`
  (`secrets.go:666`).
- Non-destructive invariant (CLAUDE.md): never delete user data. Removal
  untracks — it must NOT delete the project folder, and it should NOT delete
  secret files either (they may be the only copy); it prints where they remain.

### Behavior specification

`devspace project remove <project>`:

1. Resolve `<project>` via `findProject`; error `project %q not found` if absent.
2. Remove the `Project` from `Manifest.Projects`.
3. Cascade: remove every `ProjectAccess` whose `ProjectID` matches.
4. Remove `State.Projects[id]`.
5. Validate the resulting manifest (`ValidateManifest`) before saving; save
   manifest then state.
6. Never touch the project folder. If secret profiles exist for the ID, print:
   `Note: encrypted env profiles remain at <path>; delete them manually if no longer needed.`
7. Print `Removed project <name> (<path>) from the manifest. Files on disk were not touched.`

## Commands you will need

| Purpose | Command | Expected on success |
|---------|---------|---------------------|
| Full gate | `make verify` | exit 0 |
| Targeted | `go test ./internal/devspace -run 'RemoveProject|ProjectRemove' -v` | PASS |

## Scope

**In scope**:

- `internal/devspace/workspace.go` (new `RemoveProject(ref string) (Project, error)`)
- `internal/devspace/commands.go` (new subcommand in `newProjectCommand`)
- `internal/devspace/devspace_test.go` (tests)

**Out of scope**:

- Deleting anything on disk (folders, secrets) — spec forbids it.
- `ScanWorkspace` auto-pruning — scan stays non-destructive.
- `Users`/`Teams` entries — they are workspace-level, not project-level; do
  not cascade into them.

## Git workflow

- Branch: `advisor/012-project-remove`
- Conventional commit, e.g. `feat(project): add project remove to untrack a project`

## Steps

### Step 1: `RemoveProject` in `workspace.go`

Follow `AddProject`'s load/save shape. Implement the 7-point behavior spec.
Removal from the slice: filter, don't splice-in-place cleverness. Return the
removed `Project` for the command's output.

**Verify**: `go build ./...` → exit 0

### Step 2: Wire the subcommand

In `newProjectCommand`, mirroring `add`:

```go
Use:   "remove <project>",
Short: "Untrack a project (files on disk are not touched)",
```

If Plan 009 landed (check `grep -n withAppLock internal/devspace/commands.go`),
wrap the RunE body like its siblings.

**Verify**: `go run ./cmd/devspace project remove --help` → shows the command

### Step 3: Tests

In `devspace_test.go` (isolation pattern per `TestInitWorkspaceIsIdempotent`,
line 16):

1. `TestRemoveProjectUntracksAndCascades`: init workspace, add a project, give
   it an `Access` entry (hand-edit manifest via `LoadManifest`/`SaveManifest`
   with a valid `User` + `ProjectAccess` — see `manifest.go:77-129` for what
   validates), then `RemoveProject` by name. Assert: project gone from
   manifest; its `Access` entries gone; `User` entries intact;
   `State.Projects[id]` gone; `ValidateManifest` passes on the saved result;
   the project folder still exists on disk.
2. `TestRemoveProjectByPathAndID`: removal resolves by path ref and by ID ref.
3. `TestRemoveProjectNotFound`: unknown ref errors, manifest unchanged.
4. Rescan interaction: after removal, if the folder still exists on disk with
   its marker, `ScanWorkspace` will re-track it — assert that this is the
   case, and that removal of a folder-deleted project stays removed after
   rescan. (This documents intended semantics: `remove` untracks; keeping a
   live project untracked is not supported — that's what this test pins down.)

**Verify**: `go test ./internal/devspace -run RemoveProject -v` → PASS

### Step 4: Full gate

**Verify**: `make verify` → exit 0

## Test plan

Step 3's four tests. Note test 4 encodes a real product decision (remove ≠
permanent ignore); it belongs in the PR description.

## Done criteria

- [ ] `make verify` exits 0
- [ ] `devspace project --help` lists `remove`
- [ ] All four new tests pass
- [ ] `grep -n "os.Remove\|os.RemoveAll" internal/devspace/workspace.go` shows no NEW deletion calls added by this plan (hydrate's temp-dir cleanup at line ~355 already exists)
- [ ] `plans/README.md` status row updated

## STOP conditions

- Cascading removal makes `ValidateManifest` fail for a reason other than
  `Access` references (e.g. teams referencing projects in a way this plan
  didn't map) — report the validation error.
- The rescan-re-tracks behavior (test 4) is deemed wrong mid-implementation —
  that's a product decision; STOP and ask rather than inventing an ignore
  mechanism.

## Maintenance notes

- If users need "remove and don't re-track", that's a manifest-level ignore
  list — a separate, larger feature (deliberately not built here).
- Reviewer: check the cascade filters by `ProjectID` only and never touches
  `Users`/`Teams`.
- Follow-up candidate: `workspace prune` (bulk-remove all `Missing` projects)
  once this lands and semantics are proven.
