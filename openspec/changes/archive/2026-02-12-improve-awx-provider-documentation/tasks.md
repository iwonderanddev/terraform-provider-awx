# Tasks: improve-awx-provider-documentation

## 1. Docs Metadata Coverage

- [x] 1.1 Extend docs-enrichment schema/validation so every managed object resource and data source requires `officialAwxUrl` and `verifiedOn`.
- [x] 1.2 Backfill missing managed-object metadata entries with official AWX 24.6.1 URLs and verification dates.
- [x] 1.3 Add/adjust generator tests for metadata parsing and required provenance fields.

## 2. Online Verification Workflow

- [x] 2.1 Implement a dedicated online docs verification command/path that validates official AWX links and concept mappings.
- [x] 2.2 Ensure default offline docs validation remains deterministic and does not require network access.
- [x] 2.3 Add tests for online verification behavior, including failure modes for missing or invalid official links.
- [x] 2.4 Define and enforce a per-object online research checklist that covers object behavior, related-object interactions, and parameter semantics from official AWX docs.

## 3. Example Quality Improvements

- [x] 3.1 Update docs generation rules to enforce runnable prioritized resource examples with required supporting context.
- [x] 3.2 Curate prioritized resource examples (`awx_job_template`, `awx_workflow_job_template`, `awx_project`, `awx_credential`, `awx_inventory`, `awx_inventory_source`) to meet 1-3 example constraints, runnability guidance, and verified interaction behavior.
- [x] 3.3 Add validation checks for unresolved cross-resource references in generated examples.
- [x] 3.4 Add validation checks that curated examples include interaction-correct reference wiring derived from official AWX online docs.

## 4. Enum Rendering Normalization

- [x] 4.1 Normalize enum rendering in generator output to canonical multiline Markdown bullet formatting.
- [x] 4.2 Add validation to fail docs generation on malformed enum output patterns (for example literal `\n*` or inline collapsed bullets).
- [x] 4.3 Add regression tests covering representative enum fields across resources and data sources.

## 5. Field Description Quality Gates

- [x] 5.1 Update description fallback logic to prefer curated text, then OpenAPI text, then typed contextual wording.
- [x] 5.2 Add validation gates for low-information placeholder patterns (`Managed field from AWX OpenAPI schema`, `Value for`, `Numeric setting for`).
- [x] 5.3 Curate high-impact field descriptions in targeted docs where generic wording remains.
- [x] 5.4 Curate relationship-driving parameter descriptions to explain cross-object behavior using online-verified AWX interaction documentation.

## 6. Regeneration and Validation

- [x] 6.1 Run `make generate` and `make validate-manifest` after metadata and generator updates.
- [x] 6.2 Run `make docs` and `make docs-validate` to regenerate and verify documentation outputs.
- [x] 6.3 Run `make test` and fix any regressions in unit/integration coverage.

## 7. Change Verification and Readiness

- [x] 7.1 Confirm the updated docs satisfy the modified `awx-provider-documentation-and-examples` scenarios.
- [x] 7.2 Review generated docs for prioritized resources/data sources and spot-check non-prioritized pages for formatting consistency.
- [x] 7.3 Prepare a concise implementation summary with validation evidence for review.

## 8. Quality Analysis Pass Loop (Max 3)

- [x] 8.1 Run end-of-implementation documentation quality analysis pass 1 and record findings.
- [x] 8.2 Pass 2 not required; pass 1 quality analysis was sufficient.
- [x] 8.3 Pass 3 not required; pass 1 quality analysis was sufficient.
- [x] 8.4 Fourth pass follow-up trigger not reached because pass 1 was sufficient.
