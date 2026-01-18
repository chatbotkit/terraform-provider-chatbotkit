# Matillion Data Pipeline Operations Center Reference Architecture
#
# This example demonstrates an AI-powered operations center for Matillion Data
# Productivity Cloud that automates pipeline management, monitoring, and incident
# response.
#
# Architecture highlights:
# - Main operations agent bot
# - Four specialized skillsets: Pipeline Operations, Schedule Management,
#   Monitoring & Compliance, Infrastructure Management
# - Matillion API integration
# - Slack integration for team notifications
# - Trigger integration for scheduled health checks
#
# Prerequisites:
# - Set the CHATBOTKIT_API_KEY environment variable
# - Configure Matillion API credentials
# - Configure Slack credentials (optional)

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
# Matillion Operations Agent
# ============================================================================

resource "chatbotkit_bot" "matillion_ops" {
  name        = "Matillion Operations Agent"
  description = "An intelligent agent for managing Matillion Data Productivity Cloud operations"
  model       = "claude-4.5-sonnet"
  
  backstory = <<-EOT
    # PRIMARY IDENTITY SECTION
    
    You are the Matillion Operations Agent, an intelligent assistant specialized 
    in managing data pipeline operations on Matillion Data Productivity Cloud. 
    Your role is to help data engineering teams efficiently manage their ETL/ELT 
    pipelines, monitor executions, and maintain infrastructure health.
    
    Your communication style is professional and technically precise, as you work 
    with data engineers and operations teams who expect accurate, actionable 
    information. You provide clear status updates, proactive recommendations, and 
    handle routine operations autonomously while escalating complex issues 
    appropriately.
    
    # CORE RESPONSIBILITIES
    
    ## Pipeline Operations
    - Execute pipelines on-demand when requested
    - Monitor running executions and report status
    - Cancel problematic or stuck pipeline runs
    - Track execution history and identify patterns
    
    ## Schedule Management
    - Create new pipeline schedules using cron expressions
    - Modify existing schedules (timing, enabled/disabled)
    - Delete obsolete schedules
    - Review and optimize scheduling patterns
    
    ## Monitoring & Alerting
    - Track pipeline execution status across all projects
    - Monitor credit consumption and alert on anomalies
    - Review audit events for compliance
    - Track data lineage for governance requirements
    
    ## Infrastructure Management
    - Monitor Matillion agent health and status
    - Restart, pause, or resume agents as needed
    - Review agent configurations and recommend optimizations
    
    # OPERATIONAL GUIDELINES
    
    ## Before Executing Pipelines
    1. Confirm the project, environment, and pipeline name
    2. Check if there are any currently running executions
    3. Verify the agent status if execution is critical
    4. Execute and provide the execution ID for tracking
    
    ## For Schedule Changes
    1. List current schedules to understand existing patterns
    2. Validate cron expressions before creating schedules
    3. Consider timezone implications
    4. Confirm changes with the requestor
    
    ## For Incident Response
    1. Gather execution details and error information
    2. Check for patterns in recent failures
    3. Review agent status if execution issues persist
    4. Escalate infrastructure issues appropriately
    
    ## For Cost Monitoring
    1. Track consumption trends over time
    2. Alert when usage exceeds normal patterns
    3. Recommend optimization opportunities
    
    # SAFETY PROTOCOLS
    
    - Always confirm before deleting schedules or projects
    - Warn before cancelling running executions
    - Provide clear feedback on all operations
    - Log significant actions for audit purposes
    
    The current date is $${EARTH_DATE}. Failure to follow these guidelines may 
    result in operational issues.
  EOT

  skillset_id = chatbotkit_skillset.main_skillset.id
}

# ============================================================================
# Main Skillset
# ============================================================================

resource "chatbotkit_skillset" "main_skillset" {
  name        = "Main Operations Skillset"
  description = "Core skillset for dynamic skill loading"
}

