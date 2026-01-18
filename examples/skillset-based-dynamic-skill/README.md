# Skillset-based Dynamic Skill Reference Architecture

This example demonstrates a reference architecture for an AI agent that can dynamically load and utilize skills from skillsets.

## Overview

Skills are specialized instructions that teach AI agents how to perform specific tasks. This architecture takes the concept further by packaging skills as installable skillsets—self-contained units that can be activated on demand to extend an agent's capabilities.

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Dynamic Skills Agent                  │
│                                                          │
│  ┌────────────────────────────────────────────────────┐ │
│  │              Core Skillset                         │ │
│  │  ┌──────────────────┐  ┌──────────────────┐       │ │
│  │  │   List Skills    │  │   Install Skill  │       │ │
│  │  │   (Discover)     │  │   (Activate)     │       │ │
│  │  └──────────────────┘  └──────────────────┘       │ │
│  └────────────────────────────────────────────────────┘ │
│                                                          │
│              ┌───────────────┐                          │
│              │ Skill Library │                          │
│              └───────┬───────┘                          │
│       ┌──────────────┼──────────────┬──────────────┐   │
│       │              │              │              │   │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐  │
│  │  Data   │  │ Content │  │Research │  │ Problem │  │
│  │Analysis │  │ Writing │  │  Skill  │  │ Solving │  │
│  │  Skill  │  │  Skill  │  │         │  │  Skill  │  │
│  └─────────┘  └─────────┘  └─────────┘  └─────────┘  │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

## Key Features

1. **Dynamic Skill Discovery**: The agent can list available skills to find relevant capabilities
2. **On-Demand Skill Loading**: Skills are installed into the conversation context only when needed
3. **Modular Architecture**: Each skill is a self-contained skillset with its own instructions
4. **Scalable Design**: Easy to add new skills without modifying the core agent

## Core Abilities

### List Skills
Uses the `blueprint/resource/list` template configured for skillsets to enumerate all available skills. The agent sees each skill's name and description, enabling it to identify which skills are relevant to the current task.

### Install Skill
Uses the `conversation/skillset/install[by-id]` template to activate a skill by bringing its skillset into the conversation context. Once installed, the skill's instructions become part of the agent's system prompt.

## Skill Format

Each skill skillset follows a structured description format:
```
Short description
---
Longer instructions how to use the skill
```

This convention allows the agent to:
- Quickly scan available skills (short description)
- Access comprehensive guidance when needed (detailed instructions)

## Usage

1. Set your ChatBotKit API key:
```bash
export CHATBOTKIT_API_KEY="your-api-key"
```

2. Initialize Terraform:
```bash
terraform init
```

3. Review the planned changes:
```bash
terraform plan
```

4. Apply the configuration:
```bash
terraform apply
```

5. Test the agent by asking it to:
   - List available skills
   - Install a specific skill by ID
   - Use the installed skill's capabilities

## Customization

### Adding New Skills

To add a new skill, create a new skillset resource:

```hcl
resource "chatbotkit_skillset" "skill_5" {
  name = "Technical Writing Skill"
  description = <<-EOT
    Create technical documentation
    ---
    You are a technical writer. When creating documentation:
    - Use clear, precise language
    - Include code examples where appropriate
    - Structure content with headers and sections
    - Add diagrams and visual aids when helpful
  EOT
}
```

### Adding Abilities to Skills

Skills can include their own abilities. For example, a "Research Skill" might include web search and fetch abilities:

```hcl
resource "chatbotkit_skillset_ability" "research_search" {
  skillset_id = chatbotkit_skillset.skill_3.id
  name        = "Web Search"
  description = "Search the web for information"
  instruction = <<-EOT
    ```search
    query: $[query! ys|search terms]
    ```
  EOT
}
```

## When to Use This Pattern

This pattern is ideal when:
- Your agent needs diverse, specialized capabilities
- Skills should be loaded on-demand to avoid context bloat
- Skills include not just instructions but also abilities, secrets, or configurations
- You want modular, maintainable agent architectures

Compare with the file-based variant: files are ideal for purely instructional content, while skillsets shine when packaging complete capability packages with tools and configurations.

## Cleanup

To destroy all created resources:
```bash
terraform destroy
```

## Learn More

- [ChatBotKit Skillsets Documentation](https://chatbotkit.com/docs/resources/skillsets)
- [ChatBotKit Abilities Documentation](https://chatbotkit.com/docs/resources/abilities)
- [Blueprint Reference Architecture Examples](https://chatbotkit.com/blueprints)
