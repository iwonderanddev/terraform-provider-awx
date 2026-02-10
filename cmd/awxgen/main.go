package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/damien/terraform-awx-provider/internal/manifest"
	"github.com/damien/terraform-awx-provider/internal/openapi"
)

const (
	defaultSchemaPath         = "external/awx-openapi/schema.json"
	defaultManagedPath        = "internal/manifest/managed_objects.json"
	defaultRelationshipsPath  = "internal/manifest/relationships.json"
	defaultExclusionsPath     = "internal/manifest/runtime_exclusions.json"
	defaultPrioritiesPath     = "internal/manifest/relationship_priorities.json"
	defaultOverridesPath      = "internal/manifest/field_overrides.json"
	defaultCoverageReportPath = "internal/manifest/coverage_report.json"
)

// CoverageReport summarizes manifest and exclusion coverage.
type CoverageReport struct {
	GeneratedAt               string   `json:"generatedAt"`
	SchemaPath                string   `json:"schemaPath"`
	TotalCandidates           int      `json:"totalCandidates"`
	ResourceEligible          int      `json:"resourceEligible"`
	DataSourceEligible        int      `json:"dataSourceEligible"`
	RuntimeExcluded           int      `json:"runtimeExcluded"`
	ManagedResourceObjects    []string `json:"managedResourceObjects"`
	ManagedDataSourceObjects  []string `json:"managedDataSourceObjects"`
	MissingRuntimeExclusions  []string `json:"missingRuntimeExclusions"`
	RelationshipResourceCount int      `json:"relationshipResourceCount"`
	RelationshipResourceNames []string `json:"relationshipResourceNames"`
}

func main() {
	if len(os.Args) < 2 {
		exitWithError(errors.New("usage: awxgen <generate|validate|docs|docs-validate|report>"))
	}

	var err error
	switch os.Args[1] {
	case "generate":
		err = runGenerate(os.Args[2:])
	case "validate":
		err = runValidate(os.Args[2:])
	case "docs":
		err = runDocs(os.Args[2:])
	case "docs-validate":
		err = runDocsValidate(os.Args[2:])
	case "report":
		err = runReport(os.Args[2:])
	default:
		err = fmt.Errorf("unknown command %q", os.Args[1])
	}

	if err != nil {
		exitWithError(err)
	}
}

func runGenerate(args []string) error {
	fs := flag.NewFlagSet("generate", flag.ContinueOnError)
	schemaPath := fs.String("schema", defaultSchemaPath, "Path to AWX OpenAPI schema JSON")
	exclusionsPath := fs.String("exclusions", defaultExclusionsPath, "Path to runtime exclusions JSON")
	prioritiesPath := fs.String("relationship-priorities", defaultPrioritiesPath, "Path to relationship priority JSON")
	overridesPath := fs.String("overrides", defaultOverridesPath, "Path to field override JSON")
	managedPath := fs.String("managed", defaultManagedPath, "Output path for managed object manifest")
	relationshipsPath := fs.String("relationships", defaultRelationshipsPath, "Output path for relationship manifest")
	reportPath := fs.String("report", defaultCoverageReportPath, "Output path for coverage report")
	if err := fs.Parse(args); err != nil {
		return err
	}

	runtimeExclusions, priorities, objects, relationships, report, overrideCount, err := generate(*schemaPath, *exclusionsPath, *prioritiesPath, *overridesPath)
	if err != nil {
		return err
	}

	if err := writePrettyJSON(*managedPath, objects); err != nil {
		return fmt.Errorf("write managed object manifest: %w", err)
	}
	if err := writePrettyJSON(*relationshipsPath, relationships); err != nil {
		return fmt.Errorf("write relationship manifest: %w", err)
	}
	if err := writePrettyJSON(*reportPath, report); err != nil {
		return fmt.Errorf("write coverage report: %w", err)
	}

	fmt.Printf("Generated managed objects: %d (resource eligible: %d, data source eligible: %d)\n", len(objects), report.ResourceEligible, report.DataSourceEligible)
	fmt.Printf("Generated relationship resources: %d\n", len(relationships))
	fmt.Printf("Runtime exclusions loaded: %d\n", len(runtimeExclusions))
	if len(priorities) > 0 {
		fmt.Printf("Relationship priorities loaded: %d\n", len(priorities))
	}
	if overrideCount > 0 {
		fmt.Printf("Field overrides loaded: %d\n", overrideCount)
	}
	if len(report.MissingRuntimeExclusions) > 0 {
		fmt.Printf("Missing runtime exclusions (%d): %s\n", len(report.MissingRuntimeExclusions), strings.Join(report.MissingRuntimeExclusions, ", "))
		return errors.New("coverage validation failed: runtime exclusions are incomplete")
	}
	return nil
}

