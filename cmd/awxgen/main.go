package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/damien/terraform-provider-awx-iwd/internal/manifest"
	"github.com/damien/terraform-provider-awx-iwd/internal/openapi"
)

const (
	defaultSchemaPath          = "external/awx-openapi/schema.json"
	defaultManagedPath         = "internal/manifest/managed_objects.json"
	defaultRelationshipsPath   = "internal/manifest/relationships.json"
	defaultExclusionsPath      = "internal/manifest/runtime_exclusions.json"
	defaultDeprecatedPath      = "internal/manifest/deprecated_exclusions.json"
	defaultPrioritiesPath      = "internal/manifest/relationship_priorities.json"
	defaultOverridesPath       = "internal/manifest/field_overrides.json"
	defaultCoverageReportPath  = "internal/manifest/coverage_report.json"
	defaultDocsEnrichmentPath  = "internal/manifest/docs_enrichment.json"
	defaultOnlineVerifyTimeout = 10 * time.Second
	defaultMaxQualityPasses    = 3
	openAPIPlaceholderText     = "Managed field from AWX OpenAPI schema"
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

type docsEnrichmentCatalog struct {
	PriorityResources []string                        `json:"priorityResources"`
	Objects           map[string]objectDocsEnrichment `json:"objects"`
}

type objectDocsEnrichment struct {
	Overview                   string                   `json:"overview,omitempty"`
	Complex                    bool                     `json:"complex,omitempty"`
	ConceptPrimer              string                   `json:"conceptPrimer,omitempty"`
	FieldDescriptions          map[string]string        `json:"fieldDescriptions,omitempty"`
	FurtherReading             []docsLink               `json:"furtherReading,omitempty"`
	CurationSource             *docsSource              `json:"curationSource,omitempty"`
	OnlineResearchChecklist    *onlineResearchChecklist `json:"onlineResearchChecklist,omitempty"`
	InteractionReferenceFields []string                 `json:"interactionReferenceFields,omitempty"`
	PrimaryExample             *docsExample             `json:"primaryExample,omitempty"`
	ExtraExamples              []docsExample            `json:"extraExamples,omitempty"`
}

type docsSource struct {
	OfficialAWXURL string `json:"officialAwxUrl"`
	VerifiedOn     string `json:"verifiedOn"`
}

type onlineResearchChecklist struct {
	ObjectBehavior      string `json:"objectBehavior"`
	RelatedInteractions string `json:"relatedInteractions"`
	ParameterSemantics  string `json:"parameterSemantics"`
}

type docsLink struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type docsExample struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	HCL         string `json:"hcl"`
}

var markdownLinkPattern = regexp.MustCompile(`\[[^\]]+\]\(([^)]+)\)`)
var hclFencePattern = regexp.MustCompile("(?m)^```hcl\\s*$")
var inlineEnumBulletPattern = regexp.MustCompile(`\s+\*\s+`)
var unresolvedReferencePattern = regexp.MustCompile(`(?m)(?:^|[^\w])(data\.)?(awx_[a-z0-9_]+)\.([A-Za-z0-9_]+)\.id(?:$|[^\w])`)
var resourceBlockPattern = regexp.MustCompile(`(?m)^\s*resource\s+"(awx_[a-z0-9_]+)"\s+"([A-Za-z0-9_]+)"\s*\{`)
var dataBlockPattern = regexp.MustCompile(`(?m)^\s*data\s+"(awx_[a-z0-9_]+)"\s+"([A-Za-z0-9_]+)"\s*\{`)
var curatedReferenceAssignmentPattern = regexp.MustCompile(`(?m)^\s*([a-z0-9_]+)\s*=\s*(?:data\.)?awx_[a-z0-9_]+\.[A-Za-z0-9_]+\.id\s*$`)
var malformedEnumInlinePattern = regexp.MustCompile("(?m)^- `[^`]+` \\([^)]+\\) [^\\n]*\\* ")
var qualityAnalysisPassHeadingPattern = regexp.MustCompile(`(?m)^## Quality Analysis Pass ([1-9][0-9]*)\s*$`)

