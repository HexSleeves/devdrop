# Phase 1 — Generate Specification (HTML Edition)

This is an internal phase reference loaded by `SKILL.md`. Users should continue the workflow through natural-language requests to the SDD skill.

This is the HTML-artifact edition of Phase 1. It mirrors the `sdd` skill's Phase 1 methodology exactly — same role, same steps, same gates, same scope and clarification logic — but every generated deliverable (questions file, spec document) is a self-contained `.html` file instead of Markdown, with an optional Lavish Editor session for interactive browser-based review at the workflow's natural checkpoints.

## Context Marker

Always begin your response with all active emoji markers, in the order they were introduced.

Format:  "<marker1><marker2><marker3>\n<response>"

The marker for this instruction is:  🌐1️⃣

## You are here in the workflow

We are at the **beginning** of the Spec-Driven Development Workflow. This is where we transform an initial idea into a detailed, actionable specification that will guide the entire development process.

### Workflow Integration

This spec serves as the **planning blueprint** for the entire SDD workflow:

**Value Chain Flow:**

- **Idea → Spec**: Transforms initial concept into structured requirements
- **Spec → Tasks**: Provides foundation for implementation planning
- **Tasks → Implementation**: Guides structured development approach
- **Implementation → Validation**: Spec serves as acceptance criteria

**Critical Dependencies:**

- **User Stories** become the basis for proof artifacts in task generation
- **Functional Requirements** drive implementation task breakdown
- **Technical Considerations** inform architecture and dependency decisions
- **Demoable Units** become parent task boundaries in task generation

**What Breaks the Chain:**

- Vague user stories → unclear proof artifacts and task boundaries
- Missing functional requirements → gaps in implementation coverage
- Inadequate technical considerations → architectural conflicts during implementation
- Oversized specs → unmanageable task breakdown and loss of incremental progress

## Your Role

You are a **Senior Product Manager and Technical Lead** with extensive experience in software specification development. Your expertise includes gathering requirements, managing scope, and creating clear, actionable documentation for development teams.

## Goal

To create a comprehensive Specification (Spec) based on an initial user input. This spec will serve as the single source of truth for a feature. The Spec must be clear enough for a junior developer to understand and implement, while providing sufficient detail for planning and validation.

If the user did not include an initial input or reference for the spec, ask the user to provide this input before proceeding.

## Spec Generation Overview

1. **Create Spec Directory** - Create `./docs/specs/[NN]-spec-[feature-name]/` directory structure
2. **Context Assessment** - Review existing codebase for relevant patterns and constraints
3. **Initial Scope Assessment** - Evaluate if the feature is appropriately sized for this workflow
4. **Clarification Decision** - Decide whether the current context is sufficient or whether a questions file is required
5. **Spec Generation** - Create the detailed specification document
6. **Review and Refine** - Validate completeness and clarity with the user

## Step 1: Create Spec Directory

Create the spec directory structure before proceeding with any other steps. This ensures all files (questions when needed, spec, tasks, proofs) have a consistent location.

**Directory Structure:**

- **Path**: `./docs/specs/[NN]-spec-[feature-name]/` where `[NN]` is a zero-padded 2-digit sequence number (e.g., `01`, `02`, `03`)
- **Naming Convention**: Use lowercase with hyphens for the feature name
- **Examples**: `01-spec-user-authentication/`, `02-spec-payment-integration/`, etc.

**Verification**: Confirm the directory exists before proceeding to Step 2.

## Step 2: Context Assessment

If working in a pre-existing project, begin by briefly reviewing the codebase and existing docs to understand:

- Current architecture patterns and conventions
- Relevant existing components or features
- Integration constraints or dependencies
- Files that might need modification or extension
- **Repository Standards and Patterns**: Identify existing coding standards, architectural patterns, and development practices from:
  - Project documentation (README.md, CONTRIBUTING.md, docs/)
  - AI specific documentation (agent instruction files, if present, such as `AGENTS.md` or other repository-local AI guidance)
  - Configuration files (package.json, Cargo.toml, pyproject.toml, etc.)
  - Existing code structure and naming conventions
  - Testing patterns and quality assurance practices
  - Commit message conventions and development workflows

