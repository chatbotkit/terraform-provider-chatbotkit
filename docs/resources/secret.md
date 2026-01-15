---
page_title: "chatbotkit_secret Resource - terraform-provider-chatbotkit"
subcategory: ""
description: |-
  Manages a ChatBotKit Secret resource.
---

# chatbotkit_secret (Resource)

Manages a ChatBotKit Secret. Secrets store sensitive credentials like API keys, tokens, and passwords that are used by bots and abilities to authenticate with external services.

## Example Usage

### Basic Secret (Bearer Token)

```terraform
resource "chatbotkit_secret" "api_key" {
  name        = "External API Key"
  description = "API key for external service"
  type        = "bearer"
  value       = var.api_key
}
```

### Personal Secret

```terraform
resource "chatbotkit_secret" "personal_token" {
  name        = "User Token"
  description = "Personal access token"
  type        = "bearer"
  kind        = "personal"
  value       = var.user_token
}
```

### Secret with Configuration

```terraform
resource "chatbotkit_secret" "oauth" {
  name        = "OAuth Credentials"
  description = "OAuth configuration for service"
  type        = "oauth"
  
  config = {
    client_id     = var.client_id
    client_secret = var.client_secret
    scope         = "read write"
  }
}
```

### Secret with Blueprint

```terraform
resource "chatbotkit_blueprint" "secrets_template" {
  name        = "Secrets Template"
  description = "Template for secret configurations"
}

resource "chatbotkit_secret" "with_blueprint" {
  name         = "Service Credentials"
  description  = "Credentials linked to template"
  type         = "bearer"
  blueprint_id = chatbotkit_blueprint.secrets_template.id
  value        = var.service_token
}
```

### Using Secret with Skillset Ability

```terraform
resource "chatbotkit_secret" "weather_api" {
  name        = "Weather API Key"
  description = "API key for weather service"
  type        = "bearer"
  value       = var.weather_api_key
}

resource "chatbotkit_skillset" "tools" {
  name        = "Weather Tools"
  description = "Tools for weather queries"
}

resource "chatbotkit_skillset_ability" "get_weather" {
  skillset_id = chatbotkit_skillset.tools.id
  name        = "get_weather"
  description = "Get current weather for a location"
  secret_id   = chatbotkit_secret.weather_api.id
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional) The name of the secret. This is displayed in the ChatBotKit dashboard.
- `description` - (Optional) A description of what the secret is used for.
- `type` - (Optional) The type of secret. Common values include "bearer", "oauth", or "plain".
- `kind` - (Optional) The kind of secret. Can be "personal" for user-specific secrets or left empty for shared secrets.
- `value` - (Optional) The actual secret value (e.g., API key, token).
- `blueprint_id` - (Optional) The ID of a blueprint to associate with this secret.
- `config` - (Optional) A map of additional configuration for the secret (e.g., OAuth settings).
- `visibility` - (Optional) The visibility level of the secret. Can be "private" or "public".
- `meta` - (Optional) A map of metadata key-value pairs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier of the secret.
- `created_at` - The timestamp when the secret was created.
- `updated_at` - The timestamp when the secret was last updated.

## Import

Secrets can be imported using their ID:

```bash
terraform import chatbotkit_secret.example secret_abc123def456
```

~> **Note:** When importing secrets, the `value` field will not be imported for security reasons. You will need to set it again in your Terraform configuration.
