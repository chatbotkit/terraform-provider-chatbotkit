---
page_title: "chatbotkit_discord_integration Resource - terraform-provider-chatbotkit"
subcategory: ""
description: |-
  Manages a ChatBotKit Discord Integration resource.
---

# chatbotkit_discord_integration (Resource)

Manages a ChatBotKit Discord Integration. This integration allows you to connect your ChatBotKit bot to Discord servers, enabling AI-powered conversations in Discord channels and direct messages.

## Example Usage

### Basic Discord Integration

```terraform
resource "chatbotkit_bot" "assistant" {
  name        = "Discord Assistant"
  description = "AI assistant for Discord"
  backstory   = "You are a helpful assistant for the Discord server."
}

resource "chatbotkit_discord_integration" "example" {
  name        = "Server Bot"
  description = "Connect bot to Discord server"
  bot_id      = chatbotkit_bot.assistant.id
  bot_token   = var.discord_bot_token
  app_id      = var.discord_app_id
  public_key  = var.discord_public_key
}
```

### Full Configuration

```terraform
resource "chatbotkit_discord_integration" "advanced" {
  name        = "Advanced Discord Integration"
  description = "Full-featured Discord integration"
  bot_id      = chatbotkit_bot.assistant.id
  
  bot_token   = var.discord_bot_token
  app_id      = var.discord_app_id
  public_key  = var.discord_public_key
  handle      = "mybot"
  
  session_duration   = 3600000  # 1 hour in milliseconds
  contact_collection = true
}
```

### With Blueprint

```terraform
resource "chatbotkit_blueprint" "discord_template" {
  name        = "Discord Integration Template"
  description = "Template for Discord integrations"
}

resource "chatbotkit_discord_integration" "from_template" {
  name         = "Community Bot"
  description  = "Created from template"
  blueprint_id = chatbotkit_blueprint.discord_template.id
  bot_id       = chatbotkit_bot.assistant.id
  bot_token    = var.discord_bot_token
  app_id       = var.discord_app_id
  public_key   = var.discord_public_key
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional) The name of the integration. This is displayed in the ChatBotKit dashboard.
- `description` - (Optional) A description of the integration's purpose.
- `bot_id` - (Optional) The ID of the ChatBotKit bot to connect.
- `bot_token` - (Optional, Sensitive) The Discord Bot Token for API access.
- `app_id` - (Optional) The Discord Application ID.
- `public_key` - (Optional) The Discord Public Key for request verification.
- `handle` - (Optional) The bot handle or username.
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

Discord integrations can be imported using their ID:

```bash
terraform import chatbotkit_discord_integration.example discord_abc123def456
```