func main() {
	if len(os.Args) < 2 {
		exitWithError(errors.New("usage: awxgen <generate|validate|docs|docs-validate|docs-verify-online|docs-validate-quality|report>"))
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
	case "docs-verify-online":
		err = runDocsVerifyOnline(os.Args[2:])
	case "docs-validate-quality":
		err = runDocsValidateQuality(os.Args[2:])
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
	deprecatedPath := fs.String("deprecated", defaultDeprecatedPath, "Path to deprecated endpoint exclusions JSON")
	prioritiesPath := fs.String("relationship-priorities", defaultPrioritiesPath, "Path to relationship priority JSON")
	overridesPath := fs.String("overrides", defaultOverridesPath, "Path to field override JSON")
	managedPath := fs.String("managed", defaultManagedPath, "Output path for managed object manifest")
	relationshipsPath := fs.String("relationships", defaultRelationshipsPath, "Output path for relationship manifest")
	reportPath := fs.String("report", defaultCoverageReportPath, "Output path for coverage report")
	if err := fs.Parse(args); err != nil {
		return err
	}

	runtimeExclusions, priorities, objects, relationships, report, overrideCount, err := generate(*schemaPath, *exclusionsPath, *deprecatedPath, *prioritiesPath, *overridesPath)
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
	deprecatedPath := fs.String("deprecated", defaultDeprecatedPath, "Path to deprecated endpoint exclusions JSON")
	prioritiesPath := fs.String("relationship-priorities", defaultPrioritiesPath, "Path to relationship priority JSON")
	overridesPath := fs.String("overrides", defaultOverridesPath, "Path to field override JSON")
	managedPath := fs.String("managed", defaultManagedPath, "Managed object manifest path")
	relationshipsPath := fs.String("relationships", defaultRelationshipsPath, "Relationship manifest path")
	reportPath := fs.String("report", defaultCoverageReportPath, "Coverage report path")
	if err := fs.Parse(args); err != nil {
		return err
	}

	_, _, generatedObjects, generatedRelationships, report, _, err := generate(*schemaPath, *exclusionsPath, *deprecatedPath, *prioritiesPath, *overridesPath)
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
	docsEnrichmentPath := fs.String("docs-enrichment", defaultDocsEnrichmentPath, "Documentation enrichment metadata path")
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
	enrichment, err := readDocsEnrichment(*docsEnrichmentPath)
	if err != nil {
		return err
	}
	if err := validateDocsEnrichmentTargets(enrichment, objects); err != nil {
		return err
	}

	if err := generateDocs(*outputDir, objects, relationships, enrichment); err != nil {
		return err
	}
	fmt.Printf("Generated docs in %s\n", *outputDir)
	return nil
}

func runDocsValidate(args []string) error {
	fs := flag.NewFlagSet("docs-validate", flag.ContinueOnError)
	managedPath := fs.String("managed", defaultManagedPath, "Managed object manifest path")
	relationshipsPath := fs.String("relationships", defaultRelationshipsPath, "Relationship manifest path")
	docsEnrichmentPath := fs.String("docs-enrichment", defaultDocsEnrichmentPath, "Documentation enrichment metadata path")
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
	enrichment, err := readDocsEnrichment(*docsEnrichmentPath)
	if err != nil {
		return err
	}
	if err := validateDocsEnrichmentTargets(enrichment, objects); err != nil {
		return err
	}

	if err := validateGeneratedDocs(*docsDir, objects, relationships, enrichment); err != nil {
		return err
	}

	fmt.Println("Documentation validation passed.")
	return nil
}

func runDocsVerifyOnline(args []string) error {
	fs := flag.NewFlagSet("docs-verify-online", flag.ContinueOnError)
	managedPath := fs.String("managed", defaultManagedPath, "Managed object manifest path")
	docsEnrichmentPath := fs.String("docs-enrichment", defaultDocsEnrichmentPath, "Documentation enrichment metadata path")
	timeout := fs.Duration("timeout", defaultOnlineVerifyTimeout, "HTTP timeout per official AWX URL check")
	if err := fs.Parse(args); err != nil {
		return err
	}

	objects, err := readManagedObjects(*managedPath)
	if err != nil {
		return err
	}
	enrichment, err := readDocsEnrichment(*docsEnrichmentPath)
	if err != nil {
		return err
	}
	if err := validateDocsEnrichmentTargets(enrichment, objects); err != nil {
		return err
	}
	client := &http.Client{Timeout: *timeout}
	if err := verifyDocsEnrichmentOnline(client, objects, enrichment); err != nil {
		return err
	}

	fmt.Println("Online documentation verification passed.")
	return nil
}

func runDocsValidateQuality(args []string) error {
	fs := flag.NewFlagSet("docs-validate-quality", flag.ContinueOnError)
	summaryPath := fs.String("summary", "", "Path to implementation-summary.md quality analysis artifact")
	maxPasses := fs.Int("max-passes", defaultMaxQualityPasses, "Maximum allowed quality analysis passes")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if strings.TrimSpace(*summaryPath) == "" {
		return errors.New("docs-validate-quality requires --summary")
	}
	if *maxPasses < 1 {
		return errors.New("docs-validate-quality requires --max-passes >= 1")
	}
	if err := validateQualityAnalysisSummary(*summaryPath, *maxPasses); err != nil {
		return err
	}
	fmt.Println("Quality analysis summary validation passed.")
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

func generate(schemaPath string, exclusionsPath string, deprecatedPath string, prioritiesPath string, overridesPath string) (map[string]manifest.RuntimeExclusion, map[string]int, []manifest.ManagedObject, []manifest.Relationship, CoverageReport, int, error) {
	doc, err := openapi.LoadDocument(schemaPath)
	if err != nil {
		return nil, nil, nil, nil, CoverageReport{}, 0, err
	}

	runtimeExclusions, err := openapi.LoadRuntimeExclusions(exclusionsPath)
	if err != nil {
		return nil, nil, nil, nil, CoverageReport{}, 0, err
	}
	deprecatedObjects, deprecatedRelationships, err := openapi.LoadDeprecatedExclusions(deprecatedPath)
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

	objects := openapi.DeriveManagedObjects(doc, runtimeExclusions, deprecatedObjects)
	objects = openapi.ApplyFieldOverrides(objects, fieldOverrides)
	if err := validateTerraformFieldNameCollisions(objects); err != nil {
		return runtimeExclusions, priorities, objects, nil, buildReport(schemaPath, objects, nil, runtimeExclusions), len(fieldOverrides), err
	}
	if err := openapi.ValidateCoverage(objects, runtimeExclusions); err != nil {
		return runtimeExclusions, priorities, objects, nil, buildReport(schemaPath, objects, nil, runtimeExclusions), len(fieldOverrides), err
	}

	relationships := openapi.DeriveRelationships(doc, objects, priorities, deprecatedRelationships)
	report := buildReport(schemaPath, objects, relationships, runtimeExclusions)
	return runtimeExclusions, priorities, objects, relationships, report, len(fieldOverrides), nil
}

func validateTerraformFieldNameCollisions(objects []manifest.ManagedObject) error {
	for _, object := range objects {
		if object.RuntimeExcluded || (!object.ResourceEligible && !object.DataSourceElig) {
			continue
		}

		seen := make(map[string]string, len(object.Fields))
		for _, field := range object.Fields {
			tfName := manifest.TerraformAttributeNameForField(object.Name, field)
			if existing, ok := seen[tfName]; ok && existing != field.Name {
				return fmt.Errorf(
					"field naming collision: %s.%s and %s.%s map to Terraform attribute %q",
					object.Name, existing, object.Name, field.Name, tfName,
				)
			}
			seen[tfName] = field.Name
		}
	}
	return nil
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

var genericAWXIndexURLs = map[string]struct{}{
	"https://ansible.readthedocs.io/projects/awx/en/24.6.1/userguide/index.html": {},
	"https://docs.ansible.com/projects/awx/en/24.6.1/userguide/index.html":       {},
}

var awxOfficialDocsByObject = map[string]docsLink{
	"constructed_inventories": {
		Title: "AWX Inventories",
		URL:   "https://docs.ansible.com/projects/awx/en/24.6.1/userguide/inventories.html",
	},
	"credential_input_sources": {
		Title: "AWX Secret Management System",
		URL:   "https://docs.ansible.com/projects/awx/en/24.6.1/userguide/credential_plugins.html",
	},
	"credential_types": {
		Title: "AWX Credential Types",
		URL:   "https://docs.ansible.com/projects/awx/en/24.6.1/userguide/credential_types.html",
	},
	"credentials": {
		Title: "AWX Credentials",
		URL:   "https://docs.ansible.com/projects/awx/en/24.6.1/userguide/credentials.html",
	},
	"execution_environments": {
		Title: "AWX Execution Environments",
		URL:   "https://docs.ansible.com/projects/awx/en/24.6.1/userguide/execution_environments.html",
	},
	"groups": {
		Title: "AWX Inventories",
		URL:   "https://docs.ansible.com/projects/awx/en/24.6.1/userguide/inventories.html",
	},
	"hosts": {
		Title: "AWX Inventories",
		URL:   "https://docs.ansible.com/projects/awx/en/24.6.1/userguide/inventories.html",
	},
	"instance_groups": {
		Title: "AWX Instance Groups",
		URL:   "https://docs.ansible.com/projects/awx/en/24.6.1/userguide/instance_groups.html",
	},
	"instances": {
		Title: "AWX Instances",
		URL:   "https://docs.ansible.com/projects/awx/en/24.6.1/administration/instances.html",
	},
	"inventories": {
		Title: "AWX Inventories",
		URL:   "https://docs.ansible.com/projects/awx/en/24.6.1/userguide/inventories.html",
	},
	"inventory_sources": {
		Title: "AWX Inventory Sources",
		URL:   "https://docs.ansible.com/projects/awx/en/24.6.1/userguide/inventories.html",
	},
	"job_templates": {
		Title: "AWX Job Templates",
		URL:   "https://docs.ansible.com/projects/awx/en/24.6.1/userguide/job_templates.html",
	},
	"labels": {
		Title: "AWX Job Templates",
		URL:   "https://docs.ansible.com/projects/awx/en/24.6.1/userguide/job_templates.html",
	},
	"notification_templates": {
		Title: "AWX Notifications",
		URL:   "https://docs.ansible.com/projects/awx/en/24.6.1/userguide/notifications.html",
	},
	"organizations": {
		Title: "AWX Organizations",
		URL:   "https://docs.ansible.com/projects/awx/en/24.6.1/userguide/organizations.html",
	},
	"projects": {
		Title: "AWX Projects",
		URL:   "https://docs.ansible.com/projects/awx/en/24.6.1/userguide/projects.html",
	},
	"role_definitions": {
		Title: "AWX Role-Based Access Controls",
		URL:   "https://docs.ansible.com/projects/awx/en/24.6.1/userguide/rbac.html",
	},
	"role_team_assignments": {
		Title: "AWX Role-Based Access Controls",
		URL:   "https://docs.ansible.com/projects/awx/en/24.6.1/userguide/rbac.html",
	},
	"role_user_assignments": {
		Title: "AWX Role-Based Access Controls",
		URL:   "https://docs.ansible.com/projects/awx/en/24.6.1/userguide/rbac.html",
	},
	"schedules": {
		Title: "AWX Schedules",
		URL:   "https://docs.ansible.com/projects/awx/en/24.6.1/userguide/scheduling.html",
	},
	"settings": {
		Title: "AWX Configuration",
		URL:   "https://docs.ansible.com/projects/awx/en/24.6.1/administration/configure_awx.html",
	},
	"survey_spec": {
		Title: "AWX Job Template Surveys",
		URL:   "https://docs.ansible.com/projects/awx/en/24.6.1/userguide/job_templates.html",
	},
	"system_job_templates": {
		Title: "AWX Jobs",
		URL:   "https://docs.ansible.com/projects/awx/en/24.6.1/userguide/jobs.html",
	},
	"teams": {
		Title: "AWX Teams",
		URL:   "https://docs.ansible.com/projects/awx/en/24.6.1/userguide/teams.html",
	},
	"users": {
		Title: "AWX Users",
		URL:   "https://docs.ansible.com/projects/awx/en/24.6.1/userguide/users.html",
	},
	"workflow_job_nodes": {
		Title: "AWX Workflow Job Templates",
		URL:   "https://docs.ansible.com/projects/awx/en/24.6.1/userguide/workflow_templates.html",
	},
	"workflow_job_template_nodes": {
		Title: "AWX Workflow Job Templates",
		URL:   "https://docs.ansible.com/projects/awx/en/24.6.1/userguide/workflow_templates.html",
	},
	"workflow_job_templates": {
		Title: "AWX Workflow Job Templates",
		URL:   "https://docs.ansible.com/projects/awx/en/24.6.1/userguide/workflow_templates.html",
	},
}

func generateDocs(outputDir string, objects []manifest.ManagedObject, relationships []manifest.Relationship, enrichment docsEnrichmentCatalog) error {
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

	managedResourceBySingular := make(map[string]struct{})
	for _, obj := range objects {
		if !obj.ResourceEligible || obj.RuntimeExcluded {
			continue
		}
		managedResourceBySingular[obj.SingularName] = struct{}{}
	}

	for _, obj := range objects {
		objEnrichment := objectEnrichmentFor(obj, enrichment)
		if obj.RuntimeExcluded || (!obj.ResourceEligible && !obj.DataSourceElig) {
			continue
		}
		awxLinks, err := awxOfficialLinksForObject(obj.Name)
		if err != nil {
			return err
		}
		if obj.ResourceEligible && !obj.RuntimeExcluded {
			if err := writeResourceDoc(resourceDir, obj, objEnrichment, awxLinks, managedResourceBySingular); err != nil {
				return err
			}
		}
		if obj.DataSourceElig && !obj.RuntimeExcluded {
			if err := writeDataSourceDoc(dataSourceDir, obj, objEnrichment, awxLinks); err != nil {
				return err
			}
		}
	}
	for _, rel := range relationships {
		awxLinks, err := awxOfficialLinksForRelationship(rel)
		if err != nil {
			return err
		}
		if err := writeRelationshipDoc(resourceDir, rel, awxLinks); err != nil {
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
  hostname  = var.awx_hostname
  username  = var.awx_username
  password  = var.awx_password
}
` + "```" + `

## Schema

### Required

- ` + "`hostname`" + ` (String) AWX base URL, for example ` + "`https://awx.example.com`" + `.
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

## Breaking Changes

Reference fields that link one AWX object to another use an explicit ` + "`_id`" + ` suffix in Terraform.
If upgrading from older provider releases, rename unsuffixed link fields (for example, ` + "`organization`" + ` -> ` + "`organization_id`" + `) in resources and data sources.

## Compatibility

This provider targets AWX 24.6.1 API v2 only. Runtime-only objects are excluded from managed resources.
`

	return os.WriteFile(filepath.Join(outputDir, "index.md"), []byte(contents), 0o644)
}

func writeResourceDoc(resourceDir string, obj manifest.ManagedObject, objEnrichment objectDocsEnrichment, awxLinks []docsLink, managedResourceBySingular map[string]struct{}) error {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("# Resource: %s\n\n", obj.ResourceName))
	overview := strings.TrimSpace(objEnrichment.Overview)
	if overview == "" {
		overview = fmt.Sprintf("Manages AWX `%s` objects.", obj.Name)
	}
	builder.WriteString(overview + "\n\n")
	if isSettingsObject(obj) {
		builder.WriteString("Use `id = \"all\"` as the default and recommended settings scope.\n")
		builder.WriteString("Category-scoped IDs (for example `system`, `authentication`, and `bulk`) remain supported for optional advanced scoping.\n")
		builder.WriteString("Avoid overlapping ownership of the same setting key across `id = \"all\"` and category-scoped resources, because overlaps can cause configuration conflicts.\n\n")
	}
	if !obj.UpdateSupported {
		builder.WriteString("This endpoint does not support in-place updates; Terraform replaces the resource when arguments change.\n\n")
	}
	if objEnrichment.Complex && strings.TrimSpace(objEnrichment.ConceptPrimer) != "" {
		builder.WriteString("## AWX Concepts\n\n")
		builder.WriteString(strings.TrimSpace(objEnrichment.ConceptPrimer))
		builder.WriteString("\n\n")
	}

	builder.WriteString("## Example Usage\n\n")
	renderExamples(&builder, buildResourceExamples(obj, objEnrichment, managedResourceBySingular))

	builder.WriteString("## Schema\n\n")
	builder.WriteString("### Qualifiers\n\n")
	builder.WriteString("- `Required`: Must be set in configuration.\n")
	builder.WriteString("- `Optional`: May be omitted.\n")
	builder.WriteString("- `Computed`: AWX sets the value during create or refresh.\n")
	builder.WriteString("- `Sensitive`: Terraform redacts the value in normal CLI output.\n")
	builder.WriteString("- `Write-Only`: Sent to AWX during create/update and not read back.\n\n")

	builder.WriteString("### Required\n\n")
	requiredCount := 0
	if !obj.CollectionCreate {
		builder.WriteString("- `id` (String, Required) AWX detail-path identifier for this object.\n")
		requiredCount++
	}
	for _, field := range obj.Fields {
		if field.Required {
			tfName := manifest.TerraformAttributeNameForField(obj.Name, field)
			description := resolveFieldDescription(tfName, field, objEnrichment)
			builder.WriteString(fmt.Sprintf("- `%s` (%s, Required) %s\n", tfName, terraformTypeLabel(field), formatListItemDescription(description)))
			requiredCount++
		}
	}
	if requiredCount == 0 {
		builder.WriteString("- None.\n")
	}

	builder.WriteString("\n### Optional\n\n")
	optionalCount := 0
	for _, field := range obj.Fields {
		if field.Required {
			continue
		}
		tfName := manifest.TerraformAttributeNameForField(obj.Name, field)
		description := resolveFieldDescription(tfName, field, objEnrichment)
		qualifiers := []string{"Optional"}
		if field.Computed {
			qualifiers = append(qualifiers, "Computed")
		}
		if field.Sensitive {
			qualifiers = append(qualifiers, "Sensitive")
		}
		if field.WriteOnly {
			qualifiers = append(qualifiers, "Write-Only")
		}
		builder.WriteString(fmt.Sprintf("- `%s` (%s, %s) %s\n", tfName, terraformTypeLabel(field), strings.Join(qualifiers, ", "), formatListItemDescription(description)))
		optionalCount++
	}
	if optionalCount == 0 {
		builder.WriteString("- None.\n")
	}

	builder.WriteString("\n### Read-Only\n\n")
	if obj.CollectionCreate {
		builder.WriteString("- `id` (Number, Read-Only) Numeric AWX object identifier.\n")
	} else {
		builder.WriteString("- `id` (String, Read-Only) AWX detail-path identifier for this object.\n")
	}

	builder.WriteString("\n## Import\n\n")
	builder.WriteString("```bash\n")
	importID := "42"
	if !obj.CollectionCreate {
		importID = "example"
		if isSettingsObject(obj) {
			importID = "all"
		}
	}
	builder.WriteString(fmt.Sprintf("terraform import %s.example %s\n", obj.ResourceName, importID))
	builder.WriteString("```\n\n")

	builder.WriteString("## Further Reading\n\n")
	writeFurtherReading(&builder, furtherReadingLinks(awxLinks, objEnrichment.FurtherReading))

	return os.WriteFile(filepath.Join(resourceDir, fmt.Sprintf("%s.md", obj.ResourceName)), []byte(builder.String()), 0o644)
}

func writeDataSourceDoc(dataSourceDir string, obj manifest.ManagedObject, objEnrichment objectDocsEnrichment, awxLinks []docsLink) error {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("# Data Source: %s\n\n", obj.DataSourceName))
	builder.WriteString(fmt.Sprintf("Reads AWX `%s` objects.\n\n", obj.Name))
	if isSettingsObject(obj) {
		builder.WriteString("Use `id = \"all\"` as the default and recommended settings scope.\n")
		builder.WriteString("Category-scoped IDs (for example `system`, `authentication`, and `bulk`) remain supported for optional advanced scoping.\n")
		builder.WriteString("Avoid overlapping ownership of the same setting key across `id = \"all\"` and category-scoped resources, because overlaps can cause configuration conflicts.\n\n")
	}
	builder.WriteString("## Example Usage\n\n")
	renderExamples(&builder, []docsExample{dataSourceExample(obj)})

	builder.WriteString("## Schema\n\n")
	builder.WriteString("### Optional\n\n")
	if obj.CollectionCreate {
		builder.WriteString("- `id` (Number, Optional) Numeric AWX object ID.\n")
	} else {
		builder.WriteString("- `id` (String, Optional) AWX object identifier used in the detail endpoint path.\n")
	}
	if hasField(obj.Fields, "name") {
		builder.WriteString("- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.\n")
	}
	builder.WriteString("\n### Read-Only\n\n")
	if obj.CollectionCreate {
		builder.WriteString("- `id` (Number, Read-Only) Numeric AWX object ID.\n")
	} else {
		builder.WriteString("- `id` (String, Read-Only) AWX detail-path identifier for this object.\n")
	}
	for _, field := range obj.Fields {
		tfName := manifest.TerraformAttributeNameForField(obj.Name, field)
		description := resolveFieldDescription(tfName, field, objEnrichment)
		qualifiers := []string{"Read-Only"}
		if field.Sensitive {
			qualifiers = append(qualifiers, "Sensitive")
		}
		builder.WriteString(fmt.Sprintf("- `%s` (%s, %s) %s\n", tfName, terraformTypeLabel(field), strings.Join(qualifiers, ", "), formatListItemDescription(description)))
	}
	builder.WriteString("\n## Further Reading\n\n")
	writeFurtherReading(&builder, furtherReadingLinks(awxLinks, objEnrichment.FurtherReading))

	return os.WriteFile(filepath.Join(dataSourceDir, fmt.Sprintf("%s.md", obj.DataSourceName)), []byte(builder.String()), 0o644)
}

func writeRelationshipDoc(resourceDir string, rel manifest.Relationship, awxLinks []docsLink) error {
	builder := strings.Builder{}
	parentIDAttribute := manifest.RelationshipParentIDAttribute(rel)
	childIDAttribute := manifest.RelationshipChildIDAttribute(rel)

	builder.WriteString(fmt.Sprintf("# Resource: %s\n\n", rel.ResourceName))
	if isSurveySpecRelationship(rel) {
		builder.WriteString(fmt.Sprintf("Manages the AWX survey specification for `%s` objects.\n\n", rel.ParentObject))
		builder.WriteString("## Example Usage\n\n")
		renderExamples(&builder, []docsExample{{
			HCL: fmt.Sprintf("resource %q %q {\n  %s = 12\n  spec = jsonencode({\n    name        = \"Example survey\"\n    description = \"Managed by Terraform\"\n    spec        = []\n  })\n}", rel.ResourceName, "example", parentIDAttribute),
		}})
		builder.WriteString("## Schema\n\n")
		builder.WriteString("### Required\n\n")
		builder.WriteString(fmt.Sprintf("- `%s` (Number, Required) Parent object numeric ID.\n", parentIDAttribute))
		builder.WriteString("\n### Optional\n\n")
		builder.WriteString("- `spec` (String, Optional) JSON-encoded survey specification payload.\n")
		builder.WriteString("\n### Read-Only\n\n")
		builder.WriteString(fmt.Sprintf("- `id` (String, Read-Only) Survey specification ID (same as `%s`).\n", parentIDAttribute))
		builder.WriteString(fmt.Sprintf("- `%s` (Number, Read-Only) Parent object numeric ID.\n", parentIDAttribute))
		builder.WriteString("## Import\n\n")
		builder.WriteString("```bash\n")
		builder.WriteString(fmt.Sprintf("terraform import %s.example <resource_id>\n", rel.ResourceName))
		builder.WriteString("```\n\n")
		builder.WriteString("## Further Reading\n\n")
		writeFurtherReading(&builder, furtherReadingLinks(awxLinks))
	} else {
		builder.WriteString(fmt.Sprintf("Manages AWX associations between `%s` and `%s` objects.\n\n", rel.ParentObject, rel.ChildObject))
		builder.WriteString("## Example Usage\n\n")
		renderExamples(&builder, []docsExample{{
			HCL: fmt.Sprintf("resource %q %q {\n  %s = 12\n  %s = 34\n}", rel.ResourceName, "example", parentIDAttribute, childIDAttribute),
		}})
		builder.WriteString("## Schema\n\n")
		builder.WriteString("### Required\n\n")
		builder.WriteString(fmt.Sprintf("- `%s` (Number, Required) Parent object numeric ID.\n", parentIDAttribute))
		builder.WriteString(fmt.Sprintf("- `%s` (Number, Required) Child object numeric ID.\n", childIDAttribute))
		builder.WriteString("\n### Read-Only\n\n")
		builder.WriteString("- `id` (String, Read-Only) Composite ID in `<primary_id>:<related_id>` format.\n")
		builder.WriteString(fmt.Sprintf("- `%s` (Number, Read-Only) Parent object numeric ID.\n", parentIDAttribute))
		builder.WriteString(fmt.Sprintf("- `%s` (Number, Read-Only) Child object numeric ID.\n", childIDAttribute))
		builder.WriteString("## Import\n\n")
		builder.WriteString("```bash\n")
		builder.WriteString(fmt.Sprintf("terraform import %s.example <primary_id>:<related_id>\n", rel.ResourceName))
		builder.WriteString("```\n\n")
		builder.WriteString("## Further Reading\n\n")
		writeFurtherReading(&builder, furtherReadingLinks(awxLinks))
	}

	return os.WriteFile(filepath.Join(resourceDir, fmt.Sprintf("%s.md", rel.ResourceName)), []byte(builder.String()), 0o644)
}

func objectEnrichmentFor(obj manifest.ManagedObject, enrichment docsEnrichmentCatalog) objectDocsEnrichment {
	if out, ok := enrichment.Objects[obj.ResourceName]; ok {
		return out
	}
	if out, ok := enrichment.Objects[obj.Name]; ok {
		return out
	}
	if out, ok := enrichment.Objects[obj.DataSourceName]; ok {
		return out
	}
	return objectDocsEnrichment{}
}

func buildResourceExamples(obj manifest.ManagedObject, objEnrichment objectDocsEnrichment, managedResourceBySingular map[string]struct{}) []docsExample {
	out := make([]docsExample, 0, 3)
	if objEnrichment.PrimaryExample != nil {
		out = append(out, *objEnrichment.PrimaryExample)
	} else {
		out = append(out, defaultResourceExample(obj, managedResourceBySingular))
	}
	out = append(out, objEnrichment.ExtraExamples...)
	if len(out) > 3 {
		return out[:3]
	}
	return out
}

func defaultResourceExample(obj manifest.ManagedObject, managedResourceBySingular map[string]struct{}) docsExample {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("resource %q %q {\n", obj.ResourceName, "example"))
	if !obj.CollectionCreate {
		idExample := "\"example\""
		if isSettingsObject(obj) {
			idExample = "\"all\""
		}
		builder.WriteString("  id = " + idExample + "\n")
	}
	exampleFields := map[string]struct{}{}
	for _, field := range obj.Fields {
		if !field.Required {
			continue
		}
		tfName := manifest.TerraformAttributeNameForField(obj.Name, field)
		builder.WriteString(fmt.Sprintf("  %s = %s\n", tfName, sampleDocValue(field, tfName, managedResourceBySingular)))
		exampleFields[field.Name] = struct{}{}
	}
	for _, field := range obj.Fields {
		if field.Type != manifest.FieldTypeObject || field.WriteOnly {
			continue
		}
		if _, alreadyIncluded := exampleFields[field.Name]; alreadyIncluded {
			continue
		}
		tfName := manifest.TerraformAttributeNameForField(obj.Name, field)
		builder.WriteString(fmt.Sprintf("  %s = %s\n", tfName, sampleValue(field.Type)))
		break
	}
	builder.WriteString("}")
	return docsExample{
		Title: "Basic configuration",
		HCL:   builder.String(),
	}
}

