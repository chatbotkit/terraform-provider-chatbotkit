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
	_ datasource.DataSource              = &filesDataSource{}
	_ datasource.DataSourceWithConfigure = &filesDataSource{}
)

func NewFilesDataSource() datasource.DataSource {
	return &filesDataSource{}
}

type filesDataSource struct {
	client *client.Client
}

type filesDataSourceModel struct {
	Files []fileDataSourceModel `tfsdk:"files"`
}

func (d *filesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_files"
}

func (d *filesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists ChatBotKit files.",
		Attributes: map[string]schema.Attribute{
			"files": schema.ListNestedAttribute{
				Description: "List of files",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "File identifier",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "File name",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "File type",
							Computed:    true,
						},
						"source": schema.StringAttribute{
							Description: "File source",
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

func (d *filesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *filesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data filesDataSourceModel

	result, err := d.client.ListFiles(ctx, "")
	if err != nil {
		resp.Diagnostics.AddError("Error listing files", err.Error())
		return
	}

	for _, file := range result.Items {
		fileModel := fileDataSourceModel{
			ID:        types.StringValue(file.ID),
			Name:      types.StringValue(file.Name),
			Type:      types.StringValue(file.Type),
			Source:    types.StringValue(file.Source),
			CreatedAt: types.Int64Value(file.CreatedAt),
			UpdatedAt: types.Int64Value(file.UpdatedAt),
		}
		data.Files = append(data.Files, fileModel)
	}

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
