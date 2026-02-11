## 1. Reference Field Inventory and Naming Rules

- [x] 1.1 Enumerate generated resource/data source fields that represent AWX object links and confirm canonical `_id` target names.
- [x] 1.2 Define deterministic rename rules in generation metadata, including handling for fields already ending in `_id`.
- [x] 1.3 Add collision detection rules for cases where both `<name>` and `<name>_id` could exist.

## 2. Generator and Manifest Changes

- [x] 2.1 Update openapi/manifest derivation so Terraform-facing object-link field names are generated with `_id` suffixes.
- [x] 2.2 Preserve AWX payload key mapping so renamed Terraform fields still map to the correct AWX API fields.
- [x] 2.3 Regenerate manifests and confirm reference-field naming changes in `internal/manifest/*.json` outputs.

## 3. Provider Schema and Runtime Mapping

- [x] 3.1 Update object resource schema construction to expose canonical `_id` link fields and remove unsuffixed aliases.
- [x] 3.2 Update object data source schema construction to expose canonical `_id` link attributes and remove unsuffixed aliases.
- [x] 3.3 Ensure create/update/read conversion logic uses suffixed Terraform names while preserving existing typed-ID behavior and import identity contracts.

## 4. Documentation and Examples

- [x] 4.1 Regenerate docs so arguments/attributes for link fields consistently use `_id` names.
- [x] 4.2 Update examples to demonstrate resource/data source wiring via `_id` fields (for example, `organization_id = awx_organization.example.id`).
- [x] 4.3 Verify docs include breaking-change migration guidance from unsuffixed to suffixed link field names.

## 5. Tests and Validation

- [x] 5.1 Add or update generator/provider unit tests validating `_id` naming for reference fields and absence of unsuffixed duplicates.
- [x] 5.2 Add or update tests ensuring renamed fields preserve numeric/string ID type behavior from existing contracts.
- [x] 5.3 Run validation chain: `make generate`, `make validate-manifest`, `make docs`, `make docs-validate`, and `make test`.
