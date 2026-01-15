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
	_ resource.Resource                = &BlueprintResource{}
	_ resource.ResourceWithImportState = &BlueprintResource{}
)

func NewBlueprintResource() resource.Resource {
	return &BlueprintResource{}
}

// BlueprintResource defines the resource implementation.
type BlueprintResource struct {
	client *Client
}

// BlueprintResourceModel describes the resource data model.
type BlueprintResourceModel struct {
	ID types.String `tfsdk:"id"`

	Description types.String `tfsdk:"description"`
	Meta types.Map `tfsdk:"meta"`
	Name types.String `tfsdk:"name"`
	Visibility types.String `tfsdk:"visibility"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

// Metadata returns the resource type name.
func (r *BlueprintResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_blueprint"
}

// Schema defines the schema for the resource.
func (r *BlueprintResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Input parameters for creating a new blueprint",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the blueprint",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the blueprint",
				Optional:            true,
			},
			"meta": schema.MapAttribute{
				MarkdownDescription: "Additional metadata for the blueprint",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the blueprint",
				Optional:            true,
			},
			"visibility": schema.StringAttribute{
				MarkdownDescription: "The visibility level of the blueprint",
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
func (r *BlueprintResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *BlueprintResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data BlueprintResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to create blueprint
	result, err := r.client.CreateBlueprint(ctx, CreateBlueprintInput{

		Description: data.Description.ValueStringPointer(),
		Meta: convertMapToInterface(ctx, data.Meta),
		Name: data.Name.ValueStringPointer(),
		Visibility: data.Visibility.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create blueprint: %s", err))
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
func (r *BlueprintResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data BlueprintResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to read blueprint
	result, err := r.client.GetBlueprint(ctx, data.ID.ValueString())
	if err != nil {
		// Check if resource was deleted outside of Terraform
		if strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read blueprint: %s", err))
		return
	}

	// Update data model with response values

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
func (r *BlueprintResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data BlueprintResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to update blueprint
	_, err := r.client.UpdateBlueprint(ctx, data.ID.ValueString(), UpdateBlueprintInput{

		Description: data.Description.ValueStringPointer(),
		Meta: convertMapToInterface(ctx, data.Meta),
		Name: data.Name.ValueStringPointer(),
		Visibility: data.Visibility.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update blueprint: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *BlueprintResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data BlueprintResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to delete blueprint
	_, err := r.client.DeleteBlueprint(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete blueprint: %s", err))
		return
	}
}

// ImportState imports the resource state from Terraform.
func (r *BlueprintResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
