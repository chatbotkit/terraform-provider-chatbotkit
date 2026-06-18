# Architecture 2 — the SAME agent for every customer, via provider aliases + run_as.
#
# One partner/master token (set once via CHATBOTKIT_API_KEY) operates on every
# sub-account; each provider alias selects a customer with `run_as = <account id>`
# (the X-RunAs-UserId header). No per-customer tokens. Provider aliases are
# Terraform's native multi-account mechanism — the same way the AWS provider
# targets many accounts.
#
# Every customer is deployed from the SAME ./modules/agent module, so you maintain
# the agent in one place and it ships to everyone. Put the per-customer account
# IDs in terraform.tfvars (see terraform.tfvars.example), then a single
# `terraform apply` deploys the agent into every sub-account.
#
# To add a customer: add an account-id variable, a provider alias, and a module call.

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
# just selects which sub-account to operate on.
provider "chatbotkit" {
  alias  = "acme"
  run_as = var.acme_account_id
}

provider "chatbotkit" {
  alias  = "globex"
  run_as = var.globex_account_id
}

# The same agent, deployed into each customer's sub-account.
module "acme" {
  source    = "./modules/agent"
  providers = { chatbotkit = chatbotkit.acme }

  customer_name = "Acme Corporation"
}

module "globex" {
  source    = "./modules/agent"
  providers = { chatbotkit = chatbotkit.globex }

  customer_name = "Globex Inc"
}

output "bots" {
  description = "Bot IDs per customer sub-account"
  value = {
    acme   = module.acme.bot_id
    globex = module.globex.bot_id
  }
}
