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
	_ resource.Resource                = &SlackIntegrationResource{}
	_ resource.ResourceWithImportState = &SlackIntegrationResource{}
)

func NewSlackIntegrationResource() resource.Resource {
	return &SlackIntegrationResource{}
}

// SlackIntegrationResource defines the resource implementation.
type SlackIntegrationResource struct {
	client *Client
}

// SlackIntegrationResourceModel describes the resource data model.
type SlackIntegrationResourceModel struct {
	ID types.String `tfsdk:"id"`

	AutoRespond types.String `tfsdk:"auto_respond"`
	BlueprintId types.String `tfsdk:"blueprint_id"`
	BotId types.String `tfsdk:"bot_id"`
	BotToken types.String `tfsdk:"bot_token"`
	ContactCollection types.Bool `tfsdk:"contact_collection"`
	Description types.String `tfsdk:"description"`
	Meta types.Map `tfsdk:"meta"`
	Name types.String `tfsdk:"name"`
	Ratings types.Bool `tfsdk:"ratings"`
	References types.Bool `tfsdk:"references"`
	SessionDuration types.Int64 `tfsdk:"session_duration"`
	SigningSecret types.String `tfsdk:"signing_secret"`
	UserToken types.String `tfsdk:"user_token"`
	VisibleMessages types.Int64 `tfsdk:"visible_messages"`
}

// Metadata returns the resource type name.
func (r *SlackIntegrationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_slack_integration"
}

// Schema defines the schema for the resource.
func (r *SlackIntegrationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Input parameters for creating a new Slack integration",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the slackintegration",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"auto_respond": schema.StringAttribute{
				MarkdownDescription: "Auto-respond configuration for the integration",
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
			"bot_token": schema.StringAttribute{
				MarkdownDescription: "The Slack bot token for API access",
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
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the integration",
				Optional:            true,
			},
			"ratings": schema.BoolAttribute{
				MarkdownDescription: "Whether to enable message ratings",
				Optional:            true,
			},
			"references": schema.BoolAttribute{
				MarkdownDescription: "Whether to include message references",
				Optional:            true,
			},
			"session_duration": schema.Int64Attribute{
				MarkdownDescription: "The duration of the session in milliseconds",
				Optional:            true,
			},
			"signing_secret": schema.StringAttribute{
				MarkdownDescription: "The Slack signing secret for request verification",
				Optional:            true,
			},
			"user_token": schema.StringAttribute{
				MarkdownDescription: "The Slack user token for additional permissions",
				Optional:            true,
			},
			"visible_messages": schema.Int64Attribute{
				MarkdownDescription: "The number of visible messages in the conversation",
				Optional:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *SlackIntegrationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *SlackIntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SlackIntegrationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to create slackintegration
	result, err := r.client.CreateSlackIntegration(ctx, CreateSlackIntegrationInput{

		AutoRespond: data.AutoRespond.ValueStringPointer(),
		BlueprintId: data.BlueprintId.ValueStringPointer(),
		BotId: data.BotId.ValueStringPointer(),
		BotToken: data.BotToken.ValueStringPointer(),
		ContactCollection: data.ContactCollection.ValueBoolPointer(),
		Description: data.Description.ValueStringPointer(),
		// Meta: TODO: convert map type,
		Name: data.Name.ValueStringPointer(),
		Ratings: data.Ratings.ValueBoolPointer(),
		References: data.References.ValueBoolPointer(),
		SessionDuration: data.SessionDuration.ValueInt64Pointer(),
		SigningSecret: data.SigningSecret.ValueStringPointer(),
		UserToken: data.UserToken.ValueStringPointer(),
		VisibleMessages: data.VisibleMessages.ValueInt64Pointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create slackintegration: %s", err))
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
func (r *SlackIntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SlackIntegrationResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to read slackintegration
	result, err := r.client.GetSlackIntegration(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read slackintegration: %s", err))
		return
	}

	// Update data model with response values

	if result.AutoRespond != nil {
		data.AutoRespond = types.StringPointerValue(result.AutoRespond)
	}
	if result.BlueprintId != nil {
		data.BlueprintId = types.StringPointerValue(result.BlueprintId)
	}
	if result.BotId != nil {
		data.BotId = types.StringPointerValue(result.BotId)
	}
	if result.BotToken != nil {
		data.BotToken = types.StringPointerValue(result.BotToken)
	}
	if result.ContactCollection != nil {
		data.ContactCollection = types.BoolPointerValue(result.ContactCollection)
	}
	if result.Description != nil {
		data.Description = types.StringPointerValue(result.Description)
	}
	// Meta: TODO: set from response
	if result.Name != nil {
		data.Name = types.StringPointerValue(result.Name)
	}
	if result.Ratings != nil {
		data.Ratings = types.BoolPointerValue(result.Ratings)
	}
	if result.References != nil {
		data.References = types.BoolPointerValue(result.References)
	}
	if result.SessionDuration != nil {
		data.SessionDuration = types.Int64PointerValue(result.SessionDuration)
	}
	if result.SigningSecret != nil {
		data.SigningSecret = types.StringPointerValue(result.SigningSecret)
	}
	if result.UserToken != nil {
		data.UserToken = types.StringPointerValue(result.UserToken)
	}
	if result.VisibleMessages != nil {
		data.VisibleMessages = types.Int64PointerValue(result.VisibleMessages)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *SlackIntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SlackIntegrationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to update slackintegration
	_, err := r.client.UpdateSlackIntegration(ctx, data.ID.ValueString(), UpdateSlackIntegrationInput{

		AutoRespond: data.AutoRespond.ValueStringPointer(),
		BlueprintId: data.BlueprintId.ValueStringPointer(),
		BotId: data.BotId.ValueStringPointer(),
		BotToken: data.BotToken.ValueStringPointer(),
		ContactCollection: data.ContactCollection.ValueBoolPointer(),
		Description: data.Description.ValueStringPointer(),
		// Meta: TODO: convert map type,
		Name: data.Name.ValueStringPointer(),
		Ratings: data.Ratings.ValueBoolPointer(),
		References: data.References.ValueBoolPointer(),
		SessionDuration: data.SessionDuration.ValueInt64Pointer(),
		SigningSecret: data.SigningSecret.ValueStringPointer(),
		UserToken: data.UserToken.ValueStringPointer(),
		VisibleMessages: data.VisibleMessages.ValueInt64Pointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update slackintegration: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *SlackIntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SlackIntegrationResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to delete slackintegration
	_, err := r.client.DeleteSlackIntegration(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete slackintegration: %s", err))
		return
	}
}

// ImportState imports the resource state from Terraform.
func (r *SlackIntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
