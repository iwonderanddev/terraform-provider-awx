## Context

This change builds a Terraform provider for AWX using the AWX REST API (`/api/v2`) and the vendored OpenAPI schema at `/Users/damien/git/terraform-awx-provider/external/awx-openapi/schema.json`.

Confirmed constraints:

- Compatibility target: AWX `24.6.1` using API v2 only
- Backward compatibility: not required for earlier AWX versions
- Authentication at GA: HTTP Basic only
- Modeling: explicit resources (no inline nested sub-object lifecycle management)
- Relationship handling: explicit relationship resources
- Runtime-only objects: excluded from managed scope
- Sensitive fields: write-only/sensitive handling in Terraform
- Acceptance/e2e testing: local opt-in only (not CI-required)
- Documentation: Terraform Registry-style provider/resource/data source docs with examples
- Naming: keep AWX object naming in resource/data source names for consistency
- Import format: endpoint-aligned identifiers (numeric or detail-key for object resources; composite or parent-key for relationship resources)

This design follows HashiCorp provider design principles around predictable resource boundaries and composable configurations.

## Goals / Non-Goals

**Goals:**

- Build on `terraform-plugin-framework` with a reusable AWX API client layer (Basic auth, pagination, retries, error mapping).
- Deliver resources and data sources for all AWX API-managed configuration objects in scope.
- Map each AWX object to its own Terraform resource and map object links to explicit relationship resources.
- Provide stable CRUD/import/state behavior and deterministic drift handling.
- Normalize optional AWX server-defaulted fields as `Optional + Computed` to prevent null-to-default inconsistencies after apply.
- Enforce sensitive-field safety by avoiding secret value round-tripping in Terraform state.
- Deliver registry-grade documentation and examples for provider/resources/data sources.
- Deliver unit tests and local opt-in acceptance/e2e tests against AWX `24.6.1`.

**Non-Goals:**

- Support for AWX versions older than `24.6.1` in initial GA scope.
- Inline nested lifecycle management of sub-objects inside parent resources.
- Managing runtime-only AWX records as Terraform CRUD resources.
- OAuth2/session/token auth support at GA.
- Mandatory CI-hosted acceptance testing in initial release.

## Decisions

### 1) Provider runtime and architecture

Decision:
- Use `terraform-plugin-framework`.
- Implement shared resource/data source CRUD scaffolding and an object metadata layer.

Rationale:
- Improves maintainability for large object coverage and provides modern provider ergonomics.

Alternatives considered:
- `terraform-plugin-sdk/v2`: viable but less aligned with current provider direction.
- Fully bespoke runtime: unnecessary maintenance burden.

### 2) Resource modeling contract: one API object per resource

Decision:
- Follow one-resource-per-API-object modeling.
- Do not implement inline nested sub-object lifecycle blocks.
- Represent object relationships using dedicated resources and explicit references.

Rationale:
- Predictable plans/state, fewer hidden side effects, and better Terraform graph composability.

Alternatives considered:
- Inline nested lifecycle blocks: rejected for complexity and drift ambiguity.

### 3) API client strategy: generated metadata + handwritten transport

Decision:
- Generate object/field metadata from vendored OpenAPI.
- Keep HTTP transport, auth headers, retries, pagination, and error normalization handwritten.
- Maintain an override registry for schema/runtime mismatches.

Rationale:
- Scales object coverage while preserving control of edge-case behavior.

Alternatives considered:
- Fully handwritten provider: too slow for broad coverage.
- Fully generated provider: brittle on Terraform typing and AWX quirks.

### 4) Authentication scope for GA

Decision:
- Implement HTTP Basic authentication only for GA.

Rationale:
- Matches confirmed requirements and minimizes initial auth complexity.

Alternatives considered:
- Multiple auth modes at GA: unnecessary expansion for initial scope.

### 5) Coverage and exclusions policy

Decision:
- Include all AWX API-managed configuration objects.
- Exclude runtime-only objects via explicit deny-list in a coverage manifest.

Rationale:
- Delivers broad coverage while avoiding resources that do not represent desired-state configuration.

Alternatives considered:
- Common-object-only scope: rejected.
- Blind inclusion of every endpoint: would include non-manageable runtime records.

### 6) Sensitive field handling policy

Decision:
- Mark secrets/password-like inputs as sensitive.
- Treat secret values as write-only when possible and do not repopulate cleartext into state on read.

Rationale:
- Aligns with provider security best practices and avoids state leakage risks.

Alternatives considered:
- Round-trip secret fields from API reads: rejected due security and drift-noise risk.

### 7) Documentation model

Decision:
- Use Terraform Registry-compatible docs structure.
- Each resource/data source doc includes:
  - Minimal example
  - Argument/attribute reference
  - Import example and ID format (resources)

