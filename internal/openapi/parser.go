package openapi

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/damien/terraform-awx-provider/internal/manifest"
)

var (
	collectionPathPattern   = regexp.MustCompile(`^/api/v2/([a-z0-9_]+)/$`)
	relationshipPathPattern = regexp.MustCompile(`^/api/v2/([a-z0-9_]+)/\{id\}/([a-z0-9_]+)/$`)
)

// Document is the AWX OpenAPI schema subset used by generation.
type Document struct {
	Paths      map[string]PathItem `json:"paths"`
	Components Components          `json:"components"`
}

// Components wraps OpenAPI component schemas.
type Components struct {
	Schemas map[string]*Schema `json:"schemas"`
}

// PathItem captures endpoint operations.
type PathItem struct {
	Get    *Operation `json:"get"`
	Post   *Operation `json:"post"`
	Put    *Operation `json:"put"`
	Patch  *Operation `json:"patch"`
	Delete *Operation `json:"delete"`
}

// Operation captures operation details used by manifest generation.
type Operation struct {
	OperationID string              `json:"operationId"`
	RequestBody *RequestBody        `json:"requestBody"`
	Responses   map[string]Response `json:"responses"`
}

// RequestBody captures JSON request body schema references.
type RequestBody struct {
	Content map[string]MediaType `json:"content"`
}

// Response captures JSON response schema references.
type Response struct {
	Content map[string]MediaType `json:"content"`
}

// MediaType wraps schema references in operation bodies.
type MediaType struct {
	Schema *Schema `json:"schema"`
}

// Schema captures the subset of OpenAPI schema metadata needed for generation.
type Schema struct {
	Ref         string             `json:"$ref"`
	Type        string             `json:"type"`
	Format      string             `json:"format"`
	Default     any                `json:"default"`
	WriteOnly   bool               `json:"writeOnly"`
	Description string             `json:"description"`
	Properties  map[string]*Schema `json:"properties"`
	Required    []string           `json:"required"`
	Items       *Schema            `json:"items"`
	AllOf       []*Schema          `json:"allOf"`
	AnyOf       []*Schema          `json:"anyOf"`
	OneOf       []*Schema          `json:"oneOf"`
}

// RelationshipPriority controls generation order for association resources.
type RelationshipPriority struct {
	Name     string `json:"name"`
	Priority int    `json:"priority"`
}

// RuntimeExclusionFile stores explicit runtime-only exclusions.
type RuntimeExclusionFile struct {
	Exclusions []manifest.RuntimeExclusion `json:"exclusions"`
}

// LoadDocument reads an OpenAPI document from disk.
func LoadDocument(path string) (*Document, error) {
	raw, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("read OpenAPI document: %w", err)
	}

	var doc Document
	if err := json.Unmarshal(raw, &doc); err != nil {
		return nil, fmt.Errorf("parse OpenAPI document: %w", err)
	}
	return &doc, nil
}

// LoadRuntimeExclusions reads runtime-only exclusion definitions from disk.
func LoadRuntimeExclusions(path string) (map[string]manifest.RuntimeExclusion, error) {
	raw, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]manifest.RuntimeExclusion{}, nil
		}
		return nil, fmt.Errorf("read runtime exclusions: %w", err)
	}

	if len(strings.TrimSpace(string(raw))) == 0 {
		return map[string]manifest.RuntimeExclusion{}, nil
	}

	var payload RuntimeExclusionFile
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, fmt.Errorf("parse runtime exclusions JSON: %w", err)
	}

	exclusions := make(map[string]manifest.RuntimeExclusion, len(payload.Exclusions))
	for _, ex := range payload.Exclusions {
		exclusions[ex.Object] = ex
	}
	return exclusions, nil
}

// LoadRelationshipPriorities reads prioritized association definitions from disk.
func LoadRelationshipPriorities(path string) (map[string]int, error) {
	raw, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]int{}, nil
		}
		return nil, fmt.Errorf("read relationship priorities: %w", err)
	}
	if len(strings.TrimSpace(string(raw))) == 0 {
		return map[string]int{}, nil
	}

	var items []RelationshipPriority
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, fmt.Errorf("parse relationship priorities JSON: %w", err)
	}

	out := make(map[string]int, len(items))
	for _, item := range items {
		out[item.Name] = item.Priority
	}
	return out, nil
}