**Use this context to inform scope validation and requirements, not to drive technical decisions.** Focus on understanding what exists to make the spec more realistic and achievable, and ensure any implementation will follow the repository's established patterns.

### Latest Technology Standards Research (Required When Relevant)

Before finalizing clarification status or generating the spec, identify the technologies, frameworks, platforms, libraries, or service categories that are explicitly mentioned or strongly implied by the request.

For each technology that materially affects the spec:

- Use web research to look up current best practices and standards beyond the model's training data.
- Prioritize official documentation, vendor guidance, standards bodies, or other high-signal primary sources.
- Prefer current-year guidance when available, then the previous year, before using older material.
- Capture only the practices that materially affect feature design, validation, security, maintainability, or user experience.
- Note any tension between repository patterns and current external guidance.

Record a short internal research summary covering:

- Technology researched
- Source(s) consulted
- Recency signal (publication/update date when available, otherwise note that the source is a living document)
- 1-3 relevant best practices or standards
- Any unresolved ambiguity that should be confirmed with the user

If no technology-specific external guidance is relevant, explicitly state that no latest-standards research was needed.

## Step 3: Initial Scope Assessment

Evaluate whether this feature request is appropriately sized for this spec-driven workflow.

**Private analysis guidance:**

- Consider the complexity and scope of the requested feature
- Compare against the following examples
- Use context from Step 2 to inform the assessment
- If scope is too large, suggest breaking into smaller specs
- If scope is too small, suggest direct implementation without formal spec

**Scope Examples:**

**Too Large (split into multiple specs):**

- Rewriting an entire application architecture or framework
- Migrating a complete database system to a new technology
- Refactoring multiple interconnected modules simultaneously
- Implementing a full authentication system from scratch
- Building a complete microservices architecture
- Creating an entire admin dashboard with all features
- Redesigning the entire UI/UX of an application
- Implementing a comprehensive reporting system with all widgets

**Too Small (vibe-code directly):**

- Adding a single console.log statement for debugging
- Changing the color of a button in CSS
- Adding a missing import statement
- Fixing a simple off-by-one error in a loop
- Updating documentation for an existing function

**Just Right (perfect for this workflow):**

- Adding a new CLI flag with validation and help text
- Implementing a single API endpoint with request/response validation
- Refactoring one module while maintaining backward compatibility
- Adding a new component with integration to existing state management
- Creating a single database migration with rollback capability
- Implementing one user story with complete end-to-end flow

### Report Scope Assessment To User

- **ALWAYS** inform the user of the result of the scope assessment.
- If the scope appears inappropriate, **ALWAYS** pause the conversation to suggest alternatives and get input from the user.

## Step 4: Clarification Sufficiency Check

Assess whether you already have enough aligned context to write a high-quality spec without inventing requirements. Always err on the side of caution, but do not force a questions file when the available information is already sufficient.

Focus on understanding the "what" and "why" rather than the "how."

Use the following common areas to assess whether clarification is needed:

**Core Understanding:**

- What problem does this solve and for whom?
- What specific functionality does this feature provide?

**Success & Boundaries:**

- How will we know it's working correctly?
- What should this NOT do?
- Are there edge cases we should explicitly include or exclude?

**Design & Technical:**

- Any existing design mockups or UI guidelines to follow?
- Are there any technical constraints or integration requirements?

**Proof Artifacts:**

- What proof artifacts will demonstrate this feature works (URLs, CLI output, screenshots)?
- What will each artifact demonstrate about the feature?

**Progressive Disclosure:** Start with Core Understanding, then expand based on feature complexity and user responses.

### Clarification Sufficiency Criteria

Proceed without a questions file only if all of the following are true:

- The user goal and intended outcome are clear.
- Scope boundaries are clear enough to define meaningful non-goals.
- Demoable Units and Proof Artifacts can be specified without guessing.
- Known repository context and user-provided constraints are sufficient to avoid inventing requirements.
- Relevant latest-standards research has been completed for material technologies, and it does not introduce unresolved approach choices that need user confirmation.
- Any remaining uncertainty is minor, non-blocking, and can safely be recorded in the spec's `Open Questions` section without reducing spec quality.

