# The SHARED "tools" account.
#
# This is the master toolbox every coding sub-account borrows from. It holds:
#   - a GitHub bot that mints short-lived, repository-scoped GitHub App tokens
#   - a "Coding Tools" skillset, exported account-wide as `global-coding-tools`
#     (visibility = protected, alias = global-coding-tools) so sub-accounts can
#     install it cross-account with `@shared@global-coding-tools`
#   - shared Design and Coding spaces (design files + coding skills)
#   - a Designs Manager bot + Sync trigger that keep the Design space in sync
#
# This module is applied via a provider alias whose `run_as` targets the shared
# account. That account MUST have the alias `shared` (set on the partner user),
# because sub-accounts reference its skillset as `@shared@global-coding-tools`.

terraform {
  required_providers {
    chatbotkit = {
      source = "chatbotkit/chatbotkit"
    }
  }
}

variable "github_app_id" {
  description = "The GitHub App ID (used as the JWT `iss` claim)"
  type        = string
}

variable "github_app_private_key" {
  description = "The GitHub App private key (PEM) used to sign App JWTs"
  type        = string
  sensitive   = true
}

# ============================================================================
# GitHub — the bot that mints repository-scoped tokens
# ============================================================================

resource "chatbotkit_secret" "github_app" {
  name        = "GitHub App Private Key"
  description = "The private key for signing GitHub App JWTs."
  type        = "jwt"
  kind        = "shared"
  value       = var.github_app_private_key

  # @note config is map(string); the structured JWT claims are JSON-encoded.
  config = {
    algorithm        = "RS256"
    schema           = "Bearer"
    expiresInSeconds = "600"
    claims           = jsonencode({ iss = var.github_app_id })
  }
}

resource "chatbotkit_skillset" "github_tools" {
  name        = "GitHub Tools"
  description = "Mint GitHub App installation tokens"
}

resource "chatbotkit_skillset_ability" "create_github_repository_token" {
  skillset_id = chatbotkit_skillset.github_tools.id
  secret_id   = chatbotkit_secret.github_app.id
  name        = "Create GitHub Repository Token"
  description = "Mint a GitHub App installation access token scoped to a specific repository."
  instruction = <<-EOT
    template: "github/repository/token/create"
    parameters:
      owner: !string
        name: "owner"
        description: "the repository owner"
        optional: false
        placeholder: true
      repo: !string
        name: "repo"
        description: "the repository name"
        optional: false
        placeholder: true
  EOT
}

resource "chatbotkit_bot" "github" {
  name        = "GitHub"
  description = "Mints repository-scoped GitHub App tokens on demand"
  skillset_id = chatbotkit_skillset.github_tools.id
}

# ============================================================================
# Coding Tools — the exported, account-wide skillset (global-coding-tools)
# ============================================================================
# This is what sub-accounts install. visibility=protected + a stable alias make
# it referenceable cross-account as `@shared@global-coding-tools`.

resource "chatbotkit_skillset" "coding_tools" {
  name        = "Coding Tools"
  visibility  = "protected"
  alias       = "global-coding-tools"
  description = "The shared coding toolset borrowed by every coding sub-account"
}

# Shell tools (ephemeral per-conversation sandbox).
resource "chatbotkit_skillset_ability" "bash" {
  skillset_id = chatbotkit_skillset.coding_tools.id
  name        = "Bash"
  description = "Execute a shell command or script"
  instruction = "template: \"shell/exec\""
}

resource "chatbotkit_skillset_ability" "rw" {
  skillset_id = chatbotkit_skillset.coding_tools.id
  name        = "RW"
  description = "Read or write file content in the shell environment (mode read|write, optional line ranges)."
  instruction = "template: \"shell/rw\""
}

resource "chatbotkit_skillset_ability" "import_url" {
  skillset_id = chatbotkit_skillset.coding_tools.id
  name        = "Import URL to Shell Environment"
  description = "Import data from a URL into a file in the shell environment."
  instruction = "template: \"shell/import\""
}

