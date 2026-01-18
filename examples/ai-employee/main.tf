# AI Employee Reference Architecture
#
# This example demonstrates a comprehensive AI Employee designed to work
# autonomously within professional environments. The AI Employee can perform
# tasks such as research, content creation, email management, and more.
#
# Architecture highlights:
# - Main bot with extensive professional backstory
# - Core skillset with dynamic skill loading capabilities
# - Space for sandboxed operations
# - Shell execution ability for command-line tasks
# - Email management skillset with Gmail integration
# - Notion integration for knowledge management
#
# Prerequisites:
# - Set the CHATBOTKIT_API_KEY environment variable
# - Configure Google Mail OAuth2 credentials (platform/google/mail)
# - Configure Notion OAuth2 credentials (platform/notion)

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
# AI Employee Bot
# ============================================================================
# The main AI Employee with professional capabilities

resource "chatbotkit_bot" "ai_employee" {
  name        = "AI Employee"
  description = "An AI Employee that automates tasks and enhances workplace efficiency"
  model       = "claude-4.5-opus"
  
  backstory = <<-EOT
    # PRIMARY IDENTITY SECTION
    
    You are an AI Employee, a dedicated digital team member designed to work 
    collaboratively within professional environments. Your role is to function 
    as a reliable, efficient, and knowledgeable colleague who can assist with 
    a wide range of workplace tasks, from project management and research to 
    content creation and problem-solving.
    
    Your communication style should be professional yet approachable, maintaining 
    the tone of a competent colleague rather than a subservient assistant. You 
    demonstrate initiative, provide thoughtful insights, and contribute meaningfully 
    to workplace objectives while respecting organizational hierarchies and protocols.
    
    # CORE RESPONSIBILITIES
    
    ## Research and Information Retrieval
    - Search for relevant data, market insights, industry trends
    - Gather competitive intelligence from multiple sources
    - Verify information accuracy before presenting findings
    
    ## Document Management
    - Create, edit, and organize business documents
    - Generate reports, presentations, and collaborative materials
    - Maintain version control and documentation standards
    
    ## Communication Facilitation
    - Draft professional correspondence and meeting notes
    - Manage email communications effectively
    - Provide clear project updates and status reports
    
    ## Data Analysis and Project Support
    - Process and analyze datasets for insights
    - Track tasks, timelines, and resource allocation
    - Monitor progress and identify potential roadblocks
    
    # QUALITY STANDARDS
    
    - Always verify information through multiple credible sources
    - Use proper markdown formatting for all content
    - Include citations and references where appropriate
    - Maintain professional tone in all communications
    - Respect confidentiality and privacy requirements
  EOT

  skillset_id = chatbotkit_skillset.core_skillset.id
}

# ============================================================================
# Core Skillset
# ============================================================================
# Contains core abilities for skill management and basic operations

resource "chatbotkit_skillset" "core_skillset" {
  name        = "Core Skills"
  description = "Core workplace abilities for the AI Employee"
}

# ============================================================================
# Space for Sandboxed Operations
# ============================================================================
# Provides a secure environment for file operations and task execution

resource "chatbotkit_space" "workspace" {
  name        = "AI Employee Workspace"
  description = "Secure workspace for the AI Employee to manage files and execute tasks"
}

# ============================================================================
# Core Abilities
# ============================================================================

# Shell execution ability for command-line tasks
resource "chatbotkit_skillset_ability" "shell_exec" {
  skillset_id = chatbotkit_skillset.core_skillset.id
  space_id    = chatbotkit_space.workspace.id
  
  name        = "Shell Execution"
  description = "Execute shell commands or scripts in the secure workspace"
  instruction = <<-EOT
    template: shell/exec
    parameters: {}
  EOT
}

# List available skillsets
resource "chatbotkit_skillset_ability" "list_skillsets" {
  skillset_id = chatbotkit_skillset.core_skillset.id
  
  name        = "List Skillsets"
  description = "Display available skillsets that can be installed"
  instruction = <<-EOT
    template: blueprint/resource/list
    parameters:
      type: skillset
  EOT
}

# Install skillset by ID
resource "chatbotkit_skillset_ability" "install_skillset" {
  skillset_id = chatbotkit_skillset.core_skillset.id
  
  name        = "Install Skillset"
  description = "Bring a skillset into context by its ID"
  instruction = <<-EOT
    template: conversation/skillset/install[by-id]
    parameters:
      skillsetId: ''
  EOT
}

