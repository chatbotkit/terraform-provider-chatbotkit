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
	_ datasource.DataSource              = &secretsDataSource{}
	_ datasource.DataSourceWithConfigure = &secretsDataSource{}
)

func NewSecretsDataSource() datasource.DataSource {
	return &secretsDataSource{}
}

type secretsDataSource struct {
	client *client.Client
}

type secretsDataSourceModel struct {
	Secrets []secretDataSourceModel `tfsdk:"secrets"`
}

func (d *secretsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_secrets"
}

func (d *secretsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists ChatBotKit secrets.",
		Attributes: map[string]schema.Attribute{
			"secrets": schema.ListNestedAttribute{
				Description: "List of secrets",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Secret identifier",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Secret name",
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

func (d *secretsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *secretsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data secretsDataSourceModel

	result, err := d.client.ListSecrets(ctx, "")
	if err != nil {
		resp.Diagnostics.AddError("Error listing secrets", err.Error())
		return
	}

	for _, secret := range result.Items {
		secretModel := secretDataSourceModel{
			ID:        types.StringValue(secret.ID),
			Name:      types.StringValue(secret.Name),
			CreatedAt: types.Int64Value(secret.CreatedAt),
			UpdatedAt: types.Int64Value(secret.UpdatedAt),
		}
		data.Secrets = append(data.Secrets, secretModel)
	}

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
