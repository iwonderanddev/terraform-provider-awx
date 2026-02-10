# awx-runtime-object-exclusions Specification

## Purpose
TBD - created by archiving change create-awx-terraform-provider. Update Purpose after archive.
## Requirements
### Requirement: Runtime-only object exclusion policy
The provider SHALL exclude runtime-only AWX objects from managed Terraform resources and SHALL document excluded object categories.

#### Scenario: Runtime endpoint classification
- **WHEN** an AWX object type is classified as runtime-only with no desired-state lifecycle semantics
- **THEN** no managed resource is generated for that object type

### Requirement: Exclusion manifest enforcement
The provider build and validation pipeline SHALL enforce explicit exclusion entries for runtime-only object categories.

#### Scenario: Undeclared runtime exclusion
- **WHEN** a runtime-only object appears in generated resource candidates without an exclusion entry
- **THEN** validation fails and reports the required exclusion manifest update

