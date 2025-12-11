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
	_ datasource.DataSource              = &botDataSource{}
	_ datasource.DataSourceWithConfigure = &botDataSource{}
)

func NewBotDataSource() datasource.DataSource {
	return &botDataSource{}
}

type botDataSource struct {
	client *client.Client
}

type botDataSourceModel struct {
	ID           types.String  `tfsdk:"id"`
	Name         types.String  `tfsdk:"name"`
	Description  types.String  `tfsdk:"description"`
	Model        types.String  `tfsdk:"model"`
	DatasetID    types.String  `tfsdk:"dataset_id"`
	SkillsetID   types.String  `tfsdk:"skillset_id"`
	Backstory    types.String  `tfsdk:"backstory"`
	Temperature  types.Float64 `tfsdk:"temperature"`
	Instructions types.String  `tfsdk:"instructions"`
	Moderation   types.Bool    `tfsdk:"moderation"`
	Privacy      types.Bool    `tfsdk:"privacy"`
	CreatedAt    types.Int64   `tfsdk:"created_at"`
	UpdatedAt    types.Int64   `tfsdk:"updated_at"`
}

func (d *botDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bot"
}

func (d *botDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a ChatBotKit bot.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Bot identifier",
				Required:    true,
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
	}
}

func (d *botDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *botDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data botDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	bot, err := d.client.GetBot(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading bot", err.Error())
		return
	}

	data.Name = types.StringValue(bot.Name)
	data.Description = types.StringValue(bot.Description)
	data.Model = types.StringValue(bot.Model)
	data.DatasetID = types.StringValue(bot.DatasetID)
	data.SkillsetID = types.StringValue(bot.SkillsetID)
	data.Backstory = types.StringValue(bot.Backstory)
	data.Temperature = types.Float64Value(bot.Temperature)
	data.Instructions = types.StringValue(bot.Instructions)
	data.Moderation = types.BoolValue(bot.Moderation)
	data.Privacy = types.BoolValue(bot.Privacy)
	data.CreatedAt = types.Int64Value(bot.CreatedAt)
	data.UpdatedAt = types.Int64Value(bot.UpdatedAt)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
