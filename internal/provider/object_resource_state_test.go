package provider

import (
	"context"
	"strings"
	"testing"

	"github.com/damien/terraform-provider-awx-iwd/internal/manifest"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestNormalizeOptionalEmptyStringToNull(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		field    manifest.FieldSpec
		value    any
		prior    map[string]types.String
		wantOK   bool
		wantNull bool
	}{
		{
			name: "optional empty string with prior null normalizes",
			field: manifest.FieldSpec{
				Name:     "description",
				Type:     manifest.FieldTypeString,
				Required: false,
			},
			value:    "",
			prior:    map[string]types.String{"description": types.StringNull()},
			wantOK:   true,
			wantNull: true,
		},
		{
			name: "optional empty string with explicit prior empty does not normalize",
			field: manifest.FieldSpec{
				Name:     "description",
				Type:     manifest.FieldTypeString,
				Required: false,
			},
			value:  "",
			prior:  map[string]types.String{"description": types.StringValue("")},
			wantOK: false,
		},
		{
			name: "optional non-empty string does not normalize",
			field: manifest.FieldSpec{
				Name:     "description",
				Type:     manifest.FieldTypeString,
				Required: false,
			},
			value:  "x",
			prior:  map[string]types.String{"description": types.StringNull()},
			wantOK: false,
		},
		{
			name: "required empty string does not normalize",
			field: manifest.FieldSpec{
				Name:     "name",
				Type:     manifest.FieldTypeString,
				Required: true,
			},
			value:  "",
			prior:  map[string]types.String{"name": types.StringNull()},
			wantOK: false,
		},
		{
			name: "missing prior value does not normalize",
			field: manifest.FieldSpec{
				Name:     "description",
				Type:     manifest.FieldTypeString,
				Required: false,
			},
			value:  "",
			prior:  map[string]types.String{},
			wantOK: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := normalizeOptionalEmptyStringToNull(tc.field, tc.value, tc.prior)
			if ok != tc.wantOK {
				t.Fatalf("unexpected normalize result: got=%v want=%v", ok, tc.wantOK)
			}
			if !ok {
				return
			}
			if got.IsNull() != tc.wantNull {
				t.Fatalf("unexpected normalized nullness: got=%v want=%v", got.IsNull(), tc.wantNull)
			}
		})
	}
}

func TestPreserveKnownNormalizedStringField(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		objectName string
		field      manifest.FieldSpec
		value      any
		prior      map[string]types.String
		wantOK     bool
		wantValue  string
	}{
		{
			name:       "instance group pod spec preserves trailing newline delta",
			objectName: "instance_groups",
			field: manifest.FieldSpec{
				Name: "pod_spec_override",
				Type: manifest.FieldTypeString,
			},
			value:     "kind: Pod",
			prior:     map[string]types.String{"pod_spec_override": types.StringValue("kind: Pod\n")},
			wantOK:    true,
			wantValue: "kind: Pod\n",
		},
		{
			name:       "instance group pod spec preserves crlf trailing newline delta",
			objectName: "instance_groups",
			field: manifest.FieldSpec{
				Name: "pod_spec_override",
				Type: manifest.FieldTypeString,
			},
			value:     "kind: Pod",
			prior:     map[string]types.String{"pod_spec_override": types.StringValue("kind: Pod\r\n")},
			wantOK:    true,
			wantValue: "kind: Pod\r\n",
		},
		{
			name:       "non-target object field does not preserve",
			objectName: "instance_groups",
			field: manifest.FieldSpec{
				Name: "description",
				Type: manifest.FieldTypeString,
			},
			value:  "kind: Pod",
			prior:  map[string]types.String{"description": types.StringValue("kind: Pod\n")},
			wantOK: false,
		},
		{
			name:       "target field on different object does not preserve",
			objectName: "inventories",
			field: manifest.FieldSpec{
				Name: "pod_spec_override",
				Type: manifest.FieldTypeString,
			},
			value:  "kind: Pod",
			prior:  map[string]types.String{"pod_spec_override": types.StringValue("kind: Pod\n")},
			wantOK: false,
		},
		{
			name:       "different content does not preserve",
			objectName: "instance_groups",
			field: manifest.FieldSpec{
				Name: "pod_spec_override",
				Type: manifest.FieldTypeString,
			},
			value:  "kind: Job",
			prior:  map[string]types.String{"pod_spec_override": types.StringValue("kind: Pod\n")},
			wantOK: false,
		},
		{
			name:       "null prior does not preserve",
			objectName: "instance_groups",
			field: manifest.FieldSpec{
				Name: "pod_spec_override",
				Type: manifest.FieldTypeString,
			},
			value:  "kind: Pod",
			prior:  map[string]types.String{"pod_spec_override": types.StringNull()},
			wantOK: false,
		},
		{
			name:       "non-string api value does not preserve",
			objectName: "instance_groups",
			field: manifest.FieldSpec{
				Name: "pod_spec_override",
				Type: manifest.FieldTypeString,
			},
			value:  map[string]any{"kind": "Pod"},
			prior:  map[string]types.String{"pod_spec_override": types.StringValue("kind: Pod\n")},
			wantOK: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, ok := preserveKnownNormalizedStringField(tc.objectName, tc.field, tc.value, tc.prior)
			if ok != tc.wantOK {
				t.Fatalf("unexpected preserve result: got=%v want=%v", ok, tc.wantOK)
			}
			if !ok {
				return
			}
			if got.IsNull() {
				t.Fatalf("expected non-null preserved value")
			}
			if got.ValueString() != tc.wantValue {
				t.Fatalf("unexpected preserved value: got=%q want=%q", got.ValueString(), tc.wantValue)
			}
		})
	}
}

