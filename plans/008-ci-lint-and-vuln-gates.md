# Plan 008: Add lint, format, and vulnerability gates to make verify and CI

> **Executor instructions**: Follow this plan step by step. Run every
> verification command and confirm the expected result before moving on. On any
> STOP condition, stop and report. When done, update this plan's status row in
> `plans/README.md` — unless a reviewer told you they maintain the index.
>
> **Drift check (run first)**: `git diff --stat 595d158..HEAD -- Makefile .github/workflows/ci.yml`
> On drift, re-read both files fully before proceeding.

## Status

- **Priority**: P2
- **Effort**: S
- **Risk**: LOW (tooling only; first run may surface pre-existing findings to triage)
- **Depends on**: none (but landing it EARLY means later plans get linted)
- **Category**: dx
- **Planned at**: commit `595d158`, 2026-07-02

## Why this matters

CI and `make verify` run only `go test` + `go vet` + build. There is no
`errcheck`-class linting (ignored errors are a real hazard in a tool whose job
is not losing user state), no format enforcement, and no `govulncheck` — while
the repo ships a public-facing hosted server as a container image. A CVE in a
dependency or an unchecked error can land on `main` unnoticed.

## Current state

- `Makefile` (repo root) — `verify: test vet build`; targets `test`, `vet`,
  `build`, `clean`. `SHELL := /usr/bin/env bash`.
- `.github/workflows/ci.yml` — single `verify` job on `ubuntu-latest`:
  checkout@v7, setup-go@v6 with `go-version-file: go.mod`, then Test / Vet /
  Build steps.
- No `.golangci.yml`, no dependabot/renovate config, no format check anywhere.
- `gofmt -l internal cmd` returns clean at `595d158`.
- Go toolchain: `go 1.26` / `toolchain go1.26.4` (`go.mod`).

## Commands you will need

| Purpose | Command | Expected on success |
|---------|---------|---------------------|
| Full gate | `make verify` | exit 0 |
| Lint | `make lint` (created here) | exit 0 |
| Vulncheck | `make vulncheck` (created here) | exit 0 |

## Scope

**In scope**:
- `Makefile`
- `.github/workflows/ci.yml`
- `.golangci.yml` (create)

**Out of scope**:
- Fixing lint findings beyond the triage rule in Step 3.
- Adding dependabot config (tracked separately in the plans index backlog).
- `.goreleaser.yaml` and release workflows.

## Git workflow

- Branch: `advisor/008-ci-lint-vuln-gates`
- Conventional commit, e.g. `ci: add golangci-lint, gofmt check, and govulncheck gates`

## Steps

### Step 1: `.golangci.yml`

Create a minimal, low-noise config. golangci-lint v2 config format (verify the
installed major version first: `golangci-lint version`; if v1, use the v1
schema equivalent — same linter set):

```yaml
version: "2"
linters:
  default: standard   # govet, errcheck, staticcheck, ineffassign, unused
  settings:
    errcheck:
      exclude-functions:
        - fmt.Fprintf
        - fmt.Fprintln
        - fmt.Fprint
```

Rationale: `fmt.Fprint*` to command output writers is pervasive in
`commands.go` and not worth annotating.

### Step 2: Makefile targets

Add:

```make
lint:
	golangci-lint run ./...
	test -z "$$(gofmt -l cmd internal)" || (gofmt -l cmd internal && exit 1)

vulncheck:
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
```

Change `verify: test vet build` → `verify: test vet lint build`. Keep
`vulncheck` OUT of `verify` (network-dependent; CI-only), and add both to
`.PHONY`.

**Verify**: `make lint` → exit 0 after Step 3 triage; `make vulncheck` → exit 0
(or a real advisory — see STOP conditions).

### Step 3: Triage the first lint run

Run `golangci-lint run ./...`. Expected: a modest number of findings (repo is
`go vet`-clean and gofmt-clean).

Triage rule: fix mechanically-safe findings (unused variables, ineffectual
assignments, genuinely unchecked errors where the enclosing function returns
`error`); for anything requiring a judgment call (e.g. an intentionally
ignored error), add `//nolint:<linter> // <reason>` with a real reason.
**If findings exceed ~30 or any fix is behavior-changing, STOP and report the
list instead of fixing.**

**Verify**: `make lint` → exit 0

### Step 4: CI wiring

In `.github/workflows/ci.yml`, add after the Vet step:

```yaml
      - name: Lint
        uses: golangci/golangci-lint-action@v7
        with:
          version: latest

      - name: Vulncheck
        run: go run golang.org/x/vuln/cmd/govulncheck@latest ./...
```

(Pin the action major version that matches the config schema you used in
Step 1 — golangci-lint-action@v7+ runs golangci-lint v2.)

**Verify**: `make verify` → exit 0 locally. CI proof is on push; if you cannot
push, state that the CI step is unverified-by-CI in your report.

## Test plan

No Go tests; the gates are the tests. Local proof: `make verify`, `make lint`,
`make vulncheck` all exit 0.

## Done criteria

- [ ] `make verify` exits 0 and includes lint
- [ ] `.golangci.yml` exists; `golangci-lint run ./...` exits 0
- [ ] `make vulncheck` exits 0 (or reported advisory, see below)
- [ ] `ci.yml` contains Lint and Vulncheck steps
- [ ] `plans/README.md` status row updated

## STOP conditions

- `golangci-lint` is not installed and cannot be installed in the environment —
  report; don't hand-roll a substitute.
- First lint run yields > ~30 findings or requires behavior-changing fixes —
  report the finding list for human triage.
- `govulncheck` reports a reachable vulnerability — that's a *real* security
  finding: report it prominently (module, CVE/GO-ID, reachable symbol); do not
  bump dependencies under this plan.

## Maintenance notes

- Executors of later plans in this series must keep `make verify` (now with
  lint) green — that's the point.
- Follow-ups deferred: dependabot for `gomod` + the pinned ko base-image digest
  in `.goreleaser.yaml` (see plans/README.md backlog); pre-commit hooks.
