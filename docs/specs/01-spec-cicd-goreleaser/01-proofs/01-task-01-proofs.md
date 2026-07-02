# Task 01 Proofs - CI workflow gates tests, vet, and build

## Task Summary

This task adds `.github/workflows/ci.yml`, giving DevDrop its first automated quality gate: `go test`, `go vet`, and a `-trimpath` build run on every pull request targeting `main` and every push to `main`, mirroring `make verify`.

## What This Task Proves

- The repository's full `make verify` gate (test → vet → build) now runs automatically in GitHub Actions.
- The workflow uses the Go version from `go.mod` (no hardcoded toolchain) and requests only `contents: read`.
- The current codebase passes the gate.

## Evidence Summary

- Local execution of the exact three CI commands passes on the implementation commit (65 tests, vet clean, build OK).
- The CI run on the feature PR completed green (URL below).

## Artifact: Local run of the exact CI commands

**What it proves:** The three commands the workflow executes pass on the current codebase, so the gate lands green on day one.

**Why it matters:** Confirms CI failures, when they appear later, indicate real regressions rather than a broken baseline.

**Command:**

~~~bash
go test ./... && go vet ./... && go build -trimpath -o bin/devspace ./cmd/devdrop
~~~

**Result summary:** All three commands succeeded: 65 tests passed across 2 packages, `go vet` reported no issues, and the build produced `bin/devspace`.

~~~text
Go test: 65 passed in 2 packages
Go vet: No issues found
Go build: Success
~~~

## Artifact: Green CI run on the feature PR

**What it proves:** The workflow triggers on `pull_request` to `main` and all three named steps (Test, Vet, Build) execute and pass on GitHub-hosted runners.

**Why it matters:** This is the end-to-end proof the CI gate is live for every future PR.

**Run URL:** <https://github.com/HexSleeves/devdrop/actions/runs/28565011052> (PR #12, conclusion: success)

**Command:**

~~~bash
gh run view 28565011052 --json jobs
~~~

**Result summary:** The `verify` job succeeded with every step green — checkout, setup-go, then the three gate steps Test, Vet, and Build.

~~~text
verify
  Run actions/checkout@v7: success
  Run actions/setup-go@v6: success
  Test: success
  Vet: success
  Build: success
~~~

## Reviewer Conclusion

The CI gate reproduces `make verify` exactly, runs with least privilege, and passes on the current codebase.
