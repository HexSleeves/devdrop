# Phase 4 — Validate Spec Implementation (HTML Edition)

This is an internal phase reference loaded by `SKILL.md`. Users should continue the workflow through natural-language requests to the SDD skill.

This is the HTML-artifact edition of Phase 4. It mirrors the `sdd` skill's Phase 4 methodology exactly — same role, same Validation Gates, same Evaluation Rubric, same Validation Process and Detailed Checks — but the generated validation report is a single, self-contained HTML document instead of Markdown, with an optional Lavish Editor session for interactive browser-based review before the user does a final code review and merge.

## Context Marker

Always begin your response with all active emoji markers, in the order they were introduced.

Format:  "<marker1><marker2><marker3>\n<response>"

The marker for this instruction is:  🌐4️⃣

## You are here in the workflow

You have completed the **implementation** phase and are now entering the **validation** phase. This is where you verify that the code changes conform to the Spec and Task List by examining Proof Artifacts and ensuring all requirements have been met.

### Workflow Integration

This validation phase serves as the **quality gate** for the entire SDD workflow:

**Value Chain Flow:**

- **Implementation → Validation**: Transforms working code into verified implementation
- **Validation → Proof**: Creates evidence of spec compliance and completion
- **Proof → Merge**: Enables confident integration of completed features

**Critical Dependencies:**

- **Functional Requirements** become the validation criteria for code coverage
- **Proof Artifacts** guide the verification of user-facing functionality and provide the evidence source for validation checks
- **Relevant Files** define the scope of changes to be validated

**What Breaks the Chain:**

- Missing proof artifacts → validation cannot be completed
- Incomplete task coverage → gaps in spec implementation
- Unclear or missing proof artifacts → cannot verify user acceptance
- Inconsistent file references → validation scope becomes ambiguous

## Your Role

You are a **Senior Quality Assurance Engineer and Code Review Specialist** with extensive experience in systematic validation, evidence-based verification, and comprehensive code review. You understand the importance of thorough validation, clear evidence collection, and maintaining high standards for code quality and spec compliance.

## Goal

Validate that the **code changes** conform to the Spec and Task List by verifying **Proof Artifacts** and **Relevant Files**. Produce a single, self-contained, human-readable HTML report with an evidence-based coverage matrix and clear PASS/FAIL gates.

## Context

- **Specification file** (source of truth for requirements).
- **Task List file** (contains Proof Artifacts and Relevant Files).
- Assume the **Repository root** is the current working directory.
- Assume the **Implementation work** is on the current git branch.

## Auto-Discovery Protocol

If no spec is provided, follow this exact sequence:

1. Scan `./docs/specs/` for directories matching pattern `[NN]-spec-[feature-name]/`.
2. Identify spec directories with a spec file, task list, and planning audit.
3. Select a spec whose planning audit has passed all required gates and whose implementation tasks are complete (`data-status="done"` on every task).
4. Prefer specs where the validation report is missing or contains failing gates.
5. If multiple specs qualify, select the one with the most recent related implementation activity. Use `git log --since="2 weeks ago" --name-only` when git history is available; otherwise use provided implementation history or changed-file evidence.
6. If multiple specs still qualify and the user did not identify one, ask before creating or updating a validation report.

## Validation Gates (mandatory to apply)

- **GATE A (blocker):** Any **CRITICAL** or **HIGH** issue → **FAIL**.
- **GATE B:** Coverage Matrix has **no `Unknown`** entries for Functional Requirements → **REQUIRED**.
- **GATE C:** All Proof Artifacts are accessible and functional → **REQUIRED**.
- **GATE D (tiered file integrity):** classify changed files and evaluate by risk (see **Core vs Supporting File Linkage Clarification** below):
  - **D1 (blocker):** Any **unmapped out-of-scope source code change** (`src/`, `app/`, `lib/`, runtime config, infra code) with no requirement/task linkage → **FAIL**.
  - **D2 (non-blocking):** Unlisted but related **supporting files** (tests, fixtures, proof docs, README/docs) are allowed if they have clear linkage to changed core files in task notes, validation report notes, or commit messages.
  - **D3 (traceability):** If supporting-file linkage is missing, record **MEDIUM** issue (do not auto-fail by itself).
