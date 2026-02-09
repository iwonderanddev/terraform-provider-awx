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
