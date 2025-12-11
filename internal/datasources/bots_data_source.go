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
	_ datasource.DataSource              = &botsDataSource{}
	_ datasource.DataSourceWithConfigure = &botsDataSource{}
)

func NewBotsDataSource() datasource.DataSource {
	return &botsDataSource{}
}

type botsDataSource struct {
	client *client.Client
}

type botsDataSourceModel struct {
	Bots []botDataSourceModel `tfsdk:"bots"`
}

func (d *botsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bots"
}

func (d *botsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists ChatBotKit bots.",
		Attributes: map[string]schema.Attribute{
			"bots": schema.ListNestedAttribute{
				Description: "List of bots",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Bot identifier",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Bot name",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Bot description",
							Computed:    true,
						},
						"model": schema.StringAttribute{
							Description: "AI model",
							Computed:    true,
						},
						"dataset_id": schema.StringAttribute{
							Description: "Dataset ID",
							Computed:    true,
						},
						"skillset_id": schema.StringAttribute{
							Description: "Skillset ID",
							Computed:    true,
						},
						"backstory": schema.StringAttribute{
							Description: "Bot backstory",
							Computed:    true,
						},
						"temperature": schema.Float64Attribute{
							Description: "Temperature",
							Computed:    true,
						},
						"instructions": schema.StringAttribute{
							Description: "Instructions",
							Computed:    true,
						},
						"moderation": schema.BoolAttribute{
							Description: "Moderation enabled",
							Computed:    true,
						},
						"privacy": schema.BoolAttribute{
							Description: "Privacy enabled",
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

func (d *botsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *botsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data botsDataSourceModel

	result, err := d.client.ListBots(ctx, "")
	if err != nil {
		resp.Diagnostics.AddError("Error listing bots", err.Error())
		return
	}

	for _, bot := range result.Items {
		botModel := botDataSourceModel{
			ID:           types.StringValue(bot.ID),
			Name:         types.StringValue(bot.Name),
			Description:  types.StringValue(bot.Description),
			Model:        types.StringValue(bot.Model),
			DatasetID:    types.StringValue(bot.DatasetID),
			SkillsetID:   types.StringValue(bot.SkillsetID),
			Backstory:    types.StringValue(bot.Backstory),
			Temperature:  types.Float64Value(bot.Temperature),
			Instructions: types.StringValue(bot.Instructions),
			Moderation:   types.BoolValue(bot.Moderation),
			Privacy:      types.BoolValue(bot.Privacy),
			CreatedAt:    types.Int64Value(bot.CreatedAt),
			UpdatedAt:    types.Int64Value(bot.UpdatedAt),
		}
		data.Bots = append(data.Bots, botModel)
	}

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
