package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/damien/terraform-provider-awx-iwd/internal/manifest"
)

func TestBuildReportExcludesRuntimeDataSourcesFromManagedCoverage(t *testing.T) {
	t.Parallel()

	objects := []manifest.ManagedObject{
		{
			Name:             "projects",
			ResourceName:     "awx_project",
			DataSourceName:   "awx_project",
			ResourceEligible: true,
			DataSourceElig:   true,
			RuntimeExcluded:  false,
		},
		{
			Name:             "jobs",
			ResourceName:     "awx_job",
			DataSourceName:   "awx_job",
			ResourceEligible: false,
			DataSourceElig:   true,
			RuntimeExcluded:  true,
		},
	}

	report := buildReport("external/awx-openapi/schema.json", objects, nil, map[string]manifest.RuntimeExclusion{
		"jobs": {Object: "jobs", Reason: "runtime"},
	})

	if report.DataSourceEligible != 1 {
		t.Fatalf("unexpected data source eligible count: got=%d want=1", report.DataSourceEligible)
	}
	if report.RuntimeExcluded != 1 {
		t.Fatalf("unexpected runtime excluded count: got=%d want=1", report.RuntimeExcluded)
	}
	if len(report.ManagedDataSourceObjects) != 1 {
		t.Fatalf("unexpected managed data source object count: got=%d want=1", len(report.ManagedDataSourceObjects))
	}
	if report.ManagedDataSourceObjects[0] != "awx_project" {
		t.Fatalf("unexpected managed data source name: got=%q want=%q", report.ManagedDataSourceObjects[0], "awx_project")
	}
}

func TestWriteResourceDocIncludesComputedArgumentMarker(t *testing.T) {
	t.Parallel()

	resourceDir := t.TempDir()
	object := manifest.ManagedObject{
		Name:         "organizations",
		ResourceName: "awx_organization",
		Fields: []manifest.FieldSpec{
			{Name: "name", Type: manifest.FieldTypeString, Required: true},
			{Name: "max_hosts", Type: manifest.FieldTypeInt, Computed: true},
		},
	}

	if err := writeResourceDoc(resourceDir, object, map[string]struct{}{}); err != nil {
		t.Fatalf("writeResourceDoc returned error: %v", err)
	}

	docPath := filepath.Join(resourceDir, "awx_organization.md")
	raw, err := os.ReadFile(docPath)
	if err != nil {
		t.Fatalf("failed to read generated resource doc: %v", err)
	}
	content := string(raw)
	if !strings.Contains(content, "`max_hosts` (Optional, Computed)") {
		t.Fatalf("expected computed marker in generated docs, got:\n%s", content)
	}
}

func TestWriteProviderDocIncludesQualifierGuidance(t *testing.T) {
	t.Parallel()

	outputDir := t.TempDir()
	if err := writeProviderDoc(outputDir); err != nil {
		t.Fatalf("writeProviderDoc returned error: %v", err)
	}

	docPath := filepath.Join(outputDir, "index.md")
	raw, err := os.ReadFile(docPath)
	if err != nil {
		t.Fatalf("failed to read generated provider doc: %v", err)
	}
	content := string(raw)
	if !strings.Contains(content, "### Resource Argument Qualifiers") {
		t.Fatalf("expected resource argument qualifier section in provider docs, got:\n%s", content)
	}
	if !strings.Contains(content, "`Optional, Computed`") {
		t.Fatalf("expected Optional, Computed guidance in provider docs, got:\n%s", content)
	}
}

func TestFormatListItemDescriptionConvertsNestedBullets(t *testing.T) {
	t.Parallel()

	description := "Allowed values:\n\n* `always` - Always\n* `never` - Never"
	formatted := formatListItemDescription(description)

	if !strings.Contains(formatted, "Allowed values:") {
		t.Fatalf("expected primary description line, got=%q", formatted)
	}
	if !strings.Contains(formatted, "\n  - `always` - Always") {
		t.Fatalf("expected nested bullet conversion for first value, got=%q", formatted)
	}
	if !strings.Contains(formatted, "\n  - `never` - Never") {
		t.Fatalf("expected nested bullet conversion for second value, got=%q", formatted)
	}
}

