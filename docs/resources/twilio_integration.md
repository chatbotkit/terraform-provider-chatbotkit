---
page_title: "chatbotkit_twilio_integration Resource - terraform-provider-chatbotkit"
subcategory: ""
description: |-
  Manages a ChatBotKit Twilio Integration resource.
---

# chatbotkit_twilio_integration (Resource)

Manages a ChatBotKit Twilio Integration. This integration allows you to connect your ChatBotKit bot to Twilio, enabling AI-powered SMS and voice conversations.

## Example Usage

### Basic Twilio Integration

```terraform
resource "chatbotkit_bot" "assistant" {
  name        = "SMS Assistant"
  description = "AI assistant for SMS"
  backstory   = "You are a helpful assistant responding to text messages."
}

resource "chatbotkit_twilio_integration" "example" {
  name        = "SMS Bot"
  description = "Connect bot to Twilio SMS"
  bot_id      = chatbotkit_bot.assistant.id
}
```

### Full Configuration

```terraform
resource "chatbotkit_twilio_integration" "advanced" {
  name        = "Advanced Twilio Integration"
  description = "Full-featured Twilio integration"
  bot_id      = chatbotkit_bot.assistant.id
  
  session_duration   = 3600000  # 1 hour in milliseconds
  contact_collection = true
}
```

### With Blueprint

```terraform
resource "chatbotkit_blueprint" "twilio_template" {
  name        = "Twilio Integration Template"
  description = "Template for Twilio integrations"
}

resource "chatbotkit_twilio_integration" "from_template" {
  name         = "Customer SMS"
  description  = "Created from template"
  blueprint_id = chatbotkit_blueprint.twilio_template.id
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
- `blueprint_id` - (Optional) The ID of a blueprint to associate with this integration.
- `meta` - (Optional) A map of metadata key-value pairs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier of the integration.
- `created_at` - The timestamp when the integration was created.
- `updated_at` - The timestamp when the integration was last updated.

## Import

Twilio integrations can be imported using their ID:

```bash
terraform import chatbotkit_twilio_integration.example twilio_abc123def456
```