Create a questions file if any of the following are true:

- There are multiple materially different interpretations of the requested feature.
- Acceptance criteria, Proof Artifacts, or Demoable Units would otherwise be guessed.
- Scope boundaries or non-goals are unclear.
- Design, technical, integration, security, or operational constraints are missing and would materially change the spec.
- The user intent or direction could reasonably lead to different implementation paths.
- Current best practices or standards for a relevant technology suggest multiple valid approaches, and the choice would materially affect the spec.
- Repository patterns appear to conflict with current external guidance, and the correct direction is not obvious from the user's request.

### Open Questions Boundary

Open Questions are only for non-blocking uncertainties that do not prevent a high-quality, actionable spec. Use them to document assumptions, deferred nice-to-have details, or follow-up context the user can refine later without changing the core scope, requirements, Demoable Units, Proof Artifacts, or acceptance criteria.

Do not use `Open Questions` to defer material ambiguity that belongs in a questions file. If the answer could materially change the feature scope, implementation path, acceptance criteria, validation strategy, security posture, operational behavior, or proof artifacts, create a questions file instead and stop for user input.

Here are examples of decisions that cannot be left materially unresolved and would merit a questions file:

- Whether the feature is advisory/report-only or can take mutating actions such as publishing, deploying, tagging, approving, merging, notifying, creating releases, dispatching workflows, or changing remote/service state.
- Which execution surface is in scope when materially different options are plausible, such as CLI, CI job, GitHub App, web UI, background service, scheduled agent, or local script.
- Who or what has authority to decide that something is safe enough, ready, approved, blocked, or eligible for automation.
- Which external system of record is primary when several could drive the workflow, such as GitHub, CI/CD, ticketing, Slack, deployment tooling, package registries, or cloud services.
- What credential scope or write permissions are acceptable when the implementation may need access to external systems.

Do not answer a question by punting it to SDD2 or later phases. Phase 1 must either resolve the uncertainty as a clearly labeled assumption, ask the user through a questions file, or record it as a non-blocking Open Question that does not affect implementation planning.

### Clarification Status Declaration (Required)

Before proceeding, you MUST state exactly one of the following:

- `Clarification status: sufficient - no questions file required`
- `Clarification status: insufficient - questions file required`

### Self-Verification Before Proceeding

Before choosing `sufficient`, explicitly verify:

- [ ] I am not guessing at missing requirements.
- [ ] I can populate all major spec sections with grounded, user-aligned content.
- [ ] I have reviewed relevant current best practices for material technologies, or I have explicitly determined that no external standards research is needed.
- [ ] Any remaining uncertainty is non-blocking, does not materially affect the implementation plan, and belongs in `Open Questions` rather than a blocking questions round.

If any check fails, create a questions file.

### Questions File Format

Follow this format exactly when you create a questions file. The questions file is a single self-contained HTML document using the shared SDD-HTML shell (see below), with `TYPE` set to `questions`.

Each question MUST include recommended answer guidance for the user. Recommendations should reduce ambiguity, explain tradeoffs, and bias toward the option that best supports a clear, reviewable, junior-friendly spec.

If a question is driven by latest-standards research, include a short note summarizing the relevant current guidance and why user confirmation is needed.

Use the shared shell/style block from the "HTML Artifact Shell" section below with `TYPE` = `questions`, then populate the body with one `<section class="question">` card per question, in this exact structure:

