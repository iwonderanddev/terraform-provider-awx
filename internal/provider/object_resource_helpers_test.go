package provider

import (
	"encoding/json"
	"testing"

	"github.com/damien/terraform-awx-provider/internal/manifest"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestToTerraformValueNilReturnsTypedNull(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		fieldType manifest.FieldType
	}{
		{name: "integer", fieldType: manifest.FieldTypeInt},
		{name: "boolean", fieldType: manifest.FieldTypeBool},
		{name: "number", fieldType: manifest.FieldTypeFloat},
		{name: "string", fieldType: manifest.FieldTypeString},
		{name: "array", fieldType: manifest.FieldTypeArray},
		{name: "object", fieldType: manifest.FieldTypeObject},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			value, diags := toTerraformValue(manifest.FieldSpec{Name: "test", Type: tc.fieldType}, nil)
			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			switch tc.fieldType {
			case manifest.FieldTypeInt:
				v, ok := value.(types.Int64)
				if !ok || !v.IsNull() {
					t.Fatalf("expected null types.Int64, got %T (%#v)", value, value)
				}
			case manifest.FieldTypeBool:
				v, ok := value.(types.Bool)
				if !ok || !v.IsNull() {
					t.Fatalf("expected null types.Bool, got %T (%#v)", value, value)
				}
			case manifest.FieldTypeFloat:
				v, ok := value.(types.Float64)
				if !ok || !v.IsNull() {
					t.Fatalf("expected null types.Float64, got %T (%#v)", value, value)
				}
			default:
				v, ok := value.(types.String)
				if !ok || !v.IsNull() {
					t.Fatalf("expected null types.String, got %T (%#v)", value, value)
				}
			}
		})
	}
}

func TestToTerraformValueIntegerConversionError(t *testing.T) {
	t.Parallel()

	value, diags := toTerraformValue(manifest.FieldSpec{Name: "max_hosts", Type: manifest.FieldTypeInt}, "not-a-number")
	if !diags.HasError() {
		t.Fatalf("expected conversion diagnostics for invalid integer input")
	}

	intValue, ok := value.(types.Int64)
	if !ok {
		t.Fatalf("expected types.Int64 on conversion error, got %T", value)
	}
	if !intValue.IsNull() {
		t.Fatalf("expected null int value when conversion fails")
	}
}

func TestToTerraformValueEncodesComplexValuesAsJSON(t *testing.T) {
	t.Parallel()

	value, diags := toTerraformValue(
		manifest.FieldSpec{Name: "settings", Type: manifest.FieldTypeObject},
		map[string]any{"enabled": true},
	)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	stringValue, ok := value.(types.String)
	if !ok {
		t.Fatalf("expected types.String for encoded object, got %T", value)
	}
	if stringValue.IsNull() {
		t.Fatalf("expected non-null encoded JSON string")
	}

	var decoded map[string]any
	if err := json.Unmarshal([]byte(stringValue.ValueString()), &decoded); err != nil {
		t.Fatalf("encoded value is not valid JSON: %v", err)
	}
	if got, ok := decoded["enabled"].(bool); !ok || !got {
		t.Fatalf("unexpected decoded value: %#v", decoded["enabled"])
	}
}

func TestParseFloatSupportsMultipleInputTypes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   any
		want    float64
		wantErr bool
	}{
		{name: "float64", input: float64(1.5), want: 1.5},
		{name: "int", input: int(2), want: 2.0},
		{name: "json number", input: json.Number("3.25"), want: 3.25},
		{name: "string", input: "4.5", want: 4.5},
		{name: "invalid", input: true, wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseFloat(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected parseFloat error")
				}
				return
			}
			if err != nil {
				t.Fatalf("parseFloat returned error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("unexpected parseFloat value: got=%v want=%v", got, tc.want)
			}
		})
	}
}

func TestDecodeJSONStringCases(t *testing.T) {
	t.Parallel()

	decoded, err := decodeJSONString("  ")
	if err != nil {
		t.Fatalf("unexpected error for empty JSON input: %v", err)
	}
	if decoded != nil {
		t.Fatalf("expected nil decoded value for empty JSON input")
	}

	decoded, err = decodeJSONString(`{"enabled":true}`)
	if err != nil {
		t.Fatalf("unexpected error for valid JSON input: %v", err)
	}
	obj, ok := decoded.(map[string]any)
	if !ok {
		t.Fatalf("expected decoded object, got %T", decoded)
	}
	if got, ok := obj["enabled"].(bool); !ok || !got {
		t.Fatalf("unexpected decoded enabled value: %#v", obj["enabled"])
	}

	if _, err := decodeJSONString("{"); err == nil {
		t.Fatalf("expected JSON decode error for invalid input")
	}
}
