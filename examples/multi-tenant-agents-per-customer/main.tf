# Architecture 1: a different agent per customer in one shared state.
#
# Each customer is its own folder/module (./acme, ./globex) with a bespoke agent.
# This root composes them all into a single state and deploys them in one apply,
# wiring each module to its customer's child User with a provider alias and run_as.
#
# One parent User API token (set via CHATBOTKIT_API_TOKEN) operates on every
# child User; run_as selects which one (the X-RunAs-UserId header). No
# per-customer tokens, no for_each.
#
# To add a customer: add a folder/module, an account-id variable, a provider
# alias, and a module call.

terraform {
  required_providers {
    chatbotkit = {
      source = "chatbotkit/chatbotkit"
    }
  }
}

variable "acme_account_id" {
  description = "Acme's child User ID"
  type        = string
}

variable "globex_account_id" {
  description = "Globex's child User ID"
  type        = string
}

# api_token comes from CHATBOTKIT_API_TOKEN (one parent User API token). Each alias
# selects which child User to operate on.
provider "chatbotkit" {
  alias  = "acme"
  run_as = var.acme_account_id
}

provider "chatbotkit" {
  alias  = "globex"
  run_as = var.globex_account_id
}

# Each customer's bespoke agent, deployed into their child User.
module "acme" {
  source    = "./acme"
  providers = { chatbotkit = chatbotkit.acme }
}

module "globex" {
  source    = "./globex"
  providers = { chatbotkit = chatbotkit.globex }
}

output "bots" {
  description = "Bot IDs per customer child User"
  value = {
    acme   = module.acme.bot_id
    globex = module.globex.bot_id
  }
}
