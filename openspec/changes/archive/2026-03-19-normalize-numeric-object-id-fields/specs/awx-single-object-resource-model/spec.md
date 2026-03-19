# awx-single-object-resource-model Delta

## MODIFIED Requirements

### Requirement: One AWX object per resource boundary

The provider SHALL define resources so each resource instance manages exactly
one AWX API object and SHALL NOT embed lifecycle management for additional
nested objects. Resource schemas SHALL represent AWX numeric reference fields
using a consistent Terraform number type across both arguments and computed
attributes. Resource `id` SHALL be numeric for collection-created objects and
string for detail-path keyed objects. Any generated AWX field with manifest
type `object` SHALL be represented as a Terraform object value and MUST NOT use
JSON-string configuration or state encoding. Generated resource fields named
`object_id` that semantically carry an AWX numeric object primary key SHALL be
represented as Terraform numbers even when the field itself is optional.

#### Scenario: Parent resource update

- **WHEN** a parent object resource is updated
- **THEN** only the corresponding AWX parent object is modified by that
  resource operation

#### Scenario: Numeric reference field typing is consistent

- **WHEN** a managed object schema includes a numeric foreign-key reference
  field
- **THEN** that field is exposed with matching Terraform type semantics for
  both configuration input and read-back state

#### Scenario: Generated object field typing is consistent

- **WHEN** a managed resource schema includes a generated object field
- **THEN** the field is configured and stored as Terraform object data, not a
  JSON-encoded string

#### Scenario: Semantic numeric object_id resource fields are numeric

- **WHEN** a managed resource exposes an `object_id` field representing the
  numeric primary key of another AWX object
- **THEN** the resource schema uses Terraform number typing for that field
  while still allowing the attribute to be omitted where AWX permits a
  system-scoped assignment
