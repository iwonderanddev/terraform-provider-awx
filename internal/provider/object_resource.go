package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/damien/terraform-provider-awx-iwd/internal/client"
	"github.com/damien/terraform-provider-awx-iwd/internal/manifest"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/dynamicplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = (*objectResource)(nil)
	_ resource.ResourceWithConfigure   = (*objectResource)(nil)
	_ resource.ResourceWithImportState = (*objectResource)(nil)

	numericIDPattern = regexp.MustCompile(`^[0-9]+$`)
)

type objectResource struct {
	object manifest.ManagedObject
	data   *configuredProvider
}

type attributeSource interface {
	GetAttribute(context.Context, path.Path, any) diag.Diagnostics
}

type attributeTarget interface {
	SetAttribute(context.Context, path.Path, any) diag.Diagnostics
}

// NewObjectResource returns an AWX object resource implementation for a managed object definition.
func NewObjectResource(object manifest.ManagedObject) resource.Resource {
	return &objectResource{object: object}
}

func (r *objectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.object.ResourceName
}

func (r *objectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	attributes := map[string]resourceschema.Attribute{}
	if r.object.CollectionCreate {
		attributes["id"] = resourceschema.Int64Attribute{
			Computed:    true,
			Description: "Numeric AWX ID for this object.",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		}
	} else {
		attributes["id"] = resourceschema.StringAttribute{
			Required:    true,
			Description: "AWX detail-path identifier for this object.",
		}
	}

	for _, field := range r.object.Fields {
		attributes[manifest.TerraformAttributeNameForField(r.object.Name, field)] = newResourceFieldAttribute(r.object.Name, field, r.object.UpdateSupported)
	}

	resp.Schema = resourceschema.Schema{
		Description: fmt.Sprintf("Manages AWX `%s` objects.", r.object.Name),
		Attributes:  attributes,
	}
}

func (r *objectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerData, ok := req.ProviderData.(*configuredProvider)
	if !ok {
		resp.Diagnostics.AddError("Unexpected resource configure type", fmt.Sprintf("expected *configuredProvider, got %T", req.ProviderData))
		return
	}
	r.data = providerData
}

func (r *objectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.data == nil || r.data.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "Expected configured AWX client but provider data was not available.")
		return
	}

	payload, plannedValues, diags := r.payloadFromConfig(ctx, req.Plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plannedStrings := r.stringValuesFromSource(ctx, req.Plan)
	resp.Diagnostics.Append(plannedStrings.Diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}
	plannedJSONArrayStrings := r.jsonEncodedArrayStringValuesFromSource(ctx, req.Plan)
	resp.Diagnostics.Append(plannedJSONArrayStrings.Diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}
	plannedObjects := r.objectValuesFromSource(ctx, req.Plan)
	resp.Diagnostics.Append(plannedObjects.Diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !r.object.CollectionCreate {
		id, idDiags := getResourceID(ctx, req.Plan, r.object.CollectionCreate)
		resp.Diagnostics.Append(idDiags...)
		if resp.Diagnostics.HasError() {
			return
		}

		_, err := r.data.client.UpdateObject(ctx, r.object.DetailPath, id, payload)
		if err != nil {
			resp.Diagnostics.AddError("Failed to create AWX object", err.Error())
			return
		}

		refreshed, refreshErr := r.data.client.GetObject(ctx, r.object.DetailPath, id)
		if refreshErr != nil {
			resp.Diagnostics.AddWarning(
				"Read-after-create refresh failed",
				fmt.Sprintf("Falling back to planned state because post-create read failed: %s", refreshErr.Error()),
			)
			refreshed = map[string]any{}
		}

		resp.Diagnostics.Append(r.setState(ctx, &resp.State, id, refreshed, plannedValues, plannedStrings.Values, plannedObjects.Values, plannedJSONArrayStrings.Values)...)
		return
	}

	created, err := r.data.client.CreateObject(ctx, r.object.CollectionPath, payload)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create AWX object", err.Error())
		return
	}

	parsedID, err := parseNumericID(created["id"])
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse created object ID", err.Error())
		return
	}
	id := strconv.FormatInt(parsedID, 10)

	refreshed, refreshErr := r.data.client.GetObject(ctx, r.object.DetailPath, id)
	if refreshErr != nil {
		resp.Diagnostics.AddWarning(
			"Read-after-create refresh failed",
			fmt.Sprintf("Falling back to create response state because post-create read failed: %s", refreshErr.Error()),
		)
		refreshed = created
	}

	resp.Diagnostics.Append(r.setState(ctx, &resp.State, id, refreshed, plannedValues, plannedStrings.Values, plannedObjects.Values, plannedJSONArrayStrings.Values)...)
}

