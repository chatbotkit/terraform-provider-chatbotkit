# Community Scout — a product-led-growth agent that monitors Reddit for
# conversations worth contributing to, drafts a helpful reply, and suggests it to
# the team on Slack for a human to actually post.
#
# It reuses the "agent as a project of files" shape (instructions + skills +
# scripts) and the monitor-on-a-cycle pattern, applied to social listening:
#   - Deterministic (script) -> fetch the Reddit firehose and dedup it. Doing this in
#                           a script keeps the bulk of the data out of the model's
#                           context (cheaper, faster) — the agent only sees a summary.
#   - Tool calls (read)   -> the agent reads a *specific* thread it is judging with
#                           the read-only Reddit tools (small, on-demand).
#   - Judgment            -> SKILL.md playbooks decide what is worth engaging and
#                           what a genuinely helpful, disclosed reply says.
#   - Hand-off            -> the scout suggests the thread + draft to the team on
#                           Slack (slack/conversation/start); a human posts it.
#
# The gate is structural, not just a prompt: the scout has read-only Reddit tools
# and a Slack-suggest tool — and no Reddit-post tool at all. It physically cannot
# reply on Reddit. A human reads the Slack suggestion and posts if they choose.
#
# How the pieces map:
#   instructions.md         -> chatbotkit_bot.backstory
#   watchlist.md            -> uploaded to the workspace; the engagement + team config
#   agent/skills/*/SKILL.md -> uploaded under .skills/, read on demand
#   monitor scripts         -> deterministic dedup of discovered threads
#   workspace (space)       -> mentions/ (store), knowledge/ (targeting notes)
#   heartbeat.md            -> the cycle tick (the trigger description)
#
# Prerequisites:
# - Set the CHATBOTKIT_API_KEY environment variable
# - Provide Slack credentials (variables below) so the scout can suggest to the team

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
# workspace  — the mention store, knowledge base, scripts, and skills
# ============================================================================

resource "chatbotkit_space" "workspace" {
  name        = "Community Scout Workspace"
  description = "Mention store (mentions/), knowledge (knowledge/), watchlist, scripts, and skills"
}

# ============================================================================
# the scout's toolset
# ============================================================================

resource "chatbotkit_skillset" "scout_tools" {
  name        = "Community Scout Tools"
  description = "Read-only Reddit, Slack-suggest, shell, skills, and web tools for a community scout"
}

# Read-only Reddit tools for reading a *specific* thread the agent is judging (its
# comments and context). Bulk monitoring is the fetch_mentions.py script's job, not
# this — so the firehose never round-trips through the model. No auth required.
resource "chatbotkit_skillset_ability" "reddit_browse" {
  skillset_id = chatbotkit_skillset.scout_tools.id
  name        = "reddit_browse"
  description = "Read a specific Reddit thread's comments and context during scoring/drafting"
  instruction = <<-EOT
    template: pack/reddit[read-only]
    parameters: {}
  EOT
}

# Suggest a thread to the team on Slack. This is the scout's only outbound action
# — there is deliberately no Reddit-post tool, so it cannot reply on Reddit; it
# hands the draft to a human. The Slack integration id is baked in; the agent
# supplies the channel and the message.
resource "chatbotkit_skillset_ability" "suggest_to_team" {
  skillset_id = chatbotkit_skillset.scout_tools.id
  name        = "suggest_to_team"
  description = "Send a suggestion to the team's Slack channel: a thread worth replying to, why, a suggested draft, and who might take it"
  instruction = <<-EOT
    template: slack/conversation/start[by-id]
    parameters:
      slackIntegrationId: ${chatbotkit_slack_integration.team.id}
  EOT
}

# Shell scoped to the workspace: run the dedup script, read/write mention files.
resource "chatbotkit_skillset_ability" "shell" {
  skillset_id = chatbotkit_skillset.scout_tools.id
  space_id    = chatbotkit_space.workspace.id
  name        = "shell"
  description = "Run scripts and read/write mention, knowledge, and config files"
  instruction = <<-EOT
    template: pack/shell
    parameters: {}
  EOT
}

# Discover and read SKILL.md playbooks from the workspace.
resource "chatbotkit_skillset_ability" "skills" {
  skillset_id = chatbotkit_skillset.scout_tools.id
  space_id    = chatbotkit_space.workspace.id
  name        = "skills"
  description = "List and read the scout's skills from the workspace (.skills/)"
  instruction = <<-EOT
    template: pack/cbk/space/skills
    parameters: {}
  EOT
}

