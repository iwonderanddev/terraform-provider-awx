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
			value, diags := toTerraformValue("", manifest.FieldSpec{Name: "test", Type: tc.fieldType}, nil)
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
			case manifest.FieldTypeObject:
				v, ok := value.(types.Dynamic)
				if !ok || !v.IsNull() {
					t.Fatalf("expected null types.Dynamic, got %T (%#v)", value, value)
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

	value, diags := toTerraformValue("", manifest.FieldSpec{Name: "max_hosts", Type: manifest.FieldTypeInt}, "not-a-number")
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

func TestToTerraformValueConvertsObjectToDynamic(t *testing.T) {
	t.Parallel()

	value, diags := toTerraformValue(
		"",
		manifest.FieldSpec{Name: "settings", Type: manifest.FieldTypeObject},
		map[string]any{"enabled": true},
	)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	dynamicValue, ok := value.(types.Dynamic)
	if !ok {
		t.Fatalf("expected types.Dynamic for object field, got %T", value)
	}
	if dynamicValue.IsNull() {
		t.Fatalf("expected non-null dynamic value")
	}

	underlying := dynamicValue.UnderlyingValue()
	objectValue, ok := underlying.(types.Object)
	if !ok {
		t.Fatalf("expected underlying object value, got %T", underlying)
	}

	attrs := objectValue.Attributes()
	enabledAttr, ok := attrs["enabled"]
	if !ok {
		t.Fatalf("expected object attribute enabled")
	}
	enabledValue, ok := enabledAttr.(types.Bool)
	if !ok {
		t.Fatalf("expected enabled attribute as types.Bool")
	}
	if !enabledValue.ValueBool() {
		t.Fatalf("expected enabled to be true")
	}
}

func TestFieldUsesStringObjectTransport(t *testing.T) {
	t.Parallel()

	if !fieldUsesStringObjectTransport("job_templates", "extra_vars") {
		t.Fatalf("expected extra_vars string transport for job_templates")
	}
	if !fieldUsesStringObjectTransport("schedules", "extra_data") {
		t.Fatalf("expected extra_data string transport for schedules")
	}
	if fieldUsesStringObjectTransport("schedules", "settings") {
		t.Fatalf("expected non-extra_vars fields to skip string transport")
	}
}

func TestToTerraformValueParsesExtraVarsJSON(t *testing.T) {
	t.Parallel()

	tests := []string{"job_templates", "workflow_job_templates"}
	for _, objectName := range tests {
		objectName := objectName
		t.Run(objectName, func(t *testing.T) {
			t.Parallel()

			value, diags := toTerraformValue(
				objectName,
				manifest.FieldSpec{Name: "extra_vars", Type: manifest.FieldTypeObject},
				`{"foo":"bar","nested":{"enabled":true}}`,
			)
			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			dynamicValue, ok := value.(types.Dynamic)
			if !ok {
				t.Fatalf("expected types.Dynamic for extra_vars, got %T", value)
			}
			underlying := dynamicValue.UnderlyingValue()
			objectValue, ok := underlying.(types.Object)
			if !ok {
				t.Fatalf("expected underlying object value, got %T", underlying)
			}

			fooValue, ok := objectValue.Attributes()["foo"].(types.String)
			if !ok {
				t.Fatalf("expected foo attribute as types.String")
			}
			if got := fooValue.ValueString(); got != "bar" {
				t.Fatalf("unexpected foo value: got=%q want=%q", got, "bar")
			}
		})
	}
}

func TestToTerraformValueParsesExtraVarsWithYAMLFallback(t *testing.T) {
	t.Parallel()

	tests := []string{"job_templates", "workflow_job_templates"}
	for _, objectName := range tests {
		objectName := objectName
		t.Run(objectName, func(t *testing.T) {
			t.Parallel()

			value, diags := toTerraformValue(
				objectName,
				manifest.FieldSpec{Name: "extra_vars", Type: manifest.FieldTypeObject},
				"foo: bar\nnested:\n  enabled: true\n",
			)
			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			dynamicValue, ok := value.(types.Dynamic)
			if !ok {
				t.Fatalf("expected types.Dynamic for extra_vars, got %T", value)
			}
			underlying := dynamicValue.UnderlyingValue()
			objectValue, ok := underlying.(types.Object)
			if !ok {
				t.Fatalf("expected underlying object value, got %T", underlying)
			}

			fooValue, ok := objectValue.Attributes()["foo"].(types.String)
			if !ok {
				t.Fatalf("expected foo attribute as types.String")
			}
			if got := fooValue.ValueString(); got != "bar" {
				t.Fatalf("unexpected foo value: got=%q want=%q", got, "bar")
			}
		})
	}
}

func TestToTerraformValueRejectsNonObjectExtraVarsRoot(t *testing.T) {
	t.Parallel()

	tests := []string{"job_templates", "workflow_job_templates"}
	for _, objectName := range tests {
		objectName := objectName
		t.Run(objectName, func(t *testing.T) {
			t.Parallel()

			value, diags := toTerraformValue(
				objectName,
				manifest.FieldSpec{Name: "extra_vars", Type: manifest.FieldTypeObject},
				"- one\n- two\n",
			)
			if !diags.HasError() {
				t.Fatalf("expected diagnostics for non-object extra_vars root")
			}

			dynamicValue, ok := value.(types.Dynamic)
			if !ok {
				t.Fatalf("expected types.Dynamic for failed object conversion, got %T", value)
			}
			if !dynamicValue.IsNull() {
				t.Fatalf("expected null dynamic value when conversion fails")
			}
		})
	}
}

func TestToTerraformValueParsesExtraDataJSON(t *testing.T) {
	t.Parallel()

	tests := []string{"schedules", "workflow_job_template_nodes", "workflow_job_nodes"}
	for _, objectName := range tests {
		objectName := objectName
		t.Run(objectName, func(t *testing.T) {
			t.Parallel()

			value, diags := toTerraformValue(
				objectName,
				manifest.FieldSpec{Name: "extra_data", Type: manifest.FieldTypeObject},
				`{"foo":"bar","nested":{"enabled":true}}`,
			)
			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			dynamicValue, ok := value.(types.Dynamic)
			if !ok {
				t.Fatalf("expected types.Dynamic for extra_data, got %T", value)
			}
			underlying := dynamicValue.UnderlyingValue()
			objectValue, ok := underlying.(types.Object)
			if !ok {
				t.Fatalf("expected underlying object value, got %T", underlying)
			}

			fooValue, ok := objectValue.Attributes()["foo"].(types.String)
			if !ok {
				t.Fatalf("expected foo attribute as types.String")
			}
			if got := fooValue.ValueString(); got != "bar" {
				t.Fatalf("unexpected foo value: got=%q want=%q", got, "bar")
			}
		})
	}
}

func TestToTerraformValueParsesExtraDataWithYAMLFallback(t *testing.T) {
	t.Parallel()

	tests := []string{"schedules", "workflow_job_template_nodes", "workflow_job_nodes"}
	for _, objectName := range tests {
		objectName := objectName
		t.Run(objectName, func(t *testing.T) {
			t.Parallel()

			value, diags := toTerraformValue(
				objectName,
				manifest.FieldSpec{Name: "extra_data", Type: manifest.FieldTypeObject},
				"foo: bar\nnested:\n  enabled: true\n",
			)
			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			dynamicValue, ok := value.(types.Dynamic)
			if !ok {
				t.Fatalf("expected types.Dynamic for extra_data, got %T", value)
			}
			underlying := dynamicValue.UnderlyingValue()
			objectValue, ok := underlying.(types.Object)
			if !ok {
				t.Fatalf("expected underlying object value, got %T", underlying)
			}

			fooValue, ok := objectValue.Attributes()["foo"].(types.String)
			if !ok {
				t.Fatalf("expected foo attribute as types.String")
			}
			if got := fooValue.ValueString(); got != "bar" {
				t.Fatalf("unexpected foo value: got=%q want=%q", got, "bar")
			}
		})
	}
}

func TestToTerraformValueRejectsNonObjectExtraDataRoot(t *testing.T) {
	t.Parallel()

	tests := []string{"schedules", "workflow_job_template_nodes", "workflow_job_nodes"}
	for _, objectName := range tests {
		objectName := objectName
		t.Run(objectName, func(t *testing.T) {
			t.Parallel()

			value, diags := toTerraformValue(
				objectName,
				manifest.FieldSpec{Name: "extra_data", Type: manifest.FieldTypeObject},
				"- one\n- two\n",
			)
			if !diags.HasError() {
				t.Fatalf("expected diagnostics for non-object extra_data root")
			}

			dynamicValue, ok := value.(types.Dynamic)
			if !ok {
				t.Fatalf("expected types.Dynamic for failed object conversion, got %T", value)
			}
			if !dynamicValue.IsNull() {
				t.Fatalf("expected null dynamic value when conversion fails")
			}
		})
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
