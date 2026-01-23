# Workflow Orchestrator Agent
#
# This example demonstrates a sophisticated architecture for building complex,
# multi-step automation workflows with dynamic skillset loading, workflow state
# persistence, and comprehensive execution tracing.
#
# Architecture highlights:
# - Main bot with orchestration capabilities
# - Dynamic skillset discovery and installation
# - Three specialized workflow skillsets (Data, Control, Reporting)
# - Space for persistent workflow state and traces
# - Trigger integration for workflow execution
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
# Workflow Workspace
# ============================================================================
# Persistent storage for workflow state, traces, and execution artifacts

resource "chatbotkit_space" "workflow_workspace" {
  name        = "Workflow Workspace"
  description = "Persistent storage for workflow state, traces, and execution artifacts"
}

# ============================================================================
# Orchestration Core Skillset
# ============================================================================
# Core orchestration capabilities including dynamic skillset loading

resource "chatbotkit_skillset" "orchestration_core" {
  name        = "Orchestration Core"
  description = "Core orchestration capabilities including dynamic skillset loading"
}

resource "chatbotkit_skillset_ability" "list_skillsets" {
  skillset_id = chatbotkit_skillset.orchestration_core.id
  name        = "List Available Skillsets"
  description = "Discover all workflow skillsets available in this blueprint"
  instruction = <<-EOT
    template: blueprint/resource/list
    parameters:
      type: skillset
  EOT
}

