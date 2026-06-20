---
page_title: "chatbotkit_context Resource - terraform-provider-chatbotkit"
subcategory: ""
description: |-
  Manages a ChatBotKit Context resource.
---

# chatbotkit_context (Resource)

Manages a ChatBotKit Context. A context binds the current account (or, under `run_as`, a partner sub-account) to a set of platform resources — a blueprint, bot, dataset, or skillset — along with a free-form payload. It is commonly used to scope a sub-account to a pre-defined configuration, or to attach per-account data (such as a repository or project id) that an agent reads at runtime.

## Example Usage

### Basic Context

```terraform
resource "chatbotkit_context" "project" {
  name        = "project"
  description = "Repository and project for this account's agent"

  payload = {
    githubOwner     = "acme-co"
    githubRepo      = "landing"
    githubRepoUrl   = "https://github.com/acme-co/landing"
    vercelProjectId = "prj_acme_landing"
  }
}
```

### Context linked to a bot

```terraform
resource "chatbotkit_bot" "support" {
  name      = "Support Bot"
  backstory = "You are a helpful support agent."
}

resource "chatbotkit_context" "onboarding" {
  name        = "Customer Onboarding Context"
  description = "Links the customer to the onboarding bot"
  bot_id      = chatbotkit_bot.support.id

  payload = {
    tier   = "premium"
    locale = "en-US"
  }
}
```

### Per-sub-account context (with run_as)

```terraform
# A provider alias whose run_as targets a partner sub-account; the context is
# created inside that sub-account.
provider "chatbotkit" {
  alias  = "customer"
  run_as = var.customer_account_id
}

resource "chatbotkit_context" "scoped" {
  provider = chatbotkit.customer
  name     = "project"

  payload = {
    githubOwner = "customer-co"
    githubRepo  = "site"
  }
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional) The name of the context.
- `description` - (Optional) A description of the context's purpose.
- `blueprint_id` - (Optional) The ID of a blueprint to link.
- `bot_id` - (Optional) The ID of a bot to link.
- `dataset_id` - (Optional) The ID of a dataset to link.
- `skillset_id` - (Optional) The ID of a skillset to link.
- `payload` - (Optional) A map of string key-value pairs attached to the context (e.g. a repository owner/name, a project id). Read by agents at runtime.
- `meta` - (Optional) A map of metadata key-value pairs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier of the context.
- `created_at` - The timestamp when the context was created.
- `updated_at` - The timestamp when the context was last updated.

## Import

Contexts can be imported using their ID:

```bash
terraform import chatbotkit_context.example context_abc123def456
```
