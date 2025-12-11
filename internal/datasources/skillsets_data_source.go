package datasources

import (
	"context"
	"fmt"

	"github.com/chatbotkit/terraform-provider/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &skillsetsDataSource{}
	_ datasource.DataSourceWithConfigure = &skillsetsDataSource{}
)

func NewSkillsetsDataSource() datasource.DataSource {
	return &skillsetsDataSource{}
}

type skillsetsDataSource struct {
	client *client.Client
}

type skillsetsDataSourceModel struct {
	Skillsets []skillsetDataSourceModel `tfsdk:"skillsets"`
}

func (d *skillsetsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_skillsets"
}

func (d *skillsetsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists ChatBotKit skillsets.",
		Attributes: map[string]schema.Attribute{
			"skillsets": schema.ListNestedAttribute{
				Description: "List of skillsets",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Skillset identifier",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Skillset name",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Skillset description",
							Computed:    true,
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
				},
			},
		},
	}
}

func (d *skillsetsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *skillsetsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data skillsetsDataSourceModel

	result, err := d.client.ListSkillsets(ctx, "")
	if err != nil {
		resp.Diagnostics.AddError("Error listing skillsets", err.Error())
		return
	}

	for _, skillset := range result.Items {
		skillsetModel := skillsetDataSourceModel{
			ID:          types.StringValue(skillset.ID),
			Name:        types.StringValue(skillset.Name),
			Description: types.StringValue(skillset.Description),
			CreatedAt:   types.Int64Value(skillset.CreatedAt),
			UpdatedAt:   types.Int64Value(skillset.UpdatedAt),
		}
		data.Skillsets = append(data.Skillsets, skillsetModel)
	}

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
