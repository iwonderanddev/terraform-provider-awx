# awx-reference-id-field-naming Specification

## Purpose

TBD - created by archiving change enforce-id-suffix-for-references. Update Purpose after archive.

## Requirements

### Requirement: Canonical object-link fields SHALL use `_id` suffixes

For generated AWX object resources and data sources, the provider SHALL expose
any field that represents a link to another AWX object with an explicit `_id`
suffix in Terraform schemas. When the raw AWX field name is too generic to
communicate the relationship purpose, the provider SHALL derive a more specific
canonical Terraform field name that preserves both the related object type and
the AWX interaction context.

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

### Requirement: Legacy unsuffixed link fields SHALL NOT be generated

The provider SHALL NOT emit unsuffixed Terraform field names for object-link references once a canonical `_id` field exists for that reference.

#### Scenario: Generated schema omits unsuffixed alias

- **WHEN** generation runs for an object with a known link field
- **THEN** the generated schema includes only the canonical `_id` field name for that link and does not include an unsuffixed duplicate

#### Scenario: Legacy configuration receives actionable diagnostics

- **WHEN** Terraform configuration uses a removed unsuffixed link argument after this change
- **THEN** Terraform reports an unsupported-argument diagnostic that identifies the invalid field name
