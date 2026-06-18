# Acme's bespoke agent (a module).
#
# This folder holds Acme's own agent definition. It is composed into the shared
# root (../main.tf), which wires it to Acme's sub-account via a provider alias +
# run_as — so this module just declares resources and lets the root pick the
# account. Acme wants a customer-support agent; Globex's folder is different.

terraform {
  required_providers {
    chatbotkit = {
      source = "chatbotkit/chatbotkit"
    }
  }
}

resource "chatbotkit_skillset" "support" {
  name        = "Acme Support Tools"
  description = "Tools for Acme's support agent"
}

resource "chatbotkit_skillset_ability" "web_search" {
  skillset_id = chatbotkit_skillset.support.id
  name        = "search_web"
  description = "Search the web for current information"
  instruction = <<-EOT
    ```search
    query: $[query! ys|the search query to find information on the web]
    ```
  EOT
}

resource "chatbotkit_bot" "support" {
  name        = "Acme Support"
  description = "Acme's customer-support assistant"
  model       = "claude-4.5-sonnet"
  backstory   = "You are Acme Corporation's customer-support assistant. Be warm, patient, and solution-oriented."
  skillset_id = chatbotkit_skillset.support.id
}

output "bot_id" {
  description = "Acme's support bot ID"
  value       = chatbotkit_bot.support.id
}