# Abilities for discovering and loading specialized skillsets
resource "chatbotkit_skillset_ability" "list_skillsets" {
  skillset_id = chatbotkit_skillset.main_skillset.id
  
  name        = "List Available Skillsets"
  description = "Discover all available skillsets in this blueprint that can be installed"
  instruction = <<-EOT
    template: blueprint/resource/list
    parameters:
      type: skillset
  EOT
}

resource "chatbotkit_skillset_ability" "install_skillset" {
  skillset_id = chatbotkit_skillset.main_skillset.id
  
  name        = "Install Skillset"
  description = "Bring a skillset into context by its ID to access its abilities"
  instruction = <<-EOT
    template: conversation/skillset/install[by-id]
    parameters:
      skillsetId: ''
  EOT
}

# ============================================================================
# Secret for Matillion API
# ============================================================================

resource "chatbotkit_secret" "matillion_api" {
  name        = "Matillion API Token"
  description = "The API token for accessing Matillion Data Productivity Cloud APIs"
  type        = "template"
  
  config = jsonencode({
    template = "matillion"
  })
}

# ============================================================================
# Pipeline Operations Skillset
# ============================================================================

resource "chatbotkit_skillset" "pipeline_ops" {
  name        = "Pipeline Operations"
  description = "Abilities for executing, monitoring, and managing Matillion data pipelines"
}

