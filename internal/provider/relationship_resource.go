package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

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
	_ resource.Resource                = (*relationshipResource)(nil)
	_ resource.ResourceWithConfigure   = (*relationshipResource)(nil)
	_ resource.ResourceWithImportState = (*relationshipResource)(nil)

	compositeIDPattern = regexp.MustCompile(`^([0-9]+):([0-9]+)$`)
)

type relationshipResource struct {
	relationship manifest.Relationship
	data         *configuredProvider
}

// NewRelationshipResource returns a relationship resource implementation.
func NewRelationshipResource(relationship manifest.Relationship) resource.Resource {
	return &relationshipResource{relationship: relationship}
}

func (r *relationshipResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.relationship.ResourceName
}

func (r *relationshipResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	parentIDAttribute := r.parentIDAttributeName()
	childIDAttribute := r.childIDAttributeName()

	if r.isSurveySpecRelationship() {
		resp.Schema = resourceschema.Schema{
			Description: fmt.Sprintf("Manages AWX `%s` survey specification.", r.relationship.Name),
			Attributes: map[string]resourceschema.Attribute{
				"id": resourceschema.StringAttribute{
					Description: fmt.Sprintf("Survey spec resource identifier (same as `%s`).", parentIDAttribute),
					Computed:    true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				parentIDAttribute: resourceschema.Int64Attribute{
					Description: fmt.Sprintf("Numeric ID for parent `%s` object.", r.relationship.ParentObject),
					Required:    true,
				},
				"spec": resourceschema.StringAttribute{
					Description: "JSON-encoded survey specification payload.",
					Optional:    true,
					Computed:    true,
				},
			},
		}
		return
	}

	resp.Schema = resourceschema.Schema{
		Description: fmt.Sprintf("Manages AWX `%s` relationship resources.", r.relationship.Name),
		Attributes: map[string]resourceschema.Attribute{
			"id": resourceschema.StringAttribute{
				Description: "Composite ID in <parent_id>:<child_id> format.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			parentIDAttribute: resourceschema.Int64Attribute{
				Description: fmt.Sprintf("Numeric ID for parent `%s` object.", r.relationship.ParentObject),
				Required:    true,
			},
			childIDAttribute: resourceschema.Int64Attribute{
				Description: fmt.Sprintf("Numeric ID for child `%s` object.", r.relationship.ChildObject),
				Required:    true,
			},
		},
	}
}

func (r *relationshipResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *relationshipResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.data == nil || r.data.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "Expected configured AWX client but provider data was not available.")
		return
	}
	if r.isSurveySpecRelationship() {
		parentIDAttribute := r.parentIDAttributeName()
		parentID, payload, spec, diags := surveySpecConfig(ctx, req.Plan, parentIDAttribute)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		resolvedPath := strings.Replace(r.relationship.Path, "{id}", strconv.FormatInt(parentID, 10), 1)
		if _, err := r.data.client.DoJSON(ctx, http.MethodPost, resolvedPath, nil, payload); err != nil {
			resp.Diagnostics.AddError("Failed to create survey specification", err.Error())
			return
		}

		stateSpec := spec
		if refreshed, err := r.data.client.GetObject(ctx, r.relationship.Path, strconv.FormatInt(parentID, 10)); err == nil && len(refreshed) > 0 {
			if encoded, encodeErr := json.Marshal(refreshed); encodeErr == nil {
				stateSpec = types.StringValue(string(encoded))
			}
		}

		setSurveySpecState(ctx, parentID, parentIDAttribute, stateSpec, &resp.State, &resp.Diagnostics)
		return
	}

	parentIDAttribute := r.parentIDAttributeName()
	childIDAttribute := r.childIDAttributeName()
	parentID, childID, diags := relationshipIDsFromConfig(ctx, req.Plan, parentIDAttribute, childIDAttribute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.data.client.Associate(ctx, r.relationship.Path, parentID, childID); err != nil {
		resp.Diagnostics.AddError("Failed to create relationship", err.Error())
		return
	}

	setRelationshipState(ctx, parentID, childID, parentIDAttribute, childIDAttribute, &resp.State, &resp.Diagnostics)
}

