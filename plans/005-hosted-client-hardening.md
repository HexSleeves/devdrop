# Plan 005: Re-validate the hosted endpoint at point of use and add an env-var token path

> **Executor instructions**: Follow this plan step by step. Run every
> verification command and confirm the expected result before moving on. On any
> STOP condition, stop and report. When done, update this plan's status row in
> `plans/README.md` — unless a reviewer told you they maintain the index.
>
> **Drift check (run first)**: `git diff --stat 595d158..HEAD -- internal/devspace/hosted_sync.go internal/devspace/commands.go internal/devspace/hosted_sync_test.go`
> On drift, re-verify excerpts; on mismatch, STOP.

## Status

- **Priority**: P2
- **Effort**: S
- **Risk**: LOW
- **Depends on**: none
- **Category**: security
- **Planned at**: commit `595d158`, 2026-07-02

## Why this matters

Two point-of-use gaps in the hosted-sync client:

1. The HTTPS-only rule (plain `http` only for loopback) is enforced **only** in
   `SetHostedSync` at config-write time. `GetHostedSync` — which every
   `hosted push/pull` goes through — never re-checks the scheme. Anything that
   puts a non-loopback `http://` endpoint into `config.json` outside the CLI's
   own setter (hand edit, other tooling, a future migration bug) makes every
   subsequent sync send the bearer token and full manifest over plaintext,
   silently.
2. `devspace hosted config set --token <t>` puts the bearer token in `argv`
   (visible in process listings and shell history). `hosted serve` already has
   an env fallback (`HOSTED_TOKEN`); `config set` has none.

## Current state

- `internal/devspace/hosted_sync.go:60-99` — `SetHostedSync` contains the
  scheme validation to extract:

```go
parsed, err := url.Parse(endpoint)
if err != nil || parsed.Scheme == "" || parsed.Host == "" {
    return Config{}, fmt.Errorf("hosted sync endpoint must be an absolute http(s) URL")
}
if parsed.Scheme != "http" && parsed.Scheme != "https" {
    return Config{}, fmt.Errorf("hosted sync endpoint must use http or https")
}
if parsed.Scheme == "http" && !isLoopbackHost(parsed.Hostname()) {
    return Config{}, fmt.Errorf("hosted sync endpoint must use https (plain http is only allowed for localhost)")
}
```

- `internal/devspace/hosted_sync.go:101-122` — `GetHostedSync` checks only
  non-empty endpoint/token and a well-formed workspace ID. It also implicitly
  writes a default workspace back to config; leave that behavior alone.
- `internal/devspace/commands.go:162-180` — `newHostedConfigCommand`'s `set`
  subcommand: `set.Flags().StringVar(&token, "token", "", "hosted sync bearer token")`.
- `internal/devspace/commands.go:215-217` — the serve command's existing env
  fallback, the pattern to mirror:

```go
if token == "" {
    token = strings.TrimSpace(os.Getenv("HOSTED_TOKEN"))
}
```

## Commands you will need

| Purpose | Command | Expected on success |
|---------|---------|---------------------|
| Full gate | `make verify` | exit 0 |
| Targeted | `go test ./internal/devspace -run Hosted -v` | PASS |

## Scope

**In scope**:

- `internal/devspace/hosted_sync.go`
- `internal/devspace/commands.go` (only `newHostedConfigCommand`)
- `internal/devspace/hosted_sync_test.go`

**Out of scope**:

- The hosted **server** (`NewHostedSyncServer` and below) — its hardening
  contract lives in `hardening_test.go`; nothing here should touch it.
- `SetHostedSync`'s user-facing error strings — reuse them verbatim via the
  extracted helper so existing tests keep passing.

## Git workflow

- Branch: `advisor/005-hosted-client-hardening`
- Conventional commit, e.g. `fix(hosted): re-validate endpoint scheme at point of use; env-var token input`

## Steps

### Step 1: Extract `validateHostedEndpoint`

In `hosted_sync.go`, move the three-check block quoted above into:

```go
func validateHostedEndpoint(endpoint string) error { ... }
```

Call it from `SetHostedSync` (behavior identical, same error text) and add to
`GetHostedSync`, after the empty-endpoint check:

```go
if err := validateHostedEndpoint(cfg.HostedSyncEndpoint); err != nil {
    return Config{}, fmt.Errorf("configured hosted sync endpoint is invalid: %w; re-run `devspace hosted config set`", err)
}
```

**Verify**: `go test ./internal/devspace -run Hosted -v` → existing tests PASS

### Step 2: Env-var token fallback for `hosted config set`

In `newHostedConfigCommand` (`commands.go:162`), mirror the serve pattern: when
the `--token` flag is empty, read `DEVSPACE_HOSTED_TOKEN`:

```go
if token == "" {
    token = strings.TrimSpace(os.Getenv("DEVSPACE_HOSTED_TOKEN"))
}
```

Update the flag help string to:
`"hosted sync bearer token (prefer DEVSPACE_HOSTED_TOKEN to keep the token out of shell history and process listings)"`.

(Name is deliberately `DEVSPACE_HOSTED_TOKEN`, not `HOSTED_TOKEN`: the client
setting and the server flag are different credentials in different processes;
don't unify them.)

**Verify**: `go build ./...` → exit 0

### Step 3: Tests

In `hosted_sync_test.go` (match its existing style — it builds configs via
`t.Setenv(envHome, …)` and calls the exported functions directly):

1. `TestGetHostedSyncRejectsPlainHTTPEndpoint`: write a `config.json` (via
   `SaveConfig`) whose `HostedSyncEndpoint` is `http://example.com` with a
   token and workspace set; assert `GetHostedSync` errors mentioning https.
2. `TestGetHostedSyncAllowsLoopbackHTTP`: same with
   `http://127.0.0.1:8787` → no error.
3. `TestHostedConfigSetReadsTokenFromEnv`: `t.Setenv("DEVSPACE_HOSTED_TOKEN", "…")`,
   execute the cobra command `hosted config set https://example.com` with no
   `--token` flag (build the root command via `NewRootCommand()` and set args),
   then `GetHostedSync` succeeds and the config's token equals the env value.

**Verify**: `go test ./internal/devspace -run 'TestGetHostedSync|TestHostedConfigSet' -v` → PASS

### Step 4: Full gate

**Verify**: `make verify` → exit 0

## Test plan

The three tests in Step 3. Note for test 3: never echo the token value in
assertions' failure messages beyond a placeholder.

## Done criteria

- [ ] `make verify` exits 0
- [ ] `grep -n "validateHostedEndpoint" internal/devspace/hosted_sync.go` shows both call sites (Set + Get)
- [ ] `grep -n "DEVSPACE_HOSTED_TOKEN" internal/devspace/commands.go` shows the fallback
- [ ] New tests pass; no files outside scope modified
- [ ] `plans/README.md` status row updated

## STOP conditions

- Existing hosted tests fail after Step 1 — the extraction changed behavior;
  report the diff in error text/flow rather than editing the tests.
- `GetHostedSync`'s save-back of a default workspace interferes with the new
  validation ordering — report; do not reorder its side effects.

## Maintenance notes

- Any new field read from `config.json` that has an integrity requirement
  should follow this pattern: validate at point of use, not only at write time.
- Docs follow-up (out of scope here): mention `DEVSPACE_HOSTED_TOKEN` in the
  README's hosted-sync section.
