---
page_title: "chatbotkit_trigger_integration Resource - terraform-provider-chatbotkit"
subcategory: ""
description: |-
  Manages a ChatBotKit Trigger Integration resource.
---

# chatbotkit_trigger_integration (Resource)

Manages a ChatBotKit Trigger Integration. This integration allows you to create scheduled or webhook-triggered bot interactions, enabling proactive bot communications and automated workflows.

## Example Usage

### Basic Trigger Integration

```terraform
resource "chatbotkit_bot" "assistant" {
  name        = "Scheduled Bot"
  description = "Bot for scheduled tasks"
  backstory   = "You are an assistant that handles scheduled tasks."
}

resource "chatbotkit_trigger_integration" "example" {
  name        = "Daily Report"
  description = "Trigger daily report generation"
  bot_id      = chatbotkit_bot.assistant.id
}
```

### Scheduled Trigger

```terraform
resource "chatbotkit_trigger_integration" "scheduled" {
  name        = "Hourly Check"
  description = "Run bot check every hour"
  bot_id      = chatbotkit_bot.assistant.id
  
  trigger_schedule = "0 * * * *"  # Every hour
  session_duration = 300000       # 5 minutes
}
```

### Authenticated Webhook Trigger

```terraform
resource "chatbotkit_trigger_integration" "webhook" {
  name         = "Webhook Trigger"
  description  = "Trigger via authenticated webhook"
  bot_id       = chatbotkit_bot.assistant.id
  authenticate = true
}
```

### With Blueprint

```terraform
resource "chatbotkit_blueprint" "trigger_template" {
  name        = "Trigger Integration Template"
  description = "Template for trigger integrations"
}

resource "chatbotkit_trigger_integration" "from_template" {
  name             = "Alert Trigger"
  description      = "Created from template"
  blueprint_id     = chatbotkit_blueprint.trigger_template.id
  bot_id           = chatbotkit_bot.assistant.id
  trigger_schedule = "*/30 * * * *"  # Every 30 minutes
}
```

### Full Configuration

```terraform
resource "chatbotkit_trigger_integration" "full" {
  name        = "Full Trigger Configuration"
  description = "Complete trigger setup"
  bot_id      = chatbotkit_bot.assistant.id
  
  trigger_schedule = "0 9 * * 1-5"  # 9 AM on weekdays
  session_duration = 600000         # 10 minutes
  authenticate     = true
  
  meta = {
    type     = "scheduled"
    priority = "high"
  }
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional) The name of the integration. This is displayed in the ChatBotKit dashboard.
- `description` - (Optional) A description of the integration's purpose.
- `bot_id` - (Optional) The ID of the ChatBotKit bot to connect.
- `trigger_schedule` - (Optional) A cron expression for scheduled trigger execution.
- `session_duration` - (Optional) The duration of a trigger session in milliseconds.
- `authenticate` - (Optional) Whether to require authentication for webhook triggers.
- `blueprint_id` - (Optional) The ID of a blueprint to associate with this integration.
- `meta` - (Optional) A map of metadata key-value pairs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier of the integration.
- `created_at` - The timestamp when the integration was created.
- `updated_at` - The timestamp when the integration was last updated.

## Import

Trigger integrations can be imported using their ID:

```bash
terraform import chatbotkit_trigger_integration.example trigger_abc123def456
```
