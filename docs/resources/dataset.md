---
page_title: "chatbotkit_dataset Resource - terraform-provider-chatbotkit"
subcategory: ""
description: |-
  Manages a ChatBotKit Dataset resource.
---

# chatbotkit_dataset (Resource)

Manages a ChatBotKit Dataset. Datasets are knowledge bases that store information for AI bots to retrieve and use during conversations. They support vector-based semantic search for accurate information retrieval.

## Example Usage

### Basic Dataset

```terraform
resource "chatbotkit_dataset" "example" {
  name        = "Product Knowledge Base"
  description = "Contains product documentation and FAQs"
}
```

### Dataset with Search Configuration

```terraform
resource "chatbotkit_dataset" "advanced" {
  name        = "Technical Documentation"
  description = "Technical docs with optimized search settings"
  
  search_max_records = 5
  search_max_tokens  = 2000
  search_min_score   = 0.7
  record_max_tokens  = 500
  
  match_instruction    = "Use this information to answer the user's question."
  mismatch_instruction = "I don't have specific information about that topic."
}
```

### Dataset with Blueprint

```terraform
resource "chatbotkit_blueprint" "template" {
  name        = "Knowledge Base Template"
  description = "Template for knowledge bases"
}

resource "chatbotkit_dataset" "with_blueprint" {
  name         = "Company Knowledge Base"
  description  = "Company-wide knowledge base"
  blueprint_id = chatbotkit_blueprint.template.id
  visibility   = "private"
}
```

### Dataset with Bot

```terraform
resource "chatbotkit_dataset" "knowledge" {
  name        = "Support Knowledge Base"
  description = "Knowledge base for support bot"
}

resource "chatbotkit_bot" "support" {
  name        = "Support Bot"
  description = "Customer support assistant"
  backstory   = "You are a helpful support agent."
  dataset_id  = chatbotkit_dataset.knowledge.id
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional) The name of the dataset. This is displayed in the ChatBotKit dashboard.
- `description` - (Optional) A description of the dataset's contents.
- `blueprint_id` - (Optional) The ID of a blueprint to associate with this dataset.
- `match_instruction` - (Optional) Instructions for the bot when search results are found.
- `mismatch_instruction` - (Optional) Instructions for the bot when no search results are found.
- `record_max_tokens` - (Optional) Maximum number of tokens per record when chunking data.
- `search_max_records` - (Optional) Maximum number of records to return in search results.
- `search_max_tokens` - (Optional) Maximum total tokens across all search results.
- `search_min_score` - (Optional) Minimum similarity score (0-1) for search results to be included.
- `reranker` - (Optional) The reranking model to use for improving search relevance.
- `separators` - (Optional) Custom separators for text chunking.
- `store` - (Optional) The storage backend to use.
- `visibility` - (Optional) The visibility level of the dataset. Can be "private" or "public".
- `meta` - (Optional) A map of metadata key-value pairs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier of the dataset.
- `created_at` - The timestamp when the dataset was created.
- `updated_at` - The timestamp when the dataset was last updated.

## Import

Datasets can be imported using their ID:

```bash
terraform import chatbotkit_dataset.example dataset_abc123def456
```
