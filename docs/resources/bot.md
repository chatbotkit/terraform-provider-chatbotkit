---
page_title: "chatbotkit_bot Resource - terraform-provider-chatbotkit"
subcategory: ""
description: |-
  Manages a ChatBotKit Bot resource.
---

# chatbotkit_bot (Resource)

Manages a ChatBotKit Bot. Bots are AI agents that can interact with users through various integrations such as chat widgets, Slack, Discord, and more.

## Example Usage

### Basic Bot

```terraform
resource "chatbotkit_bot" "example" {
  name        = "My Assistant"
  description = "A helpful AI assistant"
  backstory   = "You are a helpful assistant that answers questions clearly and concisely."
}
```

### Bot with Dataset and Skillset

```terraform
resource "chatbotkit_dataset" "knowledge" {
  name        = "Knowledge Base"
  description = "Product documentation"
}

resource "chatbotkit_skillset" "tools" {
  name        = "Support Tools"
  description = "Tools for customer support"
}

resource "chatbotkit_bot" "support" {
  name        = "Customer Support Bot"
  description = "Handles customer inquiries"
  backstory   = "You are a helpful customer support agent."
  model       = "gpt-4"
  
  dataset_id  = chatbotkit_dataset.knowledge.id
  skillset_id = chatbotkit_skillset.tools.id
  
  moderation = true
  privacy    = true
}
```

### Bot with Blueprint

```terraform
resource "chatbotkit_blueprint" "template" {
  name        = "Support Template"
  description = "Template for support bots"
}

resource "chatbotkit_bot" "from_blueprint" {
  name         = "Support Bot Instance"
  description  = "Created from template"
  blueprint_id = chatbotkit_blueprint.template.id
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional) The name of the bot. This is displayed in the ChatBotKit dashboard.
- `description` - (Optional) A description of the bot's purpose.
- `backstory` - (Optional) The system prompt/backstory that defines the bot's personality and behavior.
- `model` - (Optional) The AI model to use (e.g., "gpt-4", "gpt-3.5-turbo").
- `dataset_id` - (Optional) The ID of a dataset to attach for knowledge retrieval.
- `skillset_id` - (Optional) The ID of a skillset to attach for tool usage.
- `blueprint_id` - (Optional) The ID of a blueprint this bot is based on.
- `moderation` - (Optional) Whether to enable content moderation. Defaults to `false`.
- `privacy` - (Optional) Whether to enable privacy mode. Defaults to `false`.
- `visibility` - (Optional) The visibility of the bot. Can be "private" or "public".
- `meta` - (Optional) A map of metadata key-value pairs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier of the bot.
- `created_at` - The timestamp when the bot was created.
- `updated_at` - The timestamp when the bot was last updated.

## Import

Bots can be imported using their ID:

```bash
terraform import chatbotkit_bot.example bot_abc123def456
```