func (r *relationshipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.data == nil || r.data.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "Expected configured AWX client but provider data was not available.")
		return
	}
	if r.isSurveySpecRelationship() {
		parentIDAttribute := r.parentIDAttributeName()
		parentID, diags := surveySpecParentID(ctx, req.State, parentIDAttribute)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		payload, err := r.data.client.GetObject(ctx, r.relationship.Path, strconv.FormatInt(parentID, 10))
		if err != nil {
			if shouldRemoveFromStateOnReadError(err) {
				resp.State.RemoveResource(ctx)
				return
			}
			resp.Diagnostics.AddError("Failed to read survey specification", err.Error())
			return
		}

		var currentSpec types.String
		resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("spec"), &currentSpec)...)
		if resp.Diagnostics.HasError() {
			return
		}

		stateSpec := currentSpec
		if len(payload) > 0 {
			encoded, encodeErr := json.Marshal(payload)
			if encodeErr != nil {
				resp.Diagnostics.AddError("Failed to encode survey specification state", encodeErr.Error())
				return
			}
			stateSpec = types.StringValue(string(encoded))
		} else if stateSpec.IsUnknown() {
			stateSpec = types.StringNull()
		}

		setSurveySpecState(ctx, parentID, parentIDAttribute, stateSpec, &resp.State, &resp.Diagnostics)
		return
	}

	parentIDAttribute := r.parentIDAttributeName()
	childIDAttribute := r.childIDAttributeName()
	parentID, childID, diags := relationshipIDsFromConfig(ctx, req.State, parentIDAttribute, childIDAttribute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	exists, err := r.data.client.RelationshipExists(ctx, r.relationship.Path, parentID, childID)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read relationship", err.Error())
		return
	}
	if !exists {
		resp.State.RemoveResource(ctx)
		return
	}

	setRelationshipState(ctx, parentID, childID, parentIDAttribute, childIDAttribute, &resp.State, &resp.Diagnostics)
}

func (r *relationshipResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.data == nil || r.data.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "Expected configured AWX client but provider data was not available.")
		return
	}
	if r.isSurveySpecRelationship() {
		parentIDAttribute := r.parentIDAttributeName()
		parentID, payload, spec, diags := surveySpecConfig(ctx, req.Plan, parentIDAttribute)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		resolvedPath := strings.Replace(r.relationship.Path, "{id}", strconv.FormatInt(parentID, 10), 1)
		if _, err := r.data.client.DoJSON(ctx, http.MethodPost, resolvedPath, nil, payload); err != nil {
			resp.Diagnostics.AddError("Failed to update survey specification", err.Error())
			return
		}

		setSurveySpecState(ctx, parentID, parentIDAttribute, spec, &resp.State, &resp.Diagnostics)
		return
	}

	parentIDAttribute := r.parentIDAttributeName()
	childIDAttribute := r.childIDAttributeName()
	oldParentID, oldChildID, diags := relationshipIDsFromConfig(ctx, req.State, parentIDAttribute, childIDAttribute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newParentID, newChildID, planDiags := relationshipIDsFromConfig(ctx, req.Plan, parentIDAttribute, childIDAttribute)
	resp.Diagnostics.Append(planDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if oldParentID != newParentID || oldChildID != newChildID {
		if err := r.data.client.Disassociate(ctx, r.relationship.Path, oldParentID, oldChildID); err != nil {
			resp.Diagnostics.AddError("Failed to remove previous relationship", err.Error())
			return
		}
		if err := r.data.client.Associate(ctx, r.relationship.Path, newParentID, newChildID); err != nil {
			resp.Diagnostics.AddError("Failed to create updated relationship", err.Error())
			return
		}
	}

	setRelationshipState(ctx, newParentID, newChildID, parentIDAttribute, childIDAttribute, &resp.State, &resp.Diagnostics)
}

func (r *relationshipResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.data == nil || r.data.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "Expected configured AWX client but provider data was not available.")
		return
	}
	if r.isSurveySpecRelationship() {
		parentIDAttribute := r.parentIDAttributeName()
		parentID, diags := surveySpecParentID(ctx, req.State, parentIDAttribute)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		resolvedPath := strings.Replace(r.relationship.Path, "{id}", strconv.FormatInt(parentID, 10), 1)
		_, err := r.data.client.DoJSON(ctx, http.MethodDelete, resolvedPath, nil, nil)
		if apiErr := asAPIError(err); apiErr != nil && apiErr.StatusCode == http.StatusNotFound {
			return
		}
		if err != nil {
			resp.Diagnostics.AddError("Failed to delete survey specification", err.Error())
		}
		return
	}

	parentID, childID, diags := relationshipIDsFromConfig(ctx, req.State, r.parentIDAttributeName(), r.childIDAttributeName())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.data.client.Disassociate(ctx, r.relationship.Path, parentID, childID); err != nil {
		resp.Diagnostics.AddError("Failed to delete relationship", err.Error())
		return
	}
}

func (r *relationshipResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if r.isSurveySpecRelationship() {
		parentID, err := parseSurveySpecImportID(req.ID)
		if err != nil {
			resp.Diagnostics.AddError("Invalid survey specification import ID", err.Error())
			return
		}
		setSurveySpecState(ctx, parentID, r.parentIDAttributeName(), types.StringNull(), &resp.State, &resp.Diagnostics)
		return
	}

	parentID, childID, err := parseCompositeRelationshipImportID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid relationship import ID", err.Error())
		return
	}
	setRelationshipState(ctx, parentID, childID, r.parentIDAttributeName(), r.childIDAttributeName(), &resp.State, &resp.Diagnostics)
}

