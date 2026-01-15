---
page_title: "chatbotkit_portal Resource - terraform-provider-chatbotkit"
subcategory: ""
description: |-
  Manages a ChatBotKit Portal resource.
---

# chatbotkit_portal (Resource)

Manages a ChatBotKit Portal. Portals are customizable web interfaces for bots that can be embedded or shared as standalone pages. They provide a way to deploy bots with a branded user experience.

## Example Usage

### Basic Portal

```terraform
resource "chatbotkit_portal" "example" {
  name        = "Customer Support Portal"
  description = "Self-service support portal"
}
```

### Portal with Custom Slug

```terraform
resource "chatbotkit_portal" "branded" {
  name        = "Help Center"
  description = "Company help center portal"
  slug        = "help-center"
}
```

### Portal with Configuration

```terraform
resource "chatbotkit_portal" "configured" {
  name        = "Configured Portal"
  description = "Portal with custom settings"
  slug        = "my-portal"
  
  config = {
    theme       = "dark"
    primaryColor = "#007bff"
  }
  
  meta = {
    department = "support"
    version    = "2.0"
  }
}
```

### Portal with Blueprint

```terraform
resource "chatbotkit_blueprint" "portal_template" {
  name        = "Portal Template"
  description = "Template for support portals"
}

resource "chatbotkit_portal" "from_template" {
  name         = "Support Portal"
  description  = "Created from template"
  blueprint_id = chatbotkit_blueprint.portal_template.id
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional) The name of the portal. This is displayed in the ChatBotKit dashboard.
- `description` - (Optional) A description of the portal's purpose.
- `slug` - (Optional) A custom URL slug for the portal. Must be unique across all portals.
- `blueprint_id` - (Optional) The ID of a blueprint to associate with this portal.
- `config` - (Optional) A map of configuration settings for the portal.
- `meta` - (Optional) A map of metadata key-value pairs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier of the portal.
- `created_at` - The timestamp when the portal was created.
- `updated_at` - The timestamp when the portal was last updated.

## Import

Portals can be imported using their ID:

```bash
terraform import chatbotkit_portal.example portal_abc123def456
```
