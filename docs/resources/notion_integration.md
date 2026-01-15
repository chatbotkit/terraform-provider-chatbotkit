---
page_title: "chatbotkit_notion_integration Resource - terraform-provider-chatbotkit"
subcategory: ""
description: |-
  Manages a ChatBotKit Notion Integration resource.
---

# chatbotkit_notion_integration (Resource)

Manages a ChatBotKit Notion Integration. This integration allows you to automatically sync content from Notion to a dataset, keeping your bot's knowledge base up-to-date with your Notion workspace.

## Example Usage

### Basic Notion Integration

```terraform
resource "chatbotkit_dataset" "notion_content" {
  name        = "Notion Content"
  description = "Content synced from Notion"
}

resource "chatbotkit_notion_integration" "example" {
  name        = "Notion Sync"
  description = "Sync Notion workspace content"
  dataset_id  = chatbotkit_dataset.notion_content.id
  token       = var.notion_token
}
```

### Full Configuration

```terraform
resource "chatbotkit_notion_integration" "advanced" {
  name        = "Advanced Notion Sync"
  description = "Full-featured Notion integration"
  dataset_id  = chatbotkit_dataset.notion_content.id
  
  token         = var.notion_token
  sync_schedule = "0 */4 * * *"  # Every 4 hours
  expires_in    = 604800000      # 7 days in milliseconds
}
```

### With Blueprint

```terraform
resource "chatbotkit_blueprint" "notion_template" {
  name        = "Notion Integration Template"
  description = "Template for Notion integrations"
}

resource "chatbotkit_notion_integration" "from_template" {
  name         = "Docs Sync"
  description  = "Created from template"
  blueprint_id = chatbotkit_blueprint.notion_template.id
  dataset_id   = chatbotkit_dataset.notion_content.id
  token        = var.notion_token
}
```

### Knowledge Base with Notion Source

```terraform
resource "chatbotkit_dataset" "knowledge" {
  name        = "Company Knowledge Base"
  description = "Knowledge from Notion"
}

resource "chatbotkit_notion_integration" "sync" {
  name          = "Knowledge Sync"
  description   = "Sync company knowledge from Notion"
  dataset_id    = chatbotkit_dataset.knowledge.id
  token         = var.notion_token
  sync_schedule = "0 0 * * *"  # Daily
}

resource "chatbotkit_bot" "assistant" {
  name        = "Knowledge Assistant"
  description = "Assistant with Notion knowledge"
  backstory   = "You are a helpful assistant with access to company documentation."
  dataset_id  = chatbotkit_dataset.knowledge.id
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional) The name of the integration. This is displayed in the ChatBotKit dashboard.
- `description` - (Optional) A description of the integration's purpose.
- `dataset_id` - (Optional) The ID of the dataset to sync Notion content to.
- `token` - (Optional) The Notion integration token for API access.
- `sync_schedule` - (Optional) A cron expression for automatic synchronization.
- `expires_in` - (Optional) Time in milliseconds before synced content expires.
- `blueprint_id` - (Optional) The ID of a blueprint to associate with this integration.
- `meta` - (Optional) A map of metadata key-value pairs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier of the integration.
- `created_at` - The timestamp when the integration was created.
- `updated_at` - The timestamp when the integration was last updated.

## Import

Notion integrations can be imported using their ID:

```bash
terraform import chatbotkit_notion_integration.example notion_abc123def456
```
