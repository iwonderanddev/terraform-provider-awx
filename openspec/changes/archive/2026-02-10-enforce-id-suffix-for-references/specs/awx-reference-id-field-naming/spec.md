## ADDED Requirements

### Requirement: Canonical object-link fields SHALL use `_id` suffixes
For generated AWX object resources and data sources, any field that represents a link to another AWX object SHALL be exposed with an explicit `_id` suffix in Terraform schemas.

#### Scenario: Resource reference argument naming is explicit
- **WHEN** a generated resource includes a configurable link field to another AWX object
- **THEN** the Terraform argument name ends in `_id` (for example, `organization_id`)

#### Scenario: Data source reference attribute naming is explicit
- **WHEN** a generated data source includes a computed link field to another AWX object
- **THEN** the Terraform attribute name ends in `_id`

### Requirement: Legacy unsuffixed link fields SHALL NOT be generated
The provider SHALL NOT emit unsuffixed Terraform field names for object-link references once a canonical `_id` field exists for that reference.

#### Scenario: Generated schema omits unsuffixed alias
- **WHEN** generation runs for an object with a known link field
- **THEN** the generated schema includes only the canonical `_id` field name for that link and does not include an unsuffixed duplicate

#### Scenario: Legacy configuration receives actionable diagnostics
- **WHEN** Terraform configuration uses a removed unsuffixed link argument after this change
- **THEN** Terraform reports an unsupported-argument diagnostic that identifies the invalid field name
