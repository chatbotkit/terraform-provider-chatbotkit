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
	_ resource.Resource                = &TriggerIntegrationResource{}
	_ resource.ResourceWithImportState = &TriggerIntegrationResource{}
)

func NewTriggerIntegrationResource() resource.Resource {
	return &TriggerIntegrationResource{}
}

// TriggerIntegrationResource defines the resource implementation.
type TriggerIntegrationResource struct {
	client *Client
}

// TriggerIntegrationResourceModel describes the resource data model.
type TriggerIntegrationResourceModel struct {
	ID types.String `tfsdk:"id"`

	Authenticate types.Bool `tfsdk:"authenticate"`
	BlueprintId types.String `tfsdk:"blueprint_id"`
	BotId types.String `tfsdk:"bot_id"`
	Description types.String `tfsdk:"description"`
	Meta types.Map `tfsdk:"meta"`
	Name types.String `tfsdk:"name"`
	SessionDuration types.Int64 `tfsdk:"session_duration"`
	TriggerSchedule types.String `tfsdk:"trigger_schedule"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

// Metadata returns the resource type name.
func (r *TriggerIntegrationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_trigger_integration"
}

// Schema defines the schema for the resource.
func (r *TriggerIntegrationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Input parameters for creating a new Trigger integration",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the triggerintegration",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"authenticate": schema.BoolAttribute{
				MarkdownDescription: "Whether to require authentication for the trigger",
				Optional:            true,
			},
			"blueprint_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the blueprint to use",
				Optional:            true,
			},
			"bot_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the bot to connect",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the integration",
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
			"session_duration": schema.Int64Attribute{
				MarkdownDescription: "The duration of the session in milliseconds",
				Optional:            true,
			},
			"trigger_schedule": schema.StringAttribute{
				MarkdownDescription: "The schedule for automatic trigger execution",
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
func (r *TriggerIntegrationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *TriggerIntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TriggerIntegrationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to create triggerintegration
	result, err := r.client.CreateTriggerIntegration(ctx, CreateTriggerIntegrationInput{

		Authenticate: data.Authenticate.ValueBoolPointer(),
		BlueprintId: data.BlueprintId.ValueStringPointer(),
		BotId: data.BotId.ValueStringPointer(),
		Description: data.Description.ValueStringPointer(),
		Meta: convertMapToInterface(ctx, data.Meta),
		Name: data.Name.ValueStringPointer(),
		SessionDuration: data.SessionDuration.ValueInt64Pointer(),
		TriggerSchedule: data.TriggerSchedule.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create triggerintegration: %s", err))
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
func (r *TriggerIntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TriggerIntegrationResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to read triggerintegration
	result, err := r.client.GetTriggerIntegration(ctx, data.ID.ValueString())
	if err != nil {
		// Check if resource was deleted outside of Terraform
		if strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read triggerintegration: %s", err))
		return
	}

	// Update data model with response values

	if result.Authenticate != nil {
		data.Authenticate = types.BoolPointerValue(result.Authenticate)
	}
	if result.BlueprintId != nil {
		data.BlueprintId = types.StringPointerValue(result.BlueprintId)
	}
	if result.BotId != nil {
		data.BotId = types.StringPointerValue(result.BotId)
	}
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
	if result.SessionDuration != nil {
		data.SessionDuration = types.Int64PointerValue(result.SessionDuration)
	}
	if result.TriggerSchedule != nil {
		data.TriggerSchedule = types.StringPointerValue(result.TriggerSchedule)
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
func (r *TriggerIntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TriggerIntegrationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to update triggerintegration
	_, err := r.client.UpdateTriggerIntegration(ctx, data.ID.ValueString(), UpdateTriggerIntegrationInput{

		Authenticate: data.Authenticate.ValueBoolPointer(),
		BlueprintId: data.BlueprintId.ValueStringPointer(),
		BotId: data.BotId.ValueStringPointer(),
		Description: data.Description.ValueStringPointer(),
		Meta: convertMapToInterface(ctx, data.Meta),
		Name: data.Name.ValueStringPointer(),
		SessionDuration: data.SessionDuration.ValueInt64Pointer(),
		TriggerSchedule: data.TriggerSchedule.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update triggerintegration: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *TriggerIntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TriggerIntegrationResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to delete triggerintegration
	_, err := r.client.DeleteTriggerIntegration(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete triggerintegration: %s", err))
		return
	}
}

// ImportState imports the resource state from Terraform.
func (r *TriggerIntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
