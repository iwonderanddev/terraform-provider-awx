package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gopkg.in/yaml.v3"
)

type objectFieldKey struct {
	objectName string
	fieldName  string
}

var stringObjectTransportFields = map[objectFieldKey]struct{}{
	{objectName: "job_templates", fieldName: "extra_vars"}:               {},
	{objectName: "workflow_job_templates", fieldName: "extra_vars"}:      {},
	{objectName: "schedules", fieldName: "extra_data"}:                   {},
	{objectName: "workflow_job_template_nodes", fieldName: "extra_data"}: {},
	{objectName: "workflow_job_nodes", fieldName: "extra_data"}:          {},
}

func fieldUsesStringObjectTransport(objectName string, fieldName string) bool {
	_, ok := stringObjectTransportFields[objectFieldKey{objectName: objectName, fieldName: fieldName}]
	return ok
}

func terraformDynamicObjectToMap(value types.Dynamic) (map[string]any, error) {
	if value.IsNull() {
		return nil, nil
	}
	if value.IsUnknown() {
		return nil, fmt.Errorf("value is unknown")
	}

	underlying := value.UnderlyingValue()
	if underlying == nil {
		return nil, nil
	}

	native, err := terraformAttrToNativeValue(underlying)
	if err != nil {
		return nil, err
	}

	objectPayload, err := coerceObjectMap(native)
	if err != nil {
		return nil, err
	}

	if objectPayload == nil {
		return nil, nil
	}

	return objectPayload, nil
}

func terraformObjectValueFromAPIValue(objectName string, fieldName string, value any) (types.Dynamic, error) {
	if value == nil {
		return types.DynamicNull(), nil
	}

	normalized := value
	if fieldUsesStringObjectTransport(objectName, fieldName) {
		if raw, ok := value.(string); ok {
			parsed, err := parseStructuredObjectString(raw)
			if err != nil {
				return types.DynamicNull(), err
			}
			if parsed == nil {
				return types.DynamicNull(), nil
			}
			normalized = parsed
		}
	}

	objectValue, err := coerceObjectMap(normalized)
	if err != nil {
		return types.DynamicNull(), err
	}
	if objectValue == nil {
		return types.DynamicNull(), nil
	}

	attrValue, err := nativeToTerraformAttrValue(objectValue)
	if err != nil {
		return types.DynamicNull(), err
	}

	objectAttrValue, ok := attrValue.(types.Object)
	if !ok {
		return types.DynamicNull(), fmt.Errorf("expected object value, got %T", attrValue)
	}

	return types.DynamicValue(objectAttrValue), nil
}

func parseStructuredObjectString(raw string) (map[string]any, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, nil
	}

	var decoded any
	if err := json.Unmarshal([]byte(trimmed), &decoded); err == nil {
		objectValue, objectErr := coerceObjectMap(decoded)
		if objectErr != nil {
			return nil, fmt.Errorf("expected object value after JSON parsing: %w", objectErr)
		}
		if objectValue == nil {
			return nil, fmt.Errorf("expected object value after JSON parsing")
		}
		return objectValue, nil
	} else {
		var yamlDecoded any
		if yamlErr := yaml.Unmarshal([]byte(trimmed), &yamlDecoded); yamlErr != nil {
			return nil, fmt.Errorf("unable to parse structured value as JSON (%s) or YAML (%s)", err.Error(), yamlErr.Error())
		}

		normalized, normErr := normalizeYAMLValue(yamlDecoded)
		if normErr != nil {
			return nil, normErr
		}

		objectValue, objectErr := coerceObjectMap(normalized)
		if objectErr != nil {
			return nil, fmt.Errorf("expected object value after YAML parsing: %w", objectErr)
		}
		if objectValue == nil {
			return nil, fmt.Errorf("expected object value after YAML parsing")
		}

		return objectValue, nil
	}
}