func runValidate(args []string) error {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	schemaPath := fs.String("schema", defaultSchemaPath, "Path to AWX OpenAPI schema JSON")
	exclusionsPath := fs.String("exclusions", defaultExclusionsPath, "Path to runtime exclusions JSON")
	prioritiesPath := fs.String("relationship-priorities", defaultPrioritiesPath, "Path to relationship priority JSON")
	overridesPath := fs.String("overrides", defaultOverridesPath, "Path to field override JSON")
	managedPath := fs.String("managed", defaultManagedPath, "Managed object manifest path")
	relationshipsPath := fs.String("relationships", defaultRelationshipsPath, "Relationship manifest path")
	reportPath := fs.String("report", defaultCoverageReportPath, "Coverage report path")
	if err := fs.Parse(args); err != nil {
		return err
	}

	_, _, generatedObjects, generatedRelationships, report, _, err := generate(*schemaPath, *exclusionsPath, *prioritiesPath, *overridesPath)
	if err != nil {
		return err
	}

	if len(report.MissingRuntimeExclusions) > 0 {
		return fmt.Errorf("coverage validation failed: runtime exclusions are missing for: %s", strings.Join(report.MissingRuntimeExclusions, ", "))
	}

	if err := compareJSONFile(*managedPath, generatedObjects); err != nil {
		return fmt.Errorf("managed object manifest validation failed: %w", err)
	}
	if err := compareJSONFile(*relationshipsPath, generatedRelationships); err != nil {
		return fmt.Errorf("relationship manifest validation failed: %w", err)
	}
	if err := compareJSONFile(*reportPath, report); err != nil {
		return fmt.Errorf("coverage report validation failed: %w", err)
	}

	fmt.Println("Manifest validation passed.")
	return nil
}

func runDocs(args []string) error {
	fs := flag.NewFlagSet("docs", flag.ContinueOnError)
	managedPath := fs.String("managed", defaultManagedPath, "Managed object manifest path")
	relationshipsPath := fs.String("relationships", defaultRelationshipsPath, "Relationship manifest path")
	outputDir := fs.String("out", "docs", "Documentation root directory")
	if err := fs.Parse(args); err != nil {
		return err
	}

	objects, err := readManagedObjects(*managedPath)
	if err != nil {
		return err
	}
	relationships, err := readRelationships(*relationshipsPath)
	if err != nil {
		return err
	}

	if err := generateDocs(*outputDir, objects, relationships); err != nil {
		return err
	}
	fmt.Printf("Generated docs in %s\n", *outputDir)
	return nil
}

