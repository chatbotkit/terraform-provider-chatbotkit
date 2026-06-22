---
page_title: "chatbotkit_skillserver_integration Resource - terraform-provider-chatbotkit"
subcategory: ""
description: |-
  Manages a ChatBotKit Skill Server Integration resource.
---

# chatbotkit_skillserver_integration (Resource)

Manages a ChatBotKit Skill Server Integration. This integration exposes a skillset's abilities as a text-first HTTP API instead of over MCP. A `GET` to the endpoint returns a human- and agent-readable manual describing the available abilities and how to call them; a `POST` invokes an ability. The endpoint is authenticated by a single static access token (like a trigger), so an agent that holds the token can discover and call the abilities directly.

## Example Usage

### Basic Skill Server Integration

```terraform
resource "chatbotkit_skillset" "tools" {
  name        = "Shared Tools"
  description = "Tools to expose via the skill server"
}

resource "chatbotkit_skillserver_integration" "example" {
  name        = "Skill Server"
  description = "Expose tools as a text-first HTTP API"
  skillset_id = chatbotkit_skillset.tools.id
}
```

### With Abilities

```terraform
resource "chatbotkit_skillset" "api_tools" {
  name        = "API Tools"
  description = "API integration tools"
}

resource "chatbotkit_skillset_ability" "search" {
  skillset_id = chatbotkit_skillset.api_tools.id
  name        = "search"
  description = "Search for information"
  instruction = "Use this to search for relevant data"
}

resource "chatbotkit_skillserver_integration" "api_server" {
  name        = "API Skill Server"
  description = "Expose API tools over HTTP"
  skillset_id = chatbotkit_skillset.api_tools.id
}
```

### With Blueprint

```terraform
resource "chatbotkit_blueprint" "skill_template" {
  name        = "Skill Server Template"
  description = "Template for skill servers"
}

resource "chatbotkit_skillserver_integration" "from_template" {
  name         = "Custom Skill Server"
  description  = "Created from template"
  blueprint_id = chatbotkit_blueprint.skill_template.id
  skillset_id  = chatbotkit_skillset.tools.id
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional) The name of the integration. This is displayed in the ChatBotKit dashboard.
- `description` - (Optional) A description of the integration's purpose.
- `skillset_id` - (Optional) The ID of the skillset whose abilities are exposed.
- `blueprint_id` - (Optional) The ID of a blueprint to associate with this integration.
- `alias` - (Optional) A unique alias for the integration.
- `meta` - (Optional) A map of metadata key-value pairs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier of the integration.
- `created_at` - The timestamp when the integration was created.
- `updated_at` - The timestamp when the integration was last updated.

## Import

Skill Server integrations can be imported using their ID:

```bash
terraform import chatbotkit_skillserver_integration.example skillserver_abc123def456
```
