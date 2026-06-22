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
	_ resource.Resource                = &SpaceSiteResource{}
	_ resource.ResourceWithImportState = &SpaceSiteResource{}
)

func NewSpaceSiteResource() resource.Resource {
	return &SpaceSiteResource{}
}

// SpaceSiteResource defines the resource implementation.
type SpaceSiteResource struct {
	client *Client
}

// SpaceSiteResourceModel describes the resource data model.
type SpaceSiteResourceModel struct {
	ID types.String `tfsdk:"id"`

	SpaceId     types.String `tfsdk:"space_id"`
	Alias       types.String `tfsdk:"alias"`
	Description types.String `tfsdk:"description"`
	Domain      types.String `tfsdk:"domain"`
	Index       types.String `tfsdk:"index"`
	Meta        types.Map    `tfsdk:"meta"`
	Name        types.String `tfsdk:"name"`
	NotFound    types.String `tfsdk:"not_found"`
	Prefix      types.String `tfsdk:"prefix"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
}

// Metadata returns the resource type name.
func (r *SpaceSiteResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_space_site"
}

// Schema defines the schema for the resource.
func (r *SpaceSiteResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Input parameters for creating a new space site",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the spacesite",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"space_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the space to attach this site to",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"alias": schema.StringAttribute{
				MarkdownDescription: "The alias ID for the space site",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the space site",
				Optional:            true,
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "The host the site is served at (a <label>.chatbotkit.space subdomain)",
				Required:            true,
			},
			"index": schema.StringAttribute{
				MarkdownDescription: "The directory index filename",
				Optional:            true,
			},
			"meta": schema.MapAttribute{
				MarkdownDescription: "Additional metadata for the space site",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the space site",
				Optional:            true,
			},
			"not_found": schema.StringAttribute{
				MarkdownDescription: "The not found filename",
				Optional:            true,
			},
			"prefix": schema.StringAttribute{
				MarkdownDescription: "The optional folder prefix inside the space to serve from",
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
func (r *SpaceSiteResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *SpaceSiteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SpaceSiteResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to create spacesite

	result, err := r.client.CreateSpaceSite(ctx, data.SpaceId.ValueString(), CreateSpaceSiteInput{
		Alias:       data.Alias.ValueStringPointer(),
		Description: data.Description.ValueStringPointer(),
		Domain:      data.Domain.ValueString(),
		Index:       data.Index.ValueStringPointer(),
		Meta:        convertMapToInterface(ctx, data.Meta),
		Name:        data.Name.ValueStringPointer(),
		NotFound:    data.NotFound.ValueStringPointer(),
		Prefix:      data.Prefix.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create spacesite: %s", err))
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
func (r *SpaceSiteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SpaceSiteResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to read spacesite

	result, err := r.client.GetSpaceSite(ctx, data.SpaceId.ValueString(), data.ID.ValueString())
	if err != nil {
		// Check if resource was deleted outside of Terraform
		if strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read spacesite: %s", err))
		return
	}

	// Update data model with response values

	if result.Alias != nil {
		data.Alias = types.StringPointerValue(result.Alias)
	}
	if result.Description != nil {
		data.Description = types.StringPointerValue(result.Description)
	}
	data.Domain = types.StringValue(result.Domain)
	if result.Index != nil {
		data.Index = types.StringPointerValue(result.Index)
	}
	if result.Meta != nil {
		mapValue, diags := types.MapValueFrom(ctx, types.StringType, result.Meta)
		resp.Diagnostics.Append(diags...)
		data.Meta = mapValue
	}
	if result.Name != nil {
		data.Name = types.StringPointerValue(result.Name)
	}
	if result.NotFound != nil {
		data.NotFound = types.StringPointerValue(result.NotFound)
	}
	if result.Prefix != nil {
		data.Prefix = types.StringPointerValue(result.Prefix)
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
func (r *SpaceSiteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SpaceSiteResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to update spacesite

	_, err := r.client.UpdateSpaceSite(ctx, data.SpaceId.ValueString(), data.ID.ValueString(), UpdateSpaceSiteInput{
		Alias:       data.Alias.ValueStringPointer(),
		Description: data.Description.ValueStringPointer(),
		Domain:      data.Domain.ValueString(),
		Index:       data.Index.ValueStringPointer(),
		Meta:        convertMapToInterface(ctx, data.Meta),
		Name:        data.Name.ValueStringPointer(),
		NotFound:    data.NotFound.ValueStringPointer(),
		Prefix:      data.Prefix.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update spacesite: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *SpaceSiteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SpaceSiteResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the ChatBotKit GraphQL API to delete spacesite

	_, err := r.client.DeleteSpaceSite(ctx, data.SpaceId.ValueString(), data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete spacesite: %s", err))
		return
	}
}

// ImportState imports the resource state from Terraform.
func (r *SpaceSiteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
