# awx-single-object-resource-model Specification

## Purpose
TBD - created by archiving change create-awx-terraform-provider. Update Purpose after archive.
## Requirements
### Requirement: One AWX object per resource boundary
The provider SHALL define resources so each resource instance manages exactly one AWX API object and SHALL NOT embed lifecycle management for additional nested objects.

#### Scenario: Parent resource update
- **WHEN** a parent object resource is updated
- **THEN** only the corresponding AWX parent object is modified by that resource operation

### Requirement: No inline nested sub-object lifecycle blocks
The provider SHALL reject designs that manage child object lifecycle via inline nested blocks within parent resources.

#### Scenario: Attempted inline lifecycle modeling
- **WHEN** a resource schema proposal introduces inline nested lifecycle management for child objects
- **THEN** the schema validation process flags the design as unsupported