func dataSourceExample(obj manifest.ManagedObject) docsExample {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("data %q %q {\n", obj.DataSourceName, "example"))
	idExample := sampleValue(manifest.FieldTypeInt)
	if !obj.CollectionCreate {
		idExample = sampleValue(manifest.FieldTypeString)
		if isSettingsObject(obj) {
			idExample = "\"all\""
		}
	}
	builder.WriteString("  id = " + idExample + "\n")
	builder.WriteString("}")
	return docsExample{HCL: builder.String()}
}

func isSettingsObject(obj manifest.ManagedObject) bool {
	return strings.TrimSpace(obj.Name) == "settings"
}

func renderExamples(builder *strings.Builder, examples []docsExample) {
	for idx, example := range examples {
		title := strings.TrimSpace(example.Title)
		if len(examples) > 1 || title != "" {
			if title == "" {
				title = fmt.Sprintf("Example %d", idx+1)
			}
			builder.WriteString("### " + title + "\n\n")
		}
		if description := strings.TrimSpace(example.Description); description != "" {
			builder.WriteString(description + "\n\n")
		}
		builder.WriteString("```hcl\n")
		builder.WriteString(strings.TrimSpace(example.HCL))
		builder.WriteString("\n```\n\n")
	}
}

func terraformTypeLabel(field manifest.FieldSpec) string {
	switch field.Type {
	case manifest.FieldTypeBool:
		return "Boolean"
	case manifest.FieldTypeInt, manifest.FieldTypeFloat:
		return "Number"
	case manifest.FieldTypeObject:
		return "Object"
	case manifest.FieldTypeArray:
		return "String"
	default:
		return "String"
	}
}

