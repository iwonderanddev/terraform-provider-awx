package openapi

import (
	"testing"

	"github.com/damien/terraform-provider-awx-iwd/internal/manifest"
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

	objects := DeriveManagedObjects(doc, map[string]manifest.RuntimeExclusion{}, map[string]string{})
	if len(objects) != 1 {
		t.Fatalf("expected 1 object, got %d", len(objects))
	}

	object := objects[0]
	if object.Name != "projects" {
		t.Fatalf("unexpected object name: %s", object.Name)
	}
	if !object.CollectionCreate {
		t.Fatalf("expected projects object to support collection create")
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

func TestDeriveManagedObjectsSupportsNonIDDetailPath(t *testing.T) {
	t.Parallel()

	doc := &Document{
		Paths: map[string]PathItem{
			"/api/v2/settings/": {
				Get: &Operation{},
			},
			"/api/v2/settings/{category_slug}/": {
				Get:    &Operation{Responses: map[string]Response{"200": {Content: map[string]MediaType{"application/json": {Schema: &Schema{Ref: "#/components/schemas/SettingSingleton"}}}}}},
				Patch:  &Operation{RequestBody: &RequestBody{Content: map[string]MediaType{"application/json": {Schema: &Schema{Ref: "#/components/schemas/SettingSingletonRequest"}}}}},
				Delete: &Operation{},
			},
		},
		Components: Components{Schemas: map[string]*Schema{
			"SettingSingletonRequest": {
				Type:     "object",
				Required: []string{"tower_url_base"},
				Properties: map[string]*Schema{
					"tower_url_base": {Type: "string"},
				},
			},
			"SettingSingleton": {
				Type: "object",
				Properties: map[string]*Schema{
					"tower_url_base": {Type: "string"},
				},
			},
		}},
	}

	objects := DeriveManagedObjects(doc, map[string]manifest.RuntimeExclusion{}, map[string]string{})
	if len(objects) != 1 {
		t.Fatalf("expected 1 object, got %d", len(objects))
	}

	object := objects[0]
	if object.Name != "settings" {
		t.Fatalf("unexpected object name: %s", object.Name)
	}
	if object.CollectionCreate {
		t.Fatalf("expected settings to use detail-path lifecycle (no collection create)")
	}
	if !object.ResourceEligible {
		t.Fatalf("expected settings to be resource-eligible")
	}
	if object.DetailPath != "/api/v2/settings/{category_slug}/" {
		t.Fatalf("unexpected settings detail path: %s", object.DetailPath)
	}
}

func TestDeriveManagedObjectsSupportsCreateDeleteOnlyResources(t *testing.T) {
	t.Parallel()

	doc := &Document{
		Paths: map[string]PathItem{
			"/api/v2/role_user_assignments/": {
				Get:  &Operation{},
				Post: &Operation{RequestBody: &RequestBody{Content: map[string]MediaType{"application/json": {Schema: &Schema{Ref: "#/components/schemas/RoleUserAssignmentRequest"}}}}},
			},
			"/api/v2/role_user_assignments/{id}/": {
				Get:    &Operation{Responses: map[string]Response{"200": {Content: map[string]MediaType{"application/json": {Schema: &Schema{Ref: "#/components/schemas/RoleUserAssignment"}}}}}},
				Delete: &Operation{},
			},
		},
		Components: Components{Schemas: map[string]*Schema{
			"RoleUserAssignmentRequest": {
				Type:     "object",
				Required: []string{"role_definition", "user"},
				Properties: map[string]*Schema{
					"role_definition": {Type: "integer"},
					"user":            {Type: "integer"},
				},
			},
		}},
	}

	objects := DeriveManagedObjects(doc, map[string]manifest.RuntimeExclusion{}, map[string]string{})
	if len(objects) != 1 {
		t.Fatalf("expected 1 object, got %d", len(objects))
	}

	object := objects[0]
	if !object.CollectionCreate {
		t.Fatalf("expected collection create support")
	}
	if object.UpdateSupported {
		t.Fatalf("expected update support to be disabled for create/delete-only object")
	}
	if !object.ResourceEligible {
		t.Fatalf("expected create/delete-only object to remain resource-eligible")
	}
}

func TestDeriveManagedObjectsMarksReferenceFields(t *testing.T) {
	t.Parallel()

	doc := &Document{
		Paths: map[string]PathItem{
			"/api/v2/organizations/": {
				Get:  &Operation{},
				Post: &Operation{RequestBody: &RequestBody{Content: map[string]MediaType{"application/json": {Schema: &Schema{Ref: "#/components/schemas/OrganizationRequest"}}}}},
			},
			"/api/v2/organizations/{id}/": {
				Get:    &Operation{Responses: map[string]Response{"200": {Content: map[string]MediaType{"application/json": {Schema: &Schema{Ref: "#/components/schemas/Organization"}}}}}},
				Patch:  &Operation{},
				Delete: &Operation{},
			},
			"/api/v2/teams/": {
				Get:  &Operation{},
				Post: &Operation{RequestBody: &RequestBody{Content: map[string]MediaType{"application/json": {Schema: &Schema{Ref: "#/components/schemas/TeamRequest"}}}}},
			},
			"/api/v2/teams/{id}/": {
				Get:    &Operation{Responses: map[string]Response{"200": {Content: map[string]MediaType{"application/json": {Schema: &Schema{Ref: "#/components/schemas/Team"}}}}}},
				Patch:  &Operation{},
				Delete: &Operation{},
			},
		},
		Components: Components{Schemas: map[string]*Schema{
			"OrganizationRequest": {
				Type: "object",
				Properties: map[string]*Schema{
					"name": {Type: "string"},
				},
				Required: []string{"name"},
			},
			"Organization": {
				Type: "object",
				Properties: map[string]*Schema{
					"id":   {Type: "integer"},
					"name": {Type: "string"},
				},
			},
			"TeamRequest": {
				Type: "object",
				Properties: map[string]*Schema{
					"name":         {Type: "string"},
					"organization": {Type: "integer"},
					"max_hosts":    {Type: "integer"},
				},
				Required: []string{"name", "organization"},
			},
			"Team": {
				Type: "object",
				Properties: map[string]*Schema{
					"id":           {Type: "integer"},
					"name":         {Type: "string"},
					"organization": {Type: "integer"},
					"max_hosts":    {Type: "integer"},
				},
			},
		}},
	}

	objects := DeriveManagedObjects(doc, map[string]manifest.RuntimeExclusion{}, map[string]string{})
	if len(objects) != 2 {
		t.Fatalf("expected 2 objects, got %d", len(objects))
	}

	var team manifest.ManagedObject
	for _, object := range objects {
		if object.Name == "teams" {
			team = object
			break
		}
	}
	if team.Name == "" {
		t.Fatalf("expected teams object to be derived")
	}

	orgField := findField(team.Fields, "organization")
	if !orgField.Reference {
		t.Fatalf("expected teams.organization to be marked as reference")
	}

	maxHostsField := findField(team.Fields, "max_hosts")
	if maxHostsField.Reference {
		t.Fatalf("expected teams.max_hosts to remain non-reference")
	}
}

func TestDeriveManagedObjectsDoesNotInferReferenceFromIDSuffixAlone(t *testing.T) {
	t.Parallel()

	doc := &Document{
		Paths: map[string]PathItem{
			"/api/v2/teams/": {
				Get:  &Operation{},
				Post: &Operation{RequestBody: &RequestBody{Content: map[string]MediaType{"application/json": {Schema: &Schema{Ref: "#/components/schemas/TeamRequest"}}}}},
			},
			"/api/v2/teams/{id}/": {
				Get:    &Operation{Responses: map[string]Response{"200": {Content: map[string]MediaType{"application/json": {Schema: &Schema{Ref: "#/components/schemas/Team"}}}}}},
				Patch:  &Operation{},
				Delete: &Operation{},
			},
		},
		Components: Components{Schemas: map[string]*Schema{
			"TeamRequest": {
				Type: "object",
				Properties: map[string]*Schema{
					"name":      {Type: "string"},
					"legacy_id": {Type: "integer"},
				},
				Required: []string{"name"},
			},
			"Team": {
				Type: "object",
				Properties: map[string]*Schema{
					"id":        {Type: "integer"},
					"name":      {Type: "string"},
					"legacy_id": {Type: "integer"},
				},
			},
		}},
	}

	objects := DeriveManagedObjects(doc, map[string]manifest.RuntimeExclusion{}, map[string]string{})
	if len(objects) != 1 {
		t.Fatalf("expected 1 object, got %d", len(objects))
	}

	legacyField := findField(objects[0].Fields, "legacy_id")
	if legacyField.Reference {
		t.Fatalf("expected teams.legacy_id to remain non-reference without metadata classification")
	}
}

func TestIsReferenceFieldSupportsCommonAWXReferenceNamePatterns(t *testing.T) {
	t.Parallel()

	candidates := map[string]struct{}{
		"project":               {},
		"credential":            {},
		"execution_environment": {},
		"user":                  {},
	}

	cases := []struct {
		name  string
		field manifest.FieldSpec
		want  bool
	}{
		{
			name: "exact_match",
			field: manifest.FieldSpec{
				Name: "project",
				Type: manifest.FieldTypeInt,
			},
			want: true,
		},
		{
			name: "prefixed_reference_name",
			field: manifest.FieldSpec{
				Name: "source_project",
				Type: manifest.FieldTypeInt,
			},
			want: true,
		},
		{
			name: "suffixed_reference_name",
			field: manifest.FieldSpec{
				Name: "webhook_credential",
				Type: manifest.FieldTypeInt,
			},
			want: true,
		},
		{
			name: "description_backed_alias_name",
			field: manifest.FieldSpec{
				Name:        "default_environment",
				Type:        manifest.FieldTypeInt,
				Description: "The default execution environment for jobs run by this organization.",
			},
			want: true,
		},
		{
			name: "created_by_user_link",
			field: manifest.FieldSpec{
				Name:        "created_by",
				Type:        manifest.FieldTypeInt,
				Description: "The user who created this resource.",
			},
			want: true,
		},
		{
			name: "counter_field_is_not_reference",
			field: manifest.FieldSpec{
				Name: "max_hosts",
				Type: manifest.FieldTypeInt,
			},
			want: false,
		},
		{
			name: "timeout_field_is_not_reference",
			field: manifest.FieldSpec{
				Name:        "scm_update_cache_timeout",
				Type:        manifest.FieldTypeInt,
				Description: "The number of seconds after the last project update ran that a new project update will be launched.",
			},
			want: false,
		},
		{
			name: "per_user_counter_is_not_reference",
			field: manifest.FieldSpec{
				Name:        "sessions_per_user",
				Type:        manifest.FieldTypeInt,
				Description: "Maximum number of simultaneous logged in sessions a user may have.",
			},
			want: false,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := isReferenceField(tc.field, candidates)
			if got != tc.want {
				t.Fatalf("isReferenceField(%q) mismatch: got=%t want=%t", tc.field.Name, got, tc.want)
			}
		})
	}
}

func TestDeriveManagedObjectsExcludesDeprecatedObjects(t *testing.T) {
	t.Parallel()

	doc := &Document{
		Paths: map[string]PathItem{
			"/api/v2/roles/": {
				Get: &Operation{},
			},
			"/api/v2/roles/{id}/": {
				Get: &Operation{},
			},
		},
	}

	objects := DeriveManagedObjects(doc, map[string]manifest.RuntimeExclusion{}, map[string]string{
		"roles": "Deprecated in favor of role_definitions.",
	})
	if len(objects) != 1 {
		t.Fatalf("expected 1 object, got %d", len(objects))
	}
	if objects[0].ResourceEligible {
		t.Fatalf("expected deprecated object to be excluded from resources")
	}
	if objects[0].DataSourceElig {
		t.Fatalf("expected deprecated object to be excluded from data sources")
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

	relationships := DeriveRelationships(doc, objects, priorities, map[string]string{})
	if len(relationships) != 1 {
		t.Fatalf("expected 1 relationship, got %d", len(relationships))
	}
	if relationships[0].Name != "team_user_association" {
		t.Fatalf("unexpected relationship name: %s", relationships[0].Name)
	}
	if relationships[0].Priority != 10 {
		t.Fatalf("expected priority 10, got %d", relationships[0].Priority)
	}
	if relationships[0].ParentIDAttribute != "team_id" {
		t.Fatalf("unexpected parent attribute: got=%q want=%q", relationships[0].ParentIDAttribute, "team_id")
	}
	if relationships[0].ChildIDAttribute != "user_id" {
		t.Fatalf("unexpected child attribute: got=%q want=%q", relationships[0].ChildIDAttribute, "user_id")
	}
}

func TestDeriveRelationshipsDetectsNotificationTemplateVariantAndSurveySpec(t *testing.T) {
	t.Parallel()

	doc := &Document{
		Paths: map[string]PathItem{
			"/api/v2/job_templates/{id}/notification_templates_error/": {
				Get:  &Operation{},
				Post: &Operation{},
			},
			"/api/v2/job_templates/{id}/survey_spec/": {
				Get:    &Operation{},
				Post:   &Operation{},
				Delete: &Operation{},
			},
		},
	}
	objects := []manifest.ManagedObject{
		{Name: "job_templates"},
		{Name: "notification_templates"},
	}

	relationships := DeriveRelationships(doc, objects, map[string]int{}, map[string]string{})
	if len(relationships) != 2 {
		t.Fatalf("expected 2 relationships, got %d", len(relationships))
	}

	seen := map[string]manifest.Relationship{}
	for _, rel := range relationships {
		seen[rel.Name] = rel
	}

	errorRel, ok := seen["job_template_notification_template_error"]
	if !ok {
		t.Fatalf("missing notification template error relationship")
	}
	if errorRel.ChildObject != "notification_templates" {
		t.Fatalf("unexpected child object mapping: got=%q want=%q", errorRel.ChildObject, "notification_templates")
	}
	if errorRel.ParentIDAttribute != "job_template_id" {
		t.Fatalf("unexpected notification relationship parent attribute: got=%q want=%q", errorRel.ParentIDAttribute, "job_template_id")
	}
	if errorRel.ChildIDAttribute != "notification_template_id" {
		t.Fatalf("unexpected notification relationship child attribute: got=%q want=%q", errorRel.ChildIDAttribute, "notification_template_id")
	}

	surveyRel, ok := seen["job_template_survey_spec"]
	if !ok {
		t.Fatalf("missing survey spec relationship")
	}
	if surveyRel.ParentIDAttribute != "job_template_id" {
		t.Fatalf("unexpected survey-spec parent attribute: got=%q want=%q", surveyRel.ParentIDAttribute, "job_template_id")
	}
	if surveyRel.ChildIDAttribute != "" {
		t.Fatalf("expected survey-spec child attribute to be empty, got=%q", surveyRel.ChildIDAttribute)
	}
}

func TestDeriveRelationshipsExcludesDeprecatedPaths(t *testing.T) {
	t.Parallel()

	doc := &Document{
		Paths: map[string]PathItem{
			"/api/v2/users/{id}/roles/": {
				Get:  &Operation{},
				Post: &Operation{},
			},
		},
	}
	objects := []manifest.ManagedObject{
		{Name: "users"},
		{Name: "roles"},
	}

	relationships := DeriveRelationships(doc, objects, map[string]int{}, map[string]string{
		"/api/v2/users/{id}/roles/": "Deprecated in favor of role_user_assignments.",
	})
	if len(relationships) != 0 {
		t.Fatalf("expected deprecated relationship path to be excluded, got %d entries", len(relationships))
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
							"enabled":         {Type: "boolean"},
							"tags":            {Items: &Schema{Type: "string"}, Default: []any{}},
							"options":         {Properties: map[string]*Schema{"a": {Type: "string"}}},
							"computed_server": {Type: "string", ReadOnly: true},
						},
					},
				},
			},
		}},
	}

	fields := fieldsFromSchema(doc, "ProjectRequest")
	if len(fields) != 8 {
		t.Fatalf("expected 8 fields, got %d", len(fields))
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

	computedServerField := findField(fields, "computed_server")
	if !computedServerField.Computed {
		t.Fatalf("expected readOnly field to be computed")
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
			name:     "read only field",
			property: &Schema{ReadOnly: true},
			want:     true,
		},
		{
			name:     "read only write only field",
			property: &Schema{ReadOnly: true, WriteOnly: true},
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