# Web search/fetch for product-fit research (and an X fallback via site:x.com).
resource "chatbotkit_skillset_ability" "search_web" {
  skillset_id = chatbotkit_skillset.scout_tools.id
  name        = "search_web"
  description = "Search the web for context (and as an X/other-source fallback)"
  instruction = <<-EOT
    template: search/web
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "fetch_url" {
  skillset_id = chatbotkit_skillset.scout_tools.id
  name        = "fetch_url"
  description = "Fetch and read the content of a web page"
  instruction = <<-EOT
    template: fetch/text/get
    parameters: {}
  EOT
}

# --- Optional: X / Twitter ------------------------------------------------
# X has no first-class catalogue. To add it, create a secret with an X API
# bearer token and a fetch/request[with-auth] ability against the recent-search
# endpoint, writing mentions in the same shape (mentions/<id>.json). Until then,
# the agent can approximate X via search_web with site:x.com queries.
#
# resource "chatbotkit_secret" "x_api" {
#   name  = "X API Bearer Token"
#   type  = "bearer"
#   value = var.x_api_token
# }
# resource "chatbotkit_skillset_ability" "x_search" {
#   skillset_id = chatbotkit_skillset.scout_tools.id
#   secret_id   = chatbotkit_secret.x_api.id
#   name        = "x_search"
#   description = "Search recent X posts (requires an X API bearer token)"
#   instruction = <<-EOT
#     template: fetch/request[with-auth]
#     parameters: {}
#   EOT
# }

# ============================================================================
# the scout  — instructions.md is the backstory
# ============================================================================

resource "chatbotkit_bot" "scout" {
  name        = "Community Scout"
  description = "Monitors Reddit for relevant conversations, drafts helpful replies, and suggests them to the team on Slack (never posts to Reddit itself)"
  model       = "deepseek-v4-pro"

  backstory   = file("${path.module}/agent/instructions.md")
  skillset_id = chatbotkit_skillset.scout_tools.id
}

# ============================================================================
# Slack  — where the scout suggests threads to the team
# ============================================================================
# The scout starts a Slack conversation (slack/conversation/start) to hand a
# thread + draft to a human. This is the human-in-the-loop for an automated agent.

resource "chatbotkit_slack_integration" "team" {
  name           = "Community Scout for Slack"
  description    = "Where the scout posts suggestions for the team to act on"
  bot_id         = chatbotkit_bot.scout.id
  bot_token      = var.slack_bot_token
  signing_secret = var.slack_signing_secret
}

variable "slack_bot_token" {
  description = "Slack Bot Token (xoxb-...)"
  type        = string
  sensitive   = true
  default     = ""
}

variable "slack_signing_secret" {
  description = "Slack Signing Secret"
  type        = string
  sensitive   = true
  default     = ""
}

# ============================================================================
# upload the skills tree and the watchlist config into the workspace
# ============================================================================

locals {
  skills_dir = "${path.module}/agent/skills"
}

resource "chatbotkit_space_storage_file" "skill" {
  for_each = fileset(local.skills_dir, "**")

  space_id    = chatbotkit_space.workspace.id
  path        = ".skills/${each.value}"
  source      = "${local.skills_dir}/${each.value}"
  source_hash = filesha256("${local.skills_dir}/${each.value}")
}

# @note the watchlist is a *living* config: the team (or the agent) tunes it in the
# space over time. Terraform seeds it once and then leaves it alone — ignore_changes
# stops a later apply from reverting live edits back to this committed copy. To push
# a new baseline, edit it directly in the space (or remove ignore_changes for one
# apply). The skills above are code, so they stay Terraform-managed (no ignore).
resource "chatbotkit_space_storage_file" "watchlist" {
  space_id    = chatbotkit_space.workspace.id
  path        = "watchlist.md"
  source      = "${path.module}/agent/watchlist.md"
  source_hash = filesha256("${path.module}/agent/watchlist.md")

  lifecycle {
    ignore_changes = [source, source_hash]
  }
}

# ============================================================================
# the cycle  — recurring scouting, plus a daily digest to the team
# ============================================================================

resource "chatbotkit_trigger_integration" "cycle" {
  name             = "Scouting Cycle"
  description      = file("${path.module}/agent/heartbeat.md")
  bot_id           = chatbotkit_bot.scout.id
  schedule         = "0 * * * *" # hourly
  session_duration = 3600000     # 1 hour: ticks within the hour reuse one conversation
}

resource "chatbotkit_trigger_integration" "daily_digest" {
  name             = "Daily Digest"
  description      = "Send the team a Slack digest of the day's suggested threads and what came of them, and consolidate knowledge from outcomes."
  bot_id           = chatbotkit_bot.scout.id
  schedule         = "0 17 * * 1-5" # 17:00 on weekdays
  session_duration = 900000         # 15 minutes
}

# ============================================================================
# Outputs
# ============================================================================

output "scout_bot_id" {
  description = "The ID of the community scout agent"
  value       = chatbotkit_bot.scout.id
}

output "workspace_id" {
  description = "The ID of the workspace (mention store + knowledge base)"
  value       = chatbotkit_space.workspace.id
}

output "skill_paths" {
  description = "Skill files uploaded into the workspace under .skills/"
  value       = [for f in keys(chatbotkit_space_storage_file.skill) : ".skills/${f}"]
}

output "cycle_trigger_id" {
  description = "The ID of the recurring scouting cycle trigger"
  value       = chatbotkit_trigger_integration.cycle.id
}

output "slack_integration_id" {
  description = "The ID of the Slack integration the scout suggests to"
  value       = chatbotkit_slack_integration.team.id
}
