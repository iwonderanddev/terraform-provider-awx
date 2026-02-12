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

	awxLinks, err := awxOfficialLinksForObject(object.Name)
	if err != nil {
		t.Fatalf("awxOfficialLinksForObject returned error: %v", err)
	}
	if err := writeResourceDoc(resourceDir, object, objectDocsEnrichment{}, awxLinks, map[string]struct{}{}); err != nil {
		t.Fatalf("writeResourceDoc returned error: %v", err)
	}

	docPath := filepath.Join(resourceDir, "awx_organization.md")
	raw, err := os.ReadFile(docPath)
	if err != nil {
		t.Fatalf("failed to read generated resource doc: %v", err)
	}
	content := string(raw)
	if !strings.Contains(content, "`max_hosts` (Number, Optional, Computed)") {
		t.Fatalf("expected computed marker in generated docs, got:\n%s", content)
	}
	if strings.Contains(content, "Argument qualifiers used below") {
		t.Fatalf("did not expect legacy qualifier phrasing in generated docs, got:\n%s", content)
	}
	if !strings.Contains(content, "userguide/organizations.html") {
		t.Fatalf("expected resource-specific official AWX link in docs, got:\n%s", content)
	}
}

func TestWriteSettingsResourceDocDefaultsToAllAndIncludesScopeGuidance(t *testing.T) {
	t.Parallel()

	resourceDir := t.TempDir()
	object := manifest.ManagedObject{
		Name:             "settings",
		ResourceName:     "awx_setting",
		CollectionCreate: false,
	}

	awxLinks, err := awxOfficialLinksForObject(object.Name)
	if err != nil {
		t.Fatalf("awxOfficialLinksForObject returned error: %v", err)
	}
	if err := writeResourceDoc(resourceDir, object, objectDocsEnrichment{}, awxLinks, map[string]struct{}{}); err != nil {
		t.Fatalf("writeResourceDoc returned error: %v", err)
	}

	docPath := filepath.Join(resourceDir, "awx_setting.md")
	raw, err := os.ReadFile(docPath)
	if err != nil {
		t.Fatalf("failed to read generated resource doc: %v", err)
	}
	content := string(raw)
	if !strings.Contains(content, "id = \"all\"") {
		t.Fatalf("expected settings example to default to id=all, got:\n%s", content)
	}
	if !strings.Contains(content, "terraform import awx_setting.example all") {
		t.Fatalf("expected settings import guidance to default to all, got:\n%s", content)
	}
	for _, marker := range []string{
		"Category-scoped IDs",
		"optional advanced scoping",
		"overlapping ownership",
		"configuration conflicts",
	} {
		if !strings.Contains(content, marker) {
			t.Fatalf("expected settings guidance marker %q in doc, got:\n%s", marker, content)
		}
	}
}

func TestResolveFieldDescriptionPrefersCuratedThenSchemaThenFallback(t *testing.T) {
	t.Parallel()

	field := manifest.FieldSpec{
		Name:        "inventory",
		Type:        manifest.FieldTypeInt,
		Reference:   true,
		Description: "Schema description",
	}
	withCurated := resolveFieldDescription("inventory_id", field, objectDocsEnrichment{
		FieldDescriptions: map[string]string{
			"inventory_id": "Curated description",
		},
	})
	if withCurated != "Curated description" {
		t.Fatalf("expected curated description precedence, got=%q", withCurated)
	}

	withSchema := resolveFieldDescription("inventory_id", field, objectDocsEnrichment{})
	if withSchema != "Schema description" {
		t.Fatalf("expected schema description fallback, got=%q", withSchema)
	}

	field.Description = ""
	withFallback := resolveFieldDescription("inventory_id", field, objectDocsEnrichment{})
	if !strings.Contains(withFallback, "Numeric ID of the related AWX inventory object.") {
		t.Fatalf("expected typed fallback description, got=%q", withFallback)
	}
}

func TestObjectFieldDocTypeAndSampleValue(t *testing.T) {
	t.Parallel()

	label := terraformTypeLabel(manifest.FieldSpec{Type: manifest.FieldTypeObject})
	if label != "Object" {
		t.Fatalf("expected object field doc label to be Object, got=%q", label)
	}

	value := sampleValue(manifest.FieldTypeObject)
	if strings.Contains(value, "jsonencode(") {
		t.Fatalf("expected object sample value to use object literal, got=%q", value)
	}
	if value != "{ key = \"value\" }" {
		t.Fatalf("unexpected object sample value, got=%q", value)
	}
}

