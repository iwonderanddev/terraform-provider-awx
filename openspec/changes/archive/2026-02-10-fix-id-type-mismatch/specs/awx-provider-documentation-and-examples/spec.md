## MODIFIED Requirements

### Requirement: Operational examples and import guidance
Each resource documentation page SHALL include at least one runnable example and SHALL define import usage with the accepted ID format. Documentation SHALL accurately describe typed numeric reference arguments and attributes, including object `id` as `Number` for collection-created objects and `String` for detail-path keyed objects. Examples SHALL avoid unnecessary type-conversion workarounds for reference wiring.

#### Scenario: Resource documentation review
- **WHEN** a resource documentation page is generated
- **THEN** it contains usage examples, argument and attribute references, and an import section with correct syntax

#### Scenario: Reference typing documentation consistency
- **WHEN** a resource exposes numeric reference fields
- **THEN** argument and attribute sections describe matching numeric usage and examples demonstrate direct assignment without `tonumber(...)`
