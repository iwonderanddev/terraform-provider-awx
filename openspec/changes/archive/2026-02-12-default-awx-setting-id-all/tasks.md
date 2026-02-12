# Tasks: Default awx_setting ID to all in Docs/UX

## 1. Generator and Content Updates

- [x] 1.1 Update docs generation inputs/logic so `awx_setting` example usage
      defaults to `id = "all"` for the resource and data source pages.
- [x] 1.2 Update `awx_setting` import guidance to use `all` as the primary
      identifier example.
- [x] 1.3 Add documentation text that category-scoped IDs remain supported as
      optional advanced usage.
- [x] 1.4 Add explicit warning guidance for overlapping ownership/conflicts when
      mixing `id = "all"` and category-scoped `awx_setting` resources.

## 2. Validation and Test Coverage

- [x] 2.1 Extend docs validation rules/tests to enforce `awx_setting` default
      examples/import guidance use `all`.
- [x] 2.2 Add/adjust tests to ensure category-scoped ID guidance remains
      present and marked as optional advanced scoping.
- [x] 2.3 Add/adjust tests to verify overlap/conflict warning text is included
      for `awx_setting` documentation.
- [x] 2.4 Add an integration-style docs generation test that asserts rendered
      `awx_setting` resource/data-source section ordering and canonical
      `id = "all"` guidance in generated output.

## 3. Regeneration and Verification

- [x] 3.1 Run `make docs` to regenerate documentation artifacts.
- [x] 3.2 Run `make docs-validate` and resolve failures.
- [x] 3.3 Run `make test` and resolve regressions related to docs generation.
- [x] 3.4 Run markdown lint on changed OpenSpec artifacts and fix remaining
      violations.

## 4. Final Consistency Checks

- [x] 4.1 Verify generated docs consistently present `id = "all"` as default
      for `awx_setting` while retaining scoped examples or guidance.
- [x] 4.2 Confirm no runtime/provider API contract changes were introduced and
      import ID contracts remain unchanged.
