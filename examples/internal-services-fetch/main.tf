# ============================================================================
# Internal Services Assistant — fetch + secrets reference architecture
# ============================================================================
#
# This example shows how an agent reaches your *internal* corporate services
# with the `fetch` action, authenticated by two different kinds of secret:
#
#   1. A SHARED secret (machine-to-machine).
#      A long-lived service token the agent presents to internal systems that
#      trust it as a service account. There is no human in the loop — the agent
#      *is* the principal. Use this for read-only platform data that every
#      employee may see (service catalog, status, on-call schedule, ...).
#
#   2. A PERSONAL OAuth secret (user-on-behalf).
#      Each employee signs in once through your SSO. The agent then calls
#      permission-enforcing internal systems *as that employee*, so the backend
#      applies the same access control it would for a human. Use this for
#      anything scoped to the person talking to the agent (their tickets, their
#      time-off balance, opening an incident under their name).
#
# The `open_incident` ability uses BOTH at once: the shared service token gets
# the agent past the internal API gateway (machine identity), while the personal
# OAuth token attributes the incident to the employee who asked for it.
#
# The abilities below are written as STRUCTURED fetch instructions: a `!fetch`
# YAML action tag with typed helpers instead of inline string templates.
#
#   !reference SECRET_DEFAULT  -> the secret linked to the ability via `secret_id`
#   !reference SECRET_<NAME>   -> any other secret in the account, matched by name
#                                 (e.g. !reference SECRET_ACME_INTERNAL_SERVICE_TOKEN
#                                 resolves the secret named "Acme Internal Service
#                                 Token"). Names are matched case-insensitively
#                                 with non-word characters collapsed to "_".
#   !string { name, ... }      -> a parameter the model fills in at call time.
#
# A secret already resolves to a *complete* Authorization value including its
# scheme — a `bearer`/`oauth` secret yields `Bearer <token>`. So the header is
# just `Authorization: !reference SECRET_...`; do NOT add your own `Bearer `
# prefix or you would end up with `Bearer Bearer <token>`.
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

# (1) SHARED service token — agent-to-internal-system communication.
#
# `kind = "shared"` means this is a single account-wide value (not per-user).
# `type = "bearer"` means it resolves to `Bearer <token>` wherever it is
# referenced. Every conversation the agent has uses the same token, so the
# internal service sees one stable service identity no matter who is chatting.
resource "chatbotkit_secret" "internal_service_token" {
  name        = "Acme Internal Service Token"
  description = "Service-account bearer token the assistant uses to call read-only internal platform APIs"
  type        = "bearer"
  kind        = "shared"

  # The actual token value. Keep it out of source control — pass it via
  # terraform.tfvars or TF_VAR_internal_service_token.
  value = var.internal_service_token
}

# (2) PERSONAL OAuth secret — the employee authenticates, the agent acts as them.
#
# `kind = "personal"` means the value is resolved per end-user: the first time
# an employee triggers an ability that references this secret, ChatBotKit asks
# them to sign in through your SSO and stores *their* token. From then on the
# agent calls internal services on that employee's behalf.
resource "chatbotkit_secret" "acme_sso" {
  name        = "Acme SSO (Employee)"
  description = "Per-employee OAuth token from Acme SSO, used to act on behalf of the signed-in user"
  type        = "oauth"
  kind        = "personal"

  # Standard OAuth 2.0 authorization-code configuration. Point these at your own
  # identity provider (this example uses a Keycloak-style realm). Keys are
  # camelCase — that is what the OAuth resolver reads.
  config = {
    clientId         = var.acme_sso_client_id
    clientSecret     = var.acme_sso_client_secret
    authorizationUrl = "https://sso.acme.corp/realms/employees/protocol/openid-connect/auth"
    tokenUrl         = "https://sso.acme.corp/realms/employees/protocol/openid-connect/token"
    grantType        = "authorization_code"
    scope            = "openid profile helpdesk:read helpdesk:write hr:read"
  }
}

