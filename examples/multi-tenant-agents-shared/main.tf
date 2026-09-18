# Architecture 2: the same agent for every customer via provider aliases and run_as.
#
# One parent User API token (set once via CHATBOTKIT_API_TOKEN) operates on every
# child User; each provider alias selects a customer with `run_as = <User ID>`
# (the X-RunAs-UserId header). No per-customer tokens. Provider aliases are
# Terraform's native multi-account mechanism, similar to how the AWS provider
# targets many accounts.
#
# Every customer is deployed from the SAME ./modules/agent module, so you maintain
# the agent in one place and it ships to everyone. Put the per-customer account
# IDs in terraform.tfvars (see terraform.tfvars.example), then a single
# `terraform apply` deploys the agent into every child User.
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

# The same agent, deployed into each customer's child User.
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
  description = "Bot IDs per customer child User"
  value = {
    acme   = module.acme.bot_id
    globex = module.globex.bot_id
  }
}
