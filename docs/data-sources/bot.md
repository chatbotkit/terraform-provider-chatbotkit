---
page_title: "chatbotkit_bot Data Source - terraform-provider-chatbotkit"
subcategory: ""
description: |-
  Use this data source to read information about an existing ChatBotKit Bot.
---

# chatbotkit_bot (Data Source)

Use this data source to read information about an existing ChatBotKit Bot. This is useful when you need to reference a bot that was created outside of Terraform or in a different Terraform configuration.

## Example Usage

### Read an Existing Bot

```terraform
data "chatbotkit_bot" "existing" {
  id = "bot_abc123def456"
}

output "bot_name" {
  value = data.chatbotkit_bot.existing.name
}

output "bot_model" {
  value = data.chatbotkit_bot.existing.model
}
```

### Use with Other Resources

```terraform
# Reference an existing bot for an integration
data "chatbotkit_bot" "main" {
  id = var.main_bot_id
}

resource "chatbotkit_slack_integration" "workspace" {
  name   = "Main Slack Integration"
  bot_id = data.chatbotkit_bot.main.id
}
```

## Argument Reference

The following arguments are required:

- `id` - (Required) The unique identifier of the bot to read.

## Attribute Reference

The following attributes are exported:

- `id` - The unique identifier of the bot.
- `name` - The name of the bot.
- `description` - The description of the bot.
- `backstory` - The system prompt/backstory of the bot.
- `model` - The AI model used by the bot.
- `dataset_id` - The ID of the attached dataset, if any.
- `skillset_id` - The ID of the attached skillset, if any.
- `blueprint_id` - The ID of the blueprint this bot is based on, if any.
- `moderation` - Whether content moderation is enabled.
- `privacy` - Whether privacy mode is enabled.
- `visibility` - The visibility setting of the bot.
- `meta` - A map of metadata key-value pairs.
- `created_at` - The timestamp when the bot was created.
- `updated_at` - The timestamp when the bot was last updated.
