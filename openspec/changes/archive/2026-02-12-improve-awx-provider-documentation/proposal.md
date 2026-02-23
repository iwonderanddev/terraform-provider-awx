# Proposal: improve-awx-provider-documentation

## Why

Generated provider documentation is structurally complete, but parts of it still fall short of Terraform Registry quality for practical operator use. We need to improve clarity and example usability now so users can apply resources directly without reverse engineering missing context or malformed schema details.

## What Changes

- Improve generated docs quality for resources and data sources by replacing low-information field descriptions with user-oriented AWX terminology.
- Ensure prioritized resource examples are runnable and include required supporting references where needed for comprehension.
- Require documentation curation and refresh workflows to check official AWX online documentation so each resource and data source remains accurate against AWX 24.6.1 behavior.
- Require per-object online research of AWX behavior and cross-object interactions so examples and parameter definitions reflect how objects actually relate and operate together.
- Fix schema rendering issues for enum values so option lists display correctly in Markdown/Registry output.
- Tighten docs generation and validation behavior to catch malformed enum formatting and placeholder-style description regressions.
- Add a mandatory end-of-implementation documentation quality analysis gate and iterate remediation if needed, with a maximum of three total quality passes.
- Keep current import-ID contracts, typed ID guidance, and `awx_setting` default `id = "all"` guidance intact.

## Capabilities

### New Capabilities

- None.

### Modified Capabilities

- `awx-provider-documentation-and-examples`: Strengthen documentation requirements for example runnability, enum formatting quality, field-description usefulness, and mandatory online verification of object behavior and interactions against official AWX documentation to align with Terraform Registry expectations.

## Impact

- Affected code:
  - `cmd/awxgen/*` documentation generation and docs validation logic.
  - Curated documentation metadata under `internal/manifest/*` (for example docs enrichment and field description overrides).
  - Generated docs under `docs/resources/*.md`, `docs/data-sources/*.md`, and `docs/index.md`.
- APIs/provider behavior: No AWX API or Terraform resource behavior changes.
- Backward compatibility: No breaking runtime change expected; this is documentation-quality and validation-surface improvement.
- Delivery process: Implementation sign-off requires a final documentation-quality analysis with at most three total analysis/remediation passes.
