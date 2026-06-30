# Release Readiness

## What Was Tested

- Full Go unit/regression suite with `go test ./...`.
- Static checks with `go vet ./...`.
- CLI build with `go build -o .tmp/devspace ./cmd/devdrop`.
- Top-level command help for the `devspace` command surface.
- Local two-machine simulation using temporary directories and a local bare Git remote.
- Safety cases for path traversal, invalid JSON, non-empty destination folders, dirty repos, and missing Git.

## What Passed

- Repeated `init` preserves machine identity, existing config, existing manifest projects, and age identity.
- Workspace scan ignores dependency folders and does not recurse into nested Git repos inside a parent repo.
- Manifest path validation rejects absolute paths and `..` escapes.
- Manifest writes create `.bak` backups before replacing existing files.
- `plan` creates safe/skip actions and `plan --json` returns structured JSON.
- `apply` uses the saved plan and rejects manifest drift.
- `apply` creates only safe missing directories and skips non-empty destinations.
- Hydration works against a local bare Git repository with no network access.
- Hydration refuses missing remotes and non-empty destination folders.
- Dirty repos are detected and listed as skipped.
- Workspace paths and project paths with spaces work.
- Secret values remain encrypted at rest and masked in list output.

## What Failed

- Initial audit found a temporary compile break from partial path-hardening edits.
- The original hydrate path deleted placeholder directories before cloning.
- The original sync path recomputed actions during apply and had no saved plan hash.
- The original docs still described `workspace sync` as the primary flow.

## What Was Fixed

- Replaced direct placeholder marker deletion with empty-directory placeholders.
- Added a saved plan file with manifest hash validation.
- Added top-level `scan`, `plan`, and `apply` commands for the requested workflow.
- Kept `workspace sync` only as a deprecated compatibility alias.
- Added atomic JSON writes with backups.
- Centralized workspace-relative path validation and used it at mutating call sites.
- Improved Git clone errors with remote and next-step guidance.
- Added regression tests for the hardening requirements and local two-machine simulation.
- Updated README command examples, safety guarantees, troubleshooting, and roadmap.

## Known Limitations

- Manifest exchange between machines is still manual.
- There is no hosted sync, Git-backed manifest sync, daemon, FUSE layer, partial clone, or sparse checkout.
- Secret profiles are local-only; there is no OS keychain integration, team sharing, backup, or rotation flow.
- Dependency/setup commands are detected only as hints and are never executed.
- The source package path is still `cmd/devdrop`, but the intended binary name is `devspace`.

## Remaining Risks

- The local config/state directory is still named `.devdrop`; a future rename/migration may be useful if the product name stays `devspace`.
- Plan/apply is intentionally conservative and may require manual cleanup or explicit future flags for advanced cases.
- Git inspection still avoids mutating repos, so stale/outdated remote commit detection remains shallow.
- Encrypted `.env` generation overwrites the target `.env` only when explicitly requested via `env pull`.

## Recommended Next Feature

Git-backed manifest sync between machines.
