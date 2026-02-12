# Delta Specification: awx-provider-documentation-and-examples

## MODIFIED Requirements

### Requirement: Operational examples and import guidance

Each resource documentation page SHALL include runnable usage examples and SHALL
define import usage with the accepted ID format. For
`awx_job_template`, `awx_workflow_job_template`, `awx_project`,
`awx_credential`, `awx_inventory`, and `awx_inventory_source`,
`## Example Usage` MUST contain between 1 and 3 examples, inclusive, with
supporting references shown when required for comprehension. Complex resources
SHALL include a concise AWX concept primer and link to official AWX
documentation for deeper behavior details. Field descriptions for prioritized
resources SHALL use user-oriented AWX terminology and MUST NOT use the generic
placeholder `Managed field from AWX OpenAPI schema`. Documentation SHALL
accurately describe typed numeric reference arguments and attributes, including
object `id` as `Number` for collection-created objects and `String` for
detail-path keyed objects. Documentation for relationship resources SHALL use
canonical explicit object-specific `_id` argument names (for example,
`job_template_id`, `credential_id`). Examples SHALL avoid unnecessary
type-conversion workarounds for reference wiring. For prioritized resources,
docs-enrichment metadata SHALL record curation provenance with an official AWX
link and a verification date in `YYYY-MM-DD` format. Documentation for
`awx_setting` resource and data source SHALL recommend `id = "all"` as the
default usage path, while still documenting category-scoped IDs as optional
advanced scoping.

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
- **THEN** no field description contains
  `Managed field from AWX OpenAPI schema`

#### Scenario: Online-grounded curation check

- **WHEN** curated descriptions/examples are updated for prioritized resources
- **THEN** the content reflects official AWX 24.6.1 behavior documented in
  linked official AWX pages and docs-enrichment metadata records the
  corresponding official AWX link plus verification date

#### Scenario: Reference typing documentation consistency

- **WHEN** a resource exposes numeric reference fields
- **THEN** schema sections describe matching numeric usage and examples
  demonstrate direct assignment without `tonumber(...)`

#### Scenario: Relationship argument naming documentation consistency

- **WHEN** documentation is generated for relationship resources
- **THEN** schema sections use canonical object-specific `_id` names and
  examples avoid generic directional placeholder argument names
