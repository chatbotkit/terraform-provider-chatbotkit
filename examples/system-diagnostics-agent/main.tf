# System Diagnostics Agent
#
# This example demonstrates a self-monitoring AI agent that tests and reports on
# its own capabilities, producing diagnostic logs about available skillsets,
# abilities, and system health.
#
# Architecture highlights:
# - Bot with self-introspection capabilities
# - Blueprint resource list ability for discovering skillsets
# - Space for persistent diagnostic log storage
# - Shell execution abilities for file operations
# - Scheduled trigger for automated diagnostics
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
# Diagnostic Workspace
# ============================================================================
# Persistent storage for diagnostic reports and logs

resource "chatbotkit_space" "diagnostic_workspace" {
  name        = "Diagnostic Workspace"
  description = "Persistent storage for diagnostic reports and logs"
}

# ============================================================================
# Diagnostics Skillset
# ============================================================================
# Core diagnostic abilities for system introspection and reporting

resource "chatbotkit_skillset" "diagnostics_core" {
  name        = "Diagnostics Core"
  description = "Core diagnostic abilities for system introspection and reporting"
}

# ============================================================================
# Diagnostic Abilities
# ============================================================================
# Abilities that enable the agent to introspect and report on itself

resource "chatbotkit_skillset_ability" "list_skillsets" {
  skillset_id = chatbotkit_skillset.diagnostics_core.id
  name        = "List Available Skillsets"
  description = "Discover all available skillsets in this blueprint to catalog capabilities"
  instruction = <<-EOT
    template: blueprint/resource/list
    parameters:
      type: skillset
  EOT
}

resource "chatbotkit_skillset_ability" "shell_exec" {
  skillset_id = chatbotkit_skillset.diagnostics_core.id
  space_id    = chatbotkit_space.diagnostic_workspace.id
  name        = "Execute Shell Command"
  description = "Execute shell commands in the diagnostic workspace for file operations and system checks"
  instruction = <<-EOT
    template: shell/exec
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "shell_read" {
  skillset_id = chatbotkit_skillset.diagnostics_core.id
  space_id    = chatbotkit_space.diagnostic_workspace.id
  name        = "Read File from Workspace"
  description = "Read diagnostic reports or logs from the workspace"
  instruction = <<-EOT
    template: shell/read
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "shell_write" {
  skillset_id = chatbotkit_skillset.diagnostics_core.id
  space_id    = chatbotkit_space.diagnostic_workspace.id
  name        = "Write File to Workspace"
  description = "Write diagnostic reports and logs to the workspace"
  instruction = <<-EOT
    template: shell/write
    parameters: {}
  EOT
}

# ============================================================================
# System Diagnostics Agent Bot
# ============================================================================
# The self-monitoring agent that produces diagnostic reports

resource "chatbotkit_bot" "diagnostics_agent" {
  name        = "System Diagnostics Agent"
  description = "A self-monitoring AI agent that tests and reports on its own capabilities"
  model       = "claude-4.5-sonnet"
  
  backstory = <<-EOT
    You are a System Diagnostics Agent responsible for monitoring and
    validating your own capabilities. Your role is to systematically test
    your available skillsets and abilities, and produce detailed diagnostic
    reports.

    DIAGNOSTIC WORKFLOW:

    1. List all available skillsets using the blueprint resource discovery ability
    2. For each skillset, document its name, description, and purpose
    3. Identify any missing or unexpected skillsets
    4. Document the total count of abilities and their categories
    5. Write a comprehensive diagnostic report to your workspace

    Your diagnostic reports should be structured, timestamped, and include:
    - Executive summary of system health
    - Complete inventory of skillsets and abilities
    - Any warnings or recommendations
    - Timestamp and execution metadata

    Store all diagnostic reports in your workspace using the shell execution
    ability. Create timestamped files in the format: diagnostics-YYYY-MM-DD-HH-MM.md

    The current date is $${EARTH_DATE}. Each diagnostic run should be thorough
    and produce actionable insights about system configuration and health.
  EOT

  skillset_id = chatbotkit_skillset.diagnostics_core.id
}

# ============================================================================
# Trigger Integration
# ============================================================================
# Automated trigger for periodic system diagnostics

resource "chatbotkit_trigger_integration" "diagnostic_schedule" {
  bot_id            = chatbotkit_bot.diagnostics_agent.id
  name              = "Diagnostic Schedule"
  description       = "Automated trigger for periodic system diagnostics and health checks"
  authenticate      = true
  session_duration  = 1800000
  trigger_schedule  = "daily"
}

# ============================================================================
# Outputs
# ============================================================================

output "bot_id" {
  description = "The ID of the system diagnostics agent"
  value       = chatbotkit_bot.diagnostics_agent.id
}

output "skillset_id" {
  description = "The ID of the diagnostics skillset"
  value       = chatbotkit_skillset.diagnostics_core.id
}

output "workspace_id" {
  description = "The ID of the diagnostic workspace"
  value       = chatbotkit_space.diagnostic_workspace.id
}