func resolveFieldDescription(terraformName string, field manifest.FieldSpec, objEnrichment objectDocsEnrichment) string {
	if override := lookupFieldDescription(terraformName, field.Name, objEnrichment.FieldDescriptions); override != "" {
		return sanitizeDescription(override)
	}

	description := strings.TrimSpace(field.Description)
	if description != "" && !strings.Contains(description, openAPIPlaceholderText) && !isLowInformationDescription(description) {
		return sanitizeDescription(description)
	}

	if field.Reference || strings.HasSuffix(terraformName, "_id") {
		target := strings.TrimSuffix(terraformName, "_id")
		target = strings.ReplaceAll(target, "_", " ")
		if target == "" {
			target = terraformName
		}
		return sanitizeDescription(fmt.Sprintf("Numeric ID of the related AWX %s object.", target))
	}

	if field.WriteOnly && field.Sensitive {
		return sanitizeDescription(fmt.Sprintf("Write-only sensitive object for `%s`.", terraformName))
	}

	switch field.Type {
	case manifest.FieldTypeBool:
		return sanitizeDescription(fmt.Sprintf("Controls whether `%s` is enabled in AWX.", terraformName))
	case manifest.FieldTypeInt, manifest.FieldTypeFloat:
		return sanitizeDescription(fmt.Sprintf("Numeric AWX value used for `%s`.", terraformName))
	case manifest.FieldTypeArray:
		return sanitizeDescription(fmt.Sprintf("JSON-encoded list value for `%s`.", terraformName))
	case manifest.FieldTypeObject:
		return sanitizeDescription(fmt.Sprintf("Object value for `%s`.", terraformName))
	default:
		return sanitizeDescription(fmt.Sprintf("AWX value stored in `%s`.", terraformName))
	}
}

