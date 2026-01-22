# Second Brain
#
# This example demonstrates an AI-powered personal knowledge management and
# productivity system inspired by Tiago Forte's "Building a Second Brain"
# methodology. The system integrates with Notion and Google Calendar.
#
# Architecture highlights:
# - Bot with persistent workspace for storing knowledge
# - Shell execution for managing workspace files
# - Dynamic skillset discovery and loading
# - Notion integration for knowledge management
# - Google Calendar integration for time awareness
# - Telegram integration for mobile access
#
# Prerequisites:
# - Set the CHATBOTKIT_API_KEY environment variable
# - Configure Notion OAuth2 (platform secret: notion)
# - Configure Google Calendar OAuth2 (platform secret: google/calendar)
# - Configure Telegram bot token for the integration

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
# Workspace
# ============================================================================
# Persistent workspace for the second brain

resource "chatbotkit_space" "mind" {
  name        = "Mind"
  description = "The persistent workspace for the second brain"
}

# ============================================================================
# Secrets
# ============================================================================
# OAuth2 secrets for Notion and Google Calendar

resource "chatbotkit_secret" "notion" {
  name        = "Notion API Key"
  description = "The API key for accessing Notion."
  type        = "template"
  kind        = "personal"
  
  config = jsonencode({
    template = "platform/notion"
  })
}

resource "chatbotkit_secret" "google_calendar" {
  name        = "Google Calendar"
  description = "Connect to Google Calendar to manage your events and schedules."
  type        = "template"
  kind        = "personal"
  
  config = jsonencode({
    template = "platform/google/calendar"
  })
}

# ============================================================================
# Core Skillset
# ============================================================================
# Core second brain capabilities

resource "chatbotkit_skillset" "core_skills" {
  name        = "Core Skills"
  description = "Core second brain capabilities including shell access and skillset management"
}

