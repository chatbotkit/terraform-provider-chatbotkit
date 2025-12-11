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
	_ datasource.DataSource              = &datasetsDataSource{}
	_ datasource.DataSourceWithConfigure = &datasetsDataSource{}
)

func NewDatasetsDataSource() datasource.DataSource {
	return &datasetsDataSource{}
}

type datasetsDataSource struct {
	client *client.Client
}

type datasetsDataSourceModel struct {
	Datasets []datasetDataSourceModel `tfsdk:"datasets"`
}

func (d *datasetsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_datasets"
}

func (d *datasetsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists ChatBotKit datasets.",
		Attributes: map[string]schema.Attribute{
			"datasets": schema.ListNestedAttribute{
				Description: "List of datasets",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Dataset identifier",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Dataset name",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Dataset description",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "Dataset type",
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

func (d *datasetsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *datasetsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasetsDataSourceModel

	result, err := d.client.ListDatasets(ctx, "")
	if err != nil {
		resp.Diagnostics.AddError("Error listing datasets", err.Error())
		return
	}

	for _, dataset := range result.Items {
		datasetModel := datasetDataSourceModel{
			ID:          types.StringValue(dataset.ID),
			Name:        types.StringValue(dataset.Name),
			Description: types.StringValue(dataset.Description),
			Type:        types.StringValue(dataset.Type),
			CreatedAt:   types.Int64Value(dataset.CreatedAt),
			UpdatedAt:   types.Int64Value(dataset.UpdatedAt),
		}
		data.Datasets = append(data.Datasets, datasetModel)
	}

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
