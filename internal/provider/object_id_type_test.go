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

func TestRoleAssignmentObjectIDResourceSchemasUseInt64(t *testing.T) {
	t.Parallel()

	catalog := manifest.MustLoad()
	objectNames := []string{"role_team_assignments", "role_user_assignments"}

	for _, objectName := range objectNames {
		object, ok := catalog.ObjectByName(objectName)
		if !ok {
			t.Fatalf("expected %s object in catalog", objectName)
		}

		r := &objectResource{object: object}
		resp := resource.SchemaResponse{}
		r.Schema(context.Background(), resource.SchemaRequest{}, &resp)

		attr, ok := resp.Schema.Attributes["object_id"]
		if !ok {
			t.Fatalf("expected object_id attribute for %s resource schema", objectName)
		}
		if _, ok := attr.(resourceschema.Int64Attribute); !ok {
			t.Fatalf("expected %s.resource object_id to be Int64Attribute, got %T", objectName, attr)
		}
	}
}

func TestRoleAssignmentObjectIDDataSourceSchemasUseInt64(t *testing.T) {
	t.Parallel()

	catalog := manifest.MustLoad()
	objectNames := []string{"role_team_assignments", "role_user_assignments"}

	for _, objectName := range objectNames {
		object, ok := catalog.ObjectByName(objectName)
		if !ok {
			t.Fatalf("expected %s object in catalog", objectName)
		}

		d := &objectDataSource{object: object}
		resp := datasource.SchemaResponse{}
		d.Schema(context.Background(), datasource.SchemaRequest{}, &resp)

		attr, ok := resp.Schema.Attributes["object_id"]
		if !ok {
			t.Fatalf("expected object_id attribute for %s data source schema", objectName)
		}
		if _, ok := attr.(datasourceschema.Int64Attribute); !ok {
			t.Fatalf("expected %s.data-source object_id to be Int64Attribute, got %T", objectName, attr)
		}
	}
}
