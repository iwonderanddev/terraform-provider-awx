# awx-data-sources-and-lookups Specification

## Purpose
TBD - created by archiving change create-awx-terraform-provider. Update Purpose after archive.
## Requirements
### Requirement: Data sources for managed object lookup
The provider SHALL expose data sources for managed AWX object types needed to reference existing infrastructure in Terraform configurations.

#### Scenario: Resolve existing object by filter
- **WHEN** a user configures a data source with valid lookup arguments
- **THEN** the provider returns a deterministic single object match or a clear ambiguity diagnostic

### Requirement: Stable data source query behavior
Data source queries SHALL implement normalized filtering and deterministic result handling consistent with AWX API v2 semantics.

#### Scenario: No matching object
- **WHEN** a data source lookup returns zero matches
- **THEN** the provider returns a not-found diagnostic that identifies the lookup criteria

