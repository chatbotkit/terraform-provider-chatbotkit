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
	APIToken types.String `tfsdk:"api_token"`
	APIKey   types.String `tfsdk:"api_key"`
	BaseURL  types.String `tfsdk:"base_url"`
	RunAs    types.String `tfsdk:"run_as"`
}

func (p *ChatBotKitProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "chatbotkit"
	resp.Version = p.version
}

func (p *ChatBotKitProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "ChatBotKit Terraform Provider for managing AI chatbot resources.",
		Attributes: map[string]schema.Attribute{
			"api_token": schema.StringAttribute{
				MarkdownDescription: "The API token for authenticating with the ChatBotKit API. Can also be set via the CHATBOTKIT_API_TOKEN environment variable, or CBK_API_TOKEN for short.",
				Optional:            true,
				Sensitive:           true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "The former name of `api_token`. It is used when `api_token` is not set.",
				DeprecationMessage:  "Use api_token instead. api_key is still read when api_token is not set.",
				Optional:            true,
				Sensitive:           true,
			},
			"base_url": schema.StringAttribute{
				MarkdownDescription: "The GraphQL endpoint URL. Defaults to https://api.chatbotkit.com/graphql. For a self-hosted platform use its GraphQL endpoint, e.g. http://localhost:3000/api/v1/graphql, or set the platform origin in the CHATBOTKIT_API_URL environment variable. Plain http works for local use.",
				Optional:            true,
			},
			"run_as": schema.StringAttribute{
				MarkdownDescription: "The ID of a child User to operate on behalf of. When set, requests include the X-RunAs-UserId header, so one api_token belonging to the parent User can manage many child Users by configuring one provider alias per child User. Can also be set via the CHATBOTKIT_API_RUNAS_USERID environment variable (or CBK_API_RUNAS_USERID for short); CHATBOTKIT_RUN_AS and CBK_RUN_AS are still read.",
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

	// Get the API token from the config, then its deprecated name, then the
	// environment
	apiKey := data.APIToken.ValueString()
	if apiKey == "" {
		apiKey = data.APIKey.ValueString()
	}
	if apiKey == "" {
		apiKey = tokenFromEnv()
	}

	if apiKey == "" {
		resp.Diagnostics.AddError(
			"Missing API Token",
			"The API token is required. Set api_token in the provider configuration or the CHATBOTKIT_API_TOKEN environment variable.",
		)
		return
	}

	// Get base URL from config or use default
	baseURL := data.BaseURL.ValueString()

	// Create the API client
	client := NewClient(apiKey, baseURL)

	// Optionally operate on behalf of a child User.
	runAs := data.RunAs.ValueString()
	if runAs == "" {
		runAs = runAsFromEnv()
	}
	client.RunAs = runAs

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *ChatBotKitProvider) Resources(ctx context.Context) []func() resource.Resource {
	return append([]func() resource.Resource{

		NewBlueprintResource,
		NewBotResource,
		NewContextResource,
		NewDatasetResource,
		NewDiscordIntegrationResource,
		NewEmailIntegrationResource,
		NewExtractIntegrationResource,
		NewFileResource,
		NewGooglechatIntegrationResource,
		NewInstagramIntegrationResource,
		NewMcpserverIntegrationResource,
		NewMessengerIntegrationResource,
		NewMicrosoftteamsIntegrationResource,
		NewNotionIntegrationResource,
		NewPolicyResource,
		NewPortalResource,
		NewSecretResource,
		NewSitemapIntegrationResource,
		NewSkillserverIntegrationResource,
		NewSkillsetAbilityResource,
		NewSkillsetResource,
		NewSlackIntegrationResource,
		NewSpaceResource,
		NewSpaceSiteResource,
		NewSupportIntegrationResource,
		NewTaskResource,
		NewTelegramIntegrationResource,
		NewTriggerIntegrationResource,
		NewTwilioIntegrationResource,
		NewWhatsAppIntegrationResource,
		NewWidgetIntegrationResource,
	}, manualResources()...)
}

func (p *ChatBotKitProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{

		NewBlueprintDataSource,
		NewBotDataSource,
		NewDatasetDataSource,
		NewSkillsetDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ChatBotKitProvider{
			version: version,
		}
	}
}
