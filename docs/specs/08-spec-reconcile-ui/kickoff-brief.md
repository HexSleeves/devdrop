# Kickoff Brief: Spec 08 — Reconcile UI

## Status

Input package for the `/sdd-html` workflow. Not a spec. An SDD session
should be able to start from this document alone without re-deriving scope
from the codebase, though the code files listed below remain the source of
truth for anything this brief paraphrases.

## Ordering Note

Start after spec 07 ships. Both specs 07 and 08 touch manifest semantics;
spec 07 (role warnings) is smaller and better-defined, so it should land
first to avoid the two efforts colliding on the same code paths.

## Problem

`devspace ui` (spec 05) is an interactive dashboard, but it is intentionally
safe and local-only: it shows scan/plan/apply state for the local workspace
and never surfaces sync or reconcile. Sync and reconcile (spec 06) exist as
a separate CLI surface (`devspace workspace/hosted push|pull`, `devspace
reconcile`) and are invisible in the dashboard — a user has to leave the TUI
and run separate commands to see whether their manifest is ahead, behind, or
diverged from the remote, or to resolve a conflict. Additionally,
`reconcile --force-local|--force-remote` is a single global flag: it resolves
every conflicting project the same way, with no way to decide per-project.

## Scope

- Surface sync/reconcile status in `devspace ui`: show manifest sync state
  (ahead / behind / diverged / conflicts) and expose safe actions consistent
  with the dashboard's existing safe-action set (scan, plan, apply, hydrate
  today).
- Per-project conflict resolution: let a user resolve individual project
  conflicts local-vs-remote from the dashboard, instead of being limited to
  the single global `--force-local`/`--force-remote` flag.

## Out of Scope

From `ARCHITECTURE.md`'s "Current Gaps" section (quoted):

- "Users/Teams reconciliation is record-level (whole record wins or
  conflicts), not field-level; a losing side's unrelated field changes are
  not preserved." Field-level users/teams merge is not part of this spec.
- "Machines are excluded from reconciliation entirely." Machines
  reconciliation is not part of this spec.

## Code Surface Notes

Verified directly against the code in this repo (paths and signatures
current as of this brief; re-verify before implementing, as spec 07 may land
first and touch the same files):

- `internal/devspace/ui_actions.go` establishes the dashboard command pattern:
  each command is a `func dashboardXxxCmd(...) tea.Cmd` that runs its work
  inside `runLocked(op func() error) error` where it touches shared workspace
  state (the single dashboard lock boundary, since `withAppLock` is
  non-reentrant). `dashboardScanCmd()` returns `scanLoadedMsg`;
  mutation/refresh actions such as `dashboardPlanCmd()`,
  `dashboardApplyCmd()`, `dashboardHydrateCmd(ref string)`, and
  `dashboardRefreshCmd(syncMode string)` return `actionResultMsg`; watch
  refresh paths return `watchRefreshMsg`. New sync/reconcile commands (e.g. a
  `dashboardReconcileCmd`) should follow the mutation/refresh
  lock-then-`actionResultMsg` shape.
- `internal/devspace/reconcile.go` holds the reconcile engine entry points
  this spec would wire into the dashboard:
  - `func ReconcileWorkspaceManifest(force string, apply bool) (ReconcilePlan, error)`
  - `func ReconcileHostedManifest(force string, apply bool) (ReconcilePlan, error)`
  - `func forceReconcileConflicts(...)` — the per-conflict force logic that
    a per-project resolution UI would need to call per project rather than
    globally (currently invoked with a single workspace-wide `force`
    value).
  - Plans persist via `func SaveReconcilePlan(plan ReconcilePlan) error` and
    `func LoadReconcilePlan() (ReconcilePlan, error)`.
- `internal/devspace/ui_model.go` defines `type dashboardModel struct{...}`
  and its Bubble Tea `Update`/`View` methods, which own the dashboard's
  message loop (`scanLoadedMsg`, `actionResultMsg`, `watchRefreshMsg`,
  `errMsg`) and keybindings (currently `s`/`p`/`a`/`h` for
  scan/plan/apply/hydrate). New sync/reconcile status rows and per-project
  conflict actions would extend this model and its wiring in `ui.go`.

## Suggested Verification Shape

- Unit tests around any new dashboard command following the existing
  `dashboardXxxCmd` test patterns (mock/lock behavior, `actionResultMsg`
  shape).
- Model-level tests for new `Update` cases (status rendering, per-project
  conflict selection) consistent with existing `ui_model.go` test coverage.
- `make verify` as the merge gate.

## Source of Truth

`internal/devspace/ui_actions.go`, `internal/devspace/reconcile.go`,
`internal/devspace/ui_model.go` (and `internal/devspace/ui.go` for wiring),
plus `ARCHITECTURE.md`'s "Current Gaps" section — read them directly for
anything this brief summarizes.
