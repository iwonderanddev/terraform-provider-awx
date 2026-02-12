package provider

import (
	"context"
	"testing"

	"github.com/damien/terraform-provider-awx-iwd/internal/manifest"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
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

func TestObjectResourceSchemaUsesDynamicAttributeForObjectFields(t *testing.T) {
	t.Parallel()

	r := &objectResource{
		object: manifest.ManagedObject{
			Name: "settings",
			Fields: []manifest.FieldSpec{
				{Name: "AUTH_LDAP_1_CONNECTION_OPTIONS", Type: manifest.FieldTypeObject},
			},
		},
	}

	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)

	attribute, ok := resp.Schema.Attributes["auth_ldap_1_connection_options"]
	if !ok {
		t.Fatalf("expected object field attribute")
	}
	if _, ok := attribute.(resourceschema.DynamicAttribute); !ok {
		t.Fatalf("expected resourceschema.DynamicAttribute, got %T", attribute)
	}
}

func TestObjectDataSourceSchemaUsesDynamicAttributeForObjectFields(t *testing.T) {
	t.Parallel()

	d := &objectDataSource{
		object: manifest.ManagedObject{
			Name: "settings",
			Fields: []manifest.FieldSpec{
				{Name: "AUTH_LDAP_1_CONNECTION_OPTIONS", Type: manifest.FieldTypeObject},
			},
		},
	}

	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, resp)

	attribute, ok := resp.Schema.Attributes["auth_ldap_1_connection_options"]
	if !ok {
		t.Fatalf("expected object field attribute")
	}
	if _, ok := attribute.(datasourceschema.DynamicAttribute); !ok {
		t.Fatalf("expected datasourceschema.DynamicAttribute, got %T", attribute)
	}
}

func TestObjectResourceSchemaReferenceFieldUsesIDSuffix(t *testing.T) {
	t.Parallel()

	r := &objectResource{
		object: manifest.ManagedObject{
			Name: "teams",
			Fields: []manifest.FieldSpec{
				{Name: "organization", Type: manifest.FieldTypeInt, Reference: true},
			},
		},
	}

	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)

	if _, ok := resp.Schema.Attributes["organization_id"]; !ok {
		t.Fatalf("expected suffixed Terraform attribute for reference field")
	}
	if _, ok := resp.Schema.Attributes["organization"]; ok {
		t.Fatalf("unexpected unsuffixed Terraform attribute for reference field")
	}
}

func TestObjectDataSourceSchemaReferenceFieldUsesIDSuffix(t *testing.T) {
	t.Parallel()

	d := &objectDataSource{
		object: manifest.ManagedObject{
			Name: "teams",
			Fields: []manifest.FieldSpec{
				{Name: "organization", Type: manifest.FieldTypeInt, Reference: true},
			},
		},
	}

	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, resp)

	if _, ok := resp.Schema.Attributes["organization_id"]; !ok {
		t.Fatalf("expected suffixed Terraform attribute for reference field")
	}
	if _, ok := resp.Schema.Attributes["organization"]; ok {
		t.Fatalf("unexpected unsuffixed Terraform attribute for reference field")
	}
}

func TestPayloadFromConfigReferenceFieldMapsSuffixedTerraformNameToAWXField(t *testing.T) {
	t.Parallel()

	r := &objectResource{
		object: manifest.ManagedObject{
			Name: "teams",
			Fields: []manifest.FieldSpec{
				{Name: "name", Type: manifest.FieldTypeString, Required: true},
				{Name: "organization", Type: manifest.FieldTypeInt, Reference: true, Required: true},
			},
		},
	}

	source := &mockConfigSource{
		values: map[string]any{
			"name":            types.StringValue("platform"),
			"organization_id": types.Int64Value(9),
		},
	}

	payload, _, diags := r.payloadFromConfig(context.Background(), source)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	if _, exists := payload["organization_id"]; exists {
		t.Fatalf("unexpected payload key organization_id")
	}
	if got, ok := payload["organization"].(int64); !ok || got != 9 {
		t.Fatalf("expected AWX field organization to be populated, got=%#v", payload["organization"])
	}
}