# ============================================================================
# Mail Management Skillset
# ============================================================================
# Gmail integration for email management

resource "chatbotkit_skillset" "mail_skillset" {
  name        = "Mail Management"
  description = "Expert email communication skills for effective correspondence"
}

# Secret for Google Mail OAuth2
resource "chatbotkit_secret" "google_mail" {
  name        = "Google Mail OAuth2 Token"
  description = "OAuth2 token for accessing Google Mail"
  type        = "template"
  kind        = "personal"
  
  config = jsonencode({
    template = "platform/google/mail"
  })
}

# Gmail abilities
resource "chatbotkit_skillset_ability" "list_gmail_messages" {
  skillset_id = chatbotkit_skillset.mail_skillset.id
  secret_id   = chatbotkit_secret.google_mail.id
  
  name        = "List Gmail Messages"
  description = "Get a list of all gmail messages sorted in descending order"
  instruction = <<-EOT
    template: google/mail/message/list
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "fetch_gmail_message" {
  skillset_id = chatbotkit_skillset.mail_skillset.id
  secret_id   = chatbotkit_secret.google_mail.id
  
  name        = "Fetch Gmail Message"
  description = "Get a specific gmail message by id"
  instruction = <<-EOT
    template: google/mail/message/fetch
    parameters:
      id: ''
  EOT
}

resource "chatbotkit_skillset_ability" "search_gmail_messages" {
  skillset_id = chatbotkit_skillset.mail_skillset.id
  secret_id   = chatbotkit_secret.google_mail.id
  
  name        = "Search Gmail Messages"
  description = "Search for messages in Gmail"
  instruction = <<-EOT
    template: google/mail/message/search
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "send_gmail_message" {
  skillset_id = chatbotkit_skillset.mail_skillset.id
  secret_id   = chatbotkit_secret.google_mail.id
  
  name        = "Send Gmail Message"
  description = "Send an email using Gmail"
  instruction = <<-EOT
    template: google/mail/message/send
    parameters:
      to: ''
      subject: ''
  EOT
}

# ============================================================================
# Notion Skillset
# ============================================================================
# Notion integration for knowledge management

resource "chatbotkit_skillset" "notion_skillset" {
  name        = "Notion Management"
  description = "Advanced proficiency in using Notion for organization and collaboration"
}

# Secret for Notion OAuth2
resource "chatbotkit_secret" "notion" {
  name        = "Notion OAuth2 Token"
  description = "OAuth2 token for accessing Notion"
  type        = "template"
  kind        = "personal"
  
  config = jsonencode({
    template = "platform/notion"
  })
}

# Notion abilities
resource "chatbotkit_skillset_ability" "search_notion" {
  skillset_id = chatbotkit_skillset.notion_skillset.id
  secret_id   = chatbotkit_secret.notion.id
  
  name        = "Search Notion"
  description = "Search all of Notion for specific keywords using multiple queries"
  instruction = <<-EOT
    template: notion/search[exhaustive]
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "fetch_notion_page" {
  skillset_id = chatbotkit_skillset.notion_skillset.id
  secret_id   = chatbotkit_secret.notion.id
  
  name        = "Fetch Notion Page"
  description = "Fetch details of a specific page in Notion"
  instruction = <<-EOT
    template: notion/page/fetch
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "create_notion_database_item" {
  skillset_id = chatbotkit_skillset.notion_skillset.id
  secret_id   = chatbotkit_secret.notion.id
  
  name        = "Create Notion Database Item"
  description = "Create a new item in a Notion database"
  instruction = <<-EOT
    template: notion/database/item/create
    parameters:
      databaseId: ''
  EOT
}

# ============================================================================
# Outputs
# ============================================================================

output "bot_id" {
  description = "The ID of the AI Employee bot"
  value       = chatbotkit_bot.ai_employee.id
}

output "workspace_id" {
  description = "The ID of the AI Employee workspace"
  value       = chatbotkit_space.workspace.id
}

output "skillset_ids" {
  description = "IDs of the skillsets"
  value = {
    core   = chatbotkit_skillset.core_skillset.id
    mail   = chatbotkit_skillset.mail_skillset.id
    notion = chatbotkit_skillset.notion_skillset.id
  }
}
