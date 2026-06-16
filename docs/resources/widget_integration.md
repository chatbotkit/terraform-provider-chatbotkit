---
page_title: "chatbotkit_widget_integration Resource - terraform-provider-chatbotkit"
subcategory: ""
description: |-
  Manages a ChatBotKit Widget Integration resource.
---

# chatbotkit_widget_integration (Resource)

Manages a ChatBotKit Widget Integration, the embeddable chat widget for your website.

## Example Usage

```terraform
resource "chatbotkit_widget_integration" "example" {
  name      = "Website Widget"
  bot_id    = chatbotkit_bot.assistant.id
  theme     = "light"
  title     = "Need help?"
  stream    = true
}
```

## Argument Reference

The following arguments are supported:

- `alias` - (Optional) The alias ID.
- `attachments` - (Optional) Whether attachments are enabled.
- `auto_scroll` - (Optional) Whether auto-scroll is enabled.
- `blueprint_id` - (Optional) The ID of the blueprint to use.
- `bot_id` - (Optional) The ID of the bot to connect.
- `carousel` - (Optional) Whether the carousel is enabled.
- `contact_collection` - (Optional) Whether to collect contact information.
- `description` - (Optional) The description.
- `export_conversation` - (Optional) Whether conversation export is enabled.
- `form` - (Optional) Whether forms are enabled.
- `initial` - (Optional) The initial message.
- `intro` - (Optional) The widget intro message.
- `language` - (Optional) The widget language.
- `layout` - (Optional) The widget layout.
- `math` - (Optional) Whether math rendering is enabled.
- `maximize` - (Optional) Whether the widget can be maximized.
- `message_peek` - (Optional) Whether message peek is enabled.
- `meta` - (Optional) Additional metadata for the integration.
- `name` - (Optional) The name.
- `origin` - (Optional) The allowed origin.
- `placeholder` - (Optional) The input placeholder text.
- `plugins` - (Optional) The enabled plugins.
- `powered_by` - (Optional) Whether the powered-by label is shown.
- `restart_conversation` - (Optional) Whether conversation restart is enabled.
- `session_duration` - (Optional) The duration of the session in milliseconds.
- `start_first` - (Optional) Whether to start first.
- `stream` - (Optional) Whether to stream responses.
- `theme` - (Optional) The widget theme.
- `title` - (Optional) The widget title.
- `tools` - (Optional) Whether tools are enabled.
- `unfurl` - (Optional) Whether link unfurling is enabled.
- `verbose` - (Optional) Whether verbose mode is enabled.
- `voice_in` - (Optional) Whether voice input is enabled.
- `voice_out` - (Optional) Whether voice output is enabled.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier of the widgetintegration.
- `created_at` - Timestamp when the resource was created.
- `updated_at` - Timestamp when the resource was last updated.

## Import

This resource can be imported using its ID:

```bash
terraform import chatbotkit_widget_integration.example widget_abc123def456
```
