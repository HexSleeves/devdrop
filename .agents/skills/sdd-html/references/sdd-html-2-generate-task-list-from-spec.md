# Phase 2 — Generate Task List and Planning Audit (HTML Edition)

This is the HTML-artifact edition of Phase 2 of the SDD workflow. It mirrors the `sdd` skill's Phase 2 methodology exactly — same role, same gates, same checklists, same approval checkpoints — but every deliverable it produces (`[NN]-tasks-[feature-name].html`, `[NN]-audit-[feature-name].html`) is a self-contained HTML document with `data-status` task tracking instead of Markdown checkboxes, and both approval checkpoints offer an optional Lavish Editor review session.

This is an internal phase reference loaded by `SKILL.md`. Users should continue the workflow through natural-language requests to the SDD skill.

## Context Marker

Always begin your response with all active emoji markers, in the order they were introduced.

Format:  "<marker1><marker2><marker3>\n<response>"

The marker for this instruction is:  🌐2️⃣

## You are here in the workflow

You have completed the **spec creation** phase and now need to break down the spec into actionable implementation tasks. This is the critical planning step that bridges requirements to code.

### Workflow Integration

This task list serves as the **execution blueprint** for the entire SDD workflow:

**Value Chain Flow:**

- **Spec → Tasks**: Translates requirements into implementable units
- **Tasks → Planning Audit**: Validates plan quality before implementation
- **Planning Audit → Implementation**: Prevents avoidable planning defects from reaching implementation
- **Implementation → Validation**: Proof artifacts enable verification and evidence collection

**Critical Dependencies:**

- **Parent tasks** become implementation checkpoints in the implementation phase
- **Proof Artifacts** guide implementation verification and become the evidence source for the validation phase
- **Task boundaries** determine git commit points and progress markers
- **Audit findings** determine whether planning is ready for the implementation phase

**What Breaks the Chain:**

- Poorly defined proof artifacts → implementation verification fails
- Missing proof artifacts → validation cannot be completed
- Missing requirement coverage in tasks → spec cannot be fully implemented
- Overly large tasks → loss of incremental progress and demo capability
- Unclear task dependencies → implementation sequence becomes confusing

## Your Role

You are a **Senior Software Engineer and Technical Lead** responsible for translating functional requirements into a structured implementation plan. You must think systematically about the existing codebase, architectural patterns, and deliver a task list that a junior developer can follow successfully.

## Goal

Create a detailed, step-by-step task list in HTML format based on an existing Specification (Spec). Then run a mandatory planning audit checkpoint before implementation handoff. The task list should guide a developer through implementation using **demoable units of work** that provide clear progress indicators.

## Critical Constraints

⚠️ **DO NOT** generate sub-tasks until explicitly requested by the user
⚠️ **DO NOT** begin implementation - this phase is for planning only
⚠️ **DO NOT** create tasks that are too large (multi-day) or too small (single-line changes)
⚠️ **DO NOT** skip the user confirmation step after parent task generation
⚠️ **DO NOT** apply remediation edits until the user explicitly approves the remediation plan
⚠️ **DO NOT** continue to the implementation phase while any REQUIRED audit gate is failing

## Execution Defaults (Positive Directives)

- **ALWAYS** prioritize concise, actionable output over long narrative explanation.
- **ALWAYS** map every functional requirement to at least one task and one planned test artifact.
- **ALWAYS** provide exact file sections for remediation targets.
- **ALWAYS** ask for explicit user confirmation before sub-task generation and before remediation edits.
- **ALWAYS** re-run the audit after approved remediation changes.

## Why Two-Phase Task Generation?

The two-phase approach (parent tasks first, then sub-tasks) serves critical purposes:

1. **Strategic Alignment**: Ensures high-level approach matches user expectations before diving into details
2. **Demoable Focus**: Parent tasks represent end-to-end value that can be demonstrated
3. **Adaptive Planning**: Allows course correction based on feedback before detailed work
4. **Scope Validation**: Confirms the breakdown makes sense before investing in detailed planning

## Spec-to-Task Mapping

Ensure complete spec coverage by:

1. **Trace each user story** to one or more parent tasks
2. **Verify functional requirements** are addressed in specific tasks
3. **Map technical considerations** to implementation details
4. **Identify gaps** where spec requirements aren't covered
5. **Validate acceptance criteria** are testable through proof artifacts
6. **Ensure each functional requirement** has at least one planned test artifact in tasks

## Proof Artifacts

Proof artifacts provide evidence of task completion and are essential for the upcoming validation phase. Each parent task must include artifacts that:

- **Demonstrate functionality** (screenshots, URLs, CLI output)
- **Verify quality** (test results, lint output, performance metrics)
- **Enable validation** (provide evidence for the validation phase)
- **Support troubleshooting** (logs, error messages, configuration states)

**Security Note**: When planning proof artifacts, remember that they will be committed to the repository. Artifacts should use placeholder values for API keys, tokens, and other sensitive data rather than real credentials.

## Evidence Quality Bar (Required)

For each parent task, proof artifacts must satisfy all four checks:

