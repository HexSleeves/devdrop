# Task 1 Proofs: Command Taxonomy and Status

## Result

PASS. Root help exposes exactly 14 product commands in four goal-oriented groups. Removed root commands are absent, `--version` remains available, and `status` owns aggregate, verbose, and project-specific views. Captures below use an isolated temporary application home and replace machine, path, ID, and timestamp values with placeholders.

## Automated Contract

```text
$ go test ./internal/devspace -run 'TestReleaseCommandTreeContract|TestStatusCommand' -count=1
ok  github.com/liatrio-forge/devdrop-capstone/internal/devspace
```

The contract covers the 14-command maximum, group titles, canonical names, root `--version`, help examples, rejection of `workspace`, `tui`, `mount`, `version`, `project status`, and invalid `sync` paths, verbose redaction, project selection, JSON parsing, ANSI-free JSON, and invalid status flag/argument combinations.

## Grouped Root Help

```text
$ go run ./cmd/devspace --help
EXAMPLES
  devspace init --workspace ~/code
  devspace scan
  devspace status --verbose
  devspace plan && devspace apply

CORE WORKFLOW:
  apply
  init
  plan
  scan
  status [project]

WORKSPACE MANAGEMENT:
  env
  hosted
  project
  setup
  sync

DIAGNOSTICS AND AUTOMATION:
  doctor
  ui
  watch

EXPERIMENTAL:
  experimental

FLAGS
  --version
```

Visible product commands: 14. Visible `workspace`, `tui`, `mount`, and `version` commands: 0.

## Isolated Status Workflow

Fixture: one local project at `<workspace>/apps/api`, discovered by `scan` under `<DEVSPACE_HOME>`.

```text
$ devspace --no-color status
Workspace Status
Machine: <machine>
Workspace: <workspace>

Projects tracked: 1
Hydrated: 1
Placeholders: 0
Dirty repos: 0
Missing env files: 1
Outdated repos: 0
Last scan: <timestamp>
```

```text
$ devspace --no-color status --verbose
Workspace
Root: <workspace>
Manifest version: 1
This machine: <machine> (<machine-id>)

Machines
NAME       ID            LAST SEEN
<machine>  <machine-id>  <timestamp>

Users: -
Teams: -

Sync
Manifest remote: -
Hosted endpoint: -
Last sync: -
Last scan: <timestamp>

Summary
Projects tracked: 1
Hydrated: 1
Placeholders: 0
Dirty repos: 0
Missing env files: 1
Outdated repos: 0
```

```json
$ devspace --no-color status api --json
{
  "project": {
    "id": "<project-id>",
    "name": "api",
    "path": "apps/api",
    "type": "local",
    "hydrateMode": "manual"
  },
  "state": {
    "hydrated": true,
    "exists": true,
    "dirty": false,
    "envFilePresent": false,
    "lastCheckedAt": "<timestamp>",
    "placeholder": false,
    "stale": false,
    "missing": false
  }
}
```

The JSON capture parses as one `ProjectListRow` and contains no ANSI escape bytes.
