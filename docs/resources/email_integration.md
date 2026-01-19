---
page_title: "chatbotkit_email_integration Resource - terraform-provider-chatbotkit"
subcategory: ""
description: |-
  Manages a ChatBotKit Email Integration resource.
---

# chatbotkit_email_integration (Resource)

Manages a ChatBotKit Email Integration. This integration allows you to connect your ChatBotKit bot to email, enabling AI-powered email conversations and auto-responses.

## Example Usage

### Basic Email Integration

```terraform
resource "chatbotkit_bot" "assistant" {
  name        = "Email Assistant"
  description = "AI assistant for email"
  backstory   = "You are a helpful assistant responding to email inquiries."
}

resource "chatbotkit_email_integration" "example" {
  name        = "Email Bot"
  description = "Connect bot to email"
  bot_id      = chatbotkit_bot.assistant.id
}
```

### Full Configuration

```terraform
resource "chatbotkit_email_integration" "advanced" {
  name        = "Advanced Email Integration"
  description = "Full-featured email integration"
  bot_id      = chatbotkit_bot.assistant.id
  
  session_duration   = 86400000  # 24 hours in milliseconds
  contact_collection = true
  attachments        = true
}
```

### With Blueprint

```terraform
resource "chatbotkit_blueprint" "email_template" {
  name        = "Email Integration Template"
  description = "Template for email integrations"
}

resource "chatbotkit_email_integration" "from_template" {
  name         = "Support Email"
  description  = "Created from template"
  blueprint_id = chatbotkit_blueprint.email_template.id
  bot_id       = chatbotkit_bot.assistant.id
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional) The name of the integration. This is displayed in the ChatBotKit dashboard.
- `description` - (Optional) A description of the integration's purpose.
- `bot_id` - (Optional) The ID of the ChatBotKit bot to connect.
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

Email integrations can be imported using their ID:

```bash
terraform import chatbotkit_email_integration.example email_abc123def456
```
