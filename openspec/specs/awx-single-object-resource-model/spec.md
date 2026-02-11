# awx-single-object-resource-model Specification

## Purpose
TBD - created by archiving change create-awx-terraform-provider. Update Purpose after archive.
## Requirements
### Requirement: One AWX object per resource boundary
The provider SHALL define resources so each resource instance manages exactly one AWX API object and SHALL NOT embed lifecycle management for additional nested objects. Resource schemas SHALL represent AWX numeric reference fields using a consistent Terraform number type across both arguments and computed attributes. Resource `id` SHALL be numeric for collection-created objects and string for detail-path keyed objects. Any generated AWX field with manifest type `object` SHALL be represented as a Terraform object value and MUST NOT use JSON-string configuration or state encoding.

#### Scenario: Parent resource update
- **WHEN** a parent object resource is updated
- **THEN** only the corresponding AWX parent object is modified by that resource operation

#### Scenario: Numeric reference field typing is consistent
- **WHEN** a managed object schema includes a numeric foreign-key reference field
- **THEN** that field is exposed with matching Terraform type semantics for both configuration input and read-back state

#### Scenario: Generated object field typing is consistent
- **WHEN** a managed resource schema includes a generated object field
- **THEN** the field is configured and stored as Terraform object data, not a JSON-encoded string

### Requirement: No inline nested sub-object lifecycle blocks
The provider SHALL reject designs that manage child object lifecycle via inline nested blocks within parent resources.

#### Scenario: Attempted inline lifecycle modeling
- **WHEN** a resource schema proposal introduces inline nested lifecycle management for child objects
- **THEN** the schema validation process flags the design as unsupported

