package openapi

import (
	"testing"

	"github.com/damien/terraform-awx-provider/internal/manifest"
)

func TestDeriveManagedObjects(t *testing.T) {
	t.Parallel()

	doc := &Document{
		Paths: map[string]PathItem{
			"/api/v2/projects/": {
				Get:  &Operation{},
				Post: &Operation{RequestBody: &RequestBody{Content: map[string]MediaType{"application/json": {Schema: &Schema{Ref: "#/components/schemas/ProjectRequest"}}}}},
			},
			"/api/v2/projects/{id}/": {
				Get:    &Operation{Responses: map[string]Response{"200": {Content: map[string]MediaType{"application/json": {Schema: &Schema{Ref: "#/components/schemas/Project"}}}}}},
				Patch:  &Operation{},
				Delete: &Operation{},
			},
		},
		Components: Components{Schemas: map[string]*Schema{
			"ProjectRequest": {
				Type:     "object",
				Required: []string{"name"},
				Properties: map[string]*Schema{
					"name":             {Type: "string"},
					"webhook_key":      {Type: "string"},
					"scm_branch":       {Type: "string"},
					"scm_update_cache": {Type: "boolean"},
				},
			},
		}},
	}

	objects := DeriveManagedObjects(doc, map[string]manifest.RuntimeExclusion{})
	if len(objects) != 1 {
		t.Fatalf("expected 1 object, got %d", len(objects))
	}

	object := objects[0]
	if object.Name != "projects" {
		t.Fatalf("unexpected object name: %s", object.Name)
	}
	if !object.ResourceEligible {
		t.Fatalf("expected projects resource to be resource-eligible")
	}
	if !object.DataSourceElig {
		t.Fatalf("expected projects resource to be data-source eligible")
	}

	foundSensitive := false
	for _, field := range object.Fields {
		if field.Name == "webhook_key" && field.Sensitive {
			foundSensitive = true
		}
	}
	if !foundSensitive {
		t.Fatalf("expected webhook_key field to be marked sensitive")
	}
}

func TestValidateCoverageRequiresRuntimeExclusion(t *testing.T) {
	t.Parallel()

	objects := []manifest.ManagedObject{
		{Name: "jobs", ResourceEligible: false, RuntimeExcluded: false},
	}
	exclusions := map[string]manifest.RuntimeExclusion{}

	if err := ValidateCoverage(objects, exclusions); err == nil {
		t.Fatalf("expected coverage validation to fail for missing runtime exclusion")
	}

	exclusions["jobs"] = manifest.RuntimeExclusion{Object: "jobs", Reason: "runtime"}
	if err := ValidateCoverage(objects, exclusions); err != nil {
		t.Fatalf("expected coverage validation to pass, got %v", err)
	}
}

func TestDeriveRelationships(t *testing.T) {
	t.Parallel()

	doc := &Document{
		Paths: map[string]PathItem{
			"/api/v2/teams/{id}/users/": {
				Get:  &Operation{},
				Post: &Operation{},
			},
		},
	}
	objects := []manifest.ManagedObject{
		{Name: "teams"},
		{Name: "users"},
	}
	priorities := map[string]int{"team_user_association": 10}

	relationships := DeriveRelationships(doc, objects, priorities)
	if len(relationships) != 1 {
		t.Fatalf("expected 1 relationship, got %d", len(relationships))
	}
	if relationships[0].Name != "team_user_association" {
		t.Fatalf("unexpected relationship name: %s", relationships[0].Name)
	}
	if relationships[0].Priority != 10 {
		t.Fatalf("expected priority 10, got %d", relationships[0].Priority)
	}
}

func TestFieldsFromSchemaMergesAllOfAndMarksSensitive(t *testing.T) {
	t.Parallel()

	doc := &Document{
		Components: Components{Schemas: map[string]*Schema{
			"BaseRequest": {
				Type:     "object",
				Required: []string{"name"},
				Properties: map[string]*Schema{
					"name":          {Type: "string"},
					"description":   {Type: "string", Default: ""},
					"admin_secret":  {Type: "string", Format: "password", Default: ""},
					"token_payload": {Type: "string", WriteOnly: true},
				},
			},
			"ProjectRequest": {
				AllOf: []*Schema{
					{Ref: "#/components/schemas/BaseRequest"},
					{
						Type:     "object",
						Required: []string{"enabled"},
						Properties: map[string]*Schema{
							"enabled": {Type: "boolean"},
							"tags":    {Items: &Schema{Type: "string"}, Default: []any{}},
							"options": {Properties: map[string]*Schema{"a": {Type: "string"}}},
						},
					},
				},
			},
		}},
	}

	fields := fieldsFromSchema(doc, "ProjectRequest")
	if len(fields) != 7 {
		t.Fatalf("expected 7 fields, got %d", len(fields))
	}

	nameField := findField(fields, "name")
	if !nameField.Required {
		t.Fatalf("expected name to remain required through allOf merge")
	}

	enabledField := findField(fields, "enabled")
	if !enabledField.Required {
		t.Fatalf("expected enabled to be required")
	}
	if enabledField.Computed {
		t.Fatalf("expected required field to never be computed")
	}

	secretField := findField(fields, "admin_secret")
	if !secretField.Sensitive {
		t.Fatalf("expected password-format field to be sensitive")
	}
	if !secretField.WriteOnly {
		t.Fatalf("expected sensitive field to be treated as write-only in provider schema")
	}

	descriptionField := findField(fields, "description")
	if descriptionField.Computed {
		t.Fatalf("expected empty-string default field to remain optional (not computed)")
	}

	tokenField := findField(fields, "token_payload")
	if !tokenField.Sensitive || !tokenField.WriteOnly {
		t.Fatalf("expected writeOnly field to be sensitive and write-only")
	}
	if tokenField.Computed {
		t.Fatalf("expected writeOnly field to not be computed")
	}

	tagsField := findField(fields, "tags")
	if tagsField.Type != manifest.FieldTypeArray {
		t.Fatalf("expected tags to resolve as array, got %s", tagsField.Type)
	}
	if !tagsField.Computed {
		t.Fatalf("expected optional field with default to be computed")
	}

	optionsField := findField(fields, "options")
	if optionsField.Type != manifest.FieldTypeObject {
		t.Fatalf("expected options to resolve as object, got %s", optionsField.Type)
	}
}

func findField(fields []manifest.FieldSpec, name string) manifest.FieldSpec {
	for _, field := range fields {
		if field.Name == name {
			return field
		}
	}
	return manifest.FieldSpec{}
}

func TestShouldInferComputedFromDefault(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		property *Schema
		required bool
		want     bool
	}{
		{
			name:     "nil property",
			property: nil,
			want:     false,
		},
		{
			name:     "required field",
			property: &Schema{Default: false},
			required: true,
			want:     false,
		},
		{
			name:     "write only field",
			property: &Schema{Default: false, WriteOnly: true},
			want:     false,
		},
		{
			name:     "no default",
			property: &Schema{},
			want:     false,
		},
		{
			name:     "empty string default",
			property: &Schema{Default: ""},
			want:     false,
		},
		{
			name:     "boolean default",
			property: &Schema{Default: false},
			want:     true,
		},
		{
			name:     "numeric default",
			property: &Schema{Default: float64(0)},
			want:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := shouldInferComputedFromDefault(tc.property, tc.required)
			if got != tc.want {
				t.Fatalf("unexpected computed inference: got=%v want=%v", got, tc.want)
			}
		})
	}
}
