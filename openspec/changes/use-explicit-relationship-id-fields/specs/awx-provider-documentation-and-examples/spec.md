## MODIFIED Requirements

### Requirement: Operational examples and import guidance
Each resource documentation page SHALL include at least one runnable example and SHALL define import usage with the accepted ID format. Documentation SHALL accurately describe typed numeric reference arguments and attributes, including object `id` as `Number` for collection-created objects and `String` for detail-path keyed objects. Documentation for relationship resources SHALL use canonical explicit object-specific `_id` argument names (for example, `job_template_id`, `credential_id`) and SHALL include breaking-change migration guidance for prior `parent_id`/`child_id` configurations. Examples SHALL avoid unnecessary type-conversion workarounds for reference wiring.

#### Scenario: Resource documentation review
- **WHEN** a resource documentation page is generated
- **THEN** it contains usage examples, argument and attribute references, and an import section with correct syntax

#### Scenario: Reference typing documentation consistency
- **WHEN** a resource exposes numeric reference fields
- **THEN** argument and attribute sections describe matching numeric usage and examples demonstrate direct assignment without `tonumber(...)`

#### Scenario: Relationship argument naming documentation consistency
- **WHEN** documentation is generated for relationship resources
- **THEN** argument and attribute sections use canonical object-specific `_id` names and examples avoid generic `parent_id`/`child_id` usage

#### Scenario: Relationship breaking-change documentation consistency
- **WHEN** relationship argument naming changes are released as a hard break
- **THEN** the documentation includes explicit migration guidance from legacy `parent_id`/`child_id` names to canonical object-specific `*_id` names
