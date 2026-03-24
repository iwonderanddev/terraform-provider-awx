package openapi

import (
	"testing"

	"github.com/damien/terraform-provider-awx-iwd/internal/manifest"
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

func TestApplyFieldOverridesSetsReadOnly(t *testing.T) {
	t.Parallel()

	computed := true
	readOnly := true
	objects := []manifest.ManagedObject{
		{
			Name: "projects",
			Fields: []manifest.FieldSpec{
				{Name: "name", Type: manifest.FieldTypeString, Required: true},
				{Name: "local_path", Type: manifest.FieldTypeString, Required: false, Computed: false},
			},
		},
	}

	updated := ApplyFieldOverrides(objects, map[string]FieldOverride{
		"projects.local_path": {
			Object:   "projects",
			Field:    "local_path",
			Computed: &computed,
			ReadOnly: &readOnly,
		},
	})

	var localPath manifest.FieldSpec
	for _, field := range updated[0].Fields {
		if field.Name == "local_path" {
			localPath = field
			break
		}
	}

	if !localPath.Computed {
		t.Fatalf("expected local_path to be computed after override")
	}
	if !localPath.ReadOnly {
		t.Fatalf("expected local_path to be read-only after override")
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

func TestApplyFieldOverridesSetsTerraformName(t *testing.T) {
	t.Parallel()

	objects := []manifest.ManagedObject{
		{
			Name: "projects",
			Fields: []manifest.FieldSpec{
				{Name: "credential", Type: manifest.FieldTypeInt, Reference: true},
			},
		},
	}

	updated := ApplyFieldOverrides(objects, map[string]FieldOverride{
		"projects.credential": {
			Object:        "projects",
			Field:         "credential",
			TerraformName: "scm_credential_id",
		},
	})

	credential := updated[0].Fields[0]
	if credential.TerraformName != "scm_credential_id" {
		t.Fatalf("unexpected Terraform name override: got=%q want=%q", credential.TerraformName, "scm_credential_id")
	}
}
