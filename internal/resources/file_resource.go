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

var (
	_ resource.Resource                = &fileResource{}
	_ resource.ResourceWithConfigure   = &fileResource{}
	_ resource.ResourceWithImportState = &fileResource{}
)

func NewFileResource() resource.Resource {
	return &fileResource{}
}

type fileResource struct {
	client *client.Client
}

type fileResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Type      types.String `tfsdk:"type"`
	Source    types.String `tfsdk:"source"`
	CreatedAt types.Int64  `tfsdk:"created_at"`
	UpdatedAt types.Int64  `tfsdk:"updated_at"`
}

func (r *fileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file"
}

func (r *fileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a ChatBotKit file.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "File identifier",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "File name",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "File type",
				Optional:    true,
			},
			"source": schema.StringAttribute{
				Description: "File source URL or path",
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

func (r *fileResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *fileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan fileResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	file := &client.File{
		Name:   plan.Name.ValueString(),
		Type:   plan.Type.ValueString(),
		Source: plan.Source.ValueString(),
	}

	created, err := r.client.CreateFile(ctx, file)
	if err != nil {
		resp.Diagnostics.AddError("Error creating file", err.Error())
		return
	}

	plan.ID = types.StringValue(created.ID)
	plan.CreatedAt = types.Int64Value(created.CreatedAt)
	plan.UpdatedAt = types.Int64Value(created.UpdatedAt)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *fileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state fileResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	file, err := r.client.GetFile(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading file", err.Error())
		return
	}

	state.Name = types.StringValue(file.Name)
	state.Type = types.StringValue(file.Type)
	state.Source = types.StringValue(file.Source)
	state.CreatedAt = types.Int64Value(file.CreatedAt)
	state.UpdatedAt = types.Int64Value(file.UpdatedAt)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *fileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan fileResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	file := &client.File{
		Name:   plan.Name.ValueString(),
		Type:   plan.Type.ValueString(),
		Source: plan.Source.ValueString(),
	}

	updated, err := r.client.UpdateFile(ctx, plan.ID.ValueString(), file)
	if err != nil {
		resp.Diagnostics.AddError("Error updating file", err.Error())
		return
	}

	plan.UpdatedAt = types.Int64Value(updated.UpdatedAt)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *fileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state fileResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteFile(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting file", err.Error())
		return
	}
}

func (r *fileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
