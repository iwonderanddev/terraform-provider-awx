package provider

import (
	"context"
	"fmt"
	"strconv"

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

// NewObjectDataSource returns a generated data source for one AWX object.
func NewObjectDataSource(object manifest.ManagedObject) datasource.DataSource {
	return &objectDataSource{object: object}
}

func (d *objectDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = d.object.DataSourceName
}

func (d *objectDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	attributes := map[string]datasourceschema.Attribute{
		"id": datasourceschema.StringAttribute{
			Description: "Numeric AWX object ID for deterministic lookup.",
			Optional:    true,
			Computed:    true,
		},
	}
	if hasManifestField(d.object.Fields, "name") {
		attributes["name"] = datasourceschema.StringAttribute{
			Description: "Exact object name lookup. Must uniquely resolve to a single result when id is not provided.",
			Optional:    true,
		}
	}

	for _, field := range d.object.Fields {
		attributes[field.Name] = newDataSourceFieldAttribute(field)
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

	var idValue types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("id"), &idValue)...)

	var nameValue types.String
	hasNameField := hasManifestField(d.object.Fields, "name")
	if hasNameField {
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("name"), &nameValue)...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	var target map[string]any
	if !idValue.IsNull() && !idValue.IsUnknown() {
		id, err := strconv.ParseInt(idValue.ValueString(), 10, 64)
		if err != nil {
			resp.Diagnostics.AddAttributeError(path.Root("id"), "Invalid AWX object ID", err.Error())
			return
		}

		obj, err := d.data.client.GetObject(ctx, d.object.DetailPath, id)
		if err != nil {
			resp.Diagnostics.AddError("Failed to query AWX object", err.Error())
			return
		}
		target = obj
	} else if hasNameField && !nameValue.IsNull() && !nameValue.IsUnknown() {
		matches, err := d.data.client.FindByField(ctx, d.object.CollectionPath, "name", nameValue.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Failed to query AWX object by name", err.Error())
			return
		}
		if len(matches) == 0 {
			resp.Diagnostics.AddError(
				"AWX object not found",
				fmt.Sprintf("No `%s` object matched name %q.", d.object.Name, nameValue.ValueString()),
			)
			return
		}
		if len(matches) > 1 {
			resp.Diagnostics.AddError(
				"Ambiguous AWX object lookup",
				fmt.Sprintf("Lookup for `%s` matched %d objects. Refine the query or provide id.", d.object.Name, len(matches)),
			)
			return
		}
		target = matches[0]
	} else {
		resp.Diagnostics.AddError(
			"Missing lookup input",
			"Provide either id or name for deterministic lookup.",
		)
		return
	}

	id, err := parseNumericID(target["id"])
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse AWX object ID", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), strconv.FormatInt(id, 10))...)
	for _, field := range d.object.Fields {
		if field.WriteOnly {
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(field.Name), types.StringNull())...)
			continue
		}
		value, diags := toTerraformValue(field, target[field.Name])
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			continue
		}
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(field.Name), value)...)
	}
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
