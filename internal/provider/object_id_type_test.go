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

func TestObjectResourceSchemaCollectionIDIsInt64(t *testing.T) {
	t.Parallel()

	r := &objectResource{object: manifest.ManagedObject{
		Name:             "teams",
		ResourceName:     "awx_team",
		CollectionCreate: true,
	}}

	resp := resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)

	idAttr, ok := resp.Schema.Attributes["id"]
	if !ok {
		t.Fatalf("expected id attribute in schema")
	}
	if _, ok := idAttr.(resourceschema.Int64Attribute); !ok {
		t.Fatalf("expected id to be Int64Attribute, got %T", idAttr)
	}
}

func TestObjectDataSourceSchemaCollectionIDIsInt64(t *testing.T) {
	t.Parallel()

	d := &objectDataSource{object: manifest.ManagedObject{
		Name:             "teams",
		DataSourceName:   "awx_team",
		CollectionCreate: true,
	}}

	resp := datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, &resp)

	idAttr, ok := resp.Schema.Attributes["id"]
	if !ok {
		t.Fatalf("expected id attribute in schema")
	}
	if _, ok := idAttr.(datasourceschema.Int64Attribute); !ok {
		t.Fatalf("expected id to be Int64Attribute, got %T", idAttr)
	}
}

func TestGetResourceIDCollectionUsesInt64State(t *testing.T) {
	t.Parallel()

	id, diags := getResourceID(context.Background(), &mockConfigSource{values: map[string]any{
		"id": types.Int64Value(42),
	}}, true)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if id != "42" {
		t.Fatalf("unexpected id: got=%q want=%q", id, "42")
	}
}