func TestValidateTerraformFieldNameCollisionsDetectsSuffixConflicts(t *testing.T) {
	t.Parallel()

	err := validateTerraformFieldNameCollisions([]manifest.ManagedObject{
		{
			Name:             "teams",
			ResourceEligible: true,
			Fields: []manifest.FieldSpec{
				{Name: "organization", Type: manifest.FieldTypeInt, Reference: true},
				{Name: "organization_id", Type: manifest.FieldTypeInt, Reference: true},
			},
		},
	})
	if err == nil {
		t.Fatalf("expected collision error for duplicate Terraform attribute names")
	}
}

func TestSampleDocValueUsesReferenceWiringWhenTargetResourceExists(t *testing.T) {
	t.Parallel()

	field := manifest.FieldSpec{
		Name:      "organization",
		Type:      manifest.FieldTypeInt,
		Reference: true,
	}

	got := sampleDocValue(field, "organization_id", map[string]struct{}{
		"organization": {},
	})
	if got != "awx_organization.example.id" {
		t.Fatalf("unexpected reference wiring example: got=%q", got)
	}

	fallback := sampleDocValue(field, "organization_id", map[string]struct{}{})
	if fallback != "1" {
		t.Fatalf("expected numeric fallback example when target resource is unavailable, got=%q", fallback)
	}
}

func TestWriteRelationshipDocUsesCanonicalArguments(t *testing.T) {
	t.Parallel()

	resourceDir := t.TempDir()
	rel := manifest.Relationship{
		Name:              "team_user_association",
		ResourceName:      "awx_team_user_association",
		ParentObject:      "teams",
		ChildObject:       "users",
		ParentIDAttribute: "team_id",
		ChildIDAttribute:  "user_id",
		Path:              "/api/v2/teams/{id}/users/",
	}

	if err := writeRelationshipDoc(resourceDir, rel); err != nil {
		t.Fatalf("writeRelationshipDoc returned error: %v", err)
	}

	docPath := filepath.Join(resourceDir, "awx_team_user_association.md")
	raw, err := os.ReadFile(docPath)
	if err != nil {
		t.Fatalf("failed to read generated relationship doc: %v", err)
	}
	content := string(raw)
	if !strings.Contains(content, "team_id = 12") {
		t.Fatalf("expected parent canonical argument in example, got:\n%s", content)
	}
	if !strings.Contains(content, "user_id") {
		t.Fatalf("expected child canonical argument in doc, got:\n%s", content)
	}
	if !strings.Contains(content, "legacy `parent_id` and `child_id`") {
		t.Fatalf("expected breaking-change migration guidance, got:\n%s", content)
	}
	if strings.Contains(content, "- `parent_id` (Number, Required)") {
		t.Fatalf("expected legacy parent_id argument to be removed from argument docs, got:\n%s", content)
	}
}

func TestWriteRelationshipDocUsesCanonicalSurveySpecParentArgument(t *testing.T) {
	t.Parallel()

	resourceDir := t.TempDir()
	rel := manifest.Relationship{
		Name:              "job_template_survey_spec",
		ResourceName:      "awx_job_template_survey_spec",
		ParentObject:      "job_templates",
		ChildObject:       "survey_spec",
		ParentIDAttribute: "job_template_id",
		Path:              "/api/v2/job_templates/{id}/survey_spec/",
	}

	if err := writeRelationshipDoc(resourceDir, rel); err != nil {
		t.Fatalf("writeRelationshipDoc returned error: %v", err)
	}

	docPath := filepath.Join(resourceDir, "awx_job_template_survey_spec.md")
	raw, err := os.ReadFile(docPath)
	if err != nil {
		t.Fatalf("failed to read generated survey-spec relationship doc: %v", err)
	}
	content := string(raw)
	if !strings.Contains(content, "job_template_id = 12") {
		t.Fatalf("expected canonical survey-spec parent argument in example, got:\n%s", content)
	}
	if !strings.Contains(content, "legacy `parent_id`") {
		t.Fatalf("expected survey-spec migration guidance, got:\n%s", content)
	}
	if strings.Contains(content, "- `parent_id` (Number, Required)") {
		t.Fatalf("expected legacy parent_id argument to be removed from argument docs, got:\n%s", content)
	}
}
