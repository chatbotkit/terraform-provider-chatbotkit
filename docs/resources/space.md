---
page_title: "chatbotkit_space Resource - terraform-provider-chatbotkit"
subcategory: ""
description: |-
  Manages a ChatBotKit Space resource.
---

# chatbotkit_space (Resource)

Manages a ChatBotKit Space. Spaces group conversations and resources, optionally scoped to a blueprint or contact.

## Example Usage

```terraform
resource "chatbotkit_space" "team" {
  name        = "Team Space"
  description = "Shared space for the team"
}
```

## Argument Reference

The following arguments are supported:

- `alias` - (Optional) The alias ID for the space.
- `blueprint_id` - (Optional) The ID of the blueprint to use.
- `contact_id` - (Optional) The ID of the contact to associate.
- `description` - (Optional) The description of the space.
- `meta` - (Optional) Additional metadata for the space.
- `name` - (Optional) The name of the space.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier of the space.
- `created_at` - Timestamp when the resource was created.
- `updated_at` - Timestamp when the resource was last updated.

## Import

This resource can be imported using its ID:

```bash
terraform import chatbotkit_space.example space_abc123def456
```
