## 1. Provider Foundation

- [ ] 1.1 Initialize provider module with `terraform-plugin-framework`, module layout, and development tooling configuration
- [ ] 1.2 Implement provider configuration schema for AWX base URL, HTTP Basic credentials, and TLS settings
- [ ] 1.3 Build shared API client transport layer with timeout, retry, pagination, and normalized error diagnostics
- [ ] 1.4 Add provider-level validation to fail fast on missing or invalid connectivity/auth configuration

## 2. OpenAPI Ingestion and Coverage Controls

- [ ] 2.1 Implement OpenAPI ingestion pipeline using `external/awx-openapi/schema.json` as input
- [ ] 2.2 Generate a managed-object manifest that maps AWX API v2 objects to resource/data source candidates
- [ ] 2.3 Implement explicit runtime-only exclusion manifest and validation checks for excluded object classes
- [ ] 2.4 Add CI/local validation command that fails when manifest coverage and generated mappings diverge

## 3. Core Resource and Data Source Scaffolding

- [ ] 3.1 Implement shared CRUD scaffolding for one-resource-per-AWX-object modeling with AWX-native naming
- [ ] 3.2 Implement shared data source scaffolding with deterministic lookup behavior and ambiguity/not-found diagnostics
- [ ] 3.3 Implement import/state normalization for object resources using numeric AWX IDs
- [ ] 3.4 Implement generated field mapping and override hooks for endpoint/schema mismatches

## 4. Relationship Resources and Import Semantics

- [ ] 4.1 Design relationship resource pattern for AWX associations with independent lifecycle semantics
- [ ] 4.2 Implement relationship resources for prioritized association types using explicit parent/child references
- [ ] 4.3 Implement composite import ID handling for relationship resources (`<parent_id>:<child_id>`)
- [ ] 4.4 Add relationship identity refresh logic that preserves stable composite IDs in state

## 5. Sensitive Field and State Safety

- [ ] 5.1 Mark secret and password-like schema attributes as sensitive across resources/data sources where applicable
- [ ] 5.2 Implement write-only secret handling to prevent cleartext secret round-tripping into Terraform state
- [ ] 5.3 Add guardrails for redaction-safe diagnostics and diff behavior on sensitive fields

## 6. Object Coverage Rollout

- [ ] 6.1 Generate and register resource implementations for all non-excluded managed object types in the manifest
- [ ] 6.2 Generate and register data source implementations for all supported lookup object types
- [ ] 6.3 Validate consistent CRUD and read-after-write convergence behavior across generated object resources
- [ ] 6.4 Validate remote-delete handling and state cleanup behavior across object categories

## 7. Documentation and Examples

- [ ] 7.1 Set up Terraform Registry-compatible docs structure for provider, resources, and data sources
- [ ] 7.2 Generate baseline docs per resource/data source including arguments and computed attributes
- [ ] 7.3 Add runnable usage examples and import examples (numeric IDs for objects, composite IDs for relationships)
- [ ] 7.4 Document compatibility scope (AWX 24.6.1 API v2 only) and runtime-only exclusion policy

## 8. Test Suite Implementation

- [ ] 8.1 Implement unit tests for API client behavior (auth, retries, pagination, error mapping)
- [ ] 8.2 Implement unit tests for schema mapping, scaffold behavior, and import/state normalization
- [ ] 8.3 Implement unit tests for sensitive-field handling and relationship identity semantics
- [ ] 8.4 Add local opt-in acceptance/e2e harness with required AWX environment variable gating and skip behavior
- [ ] 8.5 Implement acceptance scenarios for CRUD/import across representative object categories on AWX 24.6.1

## 9. Compatibility Validation and Release Readiness

- [ ] 9.1 Execute local acceptance/e2e validation against AWX 24.6.1 and capture known limitations
- [ ] 9.2 Verify coverage report aligns with manifest for all non-excluded managed objects
- [ ] 9.3 Verify generated docs/examples and import workflows are complete and accurate
- [ ] 9.4 Prepare initial release notes with support contract, exclusions, and test execution guidance
