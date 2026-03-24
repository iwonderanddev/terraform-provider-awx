# Tasks: normalize-survey-spec-api-state

## Implementation

- [x] Add [`internal/provider/survey_spec_normalization.go`](../../../internal/provider/survey_spec_normalization.go) (`normalizeSurveySpecAPIMap`, `normalizeSurveySpecQuestions`, slice coercion helper).
- [x] API-only key union (no pre-seed of all known fields); ordering via `surveyQuestionFieldKeys` + sorted extras.
- [x] Call normalization from [`terraformObjectValueFromAPIValue`](../../../internal/provider/object_value_conversion.go) for `_survey_spec` / `spec`.
- [x] Survey-spec **Create**: after POST, set state from planned `spec` only ([`relationship_resource.go`](../../../internal/provider/relationship_resource.go)); remove GET re-encode for post-apply state.
- [x] Unit tests: [`internal/provider/survey_spec_normalization_test.go`](../../../internal/provider/survey_spec_normalization_test.go).

## OpenSpec

- [x] Add `proposal.md`, `design.md`, `tasks.md`, `.openspec.yaml`.
- [x] Add spec delta under `specs/awx-relationship-resources/spec.md`.

## Verification

- [x] `go test ./...`
