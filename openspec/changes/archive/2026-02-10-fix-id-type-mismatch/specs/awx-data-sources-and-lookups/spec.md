## MODIFIED Requirements

### Requirement: Data sources for managed object lookup
The provider SHALL expose data sources for managed AWX object types needed to reference existing infrastructure in Terraform configurations. Data source attributes that represent numeric AWX reference identifiers SHALL use Terraform number types that are directly consumable by resource arguments expecting the same references. Data source `id` SHALL be `Number` for collection-created objects and `String` for detail-path keyed objects.

#### Scenario: Resolve existing object by filter
- **WHEN** a user configures a data source with valid lookup arguments
- **THEN** the provider returns a deterministic single object match or a clear ambiguity diagnostic

#### Scenario: Data source reference value is directly reusable
- **WHEN** a data source returns a numeric reference attribute and a resource argument consumes that same reference
- **THEN** the configuration can connect the values without explicit type conversion functions

#### Scenario: Data source ID follows object key type
- **WHEN** a user reads a collection-created object through a data source
- **THEN** the returned `id` is numeric and can be passed directly to numeric resource arguments
