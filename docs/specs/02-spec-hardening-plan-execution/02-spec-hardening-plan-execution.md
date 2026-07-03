# 02-spec-hardening-plan-execution.md

## Introduction/Overview

This feature converts the existing `plans/` hardening bundle into a governed SDD execution workflow. The goal is to make DevSpace hardening work proceed from the audit plans through task planning, implementation, proof artifacts, and validation without duplicating work already present on the `chore/hardening-pass` branch.

## Goals

- Establish one SDD source of truth for executing and reconciling the 15 existing hardening plans.
- Preserve the intended priority order while requiring branch-drift review before implementation.
- Define proof artifacts that demonstrate security, correctness, CI, and documentation changes without exposing secrets.
- Keep implementation incremental so each landed slice can be reviewed, verified, and reflected back into `plans/README.md`.
- Prevent stale plan status from causing duplicate work or conflicting changes.

## User Stories

- **As a DevSpace maintainer**, I want the hardening plans represented in SDD artifacts so that implementation can proceed in a governed, reviewable workflow.
- **As an executor**, I want clear task boundaries and proof artifacts so that I can implement one hardening slice without guessing at scope.
- **As a reviewer**, I want each completed plan reconciled against current `main` and `chore/hardening-pass` so that duplicated or stale work is caught before merge.
- **As a security-conscious operator**, I want sensitive paths, secrets handling, hosted sync, and CI gates validated with explicit evidence so that hardening claims are defensible.

## Demoable Units of Work

### Unit 1: Plan Bundle Reconciliation

**Purpose:** Establish the accurate starting state before executing any hardening plan.

**Functional Requirements:**

- The system shall identify the current SDD spec, task list, audit, proof, and validation artifacts for hardening-plan execution.
- The system shall compare `plans/README.md` statuses against current `main` and the existing `chore/hardening-pass` branch before marking any plan ready for implementation.
- The system shall classify each plan as TODO, already implemented on a branch, drifted, blocked, rejected, or ready.
- The user shall be able to inspect which plans can be executed without duplicating existing branch work.

**Proof Artifacts:**

- `Markdown: updated SDD task list` demonstrates the reconciled execution state.
- `CLI: git diff/log output summary` demonstrates branch drift and already-implemented work were checked.
- `Markdown: planning audit report` demonstrates the reconciliation approach passed SDD planning review.

### Unit 2: Priority Hardening Execution

**Purpose:** Implement the highest-priority security and correctness plans in small, reviewable slices.

**Functional Requirements:**

- The system shall execute P1 hardening tasks before lower-priority tasks unless dependency or drift checks require a different order.
- The system shall preserve each plan's STOP conditions and scope boundaries.
- The system shall update `plans/README.md` or the SDD task list after each completed slice to prevent stale status.
- The user shall receive proof for targeted tests and full verification gates for each nontrivial slice.

**Proof Artifacts:**

- `Test: targeted Go test output` demonstrates changed behavior works for the implemented slice.
- `CLI: make verify output` demonstrates the repo-level gate passes after implementation.
- `Git: diff summary` demonstrates the implementation stayed within the intended files and scope.

### Unit 3: CI, Hosted, and Sync Reconciliation

**Purpose:** Resolve hardening work that overlaps with the existing `chore/hardening-pass` branch.

**Functional Requirements:**

- The system shall reconcile plan 008 against existing CI lint, vulnerability, and Dependabot changes before adding or modifying tooling.
- The system shall reconcile hosted-client and hosted-server plans against existing hosted hardening commits before implementing overlapping behavior.
- The system shall reconcile sync identity and output-refactor commits before changing sync or CLI presentation paths.
- The user shall see explicit reasons when a plan is marked DONE, TODO, BLOCKED, or REJECTED due to branch state.

**Proof Artifacts:**

- `Git: branch comparison summary` demonstrates overlapping branch work was reviewed.
- `CI: workflow or make target output` demonstrates lint, test, vet, build, and vulnerability checks are wired as intended when in scope.
- `Markdown: status table update` demonstrates reconciled outcomes are persisted.

### Unit 4: Final Validation and Handoff

**Purpose:** Validate that hardening execution satisfies the SDD spec and leaves the repo ready for follow-on work.

**Functional Requirements:**

- The system shall validate completed implementation against this spec's goals, demoable units, and proof artifacts.
- The system shall record any incomplete plans, rejected plans, and deferred backlog items with concise rationale.
- The system shall ensure no proof artifact includes real `.env` data, hosted tokens, age identities, or generated workspace state.
- The user shall receive a final validation result with pass/fail status and remaining gaps.

**Proof Artifacts:**

- `Markdown: validation report` demonstrates each demoable unit passed or failed with evidence.
- `CLI: final make verify output` demonstrates the repo verification gate passes.
- `Git: final status output` demonstrates tracked changes are intentional and no sensitive generated state was added.

