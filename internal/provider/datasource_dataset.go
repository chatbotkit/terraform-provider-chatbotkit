package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &DatasetDataSource{}

func NewDatasetDataSource() datasource.DataSource {
	return &DatasetDataSource{}
}

// DatasetDataSource defines the data source implementation.
type DatasetDataSource struct {
	client *Client
}

// DatasetDataSourceModel describes the data source data model.
type DatasetDataSourceModel struct {
	ID types.String `tfsdk:"id"`

	BlueprintId types.String `tfsdk:"blueprint_id"`
	Description types.String `tfsdk:"description"`
	MatchInstruction types.String `tfsdk:"match_instruction"`
	Meta types.Map `tfsdk:"meta"`
	MismatchInstruction types.String `tfsdk:"mismatch_instruction"`
	Name types.String `tfsdk:"name"`
	RecordMaxTokens types.Int64 `tfsdk:"record_max_tokens"`
	Reranker types.String `tfsdk:"reranker"`
	SearchMaxRecords types.Int64 `tfsdk:"search_max_records"`
	SearchMaxTokens types.Int64 `tfsdk:"search_max_tokens"`
	SearchMinScore types.Float64 `tfsdk:"search_min_score"`
	Separators types.String `tfsdk:"separators"`
	Store types.String `tfsdk:"store"`
	Visibility types.String `tfsdk:"visibility"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

// Metadata returns the data source type name.
func (d *DatasetDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dataset"
}

// Schema defines the schema for the data source.
func (d *DatasetDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about an existing dataset.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique identifier of the dataset to look up",
			},

			"blueprint_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the blueprint to use",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the dataset",
				Computed:            true,
			},
			"match_instruction": schema.StringAttribute{
				MarkdownDescription: "Instruction when matches are found",
				Computed:            true,
			},
			"meta": schema.MapAttribute{
				MarkdownDescription: "Additional metadata for the dataset",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"mismatch_instruction": schema.StringAttribute{
				MarkdownDescription: "Instruction when no matches are found",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the dataset",
				Computed:            true,
			},
			"record_max_tokens": schema.Int64Attribute{
				MarkdownDescription: "Maximum tokens per record",
				Computed:            true,
			},
			"reranker": schema.StringAttribute{
				MarkdownDescription: "The reranking model to use",
				Computed:            true,
			},
			"search_max_records": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of search results",
				Computed:            true,
			},
			"search_max_tokens": schema.Int64Attribute{
				MarkdownDescription: "Maximum tokens in search results",
				Computed:            true,
			},
			"search_min_score": schema.Float64Attribute{
				MarkdownDescription: "Minimum score for search results",
				Computed:            true,
			},
			"separators": schema.StringAttribute{
				MarkdownDescription: "The separators for chunking text",
				Computed:            true,
			},
			"store": schema.StringAttribute{
				MarkdownDescription: "The storage backend to use",
				Computed:            true,
			},
			"visibility": schema.StringAttribute{
				MarkdownDescription: "The visibility level of the dataset",
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
func (d *DatasetDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *DatasetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DatasetDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to read dataset
	result, err := d.client.GetDataset(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read dataset: %s", err))
		return
	}

	// Update data model with response values

	if result.BlueprintId != nil {
		data.BlueprintId = types.StringPointerValue(result.BlueprintId)
	}
	if result.Description != nil {
		data.Description = types.StringPointerValue(result.Description)
	}
	if result.MatchInstruction != nil {
		data.MatchInstruction = types.StringPointerValue(result.MatchInstruction)
	}
	if result.Meta != nil {
		mapValue, diags := types.MapValueFrom(ctx, types.StringType, result.Meta)
		resp.Diagnostics.Append(diags...)
		data.Meta = mapValue
	}
	if result.MismatchInstruction != nil {
		data.MismatchInstruction = types.StringPointerValue(result.MismatchInstruction)
	}
	if result.Name != nil {
		data.Name = types.StringPointerValue(result.Name)
	}
	if result.RecordMaxTokens != nil {
		data.RecordMaxTokens = types.Int64PointerValue(result.RecordMaxTokens)
	}
	if result.Reranker != nil {
		data.Reranker = types.StringPointerValue(result.Reranker)
	}
	if result.SearchMaxRecords != nil {
		data.SearchMaxRecords = types.Int64PointerValue(result.SearchMaxRecords)
	}
	if result.SearchMaxTokens != nil {
		data.SearchMaxTokens = types.Int64PointerValue(result.SearchMaxTokens)
	}
	if result.SearchMinScore != nil {
		data.SearchMinScore = types.Float64PointerValue(result.SearchMinScore)
	}
	if result.Separators != nil {
		data.Separators = types.StringPointerValue(result.Separators)
	}
	if result.Store != nil {
		data.Store = types.StringPointerValue(result.Store)
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