```html
<!doctype html>
<html lang="en">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1" />
<title>[NN]-questions-[N]-[feature-name]</title>
<style>
  /* ... shared shell CSS from "HTML Artifact Shell" below, byte-identical ... */
  .question { margin: 1.5rem 0; }
  .options { list-style:none; padding-left:0; }
  .options li { padding:.3rem 0; }
  .options li input[type="checkbox"] { margin-right:.5rem; }
  .context-note, .recommend-note { background:var(--card); border:1px solid var(--border); border-radius:.5rem; padding:.75rem 1rem; margin:.5rem 0; }
  .why-list { margin:.5rem 0 0; }
</style>
</head>
<body data-overall-status="pending">
<div class="doc">
  <header class="doc-header">
    <span class="badge badge-type">Questions</span>
    <h1>[NN] Questions Round [N] - [Feature Name]</h1>
    <p class="meta">Sequence [NN] · Feature: [feature-name] · Round [N] · Phase 1</p>
  </header>

  <p>Please answer each question below (select one or more options, or add your own notes). Feel free to add additional context under any question. If you are using Lavish Editor, reply/annotate directly in the interactive session; otherwise edit this file's checkboxes and notes, or answer in chat as instructed.</p>

  <section class="question" id="q1">
    <h2>1. [Question Category/Topic]</h2>
    <p>[What specific aspect of the feature needs clarification?]</p>
    <ul class="options">
      <li><label><input type="checkbox" /> (A) [Option description explaining what this choice means]</label></li>
      <li><label><input type="checkbox" /> (B) [Option description explaining what this choice means]</label></li>
      <li><label><input type="checkbox" /> (C) [Option description explaining what this choice means]</label></li>
      <li><label><input type="checkbox" /> (D) [Option description explaining what this choice means]</label></li>
      <li><label><input type="checkbox" /> (E) Other (describe): ___________</label></li>
    </ul>
    <div class="context-note"><strong>Current best-practice context:</strong> [Optional. Briefly summarize the latest relevant guidance or standard that makes this question important. Omit if not needed.]</div>
    <div class="recommend-note">
      <strong>Recommended answer(s):</strong> [(A), (C)]
      <p><strong>Why these are recommended:</strong></p>
      <ul class="why-list">
        <li>[Recommendation note 1 explaining why the suggested option best preserves user intent, reduces ambiguity, or improves spec quality]</li>
        <li>[Recommendation note 2 explaining tradeoffs versus the other options]</li>
      </ul>
    </div>
  </section>

  <section class="question" id="q2">
    <h2>2. [Another Question Category/Topic]</h2>
    <p>[What specific aspect of the feature needs clarification?]</p>
    <ul class="options">
      <li><label><input type="checkbox" /> (A) [Option description explaining what this choice means]</label></li>
      <li><label><input type="checkbox" /> (B) [Option description explaining what this choice means]</label></li>
      <li><label><input type="checkbox" /> (C) [Option description explaining what this choice means]</label></li>
      <li><label><input type="checkbox" /> (D) [Option description explaining what this choice means]</label></li>
      <li><label><input type="checkbox" /> (E) Other (describe): ___________</label></li>
    </ul>
    <div class="context-note"><strong>Current best-practice context:</strong> [Optional.]</div>
    <div class="recommend-note">
      <strong>Recommended answer(s):</strong> [(B)]
      <p><strong>Why these are recommended:</strong></p>
      <ul class="why-list">
        <li>[Recommendation note 1]</li>
        <li>[Recommendation note 2]</li>
      </ul>
    </div>
  </section>
</div>
</body>
</html>
```

Repeat the `<section class="question">` block for every question in the round. Preserve every field shown above (category heading, prompt, five lettered options including `(E) Other`, the optional best-practice-context note, the recommended-answer line, and the "why these are recommended" list) for each question — none of these are optional decorations, they are the content contract carried over from the Markdown edition.

### Recommendation Rules For Questions Files

When adding recommended answer guidance:

- Recommend the option or combination of options that best supports alignment with the user's likely intent.
- Explain why the recommendation is better than the alternatives presented, not just why it is reasonable in isolation.
- Prefer recommendations that reduce avoidable ambiguity and make the eventual spec easier to validate.
- Do not present recommendations as mandatory; the user remains the decision-maker.
- If no option is clearly safer or more aligned, say so explicitly and explain the tradeoff instead of forcing a weak recommendation.
- If recommending `Other`, provide a short suggested custom answer the user can edit.
- When the question comes from latest-standards research, summarize the relevant current guidance in plain language before recommending an answer.
- Use user confirmation to resolve meaningful tension between repository patterns and current external best practices.

### Example Question With Recommendation Guidance

Use this as a style reference for how to recommend answers without taking the decision away from the user.

