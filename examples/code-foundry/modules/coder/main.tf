# A per-customer child User.
#
# Each customer gets an isolated child User containing a Coding Agent and one
# ability that installs the shared toolset across Users. The heavy tooling lives
# once in the shared User; this module borrows it.
#
# Applied through a provider alias whose `run_as` targets this child User.
#
# Two things make the setup work end to end:
#   1. install the shared coding skillset (`@shared@global-coding-tools`)
#   2. set this child User's CONTEXT (which repo + which Vercel project) so the shared
#      GitHub bot knows which repository to mint a token for; hard-coding the
#      repo would be a security risk because each agent belongs to a child User.

terraform {
  required_providers {
    chatbotkit = {
      source = "chatbotkit/chatbotkit"
    }
  }
}

variable "user_name" {
  description = "Display name of the customer represented by this child User"
  type        = string
}

variable "repo_owner" {
  description = "The GitHub repository owner for this child User's project"
  type        = string
}

variable "repo_name" {
  description = "The GitHub repository name for this child User's project"
  type        = string
}

variable "vercel_project_id" {
  description = "The Vercel project ID this child User's agent deploys to"
  type        = string
  default     = ""
}

variable "git_email" {
  description = "The git author email the agent commits with"
  type        = string
}

variable "model" {
  description = "The model the coding agent should use"
  type        = string
  default     = "glm-5.2"
}

variable "heartbeat_schedule" {
  description = "Cron schedule for the heartbeat that keeps the agent working a task"
  type        = string
  default     = "0 9-17 * * *" # every hour, 09:00–17:00 (business hours)
}

variable "heartbeat_timezone" {
  description = "IANA timezone the heartbeat schedule is evaluated in"
  type        = string
  default     = "Europe/London"
}

# ============================================================================
# the coding agent
# ============================================================================

resource "chatbotkit_skillset" "coder_tools" {
  name        = "Coder Tools"
  description = "Bootstrap toolset: installs the shared coding skillset on demand"
}

# The only ability the agent ships with: install the shared toolset across Users.
# `@shared@global-coding-tools` resolves to the Coding Tools skillset in the
# User aliased `shared`.
resource "chatbotkit_skillset_ability" "install_coding_skillset" {
  skillset_id = chatbotkit_skillset.coder_tools.id
  name        = "Install Coding Skillset"
  description = "Enhance your coding capabilities by installing the shared coding skills."
  instruction = <<-EOT
    template: conversation/skillset/install[by-id]
    parameters:
      skillsetId: '@shared@global-coding-tools'
  EOT
}

resource "chatbotkit_bot" "coder" {
  name        = "Coding Agent"
  description = "A coding agent that knows how to build landing pages."
  model       = var.model
  skillset_id = chatbotkit_skillset.coder_tools.id
  backstory   = <<-EOT
    A coding agent that builds landing pages with nextjs.

    You must install the relevant skills / tools to obtain specific platform capabilities.

    For access to the repo use the github tools to mint a token and get the repo access.
    The repository and Vercel project for this User come from your context. Do
    not assume or hard-code them.

    Then you must use the provided shell environment to perform the actions.

    NB. use ${var.git_email} for the github email.

    NB. treat task functions like system notifications... all tools are already
    available, you just need to load them.
  EOT
}

# ============================================================================
# heartbeat - keep working the active task
# ============================================================================
# Coding tasks span many steps. This recurring tick nudges the agent to make the
# next concrete step on whatever it is currently working on, reusing one
# conversation within the session window so it keeps its place across ticks.

resource "chatbotkit_trigger_integration" "heartbeat" {
  name             = "Coding Heartbeat"
  description      = <<-EOT
    Continue the current coding task. Review your workspace and the repo for
    in-progress work and make the next concrete step toward completion (commit and
    push as you go). If there is no active task, stop and wait. Do not invent work.
  EOT
  bot_id           = chatbotkit_bot.coder.id
  schedule         = var.heartbeat_schedule
  timezone         = var.heartbeat_timezone
  session_duration = 3600000 # 1 hour: ticks within the hour continue one conversation
}

# ============================================================================
# context - which repo and Vercel project this child User's agent works on
# ============================================================================
# The shared GitHub bot reads this to know which repository to scope a token to.
# Created for this child User through the module's provider, so each agent is
# scoped to its own repository. The payload is map(string).

resource "chatbotkit_context" "project" {
  name        = "project"
  description = "Repository and Vercel project for this child User's coding agent"

  payload = {
    githubOwner     = var.repo_owner
    githubRepo      = var.repo_name
    githubRepoUrl   = "https://github.com/${var.repo_owner}/${var.repo_name}"
    vercelProjectId = var.vercel_project_id
  }
}

# ============================================================================
# Outputs
# ============================================================================

output "bot_id" {
  description = "The coding agent bot ID in this child User"
  value       = chatbotkit_bot.coder.id
}
