# DevDrop

DevDrop is a local-first "Dropbox for developers" CLI prototype. It keeps a
developer workspace structurally aligned across machines by tracking project
metadata, safe placeholder folders, Git remotes, setup hints, and encrypted env
profiles.

The intended command name is `devspace`. The current binary can also be built
from this repository:

```bash
go build -o bin/devspace ./cmd/devdrop
```

During development, you can still run the command directly from source:

```bash
go test ./...
go run ./cmd/devdrop --help
```

## Current MVP Status

This MVP is intentionally local-only. It proves the workflow before adding
hosted sync, background daemons, or filesystem-level lazy loading.

What works today:

- Initialize a workspace and local machine identity.
- Generate a versioned workspace manifest.
- Scan existing projects for Git metadata, `.env` presence, dependency files,
  and setup hints.
- Show workspace and project status.
- Add local projects to the manifest.
- Plan and apply safe missing project structure as empty placeholder folders.
- Hydrate placeholder Git projects with normal `git clone`.
- Store encrypted per-project env profiles with native age encryption.
- Generate local `.env` files with `0600` permissions.

## Supported Commands

### `devspace init`

```bash
devspace init --workspace ~/code
```

Creates:

- `~/.devdrop/config.json`
- `~/.devdrop/state.json`
- `~/.devdrop/identity.txt`
- `<workspace>/.devdrop/manifest.json`

The command is idempotent and does not rotate the machine ID or age identity on
repeat runs.

### `devspace scan`

```bash
devspace scan
```

Scans the configured workspace and updates the manifest/state with:

- Git remote URL
- Current branch
- Last commit
- Dirty working tree status
- `.env` file presence
- Dependency/setup hints

Ignored by default:

```text
node_modules/
dist/
build/
.next/
turbo/
target/
vendor/
coverage/
.cache/
.DS_Store
*.log
```

### `devspace status`

```bash
devspace status
devspace project status client-a-api
```

Shows tracked projects, hydrated projects, placeholders, dirty repos, missing env
files, stale/missing projects, and last scan/sync timestamps.

### `devspace project add`

```bash
devspace project add work/client-a-api
```

Adds a relative workspace path to the manifest. Existing Git repositories are
recorded as Git projects; otherwise the project is tracked as local-only.

### `devspace plan`

```bash
devspace plan
devspace plan --json
```

Builds a deterministic plan of safe and skipped actions, saves it to
`<workspace>/.devdrop/last-plan.json`, and prints a human-readable report.
`--json` prints the same saved plan as structured JSON for automation.

### `devspace apply`

```bash
devspace apply
```

Applies the last saved plan only if the manifest hash still matches. If the
manifest changed after `plan`, `apply` refuses to run and asks you to re-plan.
Only safe actions are executed. Skipped actions remain listed and untouched.

### `devspace workspace sync`

```bash
devspace workspace sync --dry-run
devspace workspace sync
```

Compatibility alias. Prefer `devspace plan` and `devspace apply`.

### `devspace project hydrate`

```bash
devspace project hydrate client-a-api
```

Hydrates a placeholder Git project with normal `git clone`. The placeholder is
an empty directory; hydration refuses to clone into non-empty directories and
does not delete files.

### `devspace env`

```bash
printf 'postgres://local\n' | devspace env set client-a-api DATABASE_URL
devspace env list client-a-api
devspace env pull client-a-api
```

`env set` reads from stdin when piped, or prompts for a hidden value when run in a
terminal. `env list` prints keys only, with values masked. `env pull` writes the
local project `.env` with `0600` permissions.

Encrypted profiles are stored under:

```text
<workspace>/.devdrop/secrets/<project-id>/<profile>.age
```

## Example Local Workflow

This workflow uses temporary directories and a local bare Git remote, so it does
not need network access.

```bash
go build -o bin/devspace ./cmd/devdrop

tmp="$(mktemp -d)"
export DEV_DROP_HOME="$tmp/home"
workspace_a="$tmp/workspace-a"
remote_src="$tmp/remote-src"
remote_bare="$tmp/client-a-api.git"

mkdir -p "$remote_src"
git -C "$remote_src" init -b main
git -C "$remote_src" config user.email demo@example.com
git -C "$remote_src" config user.name "Demo User"
printf '# client-a-api\n' > "$remote_src/README.md"
git -C "$remote_src" add README.md
git -C "$remote_src" commit -m "initial"
git clone --bare "$remote_src" "$remote_bare"

bin/devspace init --workspace "$workspace_a"
mkdir -p "$workspace_a/work/client-a-api"
git clone "$remote_bare" "$workspace_a/work/client-a-api"
printf '{"scripts":{"dev":"vite"}}\n' > "$workspace_a/work/client-a-api/package.json"
bin/devspace scan
bin/devspace plan
bin/devspace apply
printf 'postgres://demo\n' | bin/devspace env set client-a-api DATABASE_URL
bin/devspace env list client-a-api
bin/devspace env pull client-a-api
bin/devspace status
```