1. **Observable**: demonstrates behavior a reviewer can independently verify.
2. **Reproducible**: includes exact command/path/URL/test reference where applicable.
3. **Scope-linked**: maps to at least one functional requirement and one task section.
4. **Sanitized**: contains no secrets, credentials, or private identifiers.

Reject vague artifact language such as "works as expected" without concrete evidence.

## Private Analysis Process

Before generating any tasks, you must follow this reasoning process:

1. **Spec Analysis**: What are the core functional requirements and user stories?
2. **Current State Assessment**: What existing infrastructure, patterns, and components can we leverage?
3. **Demoable Unit Identification**: What end-to-end vertical slices can be demonstrated?
4. **Dependency Mapping**: What are the logical dependencies between components?
5. **Complexity Evaluation**: Are these tasks appropriately scoped for single implementation cycles?

## Output

- **Format:** HTML (`.html`)
- **Location:** `./docs/specs/[NN]-spec-[feature-name]/` (where `[NN]` is a zero-padded 2-digit number: 01, 02, 03, etc.)
- **Filename:** `[NN]-tasks-[feature-name].html` (e.g., if the Spec is `01-spec-user-profile-editing.html`, save as `01-tasks-user-profile-editing.html`)
- **Audit Filename:** `[NN]-audit-[feature-name].html`

## Process

### Phase 1: Analysis and Planning (Internal)

1. **Receive Spec Reference:** The user points the AI to a specific Spec file in `./docs/specs/`. If the user doesn't provide a spec reference, look for the oldest spec in `./docs/specs/` that doesn't have an accompanying tasks file (i.e., no `[NN]-tasks-[feature-name].html` file in the same directory).
2. **Analyze Spec:** Read and analyze the functional requirements, user stories, and technical constraints
3. **Assess Current State:** Review existing codebase and documentation to understand:
   - Architectural patterns and conventions
   - Existing components that can be leveraged
   - Files that will need modification
   - Testing patterns and infrastructure
   - Contribution patterns and conventions
   - **Repository Standards**: Identify coding standards, build processes, quality gates, and development workflows from project documentation and configuration
4. **Define Demoable Units:** Identify thin, end-to-end vertical slices. Each parent task must be demonstrable.
5. **Evaluate Scope:** Ensure tasks are appropriately sized (not too large, not too small)

### Repository Standards Discovery (Required)

Before task generation or audit, locate and read repository guidance files.

Required search targets (if present):

- `AGENTS.md` (repository root and nearest parent directories)
- `README.md` (repository root and relevant package/application directories)
- `CONTRIBUTING.md`
- `.github/pull_request_template.md`
- lint/format/test policy files (for example: `.pre-commit-config.yaml`, `eslint*`, `pyproject.toml`, `package.json` scripts, CI workflow files)

You MUST NOT infer repository standards from spec/tasks artifacts alone.

### Blocking Checkpoint: Standards Evidence (Required)

Do not proceed to Phase 2 until you produce a standards evidence table with:

- source file path
- read status (`yes`, `not found`, or `access error`)
- 1-3 standards extracted per file when read
- conflicts detected (if any)

### Phase 2: Parent Task Generation

1. **Generate Parent Tasks:** Create the high-level tasks based on your analysis (probably 4-6 tasks, but adjust as needed). Each task must:
   - Represent a demoable unit of work
   - Have clear completion criteria
   - Follow logical dependencies
   - Be implementable in a reasonable timeframe
2. **Save Initial Task List:** Save the parent tasks to `./docs/specs/[NN]-spec-[feature-name]/[NN]-tasks-[feature-name].html` before proceeding
3. **Present for Review**: Present the generated parent tasks to the user for review and wait for their response. If Lavish is available, offer to open the tasks HTML as an interactive Lavish session here — see "Lavish Editor Integration" below.
4. **Wait for Confirmation**: Pause and wait for explicit user approval to generate sub-tasks, such as "Generate the sub-tasks" or an equivalent natural-language confirmation (or an equivalent reply through a Lavish session)

### Phase 3: Sub-Task Generation

Wait for explicit user confirmation before generating sub-tasks. Then:

1. **Identify Relevant Files:** Capture all files that will need creation or modification in an HTML table
2. **Generate Sub-Tasks:** Break down each parent task into smaller, actionable sub-tasks
3. **Update Task List:** Update the existing `./docs/specs/[NN]-spec-[feature-name]/[NN]-tasks-[feature-name].html` file with the sub-tasks and relevant files table sections

### Phase 4: Planning Audit Gate (Required)

After sub-task generation is complete:

1. Create audit report file at `./docs/specs/[NN]-spec-[feature-name]/[NN]-audit-[feature-name].html`.
2. Evaluate and report these gates:
   - **Requirement-to-test traceability (REQUIRED):** Fail if any functional requirement has no planned test artifact mapped in tasks.
   - **Proof artifact verifiability (REQUIRED):** Fail if proof artifact language is vague or not observable.
   - **Repository standards consistency (REQUIRED):** Fail if standards conflict across discovered sources and no precedence/decision is documented. Fail if fewer than 2 repository-guideline sources were read when available. Fail if `AGENTS.md` or root `README.md` exists but was not reviewed.
   - **Open question resolution (REQUIRED):** Fail if material open questions remain unresolved without explicit assumptions.
   - **Regression-risk blind spots (FLAG):** Flag if validation only covers happy-path behavior where regression risk exists.
   - **Non-goal leakage (FLAG):** Flag tasks that exceed goals/non-goals boundaries without justification.
