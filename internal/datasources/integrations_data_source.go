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
	_ datasource.DataSource              = &integrationsDataSource{}
	_ datasource.DataSourceWithConfigure = &integrationsDataSource{}
)

func NewIntegrationsDataSource() datasource.DataSource {
	return &integrationsDataSource{}
}

type integrationsDataSource struct {
	client *client.Client
}

type integrationsDataSourceModel struct {
	Integrations []integrationDataSourceModel `tfsdk:"integrations"`
}

func (d *integrationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integrations"
}

func (d *integrationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists ChatBotKit integrations.",
		Attributes: map[string]schema.Attribute{
			"integrations": schema.ListNestedAttribute{
				Description: "List of integrations",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Integration identifier",
							Computed:    true,
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
				},
			},
		},
	}
}

func (d *integrationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *integrationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data integrationsDataSourceModel

	result, err := d.client.ListIntegrations(ctx, "")
	if err != nil {
		resp.Diagnostics.AddError("Error listing integrations", err.Error())
		return
	}

	for _, integration := range result.Items {
		integrationModel := integrationDataSourceModel{
			ID:          types.StringValue(integration.ID),
			Name:        types.StringValue(integration.Name),
			Description: types.StringValue(integration.Description),
			Type:        types.StringValue(integration.Type),
			BotID:       types.StringValue(integration.BotID),
			CreatedAt:   types.Int64Value(integration.CreatedAt),
			UpdatedAt:   types.Int64Value(integration.UpdatedAt),
		}
		data.Integrations = append(data.Integrations, integrationModel)
	}

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
