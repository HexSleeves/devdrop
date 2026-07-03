---
name: sdd-html
description: "Execute the Liatrio Spec-Driven Development (SDD) workflow with HTML-formatted artifacts (spec, tasks, audit, proofs, validation) instead of Markdown, optionally using Lavish Editor for interactive browser-based review. NOTE: this skill is NOT intended to be dynamically loaded or automatically triggered; it should only ever be explicitly called by the user."
---

# Spec-Driven Workflow Orchestrator (HTML Edition)

You are the orchestrator for the Spec-Driven Development (SDD) workflow, HTML edition. This is a variant of the `sdd` skill: identical lifecycle, identical roles, gates, and checklists, but every generated artifact (questions, spec, tasks, audit, proofs, validation) is a self-contained HTML document instead of Markdown. Determine the current lifecycle phase from the workspace, then load and follow the single matching phase reference.

**Do not mix workflows.** This skill's artifacts use the `.html` extension end-to-end. The original `sdd` skill uses `.md`. Do not generate `.md` artifacts under this skill, and do not run this skill against a spec directory that already contains `.md` SDD artifacts (offer to start a fresh `[NN]-spec-...` sequence instead).

## Skill Contract

This skill manages one SDD phase per invocation unless the loaded phase explicitly requires a smaller checkpoint. It must:

1. Assess workspace state before routing.
2. Check Lavish Editor availability before routing (see below).
3. Select exactly one active spec and one phase reference.
4. Load only that phase reference.
5. Execute the phase instructions or stop at a required approval gate.
6. Tell the user the selected spec, detected state, current phase, chosen phase reference path, created or updated artifacts, Lavish availability, and next recommended natural-language request.

If multiple active specs exist and the user request does not clearly identify one, ask the user which spec to continue before modifying files. If the selection is unambiguous, state why that spec was selected.

## State Assessment

Before routing, run the bundled assessor script while the current working directory is the target repository/workspace root:

```bash
python3 {{skill_dir}}/scripts/assess-sdd-state.py .
```

`{{skill_dir}}` locates bundled skill assets only; SDD artifacts are assessed from the workspace path argument (`.` in the command above). The script scans `docs/specs/` for `.html` artifacts and outputs JSON describing active specs and a recommended phase. Use that output as the default source of truth for routing.

Use manual inspection only if the script cannot run. When falling back manually, state that fallback explicitly, inspect `./docs/specs/` for `[NN]-spec-*.html` / `[NN]-tasks-*.html` / `[NN]-audit-*.html` / `[NN]-validation-*.html`, and apply the same phase rules below.

## Lavish Availability Check

Every phase reference includes an optional interactive-review step built on the Lavish Editor (`lavish-axi`). Determine availability once per invocation, before routing, and reuse the result across the phase:

```bash
if command -v lavish-axi >/dev/null 2>&1; then
  echo "lavish: available (lavish-axi on PATH)"
elif [ -f "$HOME/.claude/skills/lavish/SKILL.md" ]; then
  echo "lavish: available (lavish skill present, invoke via: npx -y lavish-axi)"
else
  echo "lavish: not available"
fi
```

- **If available:** phase references may open generated HTML artifacts as Lavish sessions (`npx -y lavish-axi <html-file>` if no global binary, otherwise `lavish-axi <html-file>`) so the user can visually review, annotate, and reply to specific sections instead of reading raw HTML or editing markup by hand. Follow the lavish skill's own workflow conventions (create under `.lavish/` is not required here — SDD artifacts stay under `docs/specs/`; open the artifact in place) for opening, polling, and ending sessions.
- **If not available:** skip Lavish entirely. Tell the user the plain HTML file path and that they can open it directly in any browser; collect feedback through normal conversation instead of a Lavish session.
- Never install or auto-fetch tooling on the user's behalf beyond the standard `npx -y lavish-axi` invocation documented by the lavish skill. Never fetch Lavish tools speculatively — only invoke them once you have a concrete HTML artifact ready for review.

## Phase Rules

- **Phase 1 — Spec Generation**
  - **Condition:** The user requested a new feature, but no matching `[NN]-spec-[feature-name]/` directory exists, or the spec directory exists but the spec document is incomplete or missing.
  - **Action:** Gather context and write the formal specification as HTML.

- **Phase 2 — Task List Generation and Audit**
  - **Condition:** `[NN]-spec-[feature-name].html` exists and is complete, but `[NN]-tasks-[feature-name].html` is missing, or the mandatory `[NN]-audit-[feature-name].html` planning audit has not been generated and passed.
  - **Action:** Translate the spec into parent/sub-tasks, define Proof Artifacts, and run the mandatory planning audit gate — all as HTML.

- **Phase 3 — Task Implementation**
  - **Condition:** The spec, task list, and a passing planning audit report all exist, and there are incomplete tasks (`data-status="pending"` or `data-status="in-progress"`) in the task list.
  - **Action:** Execute the task list, implement the code, and generate the required Proof Artifacts as HTML.

- **Phase 4 — Validation**
  - **Condition:** Implementation tasks are marked complete and the user asks to validate the work, or the assessor discovers an active spec where implementation appears finished and validation is missing or failing.
  - **Action:** Evaluate the codebase and Proof Artifacts against the spec requirements using strict pass/fail quality gates, reported as HTML.

## Reference Routing

After determining the phase, read only the corresponding bundled reference file:

- **If Phase 1 applies:** `{{skill_dir}}/references/sdd-html-1-generate-spec.md`
- **If Phase 2 applies:** `{{skill_dir}}/references/sdd-html-2-generate-task-list-from-spec.md`
- **If Phase 3 applies:** `{{skill_dir}}/references/sdd-html-3-manage-tasks.md`
- **If Phase 4 applies:** `{{skill_dir}}/references/sdd-html-4-validate-spec-implementation.md`

Always defer to the detailed instructions, constraints, and anti-patterns in the chosen reference file.

## User-Facing Reporting

When beginning or handing off work, state the phase and selected spec in plain language:

```markdown
Current SDD phase: Phase N — <phase name>
Selected spec: `docs/specs/NN-spec-feature/` (or `none yet — new spec will be created` when no spec exists)
Phase reference: `<chosen bundled phase reference path>`
Detected state: <state summary from assessor or manual fallback>
Lavish availability: <available | not available>
Selection reason: <why this spec/phase was selected>
```

When a phase creates or updates files, include:

- created or updated artifact paths (all `.html`)
- whether a Lavish review session was opened, and its outcome if closed during this turn
- any approval gate reached
- validation or audit outcome, when applicable
- a `How to Continue the SDD Workflow` handoff using this exact spacing pattern:

```markdown
## How to Continue the SDD Workflow

Likely next phase action: <what the skill will do next, or that the current workflow is complete>

To continue the workflow in this chat, reply with:

`<suggested natural-language reply>`

You can also continue in a new chat if you want to keep context lean; the SDD skill will reassess repository state from the persisted concrete artifact set for the current phase, such as spec, task, audit, proof, or validation artifacts.
```

Use a blank line between every distinct part of the handoff: heading, likely next phase action, same-chat instruction, suggested reply, and new-chat option. Suggested replies include:

- `Continue SDD with task planning.`
- `Generate the sub-tasks.`
- `Continue SDD with implementation.`
- `Continue SDD with validation.`
- `Start SDD for a new feature.`

Do not label the handoff as a generic next-request section. Do not put the suggested reply before the likely next phase action. Do not compress the handoff into adjacent lines. Do not refer users to slash-style phase invocations, bundled reference filenames, or direct phase command names for continuation. The skill router reassesses workspace state on the next invocation and loads the correct reference.
