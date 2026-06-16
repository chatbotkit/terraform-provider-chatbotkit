---
page_title: "chatbotkit_microsoftteams_integration Resource - terraform-provider-chatbotkit"
subcategory: ""
description: |-
  Manages a ChatBotKit Microsoft Teams Integration resource.
---

# chatbotkit_microsoftteams_integration (Resource)

Manages a ChatBotKit Microsoft Teams Integration. This integration allows you to connect your ChatBotKit bot to Microsoft Teams, enabling AI-powered conversations in Teams channels and direct messages.

## Example Usage

### Basic Microsoft Teams Integration

```terraform
resource "chatbotkit_bot" "assistant" {
  name        = "Teams Assistant"
  description = "AI assistant for Microsoft Teams"
  backstory   = "You are a helpful assistant for the team."
}

resource "chatbotkit_microsoftteams_integration" "example" {
  name                     = "Teams Integration"
  description              = "Connect bot to Microsoft Teams"
  bot_id                   = chatbotkit_bot.assistant.id
  bot_framework_app_id     = var.teams_app_id
  bot_framework_app_secret = var.teams_app_secret
  tenant_id                = var.azure_tenant_id
}
```

### Full Configuration

```terraform
resource "chatbotkit_microsoftteams_integration" "advanced" {
  name        = "Advanced Teams Integration"
  description = "Full-featured Microsoft Teams integration"
  bot_id      = chatbotkit_bot.assistant.id

  bot_framework_app_id     = var.teams_app_id
  bot_framework_app_secret = var.teams_app_secret
  tenant_id                = var.azure_tenant_id

  session_duration   = 3600000 # 1 hour in milliseconds
  contact_collection = true
  allow_from         = "example.com"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional) The name of the integration. This is displayed in the ChatBotKit dashboard.
- `description` - (Optional) A description of the integration's purpose.
- `bot_id` - (Optional) The ID of the ChatBotKit bot to connect.
- `bot_framework_app_id` - (Optional) The Microsoft Bot Framework application ID.
- `bot_framework_app_secret` - (Optional, Sensitive) The Microsoft Bot Framework application secret.
- `tenant_id` - (Optional) The Azure AD tenant ID.
- `allow_from` - (Optional) The allowed senders for this integration.
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

Microsoft Teams integrations can be imported using their ID:

```bash
terraform import chatbotkit_microsoftteams_integration.example microsoftteams_abc123def456
```
