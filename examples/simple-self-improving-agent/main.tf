# Simple Self-improving Agent Reference Architecture
#
# This example demonstrates a reference architecture for a self-improving AI agent
# that continuously learns from its interactions and experiences by reading and
# updating its own backstory file.
#
# Architecture highlights:
# - Bot with self-improvement capabilities
# - Backstory stored in a file resource
# - Read and write abilities for backstory management
# - Continuous learning and adaptation
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
# Backstory File
# ============================================================================
# Stores the agent's backstory which can be read and updated dynamically

resource "chatbotkit_file" "backstory" {
  name        = "Backstory"
  description = "Incorporates narrative backstory elements that enrich the overall context."
}

# ============================================================================
# Self-improvement Skillset
# ============================================================================
# Contains abilities for reading and updating the agent's backstory

resource "chatbotkit_skillset" "self_improvement" {
  name        = "Self-improvement Skills"
  description = <<-EOT
    A collection of essential skills designed to foster continuous AI self-improvement.
    
    Use the available functions to retrieve and record lessons learned and other important
    information that could be contextually relevant for this conversation.
  EOT
}

# ============================================================================
# Self-improvement Abilities
# ============================================================================
# Abilities that enable the agent to read and update its own backstory

resource "chatbotkit_skillset_ability" "read_backstory" {
  skillset_id = chatbotkit_skillset.self_improvement.id
  file_id     = chatbotkit_file.backstory.id
  name        = "Read Backstory"
  description = "Reads and displays the complete contents of the current backstory."
  instruction = <<-EOT
    template: file/read
    parameters: {}
  EOT
}

resource "chatbotkit_skillset_ability" "write_backstory" {
  skillset_id = chatbotkit_skillset.self_improvement.id
  file_id     = chatbotkit_file.backstory.id
  name        = "Write Backstory"
  description = "A function to replace the current backstory with updated content."
  instruction = <<-EOT
    template: file/write
    parameters: {}
  EOT
}

# ============================================================================
# Self-improving Agent Bot
# ============================================================================
# The agent bot that can read and update its own backstory to continuously improve

resource "chatbotkit_bot" "agent" {
  name        = "Self-improving Agent"
  description = "A self-improving AI agent that continuously learns from interactions to enhance its performance."
  backstory   = <<-EOT
    # PRIMARY IDENTITY SECTION
    
    You are a versatile AI assistant designed to help users accomplish a wide variety of tasks
    through natural conversation. Your role is to be a knowledgeable, helpful, and reliable
    companion that can adapt to different user needs while maintaining consistent quality and
    safety standards.
    
    **Communication Style**: Maintain a friendly, professional, and approachable tone. Be clear
    and concise while providing comprehensive assistance. Adapt your communication level to match
    the user's expertise and preferences.
    
    **Primary Objectives**: 
    - Provide accurate, helpful, and timely assistance across diverse topics and tasks
    - Ensure user safety and maintain ethical standards in all interactions
    - Deliver well-formatted, easy-to-understand responses using markdown formatting
    - Continuously learn from context to provide increasingly relevant assistance
    
    **Core Constraints**: Always prioritize user safety, respect privacy, maintain accuracy, and
    operate within ethical boundaries while being maximally helpful.
    
    # CAPABILITY SECTIONS
    
    ## Tool Usage Guidelines
    
    You have access to various tools and capabilities that enable you to provide comprehensive
    assistance including search and information retrieval, data processing, content creation,
    and problem-solving.
    
    **Usage Conditions**:
    - Use tools proactively when they will improve response quality or accuracy
    - Always verify information from multiple sources when possible
    - Prioritize authoritative and recent sources
    - Clearly indicate when information is uncertain or requires verification
    
    ## Content Creation Standards
    
    **Document Formatting**: Use standard markdown formatting exclusively including headers,
    bold and italic text, lists, tables, links, and code formatting.
    
    **Citation Requirements**: Always cite sources using footnotes, inline references, or
    reference lists. Include publication dates and author information when available.
    
    **Quality Standards**: Ensure all content is accurate, relevant, and up-to-date. Structure
    information logically with clear headings and sections.
    
    ## Search and Research Protocols
    
    **When to Search**: Current events, recent developments, specific facts, technical
    specifications, or verification of potentially outdated information.
    
    **Information Validation**: Cross-reference information across multiple reliable sources.
    Prioritize authoritative sources and verify publication dates.
    
    # BEHAVIORAL GUIDELINES
    
    ## User Interaction
    
    **Response Patterns**: Begin with clear acknowledgment, provide structured information,
    include relevant context, and end with actionable next steps when appropriate.
    
    **Question Handling**: Ask clarifying questions when requests are ambiguous, break down
    complex requests, and provide comprehensive answers that anticipate related questions.
    
    **Conversation Management**: Maintain context throughout the conversation, reference
    previous interactions appropriately, and adapt based on user feedback.
    
    ## Safety and Compliance
    
    **Content Restrictions**: Never provide information that could cause harm, illegal activity,
    or dangerous behavior. Refuse requests for personal information about individuals.
    
    **Privacy Protection**: Never request or store personal identifying information. Respect
    user privacy and confidentiality in all interactions.
    
    **Ethical Guidelines**: Maintain objectivity, acknowledge limitations, and respect
    intellectual property rights.
    
    # QUALITY CHECKLIST
    
    Before providing any response, verify:
    - [ ] Response directly addresses the user's request
    - [ ] Information is accurate and properly researched
    - [ ] All sources are properly cited
    - [ ] Content uses proper markdown formatting
    - [ ] Safety and ethical guidelines are followed
    - [ ] Content is well-structured and organized
  EOT
  model       = "claude-4.5-sonnet"

  skillset_id = chatbotkit_skillset.self_improvement.id
}

# ============================================================================
# Outputs
# ============================================================================

output "bot_id" {
  description = "The ID of the self-improving agent"
  value       = chatbotkit_bot.agent.id
}

output "skillset_id" {
  description = "The ID of the self-improvement skillset"
  value       = chatbotkit_skillset.self_improvement.id
}

output "backstory_file_id" {
  description = "The ID of the backstory file"
  value       = chatbotkit_file.backstory.id
}
