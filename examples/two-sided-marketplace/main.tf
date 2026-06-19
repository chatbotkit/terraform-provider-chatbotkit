# ============================================================================
# Two-Sided Marketplace — one agent per side of a recruitment-style market
# ============================================================================
#
# This example models a two-sided market the way a recruitment agency works:
# one side has people to place, the other has roles to fill, and the agency
# clears the market between them. It deploys TWO agents, one per side:
#
#   • Talent Partner  (SUPPLY side) — works for the people. Browses open roles
#                       and submits a candidate into the market.
#   • Client Partner  (DEMAND side) — works for the companies. Posts roles,
#                       searches the talent pool, reviews who was submitted, and
#                       advances a submission (the placement decision).
#
# THE CORE IS AN API, NOT THE AGENTS.
# The matching engine — the role index, the talent pool, submissions, the
# placement/fee transaction — lives in your own marketplace service. The agents
# never hold that logic; they reach it with the `fetch` action. This is the
# whole point: ChatBotKit is the conversational + identity tier on each side,
# your API is the system of record that both sides share.
#
# THE MARKET CLEARS THROUGH SHARED STATE, NOT AGENT-TO-AGENT CHATTER.
# Talent Partner POSTs a submission → it lands in your core → it shows up in
# Client Partner's review queue → Client Partner advances it. The two agents
# have opposite loyalties (one to the candidate, one to the employer) and never
# talk directly. They meet only through the marketplace API, and the placement
# is a deterministic server-side transaction — never something a model decides
# unilaterally. (See README "Adapting it" for an optional broker desk-lead.)
#
# IDENTITY: ONE SHARED GATEWAY TOKEN + PER-SIDE ON-BEHALF OAUTH.
#   • A SHARED service token ("Marketplace Service Token") is the platform's
#     machine identity to the marketplace gateway. BOTH agents present it, so
#     the core sees one stable caller no matter who is chatting.
#   • A PERSONAL OAuth secret per side is the *human* the action is attributed
#     to: the candidate on the Talent side, the employer contact on the Client
#     side. Each person signs in once through SSO; the core then enforces what
#     that person may see and do.
#
# HEADER CONVENTION (read this once and every ability reads the same way):
#   Authorization:   !reference SECRET_MARKETPLACE_SERVICE_TOKEN
#                    -> ALWAYS the shared gateway token, referenced by name.
#                       The core authenticates the *caller* (the platform) here.
#   X-On-Behalf-Of:  !reference SECRET_DEFAULT
#                    -> the side's human, present only on actions attributed to a
#                       person. SECRET_DEFAULT = the secret linked to the ability
#                       via `secret_id`, which is THIS desk's identity — the
#                       candidate on the Talent Desk, the employer on the Client
#                       Desk. It never means two different things in one skillset.
#
# An ability links exactly ONE secret via `secret_id` (that is SECRET_DEFAULT).
# To use a *second* secret in the same request, reference it by name
# (SECRET_<NAME>, name upper-cased with non-word characters collapsed to "_").
# That is why the action abilities link the personal identity and name the
# shared token, while the read-only browse abilities link only the shared token.
#
# A `bearer`/`oauth` secret already resolves to a complete Authorization value
# including its scheme (`Bearer <token>`), so never add your own `Bearer ` prefix.
#
# Prerequisites:
# - Set the CHATBOTKIT_API_KEY environment variable
# - Provide the secret/OAuth values in terraform.tfvars (see the .example file)

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
# Secrets
# ============================================================================

# SHARED service token — the platform's machine identity to the marketplace
# gateway. `kind = "shared"` is one account-wide value; every conversation on
# both sides uses it, so the core sees one stable service caller.
resource "chatbotkit_secret" "marketplace_service_token" {
  name        = "Marketplace Service Token"
  description = "Service-account bearer token both desks present to the marketplace gateway"
  type        = "bearer"
  kind        = "shared"

  value = var.marketplace_service_token
}

# PERSONAL OAuth — the CANDIDATE. `kind = "personal"` resolves per end-user:
# the first time a candidate triggers an ability that links this secret, they
# sign in through SSO and the Talent Partner then acts as them.
resource "chatbotkit_secret" "talent_identity" {
  name        = "Talent Identity"
  description = "Per-candidate OAuth token; the Talent Partner acts on behalf of the signed-in candidate"
  type        = "oauth"
  kind        = "personal"

  config = {
    clientId         = var.talent_oauth_client_id
    clientSecret     = var.talent_oauth_client_secret
    authorizationUrl = "https://auth.recruitmarket.example/oauth/authorize"
    tokenUrl         = "https://auth.recruitmarket.example/oauth/token"
    grantType        = "authorization_code"
    scope            = "openid profile roles:read submissions:write"
  }
}

