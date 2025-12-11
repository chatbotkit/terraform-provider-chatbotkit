package resources

import (
	"context"
	"fmt"

	"github.com/chatbotkit/terraform-provider/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &botResource{}
	_ resource.ResourceWithConfigure   = &botResource{}
	_ resource.ResourceWithImportState = &botResource{}
)

// NewBotResource is a helper function to simplify the provider implementation.
func NewBotResource() resource.Resource {
	return &botResource{}
}

// botResource is the resource implementation.
type botResource struct {
	client *client.Client
}

// botResourceModel maps the resource schema data.
type botResourceModel struct {
	ID           types.String  `tfsdk:"id"`
	Name         types.String  `tfsdk:"name"`
	Description  types.String  `tfsdk:"description"`
	Model        types.String  `tfsdk:"model"`
	DatasetID    types.String  `tfsdk:"dataset_id"`
	SkillsetID   types.String  `tfsdk:"skillset_id"`
	Backstory    types.String  `tfsdk:"backstory"`
	Temperature  types.Float64 `tfsdk:"temperature"`
	Instructions types.String  `tfsdk:"instructions"`
	Moderation   types.Bool    `tfsdk:"moderation"`
	Privacy      types.Bool    `tfsdk:"privacy"`
	CreatedAt    types.Int64   `tfsdk:"created_at"`
	UpdatedAt    types.Int64   `tfsdk:"updated_at"`
}

// Metadata returns the resource type name.
func (r *botResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bot"
}

// Schema defines the schema for the resource.
func (r *botResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a ChatBotKit bot.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Bot identifier",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Bot name",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Bot description",
				Optional:    true,
			},
			"model": schema.StringAttribute{
				Description: "AI model to use for the bot",
				Optional:    true,
			},
			"dataset_id": schema.StringAttribute{
				Description: "Dataset ID to attach to the bot",
				Optional:    true,
			},
			"skillset_id": schema.StringAttribute{
				Description: "Skillset ID to attach to the bot",
				Optional:    true,
			},
			"backstory": schema.StringAttribute{
				Description: "Bot backstory",
				Optional:    true,
			},
			"temperature": schema.Float64Attribute{
				Description: "Temperature for AI responses",
				Optional:    true,
			},
			"instructions": schema.StringAttribute{
				Description: "Instructions for the bot",
				Optional:    true,
			},
			"moderation": schema.BoolAttribute{
				Description: "Enable content moderation",
				Optional:    true,
			},
			"privacy": schema.BoolAttribute{
				Description: "Enable privacy mode",
				Optional:    true,
			},
			"created_at": schema.Int64Attribute{
				Description: "Creation timestamp",
				Computed:    true,
			},
			"updated_at": schema.Int64Attribute{
				Description: "Last update timestamp",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *botResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *botResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan botResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create bot via API
	bot := &client.Bot{
		Name:         plan.Name.ValueString(),
		Description:  plan.Description.ValueString(),
		Model:        plan.Model.ValueString(),
		DatasetID:    plan.DatasetID.ValueString(),
		SkillsetID:   plan.SkillsetID.ValueString(),
		Backstory:    plan.Backstory.ValueString(),
		Temperature:  plan.Temperature.ValueFloat64(),
		Instructions: plan.Instructions.ValueString(),
		Moderation:   plan.Moderation.ValueBool(),
		Privacy:      plan.Privacy.ValueBool(),
	}

	created, err := r.client.CreateBot(ctx, bot)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating bot",
			"Could not create bot: "+err.Error(),
		)
		return
	}

	// Map response to schema
	plan.ID = types.StringValue(created.ID)
	plan.CreatedAt = types.Int64Value(created.CreatedAt)
	plan.UpdatedAt = types.Int64Value(created.UpdatedAt)

	// Set state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *botResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state botResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get bot from API
	bot, err := r.client.GetBot(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading bot",
			"Could not read bot ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map response to schema
	state.Name = types.StringValue(bot.Name)
	state.Description = types.StringValue(bot.Description)
	state.Model = types.StringValue(bot.Model)
	state.DatasetID = types.StringValue(bot.DatasetID)
	state.SkillsetID = types.StringValue(bot.SkillsetID)
	state.Backstory = types.StringValue(bot.Backstory)
	state.Temperature = types.Float64Value(bot.Temperature)
	state.Instructions = types.StringValue(bot.Instructions)
	state.Moderation = types.BoolValue(bot.Moderation)
	state.Privacy = types.BoolValue(bot.Privacy)
	state.CreatedAt = types.Int64Value(bot.CreatedAt)
	state.UpdatedAt = types.Int64Value(bot.UpdatedAt)

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *botResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan botResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update bot via API
	bot := &client.Bot{
		Name:         plan.Name.ValueString(),
		Description:  plan.Description.ValueString(),
		Model:        plan.Model.ValueString(),
		DatasetID:    plan.DatasetID.ValueString(),
		SkillsetID:   plan.SkillsetID.ValueString(),
		Backstory:    plan.Backstory.ValueString(),
		Temperature:  plan.Temperature.ValueFloat64(),
		Instructions: plan.Instructions.ValueString(),
		Moderation:   plan.Moderation.ValueBool(),
		Privacy:      plan.Privacy.ValueBool(),
	}

	updated, err := r.client.UpdateBot(ctx, plan.ID.ValueString(), bot)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating bot",
			"Could not update bot ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map response to schema
	plan.UpdatedAt = types.Int64Value(updated.UpdatedAt)

	// Set state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *botResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state botResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete bot via API
	err := r.client.DeleteBot(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting bot",
			"Could not delete bot ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}
}

// ImportState imports the resource state.
func (r *botResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
