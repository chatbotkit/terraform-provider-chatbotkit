---
page_title: "chatbotkit_sitemap_integration Resource - terraform-provider-chatbotkit"
subcategory: ""
description: |-
  Manages a ChatBotKit Sitemap Integration resource.
---

# chatbotkit_sitemap_integration (Resource)

Manages a ChatBotKit Sitemap Integration. This integration allows you to automatically crawl and sync website content from a sitemap to a dataset, keeping your bot's knowledge base up-to-date with your website.

## Example Usage

### Basic Sitemap Integration

```terraform
resource "chatbotkit_dataset" "website_content" {
  name        = "Website Content"
  description = "Content from company website"
}

resource "chatbotkit_sitemap_integration" "example" {
  name        = "Website Crawler"
  description = "Crawl and sync website content"
  dataset_id  = chatbotkit_dataset.website_content.id
  url         = "https://example.com/sitemap.xml"
}
```

### Full Configuration

```terraform
resource "chatbotkit_sitemap_integration" "advanced" {
  name        = "Advanced Website Crawler"
  description = "Full-featured sitemap crawler"
  dataset_id  = chatbotkit_dataset.website_content.id
  
  url           = "https://example.com/sitemap.xml"
  glob          = "https://example.com/docs/**"
  selectors     = "article, .content, main"
  javascript    = true
  sync_schedule = "0 0 * * *"  # Daily at midnight
  expires_in    = 604800000    # 7 days in milliseconds
}
```

### With Blueprint

```terraform
resource "chatbotkit_blueprint" "sitemap_template" {
  name        = "Sitemap Integration Template"
  description = "Template for sitemap crawlers"
}

resource "chatbotkit_sitemap_integration" "from_template" {
  name         = "Docs Crawler"
  description  = "Created from template"
  blueprint_id = chatbotkit_blueprint.sitemap_template.id
  dataset_id   = chatbotkit_dataset.website_content.id
  url          = "https://docs.example.com/sitemap.xml"
}
```

### Multi-Site Configuration

```terraform
resource "chatbotkit_dataset" "knowledge_base" {
  name        = "Combined Knowledge Base"
  description = "Content from multiple sources"
}

resource "chatbotkit_sitemap_integration" "main_site" {
  name        = "Main Site"
  description = "Main website content"
  dataset_id  = chatbotkit_dataset.knowledge_base.id
  url         = "https://www.example.com/sitemap.xml"
  glob        = "https://www.example.com/products/**"
}

resource "chatbotkit_sitemap_integration" "blog" {
  name        = "Blog"
  description = "Blog content"
  dataset_id  = chatbotkit_dataset.knowledge_base.id
  url         = "https://blog.example.com/sitemap.xml"
  sync_schedule = "0 */6 * * *"  # Every 6 hours
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional) The name of the integration. This is displayed in the ChatBotKit dashboard.
- `description` - (Optional) A description of the integration's purpose.
- `dataset_id` - (Optional) The ID of the dataset to sync crawled content to.
- `url` - (Optional) The URL of the sitemap to crawl.
- `glob` - (Optional) A glob pattern to filter which URLs to crawl.
- `selectors` - (Optional) CSS selectors to extract specific content from pages.
- `javascript` - (Optional) Whether to enable JavaScript rendering for dynamic content.
- `sync_schedule` - (Optional) A cron expression for automatic synchronization.
- `expires_in` - (Optional) Time in milliseconds before crawled content expires.
- `blueprint_id` - (Optional) The ID of a blueprint to associate with this integration.
- `meta` - (Optional) A map of metadata key-value pairs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier of the integration.
- `created_at` - The timestamp when the integration was created.
- `updated_at` - The timestamp when the integration was last updated.

## Import

Sitemap integrations can be imported using their ID:

```bash
terraform import chatbotkit_sitemap_integration.example sitemap_abc123def456
```
