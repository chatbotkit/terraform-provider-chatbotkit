---
page_title: "chatbotkit_telegram_integration Resource - terraform-provider-chatbotkit"
subcategory: ""
description: |-
  Manages a ChatBotKit Telegram Integration resource.
---

# chatbotkit_telegram_integration (Resource)

Manages a ChatBotKit Telegram Integration. This integration allows you to connect your ChatBotKit bot to Telegram, enabling AI-powered conversations in Telegram chats and groups.

## Example Usage

### Basic Telegram Integration

```terraform
resource "chatbotkit_bot" "assistant" {
  name        = "Telegram Assistant"
  description = "AI assistant for Telegram"
  backstory   = "You are a helpful assistant on Telegram."
}

resource "chatbotkit_telegram_integration" "example" {
  name        = "Telegram Bot"
  description = "Connect bot to Telegram"
  bot_id      = chatbotkit_bot.assistant.id
  bot_token   = var.telegram_bot_token
}
```

### Full Configuration

```terraform
resource "chatbotkit_telegram_integration" "advanced" {
  name        = "Advanced Telegram Integration"
  description = "Full-featured Telegram integration"
  bot_id      = chatbotkit_bot.assistant.id
  
  bot_token          = var.telegram_bot_token
  session_duration   = 3600000  # 1 hour in milliseconds
  contact_collection = true
  attachments        = true
}
```

### With Blueprint

```terraform
resource "chatbotkit_blueprint" "telegram_template" {
  name        = "Telegram Integration Template"
  description = "Template for Telegram integrations"
}

resource "chatbotkit_telegram_integration" "from_template" {
  name         = "Support Bot"
  description  = "Created from template"
  blueprint_id = chatbotkit_blueprint.telegram_template.id
  bot_id       = chatbotkit_bot.assistant.id
  bot_token    = var.telegram_bot_token
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional) The name of the integration. This is displayed in the ChatBotKit dashboard.
- `description` - (Optional) A description of the integration's purpose.
- `bot_id` - (Optional) The ID of the ChatBotKit bot to connect.
- `bot_token` - (Optional, Sensitive) The Telegram Bot Token from @BotFather.
- `session_duration` - (Optional) The duration of a conversation session in milliseconds.
- `contact_collection` - (Optional) Whether to collect contact information from users.
- `attachments` - (Optional) Whether to enable file attachments in conversations.
- `blueprint_id` - (Optional) The ID of a blueprint to associate with this integration.
- `meta` - (Optional) A map of metadata key-value pairs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier of the integration.
- `created_at` - The timestamp when the integration was created.
- `updated_at` - The timestamp when the integration was last updated.

## Import

Telegram integrations can be imported using their ID:

```bash
terraform import chatbotkit_telegram_integration.example telegram_abc123def456
```
