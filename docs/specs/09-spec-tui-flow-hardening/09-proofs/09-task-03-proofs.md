# Task 03 Proofs - stream-safe client, visible and recoverable watcher death, legacy backoff parity, width-aware cells

## Task Summary

This task closed the TUI flow's silent failure modes per
`plans/018-tui-client-and-watch-robustness.md` (executed by Codex under
orchestrator review). The NDJSON client now decodes stdout with streaming
UTF-8 semantics (`pumpText`, shared by stdout and stderr) and buffers events
that arrive before the first listener; the reducer ignores unknown event
types and marks the watcher alive again on recovery; the ui-server emits a
`watch-error` event on every watcher exit path and restarts a dead watcher
after the next successful scan/refresh; the legacy Bubble Tea dashboard gained
the same exponential backoff / give-up policy the ui-server got in PR #39
(shared constants); rapid keypresses can no longer double-fire actions
(`busyRef`); and table cells pad/truncate by terminal display width
(`Bun.stringWidth` in the new `tui/src/text.ts`).

## What This Task Proves

- A multi-byte UTF-8 character split across pipe chunks can no longer corrupt
  a frame or hang a request.
- Watcher death is always visible (event on every exit path) and recoverable
  (restart on next successful scan/refresh), in the client indicator too
  (`watchAlive` flips back on the next watch-refresh).
- The legacy dashboard no longer hot-loops on persistent watcher errors:
  backoff doubles to a cap, gives up after `watchRetryMaxAttempts`, resets on
  success.
- Unknown server event types are a no-op instead of wiping the project table.
- CJK/emoji project names keep the table columns aligned.

## Evidence Summary

- Go: 24 race-enabled tests pass in the watch/dashboard selection, including
  the 6 new lifecycle tests below.
- TUI: 29 pass / 0 fail across 4 files (adds split-UTF-8 pump, early-event
  buffering, unknown-event no-op, watchAlive recovery, and width-aware cell
  tests); `tsc --noEmit` clean.
- `make verify` exit 0 on the main tree.
- The instant re-arm is gone: `grep "return m, m.nextWatchCmd()"` has 0 matches.

## Artifact: Watch lifecycle tests (server + legacy dashboard) under -race

**What it proves:** No silent watcher death, restart-on-next-action, and
dashboard backoff parity — the Unit 3 server-side FRs.

**Command:**

~~~bash
go test ./internal/devspace -race -run 'TestUIServerWatchClosedEmitsEvent|TestUIServerScanRestartsDeadWatcher|TestDashboardWatch' -v
~~~

**Result summary:** All 7 pass, race-clean.

~~~text
--- PASS: TestUIServerWatchClosedEmitsEvent (0.14s)
--- PASS: TestUIServerScanRestartsDeadWatcher (0.15s)
--- PASS: TestDashboardWatchRefreshUpdatesModel (0.00s)
--- PASS: TestDashboardWatchErrorRearmsWithBackoff (0.00s)
--- PASS: TestDashboardWatchErrorStopsAfterMaxAttempts (0.00s)
--- PASS: TestDashboardWatchSuccessResetsErrorCount (0.00s)
--- PASS: TestDashboardWatcherEmitsRefreshOnFileChange (0.43s)
ok  	github.com/liatrio-forge/devdrop-capstone/internal/devspace	1.747s
~~~

## Artifact: TUI client/reducer/text suite

**What it proves:** Streaming decode, early-event buffering, defensive
reducer with watchAlive recovery, and display-width cells — the Unit 3
client-side FRs.

**Command:**

~~~bash
cd tui && bun test && bun run typecheck
~~~

**Result summary:** 29 pass, 0 fail, 73 expect() calls across 4 files;
typecheck clean.

~~~text
 29 pass
 0 fail
 73 expect() calls
Ran 29 tests across 4 files. [121.00ms]
~~~

## Artifact: Instant re-arm removed from the legacy dashboard

**Command:**

~~~bash
grep -n "return m, m.nextWatchCmd()" internal/devspace/ui_model.go
~~~

**Result summary:** 0 matches (exit 1) — the error path now re-arms through a
by-value delayed command with doubling backoff, or gives up with a
"watcher stopped after N consecutive errors" message.

## Artifact: Full gates

**Command:**

~~~bash
make verify   # exit 0 (test+vet+lint+vulncheck+build)
make tui-verify  # bun typecheck + 29 tests, exit 0
~~~

**Result summary:** Both green on the main tree.

## Reviewer Conclusion

Unit 3's functional requirements are implemented and proven end to end: the
client survives real-world byte streams and input timing, watcher failure is
visible and self-healing in both frontends, and the remaining cosmetic
defects (unknown events, wide characters) are covered by tests. Reviewer
scrutiny points from the plan were verified: the restart guard uses
CompareAndSwap (no double watch loops) and the delayed re-arm closure
captures its delay and command by value.
