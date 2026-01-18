# Skillset-based Dynamic Skill Reference Architecture
#
# This example demonstrates a reference architecture for an AI agent that can
# dynamically load and utilize skills from skillsets. Skills are packaged as
# installable skillsets that can be activated on demand to extend agent capabilities.
#
# Architecture highlights:
# - Main bot with core skillset for skill management
# - Two abilities: List Skills and Install Skill
# - Multiple skill skillsets (placeholder examples)
# - Dynamic skill loading at runtime
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
# Main Agent Bot
# ============================================================================
# The agent bot that can dynamically discover and load skills

resource "chatbotkit_bot" "agent" {
  name        = "Dynamic Skills Agent"
  description = "An AI agent that can dynamically load and utilize skills from skillsets"
  backstory   = <<-EOT
    You are an intelligent agent with the ability to discover and dynamically load
    specialized skills as needed. When faced with a new task, you can:
    
    1. List available skills to find relevant capabilities
    2. Install specific skills by ID to gain new abilities
    
    Each skill is a self-contained skillset with instructions and potentially
    additional abilities. Use the List Skills ability to discover what's available,
    and Install Skill to activate capabilities you need for the current task.
  EOT
  model       = "claude-4.5-sonnet"

  skillset_id = chatbotkit_skillset.core_skills.id
}

# ============================================================================
# Core Skillset
# ============================================================================
# Contains abilities for managing skills dynamically

resource "chatbotkit_skillset" "core_skills" {
  name        = "Core Skills"
  description = "Core abilities for dynamic skill management"
}

# ============================================================================
# Core Abilities
# ============================================================================
# Abilities that enable the agent to discover and load skills

resource "chatbotkit_skillset_ability" "list_skills" {
  skillset_id = chatbotkit_skillset.core_skills.id
  name        = "List Skills"
  description = "Displays a comprehensive, organized list of skills to use during conversation when necessary"
  instruction = <<-EOT
    template: blueprint/resource/list
    parameters:
      type: skillset
  EOT
}

resource "chatbotkit_skillset_ability" "install_skill" {
  skillset_id = chatbotkit_skillset.core_skills.id
  name        = "Install Skill"
  description = "Bring a skill into context by its ID"
  instruction = <<-EOT
    template: conversation/skillset/install[by-id]
    parameters: {}
  EOT
}

# ============================================================================
# Skill Skillsets
# ============================================================================
# Individual skills packaged as skillsets (placeholder examples)
# In production, these would contain domain-specific abilities and instructions

resource "chatbotkit_skillset" "skill_1" {
  name = "Data Analysis Skill"
  description = <<-EOT
    Analyze and interpret data patterns
    ---
    You are equipped with data analysis capabilities. When analyzing data:
    - Identify key patterns and trends
    - Calculate statistical measures
    - Present findings clearly with visualizations when helpful
    - Highlight anomalies or outliers
  EOT
}

resource "chatbotkit_skillset" "skill_2" {
  name = "Content Writing Skill"
  description = <<-EOT
    Create professional written content
    ---
    You are a skilled content writer. When creating content:
    - Understand the audience and purpose
    - Use clear, engaging language
    - Structure content logically with headers and sections
    - Ensure proper grammar and style
  EOT
}

resource "chatbotkit_skillset" "skill_3" {
  name = "Research Skill"
  description = <<-EOT
    Conduct thorough research on topics
    ---
    You are equipped with research capabilities. When researching:
    - Gather information from multiple sources
    - Verify facts and cross-reference data
    - Organize findings logically
    - Cite sources appropriately
  EOT
}

resource "chatbotkit_skillset" "skill_4" {
  name = "Problem Solving Skill"
  description = <<-EOT
    Break down and solve complex problems
    ---
    You are a problem solver. When approaching problems:
    - Break complex problems into smaller components
    - Analyze root causes
    - Generate multiple potential solutions
    - Evaluate options and recommend the best approach
  EOT
}

# ============================================================================
# Outputs
# ============================================================================

output "bot_id" {
  description = "The ID of the dynamic skills agent"
  value       = chatbotkit_bot.agent.id
}

output "core_skillset_id" {
  description = "The ID of the core skillset"
  value       = chatbotkit_skillset.core_skills.id
}

output "skill_skillset_ids" {
  description = "IDs of the skill skillsets"
  value = {
    data_analysis   = chatbotkit_skillset.skill_1.id
    content_writing = chatbotkit_skillset.skill_2.id
    research        = chatbotkit_skillset.skill_3.id
    problem_solving = chatbotkit_skillset.skill_4.id
  }
}
