package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = &SystemVersionDataSource{}
	_ datasource.DataSourceWithConfigure = &SystemVersionDataSource{}
)

func NewSystemVersionDataSource() datasource.DataSource {
	return &SystemVersionDataSource{}
}

type SystemVersionDataSource struct {
	client *pfsense.Client
}

type SystemVersionDataSourceModel struct {
	Current types.String `tfsdk:"current"`
	Latest  types.String `tfsdk:"latest"`
}

func (d *SystemVersionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system_version"
}

func (d *SystemVersionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves current and latest system version.",
		Attributes: map[string]schema.Attribute{
			"current": schema.StringAttribute{
				Description: "Current pfSense system version.",
				Computed:    true,
			},
			"latest": schema.StringAttribute{
				Description: "Latest pfSense system version.",
				Computed:    true,
			},
		},
	}
}

func (d *SystemVersionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*pfsense.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *pfsense.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *SystemVersionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SystemVersionDataSourceModel

	version, err := d.client.GetSystemVersion(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get system version",
			err.Error(),
		)
		return
	}

	data.Current = types.StringValue(version.Current)
	data.Latest = types.StringValue(version.Latest)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
