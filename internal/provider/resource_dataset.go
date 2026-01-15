package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &DatasetResource{}
	_ resource.ResourceWithImportState = &DatasetResource{}
)

func NewDatasetResource() resource.Resource {
	return &DatasetResource{}
}

// DatasetResource defines the resource implementation.
type DatasetResource struct {
	client *Client
}

// DatasetResourceModel describes the resource data model.
type DatasetResourceModel struct {
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

// Metadata returns the resource type name.
func (r *DatasetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dataset"
}

// Schema defines the schema for the resource.
func (r *DatasetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Input parameters for creating a new dataset",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the dataset",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"blueprint_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the blueprint to use",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the dataset",
				Optional:            true,
			},
			"match_instruction": schema.StringAttribute{
				MarkdownDescription: "Instruction when matches are found",
				Optional:            true,
			},
			"meta": schema.MapAttribute{
				MarkdownDescription: "Additional metadata for the dataset",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"mismatch_instruction": schema.StringAttribute{
				MarkdownDescription: "Instruction when no matches are found",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the dataset",
				Optional:            true,
			},
			"record_max_tokens": schema.Int64Attribute{
				MarkdownDescription: "Maximum tokens per record",
				Optional:            true,
			},
			"reranker": schema.StringAttribute{
				MarkdownDescription: "The reranking model to use",
				Optional:            true,
			},
			"search_max_records": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of search results",
				Optional:            true,
			},
			"search_max_tokens": schema.Int64Attribute{
				MarkdownDescription: "Maximum tokens in search results",
				Optional:            true,
			},
			"search_min_score": schema.Float64Attribute{
				MarkdownDescription: "Minimum score for search results",
				Optional:            true,
			},
			"separators": schema.StringAttribute{
				MarkdownDescription: "The separators for chunking text",
				Optional:            true,
			},
			"store": schema.StringAttribute{
				MarkdownDescription: "The storage backend to use",
				Optional:            true,
			},
			"visibility": schema.StringAttribute{
				MarkdownDescription: "The visibility level of the dataset",
				Optional:            true,
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

// Configure adds the provider configured client to the resource.
func (r *DatasetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *DatasetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DatasetResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to create dataset

	result, err := r.client.CreateDataset(ctx, CreateDatasetInput{
		BlueprintId: data.BlueprintId.ValueStringPointer(),
		Description: data.Description.ValueStringPointer(),
		MatchInstruction: data.MatchInstruction.ValueStringPointer(),
		Meta: convertMapToInterface(ctx, data.Meta),
		MismatchInstruction: data.MismatchInstruction.ValueStringPointer(),
		Name: data.Name.ValueStringPointer(),
		RecordMaxTokens: data.RecordMaxTokens.ValueInt64Pointer(),
		Reranker: data.Reranker.ValueStringPointer(),
		SearchMaxRecords: data.SearchMaxRecords.ValueInt64Pointer(),
		SearchMaxTokens: data.SearchMaxTokens.ValueInt64Pointer(),
		SearchMinScore: data.SearchMinScore.ValueFloat64Pointer(),
		Separators: data.Separators.ValueStringPointer(),
		Store: data.Store.ValueStringPointer(),
		Visibility: data.Visibility.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create dataset: %s", err))
		return
	}

	// Set the ID from the response
	if result.ID != nil {
		data.ID = types.StringPointerValue(result.ID)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *DatasetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DatasetResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to read dataset

	result, err := r.client.GetDataset(ctx, data.ID.ValueString())
	if err != nil {
		// Check if resource was deleted outside of Terraform
		if strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
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

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *DatasetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DatasetResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to update dataset

	_, err := r.client.UpdateDataset(ctx, data.ID.ValueString(), UpdateDatasetInput{
		BlueprintId: data.BlueprintId.ValueStringPointer(),
		Description: data.Description.ValueStringPointer(),
		MatchInstruction: data.MatchInstruction.ValueStringPointer(),
		Meta: convertMapToInterface(ctx, data.Meta),
		MismatchInstruction: data.MismatchInstruction.ValueStringPointer(),
		Name: data.Name.ValueStringPointer(),
		RecordMaxTokens: data.RecordMaxTokens.ValueInt64Pointer(),
		Reranker: data.Reranker.ValueStringPointer(),
		SearchMaxRecords: data.SearchMaxRecords.ValueInt64Pointer(),
		SearchMaxTokens: data.SearchMaxTokens.ValueInt64Pointer(),
		SearchMinScore: data.SearchMinScore.ValueFloat64Pointer(),
		Separators: data.Separators.ValueStringPointer(),
		Store: data.Store.ValueStringPointer(),
		Visibility: data.Visibility.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update dataset: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *DatasetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DatasetResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to delete dataset

	_, err := r.client.DeleteDataset(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete dataset: %s", err))
		return
	}
}

// ImportState imports the resource state from Terraform.
func (r *DatasetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
