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

# Create a dataset
resource "chatbotkit_dataset" "knowledge_base" {
  name        = "Product Knowledge Base"
  description = "Contains information about our products"
  type        = "text"
}

# Create a skillset
resource "chatbotkit_skillset" "abilities" {
  name        = "Customer Support Abilities"
  description = "Custom abilities for customer support"
}

# Create a secret
resource "chatbotkit_secret" "api_key" {
  name  = "external_api_key"
  value = var.external_api_key
}

# Create a file
resource "chatbotkit_file" "documentation" {
  name   = "product-guide.pdf"
  type   = "application/pdf"
  source = "https://example.com/docs/product-guide.pdf"
}

# Create a bot with dataset and skillset
resource "chatbotkit_bot" "support_bot" {
  name         = "Customer Support Bot"
  description  = "Automated customer support assistant"
  model        = "gpt-4"
  backstory    = "You are a knowledgeable customer support agent helping users with product questions."
  temperature  = 0.7
  instructions = "Always be helpful, professional, and accurate. Use the knowledge base to answer questions."
  moderation   = true
  privacy      = true
  dataset_id   = chatbotkit_dataset.knowledge_base.id
  skillset_id  = chatbotkit_skillset.abilities.id
}

# Create an integration
resource "chatbotkit_integration" "slack" {
  name        = "Slack Support Channel"
  description = "Connect support bot to Slack"
  type        = "slack"
  bot_id      = chatbotkit_bot.support_bot.id
}

# Data sources - list all resources
data "chatbotkit_bots" "all" {}

data "chatbotkit_datasets" "all" {}

# Variables
variable "external_api_key" {
  description = "External API key for bot integrations"
  type        = string
  sensitive   = true
}

# Outputs
output "bot_id" {
  value       = chatbotkit_bot.support_bot.id
  description = "The ID of the support bot"
}

output "bot_created_at" {
  value       = chatbotkit_bot.support_bot.created_at
  description = "When the bot was created"
}

output "integration_id" {
  value       = chatbotkit_integration.slack.id
  description = "The ID of the Slack integration"
}

output "total_bots" {
  value       = length(data.chatbotkit_bots.all.bots)
  description = "Total number of bots"
}
