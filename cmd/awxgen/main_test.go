package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/damien/terraform-awx-provider/internal/manifest"
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

	if err := writeResourceDoc(resourceDir, object); err != nil {
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