# ============================================================================
# Skillset
# ============================================================================
# All the internal-service tools live in one skillset attached to the bot.

resource "chatbotkit_skillset" "internal_services" {
  name        = "Internal Services"
  description = "Tools for reaching Acme's internal platform, helpdesk and HR systems"
}

# ============================================================================
# Abilities — SHARED service token (machine identity)
# ============================================================================

# Read-only platform data that is the same for everyone. The agent presents the
# shared service token; no employee sign-in is required. `SECRET_DEFAULT`
# resolves to the secret linked below via `secret_id`, and yields `Bearer <token>`.
resource "chatbotkit_skillset_ability" "lookup_service_status" {
  skillset_id = chatbotkit_skillset.internal_services.id
  secret_id   = chatbotkit_secret.internal_service_token.id

  name        = "lookup_service_status"
  description = "Look up the health and ownership of an internal service from the platform catalog"
  instruction = <<-EOT
    !fetch
    method: GET
    url: https://api.internal.acme.corp/platform/v1/services
    query:
      q: !string
        name: service
        description: the name or id of the internal service, e.g. billing-api
    headers:
      Authorization: !reference SECRET_DEFAULT
      Accept: application/json
    options:
      format: json
  EOT
}

# ============================================================================
# Abilities — PERSONAL OAuth (acts as the signed-in employee)
# ============================================================================

# Lists the helpdesk tickets that belong to the employee currently chatting.
# Because the request carries the employee's own OAuth token, the helpdesk
# backend returns only the tickets they are allowed to see.
resource "chatbotkit_skillset_ability" "list_my_tickets" {
  skillset_id = chatbotkit_skillset.internal_services.id
  secret_id   = chatbotkit_secret.acme_sso.id

  name        = "list_my_tickets"
  description = "List the current employee's own IT helpdesk tickets"
  instruction = <<-EOT
    !fetch
    method: GET
    url: https://api.internal.acme.corp/helpdesk/v1/tickets
    query:
      status: !string
        name: status
        description: optional status filter; leave empty for all tickets
        optional: true
        enum:
          - open
          - pending
          - closed
    headers:
      Authorization: !reference SECRET_DEFAULT
      Accept: application/json
    options:
      format: json
  EOT
}

# Reads the signed-in employee's remaining time-off balance. The HR backend
# identifies the user from the token, so no employee id is passed in.
resource "chatbotkit_skillset_ability" "get_my_timeoff_balance" {
  skillset_id = chatbotkit_skillset.internal_services.id
  secret_id   = chatbotkit_secret.acme_sso.id

  name        = "get_my_timeoff_balance"
  description = "Get the current employee's remaining time-off (PTO) balance"
  instruction = <<-EOT
    !fetch
    method: GET
    url: https://api.internal.acme.corp/hr/v1/timeoff/balance
    headers:
      Authorization: !reference SECRET_DEFAULT
      Accept: application/json
    options:
      format: json
  EOT
}

# ============================================================================
# Ability — BOTH secrets in one request (gateway + on-behalf)
# ============================================================================

