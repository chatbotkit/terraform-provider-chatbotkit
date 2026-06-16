---
page_title: "chatbotkit_googlechat_integration Resource - terraform-provider-chatbotkit"
subcategory: ""
description: |-
  Manages a ChatBotKit Google Chat Integration resource.
---

# chatbotkit_googlechat_integration (Resource)

Manages a ChatBotKit Google Chat Integration. This integration allows you to connect your ChatBotKit bot to Google Chat, enabling AI-powered conversations in Google Chat spaces and direct messages.

## Example Usage

### Basic Google Chat Integration

```terraform
resource "chatbotkit_bot" "assistant" {
  name        = "Google Chat Assistant"
  description = "AI assistant for Google Chat"
  backstory   = "You are a helpful assistant for the team."
}

resource "chatbotkit_googlechat_integration" "example" {
  name                = "Workspace Integration"
  description         = "Connect bot to Google Chat"
  bot_id              = chatbotkit_bot.assistant.id
  project_number      = var.google_project_number
  service_account_key = var.google_service_account_key
}
```

### Full Configuration

```terraform
resource "chatbotkit_googlechat_integration" "advanced" {
  name        = "Advanced Google Chat Integration"
  description = "Full-featured Google Chat integration"
  bot_id      = chatbotkit_bot.assistant.id

  project_number      = var.google_project_number
  service_account_key = var.google_service_account_key

  session_duration   = 3600000 # 1 hour in milliseconds
  contact_collection = true
  auto_respond       = "always"
  allow_from         = "example.com"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional) The name of the integration. This is displayed in the ChatBotKit dashboard.
- `description` - (Optional) A description of the integration's purpose.
- `bot_id` - (Optional) The ID of the ChatBotKit bot to connect.
- `project_number` - (Optional) The Google Cloud project number used to verify incoming event JWT audience claims.
- `service_account_key` - (Optional, Sensitive) The Google service account JSON key for sending messages via the Chat REST API.
- `allow_from` - (Optional) The allowed senders for this integration.
- `auto_respond` - (Optional) Auto-respond configuration for the integration.
- `session_duration` - (Optional) The duration of a conversation session in milliseconds.
- `contact_collection` - (Optional) Whether to collect contact information from users.
- `blueprint_id` - (Optional) The ID of a blueprint to associate with this integration.
- `meta` - (Optional) A map of metadata key-value pairs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier of the integration.
- `created_at` - The timestamp when the integration was created.
- `updated_at` - The timestamp when the integration was last updated.

## Import

Google Chat integrations can be imported using their ID:

```bash
terraform import chatbotkit_googlechat_integration.example googlechat_abc123def456
```