func runDocsValidate(args []string) error {
	fs := flag.NewFlagSet("docs-validate", flag.ContinueOnError)
	managedPath := fs.String("managed", defaultManagedPath, "Managed object manifest path")
	relationshipsPath := fs.String("relationships", defaultRelationshipsPath, "Relationship manifest path")
	docsDir := fs.String("docs", "docs", "Documentation root directory")
	if err := fs.Parse(args); err != nil {
		return err
	}

	objects, err := readManagedObjects(*managedPath)
	if err != nil {
		return err
	}
	relationships, err := readRelationships(*relationshipsPath)
	if err != nil {
		return err
	}

	providerDoc := filepath.Join(*docsDir, "index.md")
	if _, err := os.Stat(providerDoc); err != nil {
		return fmt.Errorf("provider documentation is missing: %s", providerDoc)
	}

	for _, object := range objects {
		if object.ResourceEligible && !object.RuntimeExcluded {
			resourceDoc := filepath.Join(*docsDir, "resources", fmt.Sprintf("%s.md", object.ResourceName))
			if err := requireDocSections(resourceDoc, "## Example Usage", "## Import"); err != nil {
				return err
			}
		}
		if object.DataSourceElig && !object.RuntimeExcluded {
			dataSourceDoc := filepath.Join(*docsDir, "data-sources", fmt.Sprintf("%s.md", object.DataSourceName))
			if err := requireDocSections(dataSourceDoc, "## Example Usage", "## Attributes Reference"); err != nil {
				return err
			}
		}
	}
	for _, relationship := range relationships {
		resourceDoc := filepath.Join(*docsDir, "resources", fmt.Sprintf("%s.md", relationship.ResourceName))
		if err := requireDocSections(resourceDoc, "## Example Usage", "## Import"); err != nil {
			return err
		}
	}

	fmt.Println("Documentation validation passed.")
	return nil
}

func runReport(args []string) error {
	fs := flag.NewFlagSet("report", flag.ContinueOnError)
	reportPath := fs.String("report", defaultCoverageReportPath, "Coverage report path")
	if err := fs.Parse(args); err != nil {
		return err
	}

	raw, err := os.ReadFile(filepath.Clean(*reportPath))
	if err != nil {
		return err
	}
	var report CoverageReport
	if err := json.Unmarshal(raw, &report); err != nil {
		return err
	}

	fmt.Printf("Generated at: %s\n", report.GeneratedAt)
	fmt.Printf("Candidates: %d\n", report.TotalCandidates)
	fmt.Printf("Resource eligible: %d\n", report.ResourceEligible)
	fmt.Printf("Data source eligible: %d\n", report.DataSourceEligible)
	fmt.Printf("Runtime excluded: %d\n", report.RuntimeExcluded)
	fmt.Printf("Relationship resources: %d\n", report.RelationshipResourceCount)
	if len(report.MissingRuntimeExclusions) > 0 {
		fmt.Printf("Missing exclusions: %s\n", strings.Join(report.MissingRuntimeExclusions, ", "))
		return errors.New("coverage report indicates missing runtime exclusions")
	}
	return nil
}

func generate(schemaPath string, exclusionsPath string, prioritiesPath string, overridesPath string) (map[string]manifest.RuntimeExclusion, map[string]int, []manifest.ManagedObject, []manifest.Relationship, CoverageReport, int, error) {
	doc, err := openapi.LoadDocument(schemaPath)
	if err != nil {
		return nil, nil, nil, nil, CoverageReport{}, 0, err
	}

	runtimeExclusions, err := openapi.LoadRuntimeExclusions(exclusionsPath)
	if err != nil {
		return nil, nil, nil, nil, CoverageReport{}, 0, err
	}

	priorities, err := openapi.LoadRelationshipPriorities(prioritiesPath)
	if err != nil {
		return nil, nil, nil, nil, CoverageReport{}, 0, err
	}
	fieldOverrides, err := openapi.LoadFieldOverrides(overridesPath)
	if err != nil {
		return nil, nil, nil, nil, CoverageReport{}, 0, err
	}

	objects := openapi.DeriveManagedObjects(doc, runtimeExclusions)
	objects = openapi.ApplyFieldOverrides(objects, fieldOverrides)
	if err := openapi.ValidateCoverage(objects, runtimeExclusions); err != nil {
		return runtimeExclusions, priorities, objects, nil, buildReport(schemaPath, objects, nil, runtimeExclusions), len(fieldOverrides), err
	}

	relationships := openapi.DeriveRelationships(doc, objects, priorities)
	report := buildReport(schemaPath, objects, relationships, runtimeExclusions)
	return runtimeExclusions, priorities, objects, relationships, report, len(fieldOverrides), nil
}