func (r *objectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.data == nil || r.data.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "Expected configured AWX client but provider data was not available.")
		return
	}

	id, diags := getResourceID(ctx, req.State, r.object.CollectionCreate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	currentValues := r.writeOnlyValuesFromSource(ctx, req.State)
	if currentValues.HasError() {
		resp.Diagnostics.Append(currentValues.Diagnostics...)
		return
	}
	currentStrings := r.stringValuesFromSource(ctx, req.State)
	if currentStrings.HasError() {
		resp.Diagnostics.Append(currentStrings.Diagnostics...)
		return
	}
	currentJSONArrayStrings := r.jsonEncodedArrayStringValuesFromSource(ctx, req.State)
	if currentJSONArrayStrings.HasError() {
		resp.Diagnostics.Append(currentJSONArrayStrings.Diagnostics...)
		return
	}
	currentObjects := r.objectValuesFromSource(ctx, req.State)
	if currentObjects.HasError() {
		resp.Diagnostics.Append(currentObjects.Diagnostics...)
		return
	}

	obj, err := r.data.client.GetObject(ctx, r.object.DetailPath, id)
	if err != nil {
		if shouldRemoveFromStateOnReadError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Failed to read AWX object", err.Error())
		return
	}

	resp.Diagnostics.Append(r.setState(ctx, &resp.State, id, obj, currentValues.Values, currentStrings.Values, currentObjects.Values, currentJSONArrayStrings.Values)...)
}