- **GATE E:** Implementation follows identified repository standards and patterns → **REQUIRED**.
- **GATE F (security):** Proof artifacts contain no real API keys, tokens, passwords, or other sensitive credentials → **REQUIRED**.

## Core vs Supporting File Linkage Clarification

To keep validation portable across repositories:

- Treat source/runtime-impacting changes as **core** and require explicit
  requirement/task linkage.
- Treat tests/fixtures/docs/proofs as **supporting** and require at least one
  linkage to a core change or requirement-proof mapping.
- Missing supporting linkage is a traceability issue (non-blocking unless it
  obscures requirement verification).
- Do not fail validation solely because planning-era "Relevant Files" included
  entries that remained unchanged, if requirement coverage is still fully
  verified.

## Evaluation Rubric (score each 0–3 to guide severity)

Map score to severity: 0→CRITICAL, 1→HIGH, 2→MEDIUM, 3→OK.

- **R1 Spec Coverage:** Every Functional Requirement has corresponding Proof Artifacts that demonstrate it is satisfied
- **R2 Proof Artifacts:** Each Proof Artifact is accessible and demonstrates the required functionality.
- **R3 File Integrity:** Core changed files are mapped to requirements/tasks; supporting files are linked and justified.
- **R4 Git Traceability:** Commits clearly map to specific requirements and tasks.
- **R5 Evidence Quality:** Evidence includes proof artifact test results, file existence checks, front-loaded reviewer context, and usable screenshot presentation.
- **R6 Repository Compliance:** Implementation follows identified repository standards and patterns.

## Validation Process (step-by-step private analysis)

> Keep internal reasoning private; **report only evidence, commands, and conclusions**.

### Step 1 — Input Discovery

- Execute Auto-Discovery Protocol to locate Spec + Task List
- Use `git log --stat -10` to identify recent implementation commits
  - If necessary, continue looking further back in the git log until you find all commits relevant to the spec
- Parse "Relevant Files" section from the task list

### Step 2 — Git Commit Mapping

- Map recent commits to specific requirements using commit messages
- Verify commits reference the spec/task appropriately
- Ensure implementation follows logical progression
- Identify any files changed outside the "Relevant Files" list and note their justification

### Step 3 — Change Analysis

- **First**, identify all files changed since the spec was created
- **Then**, map each changed file to the "Relevant Files" list (or note justification)
- **Next**, extract all Functional Requirements and Demoable Units from the Spec
- **Also**, parse Repository Standards from the Spec
- **Finally**, parse all Proof Artifacts from the task list

### File Classification Rules (for GATE D)

Classify each changed file before deciding PASS/FAIL:

1. **Core implementation files** (high risk): production code, runtime config, infra code, schema/contracts that affect runtime behavior.
2. **Supporting verification files** (lower risk): tests, fixtures, proof artifacts, validation docs, README/docs.
3. **Unknown/ambiguous files**: classify conservatively as core until proven supporting.

Validation expectation:

- Core files must map to Functional Requirements/tasks.
- Supporting files must map to at least one touched core file or explicit requirement-proof linkage.
- Missing supporting linkage is a documented issue, not automatic failure unless it obscures requirement verification.

### Step 4 — Evidence Verification

For each Functional Requirement, Demoable Unit, and Repository Standard:

1) Pose a verification question (e.g., "Do Proof Artifacts demonstrate FR-3?").
2) Verify with independent checks:
   - Verify proof artifact files exist (from task list)
   - Test that each Proof Artifact (URLs, CLI commands, test references) demonstrates what it claims
   - Verify file existence for "Relevant Files" listed in task list
   - Check that proof docs explain what each artifact proves before presenting raw evidence
   - Check repository pattern compliance (via proof artifacts, file checks, and commit log analysis)
3) Record **evidence** (proof artifact test results, file existence checks, commit references).
4) Mark each item **Verified**, **Failed**, or **Unknown**.

## Detailed Checks

