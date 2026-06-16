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
	_ resource.Resource                = &WidgetIntegrationResource{}
	_ resource.ResourceWithImportState = &WidgetIntegrationResource{}
)

func NewWidgetIntegrationResource() resource.Resource {
	return &WidgetIntegrationResource{}
}

// WidgetIntegrationResource defines the resource implementation.
type WidgetIntegrationResource struct {
	client *Client
}

// WidgetIntegrationResourceModel describes the resource data model.
type WidgetIntegrationResourceModel struct {
	ID types.String `tfsdk:"id"`

	Alias               types.String `tfsdk:"alias"`
	Attachments         types.Bool   `tfsdk:"attachments"`
	AutoScroll          types.Bool   `tfsdk:"auto_scroll"`
	BlueprintId         types.String `tfsdk:"blueprint_id"`
	BotId               types.String `tfsdk:"bot_id"`
	Carousel            types.Bool   `tfsdk:"carousel"`
	ContactCollection   types.Bool   `tfsdk:"contact_collection"`
	Description         types.String `tfsdk:"description"`
	ExportConversation  types.Bool   `tfsdk:"export_conversation"`
	Form                types.Bool   `tfsdk:"form"`
	Initial             types.String `tfsdk:"initial"`
	Intro               types.String `tfsdk:"intro"`
	Language            types.String `tfsdk:"language"`
	Layout              types.String `tfsdk:"layout"`
	Math                types.Bool   `tfsdk:"math"`
	Maximize            types.Bool   `tfsdk:"maximize"`
	MessagePeek         types.Bool   `tfsdk:"message_peek"`
	Meta                types.Map    `tfsdk:"meta"`
	Name                types.String `tfsdk:"name"`
	Origin              types.String `tfsdk:"origin"`
	Placeholder         types.String `tfsdk:"placeholder"`
	Plugins             types.String `tfsdk:"plugins"`
	PoweredBy           types.Bool   `tfsdk:"powered_by"`
	RestartConversation types.Bool   `tfsdk:"restart_conversation"`
	SessionDuration     types.Int64  `tfsdk:"session_duration"`
	StartFirst          types.Bool   `tfsdk:"start_first"`
	Stream              types.Bool   `tfsdk:"stream"`
	Theme               types.String `tfsdk:"theme"`
	Title               types.String `tfsdk:"title"`
	Tools               types.Bool   `tfsdk:"tools"`
	Unfurl              types.Bool   `tfsdk:"unfurl"`
	Verbose             types.Bool   `tfsdk:"verbose"`
	VoiceIn             types.Bool   `tfsdk:"voice_in"`
	VoiceOut            types.Bool   `tfsdk:"voice_out"`
	CreatedAt           types.String `tfsdk:"created_at"`
	UpdatedAt           types.String `tfsdk:"updated_at"`
}

// Metadata returns the resource type name.
func (r *WidgetIntegrationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_widget_integration"
}

