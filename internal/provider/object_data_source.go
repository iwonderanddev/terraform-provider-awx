package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/damien/terraform-awx-provider/internal/manifest"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = (*objectDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*objectDataSource)(nil)
)

type objectDataSource struct {
	object manifest.ManagedObject
	data   *configuredProvider
}

type objectLookupClient interface {
	GetObject(context.Context, string, string) (map[string]any, error)
	FindByField(context.Context, string, string, string) ([]map[string]any, error)
}

type dataSourceLookupInput struct {
	NumericID    types.Int64
	Identifier   types.String
	Name         types.String
	HasNameField bool
}

// NewObjectDataSource returns a generated data source for one AWX object.
func NewObjectDataSource(object manifest.ManagedObject) datasource.DataSource {
	return &objectDataSource{object: object}
}

func (d *objectDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = d.object.DataSourceName
}

func (d *objectDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	attributes := map[string]datasourceschema.Attribute{}
	if d.object.CollectionCreate {
		attributes["id"] = datasourceschema.Int64Attribute{
			Description: "Numeric AWX object identifier for deterministic lookup.",
			Optional:    true,
			Computed:    true,
		}
	} else {
		attributes["id"] = datasourceschema.StringAttribute{
			Description: "AWX object identifier for deterministic lookup.",
			Optional:    true,
			Computed:    true,
		}
	}
	if hasManifestField(d.object.Fields, "name") {
		attributes["name"] = datasourceschema.StringAttribute{
			Description: "Exact object name lookup. Must uniquely resolve to a single result when id is not provided.",
			Optional:    true,
		}
	}

	for _, field := range d.object.Fields {
		tfName := manifest.TerraformAttributeName(d.object.Name, field.Name)
		if _, exists := attributes[tfName]; exists {
			// Preserve explicit lookup arguments such as id and name.
			continue
		}
		attributes[tfName] = newDataSourceFieldAttribute(field)
	}

	resp.Schema = datasourceschema.Schema{
		Description: fmt.Sprintf("Reads AWX `%s` objects.", d.object.Name),
		Attributes:  attributes,
	}
}

func (d *objectDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	providerData, ok := req.ProviderData.(*configuredProvider)
	if !ok {
		resp.Diagnostics.AddError("Unexpected data source configure type", fmt.Sprintf("expected *configuredProvider, got %T", req.ProviderData))
		return
	}
	d.data = providerData
}

func (d *objectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.data == nil || d.data.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "Expected configured AWX client but provider data was not available.")
		return
	}

	lookup := dataSourceLookupInput{
		Identifier: types.StringNull(),
		NumericID:  types.Int64Null(),
	}
	if d.object.CollectionCreate {
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("id"), &lookup.NumericID)...)
	} else {
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("id"), &lookup.Identifier)...)
	}

	var nameValue types.String
	hasNameField := hasManifestField(d.object.Fields, "name")
	if hasNameField {
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("name"), &nameValue)...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	lookup.Name = nameValue
	lookup.HasNameField = hasNameField
	target, lookupDiags := resolveObjectDataSourceTarget(ctx, d.data.client, d.object, lookup)
	resp.Diagnostics.Append(lookupDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(d.setState(ctx, &resp.State, target)...)
}

func (d *objectDataSource) setState(ctx context.Context, state attributeTarget, target objectLookupResult) diag.Diagnostics {
	diags := diag.Diagnostics{}
	if d.object.CollectionCreate {
		id, err := parseNumericID(target.ID)
		if err != nil {
			diags.AddError("Failed to set AWX object ID", err.Error())
			return diags
		}
		diags.Append(state.SetAttribute(ctx, path.Root("id"), id)...)
	} else {
		diags.Append(state.SetAttribute(ctx, path.Root("id"), target.ID)...)
	}

	for _, field := range d.object.Fields {
		tfName := manifest.TerraformAttributeName(d.object.Name, field.Name)
		if field.WriteOnly {
			nullValue, nullDiags := toTerraformValue(d.object.Name, field, nil)
			diags.Append(nullDiags...)
			if nullDiags.HasError() {
				continue
			}
			diags.Append(state.SetAttribute(ctx, path.Root(tfName), nullValue)...)
			continue
		}

		value, valueDiags := toTerraformValue(d.object.Name, field, target.Object[field.Name])
		diags.Append(valueDiags...)
		if valueDiags.HasError() {
			continue
		}
		diags.Append(state.SetAttribute(ctx, path.Root(tfName), value)...)
	}

	return diags
}