func normalizeYAMLValue(value any) (any, error) {
	switch typed := value.(type) {
	case map[string]any:
		out := make(map[string]any, len(typed))
		for key, child := range typed {
			normalized, err := normalizeYAMLValue(child)
			if err != nil {
				return nil, err
			}
			out[key] = normalized
		}
		return out, nil
	case map[any]any:
		out := make(map[string]any, len(typed))
		for key, child := range typed {
			stringKey, ok := key.(string)
			if !ok {
				return nil, fmt.Errorf("YAML object key must be a string, got %T", key)
			}
			normalized, err := normalizeYAMLValue(child)
			if err != nil {
				return nil, err
			}
			out[stringKey] = normalized
		}
		return out, nil
	case []any:
		out := make([]any, len(typed))
		for i, child := range typed {
			normalized, err := normalizeYAMLValue(child)
			if err != nil {
				return nil, err
			}
			out[i] = normalized
		}
		return out, nil
	default:
		return value, nil
	}
}

func coerceObjectMap(value any) (map[string]any, error) {
	if value == nil {
		return nil, nil
	}

	if typed, ok := value.(map[string]any); ok {
		return typed, nil
	}

	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Map {
		return nil, fmt.Errorf("expected object value, got %T", value)
	}

	out := make(map[string]any, rv.Len())
	for _, key := range rv.MapKeys() {
		if key.Kind() != reflect.String {
			return nil, fmt.Errorf("object key must be a string, got %s", key.Kind())
		}
		out[key.String()] = rv.MapIndex(key).Interface()
	}

	return out, nil
}

func terraformAttrToNativeValue(value attr.Value) (any, error) {
	if value == nil {
		return nil, nil
	}
	if value.IsUnknown() {
		return nil, fmt.Errorf("value contains unknown data")
	}
	if value.IsNull() {
		return nil, nil
	}

	switch typed := value.(type) {
	case types.String:
		return typed.ValueString(), nil
	case types.Bool:
		return typed.ValueBool(), nil
	case types.Int64:
		return typed.ValueInt64(), nil
	case types.Float64:
		return typed.ValueFloat64(), nil
	case types.Number:
		number := typed.ValueBigFloat()
		if number == nil {
			return nil, nil
		}
		if number.IsInt() {
			if intValue, accuracy := number.Int64(); accuracy == big.Exact {
				return intValue, nil
			}
		}
		floatValue, _ := number.Float64()
		return floatValue, nil
	case types.Dynamic:
		if typed.IsUnknown() {
			return nil, fmt.Errorf("value contains unknown data")
		}
		if typed.IsNull() {
			return nil, nil
		}
		return terraformAttrToNativeValue(typed.UnderlyingValue())
	case types.Object:
		attributes := typed.Attributes()
		out := make(map[string]any, len(attributes))
		for name, child := range attributes {
			native, err := terraformAttrToNativeValue(child)
			if err != nil {
				return nil, err
			}
			out[name] = native
		}
		return out, nil
	case types.Map:
		elements := typed.Elements()
		out := make(map[string]any, len(elements))
		for name, child := range elements {
			native, err := terraformAttrToNativeValue(child)
			if err != nil {
				return nil, err
			}
			out[name] = native
		}
		return out, nil
	case types.List:
		elements := typed.Elements()
		out := make([]any, len(elements))
		for i, child := range elements {
			native, err := terraformAttrToNativeValue(child)
			if err != nil {
				return nil, err
			}
			out[i] = native
		}
		return out, nil
	case types.Set:
		elements := typed.Elements()
		out := make([]any, len(elements))
		for i, child := range elements {
			native, err := terraformAttrToNativeValue(child)
			if err != nil {
				return nil, err
			}
			out[i] = native
		}
		return out, nil
	case types.Tuple:
		elements := typed.Elements()
		out := make([]any, len(elements))
		for i, child := range elements {
			native, err := terraformAttrToNativeValue(child)
			if err != nil {
				return nil, err
			}
			out[i] = native
		}
		return out, nil
	default:
		return nil, fmt.Errorf("unsupported Terraform value type %T", value)
	}
}