func isLowInformationDescription(description string) bool {
	content := strings.ToLower(strings.TrimSpace(description))
	if strings.Contains(content, "managed field from awx openapi schema") {
		return true
	}
	if strings.HasPrefix(content, "value for ") {
		return true
	}
	if strings.HasPrefix(content, "numeric setting for ") {
		return true
	}
	return false
}

func sanitizeDescription(description string) string {
	sanitized := strings.ReplaceAll(description, "<", "&lt;")
	sanitized = strings.ReplaceAll(sanitized, ">", "&gt;")
	return sanitized
}

func lookupFieldDescription(terraformName string, fieldName string, overrides map[string]string) string {
	if len(overrides) == 0 {
		return ""
	}
	if out := strings.TrimSpace(overrides[terraformName]); out != "" {
		return out
	}
	return strings.TrimSpace(overrides[fieldName])
}

func awxOfficialLinksForObject(objectName string) ([]docsLink, error) {
	link, ok := awxOfficialDocsByObject[objectName]
	if !ok {
		return nil, fmt.Errorf("missing official AWX documentation mapping for object %q", objectName)
	}
	return []docsLink{link}, nil
}

func awxOfficialLinksForRelationship(rel manifest.Relationship) ([]docsLink, error) {
	parentLinks, err := awxOfficialLinksForObject(rel.ParentObject)
	if err != nil {
		return nil, err
	}
	if isSurveySpecRelationship(rel) {
		return parentLinks, nil
	}

	childLinks, err := awxOfficialLinksForObject(rel.ChildObject)
	if err != nil {
		return nil, err
	}
	return append(parentLinks, childLinks...), nil
}

func furtherReadingLinks(awxLinks []docsLink, curated ...[]docsLink) []docsLink {
	curatedCount := 0
	for _, links := range curated {
		curatedCount += len(links)
	}
	out := make([]docsLink, 0, len(awxLinks)+curatedCount)
	seen := make(map[string]struct{})
	appendLink := func(link docsLink) {
		title := strings.TrimSpace(link.Title)
		linkURL := strings.TrimSpace(link.URL)
		if title == "" || linkURL == "" {
			return
		}
		if _, ok := seen[linkURL]; ok {
			return
		}
		seen[linkURL] = struct{}{}
		out = append(out, docsLink{Title: title, URL: linkURL})
	}

	for _, link := range awxLinks {
		appendLink(link)
	}
	for _, curatedSet := range curated {
		for _, link := range curatedSet {
			appendLink(link)
		}
	}
	return out
}

func writeFurtherReading(builder *strings.Builder, links []docsLink) {
	for _, link := range links {
		builder.WriteString(fmt.Sprintf("- [%s](%s)\n", link.Title, link.URL))
	}
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
		return "{ key = \"value\" }"
	default:
		return "\"example\""
	}
}

func sampleDocValue(field manifest.FieldSpec, terraformName string, managedResourceBySingular map[string]struct{}) string {
	if !field.Reference || field.Type != manifest.FieldTypeInt {
		return sampleValue(field.Type)
	}
	if !strings.HasSuffix(terraformName, "_id") {
		return sampleValue(field.Type)
	}
	target := strings.TrimSuffix(terraformName, "_id")
	if strings.TrimSpace(target) == "" {
		return sampleValue(field.Type)
	}
	if _, ok := managedResourceBySingular[target]; !ok {
		return sampleValue(field.Type)
	}
	return fmt.Sprintf("awx_%s.example.id", target)
}

func hasField(fields []manifest.FieldSpec, name string) bool {
	for _, field := range fields {
		if field.Name == name {
			return true
		}
	}
	return false
}

func formatListItemDescription(description string) string {
	trimmed := strings.TrimSpace(description)
	if trimmed == "" {
		return ""
	}

	normalized := strings.ReplaceAll(trimmed, "\\n*", "\n*")
	normalized = strings.ReplaceAll(normalized, "\\n- ", "\n- ")
	normalized = inlineEnumBulletPattern.ReplaceAllString(normalized, "\n* ")

	lines := strings.Split(normalized, "\n")
	first := ""
	bullets := make([]string, 0, len(lines))
	notes := make([]string, 0, len(lines))

	for _, line := range lines {
		clean := strings.TrimSpace(line)
		if clean == "" {
			continue
		}
		if strings.HasPrefix(clean, "* ") {
			bullets = append(bullets, strings.TrimSpace(strings.TrimPrefix(clean, "* ")))
			continue
		}
		if strings.HasPrefix(clean, "- ") {
			bullets = append(bullets, strings.TrimSpace(strings.TrimPrefix(clean, "- ")))
			continue
		}
		if first == "" {
			first = clean
			continue
		}
		notes = append(notes, clean)
	}

	if first == "" {
		first = "Allowed values:"
	}

	if len(bullets) == 0 && len(notes) == 0 {
		return first
	}

	builder := strings.Builder{}
	builder.WriteString(first)
	for _, note := range notes {
		builder.WriteString("\n  ")
		builder.WriteString(note)
	}
	for _, bullet := range bullets {
		builder.WriteString("\n  - ")
		builder.WriteString(bullet)
	}
	return builder.String()
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

func readDocsEnrichment(path string) (docsEnrichmentCatalog, error) {
	raw, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return docsEnrichmentCatalog{}, fmt.Errorf("read docs enrichment metadata: %w", err)
	}
	var out docsEnrichmentCatalog
	if err := json.Unmarshal(raw, &out); err != nil {
		return docsEnrichmentCatalog{}, fmt.Errorf("parse docs enrichment JSON: %w", err)
	}
	if out.Objects == nil {
		out.Objects = map[string]objectDocsEnrichment{}
	}
	if err := validateDocsEnrichmentSchema(out); err != nil {
		return docsEnrichmentCatalog{}, err
	}
	return out, nil
}

