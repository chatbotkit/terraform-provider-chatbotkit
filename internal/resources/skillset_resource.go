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
	_ resource.Resource                = &skillsetResource{}
	_ resource.ResourceWithConfigure   = &skillsetResource{}
	_ resource.ResourceWithImportState = &skillsetResource{}
)

func NewSkillsetResource() resource.Resource {
	return &skillsetResource{}
}

type skillsetResource struct {
	client *client.Client
}

type skillsetResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	CreatedAt   types.Int64  `tfsdk:"created_at"`
	UpdatedAt   types.Int64  `tfsdk:"updated_at"`
}

func (r *skillsetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_skillset"
}

func (r *skillsetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a ChatBotKit skillset.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Skillset identifier",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Skillset name",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Skillset description",
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

func (r *skillsetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *skillsetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan skillsetResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	skillset := &client.Skillset{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	created, err := r.client.CreateSkillset(ctx, skillset)
	if err != nil {
		resp.Diagnostics.AddError("Error creating skillset", err.Error())
		return
	}

	plan.ID = types.StringValue(created.ID)
	plan.CreatedAt = types.Int64Value(created.CreatedAt)
	plan.UpdatedAt = types.Int64Value(created.UpdatedAt)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *skillsetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state skillsetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	skillset, err := r.client.GetSkillset(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading skillset", err.Error())
		return
	}

	state.Name = types.StringValue(skillset.Name)
	state.Description = types.StringValue(skillset.Description)
	state.CreatedAt = types.Int64Value(skillset.CreatedAt)
	state.UpdatedAt = types.Int64Value(skillset.UpdatedAt)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *skillsetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan skillsetResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	skillset := &client.Skillset{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	updated, err := r.client.UpdateSkillset(ctx, plan.ID.ValueString(), skillset)
	if err != nil {
		resp.Diagnostics.AddError("Error updating skillset", err.Error())
		return
	}

	plan.UpdatedAt = types.Int64Value(updated.UpdatedAt)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *skillsetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state skillsetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSkillset(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting skillset", err.Error())
		return
	}
}

func (r *skillsetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
