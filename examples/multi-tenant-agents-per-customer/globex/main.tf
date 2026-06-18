# Globex's bespoke agent (a module).
#
# Composed into the shared root (../main.tf), which wires it to Globex's
# sub-account via a provider alias + run_as. Globex wants a research agent with a
# sandbox — a completely different shape from Acme's support bot. Each customer
# folder evolves independently.

terraform {
  required_providers {
    chatbotkit = {
      source = "chatbotkit/chatbotkit"
    }
  }
}

resource "chatbotkit_space" "workspace" {
  name        = "Globex Workspace"
  description = "Sandbox for Globex's research agent"
}

resource "chatbotkit_skillset" "research" {
  name        = "Globex Research Tools"
  description = "Tools for Globex's research agent"
}

resource "chatbotkit_skillset_ability" "shell" {
  skillset_id = chatbotkit_skillset.research.id
  space_id    = chatbotkit_space.workspace.id
  name        = "shell"
  description = "Run commands, read/write files, and import URLs in the workspace"
  instruction = <<-EOT
    template: pack/shell
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "web_search" {
  skillset_id = chatbotkit_skillset.research.id
  name        = "search_web"
  description = "Search the web for current information"
  instruction = <<-EOT
    ```search
    query: $[query! ys|the search query to find information on the web]
    ```
  EOT
}

resource "chatbotkit_bot" "research" {
  name        = "Globex Research"
  description = "Globex's deep-research analyst"
  model       = "claude-4.5-opus"
  backstory   = "You are Globex Inc's research analyst. Investigate thoroughly, cross-check sources, and deliver concise, well-sourced findings."
  skillset_id = chatbotkit_skillset.research.id
}

output "bot_id" {
  description = "Globex's research bot ID"
  value       = chatbotkit_bot.research.id
}
