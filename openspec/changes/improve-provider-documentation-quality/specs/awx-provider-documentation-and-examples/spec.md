# Delta Specification: awx-provider-documentation-and-examples

## MODIFIED Requirements

### Requirement: Registry-compatible documentation structure

The provider SHALL ship Terraform Registry-compatible documentation for
provider configuration, each resource, and each data source. Generated resource
pages SHALL use a stable section structure compatible with HashiCorp Plugin
Framework documentation-generation guidance, including `Example Usage`,
`Schema`, `Import`, and `Further Reading` sections. Schema qualifier guidance
SHALL be rendered separately from argument entries and MUST NOT appear as a
parameter bullet.

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
- **THEN** the section includes official AWX links for behavior details and
  HashiCorp/AWS references for documentation style guidance

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
type-conversion workarounds for reference wiring.

#### Scenario: Resource documentation review

- **WHEN** a resource documentation page is generated
- **THEN** it contains usage examples, schema content, and an import section
  with correct syntax

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

#### Scenario: Reference typing documentation consistency

- **WHEN** a resource exposes numeric reference fields
- **THEN** schema sections describe matching numeric usage and examples
  demonstrate direct assignment without `tonumber(...)`

#### Scenario: Relationship argument naming documentation consistency

- **WHEN** documentation is generated for relationship resources
- **THEN** schema sections use canonical object-specific `_id` names and
  examples avoid generic directional placeholder argument names
