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
	_ datasource.DataSource              = &fileDataSource{}
	_ datasource.DataSourceWithConfigure = &fileDataSource{}
)

func NewFileDataSource() datasource.DataSource {
	return &fileDataSource{}
}

type fileDataSource struct {
	client *client.Client
}

type fileDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Type      types.String `tfsdk:"type"`
	Source    types.String `tfsdk:"source"`
	CreatedAt types.Int64  `tfsdk:"created_at"`
	UpdatedAt types.Int64  `tfsdk:"updated_at"`
}

func (d *fileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file"
}

func (d *fileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a ChatBotKit file.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "File identifier",
				Required:    true,
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
	}
}

func (d *fileDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *fileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data fileDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	file, err := d.client.GetFile(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading file", err.Error())
		return
	}

	data.Name = types.StringValue(file.Name)
	data.Type = types.StringValue(file.Type)
	data.Source = types.StringValue(file.Source)
	data.CreatedAt = types.Int64Value(file.CreatedAt)
	data.UpdatedAt = types.Int64Value(file.UpdatedAt)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
