# Tasks: rename-project-credential-id-to-scm-credential-id

## 1. Define canonical field-name override support

- [x] 1.1 Inspect the current object-field naming pipeline and identify the
  curated metadata shape needed to override a generated Terraform field name
  for a specific AWX object field.
- [x] 1.2 Implement generator and manifest support for object-field canonical
  Terraform name overrides without changing unrelated reference-field naming
  behavior.
- [x] 1.3 Add or update unit tests covering canonical-name override handling in
  openapi/manifest derivation.

## 2. Apply the project SCM credential rename

- [x] 2.1 Add a curated override for `projects.credential` so the generated
  Terraform field name becomes `scm_credential_id`.
- [x] 2.2 Ensure generated `awx_project` resource and data source metadata no
  longer expose `credential_id` for the SCM credential field.
- [x] 2.3 Add or update schema-focused tests asserting `awx_project` exposes
  `scm_credential_id` and rejects or omits the legacy field name.

## 3. Regenerate manifests and documentation

- [x] 3.1 Run `make generate` and `make validate-manifest` to refresh manifest
  outputs for the renamed field.
- [x] 3.2 Update curated docs-enrichment metadata and regenerate docs with
  `make docs`.
- [x] 3.3 Run `make docs-validate` and confirm generated `awx_project`
  resource/data source docs and examples use `scm_credential_id`
  consistently.

## 4. Verify provider behavior and readiness

- [x] 4.1 Add or update tests covering project resource/data source read/write
  behavior for the renamed SCM credential field.
- [x] 4.2 Run `make test` and resolve any regressions caused by the breaking
  rename.
- [x] 4.3 Run `make build` and confirm the change is ready for implementation
  review with the documented breaking contract.
