package provider

import "testing"

func TestValidateObjectImportIDForCollectionResource(t *testing.T) {
	t.Parallel()

	id, err := validateObjectImportID("42", true)
	if err != nil {
		t.Fatalf("expected numeric import ID to be valid, got error: %v", err)
	}
	if id != "42" {
		t.Fatalf("unexpected normalized import id: got=%q want=%q", id, "42")
	}

	if _, err := validateObjectImportID("system", true); err == nil {
		t.Fatalf("expected non-numeric import ID to fail for collection-created object")
	}
}

func TestValidateObjectImportIDForDetailKeyResource(t *testing.T) {
	t.Parallel()

	id, err := validateObjectImportID("system", false)
	if err != nil {
		t.Fatalf("expected detail-key import ID to be valid, got error: %v", err)
	}
	if id != "system" {
		t.Fatalf("unexpected normalized import id: got=%q want=%q", id, "system")
	}

	if _, err := validateObjectImportID("   ", false); err == nil {
		t.Fatalf("expected empty import ID to fail")
	}
}
