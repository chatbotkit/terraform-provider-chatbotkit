---
page_title: "chatbotkit_support_integration Resource - terraform-provider-chatbotkit"
subcategory: ""
description: |-
  Manages a ChatBotKit Support Integration resource.
---

# chatbotkit_support_integration (Resource)

Manages a ChatBotKit Support Integration, which routes escalations to a human support email.

## Example Usage

```terraform
resource "chatbotkit_support_integration" "example" {
  name        = "Support Escalation"
  description = "Escalate to the support team"
  bot_id      = chatbotkit_bot.assistant.id
  email       = "support@example.com"
}
```

## Argument Reference

The following arguments are supported:

- `alias` - (Optional) The alias ID.
- `blueprint_id` - (Optional) The ID of the blueprint to use.
- `bot_id` - (Optional) The ID of the bot to connect.
- `description` - (Optional) The description.
- `email` - (Optional) The support email address.
- `meta` - (Optional) Additional metadata for the integration.
- `name` - (Optional) The name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier of the supportintegration.
- `created_at` - Timestamp when the resource was created.
- `updated_at` - Timestamp when the resource was last updated.

## Import

This resource can be imported using its ID:

```bash
terraform import chatbotkit_support_integration.example support_abc123def456
```
