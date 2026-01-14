package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &SkillsetAbilityResource{}
	_ resource.ResourceWithImportState = &SkillsetAbilityResource{}
)

func NewSkillsetAbilityResource() resource.Resource {
	return &SkillsetAbilityResource{}
}

// SkillsetAbilityResource defines the resource implementation.
type SkillsetAbilityResource struct {
	client *Client
}

// SkillsetAbilityResourceModel describes the resource data model.
type SkillsetAbilityResourceModel struct {
	ID types.String `tfsdk:"id"`

	BlueprintId types.String `tfsdk:"blueprint_id"`
	BotId types.String `tfsdk:"bot_id"`
	Description types.String `tfsdk:"description"`
	FileId types.String `tfsdk:"file_id"`
	Instruction types.String `tfsdk:"instruction"`
	Meta types.Map `tfsdk:"meta"`
	Name types.String `tfsdk:"name"`
	SecretId types.String `tfsdk:"secret_id"`
	SpaceId types.String `tfsdk:"space_id"`
}

// Metadata returns the resource type name.
func (r *SkillsetAbilityResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_skillset_ability"
}

// Schema defines the schema for the resource.
func (r *SkillsetAbilityResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Input parameters for creating a new skillset ability",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the skillsetability",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"blueprint_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the blueprint to use",
				Optional:            true,
			},
			"bot_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the bot to use",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the ability",
				Optional:            true,
			},
			"file_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the file to use",
				Optional:            true,
			},
			"instruction": schema.StringAttribute{
				MarkdownDescription: "The instruction for the ability",
				Optional:            true,
			},
			"meta": schema.MapAttribute{
				MarkdownDescription: "Additional metadata for the ability",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the ability",
				Optional:            true,
			},
			"secret_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the secret to use for authentication",
				Optional:            true,
			},
			"space_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the space to use",
				Optional:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *SkillsetAbilityResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *SkillsetAbilityResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SkillsetAbilityResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to create skillsetability
	result, err := r.client.CreateSkillsetAbility(ctx, CreateSkillsetAbilityInput{

		BlueprintId: data.BlueprintId.ValueStringPointer(),
		BotId: data.BotId.ValueStringPointer(),
		Description: data.Description.ValueStringPointer(),
		FileId: data.FileId.ValueStringPointer(),
		Instruction: data.Instruction.ValueStringPointer(),
		// Meta: TODO: convert map type,
		Name: data.Name.ValueStringPointer(),
		SecretId: data.SecretId.ValueStringPointer(),
		SpaceId: data.SpaceId.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create skillsetability: %s", err))
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
func (r *SkillsetAbilityResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SkillsetAbilityResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to read skillsetability
	result, err := r.client.GetSkillsetAbility(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read skillsetability: %s", err))
		return
	}

	// Update data model with response values

	if result.BlueprintId != nil {
		data.BlueprintId = types.StringPointerValue(result.BlueprintId)
	}
	if result.BotId != nil {
		data.BotId = types.StringPointerValue(result.BotId)
	}
	if result.Description != nil {
		data.Description = types.StringPointerValue(result.Description)
	}
	if result.FileId != nil {
		data.FileId = types.StringPointerValue(result.FileId)
	}
	if result.Instruction != nil {
		data.Instruction = types.StringPointerValue(result.Instruction)
	}
	// Meta: TODO: set from response
	if result.Name != nil {
		data.Name = types.StringPointerValue(result.Name)
	}
	if result.SecretId != nil {
		data.SecretId = types.StringPointerValue(result.SecretId)
	}
	if result.SpaceId != nil {
		data.SpaceId = types.StringPointerValue(result.SpaceId)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *SkillsetAbilityResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SkillsetAbilityResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to update skillsetability
	_, err := r.client.UpdateSkillsetAbility(ctx, data.ID.ValueString(), UpdateSkillsetAbilityInput{

		BlueprintId: data.BlueprintId.ValueStringPointer(),
		BotId: data.BotId.ValueStringPointer(),
		Description: data.Description.ValueStringPointer(),
		FileId: data.FileId.ValueStringPointer(),
		Instruction: data.Instruction.ValueStringPointer(),
		// Meta: TODO: convert map type,
		Name: data.Name.ValueStringPointer(),
		SecretId: data.SecretId.ValueStringPointer(),
		SpaceId: data.SpaceId.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update skillsetability: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *SkillsetAbilityResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SkillsetAbilityResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to delete skillsetability
	_, err := r.client.DeleteSkillsetAbility(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete skillsetability: %s", err))
		return
	}
}

// ImportState imports the resource state from Terraform.
func (r *SkillsetAbilityResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
