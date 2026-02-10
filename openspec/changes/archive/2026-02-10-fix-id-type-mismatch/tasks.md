## 1. Confirm scope and affected surfaces

- [x] 1.1 Inventory AWX object/data-source fields currently emitted as string but consumed as numeric references, and document target fields to normalize.
- [x] 1.2 Verify and document ID typing boundaries: collection-created object/data-source `id` is numeric; detail-path identifiers and relationship composite/survey import IDs remain string.

## 2. Implement type normalization in generation/runtime metadata

- [x] 2.1 Update openapi/manifest derivation logic to classify numeric foreign-key/reference fields for Terraform number typing.
- [x] 2.2 Ensure field override and curated manifest controls can preserve/override normalized typing where needed.
- [x] 2.3 Regenerate manifests (`make generate`) and validate sync (`make validate-manifest`) after normalization changes.

## 3. Apply typed references in provider schemas

- [x] 3.1 Update object resource schema construction to use normalized numeric typing for reference arguments and computed attributes.
- [x] 3.2 Update object data source schema construction to return the same numeric typing for reference attributes.
- [x] 3.3 Confirm relationship resources keep existing import/state identity contracts unchanged while consuming numeric parent/child references consistently.

## 4. Backward compatibility and state behavior

- [x] 4.1 Add/adjust conversion or read/write handling so legacy string-shaped state for affected reference fields and collection-created object IDs does not cause destructive behavior.
- [x] 4.2 Add regression checks covering mixed old/new state interactions and import flows.

## 5. Documentation and examples

- [x] 5.1 Regenerate docs (`make docs`) and validate docs (`make docs-validate`).
- [x] 5.2 Update generated examples/wording to show direct reference assignment without `tonumber(...)` where applicable.

## 6. Verification

- [x] 6.1 Add or update unit/provider tests validating cross-resource wiring of numeric references without explicit casting.
- [x] 6.2 Add or update tests validating `id` type behavior (`Number` for collection-created objects, `String` for detail-path keyed objects).
- [x] 6.3 Run full test suite (`make test`) and resolve failures.
- [x] 6.4 Record implementation evidence in the change folder and verify OpenSpec artifacts/tests are ready for apply.
