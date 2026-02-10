package provider

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/damien/terraform-awx-provider/internal/manifest"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestPayloadFromConfigDecodesJSONAndTracksWriteOnly(t *testing.T) {
	t.Parallel()

	resource := &objectResource{
		object: manifest.ManagedObject{
			Fields: []manifest.FieldSpec{
				{Name: "name", Type: manifest.FieldTypeString, Required: true},
				{Name: "settings", Type: manifest.FieldTypeObject},
				{Name: "tags", Type: manifest.FieldTypeArray},
				{Name: "max_hosts", Type: manifest.FieldTypeInt},
				{Name: "token", Type: manifest.FieldTypeString, WriteOnly: true},
			},
		},
	}

	source := &mockConfigSource{
		values: map[string]any{
			"name":      types.StringValue("demo"),
			"settings":  types.StringValue(`{"enabled":true}`),
			"tags":      types.StringValue(`["a","b"]`),
			"max_hosts": types.Int64Null(),
			"token":     types.StringValue("secret-value"),
		},
	}

	payload, writeOnlyValues, diags := resource.payloadFromConfig(context.Background(), source)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	if got, ok := payload["name"].(string); !ok || got != "demo" {
		t.Fatalf("unexpected payload name: %#v", payload["name"])
	}

	settings, ok := payload["settings"].(map[string]any)
	if !ok {
		t.Fatalf("expected settings object payload, got %T", payload["settings"])
	}
	if got, ok := settings["enabled"].(bool); !ok || !got {
		t.Fatalf("unexpected settings.enabled value: %#v", settings["enabled"])
	}

	tags, ok := payload["tags"].([]any)
	if !ok {
		t.Fatalf("expected tags array payload, got %T", payload["tags"])
	}
	if len(tags) != 2 {
		t.Fatalf("unexpected tags length: got=%d want=2", len(tags))
	}

	if _, exists := payload["max_hosts"]; exists {
		t.Fatalf("expected max_hosts to be omitted from payload when null")
	}

	if got, ok := payload["token"].(string); !ok || got != "secret-value" {
		t.Fatalf("unexpected token payload value: %#v", payload["token"])
	}
	tokenValue, ok := writeOnlyValues["token"].(types.String)
	if !ok || tokenValue.ValueString() != "secret-value" {
		t.Fatalf("expected token write-only value to be preserved")
	}
}

func TestPayloadFromConfigInvalidJSONReturnsDiagnostic(t *testing.T) {
	t.Parallel()

	resource := &objectResource{
		object: manifest.ManagedObject{
			Fields: []manifest.FieldSpec{
				{Name: "settings", Type: manifest.FieldTypeObject},
			},
		},
	}

	source := &mockConfigSource{
		values: map[string]any{
			"settings": types.StringValue("{"),
		},
	}

	_, _, diags := resource.payloadFromConfig(context.Background(), source)
	if !diags.HasError() {
		t.Fatalf("expected invalid JSON to return diagnostics")
	}
}

func TestWriteOnlyValuesFromSourcePreservesOnlyKnownWriteOnlyValues(t *testing.T) {
	t.Parallel()

	resource := &objectResource{
		object: manifest.ManagedObject{
			Fields: []manifest.FieldSpec{
				{Name: "name", Type: manifest.FieldTypeString},
				{Name: "token", Type: manifest.FieldTypeString, WriteOnly: true},
				{Name: "password", Type: manifest.FieldTypeString, WriteOnly: true},
				{Name: "team", Type: manifest.FieldTypeInt, WriteOnly: true},
			},
		},
	}

	source := &mockConfigSource{
		values: map[string]any{
			"name":     types.StringValue("demo"),
			"token":    types.StringValue("token-value"),
			"password": types.StringUnknown(),
			"team":     types.Int64Value(22),
		},
	}

	snapshot := resource.writeOnlyValuesFromSource(context.Background(), source)
	if snapshot.HasError() {
		t.Fatalf("unexpected diagnostics: %v", snapshot.Diagnostics)
	}
	if len(snapshot.Values) != 2 {
		t.Fatalf("unexpected write-only snapshot count: got=%d want=2", len(snapshot.Values))
	}
	tokenValue, ok := snapshot.Values["token"].(types.String)
	if !ok {
		t.Fatalf("expected token write-only value to be types.String, got %T", snapshot.Values["token"])
	}
	if got := tokenValue.ValueString(); got != "token-value" {
		t.Fatalf("unexpected write-only snapshot token value: got=%q want=%q", got, "token-value")
	}
	teamValue, ok := snapshot.Values["team"].(types.Int64)
	if !ok {
		t.Fatalf("expected team write-only value to be types.Int64, got %T", snapshot.Values["team"])
	}
	if got := teamValue.ValueInt64(); got != 22 {
		t.Fatalf("unexpected write-only snapshot team value: got=%d want=%d", got, 22)
	}
	if _, exists := snapshot.Values["password"]; exists {
		t.Fatalf("expected unknown write-only value to be skipped")
	}
}

