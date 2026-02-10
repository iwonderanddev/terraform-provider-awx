package openapi

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DeprecatedObjectExclusion marks an object endpoint as intentionally excluded.
type DeprecatedObjectExclusion struct {
	Object string `json:"object"`
	Reason string `json:"reason"`
}

// DeprecatedRelationshipExclusion marks a relationship endpoint path as intentionally excluded.
type DeprecatedRelationshipExclusion struct {
	Path   string `json:"path"`
	Reason string `json:"reason"`
}

// DeprecatedExclusionFile stores deprecated object and relationship exclusions.
type DeprecatedExclusionFile struct {
	Objects       []DeprecatedObjectExclusion       `json:"objects"`
	Relationships []DeprecatedRelationshipExclusion `json:"relationships"`
}

// LoadDeprecatedExclusions reads deprecated endpoint exclusions from disk.
func LoadDeprecatedExclusions(path string) (map[string]string, map[string]string, error) {
	raw, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]string{}, map[string]string{}, nil
		}
		return nil, nil, fmt.Errorf("read deprecated exclusions: %w", err)
	}
	if len(strings.TrimSpace(string(raw))) == 0 {
		return map[string]string{}, map[string]string{}, nil
	}

	var payload DeprecatedExclusionFile
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, nil, fmt.Errorf("parse deprecated exclusions JSON: %w", err)
	}

	objectExclusions := make(map[string]string, len(payload.Objects))
	for _, object := range payload.Objects {
		if strings.TrimSpace(object.Object) == "" {
			continue
		}
		objectExclusions[object.Object] = strings.TrimSpace(object.Reason)
	}

	relationshipExclusions := make(map[string]string, len(payload.Relationships))
	for _, relationship := range payload.Relationships {
		if strings.TrimSpace(relationship.Path) == "" {
			continue
		}
		relationshipExclusions[relationship.Path] = strings.TrimSpace(relationship.Reason)
	}

	return objectExclusions, relationshipExclusions, nil
}