3. Use compact exception-only reporting:
   - Gate overview first, no long narrative
   - At most 3 REQUIRED failures and 2 FLAG findings in the main report
   - Include only exceptions and conflicts; omit empty sections
4. Present findings and remediation items to the user. If Lavish is available, offer to open the audit HTML as an interactive Lavish session here — see "Lavish Editor Integration" below.
5. Wait for explicit user approval before remediation edits.
6. Re-audit after approved remediation edits.
7. Only proceed when all REQUIRED gates pass.

### Lavish Editor Integration (Optional)

`SKILL.md` performs a one-time Lavish availability check before routing to this phase and reports the result as context (`Lavish availability: available | not available`). Reuse that result at the two checkpoints below; do not re-check availability mid-phase.

- **Parent task review checkpoint** (Phase 2 step 3, "Present for Review"): when Lavish is available, after saving `[NN]-tasks-[feature-name].html`, open it as an interactive Lavish session:
  1. `lavish-axi <html-file>` (or `npx -y lavish-axi <html-file>` if no global binary) to open or resume the session in place under `docs/specs/...`.
  2. `lavish-axi poll <html-file>` to wait for the user's review and approval to proceed to sub-tasks. This may run as a background task; never kill it — re-run `poll` if it times out.
  3. Apply `--agent-reply` as needed to answer follow-up questions surfaced in the session.
  4. `lavish-axi end <html-file>` once approval is captured.

  When Lavish is **not** available, fall back to the original chat-based review: present the parent tasks in the chat and wait for a natural-language confirmation.

- **Audit findings checkpoint** (Phase 4 step 4, "Present findings and remediation items"): same open/poll/reply/end pattern against `[NN]-audit-[feature-name].html` when Lavish is available, so the user can annotate specific gate rows or findings directly. Fall back to chat-based review when Lavish is not available.

- Lavish integration is always optional and never replaces the underlying approval gate. The gate is satisfied by explicit user approval — via a Lavish reply or a chat message — not by merely opening a session. Never skip the "wait for confirmation" step in Phase 2 or the "wait for explicit approval before remediation edits" step in Phase 4.

### Phase 4A: Chain-of-Verification Check (Required Before Handoff)

Before continuing to the implementation phase, run this verification loop:

1. **Initial assessment:** complete the audit and draft findings.
2. **Self-questioning:** ask "Do all REQUIRED gates pass with explicit evidence?"
3. **Fact-checking:** verify each finding against spec, task file, and repository standards sources.
4. **Inconsistency resolution:** correct any finding that is unsupported or ambiguous.
5. **Final synthesis:** publish the final audit status and next action.

### Failure Handling

If you cannot evaluate a REQUIRED gate due to missing artifacts or unclear standards:

1. Mark gate as `FAIL` with reason `insufficient evidence`.
2. Add one concrete remediation item that resolves the evidence gap.
3. Request user clarification only when the missing evidence cannot be derived from repository artifacts.

If repository guideline files are missing or unreadable:

1. Record exact file paths searched and result (`not found` or `access error`).
2. Use fallback evidence from repository configuration and CI workflow files.
3. Mark standards confidence as low and add a remediation item for missing standards documentation.

## Task/Audit HTML Template Requirements

A bundled Python assessor script (`scripts/assess-sdd-state.py`) regex-matches literal strings in the generated HTML to determine workflow state. Follow these rules precisely in every generated artifact:

1. **Task / sub-task items:** each task or sub-task is an element carrying `data-status="pending"`, `data-status="in-progress"`, or `data-status="done"`, e.g. `<li class="task" data-task-id="1.1" data-status="pending"><span class="task-label">1.1 Sub-task description</span></li>`. The visible checkbox glyph (☐/◐/☑) renders automatically from the shared CSS `::before` rule keyed on `data-status` — do not also hardcode the glyph into the text content. Parent tasks (`1.0`, `2.0`, etc.) are their own `<section class="parent-task" data-task-id="1.0" data-status="...">` wrapping a heading, the Proof Artifact(s) sub-section, and the Tasks sub-section.
2. **Parent tasks not yet broken into sub-tasks** (Phase 2, parents-only): wrap the "Tasks" sub-section in an element with `data-tbd="true"` AND literal visible text `TBD`, e.g. `<div class="subtasks" data-tbd="true"><h4>1.0 Tasks</h4><p class="tbd-marker">TBD</p></div>`.
3. **Audit overall status:** the `<body>` tag must carry `data-overall-status="PASS"` or `data-overall-status="FAIL"`, AND the visible Executive Summary badge must contain the literal uppercase word `PASS` or `FAIL` (not "Passed"/"Failed"), e.g. `<span class="badge status-fail" data-role="overall-status">FAIL</span>`.
4. **Per-gate status in the Gateboard table:** use the literal uppercase words `PASS`, `FAIL`, or `FLAG` inside a badge, e.g. `<td><span class="badge status-fail">FAIL</span></td>` — never lowercase or past-tense for gates specifically. This applies to every row of the Gateboard table (Requirement-to-test traceability, Proof artifact verifiability, Repository standards consistency, Open question resolution, Regression-risk blind spots, Non-goal leakage).

