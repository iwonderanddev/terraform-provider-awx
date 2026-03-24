package provider

import (
	"reflect"
	"testing"
)

func TestNormalizeSurveySpecAPIMap_fillsMissingQuestionKeys(t *testing.T) {
	t.Parallel()
	in := map[string]any{
		"name":        "S",
		"description": "D",
		"spec": []any{
			map[string]any{
				"question_name": "q1",
				"type":          "text",
				"variable":      "v1",
				"required":      true,
			},
			map[string]any{
				"question_name": "q2",
				"type":          "multiplechoice",
				"variable":      "v2",
				"required":      false,
				"choices":       "a\nb",
				"default":       "a",
			},
		},
	}
	got := normalizeSurveySpecAPIMap(in)
	spec := got["spec"].([]any)
	q0 := spec[0].(map[string]any)
	if _, ok := q0["choices"]; !ok || q0["choices"] != nil {
		t.Fatalf("expected choices nil for first question: %#v", q0)
	}
	if _, hasMin := q0["min"]; hasMin && q0["min"] != nil {
		t.Fatalf("expected min absent or nil: %#v", q0)
	}
	q1 := spec[1].(map[string]any)
	if q1["choices"] != "a\nb" {
		t.Fatalf("choices: %#v", q1)
	}
}

func TestNormalizeSurveySpecAPIMap_emptySpec(t *testing.T) {
	t.Parallel()
	got := normalizeSurveySpecAPIMap(map[string]any{
		"name": "N",
	})
	spec := got["spec"].([]any)
	if len(spec) != 0 {
		t.Fatalf("want empty spec slice, got %#v", spec)
	}
}

func TestNormalizeSurveySpecAPIMap_preservesExtraRootKeys(t *testing.T) {
	t.Parallel()
	in := map[string]any{
		"name":        "x",
		"description": "",
		"spec":        []any{},
		"related":     map[string]any{"k": "v"},
	}
	got := normalizeSurveySpecAPIMap(in)
	if !reflect.DeepEqual(in["related"], got["related"]) {
		t.Fatalf("related: want %#v got %#v", in["related"], got["related"])
	}
}

func TestTerraformObjectValueFromAPIValue_surveySpec_sparseQuestions(t *testing.T) {
	t.Parallel()
	payload := map[string]any{
		"name":        "T",
		"description": "",
		"spec": []any{
			map[string]any{
				"question_name": "q1",
				"type":          "text",
				"variable":      "v1",
				"required":      true,
			},
		},
	}
	dyn, err := terraformObjectValueFromAPIValue("_survey_spec", "spec", payload)
	if err != nil {
		t.Fatalf("terraformObjectValueFromAPIValue: %v", err)
	}
	if dyn.IsNull() || dyn.IsUnknown() {
		t.Fatalf("expected non-null dynamic value")
	}
}

func TestNormalizeSurveySpecQuestions_doesNotInjectNewQuestion(t *testing.T) {
	t.Parallel()
	in := []any{
		map[string]any{
			"question_name": "a",
			"type":          "text",
			"variable":      "x",
			"required":      false,
		},
	}
	out := normalizeSurveySpecQuestions(in)
	m := out[0].(map[string]any)
	if _, ok := m["new_question"]; ok {
		t.Fatalf("new_question must not be injected when API omits it: %#v", m)
	}
}

func TestNormalizeSurveySpecQuestions_unknownKeys(t *testing.T) {
	t.Parallel()
	in := []any{
		map[string]any{
			"question_name": "a",
			"type":          "text",
			"variable":      "x",
			"future_field":  1,
		},
	}
	out := normalizeSurveySpecQuestions(in)
	m := out[0].(map[string]any)
	if m["future_field"] != 1 {
		t.Fatalf("missing extra key: %#v", m)
	}
	if _, hasChoices := m["choices"]; hasChoices {
		t.Fatalf("expected choices omitted when no question has it: %#v", m)
	}
}