func parseSurveySpecImportID(rawID string) (int64, error) {
	identifier := strings.TrimSpace(rawID)
	parentID, err := strconv.ParseInt(identifier, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Use <parent_id>, for example 12.")
	}
	return parentID, nil
}

func parseCompositeRelationshipImportID(rawID string) (int64, int64, error) {
	matches := compositeIDPattern.FindStringSubmatch(strings.TrimSpace(rawID))
	if len(matches) != 3 {
		return 0, 0, fmt.Errorf("Use <parent_id>:<child_id>, for example 12:34.")
	}

	parentID, _ := strconv.ParseInt(matches[1], 10, 64)
	childID, _ := strconv.ParseInt(matches[2], 10, 64)
	return parentID, childID, nil
}

func relationshipIDsFromConfig(ctx context.Context, source attributeSource, parentAttribute string, childAttribute string) (int64, int64, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	var parentID types.Int64
	diags.Append(source.GetAttribute(ctx, path.Root(parentAttribute), &parentID)...)
	var childID types.Int64
	diags.Append(source.GetAttribute(ctx, path.Root(childAttribute), &childID)...)

	if parentID.IsNull() || parentID.IsUnknown() {
		diags.AddAttributeError(path.Root(parentAttribute), "Missing required parent ID", fmt.Sprintf("%s is required.", parentAttribute))
	}
	if childID.IsNull() || childID.IsUnknown() {
		diags.AddAttributeError(path.Root(childAttribute), "Missing required child ID", fmt.Sprintf("%s is required.", childAttribute))
	}
	if diags.HasError() {
		return 0, 0, diags
	}

	return parentID.ValueInt64(), childID.ValueInt64(), diags
}

func surveySpecParentID(ctx context.Context, source attributeSource, parentAttribute string) (int64, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	var parentID types.Int64
	diags.Append(source.GetAttribute(ctx, path.Root(parentAttribute), &parentID)...)

	if parentID.IsNull() || parentID.IsUnknown() {
		diags.AddAttributeError(path.Root(parentAttribute), "Missing required parent ID", fmt.Sprintf("%s is required.", parentAttribute))
		return 0, diags
	}
	return parentID.ValueInt64(), diags
}

func surveySpecConfig(ctx context.Context, source attributeSource, parentAttribute string) (int64, any, types.String, diag.Diagnostics) {
	parentID, diags := surveySpecParentID(ctx, source, parentAttribute)
	if diags.HasError() {
		return 0, nil, types.StringNull(), diags
	}

	var spec types.String
	diags.Append(source.GetAttribute(ctx, path.Root("spec"), &spec)...)
	if spec.IsNull() || spec.IsUnknown() {
		diags.AddAttributeError(path.Root("spec"), "Missing survey specification", "spec is required for survey specification resources.")
		return 0, nil, types.StringNull(), diags
	}

	decoded, err := decodeJSONString(spec.ValueString())
	if err != nil {
		diags.AddAttributeError(path.Root("spec"), "Invalid JSON payload", err.Error())
		return 0, nil, types.StringNull(), diags
	}

	encoded, err := json.Marshal(decoded)
	if err != nil {
		diags.AddAttributeError(path.Root("spec"), "Invalid JSON payload", err.Error())
		return 0, nil, types.StringNull(), diags
	}

	return parentID, decoded, types.StringValue(string(encoded)), diags
}

func setRelationshipState(ctx context.Context, parentID, childID int64, parentAttribute string, childAttribute string, target attributeTarget, diags *diag.Diagnostics) {
	compositeID := fmt.Sprintf("%d:%d", parentID, childID)
	diags.Append(target.SetAttribute(ctx, path.Root("id"), compositeID)...)
	diags.Append(target.SetAttribute(ctx, path.Root(parentAttribute), types.Int64Value(parentID))...)
	diags.Append(target.SetAttribute(ctx, path.Root(childAttribute), types.Int64Value(childID))...)
}

func setSurveySpecState(ctx context.Context, parentID int64, parentAttribute string, spec types.String, target attributeTarget, diags *diag.Diagnostics) {
	diags.Append(target.SetAttribute(ctx, path.Root("id"), fmt.Sprintf("%d", parentID))...)
	diags.Append(target.SetAttribute(ctx, path.Root(parentAttribute), types.Int64Value(parentID))...)
	diags.Append(target.SetAttribute(ctx, path.Root("spec"), spec)...)
}

func (r *relationshipResource) isSurveySpecRelationship() bool {
	return strings.HasSuffix(r.relationship.Path, "/survey_spec/")
}

func (r *relationshipResource) parentIDAttributeName() string {
	return manifest.RelationshipParentIDAttribute(r.relationship)
}

func (r *relationshipResource) childIDAttributeName() string {
	return manifest.RelationshipChildIDAttribute(r.relationship)
}