func TestReadDocsEnrichmentRejectsInvalidMetadata(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "docs_enrichment.json")
	payload := `{
  "priorityResources": ["awx_project"],
  "objects": {
    "awx_project": {
      "primaryExample": {"hcl": ""},
      "furtherReading": [{"title": "bad", "url": "not-a-url"}]
    }
  }
}`
	if err := os.WriteFile(path, []byte(payload), 0o644); err != nil {
		t.Fatalf("failed to write docs enrichment fixture: %v", err)
	}

	_, err := readDocsEnrichment(path)
	if err == nil {
		t.Fatalf("expected metadata validation error")
	}
	if !strings.Contains(err.Error(), "empty hcl content") {
		t.Fatalf("unexpected metadata validation error: %v", err)
	}
}

func TestReadDocsEnrichmentRejectsInvalidCurationSourceDate(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "docs_enrichment.json")
	payload := `{
  "priorityResources": [],
  "objects": {
    "awx_project": {
      "primaryExample": {
        "hcl": "resource \"awx_project\" \"example\" { name = \"demo\" }"
      },
      "curationSource": {
        "officialAwxUrl": "https://docs.ansible.com/projects/awx/en/24.6.1/userguide/projects.html",
        "verifiedOn": "2026/02/12"
      }
    }
  }
}`
	if err := os.WriteFile(path, []byte(payload), 0o644); err != nil {
		t.Fatalf("failed to write docs enrichment fixture: %v", err)
	}

	_, err := readDocsEnrichment(path)
	if err == nil {
		t.Fatalf("expected metadata validation error")
	}
	if !strings.Contains(err.Error(), "verifiedOn must use YYYY-MM-DD") {
		t.Fatalf("unexpected metadata validation error: %v", err)
	}
}

func TestAwxOfficialLinksForObjectRequiresKnownMapping(t *testing.T) {
	t.Parallel()

	if _, err := awxOfficialLinksForObject("projects"); err != nil {
		t.Fatalf("expected mapped object to succeed, got: %v", err)
	}

	if _, err := awxOfficialLinksForObject("nonexistent_object"); err == nil {
		t.Fatalf("expected unknown object mapping to fail")
	}
}

