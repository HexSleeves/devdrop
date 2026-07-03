# Plan 014 (SPIKE): Decide and design enforcement for the Users/Teams/Access model

> **Executor instructions**: This is a design/spike plan — the deliverable is a
> **decision document**, optionally plus a warning-only prototype. Do not add
> refusing/blocking enforcement anywhere. On any STOP condition, stop and
> report. Update this plan's status row in `plans/README.md` when done.
>
> **Drift check (run first)**: `git diff --stat 595d158..HEAD -- internal/devspace/types.go internal/devspace/manifest.go internal/devspace/secrets.go internal/devspace/hosted_sync.go`

## Status

- **Priority**: P3
- **Effort**: M (spike)
- **Risk**: LOW (doc + at most warning-only code)
- **Depends on**: none
- **Category**: direction
- **Planned at**: commit `595d158`, 2026-07-02

## Why this matters

The manifest carries a full identity/authorization vocabulary — `User`, `Team`,
`TeamMember`, `ProjectAccess` with roles `owner/maintainer/developer/viewer`
(`types.go:80-124`) — and `validateManifest` enforces its referential integrity
(`manifest.go:77-129`). But **no code path ever reads a role to gate
anything**: role constants are only *written* (`secrets.go:482,495,503,509` —
invites auto-assign `developer`, the local user gets `owner`). A `viewer` can
run `env set`, `workspace push`, everything. The hosted server ignores the
model entirely: one global bearer token grants read/write to **every**
workspace ID (`hosted_sync.go` — single `Token`, auth is a constant-time
compare against `"Bearer " + s.token`). The risk today is *presentational*:
users who see roles appear after `env recipient invite --team` will assume
enforcement exists. The docs are honest about this (`docs/capstone/case-study.md`
names managed team identity as frontier work; README Roadmap lists it), so the
spike's job is to decide what the roles *mean*, not to assume they must become
a full authz system.

## Current state (read before designing)

- `internal/devspace/types.go:17-21` — role constants; `:80-124` — the model.
- `internal/devspace/manifest.go:77-129` — referential validation (the only
  consumer of the model besides secrets writes).
- `internal/devspace/secrets.go:465-515` (region) — invite/revoke populate
  `Users`/`Access`; the local identity plumbing (`localSecretRecipient`,
  `EnvRecipientExport` at `:76-82`) shows how "who am I" is derivable
  client-side (the machine's age recipient ↔ `User.AgeRecipient`).
- `internal/devspace/hosted_sync.go:437-470` — server options: one `Token`,
  one `StoreDir`. Real per-user auth server-side means token→user mapping —
  a much bigger lift than client-side checks.
- Real enforcement boundary note for the doc: anything client-side is
  *advisory* (the manifest is user-editable and the CLI is open source);
  cryptographic enforcement already exists for secrets (age recipients — a
  viewer not in the recipient list literally cannot decrypt). The design
  should lean on that existing boundary rather than duplicate it.

## Deliverables

1. **Decision doc** at `docs/access-roles.md` (match `docs/fuse-lazy-mount.md`
   structure) containing:
   - An inventory table: every mutating CLI surface (from `commands.go`) ×
     what each role *should* be allowed to do — filled with a recommendation.
   - The three candidate postures, with a recommendation and rationale:
     (a) **Document-as-advisory**: roles are bookkeeping for humans; docs say
     so explicitly; cheapest, honest.
     (b) **Client-side advisory enforcement**: CLI warns (or refuses with
     `--override`) when the local user's effective role is insufficient —
     e.g. `viewer` running `env set`. Identity = match local age recipient to
     `Manifest.Users`.
     (c) **Server-side enforcement**: per-user tokens on the hosted server
     mapped to `User.ID`. Biggest lift; requires token issuance/rotation
     design; likely post-capstone.
   - Effective-role resolution rules (direct `Access.UserID` vs via
     `TeamID` membership; revoked entries; unknown user default →
     permissive or restrictive? recommend permissive-with-warning for
     backward compat with single-user workflows that never populated Users).
   - Open questions for the maintainer, each with a recommendation.
2. **Optional prototype** (only if posture (b) is your recommendation):
   `effectiveRole(m Manifest, ageRecipient, projectID string) (string, bool)`
   in a new `internal/devspace/access.go` + table-driven tests. Warning-only,
   unwired into commands.

## Commands you will need

| Purpose | Command | Expected on success |
|---------|---------|---------------------|
| Full gate | `make verify` | exit 0 |
| Targeted | `go test ./internal/devspace -run EffectiveRole -v` | PASS (if prototype built) |

## Scope

**In scope**: `docs/access-roles.md` (create); optionally
`internal/devspace/access.go` + `access_test.go` (create, unwired).

**Out of scope**: any refusal/blocking behavior in commands; any hosted-server
auth change; any manifest schema change.

## Git workflow

- Branch: `advisor/014-spike-access-roles`
- Conventional commit, e.g. `docs: decide access-role posture + effective-role prototype`

## Done criteria

- [ ] `docs/access-roles.md` exists with the inventory table, three postures, one recommendation, and open questions
- [ ] If prototype built: `effectiveRole` unwired, tests pass
- [ ] `make verify` exits 0
- [ ] `plans/README.md` status row updated

## STOP conditions

- You start wanting to add a refusal to a command "because it's obviously
  right" — that's the feature, not the spike. STOP at warning-only.
- Effective-role resolution turns out ambiguous in the current schema (e.g.
  a user in two teams with different roles on the same project) — record it
  as an open question with a recommendation; do not extend the schema.

## Maintenance notes

- Whichever posture is chosen, the README/docs must stop implying enforcement
  exists where it doesn't — the doc should include the exact wording fix.
- If posture (c) is ever pursued, it interacts with Plan 013 (merge of
  `Users`/`Access` across machines) — note the dependency in the doc.
