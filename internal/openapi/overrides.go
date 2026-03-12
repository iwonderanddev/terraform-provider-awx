package openapi

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/damien/terraform-provider-awx-iwd/internal/manifest"
)

// FieldOverride updates derived field metadata for schema/runtime mismatches.
type FieldOverride struct {
	Object        string             `json:"object"`
	Field         string             `json:"field"`
	TerraformName string             `json:"terraformName,omitempty"`
	Type          manifest.FieldType `json:"type,omitempty"`
	Required      *bool              `json:"required,omitempty"`
	Reference     *bool              `json:"reference,omitempty"`
	Computed      *bool              `json:"computed,omitempty"`
	Sensitive     *bool              `json:"sensitive,omitempty"`
	WriteOnly     *bool              `json:"writeOnly,omitempty"`
	Description   string             `json:"description,omitempty"`
}

// FieldOverrideFile stores override declarations.
type FieldOverrideFile struct {
	Overrides []FieldOverride `json:"overrides"`
}

// LoadFieldOverrides loads override definitions from disk.
func LoadFieldOverrides(path string) (map[string]FieldOverride, error) {
	raw, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]FieldOverride{}, nil
		}
		return nil, fmt.Errorf("read field overrides: %w", err)
	}
	if len(strings.TrimSpace(string(raw))) == 0 {
		return map[string]FieldOverride{}, nil
	}

	var payload FieldOverrideFile
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, fmt.Errorf("parse field overrides JSON: %w", err)
	}

	overrides := make(map[string]FieldOverride, len(payload.Overrides))
	for _, override := range payload.Overrides {
		key := override.Object + "." + override.Field
		overrides[key] = override
	}
	return overrides, nil
}

// ApplyFieldOverrides applies override metadata to derived fields.
func ApplyFieldOverrides(objects []manifest.ManagedObject, overrides map[string]FieldOverride) []manifest.ManagedObject {
	if len(overrides) == 0 {
		return objects
	}

	updated := make([]manifest.ManagedObject, 0, len(objects))
	for _, object := range objects {
		fields := make([]manifest.FieldSpec, 0, len(object.Fields))
		existingFields := make(map[string]struct{}, len(object.Fields))
		for _, field := range object.Fields {
			existingFields[field.Name] = struct{}{}

			key := object.Name + "." + field.Name
			override, ok := overrides[key]
			if !ok {
				fields = append(fields, field)
				continue
			}

			if override.Type != "" {
				field.Type = override.Type
			}
			if override.Required != nil {
				field.Required = *override.Required
			}
			if override.Computed != nil {
				field.Computed = *override.Computed
			}
			if override.Reference != nil {
				field.Reference = *override.Reference
			}
			if override.Sensitive != nil {
				field.Sensitive = *override.Sensitive
			}
			if override.WriteOnly != nil {
				field.WriteOnly = *override.WriteOnly
			}
			if strings.TrimSpace(override.Description) != "" {
				field.Description = override.Description
			}
			if strings.TrimSpace(override.TerraformName) != "" {
				field.TerraformName = strings.TrimSpace(override.TerraformName)
			}
			fields = append(fields, field)
		}

		additions := additionsForObject(overrides, object.Name)
		for _, override := range additions {
			if _, exists := existingFields[override.Field]; exists {
				continue
			}
			fields = append(fields, fieldFromOverride(override))
		}

		sort.SliceStable(fields, func(i, j int) bool {
			return fields[i].Name < fields[j].Name
		})

		object.Fields = fields
		updated = append(updated, object)
	}

	return updated
}

func additionsForObject(overrides map[string]FieldOverride, objectName string) []FieldOverride {
	out := make([]FieldOverride, 0)
	for _, override := range overrides {
		if override.Object != objectName {
			continue
		}
		if strings.TrimSpace(override.Field) == "" {
			continue
		}
		out = append(out, override)
	}
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].Field < out[j].Field
	})
	return out
}

func fieldFromOverride(override FieldOverride) manifest.FieldSpec {
	fieldType := override.Type
	if fieldType == "" {
		fieldType = manifest.FieldTypeString
	}

	required := false
	if override.Required != nil {
		required = *override.Required
	}

	computed := false
	if override.Computed != nil {
		computed = *override.Computed
	}
	reference := false
	if override.Reference != nil {
		reference = *override.Reference
	}

	sensitive := false
	if override.Sensitive != nil {
		sensitive = *override.Sensitive
	}

	writeOnly := sensitive
	if override.WriteOnly != nil {
		writeOnly = *override.WriteOnly
	}

	return manifest.FieldSpec{
		Name:          override.Field,
		TerraformName: strings.TrimSpace(override.TerraformName),
		Type:          fieldType,
		Required:      required,
		Reference:     reference,
		Computed:      computed,
		Sensitive:     sensitive,
		WriteOnly:     writeOnly,
		Description:   strings.TrimSpace(override.Description),
	}
}
