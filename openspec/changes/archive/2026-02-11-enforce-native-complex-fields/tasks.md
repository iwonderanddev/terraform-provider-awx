## 1. Manifest and Type Source Updates

- [x] 1.1 Add curated `field_overrides` entries to type `job_templates.extra_vars` and `workflow_job_templates.extra_vars` as `object`.
- [x] 1.2 Regenerate managed manifests and confirm object typing for generated AWX `object` fields and both `extra_vars` fields.
- [x] 1.3 Validate manifest outputs are in sync with generator expectations.

## 2. Resource and Data Source Schema Typing

- [x] 2.1 Update object resource schema construction so generated `object` fields are exposed as Terraform object-capable dynamic attributes (not string attributes).
- [x] 2.2 Update object data source schema construction so generated `object` fields are exposed as Terraform object-capable dynamic attributes (not string attributes).
- [x] 2.3 Preserve existing required/optional/computed/sensitive/write-only behavior when switching object field attribute types.

## 3. Runtime Object Conversion and State Handling

- [x] 3.1 Replace JSON-string decoding paths for generated object fields in payload building with typed object-to-`map[string]any` conversion.
- [x] 3.2 Replace JSON-string encoding paths for generated object fields in state readback with typed API-to-Terraform object conversion.
- [x] 3.3 Enforce object-root validation diagnostics when configured/read values are not objects.
- [x] 3.4 Extend write-only snapshot preservation to retain typed object values for secret-bearing object fields.

## 4. Job Template extra_vars Transport Bridge

- [x] 4.1 Implement payload serialization for `awx_job_template.extra_vars` and `awx_workflow_job_template.extra_vars` from Terraform object values to AWX string transport format.
- [x] 4.2 Implement read normalization for both `extra_vars` fields from AWX responses: direct object passthrough, JSON string parse, YAML fallback parse.
- [x] 4.3 Return clear diagnostics when `extra_vars` normalization yields non-object root values.

## 5. Tests and Validation

- [x] 5.1 Update/add provider unit tests for object field schema typing and typed object conversion in resources and data sources.
- [x] 5.2 Update/add tests for typed write-only object preservation behavior.
- [x] 5.3 Update/add tests for `job_template` and `workflow_job_template` `extra_vars` normalization (JSON and YAML cases, invalid-root failures).
- [x] 5.4 Regenerate documentation and verify docs describe object field usage without JSON-string patterns for object fields.
- [x] 5.5 Run full validation chain: `make generate`, `make validate-manifest`, `make docs`, `make docs-validate`, and `make test`.
