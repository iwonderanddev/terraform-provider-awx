# Proposal: Improve Provider Documentation Quality

## Why

Provider documentation quality is currently too low for effective onboarding and
reliable day-to-day use. Many generated descriptions are generic placeholders
(for example, `Managed field from AWX OpenAPI schema`), examples are too sparse
for realistic workflows, and at least one formatting bug
(`Argument qualifiers used below`) reduces readability.

This should be corrected now to align with Terraform Plugin Framework
documentation-generation expectations and to move closer to the quality bar of
mature providers.

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
- Require each generated resource/data-source page to include at least one
  resource-specific official AWX 24.6.1 link for that exact AWX concept
  (for example, projects -> AWX Projects docs), rather than generic AWX index
  links.
- Verify curated descriptions and examples against official AWX 24.6.1
  documentation pages before finalizing generated content.
- Record curation provenance for prioritized resources in docs-enrichment
  metadata using official AWX links and verification dates.
- Fix the documentation formatting bug where `Argument qualifiers used below`
  appears in the same list as parameters; place qualifier guidance in the
  correct dedicated section.
- Align generated documentation structure and section ordering with HashiCorp
  Terraform Plugin Framework documentation-generation guidance.
- Align writing style and example clarity with established Terraform provider
  documentation quality patterns.
- Add validation coverage to prevent regressions in placeholder descriptions,
  example completeness, qualifier section formatting, and resource-specific
  official AWX link policy, with quality gates enforced first on the
  prioritized resources.

## Capabilities

### New Capabilities

- None.

### Modified Capabilities

- `awx-provider-documentation-and-examples`: Expand requirements to enforce
  high-quality field descriptions, contextual real-world examples (up to 3 when
  needed), concise AWX concept primers on complex resources, official AWX link
  policy, curation provenance metadata, correct qualifier rendering/placement,
  and HashiCorp-aligned generated documentation
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
  `internal/manifest/field_overrides.json`,
  `internal/manifest/docs_enrichment.json`, documentation generation logic in
  `cmd/awxgen`, and docs validation behavior.
- External references used as content/style constraints: official AWX 24.6.1
  user documentation and HashiCorp Plugin Framework
  documentation-generation guidance.
- Curation process requirement: verify examples and terminology against live
  official AWX documentation pages for version 24.6.1.
