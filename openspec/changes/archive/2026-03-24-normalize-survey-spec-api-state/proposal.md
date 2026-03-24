# Proposal: Normalize survey spec API maps for stable Dynamic state

## Why

After native survey `spec` was exposed as a Terraform `Dynamic` object, `terraform apply` could fail with **Provider produced inconsistent result after apply** (for example `.spec: wrong final value type: attribute "spec": tuple required`). Causes include:

1. **Create** re-encoding the AWX GET response after POST: JSON shape (omitted keys, numeric typing) can differ from Terraform’s planned Dynamic/tuple type even when the survey is logically the same.
2. **Read/refresh**: sparse question objects in the GET payload can yield different nested object types per tuple index than configuration with explicit nulls for unused fields.

The OpenSpec change [`native-survey-spec-and-role-permissions`](../native-survey-spec-and-role-permissions/design.md) listed mitigation: *Normalize from API map to a stable object representation in state*.

## What Changes

1. **Normalize** survey-spec GET payloads in [`terraformObjectValueFromAPIValue`](../../../internal/provider/object_value_conversion.go) when handling `_survey_spec` / `spec` ([`survey_spec_normalization.go`](../../../internal/provider/survey_spec_normalization.go)): for the `spec` question array, build the **union of attribute names that appear on at least one question in the API response** (do not inject a fixed global list of keys onto every question—that can add attributes neither API nor config has, e.g. `new_question`, and break Dynamic typing vs plan). For each question, output a full map with that union, using `nil` for missing keys. Order known keys using `surveyQuestionFieldKeys`, then append any other keys sorted.
2. **Create** ([`relationship_resource.go`](../../../internal/provider/relationship_resource.go)): after a successful POST, set state `spec` from the **plan** (same as **Update**), not from a GET re-encode—this is what guarantees post-apply consistency with Terraform’s plan for create.
3. **Read** continues to use GET + normalization for drift detection.
4. Unit tests in [`survey_spec_normalization_test.go`](../../../internal/provider/survey_spec_normalization_test.go).

This does **not** replace [`fix-post-apply-null-json-array`](../fix-post-apply-null-json-array/proposal.md): that change targets optional **JSON-encoded array string** fields on manifest object resources, not survey `Dynamic` `spec`.

## Capabilities

### New Capabilities

- None.

### Modified Capabilities

- `awx-relationship-resources`: Survey-spec resources SHALL apply stable GET normalization on read; SHALL write planned `spec` to state on create after POST.

## Impact

- Runtime: [`internal/provider/survey_spec_normalization.go`](../../../internal/provider/survey_spec_normalization.go), [`internal/provider/object_value_conversion.go`](../../../internal/provider/object_value_conversion.go), [`internal/provider/relationship_resource.go`](../../../internal/provider/relationship_resource.go), tests.
- **Breaking:** None (behavioral fix).
