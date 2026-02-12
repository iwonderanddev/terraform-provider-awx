package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/damien/terraform-provider-awx-iwd/internal/manifest"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestPayloadFromConfigConvertsObjectValuesAndTracksWriteOnly(t *testing.T) {
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
			"settings":  types.DynamicValue(types.ObjectValueMust(map[string]attr.Type{"enabled": types.BoolType}, map[string]attr.Value{"enabled": types.BoolValue(true)})),
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

func TestPayloadFromConfigEncodesExtraVarsObjectForTransport(t *testing.T) {
	t.Parallel()

	tests := []string{"job_templates", "workflow_job_templates"}
	for _, objectName := range tests {
		objectName := objectName
		t.Run(objectName, func(t *testing.T) {
			t.Parallel()

			resource := &objectResource{
				object: manifest.ManagedObject{
					Name: objectName,
					Fields: []manifest.FieldSpec{
						{Name: "name", Type: manifest.FieldTypeString, Required: true},
						{Name: "extra_vars", Type: manifest.FieldTypeObject},
					},
				},
			}

			source := &mockConfigSource{
				values: map[string]any{
					"name": types.StringValue("demo"),
					"extra_vars": types.DynamicValue(types.ObjectValueMust(
						map[string]attr.Type{
							"enabled": types.BoolType,
						},
						map[string]attr.Value{
							"enabled": types.BoolValue(true),
						},
					)),
				},
			}

			payload, _, diags := resource.payloadFromConfig(context.Background(), source)
			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			raw, ok := payload["extra_vars"].(string)
			if !ok {
				t.Fatalf("expected %s.extra_vars payload to be string transport, got %T", objectName, payload["extra_vars"])
			}

			var decoded map[string]any
			if err := json.Unmarshal([]byte(raw), &decoded); err != nil {
				t.Fatalf("expected valid JSON transport string, got error: %v", err)
			}
			if got, ok := decoded["enabled"].(bool); !ok || !got {
				t.Fatalf("unexpected decoded extra_vars.enabled value: %#v", decoded["enabled"])
			}
		})
	}
}

func TestPayloadFromConfigEncodesExtraDataObjectForTransport(t *testing.T) {
	t.Parallel()

	tests := []string{"schedules", "workflow_job_template_nodes", "workflow_job_nodes"}
	for _, objectName := range tests {
		objectName := objectName
		t.Run(objectName, func(t *testing.T) {
			t.Parallel()

			resource := &objectResource{
				object: manifest.ManagedObject{
					Name: objectName,
					Fields: []manifest.FieldSpec{
						{Name: "name", Type: manifest.FieldTypeString, Required: true},
						{Name: "extra_data", Type: manifest.FieldTypeObject},
					},
				},
			}

			source := &mockConfigSource{
				values: map[string]any{
					"name": types.StringValue("demo"),
					"extra_data": types.DynamicValue(types.ObjectValueMust(
						map[string]attr.Type{
							"enabled": types.BoolType,
						},
						map[string]attr.Value{
							"enabled": types.BoolValue(true),
						},
					)),
				},
			}

			payload, _, diags := resource.payloadFromConfig(context.Background(), source)
			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			raw, ok := payload["extra_data"].(string)
			if !ok {
				t.Fatalf("expected %s.extra_data payload to be string transport, got %T", objectName, payload["extra_data"])
			}

			var decoded map[string]any
			if err := json.Unmarshal([]byte(raw), &decoded); err != nil {
				t.Fatalf("expected valid JSON transport string, got error: %v", err)
			}
			if got, ok := decoded["enabled"].(bool); !ok || !got {
				t.Fatalf("unexpected decoded extra_data.enabled value: %#v", decoded["enabled"])
			}
		})
	}
}

func TestPayloadFromConfigInvalidObjectReturnsDiagnostic(t *testing.T) {
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
			"settings": types.DynamicValue(types.StringValue("not-an-object")),
		},
	}

	_, _, diags := resource.payloadFromConfig(context.Background(), source)
	if !diags.HasError() {
		t.Fatalf("expected invalid object payload to return diagnostics")
	}
}

