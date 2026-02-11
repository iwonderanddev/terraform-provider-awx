## MODIFIED Requirements

### Requirement: One AWX object per resource boundary
The provider SHALL define resources so each resource instance manages exactly one AWX API object and SHALL NOT embed lifecycle management for additional nested objects. Resource schemas SHALL represent AWX numeric reference fields using a consistent Terraform number type across both arguments and computed attributes. Resource link/reference field names SHALL use explicit `_id` suffixes for both configurable and computed schema fields, and unsuffixed link aliases SHALL NOT be generated. Resource `id` SHALL be numeric for collection-created objects and string for detail-path keyed objects.

#### Scenario: Parent resource update
- **WHEN** a parent object resource is updated
- **THEN** only the corresponding AWX parent object is modified by that resource operation

#### Scenario: Numeric reference field typing is consistent
- **WHEN** a managed object schema includes a numeric foreign-key reference field
- **THEN** that field is exposed with matching Terraform type semantics for both configuration input and read-back state

#### Scenario: Resource reference field naming is explicit
- **WHEN** a managed object schema includes a link field to another AWX object
- **THEN** the Terraform schema exposes that link field using a canonical `_id` suffix name
