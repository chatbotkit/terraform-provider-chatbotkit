# Terraform Provider for ChatBotKit

A modern Terraform provider for managing [ChatBotKit](https://chatbotkit.com) resources using the Terraform Plugin Framework.

## Features

- **Complete Resource Coverage**: Manage all ChatBotKit resources including bots, datasets, skillsets, files, integrations, and secrets
- **Data Sources**: Query existing resources with both single and list data sources
- **Modern Architecture**: Built with the latest Terraform Plugin Framework
- **Type Safety**: Strongly typed Go implementation synchronized with the ChatBotKit API
- **Security**: Proper handling of sensitive data (secrets, API tokens)

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.24 (for development)
- ChatBotKit API token

## Installation

### From Source

```bash
git clone https://github.com/chatbotkit/terraform-provider
cd terraform-provider
go build -o terraform-provider-chatbotkit
```

### Using Terraform

Add the following to your Terraform configuration:

```hcl
terraform {
  required_providers {
    chatbotkit = {
      source = "chatbotkit/chatbotkit"
    }
  }
}

provider "chatbotkit" {
  token = var.chatbotkit_token # or set CHATBOTKIT_TOKEN environment variable
}
```

## Usage

### Provider Configuration

```hcl
provider "chatbotkit" {
  # Configuration options
  token = "your-api-token" # Can also be set via CHATBOTKIT_TOKEN env var
}
```

### Resources

#### Bot Resource

```hcl
resource "chatbotkit_bot" "example" {
  name         = "My Chatbot"
  description  = "A helpful assistant"
  model        = "gpt-4"
  backstory    = "You are a helpful assistant."
  temperature  = 0.7
  moderation   = true
  privacy      = true
  dataset_id   = chatbotkit_dataset.example.id
  skillset_id  = chatbotkit_skillset.example.id
}
```

#### Dataset Resource

```hcl
resource "chatbotkit_dataset" "example" {
  name        = "My Dataset"
  description = "Knowledge base for my bot"
  type        = "text"
}
```

#### Skillset Resource

```hcl
resource "chatbotkit_skillset" "example" {
  name        = "My Skillset"
  description = "Custom abilities for my bot"
}
```

#### File Resource

```hcl
resource "chatbotkit_file" "example" {
  name   = "example.pdf"
  type   = "application/pdf"
  source = "https://example.com/document.pdf"
}
```

#### Integration Resource

```hcl
resource "chatbotkit_integration" "example" {
  name        = "Slack Integration"
  description = "Connect bot to Slack"
  type        = "slack"
  bot_id      = chatbotkit_bot.example.id
}
```

#### Secret Resource

```hcl
resource "chatbotkit_secret" "example" {
  name  = "api_key"
  value = var.secret_value
}
```

### Data Sources

#### Single Resource Data Sources

```hcl
data "chatbotkit_bot" "existing" {
  id = "bot-id-here"
}

data "chatbotkit_dataset" "existing" {
  id = "dataset-id-here"
}

data "chatbotkit_skillset" "existing" {
  id = "skillset-id-here"
}

data "chatbotkit_file" "existing" {
  id = "file-id-here"
}

data "chatbotkit_integration" "existing" {
  id = "integration-id-here"
}

data "chatbotkit_secret" "existing" {
  id = "secret-id-here"
}
```

#### List Data Sources

```hcl
data "chatbotkit_bots" "all" {}

data "chatbotkit_datasets" "all" {}

data "chatbotkit_skillsets" "all" {}

data "chatbotkit_files" "all" {}

data "chatbotkit_integrations" "all" {}

data "chatbotkit_secrets" "all" {}
```

## Development

### Building

```bash
go build -o terraform-provider-chatbotkit
```

### Testing

```bash
go test ./...
```

### Local Development

To use a locally built provider, create a `.terraformrc` file in your home directory:

```hcl
provider_installation {
  dev_overrides {
    "chatbotkit/chatbotkit" = "/path/to/your/provider/binary"
  }
  direct {}
}
```

## API Synchronization

This provider is designed to stay synchronized with the ChatBotKit API specification. The provider uses the official ChatBotKit API v1 endpoints:

- Base URL: `https://api.chatbotkit.com/v1`
- API Spec: `https://api.chatbotkit.com/v1/spec`

All resources follow the standard CRUD pattern:
- List: `GET /resource/list`
- Fetch: `GET /resource/{id}/fetch`
- Create: `POST /resource/create`
- Update: `POST /resource/{id}/update`
- Delete: `POST /resource/{id}/delete`

## Resources Excluded

As per the design requirements, the following resources are explicitly excluded:
- Contacts
- Conversations
- Tasks
- Memory
- Spaces
- Ratings

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

See [LICENSE](LICENSE) for more information.

## Support

- [ChatBotKit Documentation](https://chatbotkit.com/docs)
- [ChatBotKit API Reference](https://chatbotkit.com/docs/api)
- [Issue Tracker](https://github.com/chatbotkit/terraform-provider/issues)