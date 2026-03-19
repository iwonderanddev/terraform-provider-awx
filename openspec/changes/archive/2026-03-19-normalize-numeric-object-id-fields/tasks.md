# Tasks: Normalize Numeric `object_id` Fields

## 1. Generator And Manifest Typing

- [x] 1.1 Identify every generated resource and data source where `object_id`
  semantically represents an AWX numeric primary key.
- [x] 1.2 Update generator or curated manifest typing so affected `object_id`
  fields emit Terraform numeric metadata instead of string metadata.
- [x] 1.3 Preserve string typing for UUID-style alternatives such as
  `object_ansible_id` and for any `object_id` fields that are not numeric.

## 2. Provider Tests And Generated Outputs

- [x] 2.1 Add or update generator/manfiest tests to assert affected
  `object_id` fields emit numeric manifest types.
- [x] 2.2 Add or update provider schema tests to assert
  `awx_role_team_assignment.object_id` is numeric.
- [x] 2.3 Add or update provider schema tests for any additional affected
  surfaces discovered during implementation, including
  `awx_role_user_assignment` if it remains in scope.
- [x] 2.4 Run `make generate`.
- [x] 2.5 Run `make validate-manifest`.
- [x] 2.6 Run `make docs`.
- [x] 2.7 Run `make docs-validate`.
- [x] 2.8 Run `make test`.
- [x] 2.9 Run `make build`.

## 3. Documentation And Compatibility

- [x] 3.1 Verify generated docs describe affected `object_id` fields as numbers
  rather than strings.
- [x] 3.2 Document the breaking configuration/state migration impact in change
  artifacts and final implementation notes.
