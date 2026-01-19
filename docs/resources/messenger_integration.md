---
page_title: "chatbotkit_messenger_integration Resource - terraform-provider-chatbotkit"
subcategory: ""
description: |-
  Manages a ChatBotKit Messenger Integration resource.
---

# chatbotkit_messenger_integration (Resource)

Manages a ChatBotKit Messenger Integration. This integration allows you to connect your ChatBotKit bot to Facebook Messenger, enabling AI-powered conversations with users on Messenger.

## Example Usage

### Basic Messenger Integration

```terraform
resource "chatbotkit_bot" "assistant" {
  name        = "Messenger Assistant"
  description = "AI assistant for Messenger"
  backstory   = "You are a helpful assistant on Facebook Messenger."
}

resource "chatbotkit_messenger_integration" "example" {
  name         = "Facebook Page Bot"
  description  = "Connect bot to Facebook Messenger"
  bot_id       = chatbotkit_bot.assistant.id
  access_token = var.messenger_access_token
}
```

### Full Configuration

```terraform
resource "chatbotkit_messenger_integration" "advanced" {
  name        = "Advanced Messenger Integration"
  description = "Full-featured Messenger integration"
  bot_id      = chatbotkit_bot.assistant.id
  
  access_token     = var.messenger_access_token
  session_duration = 3600000  # 1 hour in milliseconds
  attachments      = true
}
```

### With Blueprint

```terraform
resource "chatbotkit_blueprint" "messenger_template" {
  name        = "Messenger Integration Template"
  description = "Template for Messenger integrations"
}

resource "chatbotkit_messenger_integration" "from_template" {
  name         = "Page Support"
  description  = "Created from template"
  blueprint_id = chatbotkit_blueprint.messenger_template.id
  bot_id       = chatbotkit_bot.assistant.id
  access_token = var.messenger_access_token
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional) The name of the integration. This is displayed in the ChatBotKit dashboard.
- `description` - (Optional) A description of the integration's purpose.
- `bot_id` - (Optional) The ID of the ChatBotKit bot to connect.
- `access_token` - (Optional, Sensitive) The Facebook Messenger page access token.
- `session_duration` - (Optional) The duration of a conversation session in milliseconds.
- `attachments` - (Optional) Whether to enable file attachments in conversations.
- `blueprint_id` - (Optional) The ID of a blueprint to associate with this integration.
- `meta` - (Optional) A map of metadata key-value pairs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier of the integration.
- `created_at` - The timestamp when the integration was created.
- `updated_at` - The timestamp when the integration was last updated.

## Import

Messenger integrations can be imported using their ID:

```bash
terraform import chatbotkit_messenger_integration.example messenger_abc123def456
```