// DeriveManagedObjects derives managed object metadata from OpenAPI paths and component schemas.
func DeriveManagedObjects(doc *Document, exclusions map[string]manifest.RuntimeExclusion) []manifest.ManagedObject {
	if doc == nil {
		return nil
	}

	paths := sortedPathKeys(doc.Paths)
	objects := make([]manifest.ManagedObject, 0)

	for _, endpointPath := range paths {
		matches := collectionPathPattern.FindStringSubmatch(endpointPath)
		if len(matches) != 2 {
			continue
		}

		collectionName := matches[1]
		collectionOps := doc.Paths[endpointPath]
		detailPath := fmt.Sprintf("/api/v2/%s/{id}/", collectionName)
		detailOps, ok := doc.Paths[detailPath]
		if !ok || detailOps.Get == nil {
			continue
		}

		requestSchema := requestSchemaName(collectionOps.Post)
		if requestSchema == "" {
			requestSchema = requestSchemaName(detailOps.Patch)
		}
		if requestSchema == "" {
			requestSchema = requestSchemaName(detailOps.Put)
		}

		responseSchema := responseSchemaName(detailOps.Get)
		fields := fieldsFromSchema(doc, requestSchema)
		if len(fields) == 0 {
			fields = fieldsFromSchema(doc, responseSchema)
		}

		exclusion, runtimeExcluded := exclusions[collectionName]
		resourceEligible := collectionOps.Post != nil && detailOps.Delete != nil && (detailOps.Patch != nil || detailOps.Put != nil)
		if runtimeExcluded {
			resourceEligible = false
		}

		obj := manifest.ManagedObject{
			Name:             collectionName,
			SingularName:     singularize(collectionName),
			ResourceName:     fmt.Sprintf("awx_%s", singularize(collectionName)),
			DataSourceName:   fmt.Sprintf("awx_%s", singularize(collectionName)),
			CollectionPath:   endpointPath,
			DetailPath:       detailPath,
			RequestSchema:    requestSchema,
			ResponseSchema:   responseSchema,
			ResourceEligible: resourceEligible,
			DataSourceElig:   collectionOps.Get != nil && detailOps.Get != nil,
			RuntimeExcluded:  runtimeExcluded,
			ExclusionReason:  exclusion.Reason,
			Fields:           fields,
		}
		objects = append(objects, obj)
	}

	sort.SliceStable(objects, func(i, j int) bool {
		return objects[i].Name < objects[j].Name
	})
	return objects
}

// DeriveRelationships derives relationship resource candidates.
func DeriveRelationships(doc *Document, managedObjects []manifest.ManagedObject, priorities map[string]int) []manifest.Relationship {
	if doc == nil {
		return nil
	}

	objectByCollection := make(map[string]manifest.ManagedObject, len(managedObjects))
	for _, obj := range managedObjects {
		objectByCollection[obj.Name] = obj
	}

	relationships := make([]manifest.Relationship, 0)
	seen := make(map[string]struct{})

	for _, endpointPath := range sortedPathKeys(doc.Paths) {
		matches := relationshipPathPattern.FindStringSubmatch(endpointPath)
		if len(matches) != 3 {
			continue
		}
		parentCollection := matches[1]
		childCollection := matches[2]

		if !isRelationshipCandidate(childCollection) {
			continue
		}

		if _, ok := objectByCollection[parentCollection]; !ok {
			continue
		}
		if _, ok := objectByCollection[childCollection]; !ok {
			continue
		}

		ops := doc.Paths[endpointPath]
		if ops.Get == nil || ops.Post == nil {
			continue
		}

		name := fmt.Sprintf("%s_%s_association", singularize(parentCollection), singularize(childCollection))
		if _, exists := seen[name]; exists {
			continue
		}
		seen[name] = struct{}{}

		priority := 100
		if explicit, ok := priorities[name]; ok {
			priority = explicit
		}

		relationships = append(relationships, manifest.Relationship{
			Name:         name,
			ResourceName: fmt.Sprintf("awx_%s", name),
			ParentObject: parentCollection,
			ChildObject:  childCollection,
			Path:         endpointPath,
			Priority:     priority,
		})
	}

	sort.SliceStable(relationships, func(i, j int) bool {
		if relationships[i].Priority == relationships[j].Priority {
			return relationships[i].Name < relationships[j].Name
		}
		return relationships[i].Priority < relationships[j].Priority
	})
	return relationships
}