func (r *objectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.data == nil || r.data.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "Expected configured AWX client but provider data was not available.")
		return
	}
	if !r.object.UpdateSupported {
		resp.Diagnostics.AddError(
			"In-place update not supported",
			fmt.Sprintf("AWX `%s` objects do not support PATCH/PUT updates. Configure changes require replacement.", r.object.Name),
		)
		return
	}

	id, diags := getResourceID(ctx, req.State, r.object.CollectionCreate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload, plannedValues, planDiags := r.payloadFromConfig(ctx, req.Plan)
	resp.Diagnostics.Append(planDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(r.pruneUnchangedFieldsFromPayload(ctx, payload, req.Plan, req.State)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plannedStrings := r.stringValuesFromSource(ctx, req.Plan)
	resp.Diagnostics.Append(plannedStrings.Diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}
	plannedJSONArrayStrings := r.jsonEncodedArrayStringValuesFromSource(ctx, req.Plan)
	resp.Diagnostics.Append(plannedJSONArrayStrings.Diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}
	plannedObjects := r.objectValuesFromSource(ctx, req.Plan)
	resp.Diagnostics.Append(plannedObjects.Diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(payload) > 0 {
		_, err := r.data.client.UpdateObject(ctx, r.object.DetailPath, id, payload)
		if err != nil {
			resp.Diagnostics.AddError("Failed to update AWX object", err.Error())
			return
		}
	}

	refreshed, refreshErr := r.data.client.GetObject(ctx, r.object.DetailPath, id)
	if refreshErr != nil {
		resp.Diagnostics.AddError("Failed to refresh AWX object after update", refreshErr.Error())
		return
	}

	resp.Diagnostics.Append(r.setState(ctx, &resp.State, id, refreshed, plannedValues, plannedStrings.Values, plannedObjects.Values, plannedJSONArrayStrings.Values)...)
}

func (r *objectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.data == nil || r.data.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "Expected configured AWX client but provider data was not available.")
		return
	}

	id, diags := getResourceID(ctx, req.State, r.object.CollectionCreate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.data.client.DeleteObject(ctx, r.object.DetailPath, id); err != nil {
		resp.Diagnostics.AddError("Failed to delete AWX object", err.Error())
		return
	}
}

func (r *objectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := validateObjectImportID(req.ID, r.object.CollectionCreate)
	if err != nil {
		resp.Diagnostics.AddError("Invalid import ID", err.Error())
		return
	}

	if r.object.CollectionCreate {
		parsedID, parseErr := parseNumericID(id)
		if parseErr != nil {
			resp.Diagnostics.AddError("Invalid import ID", parseErr.Error())
			return
		}
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parsedID)...)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func validateObjectImportID(rawID string, collectionCreate bool) (string, error) {
	id := strings.TrimSpace(rawID)
	if id == "" {
		return "", fmt.Errorf("Object resources require a non-empty import identifier.")
	}
	if collectionCreate && !numericIDPattern.MatchString(id) {
		return "", fmt.Errorf("Object resources use numeric AWX IDs. Example: 42")
	}
	return id, nil
}

func (r *objectResource) payloadFromConfig(ctx context.Context, config attributeSource) (map[string]any, map[string]any, diag.Diagnostics) {
	payload := make(map[string]any)
	plannedValues := make(map[string]any)
	diags := diag.Diagnostics{}

	for _, field := range r.object.Fields {
		if field.ReadOnly {
			continue
		}
		tfName := manifest.TerraformAttributeNameForField(r.object.Name, field)
		switch field.Type {
		case manifest.FieldTypeInt:
			var value types.Int64
			diags.Append(config.GetAttribute(ctx, path.Root(tfName), &value)...)
			if value.IsNull() || value.IsUnknown() {
				continue
			}
			if field.WriteOnly {
				plannedValues[field.Name] = value
			}
			payload[field.Name] = value.ValueInt64()
		case manifest.FieldTypeBool:
			var value types.Bool
			diags.Append(config.GetAttribute(ctx, path.Root(tfName), &value)...)
			if value.IsNull() || value.IsUnknown() {
				continue
			}
			if field.WriteOnly {
				plannedValues[field.Name] = value
			}
			payload[field.Name] = value.ValueBool()
		case manifest.FieldTypeFloat:
			var value types.Float64
			diags.Append(config.GetAttribute(ctx, path.Root(tfName), &value)...)
			if value.IsNull() || value.IsUnknown() {
				continue
			}
			if field.WriteOnly {
				plannedValues[field.Name] = value
			}
			payload[field.Name] = value.ValueFloat64()
		case manifest.FieldTypeObject:
			var value types.Dynamic
			diags.Append(config.GetAttribute(ctx, path.Root(tfName), &value)...)
			if value.IsNull() || value.IsUnknown() {
				continue
			}
			if field.WriteOnly {
				plannedValues[field.Name] = value
			}

			objectPayload, objectErr := terraformDynamicObjectToMap(value)
			if objectErr != nil {
				diags.AddAttributeError(path.Root(tfName), "Invalid object payload", objectErr.Error())
				continue
			}
			if objectPayload == nil {
				continue
			}

			if fieldUsesStringObjectTransport(r.object.Name, field.Name) {
				encoded, encodeErr := json.Marshal(objectPayload)
				if encodeErr != nil {
					diags.AddAttributeError(path.Root(tfName), "Invalid object payload", encodeErr.Error())
					continue
				}
				payload[field.Name] = string(encoded)
				continue
			}

			payload[field.Name] = objectPayload
		case manifest.FieldTypeArray:
			if isNativeStringListArrayField(r.object.Name, field.Name) {
				var list types.List
				diags.Append(config.GetAttribute(ctx, path.Root(tfName), &list)...)
				if list.IsNull() || list.IsUnknown() {
					continue
				}
				if field.WriteOnly {
					plannedValues[field.Name] = list
				}
				elems := list.Elements()
				out := make([]any, len(elems))
				okElems := true
				for i, el := range elems {
					sv, ok := el.(types.String)
					if !ok {
						diags.AddAttributeError(path.Root(tfName), "Invalid list element", fmt.Sprintf("expected string at index %d, got %T", i, el))
						okElems = false
						break
					}
					if sv.IsNull() || sv.IsUnknown() {
						diags.AddAttributeError(path.Root(tfName), "Invalid list element", fmt.Sprintf("index %d must be a known string", i))
						okElems = false
						break
					}
					out[i] = sv.ValueString()
				}
				if !okElems {
					continue
				}
				payload[field.Name] = out
				continue
			}
			var value types.String
			diags.Append(config.GetAttribute(ctx, path.Root(tfName), &value)...)
			if value.IsNull() || value.IsUnknown() {
				continue
			}
			if field.WriteOnly {
				plannedValues[field.Name] = value
			}

			decoded, decodeErr := decodeJSONString(value.ValueString())
			if decodeErr != nil {
				diags.AddAttributeError(path.Root(tfName), "Invalid JSON payload", decodeErr.Error())
				continue
			}
			payload[field.Name] = decoded
		default:
			var value types.String
			diags.Append(config.GetAttribute(ctx, path.Root(tfName), &value)...)
			if value.IsNull() || value.IsUnknown() {
				continue
			}
			if field.WriteOnly {
				plannedValues[field.Name] = value
			}

			payload[field.Name] = value.ValueString()
		}
	}

	return payload, plannedValues, diags
}

func (r *objectResource) pruneUnchangedFieldsFromPayload(
	ctx context.Context,
	payload map[string]any,
	plan attributeSource,
	state attributeSource,
) diag.Diagnostics {
	diags := diag.Diagnostics{}

	for _, field := range r.object.Fields {
		if _, exists := payload[field.Name]; !exists {
			continue
		}

		tfName := manifest.TerraformAttributeNameForField(r.object.Name, field)
		attributePath := path.Root(tfName)

		switch field.Type {
		case manifest.FieldTypeInt:
			var planned types.Int64
			var prior types.Int64
			diags.Append(plan.GetAttribute(ctx, attributePath, &planned)...)
			diags.Append(state.GetAttribute(ctx, attributePath, &prior)...)
			if planned.IsNull() || planned.IsUnknown() || prior.IsNull() || prior.IsUnknown() {
				continue
			}
			if planned.ValueInt64() == prior.ValueInt64() {
				delete(payload, field.Name)
			}
		case manifest.FieldTypeBool:
			var planned types.Bool
			var prior types.Bool
			diags.Append(plan.GetAttribute(ctx, attributePath, &planned)...)
			diags.Append(state.GetAttribute(ctx, attributePath, &prior)...)
			if planned.IsNull() || planned.IsUnknown() || prior.IsNull() || prior.IsUnknown() {
				continue
			}
			if planned.ValueBool() == prior.ValueBool() {
				delete(payload, field.Name)
			}
		case manifest.FieldTypeFloat:
			var planned types.Float64
			var prior types.Float64
			diags.Append(plan.GetAttribute(ctx, attributePath, &planned)...)
			diags.Append(state.GetAttribute(ctx, attributePath, &prior)...)
			if planned.IsNull() || planned.IsUnknown() || prior.IsNull() || prior.IsUnknown() {
				continue
			}
			if planned.ValueFloat64() == prior.ValueFloat64() {
				delete(payload, field.Name)
			}
		case manifest.FieldTypeObject:
			var planned types.Dynamic
			var prior types.Dynamic
			diags.Append(plan.GetAttribute(ctx, attributePath, &planned)...)
			diags.Append(state.GetAttribute(ctx, attributePath, &prior)...)
			if planned.IsUnknown() || prior.IsUnknown() {
				continue
			}
			if planned.IsNull() && prior.IsNull() {
				delete(payload, field.Name)
				continue
			}
			if planned.IsNull() || prior.IsNull() {
				continue
			}
			plannedObject, plannedErr := terraformDynamicObjectToMap(planned)
			if plannedErr != nil {
				continue
			}
			priorObject, priorErr := terraformDynamicObjectToMap(prior)
			if priorErr != nil {
				continue
			}
			if reflect.DeepEqual(plannedObject, priorObject) {
				delete(payload, field.Name)
			}
		case manifest.FieldTypeArray:
			if isNativeStringListArrayField(r.object.Name, field.Name) {
				var planned types.List
				var prior types.List
				diags.Append(plan.GetAttribute(ctx, attributePath, &planned)...)
				diags.Append(state.GetAttribute(ctx, attributePath, &prior)...)
				if planned.IsNull() || planned.IsUnknown() || prior.IsNull() || prior.IsUnknown() {
					continue
				}
				if planned.Equal(prior) {
					delete(payload, field.Name)
				}
				continue
			}
			var planned types.String
			var prior types.String
			diags.Append(plan.GetAttribute(ctx, attributePath, &planned)...)
			diags.Append(state.GetAttribute(ctx, attributePath, &prior)...)
			if planned.IsNull() || planned.IsUnknown() || prior.IsNull() || prior.IsUnknown() {
				continue
			}

			plannedString := planned.ValueString()
			priorString := prior.ValueString()

			plannedArray, plannedErr := decodeJSONString(plannedString)
			priorArray, priorErr := decodeJSONString(priorString)
			if plannedErr == nil && priorErr == nil && reflect.DeepEqual(plannedArray, priorArray) {
				delete(payload, field.Name)
				continue
			}

			if plannedString == priorString {
				delete(payload, field.Name)
			}
		default:
			var planned types.String
			var prior types.String
			diags.Append(plan.GetAttribute(ctx, attributePath, &planned)...)
			diags.Append(state.GetAttribute(ctx, attributePath, &prior)...)
			if planned.IsNull() || planned.IsUnknown() || prior.IsNull() || prior.IsUnknown() {
				continue
			}

			plannedString := planned.ValueString()
			priorString := prior.ValueString()

			if fieldHasTrailingNewlineNormalization(r.object.Name, field.Name) &&
				stripSingleTrailingLineEnding(plannedString) == stripSingleTrailingLineEnding(priorString) {
				delete(payload, field.Name)
				continue
			}

			if plannedString == priorString {
				delete(payload, field.Name)
			}
		}
	}

	return diags
}

func (r *objectResource) writeOnlyValuesFromSource(ctx context.Context, source attributeSource) writeOnlyValueSnapshot {
	values := make(map[string]any)
	diags := diag.Diagnostics{}

	for _, field := range r.object.Fields {
		if !field.WriteOnly {
			continue
		}
		tfName := manifest.TerraformAttributeNameForField(r.object.Name, field)

		switch field.Type {
		case manifest.FieldTypeInt:
			var value types.Int64
			diags.Append(source.GetAttribute(ctx, path.Root(tfName), &value)...)
			if !value.IsNull() && !value.IsUnknown() {
				values[field.Name] = value
			}
		case manifest.FieldTypeBool:
			var value types.Bool
			diags.Append(source.GetAttribute(ctx, path.Root(tfName), &value)...)
			if !value.IsNull() && !value.IsUnknown() {
				values[field.Name] = value
			}
		case manifest.FieldTypeFloat:
			var value types.Float64
			diags.Append(source.GetAttribute(ctx, path.Root(tfName), &value)...)
			if !value.IsNull() && !value.IsUnknown() {
				values[field.Name] = value
			}
		case manifest.FieldTypeObject:
			var value types.Dynamic
			diags.Append(source.GetAttribute(ctx, path.Root(tfName), &value)...)
			if !value.IsNull() && !value.IsUnknown() {
				values[field.Name] = value
			}
		case manifest.FieldTypeArray:
			if isNativeStringListArrayField(r.object.Name, field.Name) {
				var value types.List
				diags.Append(source.GetAttribute(ctx, path.Root(tfName), &value)...)
				if !value.IsNull() && !value.IsUnknown() {
					values[field.Name] = value
				}
				continue
			}
			fallthrough
		default:
			var value types.String
			diags.Append(source.GetAttribute(ctx, path.Root(tfName), &value)...)
			if !value.IsNull() && !value.IsUnknown() {
				values[field.Name] = value
			}
		}
	}

	return writeOnlyValueSnapshot{Values: values, Diagnostics: diags}
}

func (r *objectResource) stringValuesFromSource(ctx context.Context, source attributeSource) valueSnapshot {
	values := make(map[string]types.String)
	diags := diag.Diagnostics{}

	for _, field := range r.object.Fields {
		if field.Type != manifest.FieldTypeString {
			continue
		}

		tfName := manifest.TerraformAttributeNameForField(r.object.Name, field)
		var value types.String
		diags.Append(source.GetAttribute(ctx, path.Root(tfName), &value)...)
		if value.IsUnknown() {
			continue
		}
		values[field.Name] = value
	}

	return valueSnapshot{Values: values, Diagnostics: diags}
}

// jsonEncodedArrayStringValuesFromSource collects Terraform string values for manifest array fields that are
// transported as JSON-encoded strings (non-native list types). Used to align read/apply state with omitted
// optional attributes when AWX returns empty JSON arrays.
func (r *objectResource) jsonEncodedArrayStringValuesFromSource(ctx context.Context, source attributeSource) valueSnapshot {
	values := make(map[string]types.String)
	diags := diag.Diagnostics{}

	for _, field := range r.object.Fields {
		if field.Type != manifest.FieldTypeArray {
			continue
		}
		if isNativeStringListArrayField(r.object.Name, field.Name) {
			continue
		}

		tfName := manifest.TerraformAttributeNameForField(r.object.Name, field)
		var value types.String
		diags.Append(source.GetAttribute(ctx, path.Root(tfName), &value)...)
		if value.IsUnknown() {
			continue
		}
		values[field.Name] = value
	}

	return valueSnapshot{Values: values, Diagnostics: diags}
}

func (r *objectResource) objectValuesFromSource(ctx context.Context, source attributeSource) objectValueSnapshot {
	values := make(map[string]types.Dynamic)
	diags := diag.Diagnostics{}

	for _, field := range r.object.Fields {
		if field.Type != manifest.FieldTypeObject {
			continue
		}

		tfName := manifest.TerraformAttributeNameForField(r.object.Name, field)
		var value types.Dynamic
		diags.Append(source.GetAttribute(ctx, path.Root(tfName), &value)...)
		if value.IsUnknown() {
			continue
		}
		values[field.Name] = value
	}

	return objectValueSnapshot{Values: values, Diagnostics: diags}
}

func (r *objectResource) setState(
	ctx context.Context,
	state attributeTarget,
	id string,
	apiObject map[string]any,
	writeOnlyValues map[string]any,
	priorStringValues map[string]types.String,
	priorObjectValues map[string]types.Dynamic,
	priorJSONEncodedArrayValues map[string]types.String,
) diag.Diagnostics {
	diags := diag.Diagnostics{}
	if r.object.CollectionCreate {
		parsedID, err := parseNumericID(id)
		if err != nil {
			diags.AddError("Failed to set resource ID in state", err.Error())
			return diags
		}
		diags.Append(state.SetAttribute(ctx, path.Root("id"), parsedID)...)
	} else {
		diags.Append(state.SetAttribute(ctx, path.Root("id"), id)...)
	}

	for _, field := range r.object.Fields {
		tfName := manifest.TerraformAttributeNameForField(r.object.Name, field)
		if field.WriteOnly {
			preserved, ok := writeOnlyValues[field.Name]
			if ok {
				diags.Append(state.SetAttribute(ctx, path.Root(tfName), preserved)...)
			} else {
				nullValue, nullDiags := toTerraformValue(r.object.Name, field, nil)
				diags.Append(nullDiags...)
				if !nullDiags.HasError() {
					diags.Append(state.SetAttribute(ctx, path.Root(tfName), nullValue)...)
				}
			}
			continue
		}

		value := apiObject[field.Name]
		if field.Computed && !field.Required {
			if field.Type == manifest.FieldTypeObject {
				if priorObject, ok := priorObjectValues[field.Name]; ok && !priorObject.IsUnknown() {
					diags.Append(state.SetAttribute(ctx, path.Root(tfName), priorObject)...)
					continue
				}
			}
			if field.Type == manifest.FieldTypeString && !field.ReadOnly {
				if priorString, ok := priorStringValues[field.Name]; ok && !priorString.IsUnknown() && priorString.IsNull() {
					diags.Append(state.SetAttribute(ctx, path.Root(tfName), priorString)...)
					continue
				}
			}
		}

		if value == nil && field.Type == manifest.FieldTypeObject {
			if priorObject, ok := priorObjectValues[field.Name]; ok && !priorObject.IsUnknown() && priorObject.IsNull() {
				diags.Append(state.SetAttribute(ctx, path.Root(tfName), priorObject)...)
				continue
			}
		}
		if normalized, ok := normalizeOptionalEmptyStringToNull(field, value, priorStringValues); ok {
			diags.Append(state.SetAttribute(ctx, path.Root(tfName), normalized)...)
			continue
		}
		if normalized, ok := normalizeOptionalEmptyJSONEncodedArrayToNull(r.object.Name, field, value, priorJSONEncodedArrayValues); ok {
			diags.Append(state.SetAttribute(ctx, path.Root(tfName), normalized)...)
			continue
		}
		if preserved, ok := preserveKnownNormalizedStringField(r.object.Name, field, value, priorStringValues); ok {
			diags.Append(state.SetAttribute(ctx, path.Root(tfName), preserved)...)
			continue
		}

		if field.Type == manifest.FieldTypeObject {
			if priorObject, ok := priorObjectValues[field.Name]; ok && !priorObject.IsUnknown() {
				preserve, preserveErr := shouldPreserveObjectValue(r.object.Name, field.Name, value, priorObject)
				if preserveErr != nil {
					diags.AddError("Failed to compare object field values", fmt.Sprintf("field=%s err=%s", field.Name, preserveErr.Error()))
					continue
				}
				if preserve {
					diags.Append(state.SetAttribute(ctx, path.Root(tfName), priorObject)...)
					continue
				}
			}
		}

		converted, convDiags := toTerraformValue(r.object.Name, field, value)
		diags.Append(convDiags...)
		if convDiags.HasError() {
			continue
		}
		diags.Append(state.SetAttribute(ctx, path.Root(tfName), converted)...)
	}

	return diags
}

func normalizeOptionalEmptyStringToNull(field manifest.FieldSpec, value any, priorStringValues map[string]types.String) (types.String, bool) {
	if field.Type != manifest.FieldTypeString || field.Required || field.ReadOnly {
		return types.String{}, false
	}

	strValue, ok := value.(string)
	if !ok || strValue != "" {
		return types.String{}, false
	}

	prior, hasPrior := priorStringValues[field.Name]
	if !hasPrior || !prior.IsNull() {
		return types.String{}, false
	}

	return types.StringNull(), true
}

func normalizeOptionalEmptyJSONEncodedArrayToNull(
	objectName string,
	field manifest.FieldSpec,
	value any,
	priorJSONEncodedArrayValues map[string]types.String,
) (types.String, bool) {
	if field.Type != manifest.FieldTypeArray {
		return types.String{}, false
	}
	if isNativeStringListArrayField(objectName, field.Name) {
		return types.String{}, false
	}
	if field.Required || field.Computed {
		return types.String{}, false
	}

	prior, hasPrior := priorJSONEncodedArrayValues[field.Name]
	if !hasPrior || !prior.IsNull() {
		return types.String{}, false
	}

	if !isEmptyJSONArrayAPIValue(value) {
		return types.String{}, false
	}

	return types.StringNull(), true
}

func isEmptyJSONArrayAPIValue(value any) bool {
	if value == nil {
		return false
	}
	switch v := value.(type) {
	case []any:
		return len(v) == 0
	default:
		rv := reflect.ValueOf(value)
		if rv.Kind() == reflect.Slice {
			return rv.Len() == 0
		}
	}
	return false
}

func preserveKnownNormalizedStringField(objectName string, field manifest.FieldSpec, value any, priorStringValues map[string]types.String) (types.String, bool) {
	if field.Type != manifest.FieldTypeString || !fieldHasTrailingNewlineNormalization(objectName, field.Name) {
		return types.String{}, false
	}

	apiString, ok := value.(string)
	if !ok {
		return types.String{}, false
	}

	prior, hasPrior := priorStringValues[field.Name]
	if !hasPrior || prior.IsNull() || prior.IsUnknown() {
		return types.String{}, false
	}

	priorString := prior.ValueString()
	if priorString == apiString {
		return prior, true
	}
	if stripSingleTrailingLineEnding(priorString) == stripSingleTrailingLineEnding(apiString) {
		return prior, true
	}

	return types.String{}, false
}

func fieldHasTrailingNewlineNormalization(objectName string, fieldName string) bool {
	return objectName == "instance_groups" && fieldName == "pod_spec_override"
}

func stripSingleTrailingLineEnding(value string) string {
	if strings.HasSuffix(value, "\r\n") {
		return strings.TrimSuffix(value, "\r\n")
	}
	if strings.HasSuffix(value, "\n") {
		return strings.TrimSuffix(value, "\n")
	}
	if strings.HasSuffix(value, "\r") {
		return strings.TrimSuffix(value, "\r")
	}
	return value
}

func shouldPreserveObjectValue(objectName string, fieldName string, apiValue any, prior types.Dynamic) (bool, error) {
	if prior.IsUnknown() {
		return false, nil
	}
	if apiValue == nil {
		return prior.IsNull(), nil
	}

	priorMap, err := terraformDynamicObjectToMap(prior)
	if err != nil {
		return false, err
	}

	apiDynamic, err := terraformObjectValueFromAPIValue(objectName, fieldName, apiValue)
	if err != nil {
		return false, err
	}
	apiMap, err := terraformDynamicObjectToMap(apiDynamic)
	if err != nil {
		return false, err
	}

	return reflect.DeepEqual(priorMap, apiMap), nil
}

func newResourceFieldAttribute(objectName string, field manifest.FieldSpec, updateSupported bool) resourceschema.Attribute {
	optional := !field.Required && !field.ReadOnly
	computed := field.ReadOnly || (!field.Required && field.Computed)
	requiresReplace := !updateSupported
	switch field.Type {
	case manifest.FieldTypeInt:
		planModifiers := []planmodifier.Int64{}
		if computed {
			planModifiers = append(planModifiers, int64planmodifier.UseStateForUnknown())
		}
		if requiresReplace {
			planModifiers = append(planModifiers, int64planmodifier.RequiresReplace())
		}
		return resourceschema.Int64Attribute{
			Description:   fieldDescription(objectName, field),
			Required:      field.Required,
			Optional:      optional,
			Computed:      computed,
			Sensitive:     field.Sensitive,
			PlanModifiers: planModifiers,
		}
	case manifest.FieldTypeBool:
		planModifiers := []planmodifier.Bool{}
		if computed {
			planModifiers = append(planModifiers, boolplanmodifier.UseStateForUnknown())
		}
		if requiresReplace {
			planModifiers = append(planModifiers, boolplanmodifier.RequiresReplace())
		}
		return resourceschema.BoolAttribute{
			Description:   fieldDescription(objectName, field),
			Required:      field.Required,
			Optional:      optional,
			Computed:      computed,
			Sensitive:     field.Sensitive,
			PlanModifiers: planModifiers,
		}
	case manifest.FieldTypeFloat:
		planModifiers := []planmodifier.Float64{}
		if computed {
			planModifiers = append(planModifiers, float64planmodifier.UseStateForUnknown())
		}
		if requiresReplace {
			planModifiers = append(planModifiers, float64planmodifier.RequiresReplace())
		}
		return resourceschema.Float64Attribute{
			Description:   fieldDescription(objectName, field),
			Required:      field.Required,
			Optional:      optional,
			Computed:      computed,
			Sensitive:     field.Sensitive,
			PlanModifiers: planModifiers,
		}
	case manifest.FieldTypeObject:
		planModifiers := []planmodifier.Dynamic{}
		if computed {
			planModifiers = append(planModifiers, dynamicplanmodifier.UseStateForUnknown())
		}
		if requiresReplace {
			planModifiers = append(planModifiers, dynamicplanmodifier.RequiresReplace())
		}
		return resourceschema.DynamicAttribute{
			Description:   fieldDescription(objectName, field),
			Required:      field.Required,
			Optional:      optional,
			Computed:      computed,
			Sensitive:     field.Sensitive,
			PlanModifiers: planModifiers,
		}
	case manifest.FieldTypeArray:
		if isNativeStringListArrayField(objectName, field.Name) {
			planModifiers := []planmodifier.List{}
			if computed {
				planModifiers = append(planModifiers, listplanmodifier.UseStateForUnknown())
			}
			if requiresReplace {
				planModifiers = append(planModifiers, listplanmodifier.RequiresReplace())
			}
			return resourceschema.ListAttribute{
				ElementType:   types.StringType,
				Description:   fieldDescription(objectName, field),
				Required:      field.Required,
				Optional:      optional,
				Computed:      computed,
				Sensitive:     field.Sensitive,
				PlanModifiers: planModifiers,
			}
		}
		fallthrough
	case manifest.FieldTypeString:
		planModifiers := []planmodifier.String{}
		if computed {
			planModifiers = append(planModifiers, stringplanmodifier.UseStateForUnknown())
		}
		if requiresReplace {
			planModifiers = append(planModifiers, stringplanmodifier.RequiresReplace())
		}
		return resourceschema.StringAttribute{
			Description:   fieldDescription(objectName, field),
			Required:      field.Required,
			Optional:      optional,
			Computed:      computed,
			Sensitive:     field.Sensitive,
			PlanModifiers: planModifiers,
		}
	default:
		planModifiers := []planmodifier.String{}
		if computed {
			planModifiers = append(planModifiers, stringplanmodifier.UseStateForUnknown())
		}
		if requiresReplace {
			planModifiers = append(planModifiers, stringplanmodifier.RequiresReplace())
		}
		return resourceschema.StringAttribute{
			Description:   fieldDescription(objectName, field),
			Required:      field.Required,
			Optional:      optional,
			Computed:      computed,
			Sensitive:     field.Sensitive,
			PlanModifiers: planModifiers,
		}
	}
}

func isNativeStringListArrayField(objectName string, fieldName string) bool {
	return objectName == "role_definitions" && fieldName == "permissions"
}

func nativeStringSliceToTerraformList(value any) (types.List, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	var elems []attr.Value
	switch v := value.(type) {
	case []any:
		elems = make([]attr.Value, len(v))
		for i, x := range v {
			s, ok := x.(string)
			if !ok {
				diags.AddError("Failed to convert array field", fmt.Sprintf("index %d: expected string, got %T", i, x))
				return types.ListNull(types.StringType), diags
			}
			elems[i] = types.StringValue(s)
		}
	case []string:
		elems = make([]attr.Value, len(v))
		for i, s := range v {
			elems[i] = types.StringValue(s)
		}
	default:
		rv := reflect.ValueOf(value)
		if rv.Kind() == reflect.Slice {
			n := rv.Len()
			elems = make([]attr.Value, n)
			for i := 0; i < n; i++ {
				el := rv.Index(i).Interface()
				s, ok := el.(string)
				if !ok {
					diags.AddError("Failed to convert array field", fmt.Sprintf("index %d: expected string, got %T", i, el))
					return types.ListNull(types.StringType), diags
				}
				elems[i] = types.StringValue(s)
			}
		} else {
			diags.AddError("Failed to convert array field", fmt.Sprintf("expected a string slice, got %T", value))
			return types.ListNull(types.StringType), diags
		}
	}
	list, ldiags := types.ListValue(types.StringType, elems)
	diags.Append(ldiags...)
	return list, diags
}

func fieldDescription(objectName string, field manifest.FieldSpec) string {
	if strings.TrimSpace(field.Description) != "" {
		return field.Description
	}
	if field.Type == manifest.FieldTypeArray {
		if isNativeStringListArrayField(objectName, field.Name) {
			return "List of permission strings for this role."
		}
		return "JSON-encoded value for this AWX field."
	}
	if field.Type == manifest.FieldTypeObject {
		return "Object value for this AWX field."
	}
	return "Managed field from AWX schema."
}

func toTerraformValue(objectName string, field manifest.FieldSpec, value any) (any, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	if value == nil {
		switch field.Type {
		case manifest.FieldTypeInt:
			return types.Int64Null(), diags
		case manifest.FieldTypeBool:
			return types.BoolNull(), diags
		case manifest.FieldTypeFloat:
			return types.Float64Null(), diags
		case manifest.FieldTypeObject:
			return types.DynamicNull(), diags
		case manifest.FieldTypeArray:
			if isNativeStringListArrayField(objectName, field.Name) {
				return types.ListNull(types.StringType), diags
			}
			return types.StringNull(), diags
		default:
			return types.StringNull(), diags
		}
	}

	switch field.Type {
	case manifest.FieldTypeInt:
		parsed, err := parseNumericID(value)
		if err != nil {
			diags.AddError("Failed to convert integer field", fmt.Sprintf("field=%s err=%s", field.Name, err.Error()))
			return types.Int64Null(), diags
		}
		return types.Int64Value(parsed), diags
	case manifest.FieldTypeBool:
		parsed, ok := value.(bool)
		if !ok {
			diags.AddError("Failed to convert boolean field", fmt.Sprintf("field=%s value=%v", field.Name, value))
			return types.BoolNull(), diags
		}
		return types.BoolValue(parsed), diags
	case manifest.FieldTypeFloat:
		parsed, err := parseFloat(value)
		if err != nil {
			diags.AddError("Failed to convert number field", fmt.Sprintf("field=%s err=%s", field.Name, err.Error()))
			return types.Float64Null(), diags
		}
		return types.Float64Value(parsed), diags
	case manifest.FieldTypeObject:
		dynamicValue, err := terraformObjectValueFromAPIValue(objectName, field.Name, value)
		if err != nil {
			diags.AddError("Failed to convert object field", fmt.Sprintf("field=%s err=%s", field.Name, err.Error()))
			return types.DynamicNull(), diags
		}
		return dynamicValue, diags
	case manifest.FieldTypeArray:
		if isNativeStringListArrayField(objectName, field.Name) {
			listVal, listDiags := nativeStringSliceToTerraformList(value)
			diags.Append(listDiags...)
			return listVal, diags
		}
		encoded, err := json.Marshal(value)
		if err != nil {
			diags.AddError("Failed to encode complex field as JSON", fmt.Sprintf("field=%s err=%s", field.Name, err.Error()))
			return types.StringNull(), diags
		}
		return types.StringValue(string(encoded)), diags
	default:
		return types.StringValue(fmt.Sprintf("%v", value)), diags
	}
}

func getResourceID(ctx context.Context, state attributeSource, collectionCreate bool) (string, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	if collectionCreate {
		var id types.Int64
		diags.Append(state.GetAttribute(ctx, path.Root("id"), &id)...)
		if id.IsUnknown() || id.IsNull() {
			diags.AddError("Missing resource ID", "Expected state to contain a numeric AWX identifier.")
			return "", diags
		}
		return strconv.FormatInt(id.ValueInt64(), 10), diags
	}

	var id types.String
	diags.Append(state.GetAttribute(ctx, path.Root("id"), &id)...)
	if id.IsUnknown() || id.IsNull() {
		diags.AddError("Missing resource ID", "Expected state to contain an AWX identifier.")
		return "", diags
	}

	identifier := strings.TrimSpace(id.ValueString())
	if identifier == "" {
		diags.AddError("Invalid resource ID", "id cannot be empty.")
		return "", diags
	}
	if collectionCreate && !numericIDPattern.MatchString(identifier) {
		diags.AddError("Invalid resource ID", "Expected state to contain a numeric AWX ID.")
		return "", diags
	}
	return identifier, diags
}

func parseNumericID(value any) (int64, error) {
	switch v := value.(type) {
	case int:
		return int64(v), nil
	case int64:
		return v, nil
	case float64:
		return int64(v), nil
	case json.Number:
		return v.Int64()
	case string:
		return strconv.ParseInt(v, 10, 64)
	default:
		return 0, fmt.Errorf("unsupported ID type %T", value)
	}
}

func parseFloat(value any) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case int:
		return float64(v), nil
	case json.Number:
		return v.Float64()
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("unsupported numeric type %T", value)
	}
}

func decodeJSONString(value string) (any, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}
	var decoded any
	if err := json.Unmarshal([]byte(value), &decoded); err != nil {
		return nil, err
	}
	return decoded, nil
}

func asAPIError(err error) *client.APIError {
	if err == nil {
		return nil
	}
	apiErr, ok := err.(*client.APIError)
	if !ok {
		return nil
	}
	return apiErr
}

func shouldRemoveFromStateOnReadError(err error) bool {
	apiErr := asAPIError(err)
	return apiErr != nil && apiErr.StatusCode == http.StatusNotFound
}

type valueSnapshot struct {
	Values      map[string]types.String
	Diagnostics diag.Diagnostics
}

func (s valueSnapshot) HasError() bool {
	return s.Diagnostics.HasError()
}

type writeOnlyValueSnapshot struct {
	Values      map[string]any
	Diagnostics diag.Diagnostics
}

func (s writeOnlyValueSnapshot) HasError() bool {
	return s.Diagnostics.HasError()
}

type objectValueSnapshot struct {
	Values      map[string]types.Dynamic
	Diagnostics diag.Diagnostics
}

func (s objectValueSnapshot) HasError() bool {
	return s.Diagnostics.HasError()
}
