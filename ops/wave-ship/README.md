# DevDrop MVP Wave-Ship Prep

## Workflow

Use:

```text
/Users/lecoqjacob/Projects/personal/inkwell/.claude/workflows/wave-ship.js
```

The requested path with `inkwell.claude` was not present locally. The checked-in
workflow path is under the `inkwell` repository's `.claude/workflows/` directory.

## Args

Pass `devdrop-mvp.args.json` as the workflow args JSON.

The prepared run uses:

- `repo`: `/Users/lecoqjacob/Projects/personal/devdrop`
- `project`: `DevDrop MVP`
- `team`: `Cypress Ink Labs`
- `base`: `main`
- `backend`: `orca`
- `serializedMerge`: `true`
- `autoContinue`: `false`
- `maxConcurrent`: `3`
- `engine`: `codex`

## Linear Cards

Wave 1:

- `CIL-217` — DevDrop: scaffold Go CLI and init command

Wave 2:

- `CIL-218` — DevDrop: manifest model, workspace scan, and status
- `CIL-219` — DevDrop: project add, workspace sync, and Git hydration
- `CIL-220` — DevDrop: encrypted env profile commands

Wave 3:

- `CIL-221` — DevDrop: README, examples, and end-to-end MVP verification

## Boundaries

- Workers open PRs and do not merge.
- `wave-ship` owns serialized merge and ticket closeout.
- Do not build hosted sync, FUSE, a daemon, team secrets, editor settings sync,
  or dependency auto-install in this run.
- Do not put real secrets in docs, fixtures, logs, or tests.

