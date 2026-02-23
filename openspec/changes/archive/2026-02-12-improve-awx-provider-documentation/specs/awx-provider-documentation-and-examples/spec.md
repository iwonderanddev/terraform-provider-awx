# Delta Specification: awx-provider-documentation-and-examples

## MODIFIED Requirements

### Requirement: Registry-compatible documentation structure

The provider SHALL ship Terraform Registry-compatible documentation for
provider configuration, each resource, and each data source. Generated resource
pages SHALL use a stable section structure compatible with HashiCorp Plugin
Framework documentation-generation guidance, including `Example Usage`,
`Schema`, `Import`, and `Further Reading` sections. Schema qualifier guidance
SHALL be rendered separately from argument entries and MUST NOT appear as a
parameter bullet. Documentation curation and refresh workflows SHALL verify
behavioral accuracy against official AWX 24.6.1 online documentation for each
managed object resource and each managed object data source. Documentation
curation workflows SHALL review official AWX pages that describe both the
object itself and its interactions with related objects before finalizing
examples or parameter definitions.

#### Scenario: Documentation completeness check

- **WHEN** documentation validation is executed
- **THEN** provider, resource, and data source docs are present in the expected
  structure

#### Scenario: Resource schema structure check

- **WHEN** a resource documentation page is generated
- **THEN** it contains `## Example Usage`, `## Schema`, `## Import`, and
  `## Further Reading` sections in consistent order

#### Scenario: Qualifier placement check

- **WHEN** qualifier guidance is rendered for resource arguments
- **THEN** qualifier guidance is presented outside parameter bullet lists

#### Scenario: Reading links check

- **WHEN** a documentation page includes a `## Further Reading` section
- **THEN** the section includes only resource-specific official AWX links for
  behavior details (not only generic AWX index pages)

#### Scenario: Official-link specificity check

- **WHEN** documentation is generated for a managed object resource/data source
- **THEN** `## Further Reading` includes at least one official AWX 24.6.1 link
  mapped to that object's concept page (or closest official concept section)

#### Scenario: Online verification coverage check

- **WHEN** curated docs content is generated or refreshed for managed object
  resources and data sources
- **THEN** each managed object resource/data source has recorded official AWX
  source provenance with an online-verified URL and verification date

#### Scenario: Interaction-informed research check

- **WHEN** documentation is curated for a managed object resource or data source
- **THEN** the curation process includes online review of official AWX
  interaction behavior with related objects and records that provenance for the
  generated examples and parameter descriptions

### Requirement: Operational examples and import guidance

Each resource documentation page SHALL include runnable usage examples and SHALL
define import usage with the accepted ID format. For `awx_job_template`,
`awx_workflow_job_template`, `awx_project`, `awx_credential`,
`awx_inventory`, and `awx_inventory_source`, `## Example Usage` MUST contain
between 1 and 3 examples, inclusive, with supporting references shown when
required for comprehension. Complex resources SHALL include a concise AWX
concept primer and link to official AWX documentation for deeper behavior
details. Field descriptions for prioritized resources SHALL use user-oriented
AWX terminology and MUST NOT use generic placeholder patterns (including
`Managed field from AWX OpenAPI schema`, `Value for`, and `Numeric setting
for`). Documentation SHALL accurately describe typed numeric reference arguments
and attributes, including object `id` as `Number` for collection-created
objects and `String` for detail-path keyed objects. Documentation for
relationship resources SHALL use canonical explicit object-specific `_id`
argument names (for example, `job_template_id`, `credential_id`). Examples
SHALL avoid unnecessary type-conversion workarounds for reference wiring. For
each managed object resource and managed object data source, docs-enrichment
metadata SHALL record curation provenance with an official AWX link and a
verification date in `YYYY-MM-DD` format. Documentation for `awx_setting`
resource and data source SHALL recommend `id = "all"` as the default usage
path, while still documenting category-scoped IDs as optional advanced scoping.
Enum option rendering in schema fields SHALL use valid Markdown list formatting
and MUST NOT emit escaped newline sequences or inline collapsed bullet output.
Examples and parameter definitions for managed object resources/data sources
SHALL be derived from online-verified AWX object behavior and documented
cross-object interactions. At the end of implementation, documentation output
SHALL undergo an explicit quality analysis before sign-off. If the analysis
finds quality gaps, documentation curation and generation SHALL iterate with
re-analysis until quality is acceptable, with a maximum of three total passes.

#### Scenario: Resource documentation review

- **WHEN** a resource documentation page is generated
- **THEN** it contains usage examples, schema content, and an import section
  with correct syntax

#### Scenario: awx_setting default ID guidance

- **WHEN** `awx_setting` resource or data source examples are generated
- **THEN** the canonical default example uses `id = "all"`

#### Scenario: awx_setting import guidance default

- **WHEN** import instructions are generated for `awx_setting`
- **THEN** the primary import example uses `all` as the identifier

#### Scenario: awx_setting scoped usage warning

- **WHEN** documentation describes category-scoped settings IDs
- **THEN** docs describe them as optional advanced scoping and warn that
  overlapping key ownership across `all` and scoped resources can cause
  configuration conflicts

#### Scenario: Prioritized resource example depth check

- **WHEN** documentation is generated for prioritized resources
- **THEN** each page contains 1 to 3 examples in `## Example Usage`

#### Scenario: Complex resource concept primer check

- **WHEN** a resource is classified as complex using documentation complexity,
  parameter count, and cross-resource interaction signals
- **THEN** its documentation includes a concise AWX concept primer and official
  AWX references

#### Scenario: Placeholder description removal check

- **WHEN** documentation is generated for prioritized resources
- **THEN** no field description contains low-information placeholder patterns,
  including `Managed field from AWX OpenAPI schema`, `Value for`, or
  `Numeric setting for`

#### Scenario: Online-grounded curation check

- **WHEN** curated descriptions/examples are updated for managed object
  resources and data sources
- **THEN** the content reflects official AWX 24.6.1 behavior documented in
  linked official AWX pages and docs-enrichment metadata records the
  corresponding official AWX link plus verification date

#### Scenario: Interaction-accurate example and field definition check

- **WHEN** examples and parameter descriptions are curated for managed object
  resources/data sources
- **THEN** they describe object interactions and reference wiring consistent
  with official AWX online behavior documentation

#### Scenario: Reference typing documentation consistency

- **WHEN** a resource exposes numeric reference fields
- **THEN** schema sections describe matching numeric usage and examples
  demonstrate direct assignment without `tonumber(...)`

#### Scenario: Relationship argument naming documentation consistency

- **WHEN** documentation is generated for relationship resources
- **THEN** schema sections use canonical object-specific `_id` names and
  examples avoid generic directional placeholder argument names

#### Scenario: Enum markdown formatting consistency

- **WHEN** schema documentation includes enumerated option values
- **THEN** option lists render as valid Markdown bullets without literal escape
  sequences such as `\n*` and without inline collapsed list formatting

#### Scenario: End-of-implementation quality analysis gate

- **WHEN** documentation implementation work is considered complete
- **THEN** a documented quality analysis is executed across generated provider,
  resource, and data source pages before implementation sign-off

#### Scenario: Bounded quality-iteration loop

- **WHEN** the quality analysis identifies documentation quality gaps
- **THEN** the team performs another documentation-improvement pass and reruns
  the quality analysis until quality is acceptable or three total passes have
  been completed

#### Scenario: Maximum pass limit enforcement

- **WHEN** two remediation passes have already been executed after the initial
  analysis pass
- **THEN** no additional remediation pass is started without an explicit
  follow-up decision, because the process limit of three total passes is
  reached
