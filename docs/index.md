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
  api_key = var.chatbotkit_api_key
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

The ChatBotKit provider requires an API key for authentication. You can obtain an API key from the [ChatBotKit Dashboard](https://chatbotkit.com).

### Configuration Options

You can configure authentication in two ways:

1. **Provider Configuration** (recommended for variables):
   ```terraform
   provider "chatbotkit" {
     api_key = var.chatbotkit_api_key
   }
   ```

2. **Environment Variable**:
   ```bash
   export CHATBOTKIT_API_KEY="your-api-key"
   ```

When both are set, the provider configuration takes precedence.

## Schema

### Optional

- `api_key` (String, Sensitive) - The API key for authenticating with the ChatBotKit API. Can also be set via the `CHATBOTKIT_API_KEY` environment variable.
- `base_url` (String) - Custom API endpoint URL. Defaults to `https://api.chatbotkit.com/graphql`. This is typically only needed for testing or enterprise deployments.
