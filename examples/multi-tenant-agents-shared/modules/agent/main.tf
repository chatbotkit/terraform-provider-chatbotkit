# The one canonical agent, shared by every customer.
#
# Instantiated once per customer by the per-tenant root (../../main.tf). The
# caller passes a chatbotkit provider authenticated with that customer's
# sub-account token, so every resource below is created inside the customer's
# isolated sub-account. Improve the agent here and it ships to all customers on
# the next deploy.

terraform {
  required_providers {
    chatbotkit = {
      source = "chatbotkit/chatbotkit"
    }
  }
}

variable "customer_name" {
  description = "Display name of the customer (used to name the agent's resources)"
  type        = string
}

variable "model" {
  description = "The model the customer's agent should use"
  type        = string
  default     = "claude-4.5-sonnet"
}

resource "chatbotkit_space" "workspace" {
  name        = "${var.customer_name} Workspace"
  description = "Isolated workspace for ${var.customer_name}'s agent"
}

resource "chatbotkit_skillset" "abilities" {
  name        = "${var.customer_name} Abilities"
  description = "Toolset for ${var.customer_name}'s agent"
}

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

resource "chatbotkit_bot" "agent" {
  name        = "${var.customer_name} Assistant"
  description = "Dedicated assistant for ${var.customer_name}"
  model       = var.model
  backstory   = "You are the dedicated AI assistant for ${var.customer_name}. Be helpful, concise, and professional."
  skillset_id = chatbotkit_skillset.abilities.id
}

output "bot_id" {
  description = "The bot ID created in the customer's sub-account"
  value       = chatbotkit_bot.agent.id
}

output "workspace_id" {
  description = "The workspace ID created in the customer's sub-account"
  value       = chatbotkit_space.workspace.id
}
