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
	_ resource.Resource                = &SkillsetResource{}
	_ resource.ResourceWithImportState = &SkillsetResource{}
)

func NewSkillsetResource() resource.Resource {
	return &SkillsetResource{}
}

// SkillsetResource defines the resource implementation.
type SkillsetResource struct {
	client *Client
}

// SkillsetResourceModel describes the resource data model.
type SkillsetResourceModel struct {
	ID types.String `tfsdk:"id"`

	BlueprintId types.String `tfsdk:"blueprint_id"`
	Description types.String `tfsdk:"description"`
	Meta types.Map `tfsdk:"meta"`
	Name types.String `tfsdk:"name"`
	Visibility types.String `tfsdk:"visibility"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

// Metadata returns the resource type name.
func (r *SkillsetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_skillset"
}

// Schema defines the schema for the resource.
func (r *SkillsetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Input parameters for creating a new skillset",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the skillset",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"blueprint_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the blueprint to use",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the skillset",
				Optional:            true,
			},
			"meta": schema.MapAttribute{
				MarkdownDescription: "Additional metadata for the skillset",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the skillset",
				Optional:            true,
			},
			"visibility": schema.StringAttribute{
				MarkdownDescription: "The visibility level of the skillset",
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
func (r *SkillsetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *SkillsetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SkillsetResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to create skillset
	result, err := r.client.CreateSkillset(ctx, CreateSkillsetInput{

		BlueprintId: data.BlueprintId.ValueStringPointer(),
		Description: data.Description.ValueStringPointer(),
		Meta: convertMapToInterface(ctx, data.Meta),
		Name: data.Name.ValueStringPointer(),
		Visibility: data.Visibility.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create skillset: %s", err))
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
func (r *SkillsetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SkillsetResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to read skillset
	result, err := r.client.GetSkillset(ctx, data.ID.ValueString())
	if err != nil {
		// Check if resource was deleted outside of Terraform
		if strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read skillset: %s", err))
		return
	}

	// Update data model with response values

	if result.BlueprintId != nil {
		data.BlueprintId = types.StringPointerValue(result.BlueprintId)
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
	if result.Visibility != nil {
		data.Visibility = types.StringPointerValue(result.Visibility)
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
func (r *SkillsetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SkillsetResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to update skillset
	_, err := r.client.UpdateSkillset(ctx, data.ID.ValueString(), UpdateSkillsetInput{

		BlueprintId: data.BlueprintId.ValueStringPointer(),
		Description: data.Description.ValueStringPointer(),
		Meta: convertMapToInterface(ctx, data.Meta),
		Name: data.Name.ValueStringPointer(),
		Visibility: data.Visibility.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update skillset: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *SkillsetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SkillsetResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to delete skillset
	_, err := r.client.DeleteSkillset(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete skillset: %s", err))
		return
	}
}

// ImportState imports the resource state from Terraform.
func (r *SkillsetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
