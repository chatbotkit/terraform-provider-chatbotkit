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
	_ resource.Resource                = &SecretResource{}
	_ resource.ResourceWithImportState = &SecretResource{}
)

func NewSecretResource() resource.Resource {
	return &SecretResource{}
}

// SecretResource defines the resource implementation.
type SecretResource struct {
	client *Client
}

// SecretResourceModel describes the resource data model.
type SecretResourceModel struct {
	ID types.String `tfsdk:"id"`

	BlueprintId types.String `tfsdk:"blueprint_id"`
	Config types.Map `tfsdk:"config"`
	Description types.String `tfsdk:"description"`
	Kind types.String `tfsdk:"kind"`
	Meta types.Map `tfsdk:"meta"`
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
	Value types.String `tfsdk:"value"`
	Visibility types.String `tfsdk:"visibility"`
}

// Metadata returns the resource type name.
func (r *SecretResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_secret"
}

// Schema defines the schema for the resource.
func (r *SecretResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Input parameters for creating a new secret",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the secret",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"blueprint_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the blueprint to use",
				Optional:            true,
			},
			"config": schema.MapAttribute{
				MarkdownDescription: "Additional configuration for the secret",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the secret",
				Optional:            true,
			},
			"kind": schema.StringAttribute{
				MarkdownDescription: "The kind of secret (personal or organizational)",
				Optional:            true,
			},
			"meta": schema.MapAttribute{
				MarkdownDescription: "Additional metadata for the secret",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the secret",
				Optional:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of secret (token or other)",
				Optional:            true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "The secret value",
				Optional:            true,
			},
			"visibility": schema.StringAttribute{
				MarkdownDescription: "The visibility level of the secret",
				Optional:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *SecretResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *SecretResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SecretResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Implement API call to create secret
	// Example:
	// result, err := r.client.CreateSecret(ctx, types.SecretCreateRequest{

	//     BlueprintId: data.BlueprintId.ValueStringPointer(),
	//     Config: data.Config.Elements(),
	//     Description: data.Description.ValueStringPointer(),
	//     Kind: data.Kind.ValueStringPointer(),
	//     Meta: data.Meta.Elements(),
	//     Name: data.Name.ValueStringPointer(),
	//     Type: data.Type.ValueStringPointer(),
	//     Value: data.Value.ValueStringPointer(),
	//     Visibility: data.Visibility.ValueStringPointer(),
	// })
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create secret: %s", err))
	//     return
	// }
	// data.ID = types.StringValue(result.ID)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *SecretResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SecretResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Implement API call to read secret
	// Example:
	// result, err := r.client.GetSecret(ctx, data.ID.ValueString())
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read secret: %s", err))
	//     return
	// }
	// Update data model with response values

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *SecretResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SecretResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Implement API call to update secret
	// Example:
	// _, err := r.client.UpdateSecret(ctx, data.ID.ValueString(), types.SecretUpdateRequest{

	//     BlueprintId: data.BlueprintId.ValueStringPointer(),
	//     Config: data.Config.Elements(),
	//     Description: data.Description.ValueStringPointer(),
	//     Kind: data.Kind.ValueStringPointer(),
	//     Meta: data.Meta.Elements(),
	//     Name: data.Name.ValueStringPointer(),
	//     Type: data.Type.ValueStringPointer(),
	//     Value: data.Value.ValueStringPointer(),
	//     Visibility: data.Visibility.ValueStringPointer(),
	// })
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update secret: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *SecretResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SecretResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Implement API call to delete secret
	// Example:
	// _, err := r.client.DeleteSecret(ctx, data.ID.ValueString())
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete secret: %s", err))
	//     return
	// }
}

// ImportState imports the resource state from Terraform.
func (r *SecretResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