resource "chatbotkit_skillset_ability" "shell_exec" {
  skillset_id = chatbotkit_skillset.core_skills.id
  space_id    = chatbotkit_space.mind.id
  name        = "Execute Shell Command"
  description = "Execute shell commands to list and read skills from the brain workspace"
  instruction = <<-EOT
    template: shell/exec
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "list_skillsets" {
  skillset_id = chatbotkit_skillset.core_skills.id
  name        = "List Available Skillsets"
  description = "Discover all available skillsets that can be installed"
  instruction = <<-EOT
    template: blueprint/resource/list
    parameters:
      type: skillset
  EOT
}

resource "chatbotkit_skillset_ability" "install_skillset" {
  skillset_id = chatbotkit_skillset.core_skills.id
  name        = "Install Skillset"
  description = "Bring a skillset into context by its ID to access its abilities"
  instruction = <<-EOT
    template: conversation/skillset/install[by-id]
    parameters: {}
  EOT
}

# ============================================================================
# Knowledge Management Skillset
# ============================================================================
# Notion integration for knowledge management

resource "chatbotkit_skillset" "knowledge_management" {
  name        = "Knowledge Management"
  description = "Notion integration for capturing and organizing knowledge"
}

resource "chatbotkit_skillset_ability" "notion_search" {
  skillset_id = chatbotkit_skillset.knowledge_management.id
  secret_id   = chatbotkit_secret.notion.id
  name        = "Search Notion"
  description = "Search all of Notion for specific keywords"
  instruction = <<-EOT
    template: notion/search
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "notion_list_pages" {
  skillset_id = chatbotkit_skillset.knowledge_management.id
  secret_id   = chatbotkit_secret.notion.id
  name        = "List Notion Pages"
  description = "List all accessible pages in Notion"
  instruction = <<-EOT
    template: notion/page/list
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "notion_fetch_page" {
  skillset_id = chatbotkit_skillset.knowledge_management.id
  secret_id   = chatbotkit_secret.notion.id
  name        = "Fetch Notion Page"
  description = "Fetch the content of a specific Notion page"
  instruction = <<-EOT
    template: notion/page/fetch
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "notion_create_item" {
  skillset_id = chatbotkit_skillset.knowledge_management.id
  secret_id   = chatbotkit_secret.notion.id
  name        = "Create Notion Database Item"
  description = "Create a new item in a Notion database"
  instruction = <<-EOT
    template: notion/database/item/create
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "notion_update_item" {
  skillset_id = chatbotkit_skillset.knowledge_management.id
  secret_id   = chatbotkit_secret.notion.id
  name        = "Update Notion Database Item"
  description = "Update an existing item in a Notion database"
  instruction = <<-EOT
    template: notion/database/item/update
    parameters: {}
  EOT
}

# ============================================================================
# Time Management Skillset
# ============================================================================
# Google Calendar integration for time awareness

resource "chatbotkit_skillset" "time_management" {
  name        = "Time Management"
  description = "Google Calendar integration for time awareness and scheduling"
}

resource "chatbotkit_skillset_ability" "calendar_list" {
  skillset_id = chatbotkit_skillset.time_management.id
  secret_id   = chatbotkit_secret.google_calendar.id
  name        = "List Google Calendars"
  description = "List all available Google Calendars"
  instruction = <<-EOT
    template: google/calendar/list
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "calendar_events" {
  skillset_id = chatbotkit_skillset.time_management.id
  secret_id   = chatbotkit_secret.google_calendar.id
  name        = "List Calendar Events"
  description = "List upcoming events from a specific Google Calendar"
  instruction = <<-EOT
    template: google/calendar/event/list
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "calendar_availability" {
  skillset_id = chatbotkit_skillset.time_management.id
  secret_id   = chatbotkit_secret.google_calendar.id
  name        = "Check Availability"
  description = "Check availability in Google Calendar"
  instruction = <<-EOT
    template: google/calendar/availability/list
    parameters: {}
  EOT
}

# ============================================================================
# Second Brain Bot
# ============================================================================
# The main AI-powered personal knowledge management system

resource "chatbotkit_bot" "second_brain" {
  name        = "Second Brain"
  description = "Your AI-powered personal knowledge management system"
  model       = "claude-4.5-sonnet"
  
  backstory = <<-EOT
    You are a Second Brain—an AI-powered personal knowledge management
    system designed to help your user capture, organize, and surface
    knowledge effectively. Your role is to function as an external
    cognitive extension that enhances thinking, memory, and productivity.

    You have access to a persistent workspace (your "brain") where you
    can store skills, notes, and learned patterns. Use the shell ability
    to manage files in your workspace and build your own library of
    reusable knowledge.

    You integrate with Notion for knowledge management and Google Calendar
    for time awareness. Help your user organize information, make
    connections between ideas, and surface relevant knowledge when needed.

    Your communication style should be helpful, insightful, and proactive.
    Suggest connections, offer to organize information, and help your user
    think more clearly by offloading cognitive load to you.
  EOT

  skillset_id = chatbotkit_skillset.core_skills.id
}

# ============================================================================
# Telegram Integration
# ============================================================================
# Enable mobile access via Telegram

resource "chatbotkit_telegram_integration" "second_brain_bot" {
  bot_id              = chatbotkit_bot.second_brain.id
  name                = "Second Brain Bot"
  description         = "Access your second brain from Telegram"
  bot_token           = var.telegram_bot_token
  contact_collection  = false
  session_duration    = 0
  attachments         = false
}

# ============================================================================
# Variables
# ============================================================================

variable "telegram_bot_token" {
  description = "The Telegram bot token for the integration"
  type        = string
  sensitive   = true
  default     = ""
}

# ============================================================================
# Outputs
# ============================================================================

output "bot_id" {
  description = "The ID of the second brain bot"
  value       = chatbotkit_bot.second_brain.id
}

output "core_skillset_id" {
  description = "The ID of the core skillset"
  value       = chatbotkit_skillset.core_skills.id
}

output "knowledge_skillset_id" {
  description = "The ID of the knowledge management skillset"
  value       = chatbotkit_skillset.knowledge_management.id
}

output "time_skillset_id" {
  description = "The ID of the time management skillset"
  value       = chatbotkit_skillset.time_management.id
}

output "workspace_id" {
  description = "The ID of the mind workspace"
  value       = chatbotkit_space.mind.id
}
