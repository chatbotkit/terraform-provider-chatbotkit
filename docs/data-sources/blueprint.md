---
page_title: "chatbotkit_blueprint Data Source - terraform-provider-chatbotkit"
subcategory: ""
description: |-
  Use this data source to read information about an existing ChatBotKit Blueprint.
---

# chatbotkit_blueprint (Data Source)

Use this data source to read information about an existing ChatBotKit Blueprint. This is useful when you need to reference a blueprint that was created outside of Terraform or in a different Terraform configuration.

## Example Usage

### Read an Existing Blueprint

```terraform
data "chatbotkit_blueprint" "existing" {
  id = "blueprint_abc123def456"
}

output "blueprint_name" {
  value = data.chatbotkit_blueprint.existing.name
}

output "blueprint_visibility" {
  value = data.chatbotkit_blueprint.existing.visibility
}
```

### Use with Bot

```terraform
variable "blueprint_id" {
  description = "The ID of an existing blueprint"
  type        = string
}

data "chatbotkit_blueprint" "template" {
  id = var.blueprint_id
}

resource "chatbotkit_bot" "from_template" {
  name         = "Bot from ${data.chatbotkit_blueprint.template.name}"
  description  = "Created using existing blueprint"
  blueprint_id = data.chatbotkit_blueprint.template.id
}
```

### Use with Dataset

```terraform
data "chatbotkit_blueprint" "shared" {
  id = var.shared_blueprint_id
}

resource "chatbotkit_dataset" "knowledge" {
  name         = "Knowledge Base"
  description  = "Dataset linked to ${data.chatbotkit_blueprint.shared.name}"
  blueprint_id = data.chatbotkit_blueprint.shared.id
}
```

## Argument Reference

The following arguments are required:

- `id` - (Required) The unique identifier of the blueprint to read.

## Attribute Reference

The following attributes are exported:

- `id` - The unique identifier of the blueprint.
- `name` - The name of the blueprint.
- `description` - The description of the blueprint.
- `visibility` - The visibility setting of the blueprint.
- `meta` - A map of metadata key-value pairs.
- `created_at` - The timestamp when the blueprint was created.
- `updated_at` - The timestamp when the blueprint was last updated.