## Non-Goals (Out of Scope)

1. **Implementing all 15 plans in Phase 1**: this spec only defines the SDD source of truth; implementation belongs to later SDD phases.
2. **Replacing the existing plan files**: `plans/*.md` remain detailed executor plans unless later reconciliation explicitly retires one.
3. **Changing product direction for spikes**: spike plans may produce decisions or prototypes, but this spec does not pre-decide access roles, manifest merge policy, or FUSE CI viability.
4. **Changing release or Railway deployment behavior**: release/deploy work is out of scope unless a hardening plan explicitly touches CI verification.
5. **Committing sensitive local state**: `.env`, hosted tokens, age identities, `.devspace/`, `.devdrop/`, and generated workspace state must remain untracked.

## Design Considerations

No specific UI design requirements identified. The workflow is documentation, CLI, test, and CI driven. Any user-facing CLI output changes should follow the existing Cobra command style and remain concise.

## Repository Standards

- DevSpace is a Go CLI with entrypoint `cmd/devspace/main.go` and most code in `internal/devspace/`.
- Tests use Go's built-in `testing` package and sit beside implementation files as `*_test.go`.
- Tests that touch app state must isolate with `t.Setenv("DEVSPACE_HOME", t.TempDir())`.
- Workspace-relative paths from user input or manifests must go through `safeWorkspacePath`.
- `plan` and `apply` must remain non-destructive and preserve existing safety tagging.
- `init` and scan-like operations must be idempotent and must not rotate machine identity or age identity.
- Use standard Go formatting and run the applicable targeted tests plus `make verify` before marking implementation complete.
- Commit subjects should use Conventional Commit style, such as `fix:`, `feat:`, `test:`, `docs:`, or `ci:`.
- Preserve `DEVSPACE_HOME` and legacy `DEV_DROP_HOME` compatibility when touching path or home resolution.

## Technical Considerations

- The audit bundle was generated against `595d158`, while current `main` is newer and `chore/hardening-pass` contains related hardening commits. Every task must run the relevant drift check before implementation.
- The SDD task list should group work by demoable unit, not simply by plan number, while preserving explicit plan dependencies: 010 after 006; 011 after 009 and 010; 012 soft-depends on 009.
- Plan 001 and Plan 002 both touch `secrets.go` and should remain serial.
- Plan 008 overlaps with existing branch work in `.github/workflows/ci.yml`, `Makefile`, `.golangci.yml`, and `.github/dependabot.yml`; it must be reconciled before implementation.
- Current external guidance supports the planned direction: Go's `govulncheck` is designed to report vulnerabilities reachable from Go code; golangci-lint supports repository config with a standard linter set; GitHub Actions hardening guidance emphasizes least privilege and pinned or controlled actions.
- Network-dependent checks such as vulnerability scanning should be treated as CI or explicit verification gates, not necessarily as always-on local developer prerequisites.

## Security Considerations

- Manifest, hosted-sync, secrets, and env-profile paths are security-sensitive and require focused tests.
- Secret profile writes, `.env` generation, age identity handling, hosted tokens, and recipient listing proofs must not expose real secrets.
- Proof artifacts may include command output, test output, and sanitized diffs; they must not include real `.env` values, age private keys, hosted bearer tokens, or local workspace state.
- Hosted server behavior has an existing hardened contract in `hardening_test.go`; changes must preserve constant-time auth, rate limiting, atomic PUT behavior, and HTTPS-only expectations unless a later approved spec changes them.
- CI hardening should prefer least-privilege permissions and stable action references consistent with repository policy.

## Success Metrics

1. **SDD artifact completeness**: spec, task list, planning audit, proof artifacts, and validation report exist for hardening-plan execution.
2. **Reconciliation accuracy**: every existing plan has a clear SDD status that accounts for `main`, `plans/README.md`, and `chore/hardening-pass`.
3. **Verification coverage**: each implemented slice includes targeted proof plus `make verify` or a documented reason the full gate could not run.
4. **Security hygiene**: no committed proof artifact or implementation change contains real secrets, tokens, identities, or generated workspace state.
5. **Reviewability**: each landed implementation slice is small enough to map back to one plan or one explicit reconciliation group.

## Open Questions

1. Should the existing `chore/hardening-pass` branch be reconciled by cherry-picking selected commits, rebasing the branch, or re-implementing only the still-valid pieces from scratch? This is non-blocking for task planning because the first task can be explicit branch reconciliation.
2. Should `plans/README.md` remain the primary status table after SDD task generation, or should the SDD task list become primary and `plans/README.md` be treated as historical executor guidance? This is non-blocking if both are kept consistent during implementation.
