# 01-tasks-cicd-goreleaser.md

Task list for `01-spec-cicd-goreleaser.md`.

## Relevant Files

| File | Why It Is Relevant |
| --- | --- |
| `.github/workflows/ci.yml` | New CI workflow: `go test`, `go vet`, build on PRs and pushes to `main`. |
| `.github/workflows/release-check.yml` | New snapshot dry-run workflow on PRs touching release config. |
| `.github/workflows/release.yml` | New tag-triggered GoReleaser release workflow with attestation. |
| `.goreleaser.yaml` | New GoReleaser v2 config: builds, archives, checksums, changelog, release. |
| `cmd/devdrop/main.go` | Existing entry point; declares `var version = "dev"` in package `main` — ldflags target `main.version` (read-only, no change expected). |
| `Makefile` | Existing local build/verify targets that CI mirrors; retained as manual fallback (read-only, no change expected). |
| `docs/release.md` | Updated: automated tag-driven release flow becomes primary; manual flow documented as fallback. |
| `README.md` | Updated: Release Packaging section points to the automated flow and prebuilt downloads. |
| `docs/specs/01-spec-cicd-goreleaser/proofs/` | New directory collecting proof artifact outputs (run URLs, command outputs). |

### Notes

- Quality gate parity: CI must run exactly what `make verify` runs (`go test ./...`, `go vet ./...`, `go build -trimpath -o bin/devspace ./cmd/devdrop`).
- Use `actions/setup-go` with `go-version-file: go.mod` — never hardcode the Go version.
- GoReleaser config must be schema `version: 2`; validate with `goreleaser check` (a local `goreleaser` binary is available at `/opt/homebrew/bin/goreleaser`).
- Archive contents follow the existing `make release` layout: `devspace` binary + `README.md` + `RELEASE.md` (copy of `docs/release.md`). There is no `LICENSE` file at the repo root (verified), so none is packaged (spec Open Question 1 resolved).
- Archive names must keep the `v` prefix to match the existing `devspace_v0.1.0_<os>_<arch>` convention (GoReleaser `{{.Version}}` strips `v`; use a name template that restores it, e.g. `{{ .ProjectName }}_v{{ .Version }}_{{ .Os }}_{{ .Arch }}`).
- Commit messages follow the repository's conventional-commit style (`feat:`, `fix:`, ...); the changelog grouping depends on it.
- Proof artifacts are committed under `docs/specs/01-spec-cicd-goreleaser/01-proofs/` as `01-task-[TT]-proofs.md` (Phase 3 protocol naming; supersedes the `proofs/*.md` filenames referenced in sub-task text) and must contain no secrets (run URLs and command outputs only).

## Requirement Traceability

| Spec requirement (unit) | Task(s) | Planned test/proof artifact |
| --- | --- | --- |
| U1: `ci.yml` triggers on PRs to `main` + pushes to `main` | 1.1, 1.4 | Green CI run URL on PR (`proofs/ci-run.md`) |
| U1: runs `go test`, `go vet`, build with `go-version-file` | 1.2, 1.4 | `gh run view` step listing (`proofs/ci-run.md`) |
| U1: failing check fails the PR | 1.2 | CI job config (fail-fast shell steps); green run implies gate wiring |
| U1: `contents: read` only | 1.1 | `ci.yml` diff |
| U2: `version: 2` config, `devspace` binary, `-trimpath`, `main.version` ldflags | 2.1, 2.4 | `goreleaser check` output (`proofs/goreleaser-local.md`) |
| U2: exactly 4 targets (linux/darwin × amd64/arm64) | 2.1, 2.4 | Snapshot `dist/` listing (`proofs/goreleaser-local.md`) |
| U2: tar.gz named `devspace_<version>_<os>_<arch>` with README + RELEASE.md | 2.2, 2.4 | Snapshot `dist/` listing + archive content listing |
| U2: checksums + conventional-commit changelog | 2.3, 5.4 | Snapshot checksums file; release notes on validation release |
| U2: `goreleaser check` passes, no deprecated fields | 2.4 | `goreleaser check` output |
| U2: `release-check.yml` snapshot dry-run on release-config PRs | 3.1, 3.2 | Green `release-check` run URL (`proofs/release-check-run.md`) |
| U3: `release.yml` on `v*` tags, `fetch-depth: 0`, goreleaser-action `~> v2` | 4.1 | `release.yml` diff + release run URL |
| U3: minimal permissions (`contents: write`, `id-token: write`, `attestations: write`) | 4.1 | `release.yml` diff |
| U3: attestation via `actions/attest-build-provenance` over checksums | 4.2, 5.4 | `gh attestation verify` output (`proofs/validation-release.md`) |
| U3: prerelease tags auto-marked prerelease | 2.3, 5.3 | `gh release view v0.1.0-rc.1` showing `isPrerelease: true` |
| Docs: `docs/release.md` updated to tag-driven flow | 5.1 | Docs diff |