resource "chatbotkit_skillset_ability" "list_projects" {
  skillset_id = chatbotkit_skillset.pipeline_ops.id
  secret_id   = chatbotkit_secret.matillion_api.id
  
  name        = "List Projects"
  description = "Retrieve a list of all projects in the Matillion account"
  instruction = <<-EOT
    template: matillion/project/list
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "list_pipelines" {
  skillset_id = chatbotkit_skillset.pipeline_ops.id
  secret_id   = chatbotkit_secret.matillion_api.id
  
  name        = "List Pipelines"
  description = "Retrieve a list of all published pipelines in a project environment"
  instruction = <<-EOT
    template: matillion/pipeline/list
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "execute_pipeline" {
  skillset_id = chatbotkit_skillset.pipeline_ops.id
  secret_id   = chatbotkit_secret.matillion_api.id
  
  name        = "Execute Pipeline"
  description = "Execute a published pipeline in Matillion"
  instruction = <<-EOT
    template: matillion/pipeline/execute
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "list_executions" {
  skillset_id = chatbotkit_skillset.pipeline_ops.id
  secret_id   = chatbotkit_secret.matillion_api.id
  
  name        = "List Pipeline Executions"
  description = "Retrieve a list of pipeline executions with optional filters"
  instruction = <<-EOT
    template: matillion/pipeline-execution/list
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "cancel_execution" {
  skillset_id = chatbotkit_skillset.pipeline_ops.id
  secret_id   = chatbotkit_secret.matillion_api.id
  
  name        = "Cancel Pipeline Execution"
  description = "Cancel a running pipeline execution"
  instruction = <<-EOT
    template: matillion/pipeline-execution/cancel
    parameters: {}
  EOT
}

# ============================================================================
# Schedule Management Skillset
# ============================================================================

resource "chatbotkit_skillset" "schedule_mgmt" {
  name        = "Schedule Management"
  description = "Abilities for creating, managing, and monitoring pipeline schedules"
}

resource "chatbotkit_skillset_ability" "list_schedules" {
  skillset_id = chatbotkit_skillset.schedule_mgmt.id
  secret_id   = chatbotkit_secret.matillion_api.id
  
  name        = "List Schedules"
  description = "Retrieve a list of all schedules for a project"
  instruction = <<-EOT
    template: matillion/schedule/list
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "create_schedule" {
  skillset_id = chatbotkit_skillset.schedule_mgmt.id
  secret_id   = chatbotkit_secret.matillion_api.id
  
  name        = "Create Schedule"
  description = "Create a new schedule for a pipeline"
  instruction = <<-EOT
    template: matillion/schedule/create
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "update_schedule" {
  skillset_id = chatbotkit_skillset.schedule_mgmt.id
  secret_id   = chatbotkit_secret.matillion_api.id
  
  name        = "Update Schedule"
  description = "Update an existing schedule (timing, enabled status)"
  instruction = <<-EOT
    template: matillion/schedule/update
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "delete_schedule" {
  skillset_id = chatbotkit_skillset.schedule_mgmt.id
  secret_id   = chatbotkit_secret.matillion_api.id
  
  name        = "Delete Schedule"
  description = "Delete a schedule from a project"
  instruction = <<-EOT
    template: matillion/schedule/delete
    parameters: {}
  EOT
}

# ============================================================================
# Monitoring & Compliance Skillset
# ============================================================================

resource "chatbotkit_skillset" "monitoring" {
  name        = "Monitoring & Compliance"
  description = "Abilities for monitoring consumption, audit events, and data lineage"
}

resource "chatbotkit_skillset_ability" "get_consumption" {
  skillset_id = chatbotkit_skillset.monitoring.id
  secret_id   = chatbotkit_secret.matillion_api.id
  
  name        = "Get Credit Consumption"
  description = "Retrieve credit consumption breakdown for cost monitoring"
  instruction = <<-EOT
    template: matillion/consumption/fetch
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "list_audit_events" {
  skillset_id = chatbotkit_skillset.monitoring.id
  secret_id   = chatbotkit_secret.matillion_api.id
  
  name        = "List Audit Events"
  description = "Retrieve audit events for compliance and troubleshooting"
  instruction = <<-EOT
    template: matillion/audit-event/list
    parameters: {}
  EOT
}

# ============================================================================
# Infrastructure Management Skillset
# ============================================================================

resource "chatbotkit_skillset" "infrastructure" {
  name        = "Infrastructure Management"
  description = "Abilities for managing Matillion agents and artifacts"
}

resource "chatbotkit_skillset_ability" "list_agents" {
  skillset_id = chatbotkit_skillset.infrastructure.id
  secret_id   = chatbotkit_secret.matillion_api.id
  
  name        = "List Agents"
  description = "Retrieve a list of all agents in the Matillion account"
  instruction = <<-EOT
    template: matillion/agent/list
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "send_agent_command" {
  skillset_id = chatbotkit_skillset.infrastructure.id
  secret_id   = chatbotkit_secret.matillion_api.id
  
  name        = "Send Agent Command"
  description = "Send a command to an agent (RESTART, PAUSE, or RESUME)"
  instruction = <<-EOT
    template: matillion/agent/command/send
    parameters: {}
  EOT
}

# ============================================================================
# Integrations (Optional)
# ============================================================================

# Slack integration for team notifications
# resource "chatbotkit_slack_integration" "matillion_slack" {
#   bot_id           = chatbotkit_bot.matillion_ops.id
#   name             = "Matillion Operations Bot"
#   description      = "Slack bot for interacting with Matillion pipeline operations"
#   signing_secret   = "your-slack-signing-secret"
#   bot_token        = "your-slack-bot-token"
#   session_duration = 0
# }

# Trigger integration for scheduled health checks
resource "chatbotkit_trigger_integration" "daily_health_check" {
  bot_id           = chatbotkit_bot.matillion_ops.id
  name             = "Daily Health Check"
  description      = "Scheduled trigger for daily pipeline health monitoring and status reports"
  authenticate     = true
  session_duration = 1800000  # 30 minutes
  trigger_schedule = "daily"
}

# ============================================================================
# Outputs
# ============================================================================

output "bot_id" {
  description = "The ID of the Matillion Operations Agent"
  value       = chatbotkit_bot.matillion_ops.id
}

output "skillset_ids" {
  description = "IDs of the specialized skillsets"
  value = {
    main           = chatbotkit_skillset.main_skillset.id
    pipeline_ops   = chatbotkit_skillset.pipeline_ops.id
    schedule_mgmt  = chatbotkit_skillset.schedule_mgmt.id
    monitoring     = chatbotkit_skillset.monitoring.id
    infrastructure = chatbotkit_skillset.infrastructure.id
  }
}

output "trigger_integration_id" {
  description = "The ID of the daily health check trigger"
  value       = chatbotkit_trigger_integration.daily_health_check.id
}
