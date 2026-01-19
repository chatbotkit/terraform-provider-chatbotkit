# MCP Factory Reference Architectures
#
# This example demonstrates a factory-style architecture for MCP servers that
# expose multiple skillsets and abilities through separate MCP server integrations.
# Each skillset is designed to encapsulate specific functionalities, allowing for
# modular and organized management of AI capabilities.
#
# Architecture highlights:
# - Multiple skillsets, each with dedicated abilities
# - Separate MCP server integration for each skillset
# - Modular, factory-style organization
# - Clear separation of concerns
#
# Prerequisites:
# - Set the CHATBOTKIT_API_KEY environment variable

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
# Service 1: Data Analytics
# ============================================================================
# First skillset with data analytics capabilities

resource "chatbotkit_skillset" "analytics_service" {
  name        = "Analytics Service"
  description = "Data analytics and reporting capabilities"
}

resource "chatbotkit_skillset_ability" "analytics_analyze" {
  skillset_id = chatbotkit_skillset.analytics_service.id
  name        = "Analyze Data"
  description = "Perform comprehensive data analysis"
  instruction = <<-EOT
    # Analytics logic here
    # This is a placeholder for actual analytics implementation
  EOT
}

resource "chatbotkit_skillset_ability" "analytics_report" {
  skillset_id = chatbotkit_skillset.analytics_service.id
  name        = "Generate Report"
  description = "Create detailed analytics reports"
  instruction = <<-EOT
    # Report generation logic here
    # This is a placeholder for actual reporting implementation
  EOT
}

resource "chatbotkit_mcp_server_integration" "analytics_mcp" {
  skillset_id = chatbotkit_skillset.analytics_service.id
  name        = "Analytics MCP Server"
  description = "MCP server exposing analytics capabilities"
}

# ============================================================================
# Service 2: Content Management
# ============================================================================
# Second skillset with content management capabilities

resource "chatbotkit_skillset" "content_service" {
  name        = "Content Service"
  description = "Content creation and management capabilities"
}

resource "chatbotkit_skillset_ability" "content_create" {
  skillset_id = chatbotkit_skillset.content_service.id
  name        = "Create Content"
  description = "Generate new content based on specifications"
  instruction = <<-EOT
    # Content creation logic here
    # This is a placeholder for actual content creation implementation
  EOT
}

resource "chatbotkit_skillset_ability" "content_edit" {
  skillset_id = chatbotkit_skillset.content_service.id
  name        = "Edit Content"
  description = "Modify and refine existing content"
  instruction = <<-EOT
    # Content editing logic here
    # This is a placeholder for actual editing implementation
  EOT
}

resource "chatbotkit_mcp_server_integration" "content_mcp" {
  skillset_id = chatbotkit_skillset.content_service.id
  name        = "Content MCP Server"
  description = "MCP server exposing content management capabilities"
}

# ============================================================================
# Service 3: Research
# ============================================================================
# Third skillset with research capabilities

resource "chatbotkit_skillset" "research_service" {
  name        = "Research Service"
  description = "Research and information gathering capabilities"
}

resource "chatbotkit_skillset_ability" "research_search" {
  skillset_id = chatbotkit_skillset.research_service.id
  name        = "Search Information"
  description = "Search for relevant information on topics"
  instruction = <<-EOT
    # Research search logic here
    # This is a placeholder for actual search implementation
  EOT
}

resource "chatbotkit_skillset_ability" "research_summarize" {
  skillset_id = chatbotkit_skillset.research_service.id
  name        = "Summarize Findings"
  description = "Create concise summaries of research findings"
  instruction = <<-EOT
    # Summarization logic here
    # This is a placeholder for actual summarization implementation
  EOT
}

resource "chatbotkit_mcp_server_integration" "research_mcp" {
  skillset_id = chatbotkit_skillset.research_service.id
  name        = "Research MCP Server"
  description = "MCP server exposing research capabilities"
}

# ============================================================================
# Service 4: Task Management
# ============================================================================
# Fourth skillset with task management capabilities

resource "chatbotkit_skillset" "tasks_service" {
  name        = "Task Management Service"
  description = "Task planning and execution capabilities"
}

resource "chatbotkit_skillset_ability" "tasks_plan" {
  skillset_id = chatbotkit_skillset.tasks_service.id
  name        = "Plan Tasks"
  description = "Create structured task plans"
  instruction = <<-EOT
    # Task planning logic here
    # This is a placeholder for actual planning implementation
  EOT
}

resource "chatbotkit_skillset_ability" "tasks_track" {
  skillset_id = chatbotkit_skillset.tasks_service.id
  name        = "Track Progress"
  description = "Monitor and report on task progress"
  instruction = <<-EOT
    # Progress tracking logic here
    # This is a placeholder for actual tracking implementation
  EOT
}

resource "chatbotkit_mcp_server_integration" "tasks_mcp" {
  skillset_id = chatbotkit_skillset.tasks_service.id
  name        = "Tasks MCP Server"
  description = "MCP server exposing task management capabilities"
}

# ============================================================================
# Outputs
# ============================================================================

output "mcp_server_urls" {
  description = "URLs for all MCP server integrations"
  value = {
    analytics = "Use the MCP server URL from the ChatBotKit dashboard for: ${chatbotkit_mcp_server_integration.analytics_mcp.name}"
    content   = "Use the MCP server URL from the ChatBotKit dashboard for: ${chatbotkit_mcp_server_integration.content_mcp.name}"
    research  = "Use the MCP server URL from the ChatBotKit dashboard for: ${chatbotkit_mcp_server_integration.research_mcp.name}"
    tasks     = "Use the MCP server URL from the ChatBotKit dashboard for: ${chatbotkit_mcp_server_integration.tasks_mcp.name}"
  }
}

output "skillset_ids" {
  description = "IDs of all skillsets"
  value = {
    analytics = chatbotkit_skillset.analytics_service.id
    content   = chatbotkit_skillset.content_service.id
    research  = chatbotkit_skillset.research_service.id
    tasks     = chatbotkit_skillset.tasks_service.id
  }
}