Every generated artifact is a single self-contained HTML5 file: inline `<style>`, no external assets except relative-path local images placed alongside the HTML file. It must render correctly when opened directly in a browser with no server. Use the shared shell/style block established for this skill (only change `<title>`, the badge-type text/label, and the body content) so all sdd-html artifacts look consistent — extend `<style>` with extra rules a doc type needs, but never remove or rename the shared classes.

## Phase 2 Output Format (Parent Tasks Only)

When generating parent tasks in Phase 2, use this hierarchical structure with the Tasks sub-section marked as TBD (`data-tbd="true"` plus literal `TBD` text) on every parent task, and save it as `[NN]-tasks-[feature-name].html`:

```html
<!doctype html>
<html lang="en">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1" />
<title>[NN]-tasks-[feature-name]</title>
<style>
  :root {
    color-scheme: light dark;
    --bg: #ffffff; --fg: #1a1a2e; --muted: #6b7280; --border: #e5e7eb;
    --card: #f9fafb; --accent: #4f46e5; --pass: #059669; --fail: #dc2626; --warn: #d97706;
  }
  @media (prefers-color-scheme: dark) {
    :root { --bg:#0f1117; --fg:#e6e6ef; --muted:#9ca3af; --border:#262a35; --card:#161922; --accent:#818cf8; --pass:#34d399; --fail:#f87171; --warn:#fbbf24; }
  }
  * { box-sizing: border-box; }
  body { margin:0; background:var(--bg); color:var(--fg); font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; line-height:1.6; }
  .doc { max-width: 960px; margin: 0 auto; padding: 2rem 1.5rem 4rem; }
  header.doc-header { border-bottom: 1px solid var(--border); padding-bottom:1rem; margin-bottom:2rem; }
  h1 { margin: .25rem 0; font-size:1.6rem; }
  h2 { margin-top:2rem; border-bottom:1px solid var(--border); padding-bottom:.35rem; }
  .meta { color:var(--muted); font-size:.9rem; }
  .badge { display:inline-block; font-size:.72rem; font-weight:700; letter-spacing:.05em; text-transform:uppercase; padding:.25rem .6rem; border-radius:999px; }
  .badge-type { background:var(--accent); color:#fff; }
  .status-pass, .status-done, .status-verified { background:var(--pass); color:#fff; }
  .status-fail, .status-failed { background:var(--fail); color:#fff; }
  .status-flag, .status-pending, .status-unknown, .status-in-progress { background:var(--warn); color:#111; }
  table { width:100%; border-collapse: collapse; margin: 1rem 0; }
  th, td { border:1px solid var(--border); padding:.5rem .75rem; text-align:left; vertical-align:top; }
  th { background:var(--card); }
  .card { background:var(--card); border:1px solid var(--border); border-radius:.5rem; padding:1rem 1.25rem; margin:1rem 0; }
  code, pre { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; }
  pre { background:var(--card); border:1px solid var(--border); border-radius:.5rem; padding:1rem; overflow-x:auto; }
  .task-list { list-style:none; padding-left:0; }
  .task { padding:.35rem 0; }
  .task[data-status="done"] .task-label::before { content: "\2611 "; }
  .task[data-status="pending"] .task-label::before { content: "\2610 "; }
  .task[data-status="in-progress"] .task-label::before { content: "\25D1 "; }
  figure { margin: 1rem 0; }
  figure img { max-width:100%; border:1px solid var(--border); border-radius:.5rem; }
  figcaption { color:var(--muted); font-size:.85rem; margin-top:.35rem; }
  .parent-task { border:1px solid var(--border); border-radius:.5rem; padding:1rem 1.25rem; margin:1.25rem 0; }
  .parent-task > h3 { margin-top:0; }
  .proof-artifacts ul { margin:.5rem 0; }
  .tbd-marker { color:var(--muted); font-style:italic; }
</style>
</head>
<body>
<div class="doc">
  <header class="doc-header">
    <span class="badge badge-type">TASKS</span>
    <h1>[NN]-tasks-[feature-name]</h1>
    <p class="meta">Sequence [NN] · Feature: [feature-name] · Phase 2 · Parent tasks (sub-tasks pending)</p>
  </header>

  <section id="tasks">
    <h2>Tasks</h2>

    <section class="parent-task" data-task-id="1.0" data-status="pending">
      <h3>1.0 Parent Task Title</h3>

      <div class="proof-artifacts">
        <h4>1.0 Proof Artifact(s)</h4>
        <ul>
          <li>Screenshot: <code>/path</code> page showing completed X flow demonstrates end-to-end functionality</li>
          <li>URL: https://... demonstrates feature is accessible</li>
          <li>CLI: <code>command --flag</code> returns expected output demonstrates feature works</li>
          <li>Test: <code>MyFeature.test.ts</code> passes demonstrates requirement implementation</li>
        </ul>
      </div>

      <div class="subtasks" data-tbd="true">
        <h4>1.0 Tasks</h4>
        <p class="tbd-marker">TBD</p>
      </div>
    </section>

    <section class="parent-task" data-task-id="2.0" data-status="pending">
      <h3>2.0 Parent Task Title</h3>

      <div class="proof-artifacts">
        <h4>2.0 Proof Artifact(s)</h4>
        <ul>
          <li>Screenshot: User flow showing Z with persisted state demonstrates feature persistence</li>
          <li>Test: <code>UserFlow.test.ts</code> passes demonstrates state management works</li>
        </ul>
      </div>

      <div class="subtasks" data-tbd="true">
        <h4>2.0 Tasks</h4>
        <p class="tbd-marker">TBD</p>
      </div>
    </section>

    <section class="parent-task" data-task-id="3.0" data-status="pending">
      <h3>3.0 Parent Task Title</h3>

      <div class="proof-artifacts">
        <h4>3.0 Proof Artifact(s)</h4>
        <ul>
          <li>CLI: <code>config get ...</code> returns expected value demonstrates configuration is verifiable</li>
          <li>Log: Configuration loaded message demonstrates system initialization</li>
          <li>Diff: Configuration file changes demonstrates setup completion</li>
        </ul>
      </div>

      <div class="subtasks" data-tbd="true">
        <h4>3.0 Tasks</h4>
        <p class="tbd-marker">TBD</p>
      </div>
    </section>
  </section>
</div>
</body>
</html>
```