func buildReport(schemaPath string, objects []manifest.ManagedObject, relationships []manifest.Relationship, exclusions map[string]manifest.RuntimeExclusion) CoverageReport {
	report := CoverageReport{
		GeneratedAt:               time.Now().UTC().Format("2006-01-02"),
		SchemaPath:                schemaPath,
		TotalCandidates:           len(objects),
		ManagedResourceObjects:    make([]string, 0),
		ManagedDataSourceObjects:  make([]string, 0),
		MissingRuntimeExclusions:  make([]string, 0),
		RelationshipResourceCount: len(relationships),
		RelationshipResourceNames: make([]string, 0, len(relationships)),
	}

	for _, relationship := range relationships {
		report.RelationshipResourceNames = append(report.RelationshipResourceNames, relationship.ResourceName)
	}
	for _, obj := range objects {
		if obj.ResourceEligible {
			report.ResourceEligible++
			report.ManagedResourceObjects = append(report.ManagedResourceObjects, obj.ResourceName)
		} else if !obj.RuntimeExcluded {
			if _, exists := exclusions[obj.Name]; !exists && shouldRequireRuntimeExclusion(obj.Name) {
				report.MissingRuntimeExclusions = append(report.MissingRuntimeExclusions, obj.Name)
			}
		}
		if obj.DataSourceElig && !obj.RuntimeExcluded {
			report.DataSourceEligible++
			report.ManagedDataSourceObjects = append(report.ManagedDataSourceObjects, obj.DataSourceName)
		}
		if obj.RuntimeExcluded {
			report.RuntimeExcluded++
		}
	}

	sort.Strings(report.ManagedResourceObjects)
	sort.Strings(report.ManagedDataSourceObjects)
	sort.Strings(report.MissingRuntimeExclusions)
	sort.Strings(report.RelationshipResourceNames)
	return report
}

func shouldRequireRuntimeExclusion(objectName string) bool {
	name := strings.ToLower(objectName)
	keywords := []string{
		"activity",
		"ad_hoc",
		"analytics",
		"dashboard",
		"event",
		"fact",
		"history",
		"instance",
		"job",
		"metric",
		"receptor",
		"schedule_preview",
		"task",
		"workflow_approval",
	}
	for _, keyword := range keywords {
		if strings.Contains(name, keyword) {
			return true
		}
	}
	return false
}

