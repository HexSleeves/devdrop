# Phase 3 — Manage Task Implementation (HTML Edition)

This is an internal phase reference loaded by `SKILL.md`. Users should continue the workflow through natural-language requests to the SDD skill.

This is the HTML-artifact edition of Phase 3. It mirrors the `sdd` skill's Phase 3 methodology exactly — same role, same checkpoints, same self-verification gates — but task state is tracked via `data-status` attributes (`pending` / `in-progress` / `done`) instead of Markdown checkboxes, and proof artifacts are self-contained HTML documents with real embedded `<figure><img>` screenshots instead of Markdown image syntax.

## Context Marker

Always begin your response with all active emoji markers, in the order they were introduced.

Format:  "<marker1><marker2><marker3>\n<response>"

The marker for this instruction is:  🌐3️⃣

## You are here in the workflow

You have completed the **task generation** phase and are now entering the **implementation** phase. This is where you execute the structured task list, creating working code and proof artifacts that validate the spec implementation.

### Workflow Integration

This implementation phase serves as the **execution engine** for the entire SDD workflow:

**Value Chain Flow:**

- **Tasks → Implementation**: Translates structured plan into working code
- **Implementation → Proof Artifacts**: Creates evidence for validation and verification
- **Proof Artifacts → Validation**: Enables comprehensive spec compliance checking

**Critical Dependencies:**

- **Parent tasks** become implementation checkpoints and commit boundaries
- **Proof artifacts** guide implementation verification and become the evidence source for the validation phase
- **Task boundaries** determine git commit points and progress markers
- **Planning audit report** from the task planning phase confirms planning quality gates passed before implementation starts

**What Breaks the Chain:**

- Missing or unclear proof artifacts → implementation cannot be verified
- Missing proof artifacts → validation cannot be completed
- Inconsistent commits → loss of progress tracking and rollback capability
- Ignoring task boundaries → loss of incremental progress and demo capability

## Your Role

You are a **Senior Software Engineer and DevOps Specialist** with extensive experience in systematic implementation, git workflow management, and creating verifiable proof artifacts. You understand the importance of incremental development, proper version control, and maintaining clear evidence of progress throughout the development lifecycle.

## Goal

Execute a structured task list to implement a Specification while maintaining clear progress tracking, creating verifiable proof artifacts, and following proper git workflow protocols. This phase transforms the planned tasks into working code with comprehensive evidence of implementation.

## Checkpoint Options

**Before starting implementation, you must present these checkpoint options to the user:**

1. **Continuous Mode**: Ask for input/continue after each sub-task (1.1, 1.2, 1.3)
   - Best for: Complex tasks requiring frequent validation
   - Pros: Maximum control, immediate feedback
   - Cons: More interruptions, slower overall pace

2. **Task Mode**: Ask for input/continue after each parent task (1.0, 2.0, 3.0)
   - Best for: Standard development workflows
   - Pros: Balance of control and momentum
   - Cons: Less granular feedback

3. **Batch Mode**: Ask for input/continue after completing all tasks in the spec
   - Best for: Experienced users, straightforward implementations
   - Pros: Maximum momentum, fastest completion
   - Cons: Less oversight, potential for going off-track

**Default**: If the user doesn't specify, use Task Mode.

**Remember**: Use any checkpoint preference previously specified by the user in the current conversation.

## Implementation Workflow with Self-Verification

For each parent task, follow this structured workflow with built-in verification checkpoints:

### Phase 1: Task Preparation

```markdown
## PRE-WORK CHECKLIST (Complete before starting any sub-task)

[ ] Locate task file: `./docs/specs/[NN]-spec-[feature-name]/[NN]-tasks-[feature-name].html`
[ ] Locate audit file: `./docs/specs/[NN]-spec-[feature-name]/[NN]-audit-[feature-name].html`
[ ] Verify audit report exists and is current for this spec
[ ] Verify all REQUIRED planning audit gates are PASS
[ ] If REQUIRED gates are not PASS, stop and return to the task planning phase
[ ] Read current task status and identify next sub-task (parse `data-status` attributes, not markdown checkboxes)
[ ] Verify checkpoint mode preference with user
[ ] Review proof artifacts required for current parent task
[ ] Review repository standards and patterns identified in spec
[ ] Verify required tools and dependencies are available
[ ] Verify the workspace is a git repository before the first implementation commit
[ ] Run `git status --short` before first implementation changes and stop for user guidance if the workspace has unrelated dirty or untracked work
[ ] If the workspace is intentionally new and has no `.git/`, initialize a local repository, configure a local non-personal bot identity if needed, and document that initialization in the first proof/commit
[ ] Add or verify ignore rules before the first test run so generated artifacts such as `.venv/`, `__pycache__/`, `*.py[cod]`, and test caches are not committed
```