## Two-Machine Simulation Workflow

To simulate a second machine, copy the first workspace manifest into a second
workspace and run plan/apply there.

```bash
workspace_b="$tmp/workspace-b"

bin/devspace init --workspace "$workspace_b"
mkdir -p "$workspace_b/.devdrop"
cp "$workspace_a/.devdrop/manifest.json" "$workspace_b/.devdrop/manifest.json"
bin/devspace plan
bin/devspace apply
bin/devspace status
```

The second workspace now contains placeholder directories for tracked projects.
Run `bin/devspace project hydrate <project>` to clone a placeholder Git project
when the manifest includes its remote.

## Safety Guarantees

- Project paths in the manifest are relative to the workspace root.
- Absolute paths and parent-directory escapes are rejected.
- `plan` reports planned filesystem changes before writing anything.
- `plan --json` emits the same plan as machine-readable JSON.
- `apply` executes only a saved plan whose manifest hash still matches.
- `apply` creates missing directories only; it skips non-empty destinations.
- Plan reports path, dirty-repo, and remote conflicts instead of overwriting local work.
- Hydration clones only into missing or empty directories.
- Env values are encrypted at rest with age.
- `env list` masks secret values.
- `env pull` writes `.env` with `0600` permissions.
- The MVP has no hosted control plane and no background process.

## What This Tool Will Not Do Without Permission

- Delete local projects, files, or directories.
- Overwrite existing project contents during apply or hydrate.
- Auto-pull, rebase, merge, or push Git repositories.
- Resolve Git conflicts.
- Install dependencies or run project setup commands.
- Upload manifests, secrets, source code, or project metadata to a hosted
  service.
- Share env profiles with teammates.
- Rotate secrets or replace local `.env` values without an explicit command.
- Watch the filesystem in the background.
- Mount a FUSE or virtual filesystem.

## Known Limitations

- Manifest exchange is manual; there is no Git-backed or hosted manifest sync.
- Placeholder hydration uses full `git clone`; partial clone and sparse checkout
  are not implemented.
- Secret profiles are local to the workspace; there is no team sharing, OS
  keychain integration, remote backup, or rotation flow.
- Setup hints are informational only and do not install dependencies.
- Editor settings, VS Code extensions, devcontainers, Nix, mise, and asdf are
  outside the MVP.
- The command name is documented as `devspace`, while the current source package
  still lives under `cmd/devdrop`.

## Troubleshooting

- If `devspace` is not found, build it with
  `go build -o bin/devspace ./cmd/devdrop` and run `bin/devspace`.
- If commands use the wrong workspace, check `DEV_DROP_HOME` and
  `~/.devdrop/config.json`.
- If `plan` reports a path conflict, inspect the existing directory and
  decide whether it should be tracked, renamed, or left unmanaged.
- If `project hydrate` fails, confirm the project has a remote URL in
  `<workspace>/.devdrop/manifest.json` and that `git clone <remote> <path>` works
  by itself.
- If `env pull` cannot write `.env`, check project path permissions and whether a
  directory or symlink exists at `.env`.
- If secret values appear missing, run `devspace env list <project>` to verify
  the key exists before pulling.

## Roadmap

- Git-backed manifest push/pull.
- Hosted API or cloud object storage.
- Background daemon and filesystem watchers.
- FUSE or virtual filesystem lazy loading.
- Git partial clone or sparse checkout.
- Team/shared secret profiles.
- OS keychain integration.
- Secret rotation.
- Editor settings, devcontainer, Nix, mise, or asdf sync.
- Optional dependency setup flows with explicit confirmation.
- Release-readiness checklist automation.

## Manifest

See [`examples/manifest.json`](examples/manifest.json).

The manifest is a versioned JSON file at:

```text
<workspace>/.devdrop/manifest.json
```

Project paths are always relative to the workspace root. Absolute paths and
parent-directory escapes are rejected.
