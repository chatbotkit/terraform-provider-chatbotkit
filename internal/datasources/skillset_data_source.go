package datasources

import (
	"context"
	"fmt"

	"github.com/chatbotkit/terraform-provider/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &skillsetDataSource{}
	_ datasource.DataSourceWithConfigure = &skillsetDataSource{}
)

func NewSkillsetDataSource() datasource.DataSource {
	return &skillsetDataSource{}
}

type skillsetDataSource struct {
	client *client.Client
}

type skillsetDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	CreatedAt   types.Int64  `tfsdk:"created_at"`
	UpdatedAt   types.Int64  `tfsdk:"updated_at"`
}

func (d *skillsetDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_skillset"
}

func (d *skillsetDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a ChatBotKit skillset.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Skillset identifier",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Skillset name",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Skillset description",
				Computed:    true,
			},
			"created_at": schema.Int64Attribute{
				Description: "Creation timestamp",
				Computed:    true,
			},
			"updated_at": schema.Int64Attribute{
				Description: "Last update timestamp",
				Computed:    true,
			},
		},
	}
}

func (d *skillsetDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *skillsetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data skillsetDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	skillset, err := d.client.GetSkillset(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading skillset", err.Error())
		return
	}

	data.Name = types.StringValue(skillset.Name)
	data.Description = types.StringValue(skillset.Description)
	data.CreatedAt = types.Int64Value(skillset.CreatedAt)
	data.UpdatedAt = types.Int64Value(skillset.UpdatedAt)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
