package provider

import (
	"context"
	"fmt"
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
			"parent_id": resourceschema.Int64Attribute{
				Description: fmt.Sprintf("Numeric ID for parent `%s` object.", r.relationship.ParentObject),
				Required:    true,
			},
			"child_id": resourceschema.Int64Attribute{
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

	parentID, childID, diags := relationshipIDsFromConfig(ctx, req.Plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.data.client.Associate(ctx, r.relationship.Path, parentID, childID); err != nil {
		resp.Diagnostics.AddError("Failed to create relationship", err.Error())
		return
	}

	setRelationshipState(ctx, parentID, childID, &resp.State, &resp.Diagnostics)
}

func (r *relationshipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.data == nil || r.data.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "Expected configured AWX client but provider data was not available.")
		return
	}

	parentID, childID, diags := relationshipIDsFromConfig(ctx, req.State)
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

	setRelationshipState(ctx, parentID, childID, &resp.State, &resp.Diagnostics)
}

func (r *relationshipResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.data == nil || r.data.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "Expected configured AWX client but provider data was not available.")
		return
	}

	oldParentID, oldChildID, diags := relationshipIDsFromConfig(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newParentID, newChildID, planDiags := relationshipIDsFromConfig(ctx, req.Plan)
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

	setRelationshipState(ctx, newParentID, newChildID, &resp.State, &resp.Diagnostics)
}

func (r *relationshipResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.data == nil || r.data.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "Expected configured AWX client but provider data was not available.")
		return
	}

	parentID, childID, diags := relationshipIDsFromConfig(ctx, req.State)
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
	matches := compositeIDPattern.FindStringSubmatch(strings.TrimSpace(req.ID))
	if len(matches) != 3 {
		resp.Diagnostics.AddError("Invalid relationship import ID", "Use <parent_id>:<child_id>, for example 12:34.")
		return
	}

	parentID, _ := strconv.ParseInt(matches[1], 10, 64)
	childID, _ := strconv.ParseInt(matches[2], 10, 64)
	setRelationshipState(ctx, parentID, childID, &resp.State, &resp.Diagnostics)
}

func relationshipIDsFromConfig(ctx context.Context, source attributeSource) (int64, int64, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	var parentID types.Int64
	diags.Append(source.GetAttribute(ctx, path.Root("parent_id"), &parentID)...)
	var childID types.Int64
	diags.Append(source.GetAttribute(ctx, path.Root("child_id"), &childID)...)

	if parentID.IsNull() || parentID.IsUnknown() {
		diags.AddAttributeError(path.Root("parent_id"), "Missing parent ID", "parent_id is required.")
	}
	if childID.IsNull() || childID.IsUnknown() {
		diags.AddAttributeError(path.Root("child_id"), "Missing child ID", "child_id is required.")
	}
	if diags.HasError() {
		return 0, 0, diags
	}

	return parentID.ValueInt64(), childID.ValueInt64(), diags
}

func setRelationshipState(ctx context.Context, parentID, childID int64, target attributeTarget, diags *diag.Diagnostics) {
	compositeID := fmt.Sprintf("%d:%d", parentID, childID)
	diags.Append(target.SetAttribute(ctx, path.Root("id"), compositeID)...)
	diags.Append(target.SetAttribute(ctx, path.Root("parent_id"), types.Int64Value(parentID))...)
	diags.Append(target.SetAttribute(ctx, path.Root("child_id"), types.Int64Value(childID))...)
}
