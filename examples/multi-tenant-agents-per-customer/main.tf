# Architecture 1 — a DIFFERENT agent per customer, in one shared state.
#
# Each customer is its own folder/module (./acme, ./globex) with a bespoke agent.
# This root composes them all into a single state and deploys them in one apply,
# wiring each module to its customer's sub-account with a provider alias + run_as.
#
# One partner/master token (set via CHATBOTKIT_API_KEY) operates on every
# sub-account; run_as selects which one (the X-RunAs-UserId header). No
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
  description = "Acme's sub-account (partner user) ID"
  type        = string
}

variable "globex_account_id" {
  description = "Globex's sub-account (partner user) ID"
  type        = string
}

# api_key comes from CHATBOTKIT_API_KEY (one partner/master token). Each alias
# selects which sub-account to operate on.
provider "chatbotkit" {
  alias  = "acme"
  run_as = var.acme_account_id
}

provider "chatbotkit" {
  alias  = "globex"
  run_as = var.globex_account_id
}

# Each customer's bespoke agent, deployed into their sub-account.
module "acme" {
  source    = "./acme"
  providers = { chatbotkit = chatbotkit.acme }
}

module "globex" {
  source    = "./globex"
  providers = { chatbotkit = chatbotkit.globex }
}

output "bots" {
  description = "Bot IDs per customer sub-account"
  value = {
    acme   = module.acme.bot_id
    globex = module.globex.bot_id
  }
}
