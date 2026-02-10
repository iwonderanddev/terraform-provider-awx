# awx-full-object-coverage Specification

## Purpose
TBD - created by archiving change create-awx-terraform-provider. Update Purpose after archive.
## Requirements
### Requirement: Managed-object coverage manifest
The provider build process SHALL derive a managed-object manifest from AWX API v2 and SHALL include all configuration-managed AWX object types in scope.

#### Scenario: Manifest generation for AWX 24.6.1
- **WHEN** the manifest is generated from the vendored AWX 24.6.1 schema
- **THEN** each configuration-managed object type appears in the manifest with resource and data source eligibility metadata

### Requirement: Coverage enforcement
The provider SHALL enforce coverage expectations by validating generated resources and data sources against the managed-object manifest.

#### Scenario: Missing object implementation
- **WHEN** a manifest-listed object type has no mapped resource implementation
- **THEN** validation fails with an explicit coverage error identifying the missing object type