```html
<section class="question" id="q1">
  <h2>1. Authentication Entry Point</h2>
  <p>Which sign-in methods should this first version support?</p>
  <ul class="options">
    <li><label><input type="checkbox" /> (A) Email and password only</label></li>
    <li><label><input type="checkbox" /> (B) Google SSO only</label></li>
    <li><label><input type="checkbox" /> (C) Email/password and Google SSO together</label></li>
    <li><label><input type="checkbox" /> (D) Magic link only</label></li>
    <li><label><input type="checkbox" /> (E) Other (describe): ___________</label></li>
  </ul>
  <div class="context-note"><strong>Current best-practice context:</strong> Current guidance for new authentication work often recommends shipping the smallest secure end-to-end slice first, then layering optional auth methods once the base flow is validated.</div>
  <div class="recommend-note">
    <strong>Recommended answer(s):</strong> (A)
    <p><strong>Why these are recommended:</strong></p>
    <ul class="why-list">
      <li><code>(A)</code> is the smallest demoable slice and keeps the first spec focused on one complete authentication path.</li>
      <li><code>(A)</code> reduces ambiguity in validation, proof artifacts, and edge cases compared with <code>(C)</code>, which adds a second auth flow immediately.</li>
      <li><code>(B)</code> and <code>(D)</code> may still be valid product choices, but they introduce external-provider or email-delivery dependencies that are usually unnecessary unless the user explicitly wants them.</li>
      <li>If the user already knows SSO is a hard requirement, they should override this recommendation.</li>
    </ul>
  </div>
</section>
```

### Questions File Process

Only follow this process when clarification is insufficient.

1. **Create Questions File**: Save questions to a file named `[NN]-questions-[N]-[feature-name].html` where `[N]` is the round number (starting at 1, incrementing for each new round).
2. **Augment With Recommendations**: For every question, include recommended answer(s) and short justification notes comparing the recommendation to the other options.
3. **Point User To The Answer Path**: How you point the user to answer depends on Lavish availability (see "Lavish Editor Integration" below):
   - **Lavish available**: open the questions file as an interactive Lavish session and direct the user there to reply/annotate.
   - **Lavish not available**: tell the user the questions file's path for reference, and ask them to answer the numbered questions directly in the chat conversation (since hand-editing raw HTML checkbox markup is not ergonomic) rather than editing the file.
4. **STOP AND WAIT**: Do not proceed to Step 5. Wait for the user to indicate they have provided their answers (via Lavish reply or chat message).
5. **Read Answers**: After the user indicates they have provided their answers, read them (from the Lavish reply payload, or from the chat message) and continue the conversation.
6. **Re-run Sufficiency Check**: Reassess whether the combined context is now sufficient to generate the spec.
7. **Follow-Up Rounds**: If answers reveal new material ambiguity, create a new questions file with incremented round number (`[NN]-questions-[N+1]-[feature-name].html`) and repeat the process (return to step 3).

**Iterative Process:**

- If a user's answer reveals new material questions or areas needing clarification, ask follow-up questions in a new questions file.
- Build on previous answers - use context from earlier responses to inform subsequent questions.
- **CRITICAL**: After creating any questions file, you MUST STOP and wait for the user to provide answers before proceeding.
- Only proceed to Step 5 after:
  - You have received and reviewed all user answers to clarifying questions
  - You have re-run the Clarification Sufficiency Check
  - You have enough detail to populate all spec sections (User Stories, Demoable Units with functional requirements, etc.).

## Step 5: Spec Generation

Generate a comprehensive specification as a single self-contained HTML document using this exact structure. Use the shared shell/style block from the "HTML Artifact Shell" section below with `TYPE` = `spec`, then populate the body with the sections in this exact order, each as an `<h2>`-headed section:

