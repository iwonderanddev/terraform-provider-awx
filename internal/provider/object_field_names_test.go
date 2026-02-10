package provider

import (
	"context"
	"testing"

	"github.com/damien/terraform-awx-provider/internal/manifest"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestObjectResourceSchemaSettingsUsesTerraformSafeAttributeNames(t *testing.T) {
	t.Parallel()

	r := &objectResource{
		object: manifest.ManagedObject{
			Name: "settings",
			Fields: []manifest.FieldSpec{
				{Name: "AUTH_LDAP_4_SERVER_URI", Type: manifest.FieldTypeString},
			},
		},
	}

	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)

	if _, ok := resp.Schema.Attributes["auth_ldap_4_server_uri"]; !ok {
		t.Fatalf("expected lowercase Terraform attribute for settings field")
	}
	if _, ok := resp.Schema.Attributes["AUTH_LDAP_4_SERVER_URI"]; ok {
		t.Fatalf("unexpected uppercase Terraform attribute for settings field")
	}
}

func TestObjectDataSourceSchemaSettingsUsesTerraformSafeAttributeNames(t *testing.T) {
	t.Parallel()

	d := &objectDataSource{
		object: manifest.ManagedObject{
			Name: "settings",
			Fields: []manifest.FieldSpec{
				{Name: "AUTH_LDAP_4_SERVER_URI", Type: manifest.FieldTypeString},
			},
		},
	}

	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, resp)

	if _, ok := resp.Schema.Attributes["auth_ldap_4_server_uri"]; !ok {
		t.Fatalf("expected lowercase Terraform attribute for settings field")
	}
	if _, ok := resp.Schema.Attributes["AUTH_LDAP_4_SERVER_URI"]; ok {
		t.Fatalf("unexpected uppercase Terraform attribute for settings field")
	}
}

func TestPayloadFromConfigSettingsMapsTerraformAttributeToAWXField(t *testing.T) {
	t.Parallel()

	r := &objectResource{
		object: manifest.ManagedObject{
			Name: "settings",
			Fields: []manifest.FieldSpec{
				{Name: "AUTH_LDAP_4_SERVER_URI", Type: manifest.FieldTypeString},
			},
		},
	}

	source := &mockConfigSource{
		values: map[string]any{
			"auth_ldap_4_server_uri": types.StringValue("ldap://ldap.example.com:389"),
		},
	}

	payload, _, diags := r.payloadFromConfig(context.Background(), source)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	if got, ok := payload["AUTH_LDAP_4_SERVER_URI"].(string); !ok || got == "" {
		t.Fatalf("expected AWX field key in payload, got: %#v", payload)
	}
}
