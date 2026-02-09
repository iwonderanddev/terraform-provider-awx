package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/damien/terraform-awx-provider/internal/client"
	"github.com/damien/terraform-awx-provider/internal/manifest"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
	attributes := map[string]resourceschema.Attribute{
		"id": resourceschema.StringAttribute{
			Computed:    true,
			Description: "Numeric AWX ID for this object.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
	}

	for _, field := range r.object.Fields {
		attributes[field.Name] = newResourceFieldAttribute(field)
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

	created, err := r.data.client.CreateObject(ctx, r.object.CollectionPath, payload)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create AWX object", err.Error())
		return
	}

	id, err := parseNumericID(created["id"])
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse created object ID", err.Error())
		return
	}

	refreshed, refreshErr := r.data.client.GetObject(ctx, r.object.DetailPath, id)
	if refreshErr != nil {
		resp.Diagnostics.AddWarning(
			"Read-after-create refresh failed",
			fmt.Sprintf("Falling back to create response state because post-create read failed: %s", refreshErr.Error()),
		)
		refreshed = created
	}

	resp.Diagnostics.Append(r.setState(ctx, &resp.State, id, refreshed, plannedValues, plannedStrings.Values)...)
}

func (r *objectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.data == nil || r.data.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "Expected configured AWX client but provider data was not available.")
		return
	}

	id, diags := getResourceID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	currentValues := r.valuesFromConfig(ctx, req.State)
	if currentValues.HasError() {
		resp.Diagnostics.Append(currentValues.Diagnostics...)
		return
	}
	currentStrings := r.stringValuesFromSource(ctx, req.State)
	if currentStrings.HasError() {
		resp.Diagnostics.Append(currentStrings.Diagnostics...)
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

	resp.Diagnostics.Append(r.setState(ctx, &resp.State, id, obj, currentValues.Values, currentStrings.Values)...)
}

func (r *objectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.data == nil || r.data.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "Expected configured AWX client but provider data was not available.")
		return
	}

	id, diags := getResourceID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload, plannedValues, planDiags := r.payloadFromConfig(ctx, req.Plan)
	resp.Diagnostics.Append(planDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plannedStrings := r.stringValuesFromSource(ctx, req.Plan)
	resp.Diagnostics.Append(plannedStrings.Diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.data.client.UpdateObject(ctx, r.object.DetailPath, id, payload)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update AWX object", err.Error())
		return
	}

	refreshed, refreshErr := r.data.client.GetObject(ctx, r.object.DetailPath, id)
	if refreshErr != nil {
		resp.Diagnostics.AddError("Failed to refresh AWX object after update", refreshErr.Error())
		return
	}

	resp.Diagnostics.Append(r.setState(ctx, &resp.State, id, refreshed, plannedValues, plannedStrings.Values)...)
}