func generateDocs(outputDir string, objects []manifest.ManagedObject, relationships []manifest.Relationship) error {
	resourceDir := filepath.Join(outputDir, "resources")
	dataSourceDir := filepath.Join(outputDir, "data-sources")
	if err := os.MkdirAll(resourceDir, 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(dataSourceDir, 0o755); err != nil {
		return err
	}

	if err := writeProviderDoc(outputDir); err != nil {
		return err
	}

	for _, obj := range objects {
		if obj.ResourceEligible && !obj.RuntimeExcluded {
			if err := writeResourceDoc(resourceDir, obj); err != nil {
				return err
			}
		}
		if obj.DataSourceElig && !obj.RuntimeExcluded {
			if err := writeDataSourceDoc(dataSourceDir, obj); err != nil {
				return err
			}
		}
	}
	for _, rel := range relationships {
		if err := writeRelationshipDoc(resourceDir, rel); err != nil {
			return err
		}
	}
	return nil
}

func writeProviderDoc(outputDir string) error {
	contents := `# Provider: awx

The ` + "`awx`" + ` provider manages AWX 24.6.1 objects via API v2.

## Example Usage

` + "```hcl" + `
provider "awx" {
  base_url  = var.awx_base_url
  username  = var.awx_username
  password  = var.awx_password
}
` + "```" + `

## Schema

### Required

- ` + "`base_url`" + ` (String) AWX base URL, for example ` + "`https://awx.example.com`" + `.
- ` + "`username`" + ` (String) HTTP Basic username.
- ` + "`password`" + ` (String, Sensitive) HTTP Basic password.

### Optional

- ` + "`insecure_skip_tls_verify`" + ` (Boolean) Skip TLS verification.
- ` + "`ca_cert_pem`" + ` (String, Sensitive) PEM CA certificate bundle.
- ` + "`request_timeout_seconds`" + ` (Number) API request timeout.
- ` + "`retry_max_attempts`" + ` (Number) Retry attempts for retryable failures.
- ` + "`retry_backoff_millis`" + ` (Number) Initial retry backoff in milliseconds.

### Resource Argument Qualifiers

Generated resource docs under ` + "`docs/resources/*`" + ` use these qualifiers:

- ` + "`Required`" + `: Must be set in configuration.
- ` + "`Optional`" + `: May be omitted.
- ` + "`Optional, Computed`" + `: May be omitted; AWX may apply a server-side default and Terraform records the resulting value in state after apply.

## Compatibility

This provider targets AWX 24.6.1 API v2 only. Runtime-only objects are excluded from managed resources.
`

	return os.WriteFile(filepath.Join(outputDir, "index.md"), []byte(contents), 0o644)
}

func writeResourceDoc(resourceDir string, obj manifest.ManagedObject) error {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("# Resource: %s\n\n", obj.ResourceName))
	builder.WriteString(fmt.Sprintf("Manages AWX `%s` objects.\n\n", obj.Name))

	builder.WriteString("## Example Usage\n\n")
	builder.WriteString("```hcl\n")
	builder.WriteString(fmt.Sprintf("resource \"%s\" \"example\" {\n", obj.ResourceName))
	if !obj.CollectionCreate {
		builder.WriteString("  id = \"example\"\n")
	}
	for _, field := range obj.Fields {
		if field.Required {
			builder.WriteString(fmt.Sprintf("  %s = %s\n", field.Name, sampleValue(field.Type)))
		}
	}
	builder.WriteString("}\n")
	builder.WriteString("```\n\n")

	builder.WriteString("## Argument Reference\n\n")
	builder.WriteString("Argument qualifiers used below:\n")
	builder.WriteString("- `Required`: Must be set in configuration.\n")
	builder.WriteString("- `Optional`: May be omitted.\n")
	builder.WriteString("- `Optional, Computed`: May be omitted; AWX can apply a server-side default and Terraform records the resulting value after apply.\n\n")
	if !obj.CollectionCreate {
		builder.WriteString("- `id` (Required) AWX detail-path identifier for this object.\n")
	}
	for _, field := range obj.Fields {
		required := "Optional"
		if field.Required {
			required = "Required"
		} else if field.Computed {
			required = "Optional, Computed"
		}
		sensitive := ""
		if field.Sensitive {
			sensitive = ", Sensitive"
		}
		description := strings.TrimSpace(field.Description)
		if description == "" {
			description = "Managed field from AWX OpenAPI schema."
		}
		builder.WriteString(fmt.Sprintf("- `%s` (%s%s) %s\n", field.Name, required, sensitive, description))
	}
	builder.WriteString("\n## Attributes Reference\n\n")
	if obj.CollectionCreate {
		builder.WriteString("- `id` (String) Numeric AWX object identifier.\n")
	} else {
		builder.WriteString("- `id` (String) AWX detail-path identifier for this object.\n")
	}
	builder.WriteString("\n## Import\n\n")
	builder.WriteString("```bash\n")
	importID := "42"
	if !obj.CollectionCreate {
		importID = "example"
	}
	builder.WriteString(fmt.Sprintf("terraform import %s.example %s\n", obj.ResourceName, importID))
	builder.WriteString("```\n")

	return os.WriteFile(filepath.Join(resourceDir, fmt.Sprintf("%s.md", obj.ResourceName)), []byte(builder.String()), 0o644)
}