func TestSetStateOptionalStringPreservesNullForEmptyAPIValue(t *testing.T) {
	t.Parallel()

	resource := &objectResource{
		object: manifest.ManagedObject{
			Fields: []manifest.FieldSpec{
				{
					Name:     "description",
					Type:     manifest.FieldTypeString,
					Required: false,
				},
			},
		},
	}

	target := &mockAttributeTarget{values: map[string]any{}}
	diags := resource.setState(
		context.Background(),
		target,
		"42",
		map[string]any{"description": ""},
		nil,
		map[string]types.String{"description": types.StringNull()},
		nil,
	)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	value, ok := findStateAttributeValue(target.values, "description")
	if !ok {
		t.Fatalf("description attribute was not written")
	}

	stringValue, ok := value.(types.String)
	if !ok {
		t.Fatalf("expected types.String for description, got %T", value)
	}
	if !stringValue.IsNull() {
		t.Fatalf("expected null description in state when prior plan/state was null")
	}
}

func TestSetStateOptionalStringKeepsExplicitEmptyString(t *testing.T) {
	t.Parallel()

	resource := &objectResource{
		object: manifest.ManagedObject{
			Fields: []manifest.FieldSpec{
				{
					Name:     "description",
					Type:     manifest.FieldTypeString,
					Required: false,
				},
			},
		},
	}

	target := &mockAttributeTarget{values: map[string]any{}}
	diags := resource.setState(
		context.Background(),
		target,
		"42",
		map[string]any{"description": ""},
		nil,
		map[string]types.String{"description": types.StringValue("")},
		nil,
	)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	value, ok := findStateAttributeValue(target.values, "description")
	if !ok {
		t.Fatalf("description attribute was not written")
	}

	stringValue, ok := value.(types.String)
	if !ok {
		t.Fatalf("expected types.String for description, got %T", value)
	}
	if stringValue.IsNull() {
		t.Fatalf("expected explicit empty string to remain non-null")
	}
	if got := stringValue.ValueString(); got != "" {
		t.Fatalf("unexpected description value: got=%q want=\"\"", got)
	}
}

func TestSetStateInstanceGroupPodSpecOverridePreservesTrailingNewlineEquivalentValue(t *testing.T) {
	t.Parallel()

	resource := &objectResource{
		object: manifest.ManagedObject{
			Name: "instance_groups",
			Fields: []manifest.FieldSpec{
				{
					Name:     "pod_spec_override",
					Type:     manifest.FieldTypeString,
					Required: false,
				},
			},
		},
	}

	target := &mockAttributeTarget{values: map[string]any{}}
	diags := resource.setState(
		context.Background(),
		target,
		"42",
		map[string]any{"pod_spec_override": "kind: Pod"},
		nil,
		map[string]types.String{"pod_spec_override": types.StringValue("kind: Pod\n")},
		nil,
	)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	value, ok := findStateAttributeValue(target.values, "pod_spec_override")
	if !ok {
		t.Fatalf("pod_spec_override attribute was not written")
	}

	stringValue, ok := value.(types.String)
	if !ok {
		t.Fatalf("expected types.String for pod_spec_override, got %T", value)
	}
	if got := stringValue.ValueString(); got != "kind: Pod\n" {
		t.Fatalf("unexpected pod_spec_override value: got=%q want=%q", got, "kind: Pod\n")
	}
}

func TestSetStateInstanceGroupPodSpecOverrideUsesAPIValueWhenContentChanges(t *testing.T) {
	t.Parallel()

	resource := &objectResource{
		object: manifest.ManagedObject{
			Name: "instance_groups",
			Fields: []manifest.FieldSpec{
				{
					Name:     "pod_spec_override",
					Type:     manifest.FieldTypeString,
					Required: false,
				},
			},
		},
	}

	target := &mockAttributeTarget{values: map[string]any{}}
	diags := resource.setState(
		context.Background(),
		target,
		"42",
		map[string]any{"pod_spec_override": "kind: Job"},
		nil,
		map[string]types.String{"pod_spec_override": types.StringValue("kind: Pod\n")},
		nil,
	)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	value, ok := findStateAttributeValue(target.values, "pod_spec_override")
	if !ok {
		t.Fatalf("pod_spec_override attribute was not written")
	}

	stringValue, ok := value.(types.String)
	if !ok {
		t.Fatalf("expected types.String for pod_spec_override, got %T", value)
	}
	if got := stringValue.ValueString(); got != "kind: Job" {
		t.Fatalf("unexpected pod_spec_override value: got=%q want=%q", got, "kind: Job")
	}
}

func TestSetStateWriteOnlyIntegerDefaultsToTypedNull(t *testing.T) {
	t.Parallel()

	resource := &objectResource{
		object: manifest.ManagedObject{
			Fields: []manifest.FieldSpec{
				{
					Name:      "team",
					Type:      manifest.FieldTypeInt,
					WriteOnly: true,
				},
			},
		},
	}

	target := &mockAttributeTarget{values: map[string]any{}}
	diags := resource.setState(
		context.Background(),
		target,
		"42",
		map[string]any{},
		nil,
		nil,
		nil,
	)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	value, ok := findStateAttributeValue(target.values, "team")
	if !ok {
		t.Fatalf("team attribute was not written")
	}

	intValue, ok := value.(types.Int64)
	if !ok {
		t.Fatalf("expected types.Int64 for write-only team, got %T", value)
	}
	if !intValue.IsNull() {
		t.Fatalf("expected null team when no write-only value is available")
	}
}

func findStateAttributeValue(values map[string]any, attr string) (any, bool) {
	for key, value := range values {
		if strings.Contains(key, attr) {
			return value, true
		}
	}
	return nil, false
}
