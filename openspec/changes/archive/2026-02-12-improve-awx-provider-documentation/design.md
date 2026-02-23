# Design: improve-awx-provider-documentation

## Context

The proposal targets documentation quality gaps in generated Terraform docs,
while preserving provider runtime behavior and ID/import contracts. The current
pipeline already generates complete section scaffolding, but quality issues
remain in three areas:

- Some examples are not directly runnable because required supporting resources
  are referenced but not shown.
- Enum rendering is inconsistent in several pages, including inline list output
  and escaped newline artifacts.
- Many fields still use low-information fallback wording.

This change also adds a stricter accuracy requirement: documentation curation
must verify resource and data source behavior against official AWX online
documentation (24.6.1), including how each object interacts with related AWX
objects.

Constraints:

- Keep AWX compatibility target at 24.6.1.
- Preserve existing provider behavior, import ID contracts, and typed ID
  semantics.
- Apply improvements through generator logic and curated docs metadata, then
  regenerate docs.

## Goals / Non-Goals

**Goals:**

- Require online verification against official AWX documentation for each
  managed object resource and data source.
- Require object-by-object online interaction analysis so examples and
  parameter definitions reflect real AWX relationship behavior.
- Improve doc clarity by replacing weak fallback text with user-oriented AWX
  terminology.
- Ensure prioritized resource examples are runnable and complete enough for
  comprehension.
- Normalize enum formatting so Markdown renders consistently in Terraform
  Registry pages.
- Add validation gates that prevent regressions in example quality,
  field-description quality, and online-verification metadata.
- Enforce an end-of-implementation quality-analysis loop with at most three
  total passes (initial analysis plus up to two remediation passes).

**Non-Goals:**

- Changing Terraform provider runtime behavior or AWX API semantics.
- Rewriting every relationship resource description into full tutorials.
- Depending on always-online CI execution for standard test jobs.

## Decisions

### 1. Use official AWX docs as the required behavior source for every object doc

Decision:

- Extend docs-enrichment metadata so every managed object doc entry includes:
  - `officialAwxUrl`
  - `verifiedOn` (`YYYY-MM-DD`)
- Treat these fields as required for object resources and object data sources.
- Curate and refresh content by checking the linked official AWX page online
  before updating descriptions/examples.
- Require curation to review interaction-relevant official pages for each
  object (for example related concept pages for associations, prompts, and
  ownership) before finalizing examples and parameter definitions.

Rationale:

- This makes verification explicit and auditable instead of implied.
- It supports the new requirement that each resource/data source remains
  accurate against official AWX behavior.
- It ensures documentation reflects not only isolated object fields but also
  cross-object interactions that users must model in Terraform.

Alternatives considered:

- Manual reviewer memory without recorded provenance: rejected as non-auditable.
- Non-AWX references as primary source: rejected to keep behavior authority
  aligned with AWX documentation.

### 2. Keep online checks out of default CI, but provide an explicit verification step

Decision:

- Keep `make docs-validate` deterministic and offline.
- Add a dedicated online verification workflow/command used during curation and
  release prep to confirm official links resolve and match expected AWX
  concepts.

Rationale:

- Network-dependent checks in default CI can create flaky pipelines.
- A dedicated verification step still enforces online validation at the right
  checkpoints.

Alternatives considered:

- Mandatory online checks in every CI run: rejected due to reliability risk.

### 3. Enforce runnable-example quality for prioritized resources

Decision:

- Keep 1-3 examples for prioritized resources.
- Require examples to either:
  - include needed supporting resources/data blocks, or
  - clearly state prerequisite references when full scaffolding would reduce
    readability.
- Add validation for broken example patterns (for example unresolved
  `awx_*.<name>.id` references without context).
- Require curated examples to map referenced `_id` parameters to documented AWX
  interaction semantics verified online.

Rationale:

- Terraform Registry users need examples that can be applied with minimal
  interpretation.
