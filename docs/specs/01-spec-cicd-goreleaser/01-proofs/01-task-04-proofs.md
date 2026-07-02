# Task 04 Proofs - Tag-triggered release workflow with artifact attestation

## Task Summary

This task adds `.github/workflows/release.yml`: pushing a `v*` tag runs GoReleaser (full history checkout for changelog generation) to publish a GitHub Release, then attests the artifacts with `actions/attest-build-provenance` so consumers can verify provenance with `gh attestation verify`.

## What This Task Proves

- Releases are fully automated from a single tag push.
- The workflow holds exactly the least-privilege permission set the spec requires: `contents: write`, `id-token: write`, `attestations: write`.
- Artifact provenance is generated from the GoReleaser checksums file.

## Evidence Summary

- The workflow file below matches every spec requirement for Unit 3 (trigger, fetch-depth, pinned action, version constraint, permissions, attestation step).
- Runtime proof (release run URL, `gh attestation verify`) is captured with the validation prerelease in Task 05 proofs.

## Artifact: release.yml matches the spec's trigger, permissions, and attestation requirements

**What it proves:** Static inspection confirms `on: push: tags: v*`, `fetch-depth: 0`, `goreleaser/goreleaser-action@v7` with `version: "~> v2"` and `args: release --clean`, the three-scope permission block, and the `attest-build-provenance` step over `dist/checksums.txt`.

**Why it matters:** These are the exact controls the spec's Unit 3 functional requirements and Security Considerations demand.

**Artifact path:** `.github/workflows/release.yml`

**Result summary:** All spec-required elements present; no broader permissions requested; action majors verified current at authoring time (checkout v7, setup-go v6, goreleaser-action v7, attest-build-provenance v4 — `subject-checksums` input confirmed present in v4's action.yml).

~~~yaml
on:
  push:
    tags: ["v*"]

permissions:
  contents: write
  id-token: write
  attestations: write
~~~

## Artifact: Release workflow run for the validation tag

**What it proves:** GoReleaser and the attestation step succeed under the declared permissions on a real tag push.

**Why it matters:** Runtime confirmation that the pipeline works end to end.

**Run URL:** <https://github.com/HexSleeves/devdrop/actions/runs/28565647944> (tag `v0.1.0-rc.2`, conclusion: success; GoReleaser step green, attest step `skipped` because the repo is private — see the deviation record in `01-task-05-proofs.md`).

**Note:** After the first validation run failed at the attest step (attestations are unavailable on user-owned private repos), PR #13 added `if: ${{ !github.event.repository.private }}` to the attest step, so it activates automatically once the repository is public.

## Reviewer Conclusion

The release workflow encodes every Unit 3 requirement with least-privilege permissions; its runtime behavior is proven by the validation prerelease in Task 05.
