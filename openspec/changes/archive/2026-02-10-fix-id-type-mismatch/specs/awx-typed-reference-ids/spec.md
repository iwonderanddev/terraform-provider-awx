## ADDED Requirements

### Requirement: Numeric reference fields SHALL use a consistent Terraform number type
For AWX object and data source fields that represent numeric foreign-key references, the provider SHALL expose those fields as Terraform numbers for both configurable arguments and computed attributes.

#### Scenario: Wiring a numeric reference from one resource to another
- **WHEN** a resource argument expects a numeric AWX reference field and the user assigns a value from another resource or data source attribute representing that same reference
- **THEN** Terraform type checking succeeds without requiring explicit casting functions

### Requirement: Object identity ID type SHALL match AWX key semantics
For collection-created objects, the provider SHALL expose object and data source `id` as Terraform numbers. For detail-path keyed objects, the provider SHALL expose `id` as Terraform strings.

#### Scenario: Import behavior remains stable
- **WHEN** a user imports an object resource using the existing documented import syntax
- **THEN** the provider accepts the same import ID format and stores state identity using the schema-appropriate type for that object (`Number` for collection-created objects, `String` for detail-path keyed objects)
