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
	_ resource.Resource                = &MicrosoftteamsIntegrationResource{}
	_ resource.ResourceWithImportState = &MicrosoftteamsIntegrationResource{}
)

func NewMicrosoftteamsIntegrationResource() resource.Resource {
	return &MicrosoftteamsIntegrationResource{}
}

// MicrosoftteamsIntegrationResource defines the resource implementation.
type MicrosoftteamsIntegrationResource struct {
	client *Client
}

// MicrosoftteamsIntegrationResourceModel describes the resource data model.
type MicrosoftteamsIntegrationResourceModel struct {
	ID types.String `tfsdk:"id"`

	Alias                 types.String `tfsdk:"alias"`
	AllowFrom             types.String `tfsdk:"allow_from"`
	BlueprintId           types.String `tfsdk:"blueprint_id"`
	BotFrameworkAppId     types.String `tfsdk:"bot_framework_app_id"`
	BotFrameworkAppSecret types.String `tfsdk:"bot_framework_app_secret"`
	BotId                 types.String `tfsdk:"bot_id"`
	ContactCollection     types.Bool   `tfsdk:"contact_collection"`
	Description           types.String `tfsdk:"description"`
	Meta                  types.Map    `tfsdk:"meta"`
	Name                  types.String `tfsdk:"name"`
	SessionDuration       types.Int64  `tfsdk:"session_duration"`
	TenantId              types.String `tfsdk:"tenant_id"`
	CreatedAt             types.String `tfsdk:"created_at"`
	UpdatedAt             types.String `tfsdk:"updated_at"`
}

// Metadata returns the resource type name.
func (r *MicrosoftteamsIntegrationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_microsoftteams_integration"
}

// Schema defines the schema for the resource.
func (r *MicrosoftteamsIntegrationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Input parameters for creating a new Microsoft Teams integration",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the microsoftteamsintegration",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"alias": schema.StringAttribute{
				MarkdownDescription: "The alias ID for the integration",
				Optional:            true,
			},
			"allow_from": schema.StringAttribute{
				MarkdownDescription: "The allowed senders for this integration",
				Optional:            true,
			},
			"blueprint_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the blueprint to use",
				Optional:            true,
			},
			"bot_framework_app_id": schema.StringAttribute{
				MarkdownDescription: "The Microsoft Bot Framework application ID",
				Optional:            true,
			},
			"bot_framework_app_secret": schema.StringAttribute{
				MarkdownDescription: "The Microsoft Bot Framework application secret",
				Optional:            true,
				Sensitive:           true,
			},
			"bot_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the bot to connect",
				Optional:            true,
			},
			"contact_collection": schema.BoolAttribute{
				MarkdownDescription: "Whether to collect contact information",
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
			"tenant_id": schema.StringAttribute{
				MarkdownDescription: "The Azure AD tenant ID",
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
func (r *MicrosoftteamsIntegrationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *MicrosoftteamsIntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data MicrosoftteamsIntegrationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to create microsoftteamsintegration

	result, err := r.client.CreateMicrosoftteamsIntegration(ctx, CreateMicrosoftteamsIntegrationInput{
		Alias:                 data.Alias.ValueStringPointer(),
		AllowFrom:             data.AllowFrom.ValueStringPointer(),
		BlueprintId:           data.BlueprintId.ValueStringPointer(),
		BotFrameworkAppId:     data.BotFrameworkAppId.ValueStringPointer(),
		BotFrameworkAppSecret: data.BotFrameworkAppSecret.ValueStringPointer(),
		BotId:                 data.BotId.ValueStringPointer(),
		ContactCollection:     data.ContactCollection.ValueBoolPointer(),
		Description:           data.Description.ValueStringPointer(),
		Meta:                  convertMapToInterface(ctx, data.Meta),
		Name:                  data.Name.ValueStringPointer(),
		SessionDuration:       data.SessionDuration.ValueInt64Pointer(),
		TenantId:              data.TenantId.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create microsoftteamsintegration: %s", err))
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
func (r *MicrosoftteamsIntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data MicrosoftteamsIntegrationResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to read microsoftteamsintegration

	result, err := r.client.GetMicrosoftteamsIntegration(ctx, data.ID.ValueString())
	if err != nil {
		// Check if resource was deleted outside of Terraform
		if strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read microsoftteamsintegration: %s", err))
		return
	}

	// Update data model with response values

	if result.Alias != nil {
		data.Alias = types.StringPointerValue(result.Alias)
	}
	if result.AllowFrom != nil {
		data.AllowFrom = types.StringPointerValue(result.AllowFrom)
	}
	if result.BlueprintId != nil {
		data.BlueprintId = types.StringPointerValue(result.BlueprintId)
	}
	if result.BotFrameworkAppId != nil {
		data.BotFrameworkAppId = types.StringPointerValue(result.BotFrameworkAppId)
	}
	if result.BotFrameworkAppSecret != nil {
		data.BotFrameworkAppSecret = types.StringPointerValue(result.BotFrameworkAppSecret)
	}
	if result.BotId != nil {
		data.BotId = types.StringPointerValue(result.BotId)
	}
	if result.ContactCollection != nil {
		data.ContactCollection = types.BoolPointerValue(result.ContactCollection)
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
	if result.TenantId != nil {
		data.TenantId = types.StringPointerValue(result.TenantId)
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
func (r *MicrosoftteamsIntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data MicrosoftteamsIntegrationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to update microsoftteamsintegration

	_, err := r.client.UpdateMicrosoftteamsIntegration(ctx, data.ID.ValueString(), UpdateMicrosoftteamsIntegrationInput{
		Alias:                 data.Alias.ValueStringPointer(),
		AllowFrom:             data.AllowFrom.ValueStringPointer(),
		BlueprintId:           data.BlueprintId.ValueStringPointer(),
		BotFrameworkAppId:     data.BotFrameworkAppId.ValueStringPointer(),
		BotFrameworkAppSecret: data.BotFrameworkAppSecret.ValueStringPointer(),
		BotId:                 data.BotId.ValueStringPointer(),
		ContactCollection:     data.ContactCollection.ValueBoolPointer(),
		Description:           data.Description.ValueStringPointer(),
		Meta:                  convertMapToInterface(ctx, data.Meta),
		Name:                  data.Name.ValueStringPointer(),
		SessionDuration:       data.SessionDuration.ValueInt64Pointer(),
		TenantId:              data.TenantId.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update microsoftteamsintegration: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *MicrosoftteamsIntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data MicrosoftteamsIntegrationResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to delete microsoftteamsintegration

	_, err := r.client.DeleteMicrosoftteamsIntegration(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete microsoftteamsintegration: %s", err))
		return
	}
}

// ImportState imports the resource state from Terraform.
func (r *MicrosoftteamsIntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
