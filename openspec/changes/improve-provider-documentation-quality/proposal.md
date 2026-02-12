# Proposal: Improve Provider Documentation Quality

## Why

Provider documentation quality is currently too low for effective onboarding and
reliable day-to-day use. Many generated descriptions are generic placeholders
(for example, `Managed field from AWX OpenAPI schema`), examples are too sparse
for realistic workflows, and at least one formatting bug
(`Argument qualifiers used below`) reduces readability.

This should be corrected now to align with Terraform Plugin Framework
documentation-generation expectations and to move closer to the quality bar of
mature providers, using the official AWS provider docs as a concrete benchmark.

## What Changes

- Improve generated argument and attribute descriptions by replacing low-value
  placeholder text with curated, user-focused explanations grounded in official
  AWX 24.6.1 documentation.
- Add practical, minimal-but-complete examples for the first-priority
  resources: `awx_job_template`, `awx_workflow_job_template`, `awx_project`,
  `awx_credential`, `awx_inventory`, and `awx_inventory_source`.
- Support multiple examples when needed, with a maximum of 3 examples per
  documentation page.
- Add short AWX concept primers only on complex resources, with links to
  official AWX docs for deeper details.
- Fix the documentation formatting bug where `Argument qualifiers used below`
  appears in the same list as parameters; place qualifier guidance in the
  correct dedicated section.
- Align generated documentation structure and section ordering with HashiCorp
  Terraform Plugin Framework documentation-generation guidance.
- Align writing style and example clarity with established Terraform provider
  documentation quality patterns, using HashiCorp guidance and official AWS
  provider docs as concrete references.
- Add validation coverage to prevent regressions in placeholder descriptions,
  example completeness, and qualifier section formatting, with quality gates
  enforced first on the prioritized resources.

## Capabilities

### New Capabilities

- None.

### Modified Capabilities

- `awx-provider-documentation-and-examples`: Expand requirements to enforce
  high-quality field descriptions, contextual real-world examples (up to 3 when
  needed), concise AWX concept primers on complex resources, official AWX links
  plus HashiCorp/AWS quality references, correct qualifier
  rendering/placement, and HashiCorp-aligned generated documentation
  structure.

## Impact

- Affected generated outputs: `docs/index.md`, `docs/resources/*.md`, and
  `docs/data-sources/*.md`.
- First rollout scope:
  `docs/resources/awx_job_template.md`,
  `docs/resources/awx_workflow_job_template.md`,
  `docs/resources/awx_project.md`,
  `docs/resources/awx_credential.md`,
  `docs/resources/awx_inventory.md`, and
  `docs/resources/awx_inventory_source.md`.
- Likely affected curated/generated inputs and tooling:
  `internal/manifest/field_overrides.json`, documentation generation logic in
  `cmd/awxgen`, and docs validation behavior.
- External references used as content/style constraints: official AWX 24.6.1
  user documentation, HashiCorp Plugin Framework documentation-generation
  guidance, and official AWS provider docs as a documentation quality/style
  benchmark.
