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
	_ resource.Resource                = &NotionIntegrationResource{}
	_ resource.ResourceWithImportState = &NotionIntegrationResource{}
)

func NewNotionIntegrationResource() resource.Resource {
	return &NotionIntegrationResource{}
}

// NotionIntegrationResource defines the resource implementation.
type NotionIntegrationResource struct {
	client *Client
}

// NotionIntegrationResourceModel describes the resource data model.
type NotionIntegrationResourceModel struct {
	ID types.String `tfsdk:"id"`

	BlueprintId types.String `tfsdk:"blueprint_id"`
	DatasetId types.String `tfsdk:"dataset_id"`
	Description types.String `tfsdk:"description"`
	ExpiresIn types.Int64 `tfsdk:"expires_in"`
	Meta types.Map `tfsdk:"meta"`
	Name types.String `tfsdk:"name"`
	SyncSchedule types.String `tfsdk:"sync_schedule"`
	Token types.String `tfsdk:"token"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

// Metadata returns the resource type name.
func (r *NotionIntegrationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notion_integration"
}

// Schema defines the schema for the resource.
func (r *NotionIntegrationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Input parameters for creating a new Notion integration",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the notionintegration",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"blueprint_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the blueprint to use",
				Optional:            true,
			},
			"dataset_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the dataset to sync to",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the integration",
				Optional:            true,
			},
			"expires_in": schema.Int64Attribute{
				MarkdownDescription: "Time in milliseconds before the data expires",
				Optional:            true,
			},
			"meta": schema.MapAttribute{
				MarkdownDescription: "Additional metadata for the integration",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the integration",
				Optional:            true,
			},
			"sync_schedule": schema.StringAttribute{
				MarkdownDescription: "The schedule for automatic synchronization",
				Optional:            true,
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "The Notion integration token",
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
func (r *NotionIntegrationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *NotionIntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NotionIntegrationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to create notionintegration

	result, err := r.client.CreateNotionIntegration(ctx, CreateNotionIntegrationInput{
		BlueprintId: data.BlueprintId.ValueStringPointer(),
		DatasetId: data.DatasetId.ValueStringPointer(),
		Description: data.Description.ValueStringPointer(),
		ExpiresIn: data.ExpiresIn.ValueInt64Pointer(),
		Meta: convertMapToInterface(ctx, data.Meta),
		Name: data.Name.ValueStringPointer(),
		SyncSchedule: data.SyncSchedule.ValueStringPointer(),
		Token: data.Token.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create notionintegration: %s", err))
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
func (r *NotionIntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NotionIntegrationResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to read notionintegration

	result, err := r.client.GetNotionIntegration(ctx, data.ID.ValueString())
	if err != nil {
		// Check if resource was deleted outside of Terraform
		if strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read notionintegration: %s", err))
		return
	}

	// Update data model with response values

	if result.BlueprintId != nil {
		data.BlueprintId = types.StringPointerValue(result.BlueprintId)
	}
	if result.DatasetId != nil {
		data.DatasetId = types.StringPointerValue(result.DatasetId)
	}
	if result.Description != nil {
		data.Description = types.StringPointerValue(result.Description)
	}
	if result.ExpiresIn != nil {
		data.ExpiresIn = types.Int64PointerValue(result.ExpiresIn)
	}
	if result.Meta != nil {
		mapValue, diags := types.MapValueFrom(ctx, types.StringType, result.Meta)
		resp.Diagnostics.Append(diags...)
		data.Meta = mapValue
	}
	if result.Name != nil {
		data.Name = types.StringPointerValue(result.Name)
	}
	if result.SyncSchedule != nil {
		data.SyncSchedule = types.StringPointerValue(result.SyncSchedule)
	}
	if result.Token != nil {
		data.Token = types.StringPointerValue(result.Token)
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
func (r *NotionIntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data NotionIntegrationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to update notionintegration

	_, err := r.client.UpdateNotionIntegration(ctx, data.ID.ValueString(), UpdateNotionIntegrationInput{
		BlueprintId: data.BlueprintId.ValueStringPointer(),
		DatasetId: data.DatasetId.ValueStringPointer(),
		Description: data.Description.ValueStringPointer(),
		ExpiresIn: data.ExpiresIn.ValueInt64Pointer(),
		Meta: convertMapToInterface(ctx, data.Meta),
		Name: data.Name.ValueStringPointer(),
		SyncSchedule: data.SyncSchedule.ValueStringPointer(),
		Token: data.Token.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update notionintegration: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *NotionIntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NotionIntegrationResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to delete notionintegration

	_, err := r.client.DeleteNotionIntegration(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete notionintegration: %s", err))
		return
	}
}

// ImportState imports the resource state from Terraform.
func (r *NotionIntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