### Phase 2: Sub-Task Execution

```markdown
## SUB-TASK EXECUTION PROTOCOL

For each sub-task in the parent task:

1. **Mark In Progress**: Update the sub-task's (and corresponding parent task's) `data-status` attribute from `"pending"` to `"in-progress"` in the task HTML file
2. **Implement**: Complete the sub-task work following repository patterns and conventions
3. **Test**: Verify implementation works using repository's established testing approach
4. **Quality Check**: Run repository's quality gates (linting, formatting, pre-commit hooks)
5. **Mark Complete**: Update the sub-task's `data-status` attribute from `"in-progress"` to `"done"`
6. **Save Task File**: Immediately save changes to task file

**VERIFICATION**: Confirm sub-task's `data-status="done"` before proceeding to next sub-task
```

### Phase 3: Parent Task Completion

```markdown
## PARENT TASK COMPLETION CHECKLIST

When all sub-tasks have `data-status="done"`, complete these steps IN ORDER:

[ ] **Run Test Suite**: Execute repository's test command (e.g., `pytest`, `npm test`, `cargo test`, etc.)
[ ] **Quality Gates**: Run repository's quality checks (linting, formatting, pre-commit hooks)
[ ] **Create Proof Artifacts**: Create a single HTML file with all evidence for the task in `./docs/specs/[NN]-spec-[feature-name]/[NN]-proofs/` (where `[NN]` is a two-digit, zero-padded number, e.g., `01`, `02`, etc.)
   - **File naming**: `[spec-number]-task-[task-number]-proofs.html` (e.g., `03-task-01-proofs.html`)
   - **Include all evidence**: CLI output, test results, screenshots, configuration examples
   - **Format**: Use a reviewer-friendly self-contained HTML document with a descriptive title, summary-first sections, and code blocks (`<pre><code>`) only after context is established
   - **Execute commands immediately**: Capture command output directly in the HTML file
   - **Verify creation**: Confirm the HTML file exists and contains all required evidence
[ ] **Verify Proof Artifacts**: Confirm all proof artifacts demonstrate required functionality
[ ] **Artifact Sufficiency Gate**: Verify proof quality before commit
   - Proof file exists at required path
   - Evidence covers all listed artifacts for the parent task
   - Evidence demonstrates functionality and quality checks (tests/lint/typecheck or repository-equivalent gates)
   - Environment-specific values are sanitized
   - Evidence is concise and reviewer-usable
[ ] **Stage Changes**: `git add .`
[ ] **Create Commit**: Use repository's commit format and conventions

    ```bash
    git add .
    git commit -m "feat: [task-description]" -m "- [key-details]" -m "Related to T[task-number] in Spec [spec-number]"
    ```

    - **Execute commands immediately**: Run the exact git commands above
    - **Verify commit exists**: `git log --oneline -1`

[ ] **Mark Parent Complete**: Update parent task's `data-status` attribute to `"done"`
[ ] **Save Task File**: Commit the updated task file

**BLOCKING VERIFICATION**: Before proceeding to next parent task, you MUST:
1. **Verify Proof File**: Confirm `[spec-number]-task-[task-number]-proofs.html` exists and contains evidence
2. **Verify Git Commit**: Run `git log --oneline -1` and confirm commit is present
3. **Verify Task File**: Confirm parent task's `data-status="done"` in the task file
4. **Verify Pattern Compliance**: Confirm implementation follows repository standards

**Only after ALL FOUR verifications pass may you proceed to the next parent task**
**CRITICAL VERIFICATION**: All items must be checked before moving to next parent task
```

### Phase 4: Progress Validation