func TestPayloadFromConfigAllowsNullObjectValues(t *testing.T) {
	t.Parallel()

	resource := &objectResource{
		object: manifest.ManagedObject{
			Fields: []manifest.FieldSpec{
				{Name: "extra_vars", Type: manifest.FieldTypeObject},
			},
		},
	}

	source := &mockConfigSource{
		values: map[string]any{
			"extra_vars": types.DynamicValue(types.MapNull(types.DynamicType)),
		},
	}

	payload, _, diags := resource.payloadFromConfig(context.Background(), source)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics for null object value: %v", diags)
	}
	if _, exists := payload["extra_vars"]; exists {
		t.Fatalf("expected null object value to be omitted from payload")
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
				{Name: "config", Type: manifest.FieldTypeObject, WriteOnly: true},
			},
		},
	}

	source := &mockConfigSource{
		values: map[string]any{
			"name":     types.StringValue("demo"),
			"token":    types.StringValue("token-value"),
			"password": types.StringUnknown(),
			"team":     types.Int64Value(22),
			"config": types.DynamicValue(types.ObjectValueMust(
				map[string]attr.Type{
					"threshold": types.NumberType,
				},
				map[string]attr.Value{
					"threshold": types.NumberValue(big.NewFloat(2.5)),
				},
			)),
		},
	}

	snapshot := resource.writeOnlyValuesFromSource(context.Background(), source)
	if snapshot.HasError() {
		t.Fatalf("unexpected diagnostics: %v", snapshot.Diagnostics)
	}
	if len(snapshot.Values) != 3 {
		t.Fatalf("unexpected write-only snapshot count: got=%d want=3", len(snapshot.Values))
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
	configValue, ok := snapshot.Values["config"].(types.Dynamic)
	if !ok {
		t.Fatalf("expected config write-only value to be types.Dynamic, got %T", snapshot.Values["config"])
	}
	if configValue.IsNull() {
		t.Fatalf("expected config write-only value to be non-null")
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

func TestPruneUnchangedFieldsFromPayloadRemovesUnchangedStateValues(t *testing.T) {
	t.Parallel()

	resource := &objectResource{
		object: manifest.ManagedObject{
			Name: "instance_groups",
			Fields: []manifest.FieldSpec{
				{Name: "name", Type: manifest.FieldTypeString, Required: true},
				{Name: "pod_spec_override", Type: manifest.FieldTypeString},
				{Name: "policy_instance_percentage", Type: manifest.FieldTypeInt, Computed: true},
			},
		},
	}

	planSource := &mockConfigSource{
		values: map[string]any{
			"name":                       types.StringValue("default"),
			"pod_spec_override":          types.StringValue("new-spec"),
			"policy_instance_percentage": types.Int64Value(100),
		},
	}
	stateSource := &mockConfigSource{
		values: map[string]any{
			"name":                       types.StringValue("default"),
			"pod_spec_override":          types.StringValue("old-spec"),
			"policy_instance_percentage": types.Int64Value(100),
		},
	}

	payload, _, diags := resource.payloadFromConfig(context.Background(), planSource)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics building payload: %v", diags)
	}

	diags = resource.pruneUnchangedFieldsFromPayload(context.Background(), payload, planSource, stateSource)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics pruning payload: %v", diags)
	}

	if got, ok := payload["pod_spec_override"]; !ok || got != "new-spec" {
		t.Fatalf("expected pod_spec_override to be retained in payload, got %#v", payload["pod_spec_override"])
	}
	if _, exists := payload["name"]; exists {
		t.Fatalf("expected unchanged name to be pruned from payload")
	}
	if _, exists := payload["policy_instance_percentage"]; exists {
		t.Fatalf("expected unchanged computed field to be pruned from payload")
	}
}

func TestPruneUnchangedFieldsFromPayloadTreatsEquivalentArraysAsUnchanged(t *testing.T) {
	t.Parallel()

	resource := &objectResource{
		object: manifest.ManagedObject{
			Fields: []manifest.FieldSpec{
				{Name: "policy_instance_list", Type: manifest.FieldTypeArray},
			},
		},
	}

	planSource := &mockConfigSource{
		values: map[string]any{
			"policy_instance_list": types.StringValue(`[1,2,3]`),
		},
	}
	stateSource := &mockConfigSource{
		values: map[string]any{
			"policy_instance_list": types.StringValue(`[1, 2, 3]`),
		},
	}

	payload, _, diags := resource.payloadFromConfig(context.Background(), planSource)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics building payload: %v", diags)
	}

	diags = resource.pruneUnchangedFieldsFromPayload(context.Background(), payload, planSource, stateSource)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics pruning payload: %v", diags)
	}
	if _, exists := payload["policy_instance_list"]; exists {
		t.Fatalf("expected equivalent JSON arrays to be pruned from payload")
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
	case *types.Dynamic:
		v, ok := value.(types.Dynamic)
		if !ok {
			diags.AddError("mock type mismatch", fmt.Sprintf("expected types.Dynamic, got %T", value))
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
	case *types.Dynamic:
		*t = types.DynamicNull()
	default:
		diags := diag.Diagnostics{}
		diags.AddError("unsupported mock target", fmt.Sprintf("target type %T is not supported", target))
		return diags
	}
	return nil
}
