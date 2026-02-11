# awx-data-sources-and-lookups Specification

## Purpose
TBD - created by archiving change create-awx-terraform-provider. Update Purpose after archive.
## Requirements
### Requirement: Data sources for managed object lookup
The provider SHALL expose data sources for managed AWX object types needed to reference existing infrastructure in Terraform configurations. Data source attributes that represent numeric AWX reference identifiers SHALL use Terraform number types that are directly consumable by resource arguments expecting the same references. Data source `id` SHALL be `Number` for collection-created objects and `String` for detail-path keyed objects. Data source attributes for AWX fields typed as `object` SHALL be returned as Terraform object values rather than JSON strings.

#### Scenario: Resolve existing object by filter
- **WHEN** a user configures a data source with valid lookup arguments
- **THEN** the provider returns a deterministic single object match or a clear ambiguity diagnostic

#### Scenario: Data source reference value is directly reusable
- **WHEN** a data source returns a numeric reference attribute and a resource argument consumes that same reference
- **THEN** the configuration can connect the values without explicit type conversion functions

#### Scenario: Data source ID follows object key type
- **WHEN** a user reads a collection-created object through a data source
- **THEN** the returned `id` is numeric and can be passed directly to numeric resource arguments

#### Scenario: Data source object fields are not JSON strings
- **WHEN** a data source reads a managed object field with object semantics
- **THEN** the returned attribute value is Terraform object data, not JSON-encoded text

#### Scenario: Job template extra_vars response normalization
- **WHEN** job template data source reads `extra_vars` as JSON or YAML string content from AWX
- **THEN** the provider normalizes it to a Terraform object value or returns a clear diagnostic when the normalized root is not an object

### Requirement: Stable data source query behavior
Data source queries SHALL implement normalized filtering and deterministic result handling consistent with AWX API v2 semantics.

#### Scenario: No matching object
- **WHEN** a data source lookup returns zero matches
- **THEN** the provider returns a not-found diagnostic that identifies the lookup criteria