# Mint a repo token by delegating to the GitHub bot (which reads the caller's
# context to learn which repository this is about — see the README).
resource "chatbotkit_skillset_ability" "mint_github_repo_token" {
  skillset_id = chatbotkit_skillset.coding_tools.id
  bot_id      = chatbotkit_bot.github.id
  name        = "Mint Github Repo Token"
  description = "Generates a new GitHub token for the current repository. Returns repo, token and expiry."
  instruction = <<-EOT
    template: "bot/apply"
    parameters:
      intent: !string
        name: "intent"
        description: "the configured intent for the bot to apply"
        optional: false
        placeholder: true
  EOT
}

# ============================================================================
# Design space — shared design system files
# ============================================================================

resource "chatbotkit_space" "design" {
  name        = "Design Space"
  description = "Design system files synced from upstream"
}

resource "chatbotkit_skillset_ability" "list_design_files" {
  skillset_id = chatbotkit_skillset.coding_tools.id
  space_id    = chatbotkit_space.design.id
  name        = "List Design Files"
  description = "List the available design files with concise descriptions."
  instruction = "template: \"space/storage/list\""
}

resource "chatbotkit_skillset_ability" "read_design_file" {
  skillset_id = chatbotkit_skillset.coding_tools.id
  space_id    = chatbotkit_space.design.id
  name        = "Read Design File"
  description = "Read a design file containing detailed design-system information."
  instruction = "template: \"space/storage/read\""
}

# ============================================================================
# Coding space — shared coding skills (how-to playbooks)
# ============================================================================

resource "chatbotkit_space" "coding" {
  name        = "Coding"
  description = "Coding skills that explain how to perform coding tasks"
}

resource "chatbotkit_skillset_ability" "list_coding_skills" {
  skillset_id = chatbotkit_skillset.coding_tools.id
  space_id    = chatbotkit_space.coding.id
  name        = "List Coding Skills"
  description = "List coding skills that explain how to perform various coding tasks."
  instruction = "template: \"space/skill/list\""
}

resource "chatbotkit_skillset_ability" "read_coding_skills" {
  skillset_id = chatbotkit_skillset.coding_tools.id
  space_id    = chatbotkit_space.coding.id
  name        = "Read Coding Skills"
  description = "Retrieve a specific coding skill with detailed instructions."
  instruction = "template: \"space/skill/read\""
}

# ============================================================================
# Designs Manager — keeps the Design space in sync with upstream
# ============================================================================

resource "chatbotkit_skillset" "designs_manager_tools" {
  name        = "Designs Manager Tools"
  description = "Shell access to the Design space"
}

resource "chatbotkit_skillset_ability" "designs_manager_bash" {
  skillset_id = chatbotkit_skillset.designs_manager_tools.id
  space_id    = chatbotkit_space.design.id
  name        = "Bash"
  description = "Execute a shell command or script in the Design space"
  instruction = "template: \"shell/exec\""
}

resource "chatbotkit_bot" "designs_manager" {
  name        = "Designs Manager"
  model       = "kimi-k2.6"
  skillset_id = chatbotkit_skillset.designs_manager_tools.id
  backstory   = <<-EOT
    Your role is to keep our /space folder of designs in sync with upstream systems.

    Always clone outside of the /space folder. Once cloned update the files accordingly.
  EOT
}

resource "chatbotkit_trigger_integration" "sync" {
  name         = "Sync"
  description  = "Synchronise https://github.com/nexu-io/open-design/tree/main/design-systems by copying all files into /space/ preserving the upstream folder structure."
  bot_id       = chatbotkit_bot.designs_manager.id
  schedule     = "never"
  authenticate = true
}

# ============================================================================
# Outputs
# ============================================================================

output "github_bot_id" {
  description = "The shared GitHub token-minting bot"
  value       = chatbotkit_bot.github.id
}

output "coding_tools_alias" {
  description = "The cross-account reference for the exported coding toolset"
  value       = "@shared@${chatbotkit_skillset.coding_tools.alias}"
}

output "coding_tools_skillset_id" {
  description = "The Coding Tools skillset ID"
  value       = chatbotkit_skillset.coding_tools.id
}
