## 1. Relationship Metadata and Generation

- [x] 1.1 Add canonical relationship argument-name derivation from `parentObject`/`childObject` (including survey-spec parent naming) in generation/manifest logic.
- [x] 1.2 Extend relationship manifest model/output to carry canonical parent and child Terraform argument names.
- [x] 1.3 Regenerate relationship manifests and confirm canonical `*_id` argument metadata is present for all relationship resources.

## 2. Relationship Resource Schema and Runtime

- [x] 2.1 Update relationship resource schema construction to expose canonical object-specific `*_id` arguments for standard and survey-spec relationships.
- [x] 2.2 Remove generic `parent_id`/`child_id` schema arguments from relationship resources and enforce canonical object-specific argument names only.
- [x] 2.3 Update create/read/import/state mapping helpers to read canonical fields, preserve ID contracts (`<parent_id>:<child_id>` and `<parent_id>`), and keep AWX payload behavior unchanged.

## 3. Documentation and Examples

- [x] 3.1 Regenerate relationship resource docs so argument references and examples use canonical object-specific `*_id` names.
- [x] 3.2 Add breaking-change migration guidance in docs for transitioning from `parent_id`/`child_id` to canonical names.
- [x] 3.3 Update example configurations that declare relationship resources to remove generic directional arguments.

## 4. Tests and Validation

- [x] 4.1 Add/adjust generator and provider unit tests for canonical relationship argument naming and legacy `parent_id`/`child_id` rejection behavior.
- [x] 4.2 Update Terraform acceptance tests for relationship resources to use canonical `*_id` arguments and verify unchanged import/state identity behavior.
- [x] 4.3 Run `make generate`, `make validate-manifest`, `make docs`, `make docs-validate`, and `make test`.
