---
page_title: "chatbotkit_skillset Resource - terraform-provider-chatbotkit"
subcategory: ""
description: |-
  Manages a ChatBotKit Skillset resource.
---

# chatbotkit_skillset (Resource)

Manages a ChatBotKit Skillset. Skillsets are collections of abilities (tools) that bots can use to perform actions, such as searching the web, sending emails, or integrating with external APIs.

## Example Usage

### Basic Skillset

```terraform
resource "chatbotkit_skillset" "example" {
  name        = "Customer Support Tools"
  description = "Tools for customer support operations"
}
```

### Skillset with Bot

```terraform
resource "chatbotkit_skillset" "tools" {
  name        = "Support Tools"
  description = "Tools for support bot"
}

resource "chatbotkit_bot" "support" {
  name        = "Support Bot"
  description = "Customer support assistant"
  backstory   = "You are a helpful support agent with access to various tools."
  skillset_id = chatbotkit_skillset.tools.id
}
```

### Skillset with Blueprint

```terraform
resource "chatbotkit_blueprint" "template" {
  name        = "Tools Template"
  description = "Template for skillsets"
}

resource "chatbotkit_skillset" "with_blueprint" {
  name         = "Advanced Tools"
  description  = "Advanced toolset for AI agents"
  blueprint_id = chatbotkit_blueprint.template.id
  visibility   = "private"
}
```

### Complete Bot with Dataset and Skillset

```terraform
resource "chatbotkit_dataset" "knowledge" {
  name        = "Product Knowledge"
  description = "Product documentation"
}

resource "chatbotkit_skillset" "tools" {
  name        = "Product Tools"
  description = "Tools for product queries"
}

resource "chatbotkit_bot" "product_assistant" {
  name        = "Product Assistant"
  description = "Helps with product questions"
  backstory   = "You are a product expert with access to documentation and tools."
  
  dataset_id  = chatbotkit_dataset.knowledge.id
  skillset_id = chatbotkit_skillset.tools.id
  
  model = "gpt-4"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional) The name of the skillset. This is displayed in the ChatBotKit dashboard.
- `description` - (Optional) A description of the skillset's purpose and capabilities.
- `blueprint_id` - (Optional) The ID of a blueprint to associate with this skillset.
- `visibility` - (Optional) The visibility level of the skillset. Can be "private" or "public".
- `meta` - (Optional) A map of metadata key-value pairs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier of the skillset.
- `created_at` - The timestamp when the skillset was created.
- `updated_at` - The timestamp when the skillset was last updated.

## Import

Skillsets can be imported using their ID:

```bash
terraform import chatbotkit_skillset.example skillset_abc123def456
```