func (r *objectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.data == nil || r.data.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "Expected configured AWX client but provider data was not available.")
		return
	}

	id, diags := getResourceID(ctx, req.State)
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
	id := strings.TrimSpace(req.ID)
	if !numericIDPattern.MatchString(id) {
		resp.Diagnostics.AddError("Invalid import ID", "Object resources use numeric AWX IDs. Example: 42")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (r *objectResource) payloadFromConfig(ctx context.Context, config attributeSource) (map[string]any, map[string]types.String, diag.Diagnostics) {
	payload := make(map[string]any)
	plannedValues := make(map[string]types.String)
	diags := diag.Diagnostics{}

	for _, field := range r.object.Fields {
		switch field.Type {
		case manifest.FieldTypeInt:
			var value types.Int64
			diags.Append(config.GetAttribute(ctx, path.Root(field.Name), &value)...)
			if value.IsNull() || value.IsUnknown() {
				continue
			}
			payload[field.Name] = value.ValueInt64()
		case manifest.FieldTypeBool:
			var value types.Bool
			diags.Append(config.GetAttribute(ctx, path.Root(field.Name), &value)...)
			if value.IsNull() || value.IsUnknown() {
				continue
			}
			payload[field.Name] = value.ValueBool()
		case manifest.FieldTypeFloat:
			var value types.Float64
			diags.Append(config.GetAttribute(ctx, path.Root(field.Name), &value)...)
			if value.IsNull() || value.IsUnknown() {
				continue
			}
			payload[field.Name] = value.ValueFloat64()
		default:
			var value types.String
			diags.Append(config.GetAttribute(ctx, path.Root(field.Name), &value)...)
			if value.IsNull() || value.IsUnknown() {
				continue
			}
			if field.WriteOnly {
				plannedValues[field.Name] = value
			}

			if field.Type == manifest.FieldTypeArray || field.Type == manifest.FieldTypeObject {
				decoded, decodeErr := decodeJSONString(value.ValueString())
				if decodeErr != nil {
					diags.AddAttributeError(path.Root(field.Name), "Invalid JSON payload", decodeErr.Error())
					continue
				}
				payload[field.Name] = decoded
				continue
			}

			payload[field.Name] = value.ValueString()
		}
	}

	return payload, plannedValues, diags
}

func (r *objectResource) valuesFromConfig(ctx context.Context, config attributeSource) valueSnapshot {
	values := make(map[string]types.String)
	diags := diag.Diagnostics{}

	for _, field := range r.object.Fields {
		if !field.WriteOnly {
			continue
		}
		var value types.String
		diags.Append(config.GetAttribute(ctx, path.Root(field.Name), &value)...)
		if !value.IsNull() && !value.IsUnknown() {
			values[field.Name] = value
		}
	}

	return valueSnapshot{Values: values, Diagnostics: diags}
}

func (r *objectResource) stringValuesFromSource(ctx context.Context, source attributeSource) valueSnapshot {
	values := make(map[string]types.String)
	diags := diag.Diagnostics{}

	for _, field := range r.object.Fields {
		if field.Type != manifest.FieldTypeString {
			continue
		}

		var value types.String
		diags.Append(source.GetAttribute(ctx, path.Root(field.Name), &value)...)
		if value.IsUnknown() {
			continue
		}
		values[field.Name] = value
	}

	return valueSnapshot{Values: values, Diagnostics: diags}
}

func (r *objectResource) setState(
	ctx context.Context,
	state attributeTarget,
	id int64,
	apiObject map[string]any,
	writeOnlyValues map[string]types.String,
	priorStringValues map[string]types.String,
) diag.Diagnostics {
	diags := diag.Diagnostics{}
	diags.Append(state.SetAttribute(ctx, path.Root("id"), strconv.FormatInt(id, 10))...)

	for _, field := range r.object.Fields {
		if field.WriteOnly {
			preserved, ok := writeOnlyValues[field.Name]
			if ok {
				diags.Append(state.SetAttribute(ctx, path.Root(field.Name), preserved)...)
			} else {
				diags.Append(state.SetAttribute(ctx, path.Root(field.Name), types.StringNull())...)
			}
			continue
		}

		value := apiObject[field.Name]
		if normalized, ok := normalizeOptionalEmptyStringToNull(field, value, priorStringValues); ok {
			diags.Append(state.SetAttribute(ctx, path.Root(field.Name), normalized)...)
			continue
		}

		converted, convDiags := toTerraformValue(field, value)
		diags.Append(convDiags...)
		if convDiags.HasError() {
			continue
		}
		diags.Append(state.SetAttribute(ctx, path.Root(field.Name), converted)...)
	}

	return diags
}

func normalizeOptionalEmptyStringToNull(field manifest.FieldSpec, value any, priorStringValues map[string]types.String) (types.String, bool) {
	if field.Type != manifest.FieldTypeString || field.Required {
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

func newResourceFieldAttribute(field manifest.FieldSpec) resourceschema.Attribute {
	optional := !field.Required
	switch field.Type {
	case manifest.FieldTypeInt:
		return resourceschema.Int64Attribute{
			Description: fieldDescription(field),
			Required:    field.Required,
			Optional:    optional,
			Sensitive:   field.Sensitive,
		}
	case manifest.FieldTypeBool:
		return resourceschema.BoolAttribute{
			Description: fieldDescription(field),
			Required:    field.Required,
			Optional:    optional,
			Sensitive:   field.Sensitive,
		}
	case manifest.FieldTypeFloat:
		return resourceschema.Float64Attribute{
			Description: fieldDescription(field),
			Required:    field.Required,
			Optional:    optional,
			Sensitive:   field.Sensitive,
		}
	default:
		return resourceschema.StringAttribute{
			Description: fieldDescription(field),
			Required:    field.Required,
			Optional:    optional,
			Sensitive:   field.Sensitive,
		}
	}
}

func fieldDescription(field manifest.FieldSpec) string {
	if strings.TrimSpace(field.Description) != "" {
		return field.Description
	}
	if field.Type == manifest.FieldTypeArray || field.Type == manifest.FieldTypeObject {
		return "JSON-encoded value for this AWX field."
	}
	return "Managed field from AWX schema."
}

func toTerraformValue(field manifest.FieldSpec, value any) (any, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	if value == nil {
		switch field.Type {
		case manifest.FieldTypeInt:
			return types.Int64Null(), diags
		case manifest.FieldTypeBool:
			return types.BoolNull(), diags
		case manifest.FieldTypeFloat:
			return types.Float64Null(), diags
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
	case manifest.FieldTypeArray, manifest.FieldTypeObject:
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

func getResourceID(ctx context.Context, state attributeSource) (int64, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	var id types.String
	diags.Append(state.GetAttribute(ctx, path.Root("id"), &id)...)
	if id.IsUnknown() || id.IsNull() {
		diags.AddError("Missing resource ID", "Expected state to contain a numeric AWX ID.")
		return 0, diags
	}

	parsed, err := strconv.ParseInt(id.ValueString(), 10, 64)
	if err != nil {
		diags.AddError("Invalid resource ID", err.Error())
		return 0, diags
	}
	return parsed, diags
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