## Phase 3 Output Format (Complete with Sub-Tasks)

After user confirmation in Phase 3, update the same file with this complete structure: a Relevant Files table, a Notes section, and the same parent tasks with their `data-tbd="true"` placeholders replaced by real `data-status="pending"` sub-task items.

```html
<!doctype html>
<html lang="en">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1" />
<title>[NN]-tasks-[feature-name]</title>
<style>
  :root {
    color-scheme: light dark;
    --bg: #ffffff; --fg: #1a1a2e; --muted: #6b7280; --border: #e5e7eb;
    --card: #f9fafb; --accent: #4f46e5; --pass: #059669; --fail: #dc2626; --warn: #d97706;
  }
  @media (prefers-color-scheme: dark) {
    :root { --bg:#0f1117; --fg:#e6e6ef; --muted:#9ca3af; --border:#262a35; --card:#161922; --accent:#818cf8; --pass:#34d399; --fail:#f87171; --warn:#fbbf24; }
  }
  * { box-sizing: border-box; }
  body { margin:0; background:var(--bg); color:var(--fg); font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; line-height:1.6; }
  .doc { max-width: 960px; margin: 0 auto; padding: 2rem 1.5rem 4rem; }
  header.doc-header { border-bottom: 1px solid var(--border); padding-bottom:1rem; margin-bottom:2rem; }
  h1 { margin: .25rem 0; font-size:1.6rem; }
  h2 { margin-top:2rem; border-bottom:1px solid var(--border); padding-bottom:.35rem; }
  .meta { color:var(--muted); font-size:.9rem; }
  .badge { display:inline-block; font-size:.72rem; font-weight:700; letter-spacing:.05em; text-transform:uppercase; padding:.25rem .6rem; border-radius:999px; }
  .badge-type { background:var(--accent); color:#fff; }
  .status-pass, .status-done, .status-verified { background:var(--pass); color:#fff; }
  .status-fail, .status-failed { background:var(--fail); color:#fff; }
  .status-flag, .status-pending, .status-unknown, .status-in-progress { background:var(--warn); color:#111; }
  table { width:100%; border-collapse: collapse; margin: 1rem 0; }
  th, td { border:1px solid var(--border); padding:.5rem .75rem; text-align:left; vertical-align:top; }
  th { background:var(--card); }
  .card { background:var(--card); border:1px solid var(--border); border-radius:.5rem; padding:1rem 1.25rem; margin:1rem 0; }
  code, pre { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; }
  pre { background:var(--card); border:1px solid var(--border); border-radius:.5rem; padding:1rem; overflow-x:auto; }
  .task-list { list-style:none; padding-left:0; }
  .task { padding:.35rem 0; }
  .task[data-status="done"] .task-label::before { content: "\2611 "; }
  .task[data-status="pending"] .task-label::before { content: "\2610 "; }
  .task[data-status="in-progress"] .task-label::before { content: "\25D1 "; }
  figure { margin: 1rem 0; }
  figure img { max-width:100%; border:1px solid var(--border); border-radius:.5rem; }
  figcaption { color:var(--muted); font-size:.85rem; margin-top:.35rem; }
  .parent-task { border:1px solid var(--border); border-radius:.5rem; padding:1rem 1.25rem; margin:1.25rem 0; }
  .parent-task > h3 { margin-top:0; }
  .proof-artifacts ul { margin:.5rem 0; }
</style>
</head>
<body>
<div class="doc">
  <header class="doc-header">
    <span class="badge badge-type">TASKS</span>
    <h1>[NN]-tasks-[feature-name]</h1>
    <p class="meta">Sequence [NN] · Feature: [feature-name] · Phase 2 · Sub-tasks generated</p>
  </header>

  <section id="relevant-files">
    <h2>Relevant Files</h2>
    <table>
      <thead>
        <tr><th>File</th><th>Why It Is Relevant</th></tr>
      </thead>
      <tbody>
        <tr><td><code>path/to/potential/file1.ts</code></td><td>Contains the main component or implementation entry point for this feature.</td></tr>
        <tr><td><code>path/to/file1.test.ts</code></td><td>Unit tests for <code>file1.ts</code>.</td></tr>
        <tr><td><code>path/to/another/file.tsx</code></td><td>API route handler or UI entry point for data submission.</td></tr>
        <tr><td><code>path/to/another/file.test.tsx</code></td><td>Unit tests for <code>another/file.tsx</code>.</td></tr>
        <tr><td><code>lib/utils/helpers.ts</code></td><td>Utility functions needed for calculations or shared behavior.</td></tr>
        <tr><td><code>lib/utils/helpers.test.ts</code></td><td>Unit tests for <code>helpers.ts</code>.</td></tr>
      </tbody>
    </table>

    <h3>Notes</h3>
    <ul>
      <li>Unit tests should typically be placed alongside the code files they are testing (e.g., <code>MyComponent.tsx</code> and <code>MyComponent.test.tsx</code> in the same directory).</li>
      <li>Use the repository's established testing command and patterns (e.g., <code>npx jest [optional/path/to/test/file]</code>, <code>pytest [path]</code>, <code>cargo test</code>, etc.).</li>
      <li>Follow the repository's existing code organization, naming conventions, and style guidelines.</li>
      <li>Adhere to identified quality gates and pre-commit hooks.</li>
    </ul>
  </section>

  <section id="tasks">
    <h2>Tasks</h2>

    <section class="parent-task" data-task-id="1.0" data-status="pending">
      <h3>1.0 Parent Task Title</h3>

      <div class="proof-artifacts">
        <h4>1.0 Proof Artifact(s)</h4>
        <ul>
          <li>Screenshot: <code>/path</code> page showing completed X flow demonstrates end-to-end functionality</li>
          <li>URL: https://... demonstrates feature is accessible</li>
          <li>CLI: <code>command --flag</code> returns expected output demonstrates feature works</li>
          <li>Test: <code>MyFeature.test.ts</code> passes demonstrates requirement implementation</li>
        </ul>
      </div>

      <div class="subtasks">
        <h4>1.0 Tasks</h4>
        <ul class="task-list">
          <li class="task" data-task-id="1.1" data-status="pending"><span class="task-label">1.1 [Sub-task description 1.1]</span></li>
          <li class="task" data-task-id="1.2" data-status="pending"><span class="task-label">1.2 [Sub-task description 1.2]</span></li>
        </ul>
      </div>
    </section>

    <section class="parent-task" data-task-id="2.0" data-status="pending">
      <h3>2.0 Parent Task Title</h3>

      <div class="proof-artifacts">
        <h4>2.0 Proof Artifact(s)</h4>
        <ul>
          <li>Screenshot: User flow showing Z with persisted state demonstrates feature persistence</li>
          <li>Test: <code>UserFlow.test.ts</code> passes demonstrates state management works</li>
        </ul>
      </div>

      <div class="subtasks">
        <h4>2.0 Tasks</h4>
        <ul class="task-list">
          <li class="task" data-task-id="2.1" data-status="pending"><span class="task-label">2.1 [Sub-task description 2.1]</span></li>
          <li class="task" data-task-id="2.2" data-status="pending"><span class="task-label">2.2 [Sub-task description 2.2]</span></li>
        </ul>
      </div>
    </section>

    <section class="parent-task" data-task-id="3.0" data-status="pending">
      <h3>3.0 Parent Task Title</h3>

      <div class="proof-artifacts">
        <h4>3.0 Proof Artifact(s)</h4>
        <ul>
          <li>CLI: <code>config get ...</code> returns expected value demonstrates configuration is verifiable</li>
          <li>Log: Configuration loaded message demonstrates system initialization</li>
          <li>Diff: Configuration file changes demonstrates setup completion</li>
        </ul>
      </div>

      <div class="subtasks">
        <h4>3.0 Tasks</h4>
        <ul class="task-list">
          <li class="task" data-task-id="3.1" data-status="pending"><span class="task-label">3.1 [Sub-task description 3.1]</span></li>
          <li class="task" data-task-id="3.2" data-status="pending"><span class="task-label">3.2 [Sub-task description 3.2]</span></li>
        </ul>
      </div>
    </section>
  </section>
</div>
</body>
</html>
```

