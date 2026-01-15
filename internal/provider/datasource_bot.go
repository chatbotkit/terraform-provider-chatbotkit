package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &BotDataSource{}

func NewBotDataSource() datasource.DataSource {
	return &BotDataSource{}
}

// BotDataSource defines the data source implementation.
type BotDataSource struct {
	client *Client
}

// BotDataSourceModel describes the data source data model.
type BotDataSourceModel struct {
	ID types.String `tfsdk:"id"`

	Backstory types.String `tfsdk:"backstory"`
	BlueprintId types.String `tfsdk:"blueprint_id"`
	DatasetId types.String `tfsdk:"dataset_id"`
	Description types.String `tfsdk:"description"`
	Meta types.Map `tfsdk:"meta"`
	Model types.String `tfsdk:"model"`
	Moderation types.Bool `tfsdk:"moderation"`
	Name types.String `tfsdk:"name"`
	Privacy types.Bool `tfsdk:"privacy"`
	SkillsetId types.String `tfsdk:"skillset_id"`
	Visibility types.String `tfsdk:"visibility"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

// Metadata returns the data source type name.
func (d *BotDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bot"
}

// Schema defines the schema for the data source.
func (d *BotDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about an existing bot.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique identifier of the bot to look up",
			},

			"backstory": schema.StringAttribute{
				MarkdownDescription: "The backstory for the bot",
				Computed:            true,
			},
			"blueprint_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the blueprint to use",
				Computed:            true,
			},
			"dataset_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the dataset to use",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the bot",
				Computed:            true,
			},
			"meta": schema.MapAttribute{
				MarkdownDescription: "Additional metadata for the bot",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"model": schema.StringAttribute{
				MarkdownDescription: "The AI model to use for the bot",
				Computed:            true,
			},
			"moderation": schema.BoolAttribute{
				MarkdownDescription: "Whether moderation is enabled",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the bot",
				Computed:            true,
			},
			"privacy": schema.BoolAttribute{
				MarkdownDescription: "Whether privacy mode is enabled",
				Computed:            true,
			},
			"skillset_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the skillset to use",
				Computed:            true,
			},
			"visibility": schema.StringAttribute{
				MarkdownDescription: "The visibility level of the bot",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "Timestamp when the resource was created",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "Timestamp when the resource was last updated",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *BotDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

// Read refreshes the Terraform state with the latest data.
func (d *BotDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data BotDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to read bot
	result, err := d.client.GetBot(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read bot: %s", err))
		return
	}

	// Update data model with response values

	if result.Backstory != nil {
		data.Backstory = types.StringPointerValue(result.Backstory)
	}
	if result.BlueprintId != nil {
		data.BlueprintId = types.StringPointerValue(result.BlueprintId)
	}
	if result.DatasetId != nil {
		data.DatasetId = types.StringPointerValue(result.DatasetId)
	}
	if result.Description != nil {
		data.Description = types.StringPointerValue(result.Description)
	}
	if result.Meta != nil {
		mapValue, diags := types.MapValueFrom(ctx, types.StringType, result.Meta)
		resp.Diagnostics.Append(diags...)
		data.Meta = mapValue
	}
	if result.Model != nil {
		data.Model = types.StringPointerValue(result.Model)
	}
	if result.Moderation != nil {
		data.Moderation = types.BoolPointerValue(result.Moderation)
	}
	if result.Name != nil {
		data.Name = types.StringPointerValue(result.Name)
	}
	if result.Privacy != nil {
		data.Privacy = types.BoolPointerValue(result.Privacy)
	}
	if result.SkillsetId != nil {
		data.SkillsetId = types.StringPointerValue(result.SkillsetId)
	}
	if result.Visibility != nil {
		data.Visibility = types.StringPointerValue(result.Visibility)
	}
	if result.CreatedAt != nil {
		data.CreatedAt = types.StringPointerValue(result.CreatedAt)
	}
	if result.UpdatedAt != nil {
		data.UpdatedAt = types.StringPointerValue(result.UpdatedAt)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
