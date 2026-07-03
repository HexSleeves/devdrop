# Plan 013 (SPIKE): Design manifest conflict reconciliation

> **Executor instructions**: This is a design/spike plan — the deliverable is a
> **design document plus a prototype merge function with tests**, NOT a shipped
> feature. Do not wire any new CLI flags into the sync commands. On any STOP
> condition, stop and report. Update this plan's status row in
> `plans/README.md` when done.
>
> **Drift check (run first)**: `git diff --stat 595d158..HEAD -- internal/devspace/workspace_sync.go internal/devspace/hosted_sync.go internal/devspace/types.go`
> On drift, re-read the conflict-refusal paths before designing against them.

## Status

- **Priority**: P3
- **Effort**: M (spike; the eventual feature is L)
- **Risk**: LOW (spike is additive: doc + unwired prototype)
- **Depends on**: none (but read Plans 001/004 — merged manifests must still pass the tightened validation)
- **Category**: direction
- **Planned at**: commit `595d158`, 2026-07-02

## Why this matters

This is the maintainer's own #1 named gap — twice: `docs/release-readiness.md`
"Recommended Next Feature" says *"Manifest conflict reconciliation and clearer
multi-machine history"*, and README's Roadmap leads with *"Manifest conflict
resolution and force flags."* Today both sync backends detect divergence and
refuse: `PullWorkspaceManifest` errors with `"local manifest differs from
remote manifest; push or reconcile local changes before pulling"`
(`workspace_sync.go:200-202`), and the hosted path has its own
refuse-on-conflict branch. Two machines that both scanned before syncing
deadlock until a human hand-merges JSON. For the stated multi-machine persona
(consultant across clients, `docs/capstone/spec.md`), that's the workflow's
weakest joint.

## Current state (read these before designing)

- `internal/devspace/workspace_sync.go:169-215` — `PullWorkspaceManifest`:
  pulls the sync repo, validates, then refuses when
  `localHasUnpushedManifestChanges(current, previousRemote, hasPreviousRemote, remote)`.
  The refusal is a hard safety backstop with test coverage in
  `workspace_sync_test.go`.
- `internal/devspace/hosted_sync.go` — hosted push/pull compare against
  `State.HostedSyncVersion` / `State.HostedSyncManifestHash` (`types.go:57-58`)
  — i.e. **a last-synced base is already tracked** on the hosted path, and the
  git path has `previousRemote` (the sync repo's prior copy) — both backends
  therefore already have the third point needed for a three-way merge.
- Merge units are ID-keyed records: `Projects` (by `ID`), `Machines` (by `ID`),
  `Users` (by `ID`), `Teams` (by `ID`), `Access` (composite:
  `ProjectID`+`UserID`+`TeamID`), see `types.go:70-124`.
- `localizeSyncedManifest` (used in pull) rewrites machine-local fields —
  understand it before defining merge inputs.

## Deliverables

1. **Design doc** at `docs/manifest-merge.md` (same home as the existing
   design note `docs/fuse-lazy-mount.md` — match its structure: problem,
   options considered, chosen behavior, follow-up cards) covering:
   - Merge model: three-way (base = last-synced copy; ours = local; theirs =
     remote), entity-by-entity by stable ID.
   - Per-entity rules table: added-both-sides (same ID = same path hash →
     content merge; different IDs → union), modified-ours-only /
     theirs-only (take the modification), modified-both (field-level rules or
     conflict), deleted-vs-modified (KEEP — non-destructive default), and how
     `mergeProject`'s preservation semantics (Plan 006) interact.
   - The unanswerable-by-code questions, stated as open questions with a
     recommendation each: What wins when the same path has different IDs?
     Conflicting `Access` role for the same project+user? Should `--force`
     exist independently of merge (README says yes — "force flags")?
   - UX: `workspace pull --merge` / `hosted pull --merge` proposed semantics,
     what gets printed (a summary diff), and what still refuses (true
     field-level conflicts, unless `--force-theirs`/`--force-mine`).
   - Rollout: merge lands behind the flag; default behavior stays
     refuse-and-explain.
2. **Prototype**: `internal/devspace/manifest_merge.go` with
   `mergeManifests(base, ours, theirs Manifest) (Manifest, []MergeConflict, error)`
   implementing the Projects + Access rules only (Users/Teams can be follow-up),
   plus table-driven tests in `manifest_merge_test.go` covering at minimum:
   both-added-disjoint, theirs-modified, ours-modified, both-modified-conflict,
   delete-vs-modify-keeps. The function must be **unwired** — no call sites in
   sync paths.

## Commands you will need

| Purpose | Command | Expected on success |
|---------|---------|---------------------|
| Full gate | `make verify` | exit 0 (prototype + tests compile and pass) |
| Targeted | `go test ./internal/devspace -run MergeManifests -v` | PASS |

## Scope

**In scope**: `docs/manifest-merge.md` (create),
`internal/devspace/manifest_merge.go` (create),
`internal/devspace/manifest_merge_test.go` (create).

**Out of scope**: ANY edit to `workspace_sync.go`, `hosted_sync.go`,
`commands.go`, or existing tests. The refusal backstop must be untouched.

## Git workflow

- Branch: `advisor/013-spike-manifest-merge`
- Conventional commit, e.g. `docs: design manifest conflict reconciliation + unwired merge prototype`

## Done criteria

- [ ] `docs/manifest-merge.md` exists with the rules table and ≥ 3 explicitly-stated open questions, each with a recommendation
- [ ] `mergeManifests` exists, unwired (`grep -rn "mergeManifests" internal/devspace/*.go` matches only the new files)
- [ ] `go test ./internal/devspace -run MergeManifests -v` → ≥ 5 cases pass
- [ ] Merged outputs pass `ValidateManifest` in every test case
- [ ] `make verify` exits 0
- [ ] `plans/README.md` status row updated

## STOP conditions

- The base manifest turns out NOT to be reliably recoverable on one of the two
  backends (e.g. `previousRemote` is absent on first-ever sync) — document the
  gap in the doc's open questions rather than inventing a fallback merge
  without a base (two-way merges are where data loss comes from).
- Any rule you draft would delete a project/user/access entry that only one
  side still has — that violates the non-destructive invariant; the rule must
  keep-and-flag instead.

## Maintenance notes

- The follow-up implementation plan (wiring `--merge` into both backends)
  should be written AFTER the maintainer answers the doc's open questions.
- Plan 006 changed `mergeProject` semantics; the design doc must reference the
  post-006 behavior, not `595d158`'s.