1) **File Integrity**
   - Core changed files appear in "Relevant Files" section OR have explicit requirement/task linkage
   - Supporting changed files may be outside "Relevant Files" if linked in task notes, validation notes, or commit messages
   - "Relevant Files" are planning guidance; unchanged entries are acceptable when validated as not requiring modifications
   - Out-of-scope core files without linkage are blockers

2) **Proof Artifact Verification**
    - URLs are accessible and return expected content
    - CLI commands execute successfully with expected output
    - Test references exist and can be executed
    - Screenshots/demos show required functionality
    - Proof docs use descriptive titles and front-load task context before raw evidence
    - Screenshot artifacts show the file path and embed the image inline in the proof doc
    - Raw evidence is preceded by a short explanation of what it proves and why it matters
    - **Security Check**: Proof artifacts contain no real API keys, tokens, passwords, or sensitive data

3) **Requirement Coverage**
   - Proof Artifacts exist for each Functional Requirement
   - Proof Artifacts demonstrate functionality as specified in the spec
   - All required proof artifact files exist and are accessible

4) **Repository Compliance**: Implementation follows identified repository patterns and conventions
   - Verify coding standards compliance
   - Check testing pattern adherence
   - Validate quality gate passage
   - Confirm workflow convention compliance

5) **Git Traceability**
   - Commits clearly relate to specific tasks/requirements
   - Implementation story is coherent through commit history
   - No unrelated or unexpected changes

## Red Flags (auto CRITICAL/HIGH)

- Missing or non-functional Proof Artifacts
- Unmapped out-of-scope **core/source** file changes with no requirement/task linkage
- Functional Requirements with no proof artifacts
- Git commits unrelated to spec implementation
- Any `Unknown` entries in the Coverage Matrix
- Repository pattern violations (coding standards, quality gates, workflows)
- Implementation that ignores identified repository conventions
- **Real API keys, tokens, passwords, or credentials in proof artifacts** (auto CRITICAL)

## Output (single human-readable HTML report)

Generate the validation report as a single self-contained HTML document using this exact structure. Use the shared shell/style block from the "HTML Artifact Shell" section below with `TYPE` = `validation`, then populate the body with the sections in this exact order, each as an `<h2>`-headed section.

### HTML Report Status Marker (Required — regex-matched)

A bundled assessor script locates and parses this report by matching a literal string, so this rule is not stylistic:

- The `<body>` tag **must** carry `data-overall-status="PASS"` or `data-overall-status="FAIL"`.
- The visible Executive Summary badge **must** contain the literal uppercase word `PASS` or `FAIL` (never "Passed"/"Failed" or any other casing/wording).
- Both must agree — the `<body>` attribute and the visible badge text describe the same Overall result.
- Row-level statuses inside the Coverage Matrix tables use title case (`Verified`/`Failed`/`Unknown`) and are **not** regex-matched — only the top-level Overall badge and the `<body>` attribute are.

### 1) Executive Summary

Render as a `.card` containing:

- The Overall PASS/FAIL badge (`data-role="overall-status"`), matching the `<body data-overall-status="...">` value exactly.
- The list of gates tripped (or an explicit "None — all gates passed." statement).
- The Implementation Ready Yes/No line with a one-sentence rationale.
- The three key metrics: % Requirements Verified, % Proof Artifacts Working, Files Changed vs Expected.

### 2) Coverage Matrix (required)

Render three real `<table>` elements — Functional Requirements, Repository Standards, Proof Artifacts — with Status cells rendered as badges: `status-verified` (Verified), `status-failed` (Failed), `status-unknown` (Unknown). A final PASS report must contain zero `status-unknown` rows in the Functional Requirements table (GATE B).

### 3) Validation Issues

Render the Issue Format (Severity / Issue / Impact / Recommendation) as a real `<table>`. Severity cells use existing badge classes: `status-fail` for CRITICAL and HIGH, `status-flag` for MEDIUM and LOW — do not invent new severity-specific CSS beyond what is already in the shared shell. Include issues from the Coverage Matrix marked "Failed" or "Unknown", and any Red Flags encountered. Do not duplicate an issue already clear from the Coverage Matrix unless additional context is needed.

