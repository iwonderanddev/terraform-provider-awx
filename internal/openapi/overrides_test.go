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