## Tasks

### [~] 1.0 Add CI workflow gating tests, vet, and build

#### 1.0 Proof Artifact(s)

- CI run URL: green `ci.yml` run on the feature PR demonstrates the test/vet/build gate executes and passes on the current codebase.
- CLI: `gh run view <run-id>` output showing the test, vet, and build steps each succeeded demonstrates all three `make verify` checks run in CI.

#### 1.0 Tasks

- [x] 1.1 Create a feature branch (e.g., `feat/cicd-goreleaser`) and add `.github/workflows/ci.yml`: name `ci`, triggers `pull_request` (branches: `main`) and `push` (branches: `main`), top-level `permissions: contents: read`, one `verify` job on `ubuntu-latest`.
- [x] 1.2 Add job steps: `actions/checkout`, `actions/setup-go@v6` with `go-version-file: go.mod`, then `go test ./...`, `go vet ./...`, and `go build -trimpath -o bin/devspace ./cmd/devdrop` as separate named steps so a failure pinpoints the gate that broke.
- [x] 1.3 Run the same three commands locally to confirm the codebase is green before pushing (mirrors `make verify`).
- [~] 1.4 Push the branch, open a PR to `main`, wait for the `ci` run, and save the run URL plus `gh run view <run-id>` output to `docs/specs/01-spec-cicd-goreleaser/proofs/ci-run.md`.

### [~] 2.0 Add GoReleaser v2 configuration for multi-platform devspace releases

#### 2.0 Proof Artifact(s)

- CLI: `goreleaser check` output reporting a valid, deprecation-free config demonstrates schema correctness (config `version: 2`).
- CLI: `goreleaser release --snapshot --clean` output plus `ls dist/` showing four `devspace_*_{linux,darwin}_{amd64,arm64}.tar.gz` archives and a checksums file demonstrates the multi-platform build works end to end locally.

#### 2.0 Tasks

- [x] 2.1 Create `.goreleaser.yaml` with `version: 2`, `project_name: devspace`, and a single build: `main: ./cmd/devdrop`, `binary: devspace`, `env: [CGO_ENABLED=0]`, `goos: [linux, darwin]`, `goarch: [amd64, arm64]`, `flags: [-trimpath]`, `ldflags: [-s -w -X main.version=v{{ .Version }}]` (matches `var version` in `cmd/devdrop/main.go`).
- [x] 2.2 Add the `archives` section: `formats: [tar.gz]`, name template `{{ .ProjectName }}_v{{ .Version }}_{{ .Os }}_{{ .Arch }}` (preserves the existing `v`-prefixed naming), files `README.md` and `docs/release.md` mapped to `RELEASE.md` (matches the `make release` archive layout).
- [x] 2.3 Add `checksum` (default `checksums.txt`), `changelog` with `sort: asc` and conventional-commit groups (`^feat` → Features, `^fix` → Bug Fixes, catch-all Others) excluding `^docs`, `^chore`, `^test`, and `release: prerelease: auto` so `-rc` tags publish as prereleases.
- [x] 2.4 Run `goreleaser check` and `goreleaser release --snapshot --clean`; verify `dist/` contains exactly the four expected archives + checksums file and that an archive contains `devspace`, `README.md`, `RELEASE.md`; save outputs to `docs/specs/01-spec-cicd-goreleaser/proofs/goreleaser-local.md`.