func TestStringValuesFromSourceSkipsUnknownValues(t *testing.T) {
	t.Parallel()

	resource := &objectResource{
		object: manifest.ManagedObject{
			Fields: []manifest.FieldSpec{
				{Name: "name", Type: manifest.FieldTypeString},
				{Name: "description", Type: manifest.FieldTypeString},
				{Name: "max_hosts", Type: manifest.FieldTypeInt},
			},
		},
	}

	source := &mockConfigSource{
		values: map[string]any{
			"name":        types.StringUnknown(),
			"description": types.StringNull(),
			"max_hosts":   types.Int64Value(12),
		},
	}

	snapshot := resource.stringValuesFromSource(context.Background(), source)
	if snapshot.HasError() {
		t.Fatalf("unexpected diagnostics: %v", snapshot.Diagnostics)
	}
	if len(snapshot.Values) != 1 {
		t.Fatalf("unexpected string snapshot count: got=%d want=1", len(snapshot.Values))
	}
	description, ok := snapshot.Values["description"]
	if !ok {
		t.Fatalf("expected description in snapshot")
	}
	if !description.IsNull() {
		t.Fatalf("expected description snapshot to preserve null state")
	}
	if _, exists := snapshot.Values["name"]; exists {
		t.Fatalf("expected unknown name to be skipped")
	}
}

type mockConfigSource struct {
	values map[string]any
}

func (m *mockConfigSource) GetAttribute(_ context.Context, p path.Path, target any) diag.Diagnostics {
	for key, value := range m.values {
		if strings.Contains(fmt.Sprintf("%v", p), key) {
			return assignMockAttribute(target, value)
		}
	}
	return assignMockNull(target)
}

func assignMockAttribute(target any, value any) diag.Diagnostics {
	diags := diag.Diagnostics{}
	switch t := target.(type) {
	case *types.String:
		v, ok := value.(types.String)
		if !ok {
			diags.AddError("mock type mismatch", fmt.Sprintf("expected types.String, got %T", value))
			return diags
		}
		*t = v
	case *types.Int64:
		v, ok := value.(types.Int64)
		if !ok {
			diags.AddError("mock type mismatch", fmt.Sprintf("expected types.Int64, got %T", value))
			return diags
		}
		*t = v
	case *types.Bool:
		v, ok := value.(types.Bool)
		if !ok {
			diags.AddError("mock type mismatch", fmt.Sprintf("expected types.Bool, got %T", value))
			return diags
		}
		*t = v
	case *types.Float64:
		v, ok := value.(types.Float64)
		if !ok {
			diags.AddError("mock type mismatch", fmt.Sprintf("expected types.Float64, got %T", value))
			return diags
		}
		*t = v
	default:
		diags.AddError("unsupported mock target", fmt.Sprintf("target type %T is not supported", target))
	}
	return diags
}

func assignMockNull(target any) diag.Diagnostics {
	switch t := target.(type) {
	case *types.String:
		*t = types.StringNull()
	case *types.Int64:
		*t = types.Int64Null()
	case *types.Bool:
		*t = types.BoolNull()
	case *types.Float64:
		*t = types.Float64Null()
	default:
		diags := diag.Diagnostics{}
		diags.AddError("unsupported mock target", fmt.Sprintf("target type %T is not supported", target))
		return diags
	}
	return nil
}
