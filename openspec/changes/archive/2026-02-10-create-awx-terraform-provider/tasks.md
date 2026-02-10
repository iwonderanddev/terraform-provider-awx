## 1. Provider Foundation

- [x] 1.1 Initialize provider module with `terraform-plugin-framework`, module layout, and development tooling configuration
- [x] 1.2 Implement provider configuration schema for AWX base URL, HTTP Basic credentials, and TLS settings
- [x] 1.3 Build shared API client transport layer with timeout, retry, pagination, and normalized error diagnostics
- [x] 1.4 Add provider-level validation to fail fast on missing or invalid connectivity/auth configuration

## 2. OpenAPI Ingestion and Coverage Controls

- [x] 2.1 Implement OpenAPI ingestion pipeline using `external/awx-openapi/schema.json` as input
- [x] 2.2 Generate a managed-object manifest that maps AWX API v2 objects to resource/data source candidates
- [x] 2.3 Implement explicit runtime-only exclusion manifest and validation checks for excluded object classes
- [x] 2.4 Add CI/local validation command that fails when manifest coverage and generated mappings diverge

## 3. Core Resource and Data Source Scaffolding

- [x] 3.1 Implement shared CRUD scaffolding for one-resource-per-AWX-object modeling with AWX-native naming
- [x] 3.2 Implement shared data source scaffolding with deterministic lookup behavior and ambiguity/not-found diagnostics
- [x] 3.3 Implement import/state normalization for object resources using endpoint-aligned object identifiers (numeric or detail-key)
- [x] 3.4 Implement generated field mapping and override hooks for endpoint/schema mismatches

## 4. Relationship Resources and Import Semantics

- [x] 4.1 Design relationship resource pattern for AWX associations with independent lifecycle semantics
- [x] 4.2 Implement relationship resources for prioritized association types using explicit parent/child references
- [x] 4.3 Implement relationship import ID handling for both association resources (`<parent_id>:<child_id>`) and singleton parent-scoped resources (`<parent_id>`)
- [x] 4.4 Add relationship identity refresh logic that preserves stable endpoint-aligned IDs in state

## 5. Sensitive Field and State Safety

- [x] 5.1 Mark secret and password-like schema attributes as sensitive across resources/data sources where applicable
- [x] 5.2 Implement write-only secret handling to prevent cleartext secret round-tripping into Terraform state
- [x] 5.3 Add guardrails for redaction-safe diagnostics and diff behavior on sensitive fields

## 6. Object Coverage Rollout

- [x] 6.1 Generate and register resource implementations for all non-excluded managed object types in the manifest
- [x] 6.2 Generate and register data source implementations for all supported lookup object types
- [x] 6.3 Validate consistent CRUD and read-after-write convergence behavior across generated object resources
- [x] 6.4 Validate remote-delete handling and state cleanup behavior across object categories

## 7. Documentation and Examples

- [x] 7.1 Set up Terraform Registry-compatible docs structure for provider, resources, and data sources
- [x] 7.2 Generate baseline docs per resource/data source including arguments and computed attributes
- [x] 7.3 Add runnable usage examples and import examples for all supported identifier formats (numeric/detail-key for objects, composite/parent-key for relationships)
- [x] 7.4 Document compatibility scope (AWX 24.6.1 API v2 only) and runtime-only exclusion policy

## 8. Test Suite Implementation

- [x] 8.1 Implement unit tests for API client behavior (auth, retries, pagination, error mapping)
- [x] 8.2 Implement unit tests for schema mapping, scaffold behavior, and import/state normalization
- [x] 8.3 Implement unit tests for sensitive-field handling and relationship identity semantics
- [x] 8.4 Add local opt-in acceptance/e2e harness with required AWX environment variable gating and skip behavior
- [x] 8.5 Implement acceptance scenarios for CRUD/import across representative object categories on AWX 24.6.1
- [x] 8.6 Implement Terraform-driven acceptance tests (`terraform-plugin-testing`) for representative resource/data source/relationship lifecycle and import flows

## 9. Compatibility Validation and Release Readiness

- [x] 9.1 Execute local acceptance/e2e validation against AWX 24.6.1 and capture known limitations
- [x] 9.2 Verify coverage report aligns with manifest for all non-excluded managed objects
- [x] 9.3 Verify generated docs/examples and import workflows are complete and accurate
- [x] 9.4 Prepare initial release notes with support contract, exclusions, and test execution guidance
- [x] 9.5 Execute Terraform-driven acceptance validation against AWX 24.6.1 and capture discovered provider behavior fixes

## 10. Default-Value Drift Hardening

- [x] 10.1 Extend OpenAPI schema parsing to ingest request-property `default` values
- [x] 10.2 Mark optional non-write-only fields with OpenAPI defaults as `Optional + Computed` in generated provider schemas
- [x] 10.3 Add manifest regression checks for server-defaulted fields (`organizations.max_hosts`, `inventories.prevent_instance_group_fallback`)
- [x] 10.4 Add Terraform-driven acceptance regression scenarios for omitted defaulted fields using create + plan-only + import validation
- [x] 10.5 Improve acceptance test logs with explicit per-step progress messages for API and Terraform-driven flows