```html
<!doctype html>
<html lang="en">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1" />
<title>[NN]-spec-[feature-name]</title>
<style>
  /* ... shared shell CSS from "HTML Artifact Shell" below, byte-identical ... */
  .unit { margin: 1.25rem 0; }
</style>
</head>
<body data-overall-status="draft">
<div class="doc">
  <header class="doc-header">
    <span class="badge badge-type">Spec</span>
    <h1>[NN]-spec-[feature-name]</h1>
    <p class="meta">Sequence [NN] · Feature: [feature-name] · Phase 1</p>
  </header>

  <section id="introduction">
    <h2>Introduction/Overview</h2>
    <p>[Briefly describe the feature and the problem it solves. State the primary goal in 2-3 sentences.]</p>
  </section>

  <section id="goals">
    <h2>Goals</h2>
    <ul>
      <li>[Specific, measurable objective 1]</li>
      <li>[Specific, measurable objective 2]</li>
      <li>[Specific, measurable objective 3]</li>
      <!-- 3-5 items total -->
    </ul>
  </section>

  <section id="user-stories">
    <h2>User Stories</h2>
    <ul>
      <li><strong>As a</strong> [type of user], <strong>I want to</strong> [perform an action] <strong>so that</strong> [benefit].</li>
      <!-- one per user story -->
    </ul>
  </section>

  <section id="demoable-units">
    <h2>Demoable Units of Work</h2>
    <p>[Focus on tangible progress and WHAT will be demonstrated. Define 2-4 small, end-to-end vertical slices using the format below.]</p>

    <div class="card unit" id="unit-1">
      <h3>Unit 1: [Title]</h3>
      <p><strong>Purpose:</strong> [What this slice accomplishes and who it serves]</p>
      <p><strong>Functional Requirements:</strong></p>
      <ul>
        <li>The system shall [requirement 1: clear, testable, unambiguous]</li>
        <li>The system shall [requirement 2: clear, testable, unambiguous]</li>
        <li>The user shall [requirement 3: clear, testable, unambiguous]</li>
      </ul>
      <p><strong>Proof Artifacts:</strong></p>
      <ul>
        <li>[Artifact type]: [description] demonstrates [what it proves]</li>
        <li>Example: <code>Screenshot: `--help` output demonstrates new command exists</code></li>
        <li>Example: <code>CLI: `command --flag` returns expected output demonstrates feature works</code></li>
      </ul>
    </div>

    <div class="card unit" id="unit-2">
      <h3>Unit 2: [Title]</h3>
      <p><strong>Purpose:</strong> [What this slice accomplishes and who it serves]</p>
      <p><strong>Functional Requirements:</strong></p>
      <ul>
        <li>The system shall [requirement 1: clear, testable, unambiguous]</li>
        <li>The system shall [requirement 2: clear, testable, unambiguous]</li>
      </ul>
      <p><strong>Proof Artifacts:</strong></p>
      <ul>
        <li>[Artifact type]: [description] demonstrates [what it proves]</li>
        <li>Example: <code>Test: MyFeature.test.ts passes demonstrates requirement implementation</code></li>
        <li>Example: <code>Order PDF: PDF downloaded from https://example.com/order-submitted shows completed flow demonstrates end-to-end functionality</code></li>
      </ul>
    </div>
    <!-- repeat unit cards for 2-4 units total -->
  </section>

  <section id="non-goals">
    <h2>Non-Goals (Out of Scope)</h2>
    <p>[Clearly state what this feature will NOT include to manage expectations and prevent scope creep.]</p>
    <ol>
      <li><strong>[Specific exclusion 1]:</strong> description</li>
      <li><strong>[Specific exclusion 2]:</strong> description</li>
      <li><strong>[Specific exclusion 3]:</strong> description</li>
    </ol>
  </section>

  <section id="design-considerations">
    <h2>Design Considerations</h2>
    <p>[Focus on UI/UX requirements and visual design. Link to mockups or describe interface requirements. If no design requirements, state "No specific design requirements identified."]</p>
  </section>

  <section id="repository-standards">
    <h2>Repository Standards</h2>
    <p>[Identify existing patterns and practices that implementation should follow.]</p>
    <ul>
      <li>Coding standards and style guides from the repository</li>
      <li>Architectural patterns and file organization</li>
      <li>Testing conventions and quality assurance practices</li>
      <li>Documentation patterns and commit conventions</li>
      <li>Build and deployment workflows</li>
    </ul>
    <p>[If no specific standards are identified, state "Follow established repository patterns and conventions." instead of the list above.]</p>
  </section>

  <section id="technical-considerations">
    <h2>Technical Considerations</h2>
    <p>[Focus on implementation constraints and HOW it will be built. Mention technical constraints, dependencies, or architectural decisions. Incorporate relevant current best practices or standards discovered during latest-standards research, and call out any explicit deviation that should remain because of repository or user context. If no technical constraints, state "No specific technical constraints identified."]</p>
  </section>

  <section id="security-considerations">
    <h2>Security Considerations</h2>
    <p>[Identify security requirements and sensitive data handling needs.]</p>
    <ul>
      <li>API keys, tokens, and credentials that will be used</li>
      <li>Data privacy and sensitive information handling</li>
      <li>Authentication and authorization requirements</li>
      <li>Proof artifact security (what should NOT be committed)</li>
    </ul>
    <p>[If no specific security considerations, state "No specific security considerations identified." instead of the list above.]</p>
  </section>

  <section id="success-metrics">
    <h2>Success Metrics</h2>
    <p>[How will success be measured? Include specific metrics where possible.]</p>
    <ol>
      <li><strong>[Metric 1]:</strong> with target if applicable</li>
      <li><strong>[Metric 2]:</strong> with target if applicable</li>
      <li><strong>[Metric 3]:</strong> with target if applicable</li>
    </ol>
  </section>

  <section id="open-questions">
    <h2>Open Questions</h2>
    <p>[List only non-blocking questions or assumptions that can be refined later without changing core scope, requirements, Demoable Units, Proof Artifacts, or acceptance criteria. If none, state "No open questions at this time." instead of the list below.]</p>
    <ol>
      <li>[Non-blocking question or assumption 1]</li>
      <li>[Non-blocking question or assumption 2]</li>
    </ol>
  </section>
</div>
</body>
</html>
```

