## MODIFIED Requirements

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
