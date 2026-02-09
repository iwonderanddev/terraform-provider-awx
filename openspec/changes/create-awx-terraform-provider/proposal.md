## Why

AWX configuration is often managed manually or with ad-hoc scripts, which creates drift and makes environments hard to reproduce. Building a Terraform provider directly on the AWX REST API and OpenAPI contract enables consistent, versioned, and testable AWX lifecycle management with Terraform-native resource modeling.

## What Changes

- Create a new Terraform provider for AWX using `/api/v2` REST endpoints and the published OpenAPI schema as the primary API contract.
- Implement provider configuration and authentication for API-based automation using HTTP Basic authentication at GA, with clear transport and TLS options.
- Add resources and data sources for all AWX API-managed objects, not only a subset of common objects.
- Model resources according to Terraform provider best practices: one resource per AWX API object, with dedicated relationship resources where needed instead of inline nested sub-object management.
- Keep AWX object names in Terraform resource and data source names for consistency with AWX API/domain terminology.
- Exclude runtime-only AWX objects from managed resource scope for the initial release.
- Treat sensitive and secret fields as write-only in Terraform schemas/state wherever applicable.
- Use numeric import IDs for normal object resources and composite import IDs for relationship resources.
- Add comprehensive tests, including unit tests and acceptance/e2e coverage for CRUD, import, update reconciliation, and relationship management behavior.
- Add Terraform Registry-grade documentation for provider, resources, and data sources, including practical examples and import guidance like established providers.
- Define GA compatibility as AWX `24.6.1` on API v2, with no backwards-compatibility requirement for older AWX versions.

## Capabilities

### New Capabilities

- `provider-auth-and-api-client`: Configure AWX host/auth/TLS and implement reusable API client behavior (pagination, filtering, retries, and error mapping).
- `awx-full-object-coverage`: Manage all AWX API-managed objects with Terraform resources and schema/state normalization.
- `awx-resource-lifecycle-and-import`: Provide consistent CRUD, read-after-write, drift detection, and import behavior across all resources.
- `awx-single-object-resource-model`: Enforce one Terraform resource per AWX API object and avoid inline nested sub-object lifecycle management.
- `awx-relationship-resources`: Manage object relationships using explicit resources and references instead of inline sub-object declarations.
- `awx-runtime-object-exclusions`: Exclude runtime-only AWX objects from managed Terraform resource scope.
- `awx-data-sources-and-lookups`: Provide data sources and lookup strategies required to reference existing AWX entities safely in Terraform.
- `awx-sensitive-field-handling`: Mark secrets as sensitive and write-only, avoiding secret value round-tripping in state.
- `awx-provider-documentation-and-examples`: Provide Terraform Registry-style docs for provider/resources/data sources with realistic examples and import usage.
- `awx-provider-test-suite`: Provide unit and acceptance/e2e tests, with acceptance tests run as opt-in local execution using user-supplied credentials/environment.
- `awx-ga-compatibility-target`: Validate and document GA compatibility for AWX `24.6.1` API v2 only (no backwards compatibility requirement).

### Modified Capabilities

- None.

## Impact

- Introduces a new provider implementation surface in this repository (provider config, API client, resources, data sources, docs, and tests).
- Depends on AWX API compatibility and behavior defined by AWX REST documentation and OpenAPI schema.
- Requires test workflow updates for unit tests and opt-in local acceptance/e2e tests against a reachable AWX instance.
- Aligns provider/resource modeling with HashiCorp provider design principles, reducing abstraction complexity and drift ambiguity.
- Expands scope to broad AWX object coverage and therefore increases implementation and test surface area.
- Narrows compatibility promises to AWX `24.6.1` API v2 for initial GA scope.
