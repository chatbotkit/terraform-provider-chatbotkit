terraform {
  required_providers {
    chatbotkit = {
      source = "chatbotkit/chatbotkit"
    }
  }
}

provider "chatbotkit" {
  # api_key = "..." # Or set CHATBOTKIT_API_KEY env var
}

# Example: Create a new bot
resource "chatbotkit_bot" "example" {
  name        = "My Test Bot"
  description = "Created via Terraform"
  backstory   = "You are a helpful assistant."
}

# Example: Reference an existing bot by ID
# data "chatbotkit_bot" "existing" {
#   id = "bot_abc123"
# }
#
# output "existing_bot_name" {
#   value = data.chatbotkit_bot.existing.name
# }
