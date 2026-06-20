# SOC Investigator — an autonomous security-operations agent, built with Terraform
#
# This example models the agentic half of an open-source SOC platform (in the
# shape of projects like agentic-soc-platform): an analyst agent that runs on a
# cycle to pull SIEM alerts, correlate them into cases, triage, investigate the
# ones that matter, enrich indicators, and accumulate knowledge.
#
# The point of the example is the boundary between deterministic and agentic work:
#   - Deterministic  -> scripts run with the shell tools (pull/correlate alerts,
#                       threat-intel lookup). Cheap, reliable, no model.
#   - Judgment       -> SKILL.md playbooks the agent reads and follows (triage,
#                       investigate, extract knowledge). This is where it reasons.
#
# How the pieces map:
#   instructions.md        -> chatbotkit_bot.backstory (the analyst's system prompt)
#   agent/skills/*/SKILL.md -> uploaded under .skills/, read on demand
#   agent/skills/**/scripts -> deterministic Python (stdlib only), run via shell
#   agent/data/*.json      -> a mock SIEM the puller reads (swap for a real SIEM)
#   workspace (space)      -> case store (cases/), knowledge base (knowledge/)
#   heartbeat.md           -> the cycle tick (the trigger description)
#
# Self-contained: the mock SIEM and a file-based case/knowledge store mean it runs
# with no external dependencies. Production seams are marked in the scripts and
# README (real SIEM, a SIRP case DB, a dataset-backed knowledge base).
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
# workspace  — case store, knowledge base, scripts, and skills
# ============================================================================

resource "chatbotkit_space" "workspace" {
  name        = "SOC Workspace"
  description = "Case store (cases/), knowledge base (knowledge/), scripts, and skills"
}

# ============================================================================
# the analyst's toolset
# ============================================================================
# Shell (run the deterministic scripts, read/write the case and knowledge files)
# and space-skill discovery, plus web search/fetch for investigation context.

resource "chatbotkit_skillset" "soc_tools" {
  name        = "SOC Analyst Tools"
  description = "Shell + skills + web tools for an autonomous SOC analyst"
}

# Shell tools scoped to the workspace: run scripts, read/write files.
resource "chatbotkit_skillset_ability" "shell" {
  skillset_id = chatbotkit_skillset.soc_tools.id
  space_id    = chatbotkit_space.workspace.id
  name        = "shell"
  description = "Run scripts and read/write case and knowledge files in the workspace"
  instruction = <<-EOT
    template: pack/shell
    parameters: {}
  EOT
}

# Discover and read SKILL.md playbooks from the workspace (.skills/).
resource "chatbotkit_skillset_ability" "skills" {
  skillset_id = chatbotkit_skillset.soc_tools.id
  space_id    = chatbotkit_space.workspace.id
  name        = "skills"
  description = "List and read the analyst's skills from the workspace (.skills/)"
  instruction = <<-EOT
    template: pack/cbk/space/skills
    parameters: {}
  EOT
}

# Web search for unfamiliar tooling, CVEs, and attacker techniques.
resource "chatbotkit_skillset_ability" "search_web" {
  skillset_id = chatbotkit_skillset.soc_tools.id
  name        = "search_web"
  description = "Search the web for context on tooling, CVEs, and techniques"
  instruction = <<-EOT
    template: search/web
    parameters: {}
  EOT
}

# Fetch and read a page found during investigation.
resource "chatbotkit_skillset_ability" "fetch_url" {
  skillset_id = chatbotkit_skillset.soc_tools.id
  name        = "fetch_url"
  description = "Fetch and read the content of a web page"
  instruction = <<-EOT
    template: fetch/text/get
    parameters: {}
  EOT
}

# ============================================================================
# the analyst  — instructions.md is the backstory
# ============================================================================

resource "chatbotkit_bot" "analyst" {
  name        = "SOC Analyst"
  description = "Autonomous security-operations agent: pull, triage, investigate, enrich, learn"
  model       = "claude-4.5-sonnet" # capable + cost-sensible for a cycle; bump to opus for deeper investigation

  backstory   = file("${path.module}/agent/instructions.md")
  skillset_id = chatbotkit_skillset.soc_tools.id
}

# ============================================================================
# upload the skills tree and the mock SIEM data into the workspace
# ============================================================================

locals {
  skills_dir = "${path.module}/agent/skills"
  data_dir   = "${path.module}/agent/data"
}

resource "chatbotkit_space_storage_file" "skill" {
  for_each = fileset(local.skills_dir, "**")

  space_id    = chatbotkit_space.workspace.id
  path        = ".skills/${each.value}"
  source      = "${local.skills_dir}/${each.value}"
  source_hash = filesha256("${local.skills_dir}/${each.value}")
}

resource "chatbotkit_space_storage_file" "data" {
  for_each = fileset(local.data_dir, "**")

  space_id    = chatbotkit_space.workspace.id
  path        = "data/${each.value}"
  source      = "${local.data_dir}/${each.value}"
  source_hash = filesha256("${local.data_dir}/${each.value}")
}

# ============================================================================
# the cycle  — a recurring trigger that runs one SOC cycle
# ============================================================================
# The heartbeat reuses one conversation within its session window, so a cycle's
# steps share context. Tune the schedule to your alert volume.

resource "chatbotkit_trigger_integration" "cycle" {
  name             = "SOC Cycle"
  description      = file("${path.module}/agent/heartbeat.md")
  bot_id           = chatbotkit_bot.analyst.id
  schedule         = "*/15 * * * *" # every 15 minutes
  session_duration = 3600000        # 1 hour: ticks within the hour reuse one conversation
}

# A slower, deliberate pass for end-of-day review and knowledge consolidation.
resource "chatbotkit_trigger_integration" "daily_review" {
  name             = "SOC Daily Review"
  description      = "Review the day's pending_review cases, summarize the queue for the human team, and consolidate knowledge from approved cases."
  bot_id           = chatbotkit_bot.analyst.id
  schedule         = "0 18 * * 1-5" # 18:00 on weekdays
  session_duration = 900000         # 15 minutes
}

# ============================================================================
# Outputs
# ============================================================================

output "analyst_bot_id" {
  description = "The ID of the SOC analyst agent"
  value       = chatbotkit_bot.analyst.id
}

output "workspace_id" {
  description = "The ID of the workspace (case store + knowledge base)"
  value       = chatbotkit_space.workspace.id
}

output "skill_paths" {
  description = "Skill files uploaded into the workspace under .skills/"
  value       = [for f in keys(chatbotkit_space_storage_file.skill) : ".skills/${f}"]
}

output "cycle_trigger_id" {
  description = "The ID of the recurring SOC cycle trigger"
  value       = chatbotkit_trigger_integration.cycle.id
}
