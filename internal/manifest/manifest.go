package manifest

import (
	"embed"
	"encoding/json"
	"fmt"
	"sort"
)

//go:embed managed_objects.json relationships.json runtime_exclusions.json
var fs embed.FS

// FieldType describes field shape derived from OpenAPI.
type FieldType string

const (
	FieldTypeString FieldType = "string"
	FieldTypeInt    FieldType = "integer"
	FieldTypeBool   FieldType = "boolean"
	FieldTypeFloat  FieldType = "number"
	FieldTypeArray  FieldType = "array"
	FieldTypeObject FieldType = "object"
)

// FieldSpec maps OpenAPI request schema fields into Terraform schema fields.
type FieldSpec struct {
	Name        string    `json:"name"`
	Type        FieldType `json:"type"`
	Required    bool      `json:"required"`
	Sensitive   bool      `json:"sensitive"`
	WriteOnly   bool      `json:"writeOnly"`
	Description string    `json:"description,omitempty"`
}

// ManagedObject describes one AWX API object candidate.
type ManagedObject struct {
	Name             string      `json:"name"`
	SingularName     string      `json:"singularName"`
	ResourceName     string      `json:"resourceName"`
	DataSourceName   string      `json:"dataSourceName"`
	CollectionPath   string      `json:"collectionPath"`
	DetailPath       string      `json:"detailPath"`
	RequestSchema    string      `json:"requestSchema,omitempty"`
	ResponseSchema   string      `json:"responseSchema,omitempty"`
	ResourceEligible bool        `json:"resourceEligible"`
	DataSourceElig   bool        `json:"dataSourceEligible"`
	RuntimeExcluded  bool        `json:"runtimeExcluded"`
	ExclusionReason  string      `json:"exclusionReason,omitempty"`
	Fields           []FieldSpec `json:"fields"`
}

// Relationship describes a parent-child association resource candidate.
type Relationship struct {
	Name         string `json:"name"`
	ResourceName string `json:"resourceName"`
	ParentObject string `json:"parentObject"`
	ChildObject  string `json:"childObject"`
	Path         string `json:"path"`
	Priority     int    `json:"priority"`
}

// RuntimeExclusion identifies runtime-only AWX objects omitted from managed resources.
type RuntimeExclusion struct {
	Object string `json:"object"`
	Reason string `json:"reason"`
}

// Catalog bundles managed objects and relationship metadata.
type Catalog struct {
	ManagedObjects   []ManagedObject    `json:"managedObjects"`
	Relationships    []Relationship     `json:"relationships"`
	RuntimeExclusion []RuntimeExclusion `json:"runtimeExclusions"`

	managedByName map[string]ManagedObject
}

// Load reads generated manifest assets.
func Load() (*Catalog, error) {
	objectsRaw, err := fs.ReadFile("managed_objects.json")
	if err != nil {
		return nil, fmt.Errorf("read managed objects manifest: %w", err)
	}
	relationsRaw, err := fs.ReadFile("relationships.json")
	if err != nil {
		return nil, fmt.Errorf("read relationships manifest: %w", err)
	}
	exclusionsRaw, err := fs.ReadFile("runtime_exclusions.json")
	if err != nil {
		return nil, fmt.Errorf("read runtime exclusions manifest: %w", err)
	}

	catalog := &Catalog{}
	if err := json.Unmarshal(objectsRaw, &catalog.ManagedObjects); err != nil {
		return nil, fmt.Errorf("unmarshal managed objects manifest: %w", err)
	}
	if err := json.Unmarshal(relationsRaw, &catalog.Relationships); err != nil {
		return nil, fmt.Errorf("unmarshal relationships manifest: %w", err)
	}
	if err := unmarshalRuntimeExclusions(exclusionsRaw, &catalog.RuntimeExclusion); err != nil {
		return nil, fmt.Errorf("unmarshal runtime exclusions manifest: %w", err)
	}

	catalog.managedByName = make(map[string]ManagedObject, len(catalog.ManagedObjects))
	for _, obj := range catalog.ManagedObjects {
		catalog.managedByName[obj.Name] = obj
	}

	sort.SliceStable(catalog.ManagedObjects, func(i, j int) bool {
		return catalog.ManagedObjects[i].Name < catalog.ManagedObjects[j].Name
	})
	sort.SliceStable(catalog.Relationships, func(i, j int) bool {
		if catalog.Relationships[i].Priority == catalog.Relationships[j].Priority {
			return catalog.Relationships[i].Name < catalog.Relationships[j].Name
		}
		return catalog.Relationships[i].Priority < catalog.Relationships[j].Priority
	})

	return catalog, nil
}

// MustLoad panics when generated manifests are invalid.
func MustLoad() *Catalog {
	catalog, err := Load()
	if err != nil {
		panic(err)
	}
	return catalog
}

// ManagedResourceObjects returns all non-excluded resource-eligible objects.
func (c *Catalog) ManagedResourceObjects() []ManagedObject {
	out := make([]ManagedObject, 0)
	for _, obj := range c.ManagedObjects {
		if obj.ResourceEligible && !obj.RuntimeExcluded {
			out = append(out, obj)
		}
	}
	return out
}

// ManagedDataSourceObjects returns all data-source eligible objects.
func (c *Catalog) ManagedDataSourceObjects() []ManagedObject {
	out := make([]ManagedObject, 0)
	for _, obj := range c.ManagedObjects {
		if obj.DataSourceElig && !obj.RuntimeExcluded {
			out = append(out, obj)
		}
	}
	return out
}

// ObjectByName gets a managed object definition by collection name.
func (c *Catalog) ObjectByName(name string) (ManagedObject, bool) {
	obj, ok := c.managedByName[name]
	return obj, ok
}

func unmarshalRuntimeExclusions(raw []byte, target *[]RuntimeExclusion) error {
	var direct []RuntimeExclusion
	if err := json.Unmarshal(raw, &direct); err == nil {
		*target = direct
		return nil
	}

	var wrapped struct {
		Exclusions []RuntimeExclusion `json:"exclusions"`
	}
	if err := json.Unmarshal(raw, &wrapped); err != nil {
		return err
	}
	*target = wrapped.Exclusions
	return nil
}