func writeDataSourceDoc(dataSourceDir string, obj manifest.ManagedObject) error {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("# Data Source: %s\n\n", obj.DataSourceName))
	builder.WriteString(fmt.Sprintf("Reads AWX `%s` objects.\n\n", obj.Name))
	builder.WriteString("## Example Usage\n\n")
	builder.WriteString("```hcl\n")
	builder.WriteString(fmt.Sprintf("data \"%s\" \"example\" {\n", obj.DataSourceName))
	idExample := sampleValue(manifest.FieldTypeInt)
	if !obj.CollectionCreate {
		idExample = sampleValue(manifest.FieldTypeString)
	}
	builder.WriteString("  id = " + idExample + "\n")
	builder.WriteString("}\n")
	builder.WriteString("```\n\n")
	builder.WriteString("## Argument Reference\n\n")
	if obj.CollectionCreate {
		builder.WriteString("- `id` (String, Optional) Numeric AWX object ID.\n")
	} else {
		builder.WriteString("- `id` (String, Optional) AWX object identifier used in the detail endpoint path.\n")
	}
	if hasField(obj.Fields, "name") {
		builder.WriteString("- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.\n")
	}
	builder.WriteString("\n## Attributes Reference\n\n")
	if obj.CollectionCreate {
		builder.WriteString("- `id` (String) Numeric AWX object ID.\n")
	} else {
		builder.WriteString("- `id` (String) AWX detail-path identifier for this object.\n")
	}
	for _, field := range obj.Fields {
		sensitive := ""
		if field.Sensitive {
			sensitive = ", Sensitive"
		}
		builder.WriteString(fmt.Sprintf("- `%s` (%s%s)\n", field.Name, field.Type, sensitive))
	}

	return os.WriteFile(filepath.Join(dataSourceDir, fmt.Sprintf("%s.md", obj.DataSourceName)), []byte(builder.String()), 0o644)
}

func writeRelationshipDoc(resourceDir string, rel manifest.Relationship) error {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("# Resource: %s\n\n", rel.ResourceName))
	if isSurveySpecRelationship(rel) {
		builder.WriteString(fmt.Sprintf("Manages `%s` survey specification for `%s` objects.\n\n", rel.Name, rel.ParentObject))
		builder.WriteString("## Example Usage\n\n")
		builder.WriteString("```hcl\n")
		builder.WriteString(fmt.Sprintf("resource \"%s\" \"example\" {\n", rel.ResourceName))
		builder.WriteString("  parent_id = 12\n")
		builder.WriteString("  spec = jsonencode({\n")
		builder.WriteString("    name        = \"Example survey\"\n")
		builder.WriteString("    description = \"Managed by Terraform\"\n")
		builder.WriteString("    spec        = []\n")
		builder.WriteString("  })\n")
		builder.WriteString("}\n")
		builder.WriteString("```\n\n")
		builder.WriteString("## Argument Reference\n\n")
		builder.WriteString("- `parent_id` (Number, Required) Parent object numeric ID.\n")
		builder.WriteString("- `spec` (String, Optional) JSON-encoded survey specification payload.\n\n")
		builder.WriteString("## Attributes Reference\n\n")
		builder.WriteString("- `id` (String) Survey specification ID (same as `parent_id`).\n")
		builder.WriteString("- `parent_id` (Number) Parent object numeric ID.\n")
		builder.WriteString("- `spec` (String) JSON-encoded survey specification payload.\n\n")
		builder.WriteString("## Import\n\n")
		builder.WriteString("```bash\n")
		builder.WriteString(fmt.Sprintf("terraform import %s.example 12\n", rel.ResourceName))
		builder.WriteString("```\n")
	} else {
		builder.WriteString(fmt.Sprintf("Manages `%s` relationships between `%s` and `%s` objects.\n\n", rel.Name, rel.ParentObject, rel.ChildObject))
		builder.WriteString("## Example Usage\n\n")
		builder.WriteString("```hcl\n")
		builder.WriteString(fmt.Sprintf("resource \"%s\" \"example\" {\n", rel.ResourceName))
		builder.WriteString("  parent_id = 12\n")
		builder.WriteString("  child_id  = 34\n")
		builder.WriteString("}\n")
		builder.WriteString("```\n\n")
		builder.WriteString("## Argument Reference\n\n")
		builder.WriteString("- `parent_id` (Number, Required) Parent object numeric ID.\n")
		builder.WriteString("- `child_id` (Number, Required) Child object numeric ID.\n\n")
		builder.WriteString("## Attributes Reference\n\n")
		builder.WriteString("- `id` (String) Composite ID in `<parent_id>:<child_id>` format.\n\n")
		builder.WriteString("## Import\n\n")
		builder.WriteString("```bash\n")
		builder.WriteString(fmt.Sprintf("terraform import %s.example 12:34\n", rel.ResourceName))
		builder.WriteString("```\n")
	}

	return os.WriteFile(filepath.Join(resourceDir, fmt.Sprintf("%s.md", rel.ResourceName)), []byte(builder.String()), 0o644)
}

