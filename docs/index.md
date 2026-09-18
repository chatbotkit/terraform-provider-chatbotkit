---
page_title: "Provider: ChatBotKit"
description: |-
  The ChatBotKit provider allows you to manage ChatBotKit AI chatbot resources using Terraform.
---

# ChatBotKit Provider

The ChatBotKit provider enables you to manage AI chatbot resources on the [ChatBotKit](https://chatbotkit.com) platform using Terraform. You can create and manage bots, datasets, skillsets, integrations, and more through infrastructure as code.

## Example Usage

```terraform
terraform {
  required_providers {
    chatbotkit = {
      source  = "chatbotkit/chatbotkit"
      version = "~> 1.0"
    }
  }
}

provider "chatbotkit" {
  api_token = var.chatbotkit_api_token
}

# Create a knowledge base dataset
resource "chatbotkit_dataset" "knowledge" {
  name        = "Product Knowledge Base"
  description = "Contains product documentation and FAQs"
}

# Create a skillset for tools
resource "chatbotkit_skillset" "tools" {
  name        = "Customer Support Tools"
  description = "Tools for customer support operations"
}

# Create an AI bot
resource "chatbotkit_bot" "assistant" {
  name        = "Customer Support Bot"
  description = "Handles customer inquiries"
  backstory   = "You are a helpful customer support agent for our company."
  model       = "gpt-4"
  
  dataset_id  = chatbotkit_dataset.knowledge.id
  skillset_id = chatbotkit_skillset.tools.id
}
```

## Authentication

The ChatBotKit provider requires a API token for authentication. You can obtain a API token from the [ChatBotKit Dashboard](https://chatbotkit.com).

### Configuration Options

You can configure authentication in two ways:

1. **Provider Configuration** (recommended for variables):
   ```terraform
   provider "chatbotkit" {
     api_token = var.chatbotkit_api_token
   }
   ```

2. **Environment Variable**:
   ```bash
   export CHATBOTKIT_API_TOKEN="your-api-token"
   ```

When both are set, the provider configuration takes precedence.

## Schema

### Optional

- `api_token` (String, Sensitive) - The API token for authenticating with the ChatBotKit API. Can also be set via the `CHATBOTKIT_API_TOKEN` environment variable, or `CBK_API_TOKEN` for short.
- `api_key` (String, Sensitive, Deprecated) - The former name of `api_token`. It is used when `api_token` is not set.
- `base_url` (String) - The GraphQL endpoint URL. Defaults to `https://api.chatbotkit.com/graphql`. For a self-hosted platform use its GraphQL endpoint, e.g. `http://localhost:3000/api/v1/graphql`, or set the platform origin in the `CHATBOTKIT_API_URL` environment variable. Plain `http` works for local use.
- `run_as` (String) - The ID of a child User to operate on behalf of. When set, requests include the `X-RunAs-UserId` header, so one `api_token` belonging to the parent User can manage many child Users. Configure one provider alias per child User. Can also be set via the `CHATBOTKIT_API_RUNAS_USERID` environment variable (or `CBK_API_RUNAS_USERID` for short); `CHATBOTKIT_RUN_AS` and `CBK_RUN_AS` are still read.

### Operating on child Users (multi-tenancy)

A API token belonging to a parent User, combined with `run_as`, lets one configuration manage many isolated child Users. This follows the standard Terraform multi-account pattern of provider aliases, similar to the AWS provider's `assume_role`:

```hcl
provider "chatbotkit" {
  alias  = "acme"
  run_as = var.acme_account_id # api_token from CHATBOTKIT_API_TOKEN
}

provider "chatbotkit" {
  alias  = "globex"
  run_as = var.globex_account_id
}
```

See the `multi-tenant-agents-shared` and `multi-tenant-agents-per-customer` examples.