# PERSONAL OAuth — the EMPLOYER contact. Same mechanism, the demand side.
resource "chatbotkit_secret" "client_identity" {
  name        = "Client Identity"
  description = "Per-employer OAuth token; the Client Partner acts on behalf of the signed-in hiring contact"
  type        = "oauth"
  kind        = "personal"

  config = {
    clientId         = var.client_oauth_client_id
    clientSecret     = var.client_oauth_client_secret
    authorizationUrl = "https://auth.recruitmarket.example/oauth/authorize"
    tokenUrl         = "https://auth.recruitmarket.example/oauth/token"
    grantType        = "authorization_code"
    scope            = "openid profile roles:write talent:read submissions:write"
  }
}

# ============================================================================
# SUPPLY SIDE — Talent Desk skillset + abilities
# ============================================================================

resource "chatbotkit_skillset" "talent_desk" {
  name        = "Talent Desk"
  description = "Tools for representing a candidate in the marketplace: browse roles, submit, track"
}

# Browse open roles. Read-only public catalog, so it links ONLY the shared
# service token — no candidate sign-in required just to look.
resource "chatbotkit_skillset_ability" "search_open_roles" {
  skillset_id = chatbotkit_skillset.talent_desk.id
  secret_id   = chatbotkit_secret.marketplace_service_token.id

  name        = "search_open_roles"
  description = "Search currently open roles in the marketplace by keyword, location or seniority"
  instruction = <<-EOT
    !fetch
    method: GET
    url: https://api.recruitmarket.example/v1/roles
    query:
      q: !string
        name: query
        description: keywords, e.g. "senior backend engineer remote"
      location: !string
        name: location
        description: optional location or "remote"
        optional: true
      seniority: !string
        name: seniority
        description: optional seniority filter
        optional: true
        enum:
          - junior
          - mid
          - senior
          - lead
    headers:
      Authorization: !reference SECRET_MARKETPLACE_SERVICE_TOKEN
      Accept: application/json
    options:
      format: json
  EOT
}

# Submit the candidate into a role. This is attributed to the candidate, so the
# ability links the candidate identity (SECRET_DEFAULT) and forwards it on
# behalf, while the shared token gets past the gateway.
resource "chatbotkit_skillset_ability" "submit_candidacy" {
  skillset_id = chatbotkit_skillset.talent_desk.id
  secret_id   = chatbotkit_secret.talent_identity.id

  name        = "submit_candidacy"
  description = "Submit the current candidate into an open role on their behalf"
  instruction = <<-EOT
    !fetch
    method: POST
    url: https://api.recruitmarket.example/v1/submissions
    headers:
      Authorization: !reference SECRET_MARKETPLACE_SERVICE_TOKEN
      X-On-Behalf-Of: !reference SECRET_DEFAULT
      Content-Type: application/json
    body:
      roleId: !string
        name: role_id
        description: the id of the role to apply to (from search_open_roles)
      note: !string
        name: note
        description: a short note to the hiring side on why the candidate fits
        optional: true
    options:
      format: json
  EOT
}

# Track the candidate's own submissions and their status. Attributed to the
# candidate so the core returns only their submissions.
resource "chatbotkit_skillset_ability" "track_my_submissions" {
  skillset_id = chatbotkit_skillset.talent_desk.id
  secret_id   = chatbotkit_secret.talent_identity.id

  name        = "track_my_submissions"
  description = "List the current candidate's submissions and where each one stands"
  instruction = <<-EOT
    !fetch
    method: GET
    url: https://api.recruitmarket.example/v1/submissions
    query:
      stage: !string
        name: stage
        description: optional stage filter; leave empty for all
        optional: true
        enum:
          - submitted
          - reviewing
          - interview
          - offer
          - placed
          - declined
    headers:
      Authorization: !reference SECRET_MARKETPLACE_SERVICE_TOKEN
      X-On-Behalf-Of: !reference SECRET_DEFAULT
      Accept: application/json
    options:
      format: json
  EOT
}

# ============================================================================
# DEMAND SIDE — Client Desk skillset + abilities
# ============================================================================

