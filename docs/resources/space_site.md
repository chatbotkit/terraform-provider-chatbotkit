---
page_title: "chatbotkit_space_site Resource - terraform-provider-chatbotkit"
subcategory: ""
description: |-
  Manages a ChatBotKit Space Site resource.
---

# chatbotkit_space_site (Resource)

Manages a ChatBotKit Space Site. A site binds a `<label>.chatbotkit.space` subdomain to static content served from a space's storage. It is a nested resource keyed by its parent `space_id` (like `chatbotkit_skillset_ability`).

## Example Usage

### Basic Site

```terraform
resource "chatbotkit_space" "website" {
  name        = "Website"
  description = "Static site assets"
}

resource "chatbotkit_space_site" "example" {
  space_id = chatbotkit_space.website.id
  name     = "Landing Page"
  domain   = "acme.chatbotkit.space"
}
```

### Serving From a Folder With Custom Index

```terraform
resource "chatbotkit_space" "docs" {
  name        = "Docs"
  description = "Documentation assets"
}

resource "chatbotkit_space_site" "docs_site" {
  space_id    = chatbotkit_space.docs.id
  name        = "Documentation"
  description = "Public documentation site"
  domain      = "acme-docs.chatbotkit.space"
  prefix      = "public"
  index       = "index.html"
  not_found   = "404.html"
}
```

## Argument Reference

The following arguments are supported:

- `space_id` - (Required) The ID of the space to attach this site to. Changing this forces a new resource to be created.
- `domain` - (Required) The host the site is served at (a `<label>.chatbotkit.space` subdomain).
- `name` - (Optional) The name of the site. This is displayed in the ChatBotKit dashboard.
- `description` - (Optional) A description of the site's purpose.
- `prefix` - (Optional) An optional folder prefix inside the space to serve content from.
- `index` - (Optional) The directory index filename. Defaults to `index.html`.
- `not_found` - (Optional) The not-found filename. Defaults to `404.html`.
- `alias` - (Optional) A unique alias for the site.
- `meta` - (Optional) A map of metadata key-value pairs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier of the site.
- `created_at` - The timestamp when the site was created.
- `updated_at` - The timestamp when the site was last updated.

## Import

Space sites can be imported using their ID:

```bash
terraform import chatbotkit_space_site.example site_abc123def456
```
