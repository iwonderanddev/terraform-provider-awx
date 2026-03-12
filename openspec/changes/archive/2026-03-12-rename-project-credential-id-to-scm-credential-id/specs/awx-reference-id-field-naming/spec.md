# Delta Specification: awx-reference-id-field-naming

## MODIFIED Requirements

### Requirement: Canonical object-link fields SHALL use `_id` suffixes

For generated AWX object resources and data sources, any field that represents
a link to another AWX object SHALL be exposed with an explicit `_id` suffix in
Terraform schemas. When the raw AWX field name is too generic to communicate
the relationship purpose, the provider SHALL derive a more specific canonical
Terraform field name that preserves both the related object type and the AWX
interaction context.

#### Scenario: Resource reference argument naming is explicit

- **WHEN** a generated resource includes a configurable link field to another
  AWX object
- **THEN** the Terraform argument name ends in `_id` (for example,
  `organization_id`)

#### Scenario: Data source reference attribute naming is explicit

- **WHEN** a generated data source includes a computed link field to another
  AWX object
- **THEN** the Terraform attribute name ends in `_id`

#### Scenario: Semantic canonical naming resolves ambiguous link fields

- **WHEN** a generated object-link field has a raw AWX name that is too generic
  to describe the purpose of the relationship in Terraform
- **THEN** the provider emits a more specific canonical Terraform field name
  that still identifies the related object type and ends in `_id`

#### Scenario: Project SCM credential naming is explicit

- **WHEN** the provider generates Terraform schema fields for the AWX project
  `credential` reference
- **THEN** the resource and data source expose that field as
  `scm_credential_id`