// Schema defines the schema for the resource.
func (r *WidgetIntegrationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Input parameters for creating a new Widget integration",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the widgetintegration",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"alias": schema.StringAttribute{
				MarkdownDescription: "The alias ID",
				Optional:            true,
			},
			"attachments": schema.BoolAttribute{
				MarkdownDescription: "Whether attachments are enabled",
				Optional:            true,
			},
			"auto_scroll": schema.BoolAttribute{
				MarkdownDescription: "Whether auto-scroll is enabled",
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
			"carousel": schema.BoolAttribute{
				MarkdownDescription: "Whether the carousel is enabled",
				Optional:            true,
			},
			"contact_collection": schema.BoolAttribute{
				MarkdownDescription: "Whether to collect contact information",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description",
				Optional:            true,
			},
			"export_conversation": schema.BoolAttribute{
				MarkdownDescription: "Whether conversation export is enabled",
				Optional:            true,
			},
			"form": schema.BoolAttribute{
				MarkdownDescription: "Whether forms are enabled",
				Optional:            true,
			},
			"initial": schema.StringAttribute{
				MarkdownDescription: "The initial message",
				Optional:            true,
			},
			"intro": schema.StringAttribute{
				MarkdownDescription: "The widget intro message",
				Optional:            true,
			},
			"language": schema.StringAttribute{
				MarkdownDescription: "The widget language",
				Optional:            true,
			},
			"layout": schema.StringAttribute{
				MarkdownDescription: "The widget layout",
				Optional:            true,
			},
			"math": schema.BoolAttribute{
				MarkdownDescription: "Whether math rendering is enabled",
				Optional:            true,
			},
			"maximize": schema.BoolAttribute{
				MarkdownDescription: "Whether the widget can be maximized",
				Optional:            true,
			},
			"message_peek": schema.BoolAttribute{
				MarkdownDescription: "Whether message peek is enabled",
				Optional:            true,
			},
			"meta": schema.MapAttribute{
				MarkdownDescription: "Additional metadata for the integration",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name",
				Optional:            true,
			},
			"origin": schema.StringAttribute{
				MarkdownDescription: "The allowed origin",
				Optional:            true,
			},
			"placeholder": schema.StringAttribute{
				MarkdownDescription: "The input placeholder text",
				Optional:            true,
			},
			"plugins": schema.StringAttribute{
				MarkdownDescription: "The enabled plugins",
				Optional:            true,
			},
			"powered_by": schema.BoolAttribute{
				MarkdownDescription: "Whether the powered-by label is shown",
				Optional:            true,
			},
			"restart_conversation": schema.BoolAttribute{
				MarkdownDescription: "Whether conversation restart is enabled",
				Optional:            true,
			},
			"session_duration": schema.Int64Attribute{
				MarkdownDescription: "The duration of the session in milliseconds",
				Optional:            true,
			},
			"start_first": schema.BoolAttribute{
				MarkdownDescription: "Whether to start first",
				Optional:            true,
			},
			"stream": schema.BoolAttribute{
				MarkdownDescription: "Whether to stream responses",
				Optional:            true,
			},
			"theme": schema.StringAttribute{
				MarkdownDescription: "The widget theme",
				Optional:            true,
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "The widget title",
				Optional:            true,
			},
			"tools": schema.BoolAttribute{
				MarkdownDescription: "Whether tools are enabled",
				Optional:            true,
			},
			"unfurl": schema.BoolAttribute{
				MarkdownDescription: "Whether link unfurling is enabled",
				Optional:            true,
			},
			"verbose": schema.BoolAttribute{
				MarkdownDescription: "Whether verbose mode is enabled",
				Optional:            true,
			},
			"voice_in": schema.BoolAttribute{
				MarkdownDescription: "Whether voice input is enabled",
				Optional:            true,
			},
			"voice_out": schema.BoolAttribute{
				MarkdownDescription: "Whether voice output is enabled",
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
func (r *WidgetIntegrationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *WidgetIntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data WidgetIntegrationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to create widgetintegration

	result, err := r.client.CreateWidgetIntegration(ctx, CreateWidgetIntegrationInput{
		Alias:               data.Alias.ValueStringPointer(),
		Attachments:         data.Attachments.ValueBoolPointer(),
		AutoScroll:          data.AutoScroll.ValueBoolPointer(),
		BlueprintId:         data.BlueprintId.ValueStringPointer(),
		BotId:               data.BotId.ValueStringPointer(),
		Carousel:            data.Carousel.ValueBoolPointer(),
		ContactCollection:   data.ContactCollection.ValueBoolPointer(),
		Description:         data.Description.ValueStringPointer(),
		ExportConversation:  data.ExportConversation.ValueBoolPointer(),
		Form:                data.Form.ValueBoolPointer(),
		Initial:             data.Initial.ValueStringPointer(),
		Intro:               data.Intro.ValueStringPointer(),
		Language:            data.Language.ValueStringPointer(),
		Layout:              data.Layout.ValueStringPointer(),
		Math:                data.Math.ValueBoolPointer(),
		Maximize:            data.Maximize.ValueBoolPointer(),
		MessagePeek:         data.MessagePeek.ValueBoolPointer(),
		Meta:                convertMapToInterface(ctx, data.Meta),
		Name:                data.Name.ValueStringPointer(),
		Origin:              data.Origin.ValueStringPointer(),
		Placeholder:         data.Placeholder.ValueStringPointer(),
		Plugins:             data.Plugins.ValueStringPointer(),
		PoweredBy:           data.PoweredBy.ValueBoolPointer(),
		RestartConversation: data.RestartConversation.ValueBoolPointer(),
		SessionDuration:     data.SessionDuration.ValueInt64Pointer(),
		StartFirst:          data.StartFirst.ValueBoolPointer(),
		Stream:              data.Stream.ValueBoolPointer(),
		Theme:               data.Theme.ValueStringPointer(),
		Title:               data.Title.ValueStringPointer(),
		Tools:               data.Tools.ValueBoolPointer(),
		Unfurl:              data.Unfurl.ValueBoolPointer(),
		Verbose:             data.Verbose.ValueBoolPointer(),
		VoiceIn:             data.VoiceIn.ValueBoolPointer(),
		VoiceOut:            data.VoiceOut.ValueBoolPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create widgetintegration: %s", err))
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
func (r *WidgetIntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data WidgetIntegrationResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to read widgetintegration

	result, err := r.client.GetWidgetIntegration(ctx, data.ID.ValueString())
	if err != nil {
		// Check if resource was deleted outside of Terraform
		if strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read widgetintegration: %s", err))
		return
	}

	// Update data model with response values

	if result.Alias != nil {
		data.Alias = types.StringPointerValue(result.Alias)
	}
	if result.Attachments != nil {
		data.Attachments = types.BoolPointerValue(result.Attachments)
	}
	if result.AutoScroll != nil {
		data.AutoScroll = types.BoolPointerValue(result.AutoScroll)
	}
	if result.BlueprintId != nil {
		data.BlueprintId = types.StringPointerValue(result.BlueprintId)
	}
	if result.BotId != nil {
		data.BotId = types.StringPointerValue(result.BotId)
	}
	if result.Carousel != nil {
		data.Carousel = types.BoolPointerValue(result.Carousel)
	}
	if result.ContactCollection != nil {
		data.ContactCollection = types.BoolPointerValue(result.ContactCollection)
	}
	if result.Description != nil {
		data.Description = types.StringPointerValue(result.Description)
	}
	if result.ExportConversation != nil {
		data.ExportConversation = types.BoolPointerValue(result.ExportConversation)
	}
	if result.Form != nil {
		data.Form = types.BoolPointerValue(result.Form)
	}
	if result.Initial != nil {
		data.Initial = types.StringPointerValue(result.Initial)
	}
	if result.Intro != nil {
		data.Intro = types.StringPointerValue(result.Intro)
	}
	if result.Language != nil {
		data.Language = types.StringPointerValue(result.Language)
	}
	if result.Layout != nil {
		data.Layout = types.StringPointerValue(result.Layout)
	}
	if result.Math != nil {
		data.Math = types.BoolPointerValue(result.Math)
	}
	if result.Maximize != nil {
		data.Maximize = types.BoolPointerValue(result.Maximize)
	}
	if result.MessagePeek != nil {
		data.MessagePeek = types.BoolPointerValue(result.MessagePeek)
	}
	if result.Meta != nil {
		mapValue, diags := types.MapValueFrom(ctx, types.StringType, result.Meta)
		resp.Diagnostics.Append(diags...)
		data.Meta = mapValue
	}
	if result.Name != nil {
		data.Name = types.StringPointerValue(result.Name)
	}
	if result.Origin != nil {
		data.Origin = types.StringPointerValue(result.Origin)
	}
	if result.Placeholder != nil {
		data.Placeholder = types.StringPointerValue(result.Placeholder)
	}
	if result.Plugins != nil {
		data.Plugins = types.StringPointerValue(result.Plugins)
	}
	if result.PoweredBy != nil {
		data.PoweredBy = types.BoolPointerValue(result.PoweredBy)
	}
	if result.RestartConversation != nil {
		data.RestartConversation = types.BoolPointerValue(result.RestartConversation)
	}
	if result.SessionDuration != nil {
		data.SessionDuration = types.Int64PointerValue(result.SessionDuration)
	}
	if result.StartFirst != nil {
		data.StartFirst = types.BoolPointerValue(result.StartFirst)
	}
	if result.Stream != nil {
		data.Stream = types.BoolPointerValue(result.Stream)
	}
	if result.Theme != nil {
		data.Theme = types.StringPointerValue(result.Theme)
	}
	if result.Title != nil {
		data.Title = types.StringPointerValue(result.Title)
	}
	if result.Tools != nil {
		data.Tools = types.BoolPointerValue(result.Tools)
	}
	if result.Unfurl != nil {
		data.Unfurl = types.BoolPointerValue(result.Unfurl)
	}
	if result.Verbose != nil {
		data.Verbose = types.BoolPointerValue(result.Verbose)
	}
	if result.VoiceIn != nil {
		data.VoiceIn = types.BoolPointerValue(result.VoiceIn)
	}
	if result.VoiceOut != nil {
		data.VoiceOut = types.BoolPointerValue(result.VoiceOut)
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
func (r *WidgetIntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data WidgetIntegrationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to update widgetintegration

	_, err := r.client.UpdateWidgetIntegration(ctx, data.ID.ValueString(), UpdateWidgetIntegrationInput{
		Alias:               data.Alias.ValueStringPointer(),
		Attachments:         data.Attachments.ValueBoolPointer(),
		AutoScroll:          data.AutoScroll.ValueBoolPointer(),
		BlueprintId:         data.BlueprintId.ValueStringPointer(),
		BotId:               data.BotId.ValueStringPointer(),
		Carousel:            data.Carousel.ValueBoolPointer(),
		ContactCollection:   data.ContactCollection.ValueBoolPointer(),
		Description:         data.Description.ValueStringPointer(),
		ExportConversation:  data.ExportConversation.ValueBoolPointer(),
		Form:                data.Form.ValueBoolPointer(),
		Initial:             data.Initial.ValueStringPointer(),
		Intro:               data.Intro.ValueStringPointer(),
		Language:            data.Language.ValueStringPointer(),
		Layout:              data.Layout.ValueStringPointer(),
		Math:                data.Math.ValueBoolPointer(),
		Maximize:            data.Maximize.ValueBoolPointer(),
		MessagePeek:         data.MessagePeek.ValueBoolPointer(),
		Meta:                convertMapToInterface(ctx, data.Meta),
		Name:                data.Name.ValueStringPointer(),
		Origin:              data.Origin.ValueStringPointer(),
		Placeholder:         data.Placeholder.ValueStringPointer(),
		Plugins:             data.Plugins.ValueStringPointer(),
		PoweredBy:           data.PoweredBy.ValueBoolPointer(),
		RestartConversation: data.RestartConversation.ValueBoolPointer(),
		SessionDuration:     data.SessionDuration.ValueInt64Pointer(),
		StartFirst:          data.StartFirst.ValueBoolPointer(),
		Stream:              data.Stream.ValueBoolPointer(),
		Theme:               data.Theme.ValueStringPointer(),
		Title:               data.Title.ValueStringPointer(),
		Tools:               data.Tools.ValueBoolPointer(),
		Unfurl:              data.Unfurl.ValueBoolPointer(),
		Verbose:             data.Verbose.ValueBoolPointer(),
		VoiceIn:             data.VoiceIn.ValueBoolPointer(),
		VoiceOut:            data.VoiceOut.ValueBoolPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update widgetintegration: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *WidgetIntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data WidgetIntegrationResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to delete widgetintegration

	_, err := r.client.DeleteWidgetIntegration(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete widgetintegration: %s", err))
		return
	}
}

// ImportState imports the resource state from Terraform.
func (r *WidgetIntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
