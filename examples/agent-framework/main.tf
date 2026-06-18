# Agent Framework on ChatBotKit — an agent as a project of files, built with Terraform
#
# This example shows how to define an autonomous agent as a small *project of
# files* and provision the whole thing — the agent, its tools, its workspace, its
# skills, its channels, and its schedules — end to end with Terraform. It is a
# reference architecture: the ChatBotKit backend provides the primitives, and
# Terraform wires them together.
#
# How the pieces map:
#   instructions.md    -> chatbotkit_bot.backstory (the always-on system prompt)
#   abilities          -> chatbotkit_skillset + chatbotkit_skillset_ability,
#                         using ability packs (one ability installs many tools):
#                         shell tools + space-skill tools, plus web search/fetch
#   workspace          -> chatbotkit_space (an isolated, persistent filesystem)
#   skills/*/SKILL.md  -> uploaded into the workspace under .skills/ and read on
#                         demand with the space-skill tools
#   heartbeat.md       -> the description of a frequent "heartbeat" trigger
#   channels           -> *_integration resources (Slack active; others ready)
#   schedules          -> chatbotkit_trigger_integration (cron) + the heartbeat
#
# The agent project lives in ./agent. Editing those files and re-applying updates
# the live agent — the files are the source of truth.
#
# Prerequisites:
# - Set the CHATBOTKIT_API_KEY environment variable
# - (Optional) Provide Slack credentials via variables to activate the channel

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
# workspace  — an isolated, persistent filesystem
# ============================================================================
# The agent's tools run inside this space. It also holds the uploaded skills,
# under .skills/, and any state the agent writes for itself.

resource "chatbotkit_space" "workspace" {
  name        = "Atlas Workspace"
  description = "Isolated workspace where the agent runs commands and stores state"
}

# ============================================================================
# abilities  — the agent's toolset
# ============================================================================
# Abilities are defined in Terraform and grouped into one skillset attached to
# the agent. Two of them are ability "packs" — a single ability that installs a
# whole set of tools:
#   - pack/shell           -> run commands, read/write files, import URLs
#   - pack/cbk/space/skills -> list and read skills from the workspace (.skills/)
# Plus standalone web search and fetch.

resource "chatbotkit_skillset" "abilities" {
  name        = "Atlas Abilities"
  description = "The agent's common toolset"
}

# Shell tools, scoped to the workspace (exec, read/write, import).
resource "chatbotkit_skillset_ability" "shell" {
  skillset_id = chatbotkit_skillset.abilities.id
  space_id    = chatbotkit_space.workspace.id
  name        = "shell"
  description = "Run commands, read/write files, and import URLs in the workspace"
  instruction = <<-EOT
    template: pack/shell
    parameters: {}
  EOT
}

# Space-skill tools: discover and read SKILL.md files from the workspace.
resource "chatbotkit_skillset_ability" "skills" {
  skillset_id = chatbotkit_skillset.abilities.id
  space_id    = chatbotkit_space.workspace.id
  name        = "skills"
  description = "List and read the agent's skills from the workspace (.skills/)"
  instruction = <<-EOT
    template: pack/cbk/space/skills
    parameters: {}
  EOT
}

# Search the web.
resource "chatbotkit_skillset_ability" "web_search" {
  skillset_id = chatbotkit_skillset.abilities.id
  name        = "search_web"
  description = "Search the web for current information"
  instruction = <<-EOT
    ```search
    query: $[query! ys|the search query to find information on the web]
    ```
  EOT
}

# Fetch and read a web page.
resource "chatbotkit_skillset_ability" "web_fetch" {
  skillset_id = chatbotkit_skillset.abilities.id
  name        = "fetch_url"
  description = "Fetch and read the content of a web page"
  instruction = <<-EOT
    ```fetch
    url: $[url! ys|the URL of the web page to fetch and read]
    ```
  EOT
}

# ============================================================================
# the agent  — instructions.md is the backstory
# ============================================================================
# The bot ties everything together: its backstory is the always-on system prompt
# read straight from the project's instructions.md, and its skillset is the
# toolset above.

resource "chatbotkit_bot" "atlas" {
  name        = "Atlas"
  description = "An autonomous agent defined as a project of files and deployed with Terraform"
  model       = "claude-4.5-opus"

  backstory   = file("${path.module}/agent/instructions.md")
  skillset_id = chatbotkit_skillset.abilities.id
}

# ============================================================================
# skills  — uploaded into the workspace under .skills/
# ============================================================================
# Each skill in ./agent/skills is a folder with a SKILL.md (and optional
# scripts). The whole tree is uploaded into the workspace under .skills/, where
# the agent discovers and reads it on demand with the space-skill tools.

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

# ============================================================================
# channels  — the surfaces the agent runs on
# ============================================================================
# Slack is active. The other messaging channels are ready to go — uncomment a
# block (and its variables, below) and supply credentials to activate it.