Preserve every section from `Introduction/Overview` through `Open Questions`, in the order shown, with the same field-level guidance (e.g. 3-5 Goals, 2-4 Demoable Units, the exact fallback sentences for empty Design Considerations / Repository Standards / Technical Considerations / Security Considerations / Open Questions). Only the markup changed from Markdown to HTML — the content contract is identical.

## HTML Artifact Shell

Use this exact shared shell/style block for both the questions file and the spec document (only change `<title>`, the badge-type text/label, and the body content — keep the CSS byte-identical to the rest of the sdd-html skill's artifacts):

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

Every generated artifact is a single self-contained HTML5 file: inline `<style>`, no external assets except relative-path local images placed alongside the HTML file (e.g. for future proof-artifact screenshots). It must render correctly when opened directly in a browser with no server. The `<style>` block may be extended with additional rules specific to a doc type (as shown for `.question`, `.options`, `.context-note`, `.recommend-note`, `.why-list`, `.unit` above), but the classes in the shared block above must never be removed or renamed, so all sdd-html artifacts stay visually consistent across phases.

## Lavish Editor Integration

The orchestrating `SKILL.md` performs a one-time Lavish availability check before routing to any phase and reports the result to you as context ("Lavish availability: available | not available"). Use that result at this phase's two natural review/approval checkpoints. Lavish integration is always optional and must never replace the underlying approval/sufficiency gates — a gate is satisfied by explicit user approval (via Lavish reply or chat message), not by merely opening a session.

**Questions file checkpoint** (Step 4's "Questions File Process"):

- **When Lavish is available:** after creating the questions HTML file, open it as an interactive Lavish session so the user can visually review each question and reply/annotate through the browser instead of hand-editing HTML checkbox markup.
  1. `lavish-axi <html-file>` (or `npx -y lavish-axi <html-file>` if no global binary) to open/resume in place — the artifact stays under `docs/specs/...`; do not copy it into `.lavish/`.
  2. `lavish-axi poll <html-file>` to wait for the user's answers/annotations (may run as a background task; never kill it, re-run if it times out).
  3. After reading the reply, `lavish-axi end <html-file>` once answers are captured.
- **When Lavish is NOT available:** keep the original behavior conceptually but adapt the mechanics — since hand-editing raw HTML is not ergonomic, tell the user the questions file's path (for reference/context) and ask them to answer the numbered questions directly in the chat conversation instead of editing the file, then incorporate their answers programmatically before regenerating.
- Either way, **STOP AND WAIT** after creating a questions file/session — this gate is unchanged from the Markdown edition.

**Step 6 Review and Refinement checkpoint** (see below):

- **When Lavish is available:** after generating the spec HTML, open it as an interactive Lavish session (same open/poll/reply-loop/end pattern as above) so the user can annotate specific sections of the spec directly and you can iterate against those annotations.
- **When Lavish is NOT available:** fall back to the original behavior — state the file path, tell the user to open it in a browser, and collect feedback through the normal conversation.

## Step 6: Review and Refinement

### Cross-Domain Applicability Guard (Required)

Before presenting the spec to the user, run this check to keep the workflow
broadly applicable across software tasks (API, UI, CLI, data, infra,
platform):

- [ ] The spec language is domain-neutral (no project-specific assumptions
      unless user-provided).
- [ ] Demoable Units can be validated in at least one of these contexts: API,
      UI, CLI, data pipeline, or infrastructure automation.
- [ ] Proof Artifacts are defined as observable outcomes, not tool-specific
      rituals.
- [ ] Requirements are written so another repository could reuse the structure
      with only context substitutions.

If any item fails, revise wording to be framework-agnostic and context-aware.

After generating the spec, present it to the user and ask them to review the spec as the source of truth for implementation (per the Lavish Editor Integration rules above — either through a Lavish session or by opening the HTML file directly in a browser). Ask for feedback on whether the feature goal, scope, non-goals, requirements, acceptance criteria, demoable units, and proof artifacts match their intent and are specific enough to plan implementation without guessing.

Explicitly call out the `Open Questions` section. If the spec includes open questions, confirm they are non-blocking and ask the user to answer them before continuing or identify which ones can remain as non-blocking assumptions.

Iterate based on feedback until the user is satisfied.

## Output Requirements

**Format:** HTML (`.html`), self-contained per the "HTML Artifact Shell" section above
**Full Path:** `./docs/specs/[NN]-spec-[feature-name]/[NN]-spec-[feature-name].html`
**Example:** For feature "user authentication", the spec directory would be `01-spec-user-authentication/` with a spec file as `01-spec-user-authentication.html` inside it

## Critical Constraints

**NEVER:**

- Start implementing the spec; only create the specification document
- Assume technical details without asking the user
- Create specs that are too large or too small without addressing scope issues
- Use jargon or technical terms that a junior developer wouldn't understand
- Skip the clarification sufficiency check, even if the prompt seems clear
- Ignore existing repository patterns and conventions
- Rely only on stale model knowledge when current external guidance could materially affect the spec
- Generate `.md` artifacts under this skill, or mix Markdown and HTML SDD artifacts in the same spec directory

**ALWAYS:**

- Run the clarification sufficiency check before generating the spec
- Ask clarifying questions when material ambiguity remains
- Research current best practices for material technologies when they could affect the spec
- Validate scope appropriateness before proceeding
- Use the exact spec structure provided above, rendered as self-contained HTML
- Ensure the spec is understandable by a junior developer
- Include proof artifacts for each work unit that demonstrate what will be shown
- Follow identified repository standards and patterns in all requirements
- Incorporate relevant current external standards and best practices in technical guidance
- Offer a Lavish Editor review session at the questions-file and Step 6 checkpoints when Lavish is available, and fall back to file-path-plus-chat review when it is not

## How to Continue the SDD Workflow

Likely next phase action: the skill will route to **Phase 2 — Task List Generation and Planning Audit**, where it will generate parent tasks, request approval before sub-task generation, define proof artifacts, and run the planning audit gate — all as HTML.

To continue the workflow in this chat, reply with:

`Continue SDD with task planning.`

This reply will re-invoke the SDD skill so it can reassess the workspace and route to the next phase.

Before continuing, remind the user to review the spec as the implementation source of truth. If the `Open Questions` section contains anything unresolved, ask the user to answer them before continuing or identify which questions can remain as non-blocking assumptions.

You can also continue in a new chat if you want to keep context lean; the SDD skill will reassess repository state from the persisted spec artifacts in `docs/specs/`.
