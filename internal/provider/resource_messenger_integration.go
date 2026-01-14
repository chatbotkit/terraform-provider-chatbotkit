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
	_ resource.Resource                = &MessengerIntegrationResource{}
	_ resource.ResourceWithImportState = &MessengerIntegrationResource{}
)

func NewMessengerIntegrationResource() resource.Resource {
	return &MessengerIntegrationResource{}
}

// MessengerIntegrationResource defines the resource implementation.
type MessengerIntegrationResource struct {
	client *Client
}

// MessengerIntegrationResourceModel describes the resource data model.
type MessengerIntegrationResourceModel struct {
	ID types.String `tfsdk:"id"`

	AccessToken types.String `tfsdk:"access_token"`
	Attachments types.Bool `tfsdk:"attachments"`
	BlueprintId types.String `tfsdk:"blueprint_id"`
	BotId types.String `tfsdk:"bot_id"`
	Description types.String `tfsdk:"description"`
	Meta types.Map `tfsdk:"meta"`
	Name types.String `tfsdk:"name"`
	SessionDuration types.Int64 `tfsdk:"session_duration"`
}

// Metadata returns the resource type name.
func (r *MessengerIntegrationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_messenger_integration"
}

// Schema defines the schema for the resource.
func (r *MessengerIntegrationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Input parameters for creating a new Messenger integration",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the messengerintegration",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"access_token": schema.StringAttribute{
				MarkdownDescription: "The Facebook Messenger page access token",
				Optional:            true,
			},
			"attachments": schema.BoolAttribute{
				MarkdownDescription: "Whether to enable file attachments",
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
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the integration",
				Optional:            true,
			},
			"session_duration": schema.Int64Attribute{
				MarkdownDescription: "The duration of the session in milliseconds",
				Optional:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *MessengerIntegrationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *MessengerIntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data MessengerIntegrationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Implement API call to create messengerintegration
	// Example:
	// result, err := r.client.CreateMessengerIntegration(ctx, types.MessengerIntegrationCreateRequest{

	//     AccessToken: data.AccessToken.ValueStringPointer(),
	//     Attachments: data.Attachments.ValueBoolPointer(),
	//     BlueprintId: data.BlueprintId.ValueStringPointer(),
	//     BotId: data.BotId.ValueStringPointer(),
	//     Description: data.Description.ValueStringPointer(),
	//     Meta: data.Meta.Elements(),
	//     Name: data.Name.ValueStringPointer(),
	//     SessionDuration: data.SessionDuration.ValueInt64Pointer(),
	// })
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create messengerintegration: %s", err))
	//     return
	// }
	// data.ID = types.StringValue(result.ID)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *MessengerIntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data MessengerIntegrationResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Implement API call to read messengerintegration
	// Example:
	// result, err := r.client.GetMessengerIntegration(ctx, data.ID.ValueString())
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read messengerintegration: %s", err))
	//     return
	// }
	// Update data model with response values

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *MessengerIntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data MessengerIntegrationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Implement API call to update messengerintegration
	// Example:
	// _, err := r.client.UpdateMessengerIntegration(ctx, data.ID.ValueString(), types.MessengerIntegrationUpdateRequest{

	//     AccessToken: data.AccessToken.ValueStringPointer(),
	//     Attachments: data.Attachments.ValueBoolPointer(),
	//     BlueprintId: data.BlueprintId.ValueStringPointer(),
	//     BotId: data.BotId.ValueStringPointer(),
	//     Description: data.Description.ValueStringPointer(),
	//     Meta: data.Meta.Elements(),
	//     Name: data.Name.ValueStringPointer(),
	//     SessionDuration: data.SessionDuration.ValueInt64Pointer(),
	// })
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update messengerintegration: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *MessengerIntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data MessengerIntegrationResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Implement API call to delete messengerintegration
	// Example:
	// _, err := r.client.DeleteMessengerIntegration(ctx, data.ID.ValueString())
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete messengerintegration: %s", err))
	//     return
	// }
}

// ImportState imports the resource state from Terraform.
func (r *MessengerIntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