- Interaction-aware examples reduce misleading configurations where references
  compile but do not represent correct AWX behavior.

Alternatives considered:

- Keep minimal snippets with implicit prerequisites: rejected as too ambiguous.

### 4. Normalize schema enum rendering and ban malformed output patterns

Decision:

- Render enum choices using canonical multiline bullet formatting.
- Prohibit escaped newline fragments (for example literal `\n*`) and inline
  mixed bullet output in field lines.
- Add docs validation checks for malformed enum formatting patterns.

Rationale:

- Consistent rendering improves scanability and Registry presentation quality.

Alternatives considered:

- Leave enum rendering as-is and rely on manual cleanup: rejected as
  non-scalable.

### 5. Replace low-information fallback text with typed contextual descriptions

Decision:

- Keep description precedence:
  1. curated docs-enrichment override
  2. OpenAPI/manifest description
  3. typed contextual fallback phrasing
- Add quality gates for known low-value phrases (`Value for`, `Numeric setting
  for`, and similar generic patterns) in targeted docs scope.
- Require curated parameter descriptions for relationship-driving fields to
  explain interaction intent (for example owner object, related object, and
  launch/prompt behavior) based on official AWX docs.

Rationale:

- Placeholder-like text is a major gap against Terraform Registry quality
  expectations.
- Interaction-oriented wording reduces ambiguity for fields whose meaning
  depends on other AWX objects.

Alternatives considered:

- OpenAPI-only descriptions everywhere: rejected because source text quality is
  inconsistent.

### 6. Add a bounded end-of-implementation quality loop

Decision:

- Require a formal documentation quality analysis after implementation updates
  are generated.
- If quality is insufficient, run another remediation pass and reanalyze.
- Cap the process at three total passes (initial pass plus at most two
  remediation passes).
- If quality is still insufficient after pass three, stop and record follow-up
  work rather than iterating indefinitely.

Rationale:

- This creates an explicit quality gate and prevents “done” states based only
  on generation success.
- The pass cap keeps delivery predictable and avoids open-ended churn.

Alternatives considered:

- Unlimited iterative passes: rejected due to schedule risk and unclear stop
  conditions.

## Risks / Trade-offs

- [Risk] Online docs may change URL structure or section anchors.
  → Mitigation: keep explicit mapping in metadata and refresh during AWX schema
  updates.
- [Risk] Broad quality gates may block generation during transition.
  → Mitigation: phase-in strict checks where needed, then expand once curated
  metadata coverage is complete.
- [Risk] Runnable examples can become verbose.
  → Mitigation: keep a bounded example count and prefer concise prerequisite
  notes when full scaffolding is unnecessary.
- [Risk] Extra curation metadata increases maintenance overhead.
  → Mitigation: validate metadata shape and keep a single source file for
  enrichment.

## Migration Plan

1. Update the capability delta spec (`awx-provider-documentation-and-examples`)
   with explicit online-verification and formatting requirements.
2. Expand docs-enrichment metadata to include required official AWX link and
   verification date coverage for managed object docs.
3. Add a repeatable curation process that records per-object interaction notes
   from official AWX online docs and links those notes to examples/field
   descriptions.
4. Implement generator updates for description fallbacks, enum formatting, and
   example rendering guidance.
5. Add/extend docs validation rules for:
   - malformed enum rendering
   - low-information description patterns
   - example quality expectations
   - required official AWX provenance metadata
   - interaction-aware example/parameter coverage for curated objects
6. Add a dedicated online verification command/workflow for curation/release
   checkpoints.
7. Regenerate docs and run: `make docs`, `make docs-validate`, `make test`.
8. Execute end-of-implementation quality analysis pass 1; if needed, perform
   remediation and rerun analysis for passes 2 and 3 only.

## Open Questions

- Should online verification coverage include relationship resources as a strict
  requirement in this change, or remain focused on managed object
  resources/data sources first?
- What staleness policy should trigger re-verification (`verifiedOn`) outside
  schema refresh events?
