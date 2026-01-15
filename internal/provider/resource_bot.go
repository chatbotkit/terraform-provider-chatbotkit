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
	_ resource.Resource                = &BotResource{}
	_ resource.ResourceWithImportState = &BotResource{}
)

func NewBotResource() resource.Resource {
	return &BotResource{}
}

// BotResource defines the resource implementation.
type BotResource struct {
	client *Client
}

// BotResourceModel describes the resource data model.
type BotResourceModel struct {
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

// Metadata returns the resource type name.
func (r *BotResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bot"
}

// Schema defines the schema for the resource.
func (r *BotResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Input parameters for creating a new bot",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the bot",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"backstory": schema.StringAttribute{
				MarkdownDescription: "The backstory for the bot",
				Optional:            true,
			},
			"blueprint_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the blueprint to use",
				Optional:            true,
			},
			"dataset_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the dataset to use",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the bot",
				Optional:            true,
			},
			"meta": schema.MapAttribute{
				MarkdownDescription: "Additional metadata for the bot",
				Optional:            true,
			},
			"model": schema.StringAttribute{
				MarkdownDescription: "The AI model to use for the bot",
				Optional:            true,
			},
			"moderation": schema.BoolAttribute{
				MarkdownDescription: "Whether moderation is enabled",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the bot",
				Optional:            true,
			},
			"privacy": schema.BoolAttribute{
				MarkdownDescription: "Whether privacy mode is enabled",
				Optional:            true,
			},
			"skillset_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the skillset to use",
				Optional:            true,
			},
			"visibility": schema.StringAttribute{
				MarkdownDescription: "The visibility level of the bot",
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
func (r *BotResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *BotResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data BotResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to create bot
	result, err := r.client.CreateBot(ctx, CreateBotInput{

		Backstory: data.Backstory.ValueStringPointer(),
		BlueprintId: data.BlueprintId.ValueStringPointer(),
		DatasetId: data.DatasetId.ValueStringPointer(),
		Description: data.Description.ValueStringPointer(),
		Meta: convertMapToInterface(ctx, data.Meta),
		Model: data.Model.ValueStringPointer(),
		Moderation: data.Moderation.ValueBoolPointer(),
		Name: data.Name.ValueStringPointer(),
		Privacy: data.Privacy.ValueBoolPointer(),
		SkillsetId: data.SkillsetId.ValueStringPointer(),
		Visibility: data.Visibility.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create bot: %s", err))
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
func (r *BotResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data BotResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to read bot
	result, err := r.client.GetBot(ctx, data.ID.ValueString())
	if err != nil {
		// Check if resource was deleted outside of Terraform
		if strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
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

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *BotResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data BotResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to update bot
	_, err := r.client.UpdateBot(ctx, data.ID.ValueString(), UpdateBotInput{

		Backstory: data.Backstory.ValueStringPointer(),
		BlueprintId: data.BlueprintId.ValueStringPointer(),
		DatasetId: data.DatasetId.ValueStringPointer(),
		Description: data.Description.ValueStringPointer(),
		Meta: convertMapToInterface(ctx, data.Meta),
		Model: data.Model.ValueStringPointer(),
		Moderation: data.Moderation.ValueBoolPointer(),
		Name: data.Name.ValueStringPointer(),
		Privacy: data.Privacy.ValueBoolPointer(),
		SkillsetId: data.SkillsetId.ValueStringPointer(),
		Visibility: data.Visibility.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update bot: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *BotResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data BotResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to delete bot
	_, err := r.client.DeleteBot(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete bot: %s", err))
		return
	}
}

// ImportState imports the resource state from Terraform.
func (r *BotResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
