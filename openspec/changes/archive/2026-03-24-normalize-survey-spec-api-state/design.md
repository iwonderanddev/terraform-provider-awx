# Design: Normalize survey spec API maps

## Context

Survey-spec resources are implemented in [`relationship_resource.go`](../../../internal/provider/relationship_resource.go). The `spec` attribute is a Plugin Framework `DynamicAttribute`. API values are converted with [`terraformObjectValueFromAPIValue`](../../../internal/provider/object_value_conversion.go), which maps `map[string]any` and `[]any` to Terraform `Object` and `Tuple` values via [`nativeToTerraformAttrValue`](../../../internal/provider/object_value_conversion.go).

Sparse JSON (omitted keys) can produce **different** nested object types per tuple index than HCL that sets `choices = null`, `min = null`, etc., on every question.

## Decisions

1. **Normalize in one place** — After `coerceObjectMap` and before `nativeToTerraformAttrValue`, when `objectName == "_survey_spec"` and `fieldName == "spec"`, replace the root map with `normalizeSurveySpecAPIMap` ([`survey_spec_normalization.go`](../../../internal/provider/survey_spec_normalization.go)).

2. **Question key union (API-driven)** — Collect keys that appear on **any** question object in the GET `spec` array. **Do not** pre-seed the union with every “known” AWX field: that would add attributes (e.g. `new_question`) that neither the API nor the configuration supplied, breaking Terraform’s planned Dynamic/tuple type. For ordering, emit keys that appear in both the union and [`surveyQuestionFieldKeys`](../../../internal/provider/survey_spec_normalization.go) in that list’s order, then any remaining keys sorted. For each question, emit a map containing **all** union keys, using `nil` where the question omits a key.

3. **Root fields** — Ensure `name` and `description` default to empty string when absent; `spec` missing or null becomes an empty array.

4. **POST payloads unchanged** — [`surveySpecConfig`](../../../internal/provider/relationship_resource.go) still sends the user’s decoded object.

5. **Create: state = plan** — On survey-spec **Create**, after a successful POST, state `spec` is set from the **plan** (same as **Update**), not from GET + conversion. That avoids Terraform Core rejecting the apply when the AWX response’s encoded shape differs from the planned value. The next **read**/refresh uses normalization against real AWX data.

## Non-Goals

- Replacing [`normalizeOptionalEmptyJSONEncodedArrayToNull`](../../../internal/provider/object_resource.go) (manifest JSON-array null vs `[]`).

## Risks

- **Future AWX fields:** Any key present on any question in the response is included in the union; unknown root-level keys on the survey object are preserved when copying the root map.
