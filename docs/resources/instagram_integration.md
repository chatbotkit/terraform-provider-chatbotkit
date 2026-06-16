---
page_title: "chatbotkit_instagram_integration Resource - terraform-provider-chatbotkit"
subcategory: ""
description: |-
  Manages a ChatBotKit Instagram Integration resource.
---

# chatbotkit_instagram_integration (Resource)

Manages a ChatBotKit Instagram Integration. This integration allows you to connect your ChatBotKit bot to Instagram, enabling AI-powered conversations through Instagram direct messages.

## Example Usage

### Basic Instagram Integration

```terraform
resource "chatbotkit_bot" "assistant" {
  name        = "Instagram Assistant"
  description = "AI assistant for Instagram"
  backstory   = "You are a helpful assistant for our followers."
}

resource "chatbotkit_instagram_integration" "example" {
  name         = "Instagram DMs"
  description  = "Connect bot to Instagram direct messages"
  bot_id       = chatbotkit_bot.assistant.id
  access_token = var.instagram_access_token
}
```

### Full Configuration

```terraform
resource "chatbotkit_instagram_integration" "advanced" {
  name        = "Advanced Instagram Integration"
  description = "Full-featured Instagram integration"
  bot_id      = chatbotkit_bot.assistant.id

  access_token = var.instagram_access_token

  session_duration   = 3600000 # 1 hour in milliseconds
  contact_collection = true
  attachments        = true
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional) The name of the integration. This is displayed in the ChatBotKit dashboard.
- `description` - (Optional) A description of the integration's purpose.
- `bot_id` - (Optional) The ID of the ChatBotKit bot to connect.
- `access_token` - (Optional, Sensitive) The Instagram access token used to send and receive messages.
- `attachments` - (Optional) Whether to enable file attachments.
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

Instagram integrations can be imported using their ID:

```bash
terraform import chatbotkit_instagram_integration.example instagram_abc123def456
```
