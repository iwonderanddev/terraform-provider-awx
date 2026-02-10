package openapi

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDeprecatedExclusions(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "deprecated.json")
	payload := `{
  "objects": [
    {
      "object": "roles",
      "reason": "Deprecated in favor of role_definitions."
    }
  ],
  "relationships": [
    {
      "path": "/api/v2/users/{id}/roles/",
      "reason": "Deprecated in favor of role_user_assignments."
    }
  ]
}`
	if err := os.WriteFile(path, []byte(payload), 0o600); err != nil {
		t.Fatalf("write test fixture: %v", err)
	}

	objectExclusions, relationshipExclusions, err := LoadDeprecatedExclusions(path)
	if err != nil {
		t.Fatalf("LoadDeprecatedExclusions returned error: %v", err)
	}

	if reason := objectExclusions["roles"]; reason == "" {
		t.Fatalf("expected roles object exclusion to be loaded")
	}
	if reason := relationshipExclusions["/api/v2/users/{id}/roles/"]; reason == "" {
		t.Fatalf("expected users roles relationship exclusion to be loaded")
	}
}

func TestLoadDeprecatedExclusionsMissingFile(t *testing.T) {
	t.Parallel()

	objectExclusions, relationshipExclusions, err := LoadDeprecatedExclusions("/tmp/does-not-exist/deprecated.json")
	if err != nil {
		t.Fatalf("LoadDeprecatedExclusions returned error for missing file: %v", err)
	}
	if len(objectExclusions) != 0 {
		t.Fatalf("expected no object exclusions for missing file, got %d", len(objectExclusions))
	}
	if len(relationshipExclusions) != 0 {
		t.Fatalf("expected no relationship exclusions for missing file, got %d", len(relationshipExclusions))
	}
}