## Audit Report Format (Phase 4 and Later)

Use this structure in `[NN]-audit-[feature-name].html`:

```html
<!doctype html>
<html lang="en">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1" />
<title>[NN]-audit-[feature-name]</title>
<style>
  :root {
    color-scheme: light dark;
    --bg: #ffffff; --fg: #1a1a2e; --muted: #6b7280; --border: #e5e7eb;
    --card: #f9fafb; --accent: #4f46e5; --pass: #059669; --fail: #dc2626; --warn: #d97706;
  }
  @media (prefers-color-scheme: dark) {
    :root { --bg:#0f1117; --fg:#e6e6ef; --muted:#9ca3af; --border:#262a35; --card:#161922; --accent:#818cf8; --pass:#34d399; --fail:#f87171; --warn:#fbbf24; }
  }
  * { box-sizing: border-box; }
  body { margin:0; background:var(--bg); color:var(--fg); font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; line-height:1.6; }
  .doc { max-width: 960px; margin: 0 auto; padding: 2rem 1.5rem 4rem; }
  header.doc-header { border-bottom: 1px solid var(--border); padding-bottom:1rem; margin-bottom:2rem; }
  h1 { margin: .25rem 0; font-size:1.6rem; }
  h2 { margin-top:2rem; border-bottom:1px solid var(--border); padding-bottom:.35rem; }
  .meta { color:var(--muted); font-size:.9rem; }
  .badge { display:inline-block; font-size:.72rem; font-weight:700; letter-spacing:.05em; text-transform:uppercase; padding:.25rem .6rem; border-radius:999px; }
  .badge-type { background:var(--accent); color:#fff; }
  .status-pass, .status-done, .status-verified { background:var(--pass); color:#fff; }
  .status-fail, .status-failed { background:var(--fail); color:#fff; }
  .status-flag, .status-pending, .status-unknown, .status-in-progress { background:var(--warn); color:#111; }
  table { width:100%; border-collapse: collapse; margin: 1rem 0; }
  th, td { border:1px solid var(--border); padding:.5rem .75rem; text-align:left; vertical-align:top; }
  th { background:var(--card); }
  .card { background:var(--card); border:1px solid var(--border); border-radius:.5rem; padding:1rem 1.25rem; margin:1rem 0; }
  code, pre { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; }
  pre { background:var(--card); border:1px solid var(--border); border-radius:.5rem; padding:1rem; overflow-x:auto; }
  .task-list { list-style:none; padding-left:0; }
  .task { padding:.35rem 0; }
  .task[data-status="done"] .task-label::before { content: "\2611 "; }
  .task[data-status="pending"] .task-label::before { content: "\2610 "; }
  .task[data-status="in-progress"] .task-label::before { content: "\25D1 "; }
  figure { margin: 1rem 0; }
  figure img { max-width:100%; border:1px solid var(--border); border-radius:.5rem; }
  figcaption { color:var(--muted); font-size:.85rem; margin-top:.35rem; }
  .gate-row td:nth-child(2) { white-space:nowrap; }
  .finding { border-left: 3px solid var(--warn); padding-left:.75rem; margin:.75rem 0; }
  .finding-required { border-left-color: var(--fail); }
  .finding-flag { border-left-color: var(--warn); }
  .remediation-status { font-weight:700; }
</style>
</head>
<body data-overall-status="FAIL">
<div class="doc">
  <header class="doc-header">
    <span class="badge badge-type">AUDIT</span>
    <h1>[NN]-audit-[feature-name]</h1>
    <p class="meta">Sequence [NN] · Feature: [feature-name] · Phase 4 · Planning Audit · Run 1</p>
  </header>

  <section id="executive-summary" class="card">
    <h2>Executive Summary</h2>
    <p>Overall Status: <span class="badge status-fail" data-role="overall-status">FAIL</span></p>
    <p>Required Gate Failures: <strong>1</strong></p>
    <p>Flagged Risks: <strong>1</strong></p>
  </section>

  <section id="gateboard">
    <h2>Gateboard</h2>
    <table>
      <thead>
        <tr><th>Gate</th><th>Status</th><th>Why it failed (&lt;=10 words)</th><th>Exact fix target</th></tr>
      </thead>
      <tbody>
        <tr class="gate-row"><td>Requirement-to-test traceability</td><td><span class="badge status-fail">FAIL</span></td><td>FR-2 has no mapped test artifact</td><td><code>#tasks section 2.0</code></td></tr>
        <tr class="gate-row"><td>Proof artifact verifiability</td><td><span class="badge status-pass">PASS</span></td><td>-</td><td>-</td></tr>
        <tr class="gate-row"><td>Repository standards consistency</td><td><span class="badge status-pass">PASS</span></td><td>-</td><td>-</td></tr>
        <tr class="gate-row"><td>Open question resolution</td><td><span class="badge status-pass">PASS</span></td><td>-</td><td>-</td></tr>
        <tr class="gate-row"><td>Regression-risk blind spots</td><td><span class="badge status-flag">FLAG</span></td><td>Only happy-path coverage for task 2.0</td><td><code>#tasks section 2.0</code></td></tr>
        <tr class="gate-row"><td>Non-goal leakage</td><td><span class="badge status-pass">PASS</span></td><td>-</td><td>-</td></tr>
      </tbody>
    </table>
  </section>

  <section id="standards-evidence">
    <h2>Standards Evidence Table (Required)</h2>
    <table>
      <thead>
        <tr><th>Source File</th><th>Read</th><th>Standards Extracted</th><th>Conflicts</th></tr>
      </thead>
      <tbody>
        <tr><td><code>AGENTS.md</code></td><td>yes</td><td>Follow repository workflow conventions</td><td>none</td></tr>
        <tr><td><code>README.md</code></td><td>yes</td><td>Use documented workflow order and artifact paths</td><td>none</td></tr>
      </tbody>
    </table>
  </section>

  <section id="findings">
    <h2>Findings</h2>

    <h3>REQUIRED Failures (max 3 in main report)</h3>
    <div class="finding finding-required" data-finding-id="req-1">
      <p><strong>1. [Issue]</strong></p>
      <p>Missing item: </p>
      <p>File section to edit: </p>
      <p>Acceptance condition: </p>
    </div>

    <h3>FLAG Findings (max 2 in main report)</h3>
    <div class="finding finding-flag" data-finding-id="flag-1">
      <p><strong>1. [Issue]</strong></p>
      <p>Risk: </p>
      <p>Suggested remediation: </p>
    </div>
  </section>

  <section id="remediation-plan">
    <h2>User-Approved Remediation Plan</h2>
    <p class="remediation-status" data-role="remediation-status">Pending approval</p>
  </section>

  <section id="reaudit-delta">
    <h2>Re-Audit Delta (Runs 2+ only)</h2>
    <ul>
      <li>Changed gate statuses since previous run (only changed items): </li>
      <li>Still-failing REQUIRED gates: </li>
      <li>Newly introduced findings (if any): </li>
    </ul>
  </section>
</div>
</body>
</html>
```

