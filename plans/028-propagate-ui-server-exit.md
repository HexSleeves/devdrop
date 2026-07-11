# Plan 028: Exit devspace-tui with an error when ui-server dies

> **Executor instructions**: Follow each step and stop on the conditions below.
> Update `plans/README.md` when done.
>
> **Drift check (run first)**: `git diff --stat 7b521c3..HEAD -- tui/src/client.ts tui/src/app.tsx tui/src/main.tsx tui/test/client.test.ts`

## Status

- **Priority**: P2
- **Effort**: S
- **Risk**: LOW
- **Depends on**: plans/020-runtime-validate-tui-rpc-responses.md
- **Category**: bug
- **Planned at**: commit `7b521c3`, 2026-07-10

## Why this matters

The client captures ui-server stderr and supplies it to close listeners, but
`App` discards the error. `main.tsx` interprets an absent message as exit code
zero, so a backend crash after startup looks like a successful user quit. Keep
intentional quits successful and surface unexpected server termination on
stderr with a nonzero exit.

## Current state

```ts
// tui/src/client.ts:198-201
await pumpText(proc.stdout, (text) => client.feed(text));
const tail = stderrLines.join("\n").trim();
client.closed(tail ? new Error(`devspace ui-server exited: ${tail}`) : undefined);
```

```ts
// tui/src/app.tsx:94
const offClose = client.onClose(() => quit());
```

```ts
// tui/src/main.tsx:54
<App ... quit={(message) => quit(message === undefined ? 0 : 1, message)} />
```

`tui/test/client.test.ts:119-150` already proves close errors reach listeners
and duplicate close signals notify once. Preserve this behavior. Plan 020 owns
the same client dispatch area and must land first to avoid conflicting edits.

## Commands you will need

| Purpose | Command | Expected on success |
|---|---|---|
| Focused tests | `cd tui && bun test test/client.test.ts test/startup.test.ts` | all pass |
| Typecheck | `cd tui && bun run typecheck` | exit 0 |
| Full gate | `make tui-verify` | exit 0 |

## Scope

**In scope**:

- `tui/src/client.ts`
- `tui/src/app.tsx`
- `tui/src/main.tsx` only if its existing message-to-exit-code mapping must change
- `tui/test/client.test.ts`

**Out of scope**:

- Go ui-server lifecycle.
- Restarting the server automatically.
- Changing SIGINT/SIGTERM exit codes.
- Adding a React/TUI test framework.

## Git workflow

- Branch: `advisor/028-ui-server-exit-status`
- Commit: `fix(tui): report unexpected ui-server exit`
- Do not push or open a PR unless instructed.

## Steps

### Step 1: Always describe process termination

When the spawned server stdout closes, construct an Error even when captured
stderr is empty. Use `devspace ui-server exited` as the fallback and retain the
bounded stderr tail when present. Do not change `DevspaceClient.closed()` for
in-memory/manual callers.

Add a small directly testable function only if necessary to test fallback and
stderr-tail formatting; otherwise keep the change inline and use existing close
listener tests. Do not create a process-management abstraction.

**Verify**: focused client tests pass.

### Step 2: Forward the close error to quit

Change `App`'s close listener to pass the error message to `quit`. If an
undefined error can still arrive through a custom client, supply the same
fallback message. Preserve cleanup: the returned unsubscribe function remains
called by the effect cleanup.

Intentional keyboard quit still calls `quit()` directly. `main.tsx` may retain
its existing mapping: message absent → 0; message present → 1.

**Verify**: `cd tui && bun run typecheck` → exit 0.

### Step 3: Verify exit semantics without new infrastructure

Extend client tests for an empty-diagnostic unexpected close only if a helper
was introduced. Confirm by source assertion that `onClose` forwards its `err`
instead of ignoring it:

`rg -n 'onClose\(\(err\).*quit' tui/src/app.tsx` → one match.

**Verify**: focused tests and `make tui-verify` pass.

## Test plan

Reuse the close-listener tests in `client.test.ts`. Cover stderr context,
fallback context, one notification, and requests rejected on close. Do not add a
renderer dependency solely to test a one-line App callback.

## Done criteria

- [ ] Unexpected ui-server termination always has an error message.
- [ ] `App` forwards close errors to `quit`.
- [ ] Intentional `q`/Ctrl-C behavior remains successful/standard.
- [ ] Focused tests, typecheck, and full TUI gate pass.
- [ ] Only in-scope files and `plans/README.md` changed.

## STOP conditions

- Plan 020 has not landed and produces overlapping client changes.
- OpenTUI invokes close listeners during normal renderer shutdown before `quitting` is set.
- Correct behavior requires automatic restart or a new process supervisor.

## Maintenance notes

Unexpected child-process termination is an application failure. Keep intentional
user shutdown distinguishable by initiating `quit()` before closing the client.
