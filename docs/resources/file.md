---
page_title: "chatbotkit_file Resource - terraform-provider-chatbotkit"
subcategory: ""
description: |-
  Manages a ChatBotKit File resource.
---

# chatbotkit_file (Resource)

Manages a ChatBotKit File. Files are used to store documents, configurations, or other data that can be referenced by bots, abilities, or other resources.

## Example Usage

### Basic File

```terraform
resource "chatbotkit_file" "example" {
  name        = "Configuration File"
  description = "Stores configuration data"
}
```

### File with Visibility

```terraform
resource "chatbotkit_file" "shared" {
  name        = "Shared Document"
  description = "A document shared across the organization"
  visibility  = "public"
  
  meta = {
    category = "documentation"
    format   = "markdown"
  }
}
```

### File with Blueprint

```terraform
resource "chatbotkit_blueprint" "files_template" {
  name        = "Files Template"
  description = "Template for file management"
}

resource "chatbotkit_file" "with_blueprint" {
  name         = "Template File"
  description  = "File linked to template"
  blueprint_id = chatbotkit_blueprint.files_template.id
}
```

### Using File with Skillset Ability

```terraform
resource "chatbotkit_file" "ability_config" {
  name        = "Ability Configuration"
  description = "Configuration for custom ability"
}

resource "chatbotkit_skillset" "tools" {
  name        = "Custom Tools"
  description = "Custom toolset"
}

resource "chatbotkit_skillset_ability" "custom" {
  skillset_id = chatbotkit_skillset.tools.id
  name        = "custom_action"
  description = "A custom action with file configuration"
  file_id     = chatbotkit_file.ability_config.id
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional) The name of the file. This is displayed in the ChatBotKit dashboard.
- `description` - (Optional) A description of the file's contents or purpose.
- `blueprint_id` - (Optional) The ID of a blueprint to associate with this file.
- `visibility` - (Optional) The visibility level of the file. Can be "private" or "public".
- `meta` - (Optional) A map of metadata key-value pairs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier of the file.
- `created_at` - The timestamp when the file was created.
- `updated_at` - The timestamp when the file was last updated.

## Import

Files can be imported using their ID:

```bash
terraform import chatbotkit_file.example file_abc123def456
```
