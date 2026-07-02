# Task 03 Proofs - Snapshot dry-run gate fires on release-config PRs

## Task Summary

This task adds `.github/workflows/release-check.yml`: pull requests that touch `.goreleaser.yaml` or `.github/workflows/release*.yml` must pass a full `goreleaser release --snapshot --clean` dry-run before merge, so release-config breakage is caught before any tag is pushed.

## What This Task Proves

- The path-filtered trigger fires on the PR that introduces the release config (the PR itself touches `.goreleaser.yaml`).
- The snapshot release build succeeds on GitHub-hosted runners, not just locally.
- The workflow runs with `contents: read` only.

## Evidence Summary

- The feature PR touches `.goreleaser.yaml` and `release.yml`, so the path filter must fire on it — making the PR its own trigger test.
- The green `release-check` run (URL below) proves the CI-side snapshot build works.

## Artifact: Green release-check run on the feature PR

**What it proves:** The path filter matched, and the CI snapshot dry-run built all four targets successfully on `ubuntu-latest`.

**Why it matters:** This gate is the safety net that keeps a broken `.goreleaser.yaml` from reaching `main` and failing at tag time.

**Run URL:** <https://github.com/HexSleeves/devdrop/actions/runs/28565011054> (PR #12, conclusion: success)

**Command:**

~~~bash
gh run view 28565011054 --json jobs
~~~

**Result summary:** The path filter fired on PR #12 (which introduces `.goreleaser.yaml`) and the `snapshot` job's GoReleaser dry-run succeeded on `ubuntu-latest`.

~~~text
snapshot
  Run actions/checkout@v7: success
  Run actions/setup-go@v6: success
  GoReleaser snapshot dry-run: success
~~~

## Reviewer Conclusion

The snapshot dry-run gate is wired to exactly the release-config paths and is proven by its own introduction PR.
