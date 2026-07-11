# Plan 026: Verify tagged source before publishing release artifacts

> **Executor instructions**: Follow each step and stop on the conditions below.
> Update `plans/README.md` when done.
>
> **Drift check (run first)**: `git diff --stat 7b521c3..HEAD -- .github/workflows/release.yml`

## Status

- **Priority**: P1
- **Effort**: S
- **Risk**: LOW
- **Depends on**: none
- **Category**: dx / release
- **Planned at**: commit `7b521c3`, 2026-07-10

## Why this matters

The tag workflow builds and publishes without testing the exact tagged commit.
Main-branch CI usually provides coverage, but manual prerelease tags are an
explicit supported path and can reference any commit. Run the repository's
combined Go/TUI gate before logging into registries or publishing artifacts.

## Current state

- `.github/workflows/release.yml:21-53` checks out source, configures Go, logs
  into GHCR, configures Bun, builds companions, and runs GoReleaser.
- It contains no `make verify`, `make tui-verify`, or `make ci` step.
- `Makefile:115` defines `ci: verify tui-verify`; this is the maintained combined
  gate and should be reused.
- `.github/workflows/ci.yml` pins Bun 1.3.14; release already uses the same pin.

## Commands you will need

| Purpose | Command | Expected on success |
|---|---|---|
| Combined gate | `make ci` | Go and TUI gates pass |
| Workflow check | `ruby -e 'require "yaml"; YAML.load_file(".github/workflows/release.yml", aliases: true)'` | exit 0 |
| Diff check | `git diff --check` | exit 0 |

## Scope

**In scope**:

- `.github/workflows/release.yml`

**Out of scope**:

- Release-check behavior.
- GoReleaser configuration, deployment logic, and credentials.
- New actions or validation dependencies.
- Changing tag/release-please triggers.

## Git workflow

- Branch: `advisor/026-release-tag-verification`
- Commit: `fix(release): verify tagged source before publishing`
- Do not push or open a PR unless instructed.

## Steps

### Step 1: Put verification before credential use

Move registry login after Go and Bun setup. Add a named `Verify tagged source`
step running `make ci` before any credential-bearing login, TUI release build,
GoReleaser publication, or deployment dependency can proceed. Reuse existing
toolchain setup; do not duplicate the individual gate commands.

**Verify**: inspect the workflow ordering with
`rg -n 'Verify tagged source|Log in to GitHub Container Registry|Run GoReleaser' .github/workflows/release.yml`;
verification must have the lowest line number.

### Step 2: Run local syntax and repository gates

**Verify**: YAML command, `git diff --check`, and `make ci` all pass.

## Test plan

No new source test is needed. The workflow order is the behavior: checkout and
toolchain setup → `make ci` → credential login → release build/publish. Existing
CI and release-check remain unchanged.

## Done criteria

- [ ] Tag workflow runs `make ci` against checked-out tag source.
- [ ] Verification precedes all registry login and publication steps.
- [ ] YAML parses and `make ci` passes.
- [ ] Only `.github/workflows/release.yml` and `plans/README.md` changed.

## STOP conditions

- `make ci` requires release credentials or mutates tracked source.
- The tag workflow has intentionally moved verification into a required reusable workflow.
- YAML parsing differs because of a repository-specific loader requirement.

## Maintenance notes

Keep tag verification independent of branch protection: release correctness
must follow the tag's commit, including manually created prereleases.
