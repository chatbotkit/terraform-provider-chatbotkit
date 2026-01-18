# ChatBotKit Terraform Provider Examples

This directory contains example configurations demonstrating how to use the ChatBotKit Terraform Provider to create and manage various AI agent architectures.

## Available Examples

### Basic Examples

| Example | Description | Key Resources |
|---------|-------------|---------------|
| [basic](./basic/) | Minimal example showing bot creation | Bot |
| [complete](./complete/) | Full-featured setup with dataset, skillset, abilities, and integrations | Bot, Dataset, Skillset, Abilities, Trigger Integration |

### Reference Architecture Examples

These examples demonstrate production-ready architectures based on ChatBotKit blueprint patterns:

| Example | Description | Key Features |
|---------|-------------|--------------|
| [skillset-based-dynamic-skill](./skillset-based-dynamic-skill/) | Agent that dynamically discovers and loads skills | Dynamic skill loading, Modular architecture, List & Install abilities |
| [ai-employee](./ai-employee/) | Comprehensive AI Employee for professional environments | Workspace/Space, Shell execution, Gmail integration, Notion integration, Dynamic skillsets |
| [matillion-operations](./matillion-operations/) | AI-powered operations center for Matillion Data Productivity Cloud | 4 specialized skillsets, Pipeline management, Schedule management, Monitoring, Scheduled health checks |

## Getting Started

### Prerequisites

1. [Terraform](https://www.terraform.io/downloads.html) >= 1.0
2. A ChatBotKit account and API key from [chatbotkit.com](https://chatbotkit.com)

### Quick Start

1. Set your ChatBotKit API key:

```bash
export CHATBOTKIT_API_KEY="your-api-key"
```

2. Choose an example and navigate to its directory:

```bash
cd basic  # or complete, skillset-based-dynamic-skill, etc.
```

3. Initialize Terraform:

```bash
terraform init
```

4. Review the planned changes:

```bash
terraform plan
```

5. Apply the configuration:

```bash
terraform apply
```

6. When finished, clean up resources:

```bash
terraform destroy
```

## Example Details

### Basic Example
Perfect for getting started. Creates a simple bot with minimal configuration.

**What you'll learn:**
- Basic provider configuration
- Creating a bot resource
- Setting bot properties

**Use when:** You want to understand the fundamentals or quickly test the provider.

### Complete Example
Demonstrates a full-featured AI assistant with knowledge base and web capabilities.

**What you'll learn:**
- Dataset creation and bot linking
- Skillset and ability configuration
- Trigger integration setup
- Resource dependencies

**Use when:** You need a comprehensive reference or want to build a feature-rich assistant.

### Skillset-based Dynamic Skill Example
Shows how to build agents that can discover and activate skills on-demand.

**What you'll learn:**
- Dynamic skill loading patterns
- Modular agent architectures
- Blueprint resource templates
- Skill packaging best practices

**Use when:** Building agents that need diverse, specialized capabilities loaded contextually.

### AI Employee Example
Demonstrates a sophisticated AI employee with workspace and multiple integrations.

**What you'll learn:**
- Space/Workspace configuration
- Shell execution abilities
- OAuth2 integration (Gmail, Notion)
- Personal-kind secrets for user authentication
- Complex multi-skillset architecture

**Use when:** Building autonomous digital team members with access to tools and services.

### Matillion Operations Example
AI-powered operations center for data pipeline management.

**What you'll learn:**
- Multi-skillset architecture (4 specialized skillsets)
- Pipeline operations automation
- Schedule management
- Monitoring and compliance
- Infrastructure management
- Scheduled trigger integrations
- Template-based secrets

**Use when:** Automating data pipeline operations and monitoring.

## Architecture Patterns

### Single Skillset Pattern
```
Bot → Skillset → Abilities
```
Used in: `basic`, `complete`

**Best for:** Simple, focused agents with a small number of related abilities.

### Dynamic Skill Loading Pattern
```
Bot → Core Skillset → List/Install Abilities
              ↓
        Skill Skillsets (loaded on-demand)
```
Used in: `skillset-based-dynamic-skill`, `ai-employee`, `matillion-operations`

**Best for:** Agents that need diverse capabilities but want to avoid context bloat.

### Multi-Skillset with Workspace Pattern
```
Bot → Core Skillset
      ↓
      Space/Workspace → Shell Execution
      ↓
      Specialized Skillsets (Mail, Notion, etc.)
```
Used in: `ai-employee`

**Best for:** Agents that need secure sandboxed operations and multiple integrations.

## Common Patterns

### OAuth2 Integrations
Many services require OAuth2 authentication. ChatBotKit provides template-based secrets:

```hcl
resource "chatbotkit_secret" "google_mail" {
  name        = "Google Mail OAuth2 Token"
  description = "OAuth2 token for accessing Google Mail"
  type        = "template"
  kind        = "personal"  # or "shared"
  
  config = jsonencode({
    template = "platform/google/mail"
  })
}
```

**Personal vs Shared:**
- `personal`: Each end-user authenticates with their own account
- `shared`: Single shared credential for all users

### Ability Instructions
Abilities use YAML-like syntax with template references:

```hcl
instruction = <<-EOT
  template: action/template/name
  parameters:
    param1: ''
    param2: $[value! ys|description for AI]
EOT
```

### Dynamic Resource Discovery
Enable agents to discover and load resources:

```hcl
resource "chatbotkit_skillset_ability" "list_resources" {
  skillset_id = chatbotkit_skillset.core.id
  name        = "List Resources"
  description = "Discover available resources"
  instruction = <<-EOT
    template: blueprint/resource/list
    parameters:
      type: skillset  # or file, dataset, etc.
  EOT
}
```

## Tips and Best Practices

1. **Start Simple**: Begin with the `basic` example and gradually add complexity
2. **Use Variables**: Parameterize repeated values for easier maintenance
3. **Organize Skillsets**: Group related abilities into focused skillsets
4. **Document Backstories**: Clear backstories improve agent behavior
5. **Test Incrementally**: Apply and test changes in small increments
6. **Review Plans**: Always review `terraform plan` before applying
7. **Use Outputs**: Export resource IDs for use in other configurations

## Resource Dependencies

Terraform automatically handles dependencies based on resource references:

```
Dataset (independent)
   ↓
Skillset (independent)
   ↓
Abilities (depend on Skillset, optionally Secret)
   ↓
Bot (depends on Dataset/Skillset)
   ↓
Integrations (depend on Bot)
```

## Troubleshooting

### "Provider not found"
```bash
terraform init  # Re-initialize to download the provider
```

### "Invalid API key"
```bash
# Verify your API key is set
echo $CHATBOTKIT_API_KEY

# Or configure directly in the provider block
provider "chatbotkit" {
  api_key = "your-api-key"
}
```

### "Resource already exists"
If applying fails due to existing resources, import them:
```bash
terraform import chatbotkit_bot.example bot_abc123
```

## Next Steps

- Browse the [ChatBotKit Documentation](https://chatbotkit.com/docs)
- Explore [Blueprint Examples](https://chatbotkit.com/blueprints)
- Review the [Terraform Provider Documentation](https://registry.terraform.io/providers/chatbotkit/chatbotkit/latest/docs)
- Check out the [ChatBotKit SDKs](https://github.com/chatbotkit) for programmatic access

## Contributing

Found an issue or have an improvement? Please open an issue or submit a pull request on the [ChatBotKit Platform repository](https://github.com/chatbotkit/cbk-platform).
