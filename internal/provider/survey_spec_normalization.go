package provider

import (
	"reflect"
	"sort"
)

// surveyQuestionFieldKeys is the stable order for known AWX survey question fields.
// Extra keys seen in API responses are appended after these (sorted).
var surveyQuestionFieldKeys = []string{
	"question_name",
	"question_description",
	"required",
	"type",
	"variable",
	"default",
	"choices",
	"min",
	"max",
	"new_question",
}

// normalizeSurveySpecAPIMap returns a copy of the survey root map with stable
// keys on each question object so Terraform Dynamic/tuple typing matches HCL
// that uses explicit nulls for unused fields (OpenSpec native-survey-spec mitigation).
func normalizeSurveySpecAPIMap(root map[string]any) map[string]any {
	out := make(map[string]any, len(root)+3)
	for k, v := range root {
		if k != "name" && k != "description" && k != "spec" {
			out[k] = v
		}
	}
	if v, ok := root["name"]; ok && v != nil {
		out["name"] = v
	} else {
		out["name"] = ""
	}
	if v, ok := root["description"]; ok && v != nil {
		out["description"] = v
	} else {
		out["description"] = ""
	}
	rawSpec, ok := root["spec"]
	if !ok || rawSpec == nil {
		out["spec"] = []any{}
		return out
	}
	arr, ok := rawSpec.([]any)
	if !ok {
		arr = coerceSliceToAnySlice(rawSpec)
		if arr == nil {
			out["spec"] = rawSpec
			return out
		}
	}
	out["spec"] = normalizeSurveySpecQuestions(arr)
	return out
}

func coerceSliceToAnySlice(rawSpec any) []any {
	rv := reflect.ValueOf(rawSpec)
	if rv.Kind() != reflect.Slice {
		return nil
	}
	out := make([]any, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		out[i] = rv.Index(i).Interface()
	}
	return out
}

func normalizeSurveySpecQuestions(questions []any) []any {
	if len(questions) == 0 {
		return []any{}
	}
	// Union keys from the API response only. Do not inject every "known" field on
	// every question — that adds attributes absent from both API and config (e.g.
	// new_question) and breaks Terraform Dynamic/tuple typing vs plan.
	keySet := make(map[string]struct{})
	for _, q := range questions {
		m, ok := q.(map[string]any)
		if !ok {
			continue
		}
		for k := range m {
			keySet[k] = struct{}{}
		}
	}
	seen := make(map[string]struct{})
	union := make([]string, 0, len(keySet))
	for _, k := range surveyQuestionFieldKeys {
		if _, has := keySet[k]; has {
			union = append(union, k)
			seen[k] = struct{}{}
		}
	}
	var extras []string
	for k := range keySet {
		if _, ok := seen[k]; !ok {
			extras = append(extras, k)
		}
	}
	sort.Strings(extras)
	union = append(union, extras...)

	out := make([]any, len(questions))
	for i, q := range questions {
		m, ok := q.(map[string]any)
		if !ok {
			out[i] = q
			continue
		}
		full := make(map[string]any, len(union))
		for _, k := range union {
			if v, ok := m[k]; ok {
				full[k] = v
			} else {
				full[k] = nil
			}
		}
		out[i] = full
	}
	return out
}
