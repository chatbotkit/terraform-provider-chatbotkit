---
page_title: "chatbotkit_blueprint Resource - terraform-provider-chatbotkit"
subcategory: ""
description: |-
  Manages a ChatBotKit Blueprint resource.
---

# chatbotkit_blueprint (Resource)

Manages a ChatBotKit Blueprint. Blueprints are templates that define reusable configurations for bots, datasets, skillsets, and integrations. They allow you to create standardized AI solutions that can be instantiated multiple times.

## Example Usage

### Basic Blueprint

```terraform
resource "chatbotkit_blueprint" "example" {
  name        = "Customer Support Template"
  description = "A template for customer support bots"
}
```

### Blueprint with Visibility Settings

```terraform
resource "chatbotkit_blueprint" "shared" {
  name        = "Shared Support Template"
  description = "A publicly available support bot template"
  visibility  = "public"
  
  meta = {
    category = "support"
    version  = "1.0"
  }
}
```

### Using Blueprint with Bot

```terraform
resource "chatbotkit_blueprint" "template" {
  name        = "FAQ Bot Template"
  description = "Template for FAQ bots"
}

resource "chatbotkit_bot" "faq_bot" {
  name         = "FAQ Bot"
  description  = "Answers frequently asked questions"
  blueprint_id = chatbotkit_blueprint.template.id
  backstory    = "You are a helpful FAQ assistant."
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional) The name of the blueprint. This is displayed in the ChatBotKit dashboard.
- `description` - (Optional) A description of the blueprint's purpose.
- `visibility` - (Optional) The visibility level of the blueprint. Can be "private" or "public".
- `meta` - (Optional) A map of metadata key-value pairs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier of the blueprint.
- `created_at` - The timestamp when the blueprint was created.
- `updated_at` - The timestamp when the blueprint was last updated.

## Import

Blueprints can be imported using their ID:

```bash
terraform import chatbotkit_blueprint.example blueprint_abc123def456
```
