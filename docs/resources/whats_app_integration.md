---
page_title: "chatbotkit_whats_app_integration Resource - terraform-provider-chatbotkit"
subcategory: ""
description: |-
  Manages a ChatBotKit WhatsApp Integration resource.
---

# chatbotkit_whats_app_integration (Resource)

Manages a ChatBotKit WhatsApp Integration. This integration allows you to connect your ChatBotKit bot to WhatsApp Business, enabling AI-powered conversations with customers on WhatsApp.

## Example Usage

### Basic WhatsApp Integration

```terraform
resource "chatbotkit_bot" "assistant" {
  name        = "WhatsApp Assistant"
  description = "AI assistant for WhatsApp"
  backstory   = "You are a helpful customer service assistant on WhatsApp."
}

resource "chatbotkit_whats_app_integration" "example" {
  name            = "WhatsApp Business Bot"
  description     = "Connect bot to WhatsApp Business"
  bot_id          = chatbotkit_bot.assistant.id
  access_token    = var.whatsapp_access_token
  phone_number_id = var.whatsapp_phone_number_id
}
```

### Full Configuration

```terraform
resource "chatbotkit_whats_app_integration" "advanced" {
  name        = "Advanced WhatsApp Integration"
  description = "Full-featured WhatsApp integration"
  bot_id      = chatbotkit_bot.assistant.id
  
  access_token       = var.whatsapp_access_token
  phone_number_id    = var.whatsapp_phone_number_id
  session_duration   = 3600000  # 1 hour in milliseconds
  contact_collection = true
  attachments        = true
}
```

### With Blueprint

```terraform
resource "chatbotkit_blueprint" "whatsapp_template" {
  name        = "WhatsApp Integration Template"
  description = "Template for WhatsApp integrations"
}

resource "chatbotkit_whats_app_integration" "from_template" {
  name            = "Customer Support"
  description     = "Created from template"
  blueprint_id    = chatbotkit_blueprint.whatsapp_template.id
  bot_id          = chatbotkit_bot.assistant.id
  access_token    = var.whatsapp_access_token
  phone_number_id = var.whatsapp_phone_number_id
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional) The name of the integration. This is displayed in the ChatBotKit dashboard.
- `description` - (Optional) A description of the integration's purpose.
- `bot_id` - (Optional) The ID of the ChatBotKit bot to connect.
- `access_token` - (Optional, Sensitive) The WhatsApp Business API access token.
- `phone_number_id` - (Optional) The WhatsApp Business phone number ID.
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

WhatsApp integrations can be imported using their ID:

```bash
terraform import chatbotkit_whats_app_integration.example whatsapp_abc123def456
```
