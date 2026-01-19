---
page_title: "chatbotkit_mcpserver_integration Resource - terraform-provider-chatbotkit"
subcategory: ""
description: |-
  Manages a ChatBotKit MCP Server Integration resource.
---

# chatbotkit_mcpserver_integration (Resource)

Manages a ChatBotKit MCP Server Integration. This integration exposes a skillset as an MCP (Model Context Protocol) server, allowing external AI applications to use your ChatBotKit abilities as tools.

## Example Usage

### Basic MCP Server Integration

```terraform
resource "chatbotkit_skillset" "tools" {
  name        = "Shared Tools"
  description = "Tools to expose via MCP"
}

resource "chatbotkit_mcpserver_integration" "example" {
  name        = "MCP Server"
  description = "Expose tools via MCP protocol"
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

resource "chatbotkit_skillset_ability" "lookup" {
  skillset_id = chatbotkit_skillset.api_tools.id
  name        = "lookup"
  description = "Look up specific records"
  instruction = "Use this to find specific items"
}

resource "chatbotkit_mcpserver_integration" "api_server" {
  name        = "API MCP Server"
  description = "Expose API tools via MCP"
  skillset_id = chatbotkit_skillset.api_tools.id
}
```

### With Blueprint

```terraform
resource "chatbotkit_blueprint" "mcp_template" {
  name        = "MCP Server Template"
  description = "Template for MCP servers"
}

resource "chatbotkit_mcpserver_integration" "from_template" {
  name         = "Custom MCP Server"
  description  = "Created from template"
  blueprint_id = chatbotkit_blueprint.mcp_template.id
  skillset_id  = chatbotkit_skillset.tools.id
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional) The name of the integration. This is displayed in the ChatBotKit dashboard.
- `description` - (Optional) A description of the integration's purpose.
- `skillset_id` - (Optional) The ID of the skillset to expose via MCP.
- `blueprint_id` - (Optional) The ID of a blueprint to associate with this integration.
- `meta` - (Optional) A map of metadata key-value pairs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier of the integration.
- `created_at` - The timestamp when the integration was created.
- `updated_at` - The timestamp when the integration was last updated.

## Import

MCP Server integrations can be imported using their ID:

```bash
terraform import chatbotkit_mcpserver_integration.example mcpserver_abc123def456
```
