---
page_title: "chatbotkit_skillset_ability Resource - terraform-provider-chatbotkit"
subcategory: ""
description: |-
  Manages a ChatBotKit Skillset Ability resource.
---

# chatbotkit_skillset_ability (Resource)

Manages a ChatBotKit Skillset Ability. Abilities are individual tools or functions within a skillset that bots can use to perform specific actions. They can be custom functions, integrations with external APIs, or references to other bots.

## Example Usage

### Basic Ability

```terraform
resource "chatbotkit_skillset" "tools" {
  name        = "Support Tools"
  description = "Tools for support operations"
}

resource "chatbotkit_skillset_ability" "search" {
  skillset_id = chatbotkit_skillset.tools.id
  name        = "search_docs"
  description = "Search the documentation for relevant information"
  instruction = "Use this tool to search through documentation when the user asks about product features."
}
```

### Ability with Bot Reference

```terraform
resource "chatbotkit_bot" "specialist" {
  name        = "Technical Specialist"
  description = "Handles technical questions"
  backstory   = "You are a technical specialist."
}

resource "chatbotkit_skillset" "main_tools" {
  name        = "Main Tools"
  description = "Tools for the main bot"
}

resource "chatbotkit_skillset_ability" "ask_specialist" {
  skillset_id = chatbotkit_skillset.main_tools.id
  name        = "ask_technical_specialist"
  description = "Delegate technical questions to the specialist bot"
  bot_id      = chatbotkit_bot.specialist.id
}
```

### Ability with Secret for API Integration

```terraform
resource "chatbotkit_secret" "api_key" {
  name  = "External API Key"
  type  = "bearer"
  value = var.external_api_key
}

resource "chatbotkit_skillset" "integrations" {
  name        = "External Integrations"
  description = "Integration tools"
}

resource "chatbotkit_skillset_ability" "external_api" {
  skillset_id = chatbotkit_skillset.integrations.id
  name        = "query_external_api"
  description = "Query an external API for data"
  secret_id   = chatbotkit_secret.api_key.id
  instruction = "Use this to fetch data from the external service."
}
```

### Ability with File Reference

```terraform
resource "chatbotkit_file" "config" {
  name        = "Ability Configuration"
  description = "Configuration file for ability"
}

resource "chatbotkit_skillset" "configured_tools" {
  name        = "Configured Tools"
  description = "Tools with file configurations"
}

resource "chatbotkit_skillset_ability" "configured_action" {
  skillset_id = chatbotkit_skillset.configured_tools.id
  name        = "configured_action"
  description = "Action with file-based configuration"
  file_id     = chatbotkit_file.config.id
}
```

## Argument Reference

The following arguments are supported:

- `skillset_id` - (Required) The ID of the skillset to attach this ability to. Changing this forces a new resource to be created.
- `name` - (Optional) The name of the ability. This is used as the function name when the bot calls the tool.
- `description` - (Optional) A description of what the ability does. This helps the AI understand when to use the tool.
- `instruction` - (Optional) Detailed instructions for how and when to use this ability.
- `bot_id` - (Optional) The ID of a bot to delegate to when this ability is invoked.
- `secret_id` - (Optional, Sensitive) The ID of a secret to use for authentication when calling external APIs.
- `file_id` - (Optional) The ID of a file to use with this ability.
- `blueprint_id` - (Optional) The ID of a blueprint to associate with this ability.
- `space_id` - (Optional) The ID of a space to use with this ability.
- `meta` - (Optional) A map of metadata key-value pairs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier of the ability.
- `created_at` - The timestamp when the ability was created.
- `updated_at` - The timestamp when the ability was last updated.

## Import

Skillset abilities can be imported using their ID:

```bash
terraform import chatbotkit_skillset_ability.example ability_abc123def456
```
