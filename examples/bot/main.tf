terraform {
  required_providers {
    chatbotkit = {
      source = "chatbotkit/chatbotkit"
    }
  }
}

provider "chatbotkit" {
  # token can be set via CHATBOTKIT_TOKEN environment variable
}

# Create a bot
resource "chatbotkit_bot" "example" {
  name         = "Example Bot"
  description  = "An example chatbot created with Terraform"
  model        = "gpt-4"
  backstory    = "You are a helpful and friendly assistant."
  temperature  = 0.7
  instructions = "Always be polite and helpful."
  moderation   = true
  privacy      = true
}

# Output the bot ID
output "bot_id" {
  value = chatbotkit_bot.example.id
}

output "bot_name" {
  value = chatbotkit_bot.example.name
}
