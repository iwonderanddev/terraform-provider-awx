package openapi

import (
	"testing"

	"github.com/damien/terraform-awx-provider/internal/manifest"
)

func TestApplyFieldOverridesSetsComputed(t *testing.T) {
	t.Parallel()

	computed := true
	objects := []manifest.ManagedObject{
		{
			Name: "organizations",
			Fields: []manifest.FieldSpec{
				{Name: "name", Type: manifest.FieldTypeString, Required: true},
				{Name: "max_hosts", Type: manifest.FieldTypeInt, Required: false, Computed: false},
			},
		},
	}

	updated := ApplyFieldOverrides(objects, map[string]FieldOverride{
		"organizations.max_hosts": {
			Object:   "organizations",
			Field:    "max_hosts",
			Computed: &computed,
		},
	})

	var maxHosts manifest.FieldSpec
	for _, field := range updated[0].Fields {
		if field.Name == "max_hosts" {
			maxHosts = field
			break
		}
	}

	if !maxHosts.Computed {
		t.Fatalf("expected max_hosts to be computed after override")
	}
}

func TestApplyFieldOverridesAddsMissingField(t *testing.T) {
	t.Parallel()

	computed := true
	sensitive := true
	objects := []manifest.ManagedObject{
		{
			Name: "settings",
			Fields: []manifest.FieldSpec{
				{Name: "AUTH_BASIC_ENABLED", Type: manifest.FieldTypeBool},
			},
		},
	}

	updated := ApplyFieldOverrides(objects, map[string]FieldOverride{
		"settings.SOCIAL_AUTH_GOOGLE_OAUTH2_SECRET": {
			Object:      "settings",
			Field:       "SOCIAL_AUTH_GOOGLE_OAUTH2_SECRET",
			Type:        manifest.FieldTypeString,
			Computed:    &computed,
			Sensitive:   &sensitive,
			Description: "Google OAuth2 secret.",
		},
	})

	if len(updated) != 1 {
		t.Fatalf("unexpected object count: got=%d want=1", len(updated))
	}
	if len(updated[0].Fields) != 2 {
		t.Fatalf("unexpected field count: got=%d want=2", len(updated[0].Fields))
	}

	var secret manifest.FieldSpec
	for _, field := range updated[0].Fields {
		if field.Name == "SOCIAL_AUTH_GOOGLE_OAUTH2_SECRET" {
			secret = field
			break
		}
	}

	if secret.Name == "" {
		t.Fatalf("missing added field SOCIAL_AUTH_GOOGLE_OAUTH2_SECRET")
	}
	if !secret.Computed {
		t.Fatalf("expected added field to be computed")
	}
	if !secret.Sensitive {
		t.Fatalf("expected added field to be sensitive")
	}
	if !secret.WriteOnly {
		t.Fatalf("expected added field to infer writeOnly from sensitive=true")
	}
	if secret.Description != "Google OAuth2 secret." {
		t.Fatalf("unexpected description: got=%q", secret.Description)
	}
}

func TestApplyFieldOverridesSetsReference(t *testing.T) {
	t.Parallel()

	reference := true
	objects := []manifest.ManagedObject{
		{
			Name: "teams",
			Fields: []manifest.FieldSpec{
				{Name: "organization", Type: manifest.FieldTypeInt, Reference: false},
			},
		},
	}

	updated := ApplyFieldOverrides(objects, map[string]FieldOverride{
		"teams.organization": {
			Object:    "teams",
			Field:     "organization",
			Reference: &reference,
		},
	})

	org := updated[0].Fields[0]
	if !org.Reference {
		t.Fatalf("expected teams.organization to be marked as reference after override")
	}
}
