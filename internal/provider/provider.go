package provider

import (
	"context"
	"os"

	"github.com/chatbotkit/terraform-provider/internal/client"
	"github.com/chatbotkit/terraform-provider/internal/resources"
	"github.com/chatbotkit/terraform-provider/internal/datasources"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &chatbotkitProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &chatbotkitProvider{
			version: version,
		}
	}
}

// chatbotkitProvider is the provider implementation.
type chatbotkitProvider struct {
	version string
}

// chatbotkitProviderModel maps provider schema data to a Go type.
type chatbotkitProviderModel struct {
	Token types.String `tfsdk:"token"`
}

// Metadata returns the provider type name.
func (p *chatbotkitProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "chatbotkit"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *chatbotkitProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Terraform provider for ChatBotKit API",
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Description: "ChatBotKit API token. Can also be set via CHATBOTKIT_TOKEN environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

// Configure prepares a ChatBotKit API client for data sources and resources.
func (p *chatbotkitProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config chatbotkitProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get token from configuration or environment variable
	token := os.Getenv("CHATBOTKIT_TOKEN")
	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	if token == "" {
		resp.Diagnostics.AddError(
			"Missing ChatBotKit API Token",
			"The provider cannot create the ChatBotKit API client as there is a missing or empty value for the token. "+
				"Set the token value in the configuration or use the CHATBOTKIT_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
		return
	}

	// Create the ChatBotKit client
	apiClient := client.NewClient(token)

	// Make the client available during DataSource and Resource type Configure methods.
	resp.DataSourceData = apiClient
	resp.ResourceData = apiClient
}

// DataSources defines the data sources implemented in the provider.
func (p *chatbotkitProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		datasources.NewBotDataSource,
		datasources.NewBotsDataSource,
		datasources.NewDatasetDataSource,
		datasources.NewDatasetsDataSource,
		datasources.NewSkillsetDataSource,
		datasources.NewSkillsetsDataSource,
		datasources.NewFileDataSource,
		datasources.NewFilesDataSource,
		datasources.NewIntegrationDataSource,
		datasources.NewIntegrationsDataSource,
		datasources.NewSecretDataSource,
		datasources.NewSecretsDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *chatbotkitProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewBotResource,
		resources.NewDatasetResource,
		resources.NewSkillsetResource,
		resources.NewFileResource,
		resources.NewIntegrationResource,
		resources.NewSecretResource,
	}
}