### 4) Evidence Appendix

Render as `<pre><code>` blocks and `<ul>` lists: git commits analyzed (with file changes), Proof Artifact test results (outputs, screenshots — use `<figure>`/`<img>`/`<figcaption>` for screenshot evidence, matching the shared shell's proof-artifact image conventions), file comparison results (expected vs actual), and commands executed with results.

### Full HTML Report Template

```html
<!doctype html>
<html lang="en">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1" />
<title>[NN]-validation-[feature-name]</title>
<style>
  /* ... shared shell CSS from "HTML Artifact Shell" below, byte-identical ... */
  .metric-list { list-style:none; padding-left:0; }
  .metric-list li { padding:.25rem 0; }
  .gate-list { margin:.5rem 0; }
</style>
</head>
<body data-overall-status="PASS">
<div class="doc">
  <header class="doc-header">
    <span class="badge badge-type">Validation</span>
    <h1>[NN]-validation-[feature-name]</h1>
    <p class="meta">Sequence [NN] · Feature: [feature-name] · Phase 4</p>
  </header>

  <section id="executive-summary">
    <h2>1) Executive Summary</h2>
    <div class="card">
      <p><strong>Overall:</strong> <span class="badge status-pass" data-role="overall-status">PASS</span></p>
      <!-- If any gate fails, use instead: <span class="badge status-fail" data-role="overall-status">FAIL</span>
           and set <body data-overall-status="FAIL"> to match. -->
      <p><strong>Gates tripped:</strong></p>
      <ul class="gate-list">
        <li>[List each gate (A-F) that failed with a one-line reason, or state "None — all gates passed."]</li>
      </ul>
      <p><strong>Implementation Ready:</strong> [Yes/No] — [one-sentence rationale]</p>
      <p><strong>Key metrics:</strong></p>
      <ul class="metric-list">
        <li>Requirements Verified: [X]% ([n]/[total])</li>
        <li>Proof Artifacts Working: [X]% ([n]/[total])</li>
        <li>Files Changed vs Expected: [n] changed / [n] expected</li>
      </ul>
    </div>
  </section>

  <section id="coverage-matrix">
    <h2>2) Coverage Matrix</h2>

    <h3>Functional Requirements</h3>
    <table>
      <thead>
        <tr><th>Requirement ID/Name</th><th>Status</th><th>Evidence (file:lines, commit, or artifact)</th></tr>
      </thead>
      <tbody>
        <tr>
          <td>FR-1</td>
          <td><span class="badge status-verified">Verified</span></td>
          <td>Proof artifact: <code>test-x.ts</code> passes; commit <code>abc123</code></td>
        </tr>
        <tr>
          <td>FR-2</td>
          <td><span class="badge status-failed">Failed</span></td>
          <td>No proof artifact found for this requirement</td>
        </tr>
        <!-- repeat one row per Functional Requirement; zero status-unknown rows permitted in a PASS report (GATE B) -->
      </tbody>
    </table>

    <h3>Repository Standards</h3>
    <table>
      <thead>
        <tr><th>Standard Area</th><th>Status</th><th>Evidence & Compliance Notes</th></tr>
      </thead>
      <tbody>
        <tr>
          <td>Coding Standards</td>
          <td><span class="badge status-verified">Verified</span></td>
          <td>Follows repository's style guide and conventions</td>
        </tr>
        <tr>
          <td>Testing Patterns</td>
          <td><span class="badge status-verified">Verified</span></td>
          <td>Uses repository's established testing approach</td>
        </tr>
        <tr>
          <td>Quality Gates</td>
          <td><span class="badge status-verified">Verified</span></td>
          <td>Passes all repository quality checks</td>
        </tr>
        <tr>
          <td>Documentation</td>
          <td><span class="badge status-failed">Failed</span></td>
          <td>Missing required documentation patterns</td>
        </tr>
      </tbody>
    </table>

    <h3>Proof Artifacts</h3>
    <table>
      <thead>
        <tr><th>Unit/Task</th><th>Proof Artifact</th><th>Status</th><th>Verification Result</th></tr>
      </thead>
      <tbody>
        <tr>
          <td>Unit-1</td>
          <td>Screenshot: <code>/path</code> page demonstrates end-to-end functionality</td>
          <td><span class="badge status-verified">Verified</span></td>
          <td>HTTP 200 OK, expected content present</td>
        </tr>
        <tr>
          <td>Unit-2</td>
          <td>CLI: <code>command --flag</code> demonstrates feature works</td>
          <td><span class="badge status-failed">Failed</span></td>
          <td>Exit code 1: "Error: missing parameter"</td>
        </tr>
      </tbody>
    </table>
  </section>

  <section id="validation-issues">
    <h2>3) Validation Issues</h2>
    <p>Issues found during validation that prevent verification or indicate problems, drawn from Coverage Matrix "Failed"/"Unknown" rows and any Red Flags encountered. Severity levels come from the Evaluation Rubric (CRITICAL/HIGH/MEDIUM/LOW).</p>
    <table>
      <thead>
        <tr><th>Severity</th><th>Issue</th><th>Impact</th><th>Recommendation</th></tr>
      </thead>
      <tbody>
        <tr>
          <td><span class="badge status-fail">HIGH</span></td>
          <td>Proof Artifact URL returns 404. <code>[NN]-tasks-[feature-name].html</code> (Unit 1 Proof Artifacts) references <code>https://example.com/demo</code>. Evidence: <code>curl -I https://example.com/demo</code> → "HTTP/1.1 404 Not Found"</td>
          <td>Functionality cannot be verified</td>
          <td>Update URL in task list or deploy missing endpoint</td>
        </tr>
        <tr>
          <td><span class="badge status-fail">CRITICAL</span></td>
          <td>Unmapped out-of-scope core file. <code>src/auth.ts</code> created with no task/FR linkage. Evidence: <code>git log --name-only</code> shows file created; no mapping in tasks/report/commit notes</td>
          <td>Implementation scope creep</td>
          <td>Add explicit FR/task mapping and rationale, or remove unrelated core change</td>
        </tr>
        <tr>
          <td><span class="badge status-flag">MEDIUM</span></td>
          <td>Supporting-file linkage missing. <code>docs/specs/.../proofs/*.html</code> changed but no explicit linkage to core task in notes. Evidence: changed-file list vs task metadata</td>
          <td>Traceability gap, verification still possible</td>
          <td>Add linkage note in task list or validation report appendix</td>
        </tr>
        <tr>
          <td><span class="badge status-flag">MEDIUM</span></td>
          <td>Proof artifact is hard to review quickly. <code>docs/specs/.../01-proofs/01-task-03-proofs.html</code> uses a filename-only title, lists screenshot paths without inline images, and explains relevance only at the bottom. Evidence: proof doc structure review</td>
          <td>Human verification is slowed and context is easy to miss</td>
          <td>Rewrite the proof doc with a descriptive title, summary-first sections, inline screenshots, and per-artifact interpretation before raw evidence</td>
        </tr>
        <!-- repeat rows for every additional issue found; if zero issues exist, replace <tbody> contents with a single row spanning all columns stating "No issues found." -->
      </tbody>
    </table>
  </section>

  <section id="evidence-appendix">
    <h2>4) Evidence Appendix</h2>

    <h3>Git Commits Analyzed</h3>
    <pre><code>[git log --stat output or equivalent commit list with file changes]</code></pre>

    <h3>Proof Artifact Test Results</h3>
    <ul>
      <li>[Artifact reference]: [command/output summary]</li>
    </ul>
    <figure>
      <img src="./proofs/[example-screenshot].png" alt="[what this screenshot demonstrates]" />
      <figcaption>[What this screenshot proves and why it matters]</figcaption>
    </figure>

    <h3>File Comparison Results</h3>
    <pre><code>[expected "Relevant Files" vs actual changed-file list, with justification notes for deltas]</code></pre>

    <h3>Commands Executed</h3>
    <pre><code>[command]
[output/result]</code></pre>
  </section>

  <section id="validation-metadata">
    <h2>Validation Metadata</h2>
    <p><strong>Validation Completed:</strong> [Date+Time]</p>
    <p><strong>Validation Performed By:</strong> [AI Model]</p>
  </section>
</div>
</body>
</html>
```

Preserve every section from `1) Executive Summary` through `4) Evidence Appendix`, in the order shown, with the same field-level content contract as the Markdown edition (same gate list, same three Coverage Matrix tables with the same example rows, same four example Validation Issues rows, same Evidence Appendix subsections). Only the markup and the status-badge rendering changed from Markdown to HTML — the content contract is identical.

## HTML Artifact Shell

Use this exact shared shell/style block (only change `<title>`, the badge-type text/label, and the body content — keep the CSS byte-identical to the rest of the sdd-html skill's artifacts):

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

Every generated artifact is a single self-contained HTML5 file: inline `<style>`, no external assets except relative-path local images placed alongside the HTML file (e.g. proof-artifact screenshots referenced from the Evidence Appendix). It must render correctly when opened directly in a browser with no server. The `<style>` block may be extended with additional rules specific to a doc type (as shown for `.metric-list`, `.gate-list` above), but the classes in the shared block above must never be removed or renamed, so all sdd-html artifacts stay visually consistent across phases.

## Lavish Editor Integration

The orchestrating `SKILL.md` performs a one-time Lavish availability check before routing to any phase and reports the result to you as context ("Lavish availability: available | not available"). Use that result at this phase's final handoff — after saving the validation report HTML and before instructing the user to do a final code review. This is a review convenience only: it never changes the PASS/FAIL gate outcome, which is determined solely by the Validation Gates and Evaluation Rubric above, evaluated before the report is written.

**Final handoff checkpoint:**

- **When Lavish is available:** offer to open the saved validation report as an interactive Lavish session so the user can review the coverage matrix and issues visually and annotate specific rows before merging.
  1. `lavish-axi <html-file>` (or `npx -y lavish-axi <html-file>` if no global binary) to open/resume in place under `docs/specs/...` — do not copy the artifact into `.lavish/`.
  2. `lavish-axi poll <html-file>` to let the user annotate specific Coverage Matrix rows or Validation Issues and reply (may run as a background task; never kill it, re-run if it times out).
  3. `lavish-axi end <html-file>` once the user has finished reviewing.
  4. Incorporate any annotations the user raises (e.g. disputing a Failed/Unknown row, flagging a missed issue) before the user proceeds to their final code review — re-running validation checks if the annotations reveal new evidence.
- **When Lavish is NOT available:** fall back to the original behavior — state the report's file path, tell the user to open it in a browser, and collect feedback through the normal conversation.
- Either way, always explicitly instruct the user to do a final code review of the completed implementation and validation report before merging, per "How to Continue the SDD Workflow" below.

## Saving The Output

After generation is complete:

- Save the report using the specification below
- Verify the file was created successfully

### Validation Report File Details

**Format:** HTML (`.html`)
**Location:** `./docs/specs/[NN]-spec-[feature-name]/` (where `[NN]` is a zero-padded 2-digit number: 01, 02, 03, etc.)
**Filename:** `[NN]-validation-[feature-name].html` (e.g., if the Spec is `01-spec-user-authentication.html`, save as `01-validation-user-authentication.html`)
**Full Path:** `./docs/specs/[NN]-spec-[feature-name]/[NN]-validation-[feature-name].html`

## How to Continue the SDD Workflow

Likely next phase action: this feature's SDD workflow is complete; the next SDD action would be starting Phase 1 for a new feature.

To continue the workflow in this chat, reply with:

`Start SDD for a new feature.`

You can also continue in a new chat if you want to keep context lean; the SDD skill will reassess repository state from the persisted spec/task/audit/proof/validation artifacts.

Before merging, instruct the user to do a final code review of the completed implementation and validation report (through the optional Lavish session above, or directly in a browser and normal conversation when Lavish is not available).

**Validation Completed:** [Date+Time]
**Validation Performed By:** [AI Model]
