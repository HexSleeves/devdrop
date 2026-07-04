# Plan 015 (SPIKE): Prove out FUSE-capable CI, then execute the mount prototype's own backlog

> **Executor instructions**: This is a spike with a go/no-go gate. Phase A
> determines feasibility; Phase B only proceeds on GO. On any STOP condition,
> stop and report. Update this plan's status row in `plans/README.md` when done.
>
> **Drift check (run first)**: `git diff --stat 595d158..HEAD -- internal/devspace/mount.go internal/devspace/mount_test.go .github/workflows/ci.yml docs/fuse-lazy-mount.md`

## Status

- **Priority**: P3
- **Effort**: M (Phase A is S; Phase B only on GO)
- **Risk**: LOW-MED (CI-runner FUSE support is the known unknown — which is precisely why this was deferred by the maintainer)
- **Depends on**: none
- **Category**: direction
- **Planned at**: commit `595d158`, 2026-07-02

## Why this matters

The FUSE lazy-mount prototype is the most fully built of the repo's frontier
features (real loopback mounting, hydrate-on-lookup, status stub nodes in
`mount.go`, 351 lines), and its design doc already names its own backlog —
`docs/fuse-lazy-mount.md` "Follow-Up Cards": (1) an integration test job on
FUSE-capable hosts exercising real mount traversal + hydration, (2) a richer
project status view (dirty repos, missing `.env`, setup hints), (3) a
paths-vs-names virtual-directory decision, (4) unmount diagnostics. Today
`mount_test.go` (~100 lines) covers only the non-FUSE preview path
(`BuildMountEntries`/`PrintMountPreview`); the actual `fs.Mount` code path has
zero automated coverage. The single blocker for card (1) is whether GitHub's
hosted runners can mount FUSE — that's cheap to answer definitively, and the
answer unlocks (or honestly retires) the rest.

## Current state

- `internal/devspace/mount.go` — prototype; mount entry building is separate
  from the FUSE wiring (which is what makes preview testable without FUSE).
  Guarded so normal CLI flows never require FUSE (CLAUDE.md invariant — keep).
- `internal/devspace/mount_test.go` — preview-path tests only.
- `.github/workflows/ci.yml` — one `verify` job on `ubuntu-latest`.
- `docs/fuse-lazy-mount.md` — the backlog source; also records the library
  decision (`hanwen/go-fuse/v2`) and platform requirements (`/dev/fuse`,
  `fusermount3`).
- Known at planning time (verify, don't trust): GitHub-hosted `ubuntu-latest`
  runners historically expose `/dev/fuse` and ship `fuse3`; macOS runners
  cannot load macFUSE unsigned kexts — macOS CI is expected NO-GO.

## Deliverables

### Phase A — feasibility (always)

1. A branch with a **temporary probe workflow** (`.github/workflows/fuse-probe.yml`,
   `workflow_dispatch`-triggered or PR-scoped) that on `ubuntu-latest`: checks
   `test -e /dev/fuse`, `fusermount3 --version`, then builds and runs a minimal
   Go program using `hanwen/go-fuse/v2` to mount a loopback dir in `$RUNNER_TEMP`,
   list it, and unmount. Capture output.
2. A findings section appended to `docs/fuse-lazy-mount.md` (a `## CI
   Feasibility (verified)` heading): GO or NO-GO per platform, with the probe
   run's evidence (runner image version, commands, output).
3. Delete the probe workflow before finishing (its content lives on in the
   findings + the real job if GO).

### Phase B — only on Linux GO

1. A `mount-integration` job in `ci.yml`, Linux-only, running a new
   build-tagged test file (`//go:build linux && fusetest` or an env-gated
   `t.Skip`) `mount_fuse_test.go` that: mounts a manifest-backed workspace in
   a temp dir, asserts `ls` shows manifest path segments without hydration,
   traverses into an on-demand project backed by a **local** git remote
   (pattern for local bare-repo fixtures: `workspace_sync_test.go`), asserts
   hydration happened, asserts a failed-hydration project surfaces an error
   (not an empty dir — the doc's stated behavior), and unmounts cleanly.
   The normal `verify` job must remain FUSE-free.
2. Follow-up card (2) — richer `.devspace-status` content (dirty flag, missing
   `.env`, setup hint from the manifest's `Setup`) — ONLY if the integration
   job from (4) is green and time remains; otherwise record it as the next
   card in the doc.

## Commands you will need

| Purpose | Command | Expected on success |
|---------|---------|---------------------|
| Full gate | `make verify` | exit 0 (must stay FUSE-free) |
| Local probe (Linux only) | `test -e /dev/fuse && echo yes` | `yes` on FUSE-capable host |
| Gated tests | `go test ./internal/devspace -run Fuse -tags fusetest -v` | PASS on FUSE hosts; skipped elsewhere |

## Scope

**In scope**: `.github/workflows/fuse-probe.yml` (temporary),
`.github/workflows/ci.yml` (Phase B job only), `docs/fuse-lazy-mount.md`
(findings), `internal/devspace/mount_fuse_test.go` (create, Phase B).

**Out of scope**: `mount.go` behavior changes (except what card (2) needs, and
only in Phase B's step 5); macOS CI of any kind; making `make verify` or the
default CI job require FUSE (hard invariant).

## Git workflow

- Branch: `advisor/015-spike-fuse-ci`
- Conventional commits, e.g. `ci: probe FUSE support on hosted runners`, then
  `ci: add linux mount integration job` on GO.

## Done criteria

- [ ] `docs/fuse-lazy-mount.md` has a verified GO/NO-GO section with evidence
- [ ] Probe workflow deleted from the final branch state
- [ ] On GO: `mount-integration` CI job exists and passed at least once; `mount_fuse_test.go` covers traversal-without-hydration, hydration-on-lookup, and failure propagation
- [ ] On NO-GO: findings documented; ci.yml untouched; plan status set to DONE with "NO-GO" note (that is a successful spike outcome)
- [ ] `make verify` exits 0 and still requires no FUSE
- [ ] `plans/README.md` status row updated

## STOP conditions

- You cannot trigger/observe workflow runs from your environment — report;
  Phase A requires CI evidence, not local inference.
- The probe passes but the real integration test hits go-fuse behavior that
  needs `mount.go` changes beyond card (2)'s scope — report with the failing
  output instead of patching mount internals under a spike plan.
- Any change would make the default `verify` job require FUSE — hard stop.

## Maintenance notes

- If NO-GO on hosted runners, the honest alternatives (self-hosted runner, or
  scheduled container job with `--device /dev/fuse --cap-add SYS_ADMIN`) go in
  the findings as future options — evaluated, not built.
- Card (3) (paths vs names) is a product decision — leave it in the doc's
  backlog; don't decide it inside a CI spike.
