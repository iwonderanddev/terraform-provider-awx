package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	awxclient "github.com/damien/terraform-awx-provider/internal/client"
	"github.com/damien/terraform-awx-provider/internal/manifest"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func TestNewResourceFieldAttributeSensitive(t *testing.T) {
	t.Parallel()

	attribute := newResourceFieldAttribute(manifest.FieldSpec{
		Name:      "password",
		Type:      manifest.FieldTypeString,
		Sensitive: true,
	}, true)

	stringAttribute, ok := attribute.(resourceschema.StringAttribute)
	if !ok {
		t.Fatalf("expected string attribute, got %T", attribute)
	}
	if !stringAttribute.Sensitive {
		t.Fatalf("expected attribute to be sensitive")
	}
}

func TestNewResourceFieldAttributeComputed(t *testing.T) {
	t.Parallel()

	attribute := newResourceFieldAttribute(manifest.FieldSpec{
		Name:     "max_hosts",
		Type:     manifest.FieldTypeInt,
		Required: false,
		Computed: true,
	}, true)

	intAttribute, ok := attribute.(resourceschema.Int64Attribute)
	if !ok {
		t.Fatalf("expected int64 attribute, got %T", attribute)
	}
	if !intAttribute.Optional {
		t.Fatalf("expected computed field to remain optional")
	}
	if !intAttribute.Computed {
		t.Fatalf("expected attribute to be computed")
	}
}

func TestNewResourceFieldAttributeRequiresReplaceWhenUpdateUnsupported(t *testing.T) {
	t.Parallel()

	attribute := newResourceFieldAttribute(manifest.FieldSpec{
		Name: "user",
		Type: manifest.FieldTypeInt,
	}, false)

	intAttribute, ok := attribute.(resourceschema.Int64Attribute)
	if !ok {
		t.Fatalf("expected int64 attribute, got %T", attribute)
	}
	if len(intAttribute.PlanModifiers) == 0 {
		t.Fatalf("expected RequiresReplace plan modifier when update is unsupported")
	}
}

func TestParseNumericID(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		input   any
		want    int64
		wantErr bool
	}{
		{name: "string", input: "42", want: 42},
		{name: "int", input: int(7), want: 7},
		{name: "float", input: float64(9), want: 9},
		{name: "json number", input: json.Number("13"), want: 13},
		{name: "unsupported type", input: true, wantErr: true},
		{name: "invalid string", input: "abc", wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseNumericID(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected parseNumericID error")
				}
				return
			}
			if err != nil {
				t.Fatalf("parseNumericID returned error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("unexpected parsed value: got=%d want=%d", got, tc.want)
			}
		})
	}
}

func TestCompositeIDPattern(t *testing.T) {
	t.Parallel()

	valid := "12:34"
	if !compositeIDPattern.MatchString(valid) {
		t.Fatalf("expected %q to match composite id pattern", valid)
	}

	invalid := "12-34"
	if compositeIDPattern.MatchString(invalid) {
		t.Fatalf("expected %q to fail composite id pattern", invalid)
	}
}

func TestSetRelationshipStateCompositeID(t *testing.T) {
	t.Parallel()

	target := &mockAttributeTarget{values: map[string]any{}}
	diags := diag.Diagnostics{}

	setRelationshipState(context.Background(), 12, 34, target, &diags)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	foundCompositeID := false
	for _, value := range target.values {
		if fmt.Sprintf("%v", value) == "12:34" {
			foundCompositeID = true
			break
		}
	}
	if !foundCompositeID {
		t.Fatalf("expected composite id 12:34 to be written to state")
	}
}

func TestProviderRegistersAllManagedObjectResources(t *testing.T) {
	t.Parallel()

	providerInstance := New("test")().(*awxProvider) //nolint:forcetypeassert
	resourceFactories := providerInstance.Resources(context.Background())

	expected := make(map[string]struct{})
	for _, object := range providerInstance.catalog.ManagedResourceObjects() {
		expected[object.ResourceName] = struct{}{}
	}

	actual := make(map[string]struct{})
	for _, factory := range resourceFactories {
		resourceInstance := factory()
		if _, ok := resourceInstance.(*objectResource); !ok {
			continue
		}
		var metadata resource.MetadataResponse
		resourceInstance.Metadata(context.Background(), resource.MetadataRequest{}, &metadata)
		actual[metadata.TypeName] = struct{}{}
	}

	if len(actual) != len(expected) {
		t.Fatalf("unexpected managed resource registration count: got=%d want=%d", len(actual), len(expected))
	}
	for resourceName := range expected {
		if _, ok := actual[resourceName]; !ok {
			t.Fatalf("missing managed resource registration for %s", resourceName)
		}
	}
}

func TestProviderRegistersAllManagedObjectDataSources(t *testing.T) {
	t.Parallel()

	providerInstance := New("test")().(*awxProvider) //nolint:forcetypeassert
	dataSourceFactories := providerInstance.DataSources(context.Background())

	expected := make(map[string]struct{})
	for _, object := range providerInstance.catalog.ManagedDataSourceObjects() {
		expected[object.DataSourceName] = struct{}{}
	}

	actual := make(map[string]struct{})
	for _, factory := range dataSourceFactories {
		dataSourceInstance := factory()
		if _, ok := dataSourceInstance.(*objectDataSource); !ok {
			continue
		}
		var metadata datasource.MetadataResponse
		dataSourceInstance.Metadata(context.Background(), datasource.MetadataRequest{}, &metadata)
		actual[metadata.TypeName] = struct{}{}
	}

	if len(actual) != len(expected) {
		t.Fatalf("unexpected managed data source registration count: got=%d want=%d", len(actual), len(expected))
	}
	for dataSourceName := range expected {
		if _, ok := actual[dataSourceName]; !ok {
			t.Fatalf("missing managed data source registration for %s", dataSourceName)
		}
	}
}

func TestShouldRemoveFromStateOnReadError(t *testing.T) {
	t.Parallel()

	if !shouldRemoveFromStateOnReadError(&awxclient.APIError{StatusCode: http.StatusNotFound}) {
		t.Fatalf("expected not-found APIError to trigger state removal")
	}
	if shouldRemoveFromStateOnReadError(&awxclient.APIError{StatusCode: http.StatusInternalServerError}) {
		t.Fatalf("expected non-404 APIError to keep state")
	}
	if shouldRemoveFromStateOnReadError(nil) {
		t.Fatalf("expected nil error to keep state")
	}
}

func TestObjectResourceSchemaHasNoNestedLifecycleBlocks(t *testing.T) {
	t.Parallel()

	resourceInstance := NewObjectResource(manifest.ManagedObject{
		Name:             "settings",
		ResourceName:     "awx_setting",
		CollectionCreate: false,
		UpdateSupported:  true,
		Fields: []manifest.FieldSpec{
			{Name: "tower_url_base", Type: manifest.FieldTypeString},
			{Name: "ansible_callbacks_enabled", Type: manifest.FieldTypeBool},
		},
	})

	resp := resource.SchemaResponse{}
	resourceInstance.Schema(context.Background(), resource.SchemaRequest{}, &resp)

	if len(resp.Schema.Blocks) != 0 {
		t.Fatalf("expected object resources to expose only attributes, found %d nested block(s)", len(resp.Schema.Blocks))
	}
}

type mockAttributeTarget struct {
	values map[string]any
}

func (m *mockAttributeTarget) SetAttribute(_ context.Context, p path.Path, value any) diag.Diagnostics {
	if m.values == nil {
		m.values = map[string]any{}
	}
	m.values[fmt.Sprintf("%v", p)] = value
	return nil
}
