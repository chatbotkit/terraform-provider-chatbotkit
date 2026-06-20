# Deep Researcher — an orchestrator-worker research agent, built with Terraform
#
# This example builds a deep-research agent in the shape of the modern
# orchestrator-worker pattern (the same architecture as systems like
# gpt-researcher and Anthropic's multi-agent research system): a "big brain"
# orchestrator decomposes a question, dispatches focused research to workers that
# run in parallel, supervises them, and synthesizes the results into one cited
# report.
#
# The key idea — and the reason this fits ChatBotKit cleanly — is that the
# *orchestration is not a static graph baked into Terraform*. Terraform declares
# the agents and their tools; the dynamic part (spawning sub-tasks, fanning out,
# waiting, joining, synthesizing) is the orchestrator agent doing its job at
# runtime. Tasks are runtime artifacts the agents create on demand — never
# predefined here.
#
# The three roles:
#   intake        -> a small, fast agent. Takes a request, writes a clean brief,
#                    and commissions the orchestrator as a background task.
#   orchestrator  -> the big brain. Plans sub-questions, fans out worker tasks in
#                    parallel, supervises and adapts, then writes the report.
#   researcher    -> a worker. Answers one sub-question with search + fetch and
#                    writes cited findings into the shared workspace.
#
# How the pieces map:
#   agent/<role>/instructions.md   -> each bot's backstory (its system prompt)
#   fan-out / join                 -> task/* abilities the orchestrator calls at
#                                     runtime (create + run + list + fetch)
#   shared findings "blackboard"   -> a single chatbotkit_space the workers write
#                                     to and the orchestrator reads from
#   per-role models                -> model tiering: opus to orchestrate, a
#                                     lighter model to do the legwork
#
# Editing the instruction files and re-applying updates the live agents — the
# files are the source of truth.
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
# workspace — the shared findings blackboard
# ============================================================================
# Workers write their cited findings here (under findings/); the orchestrator
# reads them back and writes the final report (report.md). This is how results
# flow between the parallel agents without any direct coupling.

resource "chatbotkit_space" "workspace" {
  name        = "Deep Research Workspace"
  description = "Shared space where workers write findings and the orchestrator writes the report"
}

# ============================================================================
# researcher (worker) — answers one sub-question, with sources
# ============================================================================

resource "chatbotkit_skillset" "researcher_tools" {
  name        = "Researcher Tools"
  description = "Search, fetch, and workspace tools for a single research worker"
}

# Search the web.
resource "chatbotkit_skillset_ability" "researcher_search" {
  skillset_id = chatbotkit_skillset.researcher_tools.id
  name        = "search_web"
  description = "Search the web for current, relevant sources"
  instruction = <<-EOT
    template: search/web
    parameters: {}
  EOT
}

# Fetch and read a web page in full.
resource "chatbotkit_skillset_ability" "researcher_fetch" {
  skillset_id = chatbotkit_skillset.researcher_tools.id
  name        = "fetch_url"
  description = "Fetch and read the full content of a web page"
  instruction = <<-EOT
    template: fetch/text/get
    parameters: {}
  EOT
}

# Read/write the shared workspace (to write findings/<task-id>.md).
resource "chatbotkit_skillset_ability" "researcher_workspace" {
  skillset_id = chatbotkit_skillset.researcher_tools.id
  space_id    = chatbotkit_space.workspace.id
  name        = "workspace"
  description = "Read and write files in the shared research workspace"
  instruction = <<-EOT
    template: shell/rw
    parameters: {}
  EOT
}

resource "chatbotkit_bot" "researcher" {
  name        = "Research Worker"
  description = "Answers a single focused sub-question with cited findings"
  model       = "claude-4.5-sonnet" # the legwork — a fast, capable model

  backstory   = file("${path.module}/agent/researcher/instructions.md")
  skillset_id = chatbotkit_skillset.researcher_tools.id
}

# ============================================================================
# orchestrator (big brain) — plans, fans out, supervises, synthesizes
# ============================================================================
# Its tools are the task/* operations scoped by-bot-id, which let it create,
# launch, list, and inspect background tasks that run the researcher bot. This
# is the fan-out + join, expressed as ordinary tool calls the agent makes when
# it decides to — not as a predefined workflow.

resource "chatbotkit_skillset" "orchestrator_tools" {
  name        = "Orchestrator Tools"
  description = "Task dispatch/supervision and workspace tools for the research orchestrator"
}

