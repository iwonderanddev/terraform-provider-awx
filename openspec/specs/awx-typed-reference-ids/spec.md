# awx-typed-reference-ids Specification

## Purpose
TBD - created by archiving change fix-id-type-mismatch. Update Purpose after archive.
## Requirements
### Requirement: Numeric reference fields SHALL use a consistent Terraform number type
For AWX object, data source, and relationship-resource fields that represent numeric foreign-key references, the provider SHALL expose those fields as Terraform numbers for both configurable arguments and computed attributes. Canonical reference fields SHALL use explicit object-specific `_id` suffix names.

#### Scenario: Wiring a numeric reference from one resource to another
- **WHEN** a resource argument expects a numeric AWX reference field and the user assigns a value from another resource or data source attribute representing that same reference
- **THEN** Terraform type checking succeeds without requiring explicit casting functions

#### Scenario: Canonical numeric reference names are explicit
- **WHEN** the provider generates Terraform schema fields for numeric AWX references
- **THEN** those reference fields use explicit object-specific `_id` suffix naming

#### Scenario: Relationship numeric argument typing is consistent
- **WHEN** a relationship resource is configured using canonical object-specific `*_id` arguments
- **THEN** those arguments accept numeric IDs directly from managed object `id` attributes without conversion

### Requirement: Object identity ID type SHALL match AWX key semantics
For collection-created objects, the provider SHALL expose object and data source `id` as Terraform numbers. For detail-path keyed objects, the provider SHALL expose `id` as Terraform strings.

#### Scenario: Import behavior remains stable
- **WHEN** a user imports an object resource using the existing documented import syntax
- **THEN** the provider accepts the same import ID format and stores state identity using the schema-appropriate type for that object (`Number` for collection-created objects, `String` for detail-path keyed objects)