func isSurveySpecRelationship(rel manifest.Relationship) bool {
	return strings.HasSuffix(rel.Path, "/survey_spec/")
}

func sampleValue(fieldType manifest.FieldType) string {
	switch fieldType {
	case manifest.FieldTypeBool:
		return "true"
	case manifest.FieldTypeFloat:
		return "1.0"
	case manifest.FieldTypeInt:
		return "1"
	case manifest.FieldTypeArray:
		return "jsonencode([\"value\"])"
	case manifest.FieldTypeObject:
		return "jsonencode({ key = \"value\" })"
	default:
		return "\"example\""
	}
}

func hasField(fields []manifest.FieldSpec, name string) bool {
	for _, field := range fields {
		if field.Name == name {
			return true
		}
	}
	return false
}

func readManagedObjects(path string) ([]manifest.ManagedObject, error) {
	raw, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	var out []manifest.ManagedObject
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func readRelationships(path string) ([]manifest.Relationship, error) {
	raw, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	var out []manifest.Relationship
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func requireDocSections(path string, sections ...string) error {
	raw, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("documentation file missing: %s", path)
	}
	content := string(raw)
	for _, section := range sections {
		if !strings.Contains(content, section) {
			return fmt.Errorf("documentation file %s is missing section %q", path, section)
		}
	}
	return nil
}

func writePrettyJSON(path string, v any) error {
	raw, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	raw = append(raw, '\n')
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Clean(path), raw, 0o644)
}

func compareJSONFile(path string, expected any) error {
	expectedRaw, err := json.Marshal(expected)
	if err != nil {
		return err
	}

	actualRaw, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return err
	}

	if !jsonEqual(actualRaw, expectedRaw) {
		return fmt.Errorf("%s is out of date; run `go run ./cmd/awxgen generate`", path)
	}
	return nil
}

func jsonEqual(a []byte, b []byte) bool {
	var left any
	var right any
	if err := json.Unmarshal(a, &left); err != nil {
		return false
	}
	if err := json.Unmarshal(b, &right); err != nil {
		return false
	}
	leftRaw, _ := json.Marshal(left)
	rightRaw, _ := json.Marshal(right)
	return string(leftRaw) == string(rightRaw)
}

func exitWithError(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "error: %s\n", err)
	os.Exit(1)
}
