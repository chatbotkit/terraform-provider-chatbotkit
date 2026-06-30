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
	_ resource.Resource                = &TaskResource{}
	_ resource.ResourceWithImportState = &TaskResource{}
)

func NewTaskResource() resource.Resource {
	return &TaskResource{}
}

// TaskResource defines the resource implementation.
type TaskResource struct {
	client *Client
}

// TaskResourceModel describes the resource data model.
type TaskResourceModel struct {
	ID types.String `tfsdk:"id"`

	BlueprintId     types.String  `tfsdk:"blueprint_id"`
	BotId           types.String  `tfsdk:"bot_id"`
	ContactId       types.String  `tfsdk:"contact_id"`
	Description     types.String  `tfsdk:"description"`
	MaxCalls        types.Int64   `tfsdk:"max_calls"`
	MaxIterations   types.Int64   `tfsdk:"max_iterations"`
	MaxTime         types.Float64 `tfsdk:"max_time"`
	Meta            types.Map     `tfsdk:"meta"`
	Name            types.String  `tfsdk:"name"`
	Schedule        types.String  `tfsdk:"schedule"`
	SessionDuration types.Float64 `tfsdk:"session_duration"`
	Timezone        types.String  `tfsdk:"timezone"`
	CreatedAt       types.String  `tfsdk:"created_at"`
	UpdatedAt       types.String  `tfsdk:"updated_at"`
}

// Metadata returns the resource type name.
func (r *TaskResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_task"
}

// Schema defines the schema for the resource.
func (r *TaskResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Input parameters for creating a new task",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the task",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"blueprint_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the blueprint to assign the task to",
				Optional:            true,
			},
			"bot_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the bot the task runs",
				Optional:            true,
			},
			"contact_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the contact to scope the task to",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the task",
				Optional:            true,
			},
			"max_calls": schema.Int64Attribute{
				MarkdownDescription: "Maximum tool calls across the whole task run (0 or null for unbounded)",
				Optional:            true,
			},
			"max_iterations": schema.Int64Attribute{
				MarkdownDescription: "Maximum reasoning iterations per execution",
				Optional:            true,
			},
			"max_time": schema.Float64Attribute{
				MarkdownDescription: "Maximum wall-clock time per execution in milliseconds",
				Optional:            true,
			},
			"meta": schema.MapAttribute{
				MarkdownDescription: "Additional metadata for the task",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the task",
				Optional:            true,
			},
			"schedule": schema.StringAttribute{
				MarkdownDescription: "The schedule: now, a cron expression, a date-time, or an interval keyword such as daily",
				Optional:            true,
			},
			"session_duration": schema.Float64Attribute{
				MarkdownDescription: "Session duration in milliseconds controlling conversation reuse across runs",
				Optional:            true,
			},
			"timezone": schema.StringAttribute{
				MarkdownDescription: "The IANA timezone the schedule is evaluated in",
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
func (r *TaskResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *TaskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TaskResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to create task

	result, err := r.client.CreateTask(ctx, CreateTaskInput{
		BlueprintId:     data.BlueprintId.ValueStringPointer(),
		BotId:           data.BotId.ValueStringPointer(),
		ContactId:       data.ContactId.ValueStringPointer(),
		Description:     data.Description.ValueStringPointer(),
		MaxCalls:        data.MaxCalls.ValueInt64Pointer(),
		MaxIterations:   data.MaxIterations.ValueInt64Pointer(),
		MaxTime:         data.MaxTime.ValueFloat64Pointer(),
		Meta:            convertMapToInterface(ctx, data.Meta),
		Name:            data.Name.ValueStringPointer(),
		Schedule:        data.Schedule.ValueStringPointer(),
		SessionDuration: data.SessionDuration.ValueFloat64Pointer(),
		Timezone:        data.Timezone.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create task: %s", err))
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
func (r *TaskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TaskResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to read task

	result, err := r.client.GetTask(ctx, data.ID.ValueString())
	if err != nil {
		// Check if resource was deleted outside of Terraform
		if strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read task: %s", err))
		return
	}

	// Update data model with response values

	if result.BlueprintId != nil {
		data.BlueprintId = types.StringPointerValue(result.BlueprintId)
	}
	if result.BotId != nil {
		data.BotId = types.StringPointerValue(result.BotId)
	}
	if result.ContactId != nil {
		data.ContactId = types.StringPointerValue(result.ContactId)
	}
	if result.Description != nil {
		data.Description = types.StringPointerValue(result.Description)
	}
	if result.MaxCalls != nil {
		data.MaxCalls = types.Int64PointerValue(result.MaxCalls)
	}
	if result.MaxIterations != nil {
		data.MaxIterations = types.Int64PointerValue(result.MaxIterations)
	}
	if result.MaxTime != nil {
		data.MaxTime = types.Float64PointerValue(result.MaxTime)
	}
	if result.Meta != nil {
		mapValue, diags := types.MapValueFrom(ctx, types.StringType, result.Meta)
		resp.Diagnostics.Append(diags...)
		data.Meta = mapValue
	}
	if result.Name != nil {
		data.Name = types.StringPointerValue(result.Name)
	}
	if result.Schedule != nil {
		data.Schedule = types.StringPointerValue(result.Schedule)
	}
	if result.SessionDuration != nil {
		data.SessionDuration = types.Float64PointerValue(result.SessionDuration)
	}
	if result.Timezone != nil {
		data.Timezone = types.StringPointerValue(result.Timezone)
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
func (r *TaskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TaskResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to update task

	_, err := r.client.UpdateTask(ctx, data.ID.ValueString(), UpdateTaskInput{
		BlueprintId:     data.BlueprintId.ValueStringPointer(),
		BotId:           data.BotId.ValueStringPointer(),
		ContactId:       data.ContactId.ValueStringPointer(),
		Description:     data.Description.ValueStringPointer(),
		MaxCalls:        data.MaxCalls.ValueInt64Pointer(),
		MaxIterations:   data.MaxIterations.ValueInt64Pointer(),
		MaxTime:         data.MaxTime.ValueFloat64Pointer(),
		Meta:            convertMapToInterface(ctx, data.Meta),
		Name:            data.Name.ValueStringPointer(),
		Schedule:        data.Schedule.ValueStringPointer(),
		SessionDuration: data.SessionDuration.ValueFloat64Pointer(),
		Timezone:        data.Timezone.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update task: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *TaskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TaskResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to delete task

	_, err := r.client.DeleteTask(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete task: %s", err))
		return
	}
}

// ImportState imports the resource state from Terraform.
func (r *TaskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