resource "chatbotkit_skillset_ability" "orc_create_subtask" {
  skillset_id = chatbotkit_skillset.orchestrator_tools.id
  name        = "create_subtask"
  description = "Create a research sub-task for a worker bot (pass the researcher bot id and the sub-question)"
  instruction = <<-EOT
    template: task/create[by-bot-id]
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "orc_run_subtask" {
  skillset_id = chatbotkit_skillset.orchestrator_tools.id
  name        = "run_subtask"
  description = "Launch a created sub-task immediately (run several to work in parallel)"
  instruction = <<-EOT
    template: task/run[by-bot-id]
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "orc_list_subtasks" {
  skillset_id = chatbotkit_skillset.orchestrator_tools.id
  name        = "list_subtasks"
  description = "List dispatched sub-tasks with their status and outcome"
  instruction = <<-EOT
    template: task/list[by-bot-id]
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "orc_fetch_subtask" {
  skillset_id = chatbotkit_skillset.orchestrator_tools.id
  name        = "fetch_subtask"
  description = "Fetch full detail of a sub-task, including its result summary once complete"
  instruction = <<-EOT
    template: task/fetch[by-bot-id]
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "orc_workspace" {
  skillset_id = chatbotkit_skillset.orchestrator_tools.id
  space_id    = chatbotkit_space.workspace.id
  name        = "workspace"
  description = "Read worker findings and write the final report in the shared workspace"
  instruction = <<-EOT
    template: shell/rw
    parameters: {}
  EOT
}

resource "chatbotkit_bot" "orchestrator" {
  name        = "Research Orchestrator"
  description = "Plans, dispatches, supervises, and synthesizes deep research"
  model       = "claude-4.5-opus" # the big brain

  # The researcher bot id is wired into the backstory so the orchestrator knows
  # which bot to dispatch sub-tasks to — Terraform resolves the dependency order.
  backstory = templatefile("${path.module}/agent/orchestrator/instructions.md.tftpl", {
    researcher_bot_id = chatbotkit_bot.researcher.id
  })
  skillset_id = chatbotkit_skillset.orchestrator_tools.id
}

# ============================================================================
# intake (small entry agent) — turns a request into a commissioned task
# ============================================================================
# This is the piece that keeps tasks out of Terraform: instead of predefining a
# research job, a small agent creates one at runtime for the orchestrator.

resource "chatbotkit_skillset" "intake_tools" {
  name        = "Intake Tools"
  description = "Commission and launch a research job for the orchestrator"
}

resource "chatbotkit_skillset_ability" "intake_commission" {
  skillset_id = chatbotkit_skillset.intake_tools.id
  name        = "commission_research"
  description = "Create a research task for the orchestrator bot (pass the orchestrator bot id and the brief)"
  instruction = <<-EOT
    template: task/create[by-bot-id]
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "intake_launch" {
  skillset_id = chatbotkit_skillset.intake_tools.id
  name        = "launch_research"
  description = "Launch the commissioned research task immediately"
  instruction = <<-EOT
    template: task/run[by-bot-id]
    parameters: {}
  EOT
}

resource "chatbotkit_bot" "intake" {
  name        = "Research Intake"
  description = "Front door: turns a request into a commissioned research job"
  model       = "claude-4.5-sonnet" # small and fast

  backstory = templatefile("${path.module}/agent/intake/instructions.md.tftpl", {
    orchestrator_bot_id = chatbotkit_bot.orchestrator.id
  })
  skillset_id = chatbotkit_skillset.intake_tools.id
}

# ============================================================================
# entry surface — a chat widget on the intake agent
# ============================================================================
# The user-facing way in. Everything downstream (orchestrator, workers) runs as
# background tasks; this is just the front door.

resource "chatbotkit_widget_integration" "intake_widget" {
  name        = "Deep Researcher"
  description = "Ask for deep research; the team handles it in the background"
  bot_id      = chatbotkit_bot.intake.id
  intro       = "What would you like me to research?"
  placeholder = "e.g. Compare the leading open-source web-research agents"
}

# ============================================================================
# Outputs
# ============================================================================

output "intake_bot_id" {
  description = "The ID of the intake (front-door) agent"
  value       = chatbotkit_bot.intake.id
}

output "orchestrator_bot_id" {
  description = "The ID of the orchestrator (big-brain) agent"
  value       = chatbotkit_bot.orchestrator.id
}

output "researcher_bot_id" {
  description = "The ID of the researcher (worker) agent"
  value       = chatbotkit_bot.researcher.id
}

output "workspace_id" {
  description = "The ID of the shared findings/report workspace"
  value       = chatbotkit_space.workspace.id
}

output "widget_integration_id" {
  description = "The ID of the intake chat widget"
  value       = chatbotkit_widget_integration.intake_widget.id
}
