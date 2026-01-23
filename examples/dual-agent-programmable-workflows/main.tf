# Dual-Agent Programmable Workflows
#
# This example demonstrates a two-agent architecture where one agent (Workflow Architect)
# programs custom scripts and automations, while another agent (Task Runner) executes
# them on demand. This showcases separation of concerns in AI agent design.
#
# Architecture highlights:
# - Two specialized bots with distinct roles (architect and runner)
# - Shared automation playbook file for documentation
# - Shared script workspace for script storage and execution
# - Asymmetric access: architect has read/write, runner has read-only
# - Scheduled trigger for automated task execution
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
# Shared Resources
# ============================================================================
# Resources shared between both agents

resource "chatbotkit_file" "automation_playbook" {
  name        = "Automation Playbook"
  description = "Central documentation containing all programmed scripts and execution instructions"
}

resource "chatbotkit_space" "script_workspace" {
  name        = "Script Workspace"
  description = "Shared environment where scripts are created, stored, and executed"
}

# ============================================================================
# Workflow Architect Bot
# ============================================================================
# The architect programs scripts and maintains documentation

resource "chatbotkit_bot" "architect" {
  name        = "Workflow Architect"
  description = "Programs custom scripts, automations, and maintains the playbook documentation"
  model       = "claude-4.5-opus"
  
  backstory = <<-EOT
    You are the Workflow Architect, a specialized AI agent responsible for
    designing, programming, and documenting automation workflows. Your role
    is to create powerful, reusable scripts and maintain comprehensive
    documentation that enables seamless execution by the Task Runner.

    YOUR RESPONSIBILITIES:

    1. SCRIPT DEVELOPMENT
      - Create bash scripts for automation tasks
      - Design multi-step workflows with clear dependencies
      - Test scripts in your architect workspace before publishing
      - Optimize scripts for reliability and efficiency

    2. PLAYBOOK MAINTENANCE
      - Document every script with clear usage instructions
      - Include parameter descriptions and expected outputs
      - Provide troubleshooting guidance for common issues
      - Organize the playbook with clear sections and navigation

    3. QUALITY ASSURANCE
      - Ensure all scripts have error handling
      - Document prerequisites and dependencies
      - Include example invocations for each script
      - Version your scripts with dates and change notes

    PLAYBOOK FORMAT:

    When writing to the Automation Playbook, use this structure:

    # Automation Playbook
    Last updated: [date]

    ## Available Scripts

    ### [Script Name]
    **Purpose**: What this script does
    **Usage**: How to run it
    **Parameters**: Input requirements
    **Output**: What to expect
    **Example**: Sample invocation
    **Troubleshooting**: Common issues and solutions

    ---

    Always write complete, production-ready scripts. The Task Runner
    depends entirely on your documentation—be thorough and precise.
    Never assume the runner knows anything beyond what's in the playbook.

    The current date is $${EARTH_DATE}. Include timestamps in all updates.
  EOT

  skillset_id = chatbotkit_skillset.architect_toolkit.id
}

# ============================================================================
# Architect Skillset and Abilities
# ============================================================================
# Tools for the architect to develop and document scripts

resource "chatbotkit_skillset" "architect_toolkit" {
  name        = "Architect Toolkit"
  description = "Tools for script development and playbook management"
}