func nativeToTerraformAttrValue(value any) (attr.Value, error) {
	switch typed := value.(type) {
	case nil:
		return types.DynamicNull(), nil
	case string:
		return types.StringValue(typed), nil
	case bool:
		return types.BoolValue(typed), nil
	case int:
		return types.Int64Value(int64(typed)), nil
	case int8:
		return types.Int64Value(int64(typed)), nil
	case int16:
		return types.Int64Value(int64(typed)), nil
	case int32:
		return types.Int64Value(int64(typed)), nil
	case int64:
		return types.Int64Value(typed), nil
	case uint:
		return types.Int64Value(int64(typed)), nil
	case uint8:
		return types.Int64Value(int64(typed)), nil
	case uint16:
		return types.Int64Value(int64(typed)), nil
	case uint32:
		return types.Int64Value(int64(typed)), nil
	case uint64:
		return types.Int64Value(int64(typed)), nil
	case float32:
		return types.Float64Value(float64(typed)), nil
	case float64:
		return types.Float64Value(typed), nil
	case json.Number:
		if parsedInt, err := typed.Int64(); err == nil {
			return types.Int64Value(parsedInt), nil
		}
		parsedFloat, err := typed.Float64()
		if err != nil {
			return nil, err
		}
		return types.Float64Value(parsedFloat), nil
	case map[string]any:
		return nativeMapToTerraformObject(typed)
	case []any:
		return nativeSliceToTerraformTuple(typed)
	default:
		rv := reflect.ValueOf(value)
		switch rv.Kind() {
		case reflect.Map:
			mapValue, err := coerceObjectMap(value)
			if err != nil {
				return nil, err
			}
			return nativeMapToTerraformObject(mapValue)
		case reflect.Slice, reflect.Array:
			values := make([]any, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				values[i] = rv.Index(i).Interface()
			}
			return nativeSliceToTerraformTuple(values)
		default:
			return nil, fmt.Errorf("unsupported native value type %T", value)
		}
	}
}

func nativeMapToTerraformObject(value map[string]any) (types.Object, error) {
	attributeTypes := make(map[string]attr.Type, len(value))
	attributes := make(map[string]attr.Value, len(value))
	ctx := context.Background()

	for name, child := range value {
		attributeValue, err := nativeToTerraformAttrValue(child)
		if err != nil {
			return types.ObjectNull(nil), err
		}
		attributes[name] = attributeValue
		attributeTypes[name] = attributeValue.Type(ctx)
	}

	objectValue, diags := types.ObjectValue(attributeTypes, attributes)
	if diags.HasError() {
		return types.ObjectNull(nil), diagnosticsToError(diags)
	}

	return objectValue, nil
}

func nativeSliceToTerraformTuple(value []any) (types.Tuple, error) {
	elements := make([]attr.Value, len(value))
	elementTypes := make([]attr.Type, len(value))
	ctx := context.Background()

	for i, child := range value {
		elementValue, err := nativeToTerraformAttrValue(child)
		if err != nil {
			return types.TupleNull(nil), err
		}
		elements[i] = elementValue
		elementTypes[i] = elementValue.Type(ctx)
	}

	tupleValue, diags := types.TupleValue(elementTypes, elements)
	if diags.HasError() {
		return types.TupleNull(nil), diagnosticsToError(diags)
	}

	return tupleValue, nil
}

func diagnosticsToError(diags diag.Diagnostics) error {
	if len(diags) == 0 {
		return nil
	}

	parts := make([]string, 0, len(diags))
	for _, item := range diags {
		parts = append(parts, fmt.Sprintf("%s: %s", item.Summary(), item.Detail()))
	}

	return errors.New(strings.Join(parts, "; "))
}