func validateDocsEnrichmentSchema(enrichment docsEnrichmentCatalog) error {
	seenPriority := make(map[string]struct{}, len(enrichment.PriorityResources))
	for _, resourceName := range enrichment.PriorityResources {
		key := strings.TrimSpace(resourceName)
		if key == "" {
			return fmt.Errorf("docs enrichment priorityResources contains an empty resource name")
		}
		if _, exists := seenPriority[key]; exists {
			return fmt.Errorf("docs enrichment priorityResources contains duplicate %q", key)
		}
		seenPriority[key] = struct{}{}
	}

	validateExample := func(label string, example docsExample) error {
		if strings.TrimSpace(example.HCL) == "" {
			return fmt.Errorf("docs enrichment %s has empty hcl content", label)
		}
		return nil
	}

	for objectKey, metadata := range enrichment.Objects {
		key := strings.TrimSpace(objectKey)
		if key == "" {
			return fmt.Errorf("docs enrichment contains an empty objects key")
		}
		if metadata.Complex && strings.TrimSpace(metadata.ConceptPrimer) == "" {
			return fmt.Errorf("docs enrichment objects.%s is marked complex but conceptPrimer is empty", key)
		}
		for fieldName, description := range metadata.FieldDescriptions {
			if strings.TrimSpace(fieldName) == "" {
				return fmt.Errorf("docs enrichment objects.%s has an empty fieldDescriptions key", key)
			}
			if strings.TrimSpace(description) == "" {
				return fmt.Errorf("docs enrichment objects.%s.fieldDescriptions.%s is empty", key, fieldName)
			}
		}
		if metadata.OnlineResearchChecklist != nil {
			if strings.TrimSpace(metadata.OnlineResearchChecklist.ObjectBehavior) == "" {
				return fmt.Errorf("docs enrichment objects.%s.onlineResearchChecklist.objectBehavior is required", key)
			}
			if strings.TrimSpace(metadata.OnlineResearchChecklist.RelatedInteractions) == "" {
				return fmt.Errorf("docs enrichment objects.%s.onlineResearchChecklist.relatedInteractions is required", key)
			}
			if strings.TrimSpace(metadata.OnlineResearchChecklist.ParameterSemantics) == "" {
				return fmt.Errorf("docs enrichment objects.%s.onlineResearchChecklist.parameterSemantics is required", key)
			}
		}
		for idx, fieldName := range metadata.InteractionReferenceFields {
			if strings.TrimSpace(fieldName) == "" {
				return fmt.Errorf("docs enrichment objects.%s.interactionReferenceFields[%d] is empty", key, idx)
			}
		}
		if len(metadata.ExtraExamples) > 2 {
			return fmt.Errorf("docs enrichment objects.%s defines %d extraExamples; maximum is 2", key, len(metadata.ExtraExamples))
		}
		if metadata.PrimaryExample != nil {
			if err := validateExample("objects."+key+".primaryExample", *metadata.PrimaryExample); err != nil {
				return err
			}
		}
		for idx, example := range metadata.ExtraExamples {
			if err := validateExample(fmt.Sprintf("objects.%s.extraExamples[%d]", key, idx), example); err != nil {
				return err
			}
		}
		for idx, link := range metadata.FurtherReading {
			title := strings.TrimSpace(link.Title)
			linkURL := strings.TrimSpace(link.URL)
			if title == "" || linkURL == "" {
				return fmt.Errorf("docs enrichment objects.%s.furtherReading[%d] requires title and url", key, idx)
			}
			parsed, err := url.Parse(linkURL)
			if err != nil || parsed.Scheme == "" || parsed.Host == "" {
				return fmt.Errorf("docs enrichment objects.%s.furtherReading[%d] has invalid url %q", key, idx, linkURL)
			}
		}
		if metadata.CurationSource != nil {
			urlValue := strings.TrimSpace(metadata.CurationSource.OfficialAWXURL)
			verifiedOn := strings.TrimSpace(metadata.CurationSource.VerifiedOn)
			if urlValue == "" || verifiedOn == "" {
				return fmt.Errorf("docs enrichment objects.%s.curationSource requires officialAwxUrl and verifiedOn", key)
			}
			parsed, err := url.Parse(urlValue)
			if err != nil || parsed.Scheme == "" || parsed.Host == "" || !isOfficialAWXLink(urlValue) {
				return fmt.Errorf("docs enrichment objects.%s.curationSource.officialAwxUrl must be an official AWX documentation url", key)
			}
			parsedDate, err := time.Parse("2006-01-02", verifiedOn)
			if err != nil {
				return fmt.Errorf("docs enrichment objects.%s.curationSource.verifiedOn must use YYYY-MM-DD", key)
			}
			if parsedDate.After(time.Now().UTC()) {
				return fmt.Errorf("docs enrichment objects.%s.curationSource.verifiedOn cannot be in the future", key)
			}
		}
	}
	return nil
}

func validateDocsEnrichmentTargets(enrichment docsEnrichmentCatalog, objects []manifest.ManagedObject) error {
	validTargets := make(map[string]manifest.ManagedObject, len(objects)*3)
	validResources := make(map[string]manifest.ManagedObject, len(objects))
	requiredObjects := make([]manifest.ManagedObject, 0, len(objects))
	for _, object := range objects {
		validTargets[object.Name] = object
		validTargets[object.ResourceName] = object
		validTargets[object.DataSourceName] = object
		if object.ResourceEligible && !object.RuntimeExcluded {
			validResources[object.ResourceName] = object
		}
		if !object.RuntimeExcluded && (object.ResourceEligible || object.DataSourceElig) {
			requiredObjects = append(requiredObjects, object)
		}
	}

	for key := range enrichment.Objects {
		if _, ok := validTargets[key]; !ok {
			return fmt.Errorf("docs enrichment objects.%s does not match any managed object/resource/data-source name", key)
		}
	}

	for _, object := range requiredObjects {
		metadata, metadataKey, ok := objectEnrichmentMetadataForTarget(object, enrichment)
		if !ok {
			return fmt.Errorf("docs enrichment missing objects.%s entry for managed object documentation", object.ResourceName)
		}
		if metadata.CurationSource == nil {
			return fmt.Errorf("docs enrichment objects.%s requires curationSource for managed object documentation", metadataKey)
		}
		if metadata.OnlineResearchChecklist == nil {
			return fmt.Errorf("docs enrichment objects.%s requires onlineResearchChecklist for managed object documentation", metadataKey)
		}
		expectedAWXLinks, err := awxOfficialLinksForObject(object.Name)
		if err != nil {
			return err
		}
		if !matchesExpectedAWXConceptURL(metadata.CurationSource.OfficialAWXURL, expectedAWXLinks) {
			return fmt.Errorf("docs enrichment objects.%s.curationSource.officialAwxUrl must reference the mapped official AWX concept link", metadataKey)
		}
	}

	for _, resourceName := range enrichment.PriorityResources {
		object, ok := validResources[resourceName]
		if !ok {
			return fmt.Errorf("docs enrichment priority resource %q is not a managed resource", resourceName)
		}
		metadata, metadataKey, ok := objectEnrichmentMetadataForTarget(object, enrichment)
		if !ok {
			return fmt.Errorf("docs enrichment missing objects.%s entry for prioritized resource", resourceName)
		}
		if metadata.PrimaryExample == nil {
			return fmt.Errorf("docs enrichment objects.%s requires a primaryExample for prioritized resource documentation", metadataKey)
		}
	}

	return nil
}

func objectEnrichmentMetadataForTarget(object manifest.ManagedObject, enrichment docsEnrichmentCatalog) (objectDocsEnrichment, string, bool) {
	keys := []string{object.ResourceName, object.Name, object.DataSourceName}
	for _, key := range keys {
		if metadata, ok := enrichment.Objects[key]; ok {
			return metadata, key, true
		}
	}
	return objectDocsEnrichment{}, "", false
}

func validateGeneratedDocs(docsDir string, objects []manifest.ManagedObject, relationships []manifest.Relationship, enrichment docsEnrichmentCatalog) error {
	providerDoc := filepath.Join(docsDir, "index.md")
	if _, err := os.Stat(providerDoc); err != nil {
		return fmt.Errorf("provider documentation is missing: %s", providerDoc)
	}

	prioritySet := make(map[string]struct{}, len(enrichment.PriorityResources))
	for _, resourceName := range enrichment.PriorityResources {
		prioritySet[resourceName] = struct{}{}
	}
	seenPriorities := make(map[string]struct{}, len(prioritySet))

	for _, object := range objects {
		if object.RuntimeExcluded || (!object.ResourceEligible && !object.DataSourceElig) {
			continue
		}
		expectedObjectAWXLinks, err := awxOfficialLinksForObject(object.Name)
		if err != nil {
			return err
		}
		if object.ResourceEligible && !object.RuntimeExcluded {
			resourceDoc := filepath.Join(docsDir, "resources", fmt.Sprintf("%s.md", object.ResourceName))
			content, err := requireDocSectionsInOrder(resourceDoc, "## Example Usage", "## Schema", "## Import", "## Further Reading")
			if err != nil {
				return err
			}
			if err := ensureNoPlaceholderText(resourceDoc, content); err != nil {
				return err
			}
			if err := ensureNoLowInformationText(resourceDoc, content); err != nil {
				return err
			}
			if err := validateEnumFormatting(resourceDoc, content); err != nil {
				return err
			}
			if err := validateExampleBounds(resourceDoc, content, 1, 3); err != nil {
				return err
			}
			if err := validateFurtherReadingPolicy(resourceDoc, content, expectedObjectAWXLinks); err != nil {
				return err
			}
			if err := validateQualifierPlacement(resourceDoc, content); err != nil {
				return err
			}
			if metadata := objectEnrichmentFor(object, enrichment); metadata.Complex {
				if err := validateComplexPrimer(resourceDoc, content); err != nil {
					return err
				}
			}
			if isSettingsObject(object) {
				if err := validateSettingsResourceDocumentation(resourceDoc, content); err != nil {
					return err
				}
			}
			metadata := objectEnrichmentFor(object, enrichment)
			if _, isPriority := prioritySet[object.ResourceName]; isPriority {
				seenPriorities[object.ResourceName] = struct{}{}
				if err := validateResolvedExampleReferences(resourceDoc, content); err != nil {
					return err
				}
				if err := validateInteractionReferenceFields(resourceDoc, content, metadata); err != nil {
					return err
				}
			}
		}

		if object.DataSourceElig && !object.RuntimeExcluded {
			dataSourceDoc := filepath.Join(docsDir, "data-sources", fmt.Sprintf("%s.md", object.DataSourceName))
			content, err := requireDocSectionsInOrder(dataSourceDoc, "## Example Usage", "## Schema", "## Further Reading")
			if err != nil {
				return err
			}
			if err := ensureNoPlaceholderText(dataSourceDoc, content); err != nil {
				return err
			}
			if err := ensureNoLowInformationText(dataSourceDoc, content); err != nil {
				return err
			}
			if err := validateEnumFormatting(dataSourceDoc, content); err != nil {
				return err
			}
			if err := validateExampleBounds(dataSourceDoc, content, 1, 3); err != nil {
				return err
			}
			if err := validateFurtherReadingPolicy(dataSourceDoc, content, expectedObjectAWXLinks); err != nil {
				return err
			}
			if isSettingsObject(object) {
				if err := validateSettingsDataSourceDocumentation(dataSourceDoc, content); err != nil {
					return err
				}
			}
		}
	}

	for _, relationship := range relationships {
		expectedRelationshipAWXLinks, err := awxOfficialLinksForRelationship(relationship)
		if err != nil {
			return err
		}
		resourceDoc := filepath.Join(docsDir, "resources", fmt.Sprintf("%s.md", relationship.ResourceName))
		content, err := requireDocSectionsInOrder(resourceDoc, "## Example Usage", "## Schema", "## Import", "## Further Reading")
		if err != nil {
			return err
		}
		if err := validateExampleBounds(resourceDoc, content, 1, 3); err != nil {
			return err
		}
		if err := validateFurtherReadingPolicy(resourceDoc, content, expectedRelationshipAWXLinks); err != nil {
			return err
		}
		if isSurveySpecRelationship(relationship) {
			if !strings.Contains(content, "<resource_id>") {
				return fmt.Errorf("relationship documentation %s must render survey-spec import ID as <resource_id>", resourceDoc)
			}
			continue
		}
		if !strings.Contains(content, "<primary_id>:<related_id>") {
			return fmt.Errorf("relationship documentation %s must render composite import ID as <primary_id>:<related_id>", resourceDoc)
		}
	}

	for resourceName := range prioritySet {
		if _, found := seenPriorities[resourceName]; !found {
			return fmt.Errorf("prioritized resource documentation %q was not validated", resourceName)
		}
	}

	return nil
}

