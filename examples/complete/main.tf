# Complete ChatBotKit Terraform Example
#
# This example demonstrates a complete setup with:
# - A dataset for storing knowledge (automatically adds search functions)
# - A skillset with web search and fetch abilities
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
# A dataset stores knowledge that the bot can search and reference.
# Note: Linking a dataset to a bot automatically adds search functions.

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
# Abilities
# ============================================================================
# Abilities define specific actions or capabilities the bot can use.
# Here we add web search and fetch abilities for internet access.

# Ability to search the web for information
resource "chatbotkit_skillset_ability" "web_search" {
  skillset_id = chatbotkit_skillset.bot_skills.id
  name        = "Search Web"
  description = "Search the web for information"
  instruction = <<-EOT
    ```search
    query: $[query! ys|the search query to find information on the web]
    ```
  EOT
}

# Ability to fetch and read web pages
resource "chatbotkit_skillset_ability" "web_fetch" {
  skillset_id = chatbotkit_skillset.bot_skills.id
  name        = "Fetch Web Page"
  description = "Fetch and read the content of a web page"
  instruction = <<-EOT
    ```fetch
    url: $[url! ys|the URL of the web page to fetch and read]
    ```
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
    You are a helpful AI assistant. You have access to:
    - A knowledge base that you can search to answer questions
    - Web search capabilities to find current information online
    - The ability to fetch and read web pages

    Always be polite and helpful. If you don't know something, try searching
    the knowledge base first, then the web if needed.
  EOT

  # Link the bot to the dataset (automatically adds search functions)
  dataset_id = chatbotkit_dataset.knowledge_base.id

  # Link the bot to the skillset (includes web search and fetch abilities)
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

output "web_search_ability_id" {
  description = "The ID of the web search ability"
  value       = chatbotkit_skillset_ability.web_search.id
}

output "web_fetch_ability_id" {
  description = "The ID of the web fetch ability"
  value       = chatbotkit_skillset_ability.web_fetch.id
}

output "bot_id" {
  description = "The ID of the created bot"
  value       = chatbotkit_bot.assistant.id
}

output "trigger_integration_id" {
  description = "The ID of the created trigger integration"
  value       = chatbotkit_trigger_integration.bot_trigger.id
}
