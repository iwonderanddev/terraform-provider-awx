package provider

import (
	"context"
	"strings"
	"testing"

	"github.com/damien/terraform-awx-provider/internal/manifest"
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

func findStateAttributeValue(values map[string]any, attr string) (any, bool) {
	for key, value := range values {
		if strings.Contains(key, attr) {
			return value, true
		}
	}
	return nil, false
}