func validateSettingsResourceDocumentation(path string, content string) error {
	if err := validateSettingsDocumentationCommon(path, content); err != nil {
		return err
	}
	importSection, err := extractTopLevelSection(content, "## Import")
	if err != nil {
		return fmt.Errorf("documentation file %s: %w", path, err)
	}
	if !strings.Contains(importSection, "terraform import awx_setting.example all") {
		return fmt.Errorf("documentation file %s must use `all` as the primary awx_setting import example", path)
	}
	return nil
}

func validateSettingsDataSourceDocumentation(path string, content string) error {
	return validateSettingsDocumentationCommon(path, content)
}

func validateSettingsDocumentationCommon(path string, content string) error {
	exampleSection, err := extractTopLevelSection(content, "## Example Usage")
	if err != nil {
		return fmt.Errorf("documentation file %s: %w", path, err)
	}
	if !strings.Contains(exampleSection, "id = \"all\"") {
		return fmt.Errorf("documentation file %s must use `id = \"all\"` in the canonical awx_setting example", path)
	}
	requiredMarkers := []string{
		"category-scoped ids",
		"optional advanced scoping",
		"overlapping ownership",
		"configuration conflicts",
	}
	contentLower := strings.ToLower(content)
	for _, marker := range requiredMarkers {
		if !strings.Contains(contentLower, marker) {
			return fmt.Errorf("documentation file %s is missing required awx_setting guidance marker %q", path, marker)
		}
	}
	return nil
}

func requireDocSectionsInOrder(path string, sections ...string) (string, error) {
	raw, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return "", fmt.Errorf("documentation file missing: %s", path)
	}
	content := string(raw)
	lastIndex := -1
	for _, section := range sections {
		currentIndex := strings.Index(content, section)
		if currentIndex == -1 {
			return "", fmt.Errorf("documentation file %s is missing section %q", path, section)
		}
		if currentIndex < lastIndex {
			return "", fmt.Errorf("documentation file %s has section %q out of order", path, section)
		}
		lastIndex = currentIndex
	}
	return content, nil
}

func ensureNoPlaceholderText(path string, content string) error {
	if strings.Contains(content, openAPIPlaceholderText) {
		return fmt.Errorf("documentation file %s still contains placeholder text %q", path, openAPIPlaceholderText)
	}
	return nil
}

func ensureNoLowInformationText(path string, content string) error {
	lowInformationPatterns := []string{
		"managed field from awx openapi schema",
		") value for `",
		") numeric setting for `",
	}
	contentLower := strings.ToLower(content)
	for _, pattern := range lowInformationPatterns {
		if strings.Contains(contentLower, pattern) {
			return fmt.Errorf("documentation file %s contains low-information placeholder pattern %q", path, pattern)
		}
	}
	return nil
}

func validateEnumFormatting(path string, content string) error {
	if strings.Contains(content, "\\n*") || strings.Contains(content, "\\n-") {
		return fmt.Errorf("documentation file %s contains escaped newline enum formatting artifacts", path)
	}
	if malformedEnumInlinePattern.MatchString(content) {
		return fmt.Errorf("documentation file %s contains inline collapsed enum bullets; expected multiline markdown list formatting", path)
	}
	return nil
}

func validateResolvedExampleReferences(path string, content string) error {
	exampleSection, err := extractTopLevelSection(content, "## Example Usage")
	if err != nil {
		return fmt.Errorf("documentation file %s: %w", path, err)
	}

	declaredBlocks := make(map[string]struct{})
	for _, match := range resourceBlockPattern.FindAllStringSubmatch(exampleSection, -1) {
		key := "resource|" + match[1] + "|" + match[2]
		declaredBlocks[key] = struct{}{}
	}
	for _, match := range dataBlockPattern.FindAllStringSubmatch(exampleSection, -1) {
		key := "data|" + match[1] + "|" + match[2]
		declaredBlocks[key] = struct{}{}
	}

	for _, match := range unresolvedReferencePattern.FindAllStringSubmatch(exampleSection, -1) {
		prefix := "resource"
		if strings.TrimSpace(match[1]) != "" {
			prefix = "data"
		}
		resourceType := strings.TrimSpace(match[2])
		name := strings.TrimSpace(match[3])
		key := prefix + "|" + resourceType + "|" + name
		if _, ok := declaredBlocks[key]; !ok {
			return fmt.Errorf("documentation file %s has unresolved cross-resource reference %s.%s.id in examples", path, resourceType, name)
		}
	}
	return nil
}

func validateInteractionReferenceFields(path string, content string, metadata objectDocsEnrichment) error {
	if len(metadata.InteractionReferenceFields) == 0 {
		return nil
	}
	exampleSection, err := extractTopLevelSection(content, "## Example Usage")
	if err != nil {
		return fmt.Errorf("documentation file %s: %w", path, err)
	}
	found := make(map[string]struct{}, len(metadata.InteractionReferenceFields))
	for _, match := range curatedReferenceAssignmentPattern.FindAllStringSubmatch(exampleSection, -1) {
		found[strings.TrimSpace(match[1])] = struct{}{}
	}
	for _, field := range metadata.InteractionReferenceFields {
		key := strings.TrimSpace(field)
		if key == "" {
			continue
		}
		if _, ok := found[key]; !ok {
			return fmt.Errorf("documentation file %s must show `%s` wired to related object IDs in examples", path, key)
		}
	}
	return nil
}

func validateQualifierPlacement(path string, content string) error {
	schemaSection, err := extractTopLevelSection(content, "## Schema")
	if err != nil {
		return fmt.Errorf("documentation file %s: %w", path, err)
	}
	if !strings.Contains(schemaSection, "### Qualifiers") {
		return fmt.Errorf("documentation file %s is missing Schema/Qualifiers guidance", path)
	}
	if strings.Contains(content, "Argument qualifiers used below") {
		return fmt.Errorf("documentation file %s still includes legacy inline qualifier phrasing", path)
	}

	requiredSection, err := extractSubsection(schemaSection, "### Required")
	if err != nil {
		return fmt.Errorf("documentation file %s: %w", path, err)
	}
	optionalSection, err := extractSubsection(schemaSection, "### Optional")
	if err != nil {
		return fmt.Errorf("documentation file %s: %w", path, err)
	}
	for _, marker := range []string{"- `Required`:", "- `Optional`:", "- `Computed`:", "- `Sensitive`:", "- `Write-Only`:"} {
		if strings.Contains(requiredSection, marker) || strings.Contains(optionalSection, marker) {
			return fmt.Errorf("documentation file %s has qualifier guidance mixed into argument lists", path)
		}
	}
	return nil
}