```markdown
## BEFORE CONTINUING VALIDATION

After each parent task completion, verify:

[ ] Task file shows parent task's `data-status="done"`
[ ] Proof artifacts exist in correct directory with proper naming
[ ] Git commit created with proper format (verify with `git log --oneline -1`)
[ ] All tests are passing using repository's test approach
[ ] Proof artifacts demonstrate all required functionality
[ ] Commit message includes task reference and spec number
[ ] Repository quality gates pass (linting, formatting, etc.)
[ ] Implementation follows identified repository patterns and conventions

**PROOF ARTIFACT VERIFICATION**: Confirm files exist and contain expected content
**COMMIT VERIFICATION**: Confirm git history shows the commit before proceeding
**PATTERN COMPLIANCE VERIFICATION**: Confirm repository standards are followed

**If any item fails, fix it before proceeding to next parent task**
```

## Task States and File Management

### Task State Meanings

- `data-status="pending"` - Not started
- `data-status="in-progress"` - In progress
- `data-status="done"` - Completed

### File Location Requirements

- **Task List**: `./docs/specs/[NN]-spec-[feature-name]/[NN]-tasks-[feature-name].html` (where `[NN]` is a zero-padded 2-digit number: 01, 02, 03, etc.)
- **Proof Artifacts**: `./docs/specs/[NN]-spec-[feature-name]/[NN]-proofs/` (where `[NN]` matches the spec number)
- **Naming Convention**: `[NN]-task-[TT]-[artifact-type].[ext]` (e.g., `03-task-01-proofs.html` where NN is spec number, TT is task number; screenshot images referenced by a proof file live alongside it in the same `[NN]-proofs/` directory, e.g. `03-task-01-database-list.png`, and are referenced with a relative `src`, never a leading `/`)

### File Update Protocol

1. Update task status (`data-status` attribute) immediately after any state change
2. Save task file after each update
3. Include task file in git commits
4. Never proceed without saving task file

## HTML Artifact Shell

