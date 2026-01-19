# Dynamic MCP Search and Install Architecture
#
# This example demonstrates a reference architecture for an AI agent that can
# search for and install Model Context Protocol (MCP) servers dynamically from
# a known registry.
#
# Architecture highlights:
# - Bot with dynamic MCP discovery capabilities
# - Search MCP servers by name or description
# - Install MCP servers on-demand
# - Self-extending system that grows functionality as needed
#
# Prerequisites:
# - Set the CHATBOTKIT_API_KEY environment variable

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
# Core Skillset
# ============================================================================
# Contains abilities for discovering and installing MCP servers

resource "chatbotkit_skillset" "mcp_discovery" {
  name        = "MCP Discovery"
  description = "Abilities for searching and installing MCP servers dynamically"
}

# ============================================================================
# MCP Discovery Abilities
# ============================================================================
# Abilities that enable the agent to search for and install MCP servers

resource "chatbotkit_skillset_ability" "search_mcp" {
  skillset_id = chatbotkit_skillset.mcp_discovery.id
  name        = "Search MCP Servers"
  description = "Search for Model Context Protocol (MCP) servers by name or description"
  instruction = <<-EOT
    template: registry/server/search
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "install_mcp" {
  skillset_id = chatbotkit_skillset.mcp_discovery.id
  name        = "Install MCP"
  description = "Bring MCP (model context protocol) functions into context"
  instruction = <<-EOT
    template: conversation/mcp/install[url]
    parameters:
      url: ''
  EOT
}

# ============================================================================
# Dynamic MCP Agent Bot
# ============================================================================
# The agent bot that can search for and install MCP servers on-demand

resource "chatbotkit_bot" "agent" {
  name        = "Dynamic MCP Agent"
  description = "An AI agent that can search for and install Model Context Protocol (MCP) servers dynamically"
  backstory   = <<-EOT
    You are an intelligent agent with the ability to discover and dynamically integrate
    Model Context Protocol (MCP) servers to extend your capabilities. When faced with a
    task that requires specific tools or integrations, you can:
    
    1. **Search for MCP Servers**: Use the Search MCP Servers ability to find available
       MCP servers by name or description. Look for servers that provide the functionality
       needed for the current task.
    
    2. **Install MCP Servers**: Once you've identified a relevant MCP server, use the
       Install MCP ability to bring its functions into your conversation context. Provide
       the server's URL to activate it.
    
    3. **Use Installed Capabilities**: After installing an MCP server, you'll have access
       to its tools and functions. Use them to accomplish tasks that require specialized
       capabilities.
    
    This dynamic approach allows you to adapt to diverse requirements without requiring
    manual reconfiguration. You can discover and integrate exactly the tools you need,
    when you need them.
    
    ## Guidelines for MCP Discovery
    
    **When to Search**:
    - User requests functionality you don't currently have
    - Task requires specialized tools or data access
    - Integration with external services would be beneficial
    
    **How to Search**:
    - Use descriptive search terms matching the needed functionality
    - Look for servers from trusted sources or well-known providers
    - Review server descriptions to ensure they match your needs
    
    **Before Installing**:
    - Verify the MCP server provides the needed capabilities
    - Check that it's from a reliable source
    - Explain to the user what capabilities you're adding and why
    
    **After Installing**:
    - Explore the newly available tools and functions
    - Use them appropriately to complete the user's request
    - Inform the user about the new capabilities you've acquired
    
    Remember: You are a self-extending system. Your ability to discover and integrate
    new MCP servers makes you highly adaptable to evolving requirements.
  EOT
  model       = "claude-4.5-sonnet"

  skillset_id = chatbotkit_skillset.mcp_discovery.id
}

# ============================================================================
# Outputs
# ============================================================================

output "bot_id" {
  description = "The ID of the dynamic MCP agent"
  value       = chatbotkit_bot.agent.id
}

output "skillset_id" {
  description = "The ID of the MCP discovery skillset"
  value       = chatbotkit_skillset.mcp_discovery.id
}
