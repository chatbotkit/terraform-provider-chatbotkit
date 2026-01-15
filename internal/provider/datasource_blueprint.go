package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &BlueprintDataSource{}

func NewBlueprintDataSource() datasource.DataSource {
	return &BlueprintDataSource{}
}

// BlueprintDataSource defines the data source implementation.
type BlueprintDataSource struct {
	client *Client
}

// BlueprintDataSourceModel describes the data source data model.
type BlueprintDataSourceModel struct {
	ID types.String `tfsdk:"id"`

	Description types.String `tfsdk:"description"`
	Meta types.Map `tfsdk:"meta"`
	Name types.String `tfsdk:"name"`
	Visibility types.String `tfsdk:"visibility"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

// Metadata returns the data source type name.
func (d *BlueprintDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_blueprint"
}

// Schema defines the schema for the data source.
func (d *BlueprintDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about an existing blueprint.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique identifier of the blueprint to look up",
			},

			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the blueprint",
				Computed:            true,
			},
			"meta": schema.MapAttribute{
				MarkdownDescription: "Additional metadata for the blueprint",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the blueprint",
				Computed:            true,
			},
			"visibility": schema.StringAttribute{
				MarkdownDescription: "The visibility level of the blueprint",
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
func (d *BlueprintDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *BlueprintDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data BlueprintDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to read blueprint
	result, err := d.client.GetBlueprint(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read blueprint: %s", err))
		return
	}

	// Update data model with response values

	if result.Description != nil {
		data.Description = types.StringPointerValue(result.Description)
	}
	if result.Meta != nil {
		mapValue, diags := types.MapValueFrom(ctx, types.StringType, result.Meta)
		resp.Diagnostics.Append(diags...)
		data.Meta = mapValue
	}
	if result.Name != nil {
		data.Name = types.StringPointerValue(result.Name)
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