Proof artifact documents share the same self-contained HTML shell used by every other artifact in this skill. Use this exact shared shell/style block when generating a proof file (only change `<title>`, the badge-type text/label, and the body content — keep the CSS byte-identical to the rest of the sdd-html skill's artifacts):

```html
<!doctype html>
<html lang="en">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1" />
<title>[NN]-TYPE-[feature-name]</title>
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
</style>
</head>
<body data-overall-status="...">
<div class="doc">
  <header class="doc-header">
    <span class="badge badge-type">TYPE</span>
    <h1>[NN]-TYPE-[feature-name]</h1>
    <p class="meta">Sequence [NN] · Feature: [feature-name] · Phase N</p>
  </header>
  ...sections...
</div>
</body>
</html>
```

Every generated artifact is a single self-contained HTML5 file: inline `<style>`, no external assets except relative-path local images placed alongside the HTML file. It must render correctly when opened directly in a browser with no server. The `<style>` block may be extended with extra rules a doc type needs (for proof files, add rules for `.artifact`, `.proves`, `.matters`, `.result-summary` as shown in "Recommended Proof File Shape" below), but the classes in the shared block above must never be removed or renamed, so all sdd-html artifacts stay visually consistent across phases.

**Proof artifacts are the one exception to `data-overall-status`**: that attribute is reserved for audit/validation reports whose top-level pass/fail state matters to the router. Proof documents demonstrate evidence rather than declare a gate outcome, so omit `data-overall-status` from `<body>` in proof files.

## Proof Artifact Requirements

Each parent task must include artifacts that:

- **Demonstrate functionality** (screenshots, URLs, CLI output)
- **Verify quality** (test results, lint output, performance metrics)
- **Enable validation** (provide evidence for the validation phase)
- **Support troubleshooting** (logs, error messages, configuration states)

Proof artifacts must be optimized for fast human review, not just raw evidence storage.

- Lead with what the task proves before showing raw output.
- Use descriptive headings that name the task outcome, not just the proof filename.
- Explain why each artifact matters before presenting commands, logs, or screenshots.
- Keep raw evidence intact, but front-load interpretation so a reviewer can understand the result quickly.
- For screenshots, always show the artifact path above the image and embed the image inline in the proof file using a real `<figure><img>` element (HTML lets you show screenshots at full fidelity, zoomable in-browser, unlike flat markdown image syntax — lean into that).
- If output is long, summarize the important result first and then include the most relevant excerpt or reference to the full artifact path.

### Security Warning

**CRITICAL**: Proof artifacts will be committed to the repository. Never include sensitive data:

- Replace API keys, tokens, and secrets with placeholders like `[YOUR_API_KEY_HERE]` or `[REDACTED]`
- Sanitize configuration examples to remove credentials
- Use example or dummy values instead of real production data
- Review all proof artifact files before committing to ensure no sensitive information is present

### Proof Artifact Creation Protocol

```markdown
## PROOF ARTIFACT CREATION CHECKLIST

For each parent task completion:

[ ] **Directory Ready**: `./docs/specs/[NN]-spec-[feature-name]/[NN]-proofs/` exists
[ ] **Review Task Requirements**: Check what proof artifacts the task specifically requires
[ ] **Create Single Proof File**: Create `[spec-number]-task-[task-number]-proofs.html` using the shared HTML Artifact Shell above (badge-type `PROOFS`, no `data-overall-status` attribute needed)
[ ] **Use A Reviewable Structure**:
   - `<h1>Task [TT] Proofs - [descriptive task outcome]</h1>`
   - `<h2>Task Summary</h2>` explaining what was built and why this task matters
   - `<h2>What This Task Proves</h2>` mapping the task to the key behaviors now working
   - `<h2>Evidence Summary</h2>` giving a short reviewer-oriented overview before raw artifacts
[ ] **Document Each Artifact With Context Before Evidence**:
   - `<article class="artifact"><h2>Artifact: [descriptive name]</h2>`
   - `<p class="proves"><strong>What it proves:</strong> [specific behavior or requirement validated]</p>`
   - `<p class="matters"><strong>Why it matters:</strong> [why a reviewer should care about this artifact]</p>`
   - `<p><strong>Command</strong></p>` or `<p><strong>Artifact path</strong></p>`
   - `<p class="result-summary"><strong>Result summary:</strong> [1-3 sentence interpretation of the evidence]</p>`
   - Raw evidence block (`<pre><code>`) or inline `<figure><img>`
[ ] **Present Screenshots For Fast Review**:
   - Show the screenshot file path in its own line above the image
   - Embed every screenshot inline using a real `<figure><img src="relative/path.png" alt="..."><figcaption>...</figcaption></figure>` element
   - Use alt text that describes the visible behavior being proven
[ ] **Keep Verification Context Near The Top**:
   - Do not rely on a bottom-only `Verification` section to explain relevance
   - Place interpretation before the raw command output, logs, or screenshots
[ ] **End With A Short Reviewer Conclusion**:
   - State the final conclusion the reviewer should draw from the combined evidence
[ ] **Format with HTML**: Use `<pre><code>` blocks, headings, tables, and clear semantic structure
[ ] **Verify File Content**: Ensure the HTML file contains all required evidence
[ ] **Security Check**: Scan proof file for API keys, tokens, passwords, or other sensitive data and replace with placeholders

**SIMPLE VERIFICATION**: One file per task, all evidence included
**CONTENT VERIFICATION**: Check the HTML file contains both context-setting summary sections and raw evidence
**VERIFICATION**: Ensure proof artifact file demonstrates all required functionality
**SECURITY VERIFICATION**: Confirm no real credentials or sensitive data are present

**The single HTML proof file must be created BEFORE the parent task commit**
```

### Recommended Proof File Shape

Use the shared shell/style block from the "HTML Artifact Shell" section above with `TYPE` = `proofs` and `badge-type` = `PROOFS`. Do not add `data-overall-status` to `<body>` (proof files are evidence, not a pass/fail gate). Extend the `<style>` block with a few extra rules for the artifact cards, then populate the body with the sections shown in this complete worked example:

```html
<!doctype html>
<html lang="en">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1" />
<title>03-task-02-proofs</title>
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
  article.artifact { margin: 2rem 0; padding-top: 1.25rem; border-top: 1px solid var(--border); }
  .proves, .matters, .result-summary { margin: .5rem 0; }
</style>
</head>
<body>
<div class="doc">
  <header class="doc-header">
    <span class="badge badge-type">PROOFS</span>
    <h1>Task 02 Proofs - Metabase bootstrap and first-run setup</h1>
    <p class="meta">Sequence 03 · Feature: metabase-bootstrap · Task 02 · Phase 3</p>
  </header>

  <h2>Task Summary</h2>
  <p>This task added the Metabase bootstrap CLI script and the accompanying first-run setup flow. Before this task, standing up a working Metabase instance required a developer to manually click through Metabase's first-run setup wizard in a browser (admin account, sample database connection, and instance verification). The bootstrap script now performs that entire sequence non-interactively from the CLI, so a developer can go from a freshly started Metabase container to a fully configured, authenticated instance with a single command.</p>

  <h2>What This Task Proves</h2>
  <ul>
    <li>The bootstrap CLI completes Metabase's first-run setup wizard end-to-end — admin account creation and the sample database connection — without any manual browser interaction.</li>
    <li>The bootstrapped instance is immediately reachable and healthy, confirmed by a passing health-check request right after bootstrap completes.</li>
    <li>The bootstrap flow is covered by an automated pytest suite that exercises the same code path used in CI, so regressions are caught before merge.</li>
  </ul>

  <h2>Evidence Summary</h2>
  <ul>
    <li>CLI run: the bootstrap command exits with <code>status: ok</code> and prints the reachable Metabase URL.</li>
    <li>Health check: <code>GET /api/health</code> returns HTTP 200 immediately after bootstrap.</li>
    <li>Test suite: <code>pytest tests/test_bootstrap.py</code> passes with all first-run-flow assertions green.</li>
  </ul>

  <article class="artifact">
    <h2>Artifact: Successful bootstrap CLI run</h2>
    <p class="proves"><strong>What it proves:</strong> The bootstrap script drives Metabase's first-run setup wizard to completion — admin account and sample database connection — using only CLI flags, with no manual browser steps.</p>
    <p class="matters"><strong>Why it matters:</strong> This is the core automation this task delivers. Without it, every developer and every CI job that needs a working Metabase instance would have to click through the setup wizard by hand, which does not scale and cannot run headless.</p>
    <p><strong>Command:</strong></p>
    <pre><code>python -m devspace.metabase bootstrap --profile local --non-interactive --json</code></pre>
    <p class="result-summary"><strong>Result summary:</strong> The command completed in under 10 seconds and returned a JSON status object with <code>"status": "ok"</code> and the instance URL, confirming the wizard finished successfully.</p>
    <pre><code class="language-json">{
  "status": "ok",
  "url": "http://localhost:3000",
  "admin_email": "admin@example.com",
  "setup_token_used": true,
  "database_configured": "sample-dataset"
}</code></pre>
  </article>

  <article class="artifact">
    <h2>Artifact: Metabase health endpoint</h2>
    <p class="proves"><strong>What it proves:</strong> The Metabase instance is reachable and reports healthy immediately after bootstrap, not just that the bootstrap command exited successfully.</p>
    <p class="matters"><strong>Why it matters:</strong> A successful CLI exit code does not guarantee the underlying service is actually serving traffic. This independent health check closes that gap and is the same check used by the deployment's readiness probe.</p>
    <p><strong>Command:</strong></p>
    <pre><code>curl -s http://localhost:3000/api/health</code></pre>
    <p class="result-summary"><strong>Result summary:</strong> The endpoint returned HTTP 200 with a healthy status body on the first request after bootstrap, with no retries needed.</p>
    <pre><code class="language-json">{
  "status": "ok"
}</code></pre>
  </article>

  <article class="artifact">
    <h2>Artifact: Database list screenshot</h2>
    <p class="proves"><strong>What it proves:</strong> The sample database configured during bootstrap is visible and connected inside the Metabase admin UI, not just present in the API response.</p>
    <p class="matters"><strong>Why it matters:</strong> This is the same view a developer would check by hand to confirm the setup worked, so it gives reviewers a direct, human-verifiable confirmation that matches the CLI and health-check evidence above.</p>
    <p><strong>Artifact path:</strong> <code>docs/specs/03-spec-metabase-bootstrap/03-proofs/03-task-02-database-list.png</code></p>
    <p class="result-summary"><strong>Result summary:</strong> The Databases admin page lists exactly one connected database, "Sample Dataset," with a green "Connected" status indicator.</p>
    <figure>
      <img src="03-task-02-database-list.png" alt="Metabase admin Databases page showing the Sample Dataset database connected with a green status indicator" />
      <figcaption>docs/specs/03-spec-metabase-bootstrap/03-proofs/03-task-02-database-list.png</figcaption>
    </figure>
  </article>

  <h2>Reviewer Conclusion</h2>
  <p>Together, the CLI run, the health-check response, and the database-list screenshot show that Metabase's first-run setup is now fully automated and independently verifiable at three different layers — the automation's own report, an external health probe, and the admin UI a human would check manually. Combined with the passing pytest suite, this task is ready for review.</p>
</div>
</body>
</html>
```

Preserve every field shown above (Task Summary, What This Task Proves, Evidence Summary, one `<article class="artifact">` block per artifact with proves/matters/command-or-path/result-summary/raw-evidence, and a closing Reviewer Conclusion) for every proof file you generate — none of these are optional decorations, they are the content contract carried over from the Markdown edition, now expressed as semantic HTML instead of Markdown headings and image syntax. For screenshot artifacts, always show the artifact path both above and below the image (a `<p><strong>Artifact path:</strong>...</p>` line and a matching `<figcaption>`), and use a relative `src` with no leading `/`.

## Git Workflow Protocol

### Commit Requirements

- **Frequency**: One commit per parent task minimum
- **Format**: Conventional commits with task references
- **Content**: Include all code changes and task file updates
- **Message**:

  ```bash
  git commit -m "feat: [task-description]" -m "- [key-details]" -m "Related to T[task-number] in Spec [spec-number]"
  ```

- **Verification**: Always verify with `git log --oneline -1` after committing

### Branch Management

- Work on the appropriate branch for the spec
- Keep commits clean and atomic
- Include proof artifacts in commits when appropriate

### Commit Validation Protocol

```markdown
## COMMIT CREATION CHECKLIST

Before marking parent task as complete:

[ ] All code changes staged: `git add .`
[ ] Task file updates included in staging
[ ] Proof artifacts created and included
[ ] Commit message follows conventional format
[ ] Task reference included in commit message
[ ] Spec number included in commit message
[ ] Commit created successfully
[ ] Verification passed: `git log --oneline -1`

**Only after commit verification passes may you mark parent task's `data-status="done"`**
```

## What Happens Next

After completing all tasks in the task list:

1. **Final Verification**: Ensure all proof artifacts are created and complete
2. **Proof Artifact Validation**: Verify all proof artifacts demonstrate functionality from original spec
3. **Test Suite**: Run final comprehensive test suite
4. **Documentation**: Update any relevant documentation
5. **Handoff**: Instruct user to continue to the validation phase

The validation phase will use your proof artifacts as evidence to verify that the spec has been fully and correctly implemented.

## Instructions

1. **Locate Task File**: Find the task list in `./docs/specs/` directory
2. **Verify Planning Audit Status**: Confirm audit report exists and all REQUIRED gates passed before any implementation work
3. **Present Checkpoints**: Show checkpoint options and confirm user preference
4. **Execute Workflow**: Follow the structured workflow with self-verification checklists
5. **Validate Progress**: Use verification checkpoints before proceeding
6. **Track Progress**: Update task file immediately after any status changes
7. **Complete or Continue**:
   - If tasks remain, proceed to next parent task
   - If all complete, instruct user to proceed to validation

## Implementation Verification Sequence

**For each parent task, follow this exact sequence:**

1. Sub-tasks → 2. Demo verification → 3. Proof artifacts → 4. Git commit → 5. Parent task completion → 6. Validation → 7. Next task

**Critical checkpoints that block progression:**

- Sub-task verification before next sub-task
- Proof artifact verification before commit
- Commit verification before parent task completion
- Full validation before next parent task

## Error Recovery

If you encounter issues:

1. **Stop immediately** at the point of failure
2. **Assess the problem** using the relevant verification checklist
3. **Fix the issue** before proceeding
4. **Re-run verification** to confirm the fix
5. **Document the issue** in task comments if needed

## Success Criteria

Implementation is successful when:

- All parent tasks have `data-status="done"` in task file
- Proof artifacts exist for each parent task
- Git commits follow repository format with proper frequency
- All tests pass using repository's testing approach
- Proof artifacts demonstrate all required functionality
- Repository quality gates pass consistently
- Task file accurately reflects final status
- Implementation follows established repository patterns and conventions

## How to Continue the SDD Workflow

Likely next phase action: the skill will route to Phase 4 and validate the implementation against the spec and proof artifacts using strict pass/fail gates.

To continue the workflow in this chat, reply with:

`Continue SDD with validation.`

You can also continue in a new chat if you want to keep context lean; the SDD skill will reassess repository state from the persisted spec/task/audit/proof artifacts.