resource "chatbotkit_skillset_ability" "architect_playbook" {
  skillset_id = chatbotkit_skillset.architect_toolkit.id
  file_id     = chatbotkit_file.automation_playbook.id
  name        = "Read/Write Playbook"
  description = "Read or update the Automation Playbook with scripts and documentation"
  instruction = <<-EOT
    template: file/rw
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "architect_bash" {
  skillset_id = chatbotkit_skillset.architect_toolkit.id
  space_id    = chatbotkit_space.script_workspace.id
  name        = "Bash"
  description = "Execute bash commands in the script workspace"
  instruction = <<-EOT
    template: shell/exec
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "architect_script_rw" {
  skillset_id = chatbotkit_skillset.architect_toolkit.id
  space_id    = chatbotkit_space.script_workspace.id
  name        = "Read/Write Script"
  description = "Read or write script files in the shared workspace"
  instruction = <<-EOT
    template: shell/rw
    parameters: {}
  EOT
}

# ============================================================================
# Task Runner Bot
# ============================================================================
# The runner executes scripts according to the playbook

resource "chatbotkit_bot" "runner" {
  name        = "Task Runner"
  description = "Executes programmed workflows by reading the playbook and running scripts"
  model       = "claude-4.5-sonnet"
  
  backstory = <<-EOT
    You are the Task Runner, a specialized AI agent responsible for
    executing automation workflows. Your role is to read the Automation
    Playbook, understand available scripts, and execute them precisely
    according to the documented instructions.

    YOUR RESPONSIBILITIES:

    1. PLAYBOOK CONSULTATION
      - Always read the Automation Playbook before executing tasks
      - Understand the available scripts and their purposes
      - Follow documented instructions exactly as written
      - Check for any updates or new scripts regularly

    2. TASK EXECUTION
      - Run scripts in the Runtime Workspace
      - Provide required parameters as documented
      - Capture and report execution output
      - Handle errors according to troubleshooting guidance

    3. EXECUTION DISCIPLINE
      - Never modify scripts—only execute them
      - Follow the exact syntax and parameters specified
      - Report any issues encountered during execution
      - Log execution results for tracking

    EXECUTION WORKFLOW:

    1. Receive a task request (user message or trigger)
    2. Read the Automation Playbook to find relevant scripts
    3. Identify the correct script and its parameters
    4. Execute the script in the Runtime Workspace
    5. Report the results back clearly

    IMPORTANT CONSTRAINTS:

    - You can ONLY READ the Automation Playbook, not modify it
    - Always consult the playbook before executing—never guess
    - If a script doesn't exist for a task, inform the user
    - If documentation is unclear, report the issue rather than improvise

    The current date is $${EARTH_DATE}. Log all execution timestamps.
  EOT

  skillset_id = chatbotkit_skillset.runner_toolkit.id
}

# ============================================================================
# Runner Skillset and Abilities
# ============================================================================
# Tools for the runner to read playbook and execute scripts

resource "chatbotkit_skillset" "runner_toolkit" {
  name        = "Runner Toolkit"
  description = "Tools for playbook reading and script execution"
}

resource "chatbotkit_skillset_ability" "runner_read_playbook" {
  skillset_id = chatbotkit_skillset.runner_toolkit.id
  file_id     = chatbotkit_file.automation_playbook.id
  name        = "Read Playbook"
  description = "Read the Automation Playbook to find scripts and execution instructions"
  instruction = <<-EOT
    template: file/read
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "runner_bash" {
  skillset_id = chatbotkit_skillset.runner_toolkit.id
  space_id    = chatbotkit_space.script_workspace.id
  name        = "Bash"
  description = "Execute bash commands in the shared workspace to run scripts"
  instruction = <<-EOT
    template: shell/exec
    parameters: {}
  EOT
}

# ============================================================================
# Trigger Integration
# ============================================================================
# Automated trigger for scheduled workflow execution

resource "chatbotkit_trigger_integration" "scheduled_trigger" {
  bot_id            = chatbotkit_bot.runner.id
  name              = "Scheduled Task Trigger"
  description       = "Automated trigger for scheduled workflow execution by the Task Runner"
  authenticate      = true
  session_duration  = 1800000
  trigger_schedule  = "hourly"
}

# ============================================================================
# Outputs
# ============================================================================

output "architect_bot_id" {
  description = "The ID of the Workflow Architect bot"
  value       = chatbotkit_bot.architect.id
}

output "runner_bot_id" {
  description = "The ID of the Task Runner bot"
  value       = chatbotkit_bot.runner.id
}

output "automation_playbook_id" {
  description = "The ID of the shared Automation Playbook file"
  value       = chatbotkit_file.automation_playbook.id
}

output "script_workspace_id" {
  description = "The ID of the shared Script Workspace"
  value       = chatbotkit_space.script_workspace.id
}
