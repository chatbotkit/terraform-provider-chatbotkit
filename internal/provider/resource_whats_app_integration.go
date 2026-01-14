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
	_ resource.Resource                = &WhatsAppIntegrationResource{}
	_ resource.ResourceWithImportState = &WhatsAppIntegrationResource{}
)

func NewWhatsAppIntegrationResource() resource.Resource {
	return &WhatsAppIntegrationResource{}
}

// WhatsAppIntegrationResource defines the resource implementation.
type WhatsAppIntegrationResource struct {
	client *Client
}

// WhatsAppIntegrationResourceModel describes the resource data model.
type WhatsAppIntegrationResourceModel struct {
	ID types.String `tfsdk:"id"`

	AccessToken types.String `tfsdk:"access_token"`
	Attachments types.Bool `tfsdk:"attachments"`
	BlueprintId types.String `tfsdk:"blueprint_id"`
	BotId types.String `tfsdk:"bot_id"`
	ContactCollection types.Bool `tfsdk:"contact_collection"`
	Description types.String `tfsdk:"description"`
	Meta types.Map `tfsdk:"meta"`
	Name types.String `tfsdk:"name"`
	PhoneNumberId types.String `tfsdk:"phone_number_id"`
	SessionDuration types.Int64 `tfsdk:"session_duration"`
}

// Metadata returns the resource type name.
func (r *WhatsAppIntegrationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_whats_app_integration"
}

// Schema defines the schema for the resource.
func (r *WhatsAppIntegrationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Input parameters for creating a new WhatsApp integration",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the whatsappintegration",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"access_token": schema.StringAttribute{
				MarkdownDescription: "The WhatsApp Business API access token",
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
			"phone_number_id": schema.StringAttribute{
				MarkdownDescription: "The WhatsApp Business phone number ID",
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
func (r *WhatsAppIntegrationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *WhatsAppIntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data WhatsAppIntegrationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Implement API call to create whatsappintegration
	// Example:
	// result, err := r.client.CreateWhatsAppIntegration(ctx, types.WhatsAppIntegrationCreateRequest{

	//     AccessToken: data.AccessToken.ValueStringPointer(),
	//     Attachments: data.Attachments.ValueBoolPointer(),
	//     BlueprintId: data.BlueprintId.ValueStringPointer(),
	//     BotId: data.BotId.ValueStringPointer(),
	//     ContactCollection: data.ContactCollection.ValueBoolPointer(),
	//     Description: data.Description.ValueStringPointer(),
	//     Meta: data.Meta.Elements(),
	//     Name: data.Name.ValueStringPointer(),
	//     PhoneNumberId: data.PhoneNumberId.ValueStringPointer(),
	//     SessionDuration: data.SessionDuration.ValueInt64Pointer(),
	// })
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create whatsappintegration: %s", err))
	//     return
	// }
	// data.ID = types.StringValue(result.ID)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *WhatsAppIntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data WhatsAppIntegrationResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Implement API call to read whatsappintegration
	// Example:
	// result, err := r.client.GetWhatsAppIntegration(ctx, data.ID.ValueString())
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read whatsappintegration: %s", err))
	//     return
	// }
	// Update data model with response values

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *WhatsAppIntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data WhatsAppIntegrationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Implement API call to update whatsappintegration
	// Example:
	// _, err := r.client.UpdateWhatsAppIntegration(ctx, data.ID.ValueString(), types.WhatsAppIntegrationUpdateRequest{

	//     AccessToken: data.AccessToken.ValueStringPointer(),
	//     Attachments: data.Attachments.ValueBoolPointer(),
	//     BlueprintId: data.BlueprintId.ValueStringPointer(),
	//     BotId: data.BotId.ValueStringPointer(),
	//     ContactCollection: data.ContactCollection.ValueBoolPointer(),
	//     Description: data.Description.ValueStringPointer(),
	//     Meta: data.Meta.Elements(),
	//     Name: data.Name.ValueStringPointer(),
	//     PhoneNumberId: data.PhoneNumberId.ValueStringPointer(),
	//     SessionDuration: data.SessionDuration.ValueInt64Pointer(),
	// })
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update whatsappintegration: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *WhatsAppIntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data WhatsAppIntegrationResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Implement API call to delete whatsappintegration
	// Example:
	// _, err := r.client.DeleteWhatsAppIntegration(ctx, data.ID.ValueString())
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete whatsappintegration: %s", err))
	//     return
	// }
}

// ImportState imports the resource state from Terraform.
func (r *WhatsAppIntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
