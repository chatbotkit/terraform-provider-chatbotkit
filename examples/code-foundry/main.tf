# Code Foundry: a multi-User architecture.
#
# One shared tools User holds the expensive, sensitive machinery once: a
# GitHub bot that mints repository-scoped App tokens, a Coding Tools skillset
# exposed to child Users as `global-coding-tools`, shared Design/Coding spaces, and
# a Designs Manager that keeps designs in sync. Each customer gets a thin,
# isolated child User with a Coding Agent that installs the shared toolset
# across Users (`@shared@global-coding-tools`) and works on that customer's repo.
#
# Everything runs on one parent User API token (CHATBOTKIT_API_TOKEN). Each child User
# is selected with a provider alias and `run_as` (the X-RunAs-UserId header),
# using the same multi-User pattern as the multi-tenant examples.
#
#   shared User (alias "shared")          child Users
#   ┌───────────────────────────┐          ┌─────────────────────────┐
#   │ GitHub bot (mint tokens)  │   install│ alice: Coding Agent     │
#   │ Coding Tools              │──────────▶  + heartbeat + context  │
#   │   = global-coding-tools   │  @shared@│ bob:   Coding Agent     │
#   │ Design / Coding spaces    │          │  + heartbeat + context  │
#   │ Designs Manager + Sync    │          └─────────────────────────┘
#   └───────────────────────────┘
#
# IMPORTANT: the shared User MUST have the alias `shared`, since child Users
# reference its skillset as `@shared@global-coding-tools`.

terraform {
  required_providers {
    chatbotkit = {
      source = "chatbotkit/chatbotkit"
    }
  }
}

variable "shared_account_id" {
  description = "The shared tools User ID; it must have alias `shared`"
  type        = string
}

variable "github_app_id" {
  description = "The GitHub App ID (JWT iss claim)"
  type        = string
}

variable "github_app_private_key" {
  description = "The GitHub App private key (PEM)"
  type        = string
  sensitive   = true
}

variable "git_email" {
  description = "The git author email the coding agents commit with"
  type        = string
}

variable "alice_account_id" {
  description = "Alice's child User ID"
  type        = string
}

variable "bob_account_id" {
  description = "Bob's child User ID"
  type        = string
}

# One parent User API token via CHATBOTKIT_API_TOKEN; each alias selects a User.

provider "chatbotkit" {
  alias  = "shared"
  run_as = var.shared_account_id
}

provider "chatbotkit" {
  alias  = "alice"
  run_as = var.alice_account_id
}

provider "chatbotkit" {
  alias  = "bob"
  run_as = var.bob_account_id
}

# ============================================================================
# the shared toolbox (one User)
# ============================================================================

module "shared" {
  source    = "./modules/shared"
  providers = { chatbotkit = chatbotkit.shared }

  github_app_id          = var.github_app_id
  github_app_private_key = var.github_app_private_key
}

# ============================================================================
# per-customer coding agents (one child User each, same module)
# ============================================================================
# Each installs `@shared@global-coding-tools` at runtime, so depend on the shared
# shared User being provisioned first.
#
# @note repo_name (and vercel_project_id) are UUID-ish here so each customer's project
# is unique and non-guessable. They are placeholders for real identifiers. In a
# fuller setup Terraform could PROVISION these too and feed the results straight in,
# e.g. the `integrations/github` provider to create the repo and the `vercel`
# provider to create the project:
#
#     repo_name         = github_repository.alice.name
#     vercel_project_id = vercel_project.alice.id
#
# so the GitHub repo, the Vercel project, and the agent's context are all stood up
# in a single apply rather than wired by hand.

module "alice" {
  source    = "./modules/coder"
  providers = { chatbotkit = chatbotkit.alice }

  user_name         = "Alice"
  repo_owner        = "alice-co"
  repo_name         = "8f2a1c4e-3b6d-4a9f-bc12-7e0d5a3f9c81"
  vercel_project_id = "prj_8f2a1c4e3b6d4a9f"
  git_email         = var.git_email

  depends_on = [module.shared]
}

module "bob" {
  source    = "./modules/coder"
  providers = { chatbotkit = chatbotkit.bob }

  user_name         = "Bob"
  repo_owner        = "bob-labs"
  repo_name         = "b71e9d05-2c4a-4f83-9a6e-1d8c3b5f0a27"
  vercel_project_id = "prj_b71e9d052c4a4f83"
  git_email         = var.git_email

  depends_on = [module.shared]
}

# ============================================================================
# Outputs
# ============================================================================

output "shared_github_bot_id" {
  description = "The shared GitHub token-minting bot"
  value       = module.shared.github_bot_id
}

output "coding_tools_alias" {
  description = "Cross-User reference for the shared coding toolset"
  value       = module.shared.coding_tools_alias
}

output "coding_agents" {
  description = "Coding agent bot IDs per child User"
  value = {
    alice = module.alice.bot_id
    bob   = module.bob.bot_id
  }
}