resource "chatbotkit_slack_integration" "slack" {
  name           = "Atlas for Slack"
  description    = "Atlas in your Slack workspace"
  bot_id         = chatbotkit_bot.atlas.id
  bot_token      = var.slack_bot_token
  signing_secret = var.slack_signing_secret
}

# resource "chatbotkit_discord_integration" "discord" {
#   name        = "Atlas for Discord"
#   description = "Atlas in your Discord server"
#   bot_id      = chatbotkit_bot.atlas.id
#   bot_token   = var.discord_bot_token
#   app_id      = var.discord_app_id
#   public_key  = var.discord_public_key
# }

# resource "chatbotkit_microsoftteams_integration" "teams" {
#   name                     = "Atlas for Teams"
#   description              = "Atlas in Microsoft Teams"
#   bot_id                   = chatbotkit_bot.atlas.id
#   bot_framework_app_id     = var.teams_app_id
#   bot_framework_app_secret = var.teams_app_secret
#   tenant_id                = var.azure_tenant_id
# }

# resource "chatbotkit_whatsapp_integration" "whatsapp" {
#   name            = "Atlas for WhatsApp"
#   description     = "Atlas on WhatsApp Business"
#   bot_id          = chatbotkit_bot.atlas.id
#   access_token    = var.whatsapp_access_token
#   phone_number_id = var.whatsapp_phone_number_id
# }

# ============================================================================
# schedules  — recurring runs
# ============================================================================
# Two kinds:
#   1) Normal schedules: run the agent at specific times for a defined job.
#   2) Heartbeat: a frequent "tick" whose instructions are the heartbeat.md file,
#      passed straight through as the trigger description.

resource "chatbotkit_trigger_integration" "daily_briefing" {
  name             = "Daily Briefing"
  description      = "Weekday morning briefing produced by the agent"
  bot_id           = chatbotkit_bot.atlas.id
  schedule         = "0 9 * * 1-5" # 09:00 on weekdays
  session_duration = 600000        # 10 minutes
  authenticate     = true
}

resource "chatbotkit_trigger_integration" "weekly_review" {
  name             = "Weekly Review"
  description      = "Monday morning review of the past week"
  bot_id           = chatbotkit_bot.atlas.id
  schedule         = "0 8 * * 1" # 08:00 on Mondays
  session_duration = 900000      # 15 minutes
}

resource "chatbotkit_trigger_integration" "heartbeat" {
  name             = "Heartbeat"
  description      = file("${path.module}/agent/heartbeat.md")
  bot_id           = chatbotkit_bot.atlas.id
  schedule         = "*/5 * * * *" # every 5 minutes
  session_duration = 3600000       # 1 hour: ticks within the hour reuse one conversation
}

# ============================================================================
# Variables  — channel credentials
# ============================================================================

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

# Variables for the other channels. Uncomment alongside the matching integration.
#
# variable "discord_bot_token" {
#   description = "Discord Bot Token"
#   type        = string
#   sensitive   = true
#   default     = ""
# }
#
# variable "discord_app_id" {
#   description = "Discord Application ID"
#   type        = string
#   default     = ""
# }
#
# variable "discord_public_key" {
#   description = "Discord Public Key"
#   type        = string
#   default     = ""
# }
#
# variable "teams_app_id" {
#   description = "Microsoft Bot Framework application ID"
#   type        = string
#   default     = ""
# }
#
# variable "teams_app_secret" {
#   description = "Microsoft Bot Framework application secret"
#   type        = string
#   sensitive   = true
#   default     = ""
# }
#
# variable "azure_tenant_id" {
#   description = "Azure AD tenant ID"
#   type        = string
#   default     = ""
# }
#
# variable "whatsapp_access_token" {
#   description = "WhatsApp Business API access token"
#   type        = string
#   sensitive   = true
#   default     = ""
# }
#
# variable "whatsapp_phone_number_id" {
#   description = "WhatsApp Business phone number ID"
#   type        = string
#   default     = ""
# }

# ============================================================================
# Outputs
# ============================================================================

output "agent_bot_id" {
  description = "The ID of the agent"
  value       = chatbotkit_bot.atlas.id
}

output "abilities_skillset_id" {
  description = "The ID of the abilities skillset"
  value       = chatbotkit_skillset.abilities.id
}

output "workspace_id" {
  description = "The ID of the workspace space"
  value       = chatbotkit_space.workspace.id
}

output "skill_paths" {
  description = "Skill files uploaded into the workspace under .skills/"
  value       = [for f in keys(chatbotkit_space_storage_file.skill) : ".skills/${f}"]
}

output "slack_integration_id" {
  description = "The ID of the Slack channel"
  value       = chatbotkit_slack_integration.slack.id
}

output "heartbeat_trigger_id" {
  description = "The ID of the heartbeat trigger"
  value       = chatbotkit_trigger_integration.heartbeat.id
}
