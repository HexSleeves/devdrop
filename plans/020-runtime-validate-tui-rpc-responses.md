# Plan 020: Runtime-validate every devspace-tui RPC result and server event

> **Executor instructions**: Follow this plan step by step. Run every
> verification command and confirm the expected result before moving to the
> next step. If anything in the "STOP conditions" section occurs, stop and
> report — do not improvise. When done, update the status row for this plan in
> `plans/README.md`.
>
> **Drift check (run first)**: `git diff --stat 7b521c3..HEAD -- tui/src/client.ts tui/src/protocol.ts tui/test/client.test.ts tui/test/protocol.test.ts tui/test/startup.test.ts`
> If any in-scope file changed, compare the excerpts below with live code. Stop
> if the method map, validator contracts, or dispatch shape no longer match.

## Status

- **Priority**: P1
- **Effort**: M
- **Risk**: LOW
- **Depends on**: none
- **Category**: bug / tests
- **Planned at**: commit `7b521c3`, 2026-07-10

## Why this matters

The TypeScript TUI has runtime validators for every response DTO and unsolicited
server event, but `DevspaceClient` trusts parsed JSON after `JSON.parse`.
Version-skewed or malformed data can therefore enter React state as if it were
typed. Validate both request results and events at the single transport
boundary, with explicit errors for bad results and safe rejection of bad events.

## Current state

- `tui/src/client.ts` stores no method on a pending request and resolves the
  result directly:

```ts
// tui/src/client.ts:124-130
const pending = this.pending.get(msg.id);
if (!pending) return;
this.pending.delete(msg.id);
if (pending.timer) clearTimeout(pending.timer);
if (msg.error) pending.reject(new Error(msg.error.message));
else pending.resolve(msg.result);
```

- The same dispatch path forwards `msg.params` as `ServerEvent` without calling
  the validator:

```ts
// tui/src/client.ts:115-122
if (msg.method === "event" && msg.params) {
  if (this.eventListeners.size === 0) this.earlyEvents.push(msg.params);
  for (const listener of this.eventListeners) listener(msg.params);
}
```

- `tui/src/protocol.ts:312-371` already provides `isHello`, `isSnapshot`,
  `isSyncStatus`, `isWorkspaceOverview`, and `isServerEvent`.
- `tui/src/protocol.ts:381-393` defines `RequestMap`; it is the authoritative
  method-to-result map.
- Tests use Bun's runner and the in-memory `pair()` transport in
  `tui/test/client.test.ts`. Keep handwritten validators; do not add a schema
  dependency.

## Commands you will need

| Purpose | Command | Expected on success |
|---|---|---|
| Typecheck | `cd tui && bun run typecheck` | exit 0 |
| Focused tests | `cd tui && bun test test/client.test.ts test/protocol.test.ts test/startup.test.ts` | all pass |
| Full gate | `make tui-verify` | exit 0 |

## Scope

**In scope**:

- `tui/src/protocol.ts`
- `tui/src/client.ts`
- `tui/test/client.test.ts`
- `tui/test/protocol.test.ts`
- `tui/test/startup.test.ts` only if needed for a startup error assertion

**Out of scope**:

- Go DTO or wire-format changes.
- Regenerating fixtures unless the live Go DTO changed.
- Adding a validation dependency.
- UI rendering or exit-code behavior; plan 028 owns unexpected server exits.

## Git workflow

- Branch: `advisor/020-tui-runtime-validation`
- Commit: `fix(tui): validate ui-server messages at runtime`
- Do not push or open a PR unless instructed.

## Steps

### Step 1: Add one explicit result parser

In `tui/src/protocol.ts`, add `parseResult<M extends Method>(method, result)`.
Use a `switch` over every `RequestMap` method: hello → `isHello`; snapshot
methods → `isSnapshot`; status → `isSyncStatus`; workspace →
`isWorkspaceOverview`; lastPlan → the existing plan validator. Throw
`invalid <method> response from devspace ui-server` when validation fails.

**Verify**: `cd tui && bun run typecheck` → exit 0.

### Step 2: Validate results before resolving

Add `method: Method` to `Pending`, store it in `request`, and pass
`msg.result` through `parseResult(pending.method, msg.result)` inside a
`try/catch`. Preserve pending deletion and timer cleanup before resolve/reject.

**Verify**: `cd tui && bun run typecheck` → exit 0.

### Step 3: Validate unsolicited events

Treat parsed `params` as `unknown`, not `ServerEvent`. For `method === "event"`,
call `isServerEvent`; ignore invalid events without buffering or notifying
listeners. Do not close the stream: one malformed line must not break later
valid traffic.

**Verify**: `cd tui && bun run typecheck` → exit 0.

### Step 4: Add boundary regression tests

Using `pair()` in `tui/test/client.test.ts`, prove:

- malformed hello, status, and snapshot results reject with method-specific errors;
- a valid result still resolves and clears its timer;
- malformed and unknown events are not delivered or buffered;
- a valid event after a malformed event is still delivered.

Add a direct `parseResult` fixture assertion in `protocol.test.ts` only if it
makes the method map coverage clearer.

**Verify**: `cd tui && bun test test/client.test.ts test/protocol.test.ts test/startup.test.ts` → all pass.

### Step 5: Run the TUI gate

**Verify**: `make tui-verify` → exit 0.

## Test plan

Follow the existing `pair()` request/response tests in
`tui/test/client.test.ts:28-100`. Cover invalid results, invalid events, recovery
on the next frame, valid out-of-order responses, and timeout cleanup.

## Done criteria

- [ ] Every `RequestMap` method is handled by `parseResult`.
- [ ] `pending.resolve(msg.result)` no longer exists.
- [ ] Events reach listeners only after `isServerEvent` succeeds.
- [ ] Focused Bun tests pass.
- [ ] `make tui-verify` passes.
- [ ] Only in-scope files and `plans/README.md` changed.
- [ ] Plan 020 is marked DONE in `plans/README.md`.

## STOP conditions

- `RequestMap` differs from the listed live methods.
- Validation requires a Go DTO or protocol-version change.
- Live code intentionally accepts partial objects rejected by existing validators.
- The full gate fails for an unrelated toolchain reason after focused tests pass.

## Maintenance notes

Every future UI-server method needs a `RequestMap` entry and parser branch.
Every future event variant needs an `isServerEvent` branch. Reviewers should
reject new dispatch paths that cast parsed JSON directly to protocol types.