resource "chatbotkit_skillset" "client_desk" {
  name        = "Client Desk"
  description = "Tools for representing a hiring company: post roles, search talent, review and advance submissions"
}

# Post a new open role, attributed to the employer contact.
resource "chatbotkit_skillset_ability" "post_role" {
  skillset_id = chatbotkit_skillset.client_desk.id
  secret_id   = chatbotkit_secret.client_identity.id

  name        = "post_role"
  description = "Post a new open role into the marketplace on behalf of the employer"
  instruction = <<-EOT
    !fetch
    method: POST
    url: https://api.recruitmarket.example/v1/roles
    headers:
      Authorization: !reference SECRET_MARKETPLACE_SERVICE_TOKEN
      X-On-Behalf-Of: !reference SECRET_DEFAULT
      Content-Type: application/json
    body:
      title: !string
        name: title
        description: the role title, e.g. "Senior Backend Engineer"
      seniority: !string
        name: seniority
        description: the seniority level
        enum:
          - junior
          - mid
          - senior
          - lead
      location: !string
        name: location
        description: location or "remote"
      description: !string
        name: description
        description: a detailed description of the role and requirements
    options:
      format: json
  EOT
}

# Search the talent pool. Attributed to the employer so the core enforces which
# (consented, visible) candidate profiles this employer is allowed to see.
resource "chatbotkit_skillset_ability" "search_talent" {
  skillset_id = chatbotkit_skillset.client_desk.id
  secret_id   = chatbotkit_secret.client_identity.id

  name        = "search_talent"
  description = "Search the marketplace talent pool for candidates matching a role"
  instruction = <<-EOT
    !fetch
    method: GET
    url: https://api.recruitmarket.example/v1/talent
    query:
      q: !string
        name: query
        description: skills/keywords, e.g. "golang kubernetes 5y"
      location: !string
        name: location
        description: optional location or "remote"
        optional: true
    headers:
      Authorization: !reference SECRET_MARKETPLACE_SERVICE_TOKEN
      X-On-Behalf-Of: !reference SECRET_DEFAULT
      Accept: application/json
    options:
      format: json
  EOT
}

# Review who has been submitted to one of the employer's roles. This is the
# other side of submit_candidacy — the market clearing through shared state.
resource "chatbotkit_skillset_ability" "review_submissions" {
  skillset_id = chatbotkit_skillset.client_desk.id
  secret_id   = chatbotkit_secret.client_identity.id

  name        = "review_submissions"
  description = "List candidates submitted to one of the employer's roles"
  instruction = <<-EOT
    !fetch
    method: GET
    url: https://api.recruitmarket.example/v1/submissions
    query:
      roleId: !string
        name: role_id
        description: the employer's role id to review submissions for
    headers:
      Authorization: !reference SECRET_MARKETPLACE_SERVICE_TOKEN
      X-On-Behalf-Of: !reference SECRET_DEFAULT
      Accept: application/json
    options:
      format: json
  EOT
}

# Advance a submission — request an interview, make an offer, or place. This is
# THE market-clearing transaction: the employer accepts what the candidate
# submitted, and the core records the stage change (and any placement fee)
# server-side. The agent requests it; the API decides and records it.
resource "chatbotkit_skillset_ability" "advance_submission" {
  skillset_id = chatbotkit_skillset.client_desk.id
  secret_id   = chatbotkit_secret.client_identity.id

  name        = "advance_submission"
  description = "Advance a submission to the next stage (interview, offer, or placement) on behalf of the employer"
  instruction = <<-EOT
    !fetch
    method: POST
    url: https://api.recruitmarket.example/v1/submissions/advance
    headers:
      Authorization: !reference SECRET_MARKETPLACE_SERVICE_TOKEN
      X-On-Behalf-Of: !reference SECRET_DEFAULT
      Content-Type: application/json
    body:
      submissionId: !string
        name: submission_id
        description: the submission to advance (from review_submissions)
      stage: !string
        name: stage
        description: the stage to move the submission to
        enum:
          - interview
          - offer
          - placed
          - declined
      message: !string
        name: message
        description: a short message to the candidate side about this decision
        optional: true
    options:
      format: json
  EOT
}

# ============================================================================
# Bots — one per side of the market
# ============================================================================

