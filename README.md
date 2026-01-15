# ChatBotKit Terraform Provider

This directory contains the ChatBotKit Terraform Provider, which allows you to manage ChatBotKit resources using Terraform.

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

## Directory Structure

```
sdks/terraform/
├── main.go                    # Provider entry point
├── go.mod                     # Go module definition
├── go.sum                     # Go dependencies
├── types/
│   └── types.go              # Generated Go types
├── internal/
│   └── provider/
│       ├── client.go          # GraphQL API client
│       ├── provider.go        # Provider configuration
│       └── resource_*.go      # Resource implementations
└── examples/
    └── basic/
        └── main.tf           # Example Terraform configuration
```

## Regenerating Code

The provider code is generated from the GraphQL schema:

```bash
cd sites/main

# Generate Go types
pnpm script:sync-api-spec-to-terraform

# Generate provider resources and client
pnpm script:gen-terraform-stubs
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

| Data Source          | Description                         |
| -------------------- | ----------------------------------- |
| `chatbotkit_bot`     | Read information about an existing bot |
| `chatbotkit_dataset` | Read information about an existing dataset |
| `chatbotkit_blueprint` | Read information about an existing blueprint |
| `chatbotkit_skillset` | Read information about an existing skillset |

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
