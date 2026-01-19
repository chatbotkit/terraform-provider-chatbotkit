---
page_title: "chatbotkit_skillset Data Source - terraform-provider-chatbotkit"
subcategory: ""
description: |-
  Use this data source to read information about an existing ChatBotKit Skillset.
---

# chatbotkit_skillset (Data Source)

Use this data source to read information about an existing ChatBotKit Skillset. This is useful when you need to reference a skillset that was created outside of Terraform or in a different Terraform configuration.

## Example Usage

### Read an Existing Skillset

```terraform
data "chatbotkit_skillset" "existing" {
  id = "skillset_abc123def456"
}

output "skillset_name" {
  value = data.chatbotkit_skillset.existing.name
}

output "skillset_visibility" {
  value = data.chatbotkit_skillset.existing.visibility
}
```

### Use with Bot

```terraform
variable "shared_skillset_id" {
  description = "The ID of a shared skillset"
  type        = string
}

data "chatbotkit_skillset" "tools" {
  id = var.shared_skillset_id
}

resource "chatbotkit_bot" "assistant" {
  name        = "Assistant"
  description = "Bot using shared tools: ${data.chatbotkit_skillset.tools.name}"
  backstory   = "You are a helpful assistant with access to various tools."
  skillset_id = data.chatbotkit_skillset.tools.id
}
```

### Use with MCP Server Integration

```terraform
data "chatbotkit_skillset" "api_tools" {
  id = var.api_skillset_id
}

resource "chatbotkit_mcpserver_integration" "mcp" {
  name        = "MCP Server"
  description = "Expose ${data.chatbotkit_skillset.api_tools.name} via MCP"
  skillset_id = data.chatbotkit_skillset.api_tools.id
}
```

### Use with Skillset Ability

```terraform
data "chatbotkit_skillset" "existing" {
  id = var.skillset_id
}

resource "chatbotkit_skillset_ability" "new_ability" {
  skillset_id = data.chatbotkit_skillset.existing.id
  name        = "new_tool"
  description = "A new ability added to ${data.chatbotkit_skillset.existing.name}"
  instruction = "Use this tool when needed"
}
```

## Argument Reference

The following arguments are required:

- `id` - (Required) The unique identifier of the skillset to read.

## Attribute Reference

The following attributes are exported:

- `id` - The unique identifier of the skillset.
- `name` - The name of the skillset.
- `description` - The description of the skillset.
- `blueprint_id` - The ID of the associated blueprint, if any.
- `visibility` - The visibility setting of the skillset.
- `meta` - A map of metadata key-value pairs.
- `created_at` - The timestamp when the skillset was created.
- `updated_at` - The timestamp when the skillset was last updated.
