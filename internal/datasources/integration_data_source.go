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
	_ datasource.DataSource              = &integrationDataSource{}
	_ datasource.DataSourceWithConfigure = &integrationDataSource{}
)

func NewIntegrationDataSource() datasource.DataSource {
	return &integrationDataSource{}
}

type integrationDataSource struct {
	client *client.Client
}

type integrationDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
	BotID       types.String `tfsdk:"bot_id"`
	CreatedAt   types.Int64  `tfsdk:"created_at"`
	UpdatedAt   types.Int64  `tfsdk:"updated_at"`
}

func (d *integrationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration"
}

func (d *integrationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a ChatBotKit integration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Integration identifier",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Integration name",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Integration description",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "Integration type",
				Computed:    true,
			},
			"bot_id": schema.StringAttribute{
				Description: "Bot ID",
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

func (d *integrationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *integrationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data integrationDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	integration, err := d.client.GetIntegration(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading integration", err.Error())
		return
	}

	data.Name = types.StringValue(integration.Name)
	data.Description = types.StringValue(integration.Description)
	data.Type = types.StringValue(integration.Type)
	data.BotID = types.StringValue(integration.BotID)
	data.CreatedAt = types.Int64Value(integration.CreatedAt)
	data.UpdatedAt = types.Int64Value(integration.UpdatedAt)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
