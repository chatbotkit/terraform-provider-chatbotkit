---
page_title: "chatbotkit_slack_integration Resource - terraform-provider-chatbotkit"
subcategory: ""
description: |-
  Manages a ChatBotKit Slack Integration resource.
---

# chatbotkit_slack_integration (Resource)

Manages a ChatBotKit Slack Integration. This integration allows you to connect your ChatBotKit bot to Slack workspaces, enabling AI-powered conversations in Slack channels and direct messages.

## Example Usage

### Basic Slack Integration

```terraform
resource "chatbotkit_bot" "assistant" {
  name        = "Slack Assistant"
  description = "AI assistant for Slack"
  backstory   = "You are a helpful assistant for the team."
}

resource "chatbotkit_slack_integration" "example" {
  name        = "Workspace Integration"
  description = "Connect bot to Slack workspace"
  bot_id      = chatbotkit_bot.assistant.id
  bot_token   = var.slack_bot_token
  signing_secret = var.slack_signing_secret
}
```

### Full Configuration

```terraform
resource "chatbotkit_slack_integration" "advanced" {
  name        = "Advanced Slack Integration"
  description = "Full-featured Slack integration"
  bot_id      = chatbotkit_bot.assistant.id
  
  bot_token      = var.slack_bot_token
  signing_secret = var.slack_signing_secret
  user_token     = var.slack_user_token
  
  session_duration   = 3600000  # 1 hour in milliseconds
  visible_messages   = 10
  contact_collection = true
  ratings            = true
  references         = true
  auto_respond       = "always"
}
```

### With Blueprint

```terraform
resource "chatbotkit_blueprint" "slack_template" {
  name        = "Slack Integration Template"
  description = "Template for Slack integrations"
}

resource "chatbotkit_slack_integration" "from_template" {
  name         = "Team Slack Bot"
  description  = "Created from template"
  blueprint_id = chatbotkit_blueprint.slack_template.id
  bot_id       = chatbotkit_bot.assistant.id
  bot_token    = var.slack_bot_token
  signing_secret = var.slack_signing_secret
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional) The name of the integration. This is displayed in the ChatBotKit dashboard.
- `description` - (Optional) A description of the integration's purpose.
- `bot_id` - (Optional) The ID of the ChatBotKit bot to connect.
- `bot_token` - (Optional, Sensitive) The Slack Bot Token (xoxb-...) for API access.
- `signing_secret` - (Optional, Sensitive) The Slack Signing Secret for verifying requests.
- `user_token` - (Optional) The Slack User Token for additional permissions.
- `auto_respond` - (Optional) Auto-respond configuration for the integration.
- `session_duration` - (Optional) The duration of a conversation session in milliseconds.
- `visible_messages` - (Optional) Number of previous messages visible in the conversation context.
- `contact_collection` - (Optional) Whether to collect contact information from users.
- `ratings` - (Optional) Whether to enable message ratings.
- `references` - (Optional) Whether to include message references in responses.
- `blueprint_id` - (Optional) The ID of a blueprint to associate with this integration.
- `meta` - (Optional) A map of metadata key-value pairs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier of the integration.
- `created_at` - The timestamp when the integration was created.
- `updated_at` - The timestamp when the integration was last updated.

## Import

Slack integrations can be imported using their ID:

```bash
terraform import chatbotkit_slack_integration.example slack_abc123def456
```
