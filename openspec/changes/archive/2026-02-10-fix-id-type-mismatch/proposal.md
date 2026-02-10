## Why

The provider currently exposes many reference ID attributes as Terraform strings while corresponding input arguments expect numbers. This forces consumers to cast with `tonumber(...)`, increases configuration friction, and can cause avoidable plan-time type errors in otherwise straightforward resource wiring.

## What Changes

- Align object reference argument types with emitted attribute types so resource-to-resource wiring works without explicit casting.
- Define a consistent ID typing contract across object resources, data sources, and relationship resources.
- Align object `id` typing with AWX/OpenAPI integer semantics for collection-created objects while retaining string identifiers for detail-path keyed objects and relationship identity formats.
- Update generated documentation to clearly state typed identifier behavior and migration expectations.
- Add/adjust acceptance and provider tests that cover cross-resource references without `tonumber(...)`.

## Capabilities

### New Capabilities
- `awx-typed-reference-ids`: Establishes consistent Terraform typing for AWX numeric reference identifiers used in arguments and computed attributes.

### Modified Capabilities
- `awx-single-object-resource-model`: Clarifies field typing rules for generated object schemas when AWX fields represent numeric identifiers.
- `awx-data-sources-and-lookups`: Ensures returned object attributes used as references follow the same identifier typing contract.
- `awx-provider-documentation-and-examples`: Updates docs/examples to reflect typed identifier usage and remove unnecessary casting.

## Impact

- Affected code: schema generation and runtime schema mapping in `internal/openapi/*`, `internal/manifest/*`, and `internal/provider/*`.
- Affected outputs: generated manifests and docs under `internal/manifest/*.json` and `docs/*`.
- User impact: simpler configurations with less explicit conversion; potential compatibility considerations for modules relying on legacy string-typed object `id` values for collection-created resources.
- Validation impact: provider/unit/acceptance tests need coverage for typed references and state upgrade safety.
