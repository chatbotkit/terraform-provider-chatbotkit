---
page_title: "chatbotkit_extract_integration Resource - terraform-provider-chatbotkit"
subcategory: ""
description: |-
  Manages a ChatBotKit Extract Integration resource.
---

# chatbotkit_extract_integration (Resource)

Manages a ChatBotKit Extract Integration. This integration allows you to extract structured data from conversations based on a defined JSON schema and send it to a webhook endpoint.

## Example Usage

### Basic Extract Integration

```terraform
resource "chatbotkit_bot" "assistant" {
  name        = "Data Collection Bot"
  description = "Collects information from users"
  backstory   = "You are an assistant that helps collect user information."
}

resource "chatbotkit_extract_integration" "example" {
  name        = "Lead Capture"
  description = "Extract lead information from conversations"
  bot_id      = chatbotkit_bot.assistant.id
  request     = "https://webhook.example.com/leads"
  
  schema = {
    type = "object"
  }
}
```

### Full Configuration

```terraform
resource "chatbotkit_extract_integration" "advanced" {
  name        = "Survey Extractor"
  description = "Extract survey responses"
  bot_id      = chatbotkit_bot.assistant.id
  
  request = "https://webhook.example.com/surveys"
  
  schema = {
    type       = "object"
    properties = jsonencode({
      name = {
        type        = "string"
        description = "User's full name"
      }
      email = {
        type        = "string"
        description = "User's email address"
      }
      rating = {
        type        = "integer"
        description = "Satisfaction rating 1-5"
      }
      feedback = {
        type        = "string"
        description = "Additional feedback"
      }
    })
  }
}
```

### With Blueprint

```terraform
resource "chatbotkit_blueprint" "extract_template" {
  name        = "Extract Integration Template"
  description = "Template for data extraction"
}

resource "chatbotkit_extract_integration" "from_template" {
  name         = "Contact Extractor"
  description  = "Created from template"
  blueprint_id = chatbotkit_blueprint.extract_template.id
  bot_id       = chatbotkit_bot.assistant.id
  request      = "https://webhook.example.com/contacts"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional) The name of the integration. This is displayed in the ChatBotKit dashboard.
- `description` - (Optional) A description of the integration's purpose.
- `bot_id` - (Optional) The ID of the ChatBotKit bot to connect.
- `request` - (Optional) The webhook URL to send extracted data to.
- `schema` - (Optional) A map defining the JSON schema for data extraction.
- `blueprint_id` - (Optional) The ID of a blueprint to associate with this integration.
- `meta` - (Optional) A map of metadata key-value pairs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier of the integration.
- `created_at` - The timestamp when the integration was created.
- `updated_at` - The timestamp when the integration was last updated.

## Import

Extract integrations can be imported using their ID:

```bash
terraform import chatbotkit_extract_integration.example extract_abc123def456
```