// ValidateCoverage ensures candidate objects are either implemented or explicitly excluded.
func ValidateCoverage(objects []manifest.ManagedObject, exclusions map[string]manifest.RuntimeExclusion) error {
	missing := make([]string, 0)
	for _, obj := range objects {
		if obj.ResourceEligible || obj.RuntimeExcluded {
			continue
		}
		if !requiresRuntimeExclusion(obj.Name) {
			continue
		}
		if _, ok := exclusions[obj.Name]; !ok {
			missing = append(missing, obj.Name)
		}
	}

	if len(missing) > 0 {
		sort.Strings(missing)
		return fmt.Errorf("objects missing explicit runtime exclusion entries: %s", strings.Join(missing, ", "))
	}
	return nil
}

func requiresRuntimeExclusion(objectName string) bool {
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

func sortedPathKeys(paths map[string]PathItem) []string {
	keys := make([]string, 0, len(paths))
	for k := range paths {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func requestSchemaName(op *Operation) string {
	if op == nil || op.RequestBody == nil {
		return ""
	}
	if media, ok := op.RequestBody.Content["application/json"]; ok {
		return refName(media.Schema.Ref)
	}
	for _, media := range op.RequestBody.Content {
		if media.Schema != nil && media.Schema.Ref != "" {
			return refName(media.Schema.Ref)
		}
	}
	return ""
}

func responseSchemaName(op *Operation) string {
	if op == nil {
		return ""
	}
	for _, status := range []string{"200", "201"} {
		response, ok := op.Responses[status]
		if !ok {
			continue
		}
		if media, ok := response.Content["application/json"]; ok {
			return refName(media.Schema.Ref)
		}
		for _, media := range response.Content {
			if media.Schema != nil && media.Schema.Ref != "" {
				return refName(media.Schema.Ref)
			}
		}
	}
	return ""
}

func refName(ref string) string {
	const prefix = "#/components/schemas/"
	if strings.HasPrefix(ref, prefix) {
		return strings.TrimPrefix(ref, prefix)
	}
	return strings.TrimSpace(ref)
}

func fieldsFromSchema(doc *Document, schemaName string) []manifest.FieldSpec {
	if doc == nil || schemaName == "" {
		return nil
	}

	schema := doc.Components.Schemas[schemaName]
	if schema == nil {
		return nil
	}
	resolved := resolveSchema(doc, schema)
	if resolved == nil {
		return nil
	}

	required := make(map[string]struct{}, len(resolved.Required))
	for _, name := range resolved.Required {
		required[name] = struct{}{}
	}

	fieldNames := make([]string, 0, len(resolved.Properties))
	for name := range resolved.Properties {
		fieldNames = append(fieldNames, name)
	}
	sort.Strings(fieldNames)

	out := make([]manifest.FieldSpec, 0, len(fieldNames))
	for _, name := range fieldNames {
		property := resolveSchema(doc, resolved.Properties[name])
		if property == nil {
			continue
		}

		sensitive := property.WriteOnly || isSensitiveField(name, property)
		_, isRequired := required[name]
		computed := shouldInferComputedFromDefault(property, isRequired)
		out = append(out, manifest.FieldSpec{
			Name:        name,
			Type:        normalizeFieldType(property),
			Required:    isRequired,
			Computed:    computed,
			Sensitive:   sensitive,
			WriteOnly:   property.WriteOnly || sensitive,
			Description: strings.TrimSpace(property.Description),
		})
	}
	return out
}

func resolveSchema(doc *Document, schema *Schema) *Schema {
	if doc == nil || schema == nil {
		return nil
	}

	if schema.Ref != "" {
		resolved := doc.Components.Schemas[refName(schema.Ref)]
		if resolved == nil {
			return nil
		}
		return resolveSchema(doc, resolved)
	}

	out := &Schema{
		Type:        schema.Type,
		Format:      schema.Format,
		Default:     schema.Default,
		WriteOnly:   schema.WriteOnly,
		Description: schema.Description,
		Properties:  make(map[string]*Schema),
		Required:    append([]string{}, schema.Required...),
		Items:       schema.Items,
	}

	for name, property := range schema.Properties {
		out.Properties[name] = resolveSchema(doc, property)
	}

	mergeComposedSchema := func(items []*Schema) {
		for _, part := range items {
			resolved := resolveSchema(doc, part)
			if resolved == nil {
				continue
			}
			if out.Type == "" {
				out.Type = resolved.Type
			}
			if out.Format == "" {
				out.Format = resolved.Format
			}
			if out.Default == nil && resolved.Default != nil {
				out.Default = resolved.Default
			}
			if resolved.WriteOnly {
				out.WriteOnly = true
			}
			if out.Description == "" {
				out.Description = resolved.Description
			}
			for name, property := range resolved.Properties {
				out.Properties[name] = property
			}
			out.Required = append(out.Required, resolved.Required...)
		}
	}

	mergeComposedSchema(schema.AllOf)
	mergeComposedSchema(schema.OneOf)
	mergeComposedSchema(schema.AnyOf)

	if out.Properties == nil {
		out.Properties = map[string]*Schema{}
	}
	out.Required = uniqueSorted(out.Required)
	return out
}

func shouldInferComputedFromDefault(property *Schema, required bool) bool {
	if property == nil || required || property.WriteOnly || property.Default == nil {
		return false
	}

	// AWX commonly sets empty-string defaults for optional text fields such as
	// descriptions; keep those as plain Optional to avoid over-marking computed.
	if defaultString, ok := property.Default.(string); ok && defaultString == "" {
		return false
	}

	return true
}

func normalizeFieldType(schema *Schema) manifest.FieldType {
	if schema == nil {
		return manifest.FieldTypeString
	}
	typ := schema.Type
	if typ == "" {
		if schema.Items != nil {
			typ = "array"
		} else if len(schema.Properties) > 0 {
			typ = "object"
		}
	}

	switch typ {
	case "integer":
		return manifest.FieldTypeInt
	case "number":
		return manifest.FieldTypeFloat
	case "boolean":
		return manifest.FieldTypeBool
	case "array":
		return manifest.FieldTypeArray
	case "object":
		return manifest.FieldTypeObject
	default:
		return manifest.FieldTypeString
	}
}

func singularize(collectionName string) string {
	if strings.HasSuffix(collectionName, "ies") && len(collectionName) > 3 {
		return strings.TrimSuffix(collectionName, "ies") + "y"
	}
	if strings.HasSuffix(collectionName, "sses") {
		return strings.TrimSuffix(collectionName, "es")
	}
	if strings.HasSuffix(collectionName, "ses") && len(collectionName) > 3 {
		return strings.TrimSuffix(collectionName, "es")
	}
	if strings.HasSuffix(collectionName, "s") && !strings.HasSuffix(collectionName, "ss") && len(collectionName) > 1 {
		return strings.TrimSuffix(collectionName, "s")
	}
	return collectionName
}

func isSensitiveField(fieldName string, schema *Schema) bool {
	if schema == nil {
		return false
	}
	if strings.EqualFold(schema.Format, "password") {
		return true
	}
	name := strings.ToLower(fieldName)
	sensitiveSubstrings := []string{"password", "secret", "token", "private_key", "ssh_key", "vault", "webhook_key", "passphrase"}
	for _, needle := range sensitiveSubstrings {
		if strings.Contains(name, needle) {
			return true
		}
	}
	return false
}

func isRelationshipCandidate(childCollection string) bool {
	if childCollection == "" {
		return false
	}
	blacklist := map[string]struct{}{
		"activity_stream":          {},
		"ad_hoc_commands":          {},
		"all_groups":               {},
		"all_hosts":                {},
		"ansible_facts":            {},
		"cancel":                   {},
		"copy":                     {},
		"create_schedule":          {},
		"events":                   {},
		"health_check":             {},
		"job_events":               {},
		"job_host_summaries":       {},
		"launch":                   {},
		"notifications":            {},
		"object_roles":             {},
		"owner_teams":              {},
		"owner_users":              {},
		"potential_children":       {},
		"relaunch":                 {},
		"stdout":                   {},
		"survey_spec":              {},
		"test":                     {},
		"update":                   {},
		"update_inventory_sources": {},
		"variable_data":            {},
		"webhook_key":              {},
	}
	_, blocked := blacklist[childCollection]
	return !blocked
}

func uniqueSorted(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	set := make(map[string]struct{}, len(values))
	for _, value := range values {
		if strings.TrimSpace(value) == "" {
			continue
		}
		set[value] = struct{}{}
	}
	out := make([]string, 0, len(set))
	for value := range set {
		out = append(out, value)
	}
	sort.Strings(out)
	return out
}