func newDataSourceFieldAttribute(field manifest.FieldSpec) datasourceschema.Attribute {
	switch field.Type {
	case manifest.FieldTypeInt:
		return datasourceschema.Int64Attribute{
			Description: fieldDescription(field),
			Computed:    true,
			Sensitive:   field.Sensitive,
		}
	case manifest.FieldTypeBool:
		return datasourceschema.BoolAttribute{
			Description: fieldDescription(field),
			Computed:    true,
			Sensitive:   field.Sensitive,
		}
	case manifest.FieldTypeFloat:
		return datasourceschema.Float64Attribute{
			Description: fieldDescription(field),
			Computed:    true,
			Sensitive:   field.Sensitive,
		}
	case manifest.FieldTypeObject:
		return datasourceschema.DynamicAttribute{
			Description: fieldDescription(field),
			Computed:    true,
			Sensitive:   field.Sensitive,
		}
	default:
		return datasourceschema.StringAttribute{
			Description: fieldDescription(field),
			Computed:    true,
			Sensitive:   field.Sensitive,
		}
	}
}

func hasManifestField(fields []manifest.FieldSpec, name string) bool {
	for _, field := range fields {
		if field.Name == name {
			return true
		}
	}
	return false
}

func copyDiags(in diag.Diagnostics) diag.Diagnostics {
	out := make(diag.Diagnostics, len(in))
	copy(out, in)
	return out
}

type objectLookupResult struct {
	Object map[string]any
	ID     string
}

func resolveObjectDataSourceTarget(ctx context.Context, api objectLookupClient, object manifest.ManagedObject, lookup dataSourceLookupInput) (objectLookupResult, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	if object.CollectionCreate && !lookup.NumericID.IsNull() && !lookup.NumericID.IsUnknown() {
		id := strconv.FormatInt(lookup.NumericID.ValueInt64(), 10)
		obj, err := api.GetObject(ctx, object.DetailPath, id)
		if err != nil {
			diags.AddError("Failed to query AWX object", err.Error())
			return objectLookupResult{}, diags
		}
		return objectLookupResult{Object: obj, ID: id}, diags
	}
	if !object.CollectionCreate && !lookup.Identifier.IsNull() && !lookup.Identifier.IsUnknown() {
		id := strings.TrimSpace(lookup.Identifier.ValueString())
		if id == "" {
			diags.AddAttributeError(path.Root("id"), "Invalid AWX object ID", "id cannot be empty.")
			return objectLookupResult{}, diags
		}

		obj, err := api.GetObject(ctx, object.DetailPath, id)
		if err != nil {
			diags.AddError("Failed to query AWX object", err.Error())
			return objectLookupResult{}, diags
		}
		return objectLookupResult{Object: obj, ID: id}, diags
	}

	if lookup.HasNameField && !lookup.Name.IsNull() && !lookup.Name.IsUnknown() {
		matches, err := api.FindByField(ctx, object.CollectionPath, "name", lookup.Name.ValueString())
		if err != nil {
			diags.AddError("Failed to query AWX object by name", err.Error())
			return objectLookupResult{}, diags
		}
		if len(matches) == 0 {
			diags.AddError(
				"AWX object not found",
				fmt.Sprintf("No `%s` object matched name %q.", object.Name, lookup.Name.ValueString()),
			)
			return objectLookupResult{}, diags
		}
		if len(matches) > 1 {
			diags.AddError(
				"Ambiguous AWX object lookup",
				fmt.Sprintf("Lookup for `%s` matched %d objects. Refine the query or provide id.", object.Name, len(matches)),
			)
			return objectLookupResult{}, diags
		}

		if object.CollectionCreate {
			id, err := parseNumericID(matches[0]["id"])
			if err != nil {
				diags.AddError("Failed to parse AWX object ID", err.Error())
				return objectLookupResult{}, diags
			}
			return objectLookupResult{Object: matches[0], ID: strconv.FormatInt(id, 10)}, diags
		}

		id, err := parseAPIObjectID(matches[0]["id"])
		if err != nil {
			diags.AddError("Failed to parse AWX object ID", err.Error())
			return objectLookupResult{}, diags
		}
		return objectLookupResult{Object: matches[0], ID: id}, diags
	}

	diags.AddError(
		"Missing lookup input",
		"Provide either id or name for deterministic lookup.",
	)
	return objectLookupResult{}, diags
}

func parseAPIObjectID(value any) (string, error) {
	switch typed := value.(type) {
	case string:
		identifier := strings.TrimSpace(typed)
		if identifier == "" {
			return "", fmt.Errorf("ID is empty")
		}
		return identifier, nil
	default:
		id, err := parseNumericID(value)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%d", id), nil
	}
}
