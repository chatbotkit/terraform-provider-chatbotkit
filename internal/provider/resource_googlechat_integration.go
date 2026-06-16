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
	_ resource.Resource                = &GooglechatIntegrationResource{}
	_ resource.ResourceWithImportState = &GooglechatIntegrationResource{}
)

func NewGooglechatIntegrationResource() resource.Resource {
	return &GooglechatIntegrationResource{}
}

// GooglechatIntegrationResource defines the resource implementation.
type GooglechatIntegrationResource struct {
	client *Client
}

// GooglechatIntegrationResourceModel describes the resource data model.
type GooglechatIntegrationResourceModel struct {
	ID types.String `tfsdk:"id"`

	Alias             types.String `tfsdk:"alias"`
	AllowFrom         types.String `tfsdk:"allow_from"`
	Attachments       types.Bool   `tfsdk:"attachments"`
	AutoRespond       types.String `tfsdk:"auto_respond"`
	BlueprintId       types.String `tfsdk:"blueprint_id"`
	BotId             types.String `tfsdk:"bot_id"`
	ContactCollection types.Bool   `tfsdk:"contact_collection"`
	Description       types.String `tfsdk:"description"`
	Meta              types.Map    `tfsdk:"meta"`
	Name              types.String `tfsdk:"name"`
	ProjectNumber     types.String `tfsdk:"project_number"`
	ServiceAccountKey types.String `tfsdk:"service_account_key"`
	SessionDuration   types.Int64  `tfsdk:"session_duration"`
	CreatedAt         types.String `tfsdk:"created_at"`
	UpdatedAt         types.String `tfsdk:"updated_at"`
}

// Metadata returns the resource type name.
func (r *GooglechatIntegrationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_googlechat_integration"
}

// Schema defines the schema for the resource.
func (r *GooglechatIntegrationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Input parameters for creating a new Google Chat integration",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the googlechatintegration",
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
			"attachments": schema.BoolAttribute{
				MarkdownDescription: "Whether to enable file attachments",
				Optional:            true,
			},
			"auto_respond": schema.StringAttribute{
				MarkdownDescription: "Configure automatic response behavior",
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
			"project_number": schema.StringAttribute{
				MarkdownDescription: "The Google Cloud project number used to verify incoming event JWT audience claims",
				Optional:            true,
			},
			"service_account_key": schema.StringAttribute{
				MarkdownDescription: "The Google service account JSON key for sending messages via the Chat REST API",
				Optional:            true,
				Sensitive:           true,
			},
			"session_duration": schema.Int64Attribute{
				MarkdownDescription: "The duration of the session in milliseconds",
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
func (r *GooglechatIntegrationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *GooglechatIntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data GooglechatIntegrationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to create googlechatintegration

	result, err := r.client.CreateGooglechatIntegration(ctx, CreateGooglechatIntegrationInput{
		Alias:             data.Alias.ValueStringPointer(),
		AllowFrom:         data.AllowFrom.ValueStringPointer(),
		Attachments:       data.Attachments.ValueBoolPointer(),
		AutoRespond:       data.AutoRespond.ValueStringPointer(),
		BlueprintId:       data.BlueprintId.ValueStringPointer(),
		BotId:             data.BotId.ValueStringPointer(),
		ContactCollection: data.ContactCollection.ValueBoolPointer(),
		Description:       data.Description.ValueStringPointer(),
		Meta:              convertMapToInterface(ctx, data.Meta),
		Name:              data.Name.ValueStringPointer(),
		ProjectNumber:     data.ProjectNumber.ValueStringPointer(),
		ServiceAccountKey: data.ServiceAccountKey.ValueStringPointer(),
		SessionDuration:   data.SessionDuration.ValueInt64Pointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create googlechatintegration: %s", err))
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
func (r *GooglechatIntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data GooglechatIntegrationResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to read googlechatintegration

	result, err := r.client.GetGooglechatIntegration(ctx, data.ID.ValueString())
	if err != nil {
		// Check if resource was deleted outside of Terraform
		if strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read googlechatintegration: %s", err))
		return
	}

	// Update data model with response values

	if result.Alias != nil {
		data.Alias = types.StringPointerValue(result.Alias)
	}
	if result.AllowFrom != nil {
		data.AllowFrom = types.StringPointerValue(result.AllowFrom)
	}
	if result.Attachments != nil {
		data.Attachments = types.BoolPointerValue(result.Attachments)
	}
	if result.AutoRespond != nil {
		data.AutoRespond = types.StringPointerValue(result.AutoRespond)
	}
	if result.BlueprintId != nil {
		data.BlueprintId = types.StringPointerValue(result.BlueprintId)
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
	if result.ProjectNumber != nil {
		data.ProjectNumber = types.StringPointerValue(result.ProjectNumber)
	}
	if result.ServiceAccountKey != nil {
		data.ServiceAccountKey = types.StringPointerValue(result.ServiceAccountKey)
	}
	if result.SessionDuration != nil {
		data.SessionDuration = types.Int64PointerValue(result.SessionDuration)
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
func (r *GooglechatIntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data GooglechatIntegrationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to update googlechatintegration

	_, err := r.client.UpdateGooglechatIntegration(ctx, data.ID.ValueString(), UpdateGooglechatIntegrationInput{
		Alias:             data.Alias.ValueStringPointer(),
		AllowFrom:         data.AllowFrom.ValueStringPointer(),
		Attachments:       data.Attachments.ValueBoolPointer(),
		AutoRespond:       data.AutoRespond.ValueStringPointer(),
		BlueprintId:       data.BlueprintId.ValueStringPointer(),
		BotId:             data.BotId.ValueStringPointer(),
		ContactCollection: data.ContactCollection.ValueBoolPointer(),
		Description:       data.Description.ValueStringPointer(),
		Meta:              convertMapToInterface(ctx, data.Meta),
		Name:              data.Name.ValueStringPointer(),
		ProjectNumber:     data.ProjectNumber.ValueStringPointer(),
		ServiceAccountKey: data.ServiceAccountKey.ValueStringPointer(),
		SessionDuration:   data.SessionDuration.ValueInt64Pointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update googlechatintegration: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *GooglechatIntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data GooglechatIntegrationResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to delete googlechatintegration

	_, err := r.client.DeleteGooglechatIntegration(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete googlechatintegration: %s", err))
		return
	}
}

// ImportState imports the resource state from Terraform.
func (r *GooglechatIntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