resource "chatbotkit_skillset_ability" "install_skillset" {
  skillset_id = chatbotkit_skillset.orchestration_core.id
  name        = "Install Skillset"
  description = "Dynamically load a workflow skillset by its ID to access its capabilities"
  instruction = <<-EOT
    template: conversation/skillset/install[by-id]
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "shell_exec" {
  skillset_id = chatbotkit_skillset.orchestration_core.id
  space_id    = chatbotkit_space.workflow_workspace.id
  name        = "Execute Workflow Command"
  description = "Execute shell commands in the workflow workspace for state management and logging"
  instruction = <<-EOT
    template: shell/exec
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "shell_read" {
  skillset_id = chatbotkit_skillset.orchestration_core.id
  space_id    = chatbotkit_space.workflow_workspace.id
  name        = "Read Workflow File"
  description = "Read workflow state, logs, or configuration files from the workspace"
  instruction = <<-EOT
    template: shell/read
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "shell_write" {
  skillset_id = chatbotkit_skillset.orchestration_core.id
  space_id    = chatbotkit_space.workflow_workspace.id
  name        = "Write Workflow File"
  description = "Write workflow state, traces, or results to the workspace"
  instruction = <<-EOT
    template: shell/write
    parameters: {}
  EOT
}

# ============================================================================
# Data Processing Skillset
# ============================================================================
# Workflow abilities for data transformation, validation, and quality checks

resource "chatbotkit_skillset" "data_processing" {
  name        = "Data Processing"
  description = "Workflow abilities for data transformation, validation, and quality checks"
}

resource "chatbotkit_skillset_ability" "fetch_web_data" {
  skillset_id = chatbotkit_skillset.data_processing.id
  name        = "Fetch Web Data"
  description = "Retrieve data from web sources for workflow processing"
  instruction = <<-EOT
    template: fetch/text/get
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "search_web" {
  skillset_id = chatbotkit_skillset.data_processing.id
  name        = "Search Web for Information"
  description = "Search the web to gather data required by the workflow"
  instruction = <<-EOT
    template: search/web
    parameters: {}
  EOT
}

# ============================================================================
# Execution Control Skillset
# ============================================================================
# Workflow abilities for state management, error handling, and flow control

resource "chatbotkit_skillset" "execution_control" {
  name        = "Execution Control"
  description = "Workflow abilities for state management, error handling, and flow control"
}

resource "chatbotkit_skillset_ability" "control_script" {
  skillset_id = chatbotkit_skillset.execution_control.id
  space_id    = chatbotkit_space.workflow_workspace.id
  name        = "Execute Control Script"
  description = "Run workflow control scripts for state transitions and flow management"
  instruction = <<-EOT
    template: shell/exec
    parameters: {}
  EOT
}

# ============================================================================
# Reporting & Audit Skillset
# ============================================================================
# Workflow abilities for generating reports, logs, and audit trails

resource "chatbotkit_skillset" "reporting_audit" {
  name        = "Reporting & Audit"
  description = "Workflow abilities for generating reports, logs, and audit trails"
}

resource "chatbotkit_skillset_ability" "generate_report" {
  skillset_id = chatbotkit_skillset.reporting_audit.id
  space_id    = chatbotkit_space.workflow_workspace.id
  name        = "Generate Execution Report"
  description = "Create detailed execution reports and audit logs"
  instruction = <<-EOT
    template: shell/exec
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "send_email" {
  skillset_id = chatbotkit_skillset.reporting_audit.id
  name        = "Send Workflow Summary Email"
  description = "Email workflow execution summaries to stakeholders"
  instruction = <<-EOT
    template: email/send
    parameters: {}
  EOT
}

# ============================================================================
# Workflow Orchestrator Bot
# ============================================================================
# The main orchestrator bot that executes multi-step workflows

resource "chatbotkit_bot" "orchestrator" {
  name        = "Workflow Orchestrator"
  description = "Multi-step workflow execution engine with dynamic capability loading and comprehensive tracing"
  model       = "claude-4.5-sonnet"
  
  backstory = <<-EOT
    You are a Workflow Orchestrator Agent responsible for executing complex
    multi-step workflows with precision and comprehensive logging. Your role
    is to coordinate specialized workflow capabilities, manage execution
    state, and produce detailed traces of all operations.

    ORCHESTRATION PRINCIPLES:

    1. Dynamic Capability Loading
      - Discover available workflow skillsets at runtime
      - Load only the skillsets needed for the current workflow
      - Understand each skillset's purpose before using its abilities

    2. Workflow Execution
      - Execute workflow steps in the correct sequence
      - Validate inputs and outputs at each step
      - Handle errors gracefully with appropriate retry logic
      - Maintain idempotency where possible

    3. State Management
      - Write workflow state to persistent storage after each step
      - Store complete context to enable pause/resume
      - Track execution history for audit trails
      - Use clear, consistent state file naming conventions

    4. Comprehensive Tracing
      - Log every workflow step with timestamp and details
      - Document decision points and the reasoning behind choices
      - Capture both successful operations and errors
      - Generate execution summaries for human review

    WORKFLOW STRUCTURE:

    Your workflows should follow this pattern:
    - Initialization: Load required skillsets and validate inputs
    - Execution: Process each workflow step with state persistence
    - Validation: Verify outputs meet requirements
    - Completion: Generate execution report and clean up resources

    TRACING FORMAT:

    Create trace logs in your workspace with this structure:
    ```
    # Workflow Execution Trace
    Workflow: [workflow name]
    Started: [timestamp]
    Status: [in_progress|completed|failed]

    ## Step 1: [Step Name]
    - Time: [timestamp]
    - Action: [description]
    - Input: [input data]
    - Output: [output data]
    - Status: [success|failed]
    - Notes: [any relevant observations]

    ## Summary
    - Total Steps: [number]
    - Successful: [number]
    - Failed: [number]
    - Duration: [time]
    ```

    Store all traces in your workspace using timestamped filenames like:
    workflow-trace-[workflow-name]-[YYYY-MM-DD-HH-MM-SS].md

    The current date is $${EARTH_DATE}. Execute workflows methodically,
    document everything, and ensure complete traceability of all operations.
  EOT

  skillset_id = chatbotkit_skillset.orchestration_core.id
}

# ============================================================================
# Trigger Integration
# ============================================================================
# Manual or scheduled trigger for workflow execution

resource "chatbotkit_trigger_integration" "workflow_trigger" {
  bot_id            = chatbotkit_bot.orchestrator.id
  name              = "Workflow Execution Trigger"
  description       = "Manual or scheduled trigger for workflow execution"
  authenticate      = true
  session_duration  = 3600000
  trigger_schedule  = "never"
}

# ============================================================================
# Outputs
# ============================================================================

output "bot_id" {
  description = "The ID of the workflow orchestrator bot"
  value       = chatbotkit_bot.orchestrator.id
}

output "orchestration_skillset_id" {
  description = "The ID of the orchestration core skillset"
  value       = chatbotkit_skillset.orchestration_core.id
}

output "data_processing_skillset_id" {
  description = "The ID of the data processing skillset"
  value       = chatbotkit_skillset.data_processing.id
}

output "execution_control_skillset_id" {
  description = "The ID of the execution control skillset"
  value       = chatbotkit_skillset.execution_control.id
}

output "reporting_audit_skillset_id" {
  description = "The ID of the reporting & audit skillset"
  value       = chatbotkit_skillset.reporting_audit.id
}

output "workspace_id" {
  description = "The ID of the workflow workspace"
  value       = chatbotkit_space.workflow_workspace.id
}