### [~] 3.0 Add release-config snapshot dry-run workflow (release-check)

#### 3.0 Proof Artifact(s)

- CI run URL: green `release-check.yml` run on the PR that introduces `.goreleaser.yaml` demonstrates the snapshot dry-run gate fires on release-config changes and the release build succeeds in CI.

#### 3.0 Tasks

- [x] 3.1 Add `.github/workflows/release-check.yml`: trigger `pull_request` with `paths: ['.goreleaser.yaml', '.github/workflows/release*.yml']`, `permissions: contents: read`, job with `actions/checkout` (`fetch-depth: 0`), `actions/setup-go@v6` (`go-version-file: go.mod`), and `goreleaser/goreleaser-action@v7` (`version: "~> v2"`, `args: release --snapshot --clean`).
- [~] 3.2 Push to the open PR (it touches `.goreleaser.yaml`, so the path filter must fire), wait for the green `release-check` run, and save the run URL to `docs/specs/01-spec-cicd-goreleaser/proofs/release-check-run.md`.

### [~] 4.0 Add tag-triggered release workflow with artifact attestation

#### 4.0 Proof Artifact(s)

- Diff: `.github/workflows/release.yml` showing `on: push: tags: v*`, `fetch-depth: 0`, `goreleaser/goreleaser-action` with `version: "~> v2"`, and only `contents: write`, `id-token: write`, `attestations: write` permissions demonstrates the workflow matches the spec's trigger, changelog, and least-privilege requirements.
- CI run URL: the release workflow run for the validation tag (captured in 5.0) demonstrates GoReleaser and the attestation step succeed under the declared permissions.

#### 4.0 Tasks

- [x] 4.1 Add `.github/workflows/release.yml`: trigger `push` on tags `v*`; `permissions: {contents: write, id-token: write, attestations: write}`; job with `actions/checkout` (`fetch-depth: 0`), `actions/setup-go@v6` (`go-version-file: go.mod`), and `goreleaser/goreleaser-action@v7` (`version: "~> v2"`, `args: release --clean`, `env: GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}`).
- [x] 4.2 Append an attestation step after GoReleaser: `actions/attest-build-provenance@v3` with `subject-checksums: dist/checksums.txt`, so every released archive gets verifiable provenance.
- [x] 4.3 Re-run `goreleaser check` and re-confirm workflow YAML validity (`gh workflow list` after merge, or a clean `release-check` run on the PR), then commit; the full runtime proof lands with the validation tag in 5.0.

### [ ] 5.0 Update release docs and cut validation prerelease

#### 5.0 Proof Artifact(s)

- Diff: updated `docs/release.md` (tag-driven flow primary, manual flow as fallback) and README release section demonstrates documentation matches the automated process.
- Release URL: GitHub Release for `v0.1.0-rc.1`, auto-marked prerelease, with four archives + checksums attached demonstrates the tag-to-release path works with no manual steps.
- CLI: `gh attestation verify <downloaded-archive-or-checksums> --repo HexSleeves/devdrop` succeeding demonstrates provenance attestation works end to end.

#### 5.0 Tasks

- [ ] 5.1 Update `docs/release.md`: document the automated flow (push a `v*` tag → GitHub Release with archives, checksums, changelog, attestation; verify with `gh attestation verify` and `sha256sum -c`) as the primary process, and retitle the existing manual `make release` steps as the offline/manual fallback. Update the README "Release Packaging" section to match.
- [ ] 5.2 Get the PR green (ci + release-check), then merge it into `main`.
- [ ] 5.3 Tag the merge commit `v0.1.0-rc.1` and push the tag; watch the `release` workflow run to completion.
- [ ] 5.4 Verify the release: `gh release view v0.1.0-rc.1` shows prerelease=true with 4 archives + `checksums.txt`; download one archive and `checksums.txt`, run `sha256sum -c` and `gh attestation verify checksums.txt --repo HexSleeves/devdrop`; save all outputs + release URL to `docs/specs/01-spec-cicd-goreleaser/proofs/validation-release.md` (commit proofs to `main`).