# Opens an incident through the internal API gateway. The gateway authenticates
# the *caller* — the agent service — via the standard Authorization header, so
# that carries the shared service token. The incident must be attributed to a
# real employee, so the gateway expects the user's own token forwarded alongside
# it in X-Forwarded-Authorization. Both secrets resolve to a complete
# `Bearer <token>` value, hence no manual scheme on either header.
#
#   Authorization:             !reference SECRET_ACME_INTERNAL_SERVICE_TOKEN  <- shared, by name
#   X-Forwarded-Authorization: !reference SECRET_DEFAULT                      <- personal, linked
resource "chatbotkit_skillset_ability" "open_incident" {
  skillset_id = chatbotkit_skillset.internal_services.id
  secret_id   = chatbotkit_secret.acme_sso.id

  name        = "open_incident"
  description = "Open an incident in the internal incident tracker on behalf of the current employee"
  instruction = <<-EOT
    !fetch
    method: POST
    url: https://gateway.internal.acme.corp/incidents/v1/incidents
    headers:
      Authorization: !reference SECRET_ACME_INTERNAL_SERVICE_TOKEN
      X-Forwarded-Authorization: !reference SECRET_DEFAULT
      Content-Type: application/json
    body:
      title: !string
        name: title
        description: a short, descriptive incident title
      severity: !string
        name: severity
        description: the incident severity
        enum:
          - low
          - medium
          - high
          - critical
      service: !string
        name: service
        description: the affected internal service name or id
      description: !string
        name: description
        description: a detailed description of the incident and its impact
        optional: true
    options:
      format: json
  EOT
}

# ============================================================================
# Bot
# ============================================================================

resource "chatbotkit_bot" "assistant" {
  name        = "Acme Internal Operations Assistant"
  description = "Helps employees reach internal platform, helpdesk and HR services"
  model       = "claude-4.5-sonnet"

  backstory = <<-EOT
    You are Acme's Internal Operations Assistant. You help employees get answers
    from internal systems quickly and safely.

    You have two ways of reaching internal services:

    - For shared, read-only platform information (such as service health and
      ownership) you call internal APIs directly as a trusted service. This data
      is the same for everyone and needs no sign-in.

    - For anything specific to the person you are talking to (their helpdesk
      tickets, their time-off balance, opening an incident in their name) you
      act on that employee's behalf. The first time someone uses one of these
      tools they will be asked to sign in through Acme SSO. If you receive an
      authorization prompt, relay it to the employee and wait for them to
      authenticate before retrying — never ask them for a password directly.

    Be concise, confirm before creating or changing anything (like opening an
    incident), and never expose raw tokens or internal URLs to the employee.
  EOT

  skillset_id = chatbotkit_skillset.internal_services.id
}

# ============================================================================
# Channel (optional) — deploy the assistant into Slack
# ============================================================================
# Provide Slack credentials in terraform.tfvars to activate the channel. Leave
# them empty to skip it.

resource "chatbotkit_slack_integration" "slack" {
  count = var.slack_bot_token == "" ? 0 : 1

  name           = "Internal Assistant for Slack"
  description    = "Acme's Internal Operations Assistant in Slack"
  bot_id         = chatbotkit_bot.assistant.id
  bot_token      = var.slack_bot_token
  signing_secret = var.slack_signing_secret
}

# ============================================================================
# Variables
# ============================================================================

variable "internal_service_token" {
  description = "Service-account bearer token for read-only internal platform APIs"
  type        = string
  sensitive   = true
}

variable "acme_sso_client_id" {
  description = "OAuth client id for Acme SSO"
  type        = string
  sensitive   = true
}

variable "acme_sso_client_secret" {
  description = "OAuth client secret for Acme SSO"
  type        = string
  sensitive   = true
}

variable "slack_bot_token" {
  description = "Slack bot token (leave empty to skip the Slack channel)"
  type        = string
  default     = ""
  sensitive   = true
}

variable "slack_signing_secret" {
  description = "Slack signing secret (leave empty to skip the Slack channel)"
  type        = string
  default     = ""
  sensitive   = true
}

# ============================================================================
# Outputs
# ============================================================================

output "bot_id" {
  description = "The ID of the internal operations assistant"
  value       = chatbotkit_bot.assistant.id
}

output "skillset_id" {
  description = "The ID of the internal services skillset"
  value       = chatbotkit_skillset.internal_services.id
}

output "shared_secret_id" {
  description = "The ID of the shared service-account secret"
  value       = chatbotkit_secret.internal_service_token.id
}

output "personal_oauth_secret_id" {
  description = "The ID of the personal Acme SSO OAuth secret"
  value       = chatbotkit_secret.acme_sso.id
}
