---
page_title: "chatbotkit_policy Resource - terraform-provider-chatbotkit"
subcategory: ""
description: |-
  Manages a ChatBotKit Policy resource.
---

# chatbotkit_policy (Resource)

Manages a ChatBotKit Policy. Policies control behaviors such as data retention and usage limits for bots and blueprints.

## Example Usage

```terraform
resource "chatbotkit_policy" "retention" {
  name        = "30-day retention"
  description = "Delete conversation data after 30 days"
  type        = "retention"
  config = {
    days = "30"
  }
}
```

## Argument Reference

The following arguments are supported:

- `alias` - (Optional) The alias ID for the policy.
- `blueprint_id` - (Optional) The ID of the blueprint to use.
- `bot_id` - (Optional) The ID of the bot to associate.
- `config` - (Optional) The policy configuration as JSON.
- `description` - (Optional) The description of the policy.
- `meta` - (Optional) Additional metadata for the policy.
- `name` - (Optional) The name of the policy.
- `type` - (Required) The type of the policy.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier of the policy.
- `created_at` - Timestamp when the resource was created.
- `updated_at` - Timestamp when the resource was last updated.

## Import

This resource can be imported using its ID:

```bash
terraform import chatbotkit_policy.example policy_abc123def456
```