If all REQUIRED gates pass on the first audit run, keep the report minimal: include only the Executive Summary section (with `data-overall-status="PASS"` on `<body>` and a `status-pass` `PASS` badge) and the Gateboard section; omit the empty Findings, User-Approved Remediation Plan, and Re-Audit Delta sections entirely rather than leaving them present-but-empty.

## Interaction Model

**Critical:** This process includes explicit approval checkpoints:

1. **Phase 1 Completion:** After generating parent tasks, you must stop and present them for review
2. **Explicit Confirmation:** Only proceed to sub-tasks after explicit user approval, such as "Generate the sub-tasks" or an equivalent natural-language confirmation
3. **Audit Review:** After generating the audit report, you must present findings and wait for approval before remediation edits
4. **No Auto-progression:** Never continue to the implementation phase while REQUIRED audit gates fail

**Example interaction:**
> "I have analyzed the spec and generated [X] parent tasks that represent demoable units of work. Each task includes proof artifacts that demonstrate what will be shown. Please review these high-level tasks and confirm if you'd like me to proceed with generating detailed sub-tasks. Confirm in natural language, such as 'Generate the sub-tasks', to continue."

## Target Audience

Write tasks and sub-tasks for a **junior developer** who:

- Understands the programming language and framework
- Is familiar with the existing codebase structure
- Needs clear, actionable steps without ambiguity
- Will be implementing tasks independently
- Relies on proof artifacts to verify completion
- Must follow established repository patterns and conventions

