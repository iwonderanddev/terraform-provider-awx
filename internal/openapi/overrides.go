package openapi

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/damien/terraform-awx-provider/internal/manifest"
)

// FieldOverride updates derived field metadata for schema/runtime mismatches.
type FieldOverride struct {
	Object      string             `json:"object"`
	Field       string             `json:"field"`
	Type        manifest.FieldType `json:"type,omitempty"`
	Required    *bool              `json:"required,omitempty"`
	Sensitive   *bool              `json:"sensitive,omitempty"`
	WriteOnly   *bool              `json:"writeOnly,omitempty"`
	Description string             `json:"description,omitempty"`
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
		for _, field := range object.Fields {
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
			if override.Sensitive != nil {
				field.Sensitive = *override.Sensitive
			}
			if override.WriteOnly != nil {
				field.WriteOnly = *override.WriteOnly
			}
			if strings.TrimSpace(override.Description) != "" {
				field.Description = override.Description
			}
			fields = append(fields, field)
		}
		object.Fields = fields
		updated = append(updated, object)
	}

	return updated
}
