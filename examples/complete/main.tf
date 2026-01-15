# Complete ChatBotKit Terraform Example
#
# This example demonstrates a complete setup with:
# - A dataset for storing knowledge
# - A skillset with an ability
# - A bot linked to both the dataset and skillset
# - A trigger integration linked to the bot
#
# Prerequisites:
# - Set the CHATBOTKIT_API_KEY environment variable or configure api_key below

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

# ============================================================================
# Dataset
# ============================================================================
# A dataset stores knowledge that the bot can search and reference

resource "chatbotkit_dataset" "knowledge_base" {
  name        = "Knowledge Base"
  description = "Dataset containing company knowledge for the bot to reference"
}

# ============================================================================
# Skillset
# ============================================================================
# A skillset groups abilities that extend what the bot can do

resource "chatbotkit_skillset" "bot_skills" {
  name        = "Bot Skills"
  description = "Skills and abilities for the assistant bot"
}

# ============================================================================
# Ability
# ============================================================================
# An ability defines a specific action or capability the bot can use

resource "chatbotkit_skillset_ability" "search_ability" {
  skillset_id = chatbotkit_skillset.bot_skills.id
  name        = "Search Knowledge"
  description = "Search the knowledge base for relevant information"
  instruction = <<-EOT
    Use this ability to search for information in the knowledge base.
    When the user asks a question, use this to find relevant answers.
  EOT
}

# ============================================================================
# Bot
# ============================================================================
# The main bot that uses the dataset and skillset

resource "chatbotkit_bot" "assistant" {
  name        = "AI Assistant"
  description = "An intelligent assistant powered by ChatBotKit"
  backstory   = <<-EOT
    You are a helpful AI assistant. You have access to a knowledge base
    that you can search to answer questions. Always be polite and helpful.
    If you don't know something, say so honestly.
  EOT

  # Link the bot to the dataset
  dataset_id = chatbotkit_dataset.knowledge_base.id

  # Link the bot to the skillset
  skillset_id = chatbotkit_skillset.bot_skills.id
}

# ============================================================================
# Trigger Integration
# ============================================================================
# A trigger integration allows the bot to be invoked via webhooks or schedules

resource "chatbotkit_trigger_integration" "bot_trigger" {
  name        = "Bot Trigger"
  description = "Trigger integration for invoking the assistant bot"

  # Link to the bot
  bot_id = chatbotkit_bot.assistant.id

  # Optional: Set a schedule (e.g., "hourly", "daily", "weekly")
  # trigger_schedule = "never"
}

# ============================================================================
# Outputs
# ============================================================================
# Output the IDs of created resources for reference

output "dataset_id" {
  description = "The ID of the created dataset"
  value       = chatbotkit_dataset.knowledge_base.id
}

output "skillset_id" {
  description = "The ID of the created skillset"
  value       = chatbotkit_skillset.bot_skills.id
}

output "ability_id" {
  description = "The ID of the created ability"
  value       = chatbotkit_skillset_ability.search_ability.id
}

output "bot_id" {
  description = "The ID of the created bot"
  value       = chatbotkit_bot.assistant.id
}

output "trigger_integration_id" {
  description = "The ID of the created trigger integration"
  value       = chatbotkit_trigger_integration.bot_trigger.id
}
