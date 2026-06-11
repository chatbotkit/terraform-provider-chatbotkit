[![ChatBotKit](https://img.shields.io/badge/credits-ChatBotKit-blue.svg)](https://chatbotkit.com)
[![CBK.AI](https://img.shields.io/badge/credits-CBK.AI-blue.svg)](https://cbk.ai)
[![Email](https://img.shields.io/badge/Email-Support-blue?logo=mail.ru)](mailto:support@chatbotkit.com)
[![Discord](https://img.shields.io/badge/Discord-Support-blue?logo=discord)](https://go.cbk.ai/discord)
[![Terraform Registry](https://img.shields.io/badge/Terraform-Registry-purple.svg)](https://registry.terraform.io/providers/chatbotkit/chatbotkit/latest)
[![Follow on Twitter](https://img.shields.io/twitter/follow/chatbotkit.svg?logo=twitter)](https://twitter.com/chatbotkit)

```text
 .d8888b.  888888b.   888    d8P
d88P  Y88b 888  "88b  888   d8P
888    888 888  .88P  888  d8P
888        8888888K.  888d88K
888        888  "Y88b 8888888b
888    888 888    888 888  Y88b
Y88b  d88P 888   d88P 888   Y88b
 "Y8888P"  8888888P"  888    Y88b .ai
```

# ChatBotKit Terraform Provider

The official Terraform provider for [ChatBotKit](https://chatbotkit.com) - a
platform for building and deploying conversational AI applications. Use it to
create and manage ChatBotKit bots, datasets, skillsets, integrations, and other
resources as infrastructure as code.

**Terraform Registry:** https://registry.terraform.io/providers/chatbotkit/chatbotkit/latest

## Why ChatBotKit?

**Build lighter, future-proof AI agents.** When you build with ChatBotKit, the heavy lifting happens on our servers, not in your application. This architectural advantage delivers:

- 🪶 **Lightweight Agents**: Your agents stay lean because complex AI processing, model orchestration, and tool execution happen server-side. Less code in your app means faster load times and simpler maintenance.

- 🛡️ **Robust & Streamlined**: Server-side processing provides a more reliable experience with built-in error handling, automatic retries, and consistent behavior across all platforms.

- 🔄 **Backward & Forward Compatible**: As AI technology evolves with new models, new capabilities, and new paradigms, your agents automatically benefit. No code changes required on your end.

- 🔮 **Future-Proof**: Agents you build today will remain capable tomorrow. When we add support for new AI models or capabilities, your existing agents gain those powers without any updates to your codebase.

This means you can focus on building great user experiences while ChatBotKit handles the complexity of the ever-changing AI landscape.

## Installation

```hcl
terraform {
  required_providers {
    chatbotkit = {
      source = "chatbotkit/chatbotkit"
    }
  }
}
```

## Building

```bash
cd sdks/terraform
go build -o terraform-provider-chatbotkit
```

## Development Testing

### 1. Set up Terraform Dev Override

Create or edit `~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "chatbotkit/chatbotkit" = "/path/to/cbk-platform/sdks/terraform"
  }
  direct {}
}
```

### 2. Set API Key

```bash
export CHATBOTKIT_API_KEY="your-api-key"
```

### 3. Test with Example Configuration

```bash
cd examples/basic
terraform init
terraform plan
terraform apply
```

## Running Tests

```bash
# Run unit tests
go test -v ./internal/provider/ -run "^Test[^Acc]"

# Run acceptance tests (requires CHATBOTKIT_API_KEY)
CHATBOTKIT_API_KEY=your-api-key go test -v ./internal/provider/ -run "^TestAcc"
```

## Directory Structure

```
terraform-provider-chatbotkit/
├── main.go                          # Provider entry point
├── go.mod                           # Go module definition
├── go.sum                           # Go dependencies
├── .goreleaser.yml                  # Release configuration
├── terraform-registry-manifest.json # Registry manifest
├── docs/                            # Terraform Registry documentation
│   ├── index.md                     # Provider documentation
│   ├── resources/                   # Resource documentation
│   └── data-sources/                # Data source documentation
├── types/
│   └── types.go                     # Generated Go types
├── internal/
│   └── provider/
│       ├── client.go                # GraphQL API client
│       ├── client_test.go           # Client unit tests
│       ├── provider.go              # Provider configuration
│       ├── provider_test.go         # Provider tests
│       ├── resource_*.go            # Resource implementations
│       └── resource_*_test.go       # Resource tests
└── examples/
    └── basic/
        └── main.tf                  # Example Terraform configuration
```

## Resources

The provider supports the following resources:

| Resource                           | Description                    |
| ---------------------------------- | ------------------------------ |
| `chatbotkit_bot`                   | Manages a ChatBotKit bot       |
| `chatbotkit_dataset`               | Manages a dataset              |
| `chatbotkit_blueprint`             | Manages a blueprint            |
| `chatbotkit_skillset`              | Manages a skillset             |
| `chatbotkit_skillset_ability`      | Manages a skillset ability     |
| `chatbotkit_secret`                | Manages a secret               |
| `chatbotkit_file`                  | Manages a file                 |
| `chatbotkit_portal`                | Manages a portal               |
| `chatbotkit_discord_integration`   | Manages Discord integration    |
| `chatbotkit_email_integration`     | Manages Email integration      |
| `chatbotkit_extract_integration`   | Manages Extract integration    |
| `chatbotkit_mcpserver_integration` | Manages MCP Server integration |
| `chatbotkit_messenger_integration` | Manages Messenger integration  |
| `chatbotkit_notion_integration`    | Manages Notion integration     |
| `chatbotkit_sitemap_integration`   | Manages Sitemap integration    |
| `chatbotkit_slack_integration`     | Manages Slack integration      |
| `chatbotkit_telegram_integration`  | Manages Telegram integration   |
| `chatbotkit_trigger_integration`   | Manages Trigger integration    |
| `chatbotkit_twilio_integration`    | Manages Twilio integration     |
| `chatbotkit_whatsapp_integration`  | Manages WhatsApp integration   |

## Data Sources

The provider supports the following data sources for reading existing resources:

| Data Source            | Description                                  |
| ---------------------- | -------------------------------------------- |
| `chatbotkit_bot`       | Read information about an existing bot       |
| `chatbotkit_dataset`   | Read information about an existing dataset   |
| `chatbotkit_blueprint` | Read information about an existing blueprint |
| `chatbotkit_skillset`  | Read information about an existing skillset  |

## Example Usage

```hcl
terraform {
  required_providers {
    chatbotkit = {
      source = "chatbotkit/chatbotkit"
    }
  }
}

provider "chatbotkit" {
  # api_key = "..." # Or set CHATBOTKIT_API_KEY env var
}

# Create a new bot
resource "chatbotkit_bot" "assistant" {
  name        = "Customer Support Bot"
  description = "Handles customer inquiries"
  backstory   = "You are a helpful customer support agent..."
  model       = "gpt-4"
}

# Create a dataset
resource "chatbotkit_dataset" "knowledge" {
  name        = "Product Knowledge Base"
  description = "Contains product documentation"
}

# Reference an existing bot by ID
data "chatbotkit_bot" "existing" {
  id = "bot_abc123"
}

output "existing_bot_name" {
  value = data.chatbotkit_bot.existing.name
}
```
