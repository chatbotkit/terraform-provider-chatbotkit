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
	_ resource.Resource                = &SitemapIntegrationResource{}
	_ resource.ResourceWithImportState = &SitemapIntegrationResource{}
)

func NewSitemapIntegrationResource() resource.Resource {
	return &SitemapIntegrationResource{}
}

// SitemapIntegrationResource defines the resource implementation.
type SitemapIntegrationResource struct {
	client *Client
}

// SitemapIntegrationResourceModel describes the resource data model.
type SitemapIntegrationResourceModel struct {
	ID types.String `tfsdk:"id"`

	BlueprintId types.String `tfsdk:"blueprint_id"`
	DatasetId types.String `tfsdk:"dataset_id"`
	Description types.String `tfsdk:"description"`
	ExpiresIn types.Int64 `tfsdk:"expires_in"`
	Glob types.String `tfsdk:"glob"`
	Javascript types.Bool `tfsdk:"javascript"`
	Meta types.Map `tfsdk:"meta"`
	Name types.String `tfsdk:"name"`
	Selectors types.String `tfsdk:"selectors"`
	SyncSchedule types.String `tfsdk:"sync_schedule"`
	URL types.String `tfsdk:"url"`
}

// Metadata returns the resource type name.
func (r *SitemapIntegrationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sitemap_integration"
}

// Schema defines the schema for the resource.
func (r *SitemapIntegrationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Input parameters for creating a new Sitemap integration",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the sitemapintegration",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"blueprint_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the blueprint to use",
				Optional:            true,
			},
			"dataset_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the dataset to sync to",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the integration",
				Optional:            true,
			},
			"expires_in": schema.Int64Attribute{
				MarkdownDescription: "Time in milliseconds before the data expires",
				Optional:            true,
			},
			"glob": schema.StringAttribute{
				MarkdownDescription: "Glob pattern to filter URLs",
				Optional:            true,
			},
			"javascript": schema.BoolAttribute{
				MarkdownDescription: "Whether to enable JavaScript rendering",
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
			"selectors": schema.StringAttribute{
				MarkdownDescription: "CSS selectors to focus on specific parts of the pages",
				Optional:            true,
			},
			"sync_schedule": schema.StringAttribute{
				MarkdownDescription: "The schedule for automatic synchronization",
				Optional:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "The URL of the sitemap to crawl",
				Optional:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *SitemapIntegrationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *SitemapIntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SitemapIntegrationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Implement API call to create sitemapintegration
	// Example:
	// result, err := r.client.CreateSitemapIntegration(ctx, types.SitemapIntegrationCreateRequest{

	//     BlueprintId: data.BlueprintId.ValueStringPointer(),
	//     DatasetId: data.DatasetId.ValueStringPointer(),
	//     Description: data.Description.ValueStringPointer(),
	//     ExpiresIn: data.ExpiresIn.ValueInt64Pointer(),
	//     Glob: data.Glob.ValueStringPointer(),
	//     Javascript: data.Javascript.ValueBoolPointer(),
	//     Meta: data.Meta.Elements(),
	//     Name: data.Name.ValueStringPointer(),
	//     Selectors: data.Selectors.ValueStringPointer(),
	//     SyncSchedule: data.SyncSchedule.ValueStringPointer(),
	//     URL: data.URL.ValueStringPointer(),
	// })
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create sitemapintegration: %s", err))
	//     return
	// }
	// data.ID = types.StringValue(result.ID)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *SitemapIntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SitemapIntegrationResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Implement API call to read sitemapintegration
	// Example:
	// result, err := r.client.GetSitemapIntegration(ctx, data.ID.ValueString())
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read sitemapintegration: %s", err))
	//     return
	// }
	// Update data model with response values

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *SitemapIntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SitemapIntegrationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Implement API call to update sitemapintegration
	// Example:
	// _, err := r.client.UpdateSitemapIntegration(ctx, data.ID.ValueString(), types.SitemapIntegrationUpdateRequest{

	//     BlueprintId: data.BlueprintId.ValueStringPointer(),
	//     DatasetId: data.DatasetId.ValueStringPointer(),
	//     Description: data.Description.ValueStringPointer(),
	//     ExpiresIn: data.ExpiresIn.ValueInt64Pointer(),
	//     Glob: data.Glob.ValueStringPointer(),
	//     Javascript: data.Javascript.ValueBoolPointer(),
	//     Meta: data.Meta.Elements(),
	//     Name: data.Name.ValueStringPointer(),
	//     Selectors: data.Selectors.ValueStringPointer(),
	//     SyncSchedule: data.SyncSchedule.ValueStringPointer(),
	//     URL: data.URL.ValueStringPointer(),
	// })
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update sitemapintegration: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *SitemapIntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SitemapIntegrationResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Implement API call to delete sitemapintegration
	// Example:
	// _, err := r.client.DeleteSitemapIntegration(ctx, data.ID.ValueString())
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete sitemapintegration: %s", err))
	//     return
	// }
}

// ImportState imports the resource state from Terraform.
func (r *SitemapIntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