func TestValidateFurtherReadingPolicyRequiresOfficialAWXLinks(t *testing.T) {
	t.Parallel()

	valid := strings.Join([]string{
		"## Further Reading",
		"",
		"- [AWX Projects](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/projects.html)",
		"",
	}, "\n")
	if err := validateFurtherReadingPolicy("valid.md", valid, []docsLink{{
		Title: "AWX Projects",
		URL:   "https://docs.ansible.com/projects/awx/en/24.6.1/userguide/projects.html",
	}}); err != nil {
		t.Fatalf("expected valid Further Reading policy, got error: %v", err)
	}

	invalid := strings.Join([]string{
		"## Further Reading",
		"",
		"- [AWX Index](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/index.html)",
		"",
	}, "\n")
	if err := validateFurtherReadingPolicy("invalid.md", invalid, []docsLink{{
		Title: "AWX Projects",
		URL:   "https://docs.ansible.com/projects/awx/en/24.6.1/userguide/projects.html",
	}}); err == nil {
		t.Fatalf("expected generic AWX index to fail validation")
	}

	nonAWX := strings.Join([]string{
		"## Further Reading",
		"",
		"- [AWX Projects](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/projects.html)",
		"- [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework)",
		"",
	}, "\n")
	if err := validateFurtherReadingPolicy("non-awx.md", nonAWX, []docsLink{{
		Title: "AWX Projects",
		URL:   "https://docs.ansible.com/projects/awx/en/24.6.1/userguide/projects.html",
	}}); err == nil {
		t.Fatalf("expected non-AWX link to fail validation")
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

func TestValidateDocsEnrichmentTargetsRequiresPriorityCurationSource(t *testing.T) {
	t.Parallel()

	objects := []manifest.ManagedObject{{
		Name:             "projects",
		ResourceName:     "awx_project",
		DataSourceName:   "awx_project",
		ResourceEligible: true,
	}}

	err := validateDocsEnrichmentTargets(docsEnrichmentCatalog{
		PriorityResources: []string{"awx_project"},
		Objects: map[string]objectDocsEnrichment{
			"awx_project": {
				PrimaryExample: &docsExample{
					HCL: "resource \"awx_project\" \"example\" { name = \"demo\" }",
				},
			},
		},
	}, objects)
	if err == nil {
		t.Fatalf("expected priority curation source validation error")
	}
	if !strings.Contains(err.Error(), "requires curationSource") {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestValidateDocsEnrichmentTargetsRequiresMappedCurationSourceURL(t *testing.T) {
	t.Parallel()

	objects := []manifest.ManagedObject{{
		Name:             "projects",
		ResourceName:     "awx_project",
		DataSourceName:   "awx_project",
		ResourceEligible: true,
	}}

	err := validateDocsEnrichmentTargets(docsEnrichmentCatalog{
		PriorityResources: []string{"awx_project"},
		Objects: map[string]objectDocsEnrichment{
			"awx_project": {
				CurationSource: &docsSource{
					OfficialAWXURL: "https://docs.ansible.com/projects/awx/en/24.6.1/userguide/teams.html",
					VerifiedOn:     "2026-02-12",
				},
				PrimaryExample: &docsExample{
					HCL: "resource \"awx_project\" \"example\" { name = \"demo\" }",
				},
			},
		},
	}, objects)
	if err == nil {
		t.Fatalf("expected curation source mapping validation error")
	}
	if !strings.Contains(err.Error(), "must reference the mapped official AWX concept link") {
		t.Fatalf("unexpected validation error: %v", err)
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

func TestDataSourceExampleUsesAllForSettings(t *testing.T) {
	t.Parallel()

	object := manifest.ManagedObject{
		Name:             "settings",
		DataSourceName:   "awx_setting",
		CollectionCreate: false,
	}

	example := dataSourceExample(object)
	if !strings.Contains(example.HCL, "id = \"all\"") {
		t.Fatalf("expected settings data source example to default to id=all, got=%q", example.HCL)
	}
}

func TestValidateSettingsResourceDocumentation(t *testing.T) {
	t.Parallel()

	valid := strings.Join([]string{
		"## Example Usage",
		"",
		"```hcl",
		"resource \"awx_setting\" \"example\" {",
		"  id = \"all\"",
		"}",
		"```",
		"",
		"category-scoped IDs remain available for optional advanced scoping.",
		"Avoid overlapping ownership because this can cause configuration conflicts.",
		"",
		"## Import",
		"",
		"```bash",
		"terraform import awx_setting.example all",
		"```",
	}, "\n")

	if err := validateSettingsResourceDocumentation("settings-resource.md", valid); err != nil {
		t.Fatalf("expected valid settings resource doc, got: %v", err)
	}

	invalid := strings.ReplaceAll(valid, "id = \"all\"", "id = \"system\"")
	if err := validateSettingsResourceDocumentation("settings-resource.md", invalid); err == nil {
		t.Fatalf("expected invalid settings resource doc without id=all default")
	}
}

func TestValidateSettingsDataSourceDocumentation(t *testing.T) {
	t.Parallel()

	valid := strings.Join([]string{
		"## Example Usage",
		"",
		"```hcl",
		"data \"awx_setting\" \"example\" {",
		"  id = \"all\"",
		"}",
		"```",
		"",
		"category-scoped IDs remain available for optional advanced scoping.",
		"Avoid overlapping ownership because this can cause configuration conflicts.",
	}, "\n")

	if err := validateSettingsDataSourceDocumentation("settings-ds.md", valid); err != nil {
		t.Fatalf("expected valid settings data source doc, got: %v", err)
	}

	invalid := strings.ReplaceAll(valid, "configuration conflicts", "drift")
	if err := validateSettingsDataSourceDocumentation("settings-ds.md", invalid); err == nil {
		t.Fatalf("expected invalid settings data source doc without required conflict marker")
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

	awxLinks, err := awxOfficialLinksForRelationship(rel)
	if err != nil {
		t.Fatalf("awxOfficialLinksForRelationship returned error: %v", err)
	}
	if err := writeRelationshipDoc(resourceDir, rel, awxLinks); err != nil {
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
	if strings.Contains(content, "Breaking change:") {
		t.Fatalf("did not expect legacy breaking-change migration guidance, got:\n%s", content)
	}
	if !strings.Contains(content, "<primary_id>:<related_id>") {
		t.Fatalf("expected neutral composite ID placeholder in docs, got:\n%s", content)
	}
	if !strings.Contains(content, "## Schema") || !strings.Contains(content, "## Further Reading") {
		t.Fatalf("expected schema and further-reading sections in relationship docs, got:\n%s", content)
	}
	if !strings.Contains(content, "userguide/teams.html") || !strings.Contains(content, "userguide/users.html") {
		t.Fatalf("expected relationship docs to include parent/child official AWX links, got:\n%s", content)
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

	awxLinks, err := awxOfficialLinksForRelationship(rel)
	if err != nil {
		t.Fatalf("awxOfficialLinksForRelationship returned error: %v", err)
	}
	if err := writeRelationshipDoc(resourceDir, rel, awxLinks); err != nil {
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
	if strings.Contains(content, "Breaking change:") {
		t.Fatalf("did not expect legacy breaking-change migration guidance, got:\n%s", content)
	}
	if !strings.Contains(content, "<resource_id>") {
		t.Fatalf("expected survey-spec import placeholder in docs, got:\n%s", content)
	}
	if !strings.Contains(content, "userguide/job_templates.html") {
		t.Fatalf("expected survey-spec relationship docs to include job template official AWX link, got:\n%s", content)
	}
}

func TestGenerateDocsRendersSettingsDefaultsAndGuidance(t *testing.T) {
	t.Parallel()

	outputDir := t.TempDir()
	objects := []manifest.ManagedObject{
		{
			Name:             "settings",
			SingularName:     "setting",
			ResourceName:     "awx_setting",
			DataSourceName:   "awx_setting",
			ResourceEligible: true,
			DataSourceElig:   true,
			CollectionCreate: false,
		},
	}

	if err := generateDocs(outputDir, objects, nil, docsEnrichmentCatalog{}); err != nil {
		t.Fatalf("generateDocs returned error: %v", err)
	}

	resourcePath := filepath.Join(outputDir, "resources", "awx_setting.md")
	resourceRaw, err := os.ReadFile(resourcePath)
	if err != nil {
		t.Fatalf("failed to read generated resource doc: %v", err)
	}
	resourceContent := string(resourceRaw)
	if _, err := requireDocSectionsInOrder(resourcePath, "## Example Usage", "## Schema", "## Import", "## Further Reading"); err != nil {
		t.Fatalf("resource doc section ordering validation failed: %v", err)
	}
	if !strings.Contains(resourceContent, "id = \"all\"") {
		t.Fatalf("expected canonical settings resource example id=all, got:\n%s", resourceContent)
	}
	if !strings.Contains(resourceContent, "terraform import awx_setting.example all") {
		t.Fatalf("expected canonical settings resource import example with all, got:\n%s", resourceContent)
	}
	for _, marker := range []string{
		"Category-scoped IDs",
		"optional advanced scoping",
		"overlapping ownership",
		"configuration conflicts",
	} {
		if !strings.Contains(resourceContent, marker) {
			t.Fatalf("expected settings resource guidance marker %q, got:\n%s", marker, resourceContent)
		}
	}

	dataSourcePath := filepath.Join(outputDir, "data-sources", "awx_setting.md")
	dataSourceRaw, err := os.ReadFile(dataSourcePath)
	if err != nil {
		t.Fatalf("failed to read generated data source doc: %v", err)
	}
	dataSourceContent := string(dataSourceRaw)
	if _, err := requireDocSectionsInOrder(dataSourcePath, "## Example Usage", "## Schema", "## Further Reading"); err != nil {
		t.Fatalf("data source doc section ordering validation failed: %v", err)
	}
	if !strings.Contains(dataSourceContent, "id = \"all\"") {
		t.Fatalf("expected canonical settings data source example id=all, got:\n%s", dataSourceContent)
	}
	for _, marker := range []string{
		"Category-scoped IDs",
		"optional advanced scoping",
		"overlapping ownership",
		"configuration conflicts",
	} {
		if !strings.Contains(dataSourceContent, marker) {
			t.Fatalf("expected settings data source guidance marker %q, got:\n%s", marker, dataSourceContent)
		}
	}
}
