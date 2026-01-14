package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure ChatBotKitProvider satisfies various provider interfaces.
var _ provider.Provider = &ChatBotKitProvider{}

// ChatBotKitProvider defines the provider implementation.
type ChatBotKitProvider struct {
	version string
}

// ChatBotKitProviderModel describes the provider data model.
type ChatBotKitProviderModel struct {
	APIKey  types.String `tfsdk:"api_key"`
	BaseURL types.String `tfsdk:"base_url"`
}

// Client represents the ChatBotKit API client.
type Client struct {
	APIKey  string
	BaseURL string
}

func (p *ChatBotKitProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "chatbotkit"
	resp.Version = p.version
}

func (p *ChatBotKitProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "ChatBotKit Terraform Provider for managing AI chatbot resources.",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "The API key for authenticating with ChatBotKit API. Can also be set via CHATBOTKIT_API_KEY environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"base_url": schema.StringAttribute{
				MarkdownDescription: "The base URL for the ChatBotKit API. Defaults to https://api.chatbotkit.com",
				Optional:            true,
			},
		},
	}
}

func (p *ChatBotKitProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ChatBotKitProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Initialize the ChatBotKit API client
	// apiKey := data.APIKey.ValueString()
	// if apiKey == "" {
	//     apiKey = os.Getenv("CHATBOTKIT_API_KEY")
	// }
	// baseURL := data.BaseURL.ValueString()
	// if baseURL == "" {
	//     baseURL = "https://api.chatbotkit.com"
	// }
	// client := &Client{
	//     APIKey:  apiKey,
	//     BaseURL: baseURL,
	// }

	client := &Client{}
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *ChatBotKitProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{

		NewBlueprintResource,
		NewBotResource,
		NewDatasetResource,
		NewDiscordIntegrationResource,
		NewEmailIntegrationResource,
		NewExtractIntegrationResource,
		NewFileResource,
		NewMcpserverIntegrationResource,
		NewMessengerIntegrationResource,
		NewNotionIntegrationResource,
		NewPortalResource,
		NewSecretResource,
		NewSitemapIntegrationResource,
		NewSkillsetAbilityResource,
		NewSkillsetResource,
		NewSlackIntegrationResource,
		NewTelegramIntegrationResource,
		NewTriggerIntegrationResource,
		NewTwilioIntegrationResource,
		NewWhatsAppIntegrationResource,
	}
}

func (p *ChatBotKitProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// TODO: Add data sources here
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ChatBotKitProvider{
			version: version,
		}
	}
}
