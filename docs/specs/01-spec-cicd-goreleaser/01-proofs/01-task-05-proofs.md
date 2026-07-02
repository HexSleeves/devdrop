# Task 05 Proofs - Docs updated and validation prerelease published automatically

## Task Summary

This task updated the release documentation to make the tag-driven GoReleaser flow primary, then proved the entire pipeline by cutting the validation prerelease `v0.1.0-rc.2` from a single tag push — no manual steps.

## What This Task Proves

- One `git push origin <tag>` produces a complete GitHub Release: 4 archives + `checksums.txt` + conventional-commit-grouped changelog.
- Prerelease-style tags are auto-marked as GitHub prereleases.
- Downloaded artifacts pass checksum verification and the packaged binary reports the injected release version.
- Docs (`docs/release.md`, README) now describe the automated flow, failure recovery, and consumer verification.

## Deviation Record: attestation on a private repository

The spec planned `gh attestation verify` as runtime proof. The first validation run (tag `v0.1.0-rc.1`, run [28565440468](https://github.com/HexSleeves/devdrop/actions/runs/28565440468)) failed at the attest step: **GitHub artifact attestations are not available for user-owned private repositories**, and this repo is currently private. Resolution (PR [#13](https://github.com/HexSleeves/devdrop/pull/13)): the attest step now runs only when the repository is public (`if: ${{ !github.event.repository.private }}`), so it activates automatically if/when the repo is made public. The `v0.1.0-rc.1` release and tag were deleted and superseded by `v0.1.0-rc.2`. Attestation verification therefore remains **pending repo visibility change** — a platform limitation, not a pipeline defect.

## Evidence Summary

- Release run [28565647944](https://github.com/HexSleeves/devdrop/actions/runs/28565647944): success; attest step `skipped` (private repo), all other steps green.
- Release [v0.1.0-rc.2](https://github.com/HexSleeves/devdrop/releases/tag/v0.1.0-rc.2): `isPrerelease: true`, 5 assets.
- Local download: checksum OK, binary runs, archive layout correct.

## Artifact: Published prerelease with all assets

**What it proves:** The tag-to-release path works with no manual steps and marks `-rc` tags as prereleases.

**Why it matters:** This is the core Goal of the spec — replacing the manual `docs/release.md` assembly process.

**Command:**

~~~bash
gh release view v0.1.0-rc.2 --json name,isPrerelease,url,assets
~~~

**Result summary:** Release created automatically by the workflow with exactly the expected asset set.

~~~json
{"name":"v0.1.0-rc.2","isPrerelease":true,
 "url":"https://github.com/HexSleeves/devdrop/releases/tag/v0.1.0-rc.2",
 "assets":["checksums.txt",
  "devspace_v0.1.0-rc.2_darwin_amd64.tar.gz",
  "devspace_v0.1.0-rc.2_darwin_arm64.tar.gz",
  "devspace_v0.1.0-rc.2_linux_amd64.tar.gz",
  "devspace_v0.1.0-rc.2_linux_arm64.tar.gz"]}
~~~

## Artifact: Green release workflow run under least-privilege permissions

**What it proves:** GoReleaser ran to completion on the tag push; the attest step skipped cleanly per the private-repo condition (Task 04's runtime proof).

**Why it matters:** Confirms the workflow is reliable end to end and fails visible (the rc.1 run failed loudly before the fix).

**Run URL:** <https://github.com/HexSleeves/devdrop/actions/runs/28565647944>

~~~text
Run actions/checkout@v7: success
Run actions/setup-go@v6: success
Run GoReleaser: success
Attest release artifacts: skipped
~~~

## Artifact: Consumer verification of a downloaded artifact

**What it proves:** A consumer can download an archive plus `checksums.txt`, verify integrity, and run the binary — which reports the release version.

**Why it matters:** Validates the user story "download a prebuilt devspace binary … without a Go toolchain" including integrity checking.

**Command:**

~~~bash
gh release download v0.1.0-rc.2 --pattern 'devspace_v0.1.0-rc.2_darwin_arm64.tar.gz' --pattern 'checksums.txt'
shasum -a 256 -c --ignore-missing checksums.txt
tar -xzf devspace_v0.1.0-rc.2_darwin_arm64.tar.gz && ./devspace version
~~~

**Result summary:** Checksum verification passed; the extracted binary prints `v0.1.0-rc.2`; the archive contains `devspace`, `README.md`, `RELEASE.md`.

~~~text
devspace_v0.1.0-rc.2_darwin_arm64.tar.gz: OK
v0.1.0-rc.2
README.md
RELEASE.md
devspace
~~~

## Artifact: Conventional-commit changelog on the release

**What it proves:** The changelog is generated from git history and grouped by commit type (Features / Bug Fixes), with docs/chore/test commits excluded.

**Why it matters:** Fulfills the spec's changelog requirement with zero manual release-notes work.

**Command:**

~~~bash
gh release view v0.1.0-rc.2 --json body
~~~

**Result summary:** Body contains `## Changelog` with `### Features` (17 entries spanning full history — expected for the first tag) and `### Bug Fixes` sections.

## Artifact: Documentation updated to the automated flow

**What it proves:** `docs/release.md` now leads with the automated release process (cutting a release, failure recovery, consumer verification incl. the `checksums.txt` naming and the attestation/public-repo caveat) and retains `make release` as an offline fallback; README's Release Packaging section points to automated releases.

**Why it matters:** Closes both planning-audit FLAG findings and keeps maintainer docs truthful.

**Artifact path:** `docs/release.md`, `README.md` (see PR #12/#13 diffs)

## Reviewer Conclusion

The pipeline is proven end to end: tag push → green workflow → complete prerelease with verified checksums, working binaries, and generated changelog. The single spec deviation (attestation) is a documented GitHub platform limitation on private repos, with the workflow already wired to attest automatically once the repo is public.