resource "chatbotkit_bot" "talent_partner" {
  name        = "Talent Partner"
  description = "Represents candidates in the marketplace: finds roles and submits on their behalf"
  model       = "claude-4.5-sonnet"

  backstory = <<-EOT
    You are the Talent Partner. You work for the candidate you are talking to —
    their interests come first. You help them find suitable open roles and put
    them forward.

    What you do:
    - Understand what the candidate is looking for, then use the marketplace to
      find roles that genuinely fit. Do not pad results with poor matches.
    - Submit the candidate to a role only with their clear go-ahead, and write
      an honest note about why they fit. Never invent experience.
    - Keep them informed on where each submission stands.

    How identity works: the first time the candidate does something attributed
    to them (submitting, tracking their submissions) they are asked to sign in.
    If you receive an authorization prompt, relay it and wait for them to
    authenticate — never ask for a password directly.

    The hiring side reviews submissions and decides whether to advance them; you
    do not control that. Be honest about it. Never expose raw tokens or internal
    URLs.
  EOT

  skillset_id = chatbotkit_skillset.talent_desk.id
}

resource "chatbotkit_bot" "client_partner" {
  name        = "Client Partner"
  description = "Represents hiring companies: posts roles, sources talent, reviews and advances submissions"
  model       = "claude-4.5-sonnet"

  backstory = <<-EOT
    You are the Client Partner. You work for the hiring company you are talking
    to. You help them fill roles well and fast.

    What you do:
    - Post clear, accurate roles. Search the talent pool for real fits.
    - Review who has been submitted to the company's roles.
    - Advancing a submission (interview, offer, placement) is a real commitment
      to the candidate side. Always confirm with the employer before you advance
      or decline anyone, and pass along a brief, respectful message.

    How identity works: the first time the employer contact does something
    attributed to them they are asked to sign in. If you receive an
    authorization prompt, relay it and wait for them to authenticate — never ask
    for a password directly.

    You only see talent the marketplace allows this employer to see. Never expose
    raw tokens or internal URLs, and never fabricate candidate details.
  EOT

  skillset_id = chatbotkit_skillset.client_desk.id
}

# ============================================================================
# Channels (optional) — each side reaches a different audience
# ============================================================================
# Candidates meet the Talent Partner through a web widget; employer recruiters
# meet the Client Partner in Slack. Leave the Slack vars empty to skip it.

resource "chatbotkit_widget_integration" "talent_widget" {
  name        = "Talent Partner Widget"
  description = "Candidate-facing web widget for the Talent Partner"
  bot_id      = chatbotkit_bot.talent_partner.id

  title = "Find your next role"
  intro = "Tell me what you're looking for and I'll surface roles and put you forward."
}

resource "chatbotkit_slack_integration" "client_slack" {
  count = var.client_slack_bot_token == "" ? 0 : 1

  name           = "Client Partner for Slack"
  description    = "Employer-facing Client Partner in Slack"
  bot_id         = chatbotkit_bot.client_partner.id
  bot_token      = var.client_slack_bot_token
  signing_secret = var.client_slack_signing_secret
}

# ============================================================================
# Variables
# ============================================================================

variable "marketplace_service_token" {
  description = "Shared service-account bearer token for the marketplace gateway"
  type        = string
  sensitive   = true
}

variable "talent_oauth_client_id" {
  description = "OAuth client id for candidate (Talent) sign-in"
  type        = string
  sensitive   = true
}

variable "talent_oauth_client_secret" {
  description = "OAuth client secret for candidate (Talent) sign-in"
  type        = string
  sensitive   = true
}

variable "client_oauth_client_id" {
  description = "OAuth client id for employer (Client) sign-in"
  type        = string
  sensitive   = true
}

variable "client_oauth_client_secret" {
  description = "OAuth client secret for employer (Client) sign-in"
  type        = string
  sensitive   = true
}

variable "client_slack_bot_token" {
  description = "Slack bot token for the Client Partner (leave empty to skip Slack)"
  type        = string
  default     = ""
  sensitive   = true
}

variable "client_slack_signing_secret" {
  description = "Slack signing secret for the Client Partner (leave empty to skip Slack)"
  type        = string
  default     = ""
  sensitive   = true
}

# ============================================================================
# Outputs
# ============================================================================

output "talent_partner_bot_id" {
  description = "The ID of the supply-side Talent Partner bot"
  value       = chatbotkit_bot.talent_partner.id
}

output "client_partner_bot_id" {
  description = "The ID of the demand-side Client Partner bot"
  value       = chatbotkit_bot.client_partner.id
}

output "marketplace_service_secret_id" {
  description = "The ID of the shared marketplace service token"
  value       = chatbotkit_secret.marketplace_service_token.id
}
