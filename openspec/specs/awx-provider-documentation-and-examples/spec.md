# awx-provider-documentation-and-examples Specification

## Purpose
TBD - created by archiving change create-awx-terraform-provider. Update Purpose after archive.
## Requirements
### Requirement: Registry-compatible documentation structure
The provider SHALL ship Terraform Registry-compatible documentation for provider configuration, each resource, and each data source.

#### Scenario: Documentation completeness check
- **WHEN** documentation validation is executed
- **THEN** provider, resource, and data source docs are present in the expected structure

### Requirement: Operational examples and import guidance
Each resource documentation page SHALL include at least one runnable example and SHALL define import usage with the accepted ID format.

#### Scenario: Resource documentation review
- **WHEN** a resource documentation page is generated
- **THEN** it contains usage examples, argument and attribute references, and an import section with correct syntax