## Quality Checklist

Before finalizing your task list, verify:

- [ ] Each parent task is demoable and has clear completion criteria
- [ ] Proof Artifacts are specific and demonstrate clear functionality
- [ ] Proof Artifacts are appropriate for each task
- [ ] Tasks are appropriately scoped (not too large/small)
- [ ] Dependencies are logical and sequential
- [ ] Sub-tasks are actionable and unambiguous
- [ ] Relevant files table is comprehensive, accurate, and easy to scan
- [ ] Format follows the exact structure specified above
- [ ] Repository standards and patterns are identified and incorporated
- [ ] Implementation will follow established coding conventions and workflows
- [ ] Every functional requirement maps to planned test artifacts
- [ ] Audit report exists and is current
- [ ] REQUIRED audit gates are passing
- [ ] Any remediation edits were explicitly user-approved

## How to Continue the SDD Workflow

Likely next phase action: the skill will route to Phase 3 and execute the task list using TDD, implement the planned changes, and generate the required proof artifacts.

To continue the workflow in this chat, reply with:

`Continue SDD with implementation.`

You can also continue in a new chat if you want to keep context lean; the SDD skill will reassess the repository state from the persisted spec/task/audit artifacts.

## Final Instructions

1. Follow the Private Analysis Process before generating any tasks
2. Assess current codebase for existing patterns and reusable components
3. Generate high-level tasks that represent demoable units of work (adjust count based on spec complexity) and save them to `./docs/specs/[NN]-spec-[feature-name]/[NN]-tasks-[feature-name].html`
4. **CRITICAL**: Stop after generating parent tasks and wait for explicit natural-language confirmation before proceeding to sub-task generation.
5. Ensure every parent task has specific Proof Artifacts that demonstrate what will be shown
6. Identify all relevant files for creation/modification and present them in the required HTML table format
7. Run the planning audit gate and create `[NN]-audit-[feature-name].html`
8. Present findings and remediation plan; wait for explicit approval before remediation edits
9. Run the Chain-of-Verification check before handoff decisions
10. Re-audit until all REQUIRED gates pass
11. Guide user to the next workflow step (the implementation phase) only when audit is passing
12. Stop working once user confirms task list is complete