Rationale:
- Provider usability and operability require first-class docs from initial release.

Alternatives considered:
- Deferring docs until post-MVP: rejected.

### 8) Test strategy and execution model

Decision:
- Unit tests cover client behavior, schema mapping, CRUD scaffolding, import/state normalization, and relationship resources.
- Acceptance/e2e tests are local opt-in via environment variables and skipped by default in CI.
- Acceptance tests include Terraform-driven scenarios using `terraform-plugin-testing` to exercise provider resources/data sources through plan/apply/import lifecycle flows (not client-only calls).

Rationale:
- Preserves CI reliability while allowing true integration validation.

Alternatives considered:
- CI-mandatory e2e: rejected due environment availability constraints.

### 9) Compatibility contract

Decision:
- Compatibility target is AWX `24.6.1` API v2 only.
- No backward compatibility commitment for older AWX versions in GA scope.

Rationale:
- Establishes clear, testable support boundaries.

Alternatives considered:
- Broad multi-version promise: higher risk before baseline stabilization.

### 10) Naming convention contract

Decision:
- Keep AWX names in Terraform resource/data source identifiers for consistency.

Rationale:
- Preserves a direct mapping between AWX API/domain terms and provider objects, reducing translation ambiguity.

Alternatives considered:
- Normalizing/renaming AWX terms for Terraform ergonomics: rejected to avoid drift between AWX docs/API and provider names.

### 11) Import ID contract

Decision:
- Use endpoint-aligned identifiers for object resources.
- Use numeric AWX IDs for standard object resources with numeric identity.
- Use detail-path identifiers for singleton/detail-key object resources (for example `settings` category slugs).
- Use composite IDs for association relationship resources that bind parent and child objects.
- Use parent-key identifiers for singleton relationship resources that are scoped only by parent identity (for example survey specification subresources).

Examples:
- `terraform import awx_project.main 42`
- `terraform import awx_setting.main system`
- `terraform import awx_team_user_membership.main 12:34`
- `terraform import awx_job_template_survey_spec.main 12`

Rationale:
- Matches AWX endpoint identity semantics while keeping resource and relationship imports deterministic.

Alternatives considered:
- Composite IDs for all resources: unnecessary verbosity for object resources and incompatible with singleton/detail-key endpoints.
- Numeric IDs for relationship resources only: not possible where association endpoints require both parent and child identity.

### 12) Server-default field normalization contract

Decision:
- Parse OpenAPI `default` values from request schemas.
- Mark optional non-write-only fields with OpenAPI defaults as Terraform `Computed` (while preserving `Optional`).
- Add targeted acceptance scenarios for omitted defaulted fields to verify create + plan-only + import stability.

Rationale:
- AWX frequently returns explicit default values even when fields are omitted in configuration.
- Without computed normalization, Terraform can fail with provider inconsistency errors (planned `null`, observed concrete default).

Alternatives considered:
- Manual per-field overrides only: too brittle and reactive for broad object coverage.
- Keep fields optional-only and special-case in runtime state conversion: higher complexity and less transparent schema behavior.

## Risks / Trade-offs

- [Coverage scale] Broad object coverage increases implementation/test surface. -> Mitigation: generator-first workflow plus manifest-based coverage tracking.
- [Schema/runtime mismatch] OpenAPI may diverge from behavior of certain endpoints. -> Mitigation: override registry and fixture-backed validation against AWX `24.6.1`.
- [Default inference overreach] OpenAPI defaults may not always represent desired computed semantics for every field. -> Mitigation: explicit field overrides and regression acceptance tests for high-risk objects.
- [Runtime exclusion mistakes] Misclassification could hide useful objects or include invalid ones. -> Mitigation: explicit inclusion/exclusion manifest and review checks.
- [Local-only acceptance] No CI integration signal may delay detection of regressions. -> Mitigation: strong unit tests and documented local pre-release acceptance workflow.

## Migration Plan

1. Scaffold provider module (`terraform-plugin-framework`, provider config, diagnostics, shared API client).
2. Implement HTTP Basic auth and core API behavior (timeouts, retries, pagination, error mapping).
3. Build OpenAPI ingestion and managed-object manifest generation.
4. Add OpenAPI default-value inference and optional-field computed normalization rules.
5. Generate baseline resources/data sources and wire explicit relationship resources.
6. Apply runtime-only object exclusions and sensitive-field schema policies.
7. Implement import/state normalization and consistent read/drift behavior.
8. Generate and refine registry-style docs with examples and import sections.
9. Implement unit tests and local opt-in acceptance harness targeting AWX `24.6.1`.
10. Validate against a real AWX instance and publish compatibility/known-limitations notes.

Rollback strategy:
- Temporarily disable unstable object categories via manifest configuration.
- Avoid destructive state schema migrations in early releases.

## Open Questions

- None at this stage.
