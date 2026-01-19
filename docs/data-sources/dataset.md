---
page_title: "chatbotkit_dataset Data Source - terraform-provider-chatbotkit"
subcategory: ""
description: |-
  Use this data source to read information about an existing ChatBotKit Dataset.
---

# chatbotkit_dataset (Data Source)

Use this data source to read information about an existing ChatBotKit Dataset. This is useful when you need to reference a dataset that was created outside of Terraform or in a different Terraform configuration.

## Example Usage

### Read an Existing Dataset

```terraform
data "chatbotkit_dataset" "existing" {
  id = "dataset_abc123def456"
}

output "dataset_name" {
  value = data.chatbotkit_dataset.existing.name
}

output "search_settings" {
  value = {
    max_records = data.chatbotkit_dataset.existing.search_max_records
    max_tokens  = data.chatbotkit_dataset.existing.search_max_tokens
    min_score   = data.chatbotkit_dataset.existing.search_min_score
  }
}
```

### Use with Bot

```terraform
variable "shared_dataset_id" {
  description = "The ID of a shared dataset"
  type        = string
}

data "chatbotkit_dataset" "knowledge" {
  id = var.shared_dataset_id
}

resource "chatbotkit_bot" "assistant" {
  name        = "Assistant"
  description = "Bot using shared dataset: ${data.chatbotkit_dataset.knowledge.name}"
  backstory   = "You are a helpful assistant with access to our knowledge base."
  dataset_id  = data.chatbotkit_dataset.knowledge.id
}
```

### Use with Sitemap Integration

```terraform
data "chatbotkit_dataset" "docs" {
  id = var.docs_dataset_id
}

resource "chatbotkit_sitemap_integration" "docs_sync" {
  name        = "Documentation Sync"
  description = "Sync content to ${data.chatbotkit_dataset.docs.name}"
  dataset_id  = data.chatbotkit_dataset.docs.id
  url         = "https://docs.example.com/sitemap.xml"
}
```

## Argument Reference

The following arguments are required:

- `id` - (Required) The unique identifier of the dataset to read.

## Attribute Reference

The following attributes are exported:

- `id` - The unique identifier of the dataset.
- `name` - The name of the dataset.
- `description` - The description of the dataset.
- `blueprint_id` - The ID of the associated blueprint, if any.
- `match_instruction` - Instructions for the bot when search results are found.
- `mismatch_instruction` - Instructions for the bot when no search results are found.
- `record_max_tokens` - Maximum number of tokens per record.
- `search_max_records` - Maximum number of records in search results.
- `search_max_tokens` - Maximum total tokens in search results.
- `search_min_score` - Minimum similarity score for search results.
- `reranker` - The reranking model used.
- `separators` - Custom separators for text chunking.
- `store` - The storage backend used.
- `visibility` - The visibility setting of the dataset.
- `meta` - A map of metadata key-value pairs.
- `created_at` - The timestamp when the dataset was created.
- `updated_at` - The timestamp when the dataset was last updated.
