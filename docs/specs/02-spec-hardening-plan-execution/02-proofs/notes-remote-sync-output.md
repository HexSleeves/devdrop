# Remote Sync/Output Branch Reconciliation Notes

Scope: Spec 02 Task 3.5 analysis for branch commits `e4f7bd0`, `71919e8`, `9b4da48`, and `973b6f1`. No branch commit was implemented.

| Commit | Recommendation | Rationale |
| --- | --- | --- |
| `e4f7bd0` | defer | Configurable manifest commit identity is useful for auditability, but it changes config schema, command flags, and sync commit behavior outside Plan 004's remote validation scope. Revisit with Plan 013 or a dedicated sync-audit task. |
| `71919e8` | defer | Moving pure output helpers from `commands.go` to `output.go` is a reasonable refactor, but it is not required for Plan 004 and would add a new implementation file outside the current write scope. |
| `9b4da48` | defer | Command output and CLI wiring tests are useful only if the output-helper refactor is kept. The new `commands_test.go` file is outside the current write scope and should not be landed as part of remote clone hardening. |
| `973b6f1` | reject for this task | SHA-1 naming comments in `secrets.go` are directionally helpful documentation, but they touch secret identity code and do not support remote sync/output reconciliation or Plan 004. Reconsider in a security-doc cleanup task if needed. |

Plan 004 note: remote validation is enforced both during manifest validation and again at hydrate time.