func validateExampleBounds(path string, content string, min int, max int) error {
	exampleSection, err := extractTopLevelSection(content, "## Example Usage")
	if err != nil {
		return fmt.Errorf("documentation file %s: %w", path, err)
	}
	count := len(hclFencePattern.FindAllString(exampleSection, -1))
	if count < min || count > max {
		return fmt.Errorf("documentation file %s has %d examples; expected %d-%d", path, count, min, max)
	}
	return nil
}

func validateFurtherReadingPolicy(path string, content string, expectedAWXLinks []docsLink) error {
	readingSection, err := extractTopLevelSection(content, "## Further Reading")
	if err != nil {
		return fmt.Errorf("documentation file %s: %w", path, err)
	}
	matches := markdownLinkPattern.FindAllStringSubmatch(readingSection, -1)
	if len(matches) == 0 {
		return fmt.Errorf("documentation file %s has no links in Further Reading", path)
	}

	hasAWX := false
	hasGenericAWXIndex := false
	foundExpected := make(map[string]struct{}, len(expectedAWXLinks))
	expectedURLs := make(map[string]struct{}, len(expectedAWXLinks))
	for _, link := range expectedAWXLinks {
		expectedURLs[strings.TrimSpace(link.URL)] = struct{}{}
	}
	for _, match := range matches {
		linkURL := strings.TrimSpace(match[1])
		if isOfficialAWXLink(linkURL) {
			hasAWX = true
			if isGenericAWXIndexLink(linkURL) {
				hasGenericAWXIndex = true
			}
			if _, expected := expectedURLs[linkURL]; expected {
				foundExpected[linkURL] = struct{}{}
			}
			continue
		}
		return fmt.Errorf("documentation file %s includes non-AWX links in Further Reading; use official AWX documentation links only", path)
	}

	if !hasAWX {
		return fmt.Errorf("documentation file %s requires official AWX references in Further Reading", path)
	}
	if hasGenericAWXIndex {
		return fmt.Errorf("documentation file %s includes a generic AWX index link; use resource-specific official AWX links", path)
	}
	if len(expectedURLs) > 0 && len(foundExpected) == 0 {
		return fmt.Errorf("documentation file %s is missing expected resource-specific official AWX links", path)
	}
	return nil
}

func matchesExpectedAWXConceptURL(actualURL string, expectedAWXLinks []docsLink) bool {
	actual := strings.TrimSpace(actualURL)
	if actual == "" {
		return false
	}
	actualParsed, err := url.Parse(actual)
	if err != nil || actualParsed.Scheme == "" || actualParsed.Host == "" {
		return false
	}
	actualHost := strings.ToLower(actualParsed.Host)
	actualPath := strings.TrimSpace(actualParsed.Path)
	for _, expected := range expectedAWXLinks {
		expectedURL := strings.TrimSpace(expected.URL)
		expectedParsed, err := url.Parse(expectedURL)
		if err != nil || expectedParsed.Scheme == "" || expectedParsed.Host == "" {
			continue
		}
		expectedHost := strings.ToLower(expectedParsed.Host)
		expectedPath := strings.TrimSpace(expectedParsed.Path)
		if actualHost == expectedHost && actualPath == expectedPath {
			return true
		}
	}
	return false
}

func validateComplexPrimer(path string, content string) error {
	conceptsSection, err := extractTopLevelSection(content, "## AWX Concepts")
	if err != nil {
		return fmt.Errorf("documentation file %s requires AWX Concepts section for complex resources", path)
	}
	if strings.TrimSpace(conceptsSection) == "" {
		return fmt.Errorf("documentation file %s has an empty AWX Concepts section", path)
	}
	return nil
}

func validateQualityAnalysisSummary(path string, maxPasses int) error {
	raw, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("quality analysis summary missing: %s", path)
	}
	content := string(raw)
	matches := qualityAnalysisPassHeadingPattern.FindAllStringSubmatchIndex(content, -1)
	if len(matches) == 0 {
		return fmt.Errorf("quality analysis summary %s is missing required \"## Quality Analysis Pass <n>\" sections", path)
	}

	passNumbers := make([]int, 0, len(matches))
	for _, match := range matches {
		pass, err := strconv.Atoi(content[match[2]:match[3]])
		if err != nil {
			return fmt.Errorf("quality analysis summary %s contains an invalid pass heading", path)
		}
		passNumbers = append(passNumbers, pass)
	}
	if passNumbers[0] != 1 {
		return fmt.Errorf("quality analysis summary %s must start at pass 1; found pass %d first", path, passNumbers[0])
	}
	for i := 1; i < len(passNumbers); i++ {
		if passNumbers[i] != passNumbers[i-1]+1 {
			return fmt.Errorf("quality analysis summary %s has non-contiguous pass numbering: found pass %d after pass %d", path, passNumbers[i], passNumbers[i-1])
		}
	}
	if passNumbers[len(passNumbers)-1] > maxPasses {
		return fmt.Errorf("quality analysis summary %s exceeds the maximum of %d passes", path, maxPasses)
	}

	for idx, match := range matches {
		start := match[0]
		end := len(content)
		if idx+1 < len(matches) {
			end = matches[idx+1][0]
		}
		section := content[start:end]
		if !strings.Contains(section, "### Pass result") {
			return fmt.Errorf("quality analysis summary %s is missing \"### Pass result\" for pass %d", path, passNumbers[idx])
		}
	}
	return nil
}

func verifyDocsEnrichmentOnline(client *http.Client, objects []manifest.ManagedObject, enrichment docsEnrichmentCatalog) error {
	for _, object := range objects {
		if object.RuntimeExcluded || (!object.ResourceEligible && !object.DataSourceElig) {
			continue
		}
		metadata, metadataKey, ok := objectEnrichmentMetadataForTarget(object, enrichment)
		if !ok || metadata.CurationSource == nil {
			return fmt.Errorf("docs enrichment objects.%s missing curationSource for online verification", object.ResourceName)
		}
		expectedLinks, err := awxOfficialLinksForObject(object.Name)
		if err != nil {
			return err
		}
		if err := verifyOfficialAWXLinkOnline(client, metadata.CurationSource.OfficialAWXURL, expectedLinks); err != nil {
			return fmt.Errorf("docs enrichment objects.%s online verification failed: %w", metadataKey, err)
		}
	}
	return nil
}

func verifyOfficialAWXLinkOnline(client *http.Client, officialURL string, expectedLinks []docsLink) error {
	urlValue := strings.TrimSpace(officialURL)
	if !isOfficialAWXLink(urlValue) {
		return fmt.Errorf("official link %q is not an AWX documentation URL", urlValue)
	}
	if !matchesExpectedAWXConceptURL(urlValue, expectedLinks) {
		return fmt.Errorf("official link %q does not match expected AWX concept mapping", urlValue)
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, urlValue, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("User-Agent", "awxgen-docs-verify-online/1.0")
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("fetch official AWX URL: %w", err)
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 32*1024))
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("fetch official AWX URL returned unexpected status %d", resp.StatusCode)
	}
	return nil
}

func extractTopLevelSection(content string, heading string) (string, error) {
	start := strings.Index(content, heading)
	if start == -1 {
		return "", fmt.Errorf("missing section %q", heading)
	}
	remaining := content[start+len(heading):]
	next := strings.Index(remaining, "\n## ")
	if next == -1 {
		return strings.TrimSpace(remaining), nil
	}
	return strings.TrimSpace(remaining[:next]), nil
}

func extractSubsection(content string, heading string) (string, error) {
	start := strings.Index(content, heading)
	if start == -1 {
		return "", fmt.Errorf("missing subsection %q", heading)
	}
	remaining := content[start+len(heading):]
	next := strings.Index(remaining, "\n### ")
	if next == -1 {
		return strings.TrimSpace(remaining), nil
	}
	return strings.TrimSpace(remaining[:next]), nil
}

func isOfficialAWXLink(link string) bool {
	parsed, err := url.Parse(link)
	if err != nil {
		return false
	}
	host := strings.ToLower(parsed.Host)
	if !strings.Contains(parsed.Path, "/awx/") {
		return false
	}
	return strings.Contains(host, "ansible.readthedocs.io") || strings.Contains(host, "docs.ansible.com")
}

func isGenericAWXIndexLink(link string) bool {
	normalized := strings.TrimSpace(link)
	if normalized == "" {
		return false
	}
	parsed, err := url.Parse(normalized)
	if err == nil && parsed.Scheme != "" && parsed.Host != "" {
		normalized = fmt.Sprintf("%s://%s%s", parsed.Scheme, strings.ToLower(parsed.Host), parsed.Path)
	}
	_, exists := genericAWXIndexURLs[normalized]
	return exists
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
