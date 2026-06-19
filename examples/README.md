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
| [agent-framework](./agent-framework/) | An agent framework architecture (instructions, abilities, workspace, file-based skills, channels, schedules, heartbeat) — authored as a project of files that Terraform uploads and wires up | Bot, Skillset, Ability packs (shell + space skills), Space + file uploads, Slack integration, Scheduled triggers, Heartbeat |
| [internal-services-fetch](./internal-services-fetch/) | An agent that reaches internal corporate services with the `fetch` action, authenticated by a shared service token (machine-to-machine) and a personal OAuth secret (acting on behalf of the signed-in employee) | Fetch action, Shared bearer secret, Personal OAuth secret, `${SECRET_DEFAULT}` vs named references, Both secrets in one request, Slack integration |
| [multi-tenant-agents-shared](./multi-tenant-agents-shared/) | The same agent deployed into each customer's own sub-account, from one module. One master token + `run_as` per provider alias | Sub-accounts (partner users), `run_as` (X-RunAs-UserId), Provider aliases, Shared module, Per-tenant isolation |
| [multi-tenant-agents-per-customer](./multi-tenant-agents-per-customer/) | A bespoke agent per customer (each its own module/folder), composed into one shared state and deployed in a single apply. One master token + `run_as` per provider alias | Sub-accounts (partner users), `run_as` (X-RunAs-UserId), Provider aliases, Per-customer modules, Shared state |
| [dual-agent-programmable-workflows](./dual-agent-programmable-workflows/) | Two-agent architecture for workflow programming and execution | Multi-agent collaboration, Shared resources, Asymmetric access patterns, Scheduled triggers |
| [system-diagnostics-agent](./system-diagnostics-agent/) | Self-monitoring agent that reports on its own capabilities | Self-introspection, Blueprint resource discovery, Scheduled diagnostics, Automated reporting |
| [second-brain](./second-brain/) | Personal knowledge management system with Notion and Calendar | Persistent workspace, Notion integration, Google Calendar, Telegram bot, Dynamic skillsets |
| [workflow-orchestrator](./workflow-orchestrator/) | Multi-step workflow execution with comprehensive tracing | Dynamic skillset loading, Multiple specialized skillsets, State persistence, Execution tracing |
| [skillset-based-dynamic-skill](./skillset-based-dynamic-skill/) | Agent that dynamically discovers and loads skills | Dynamic skill loading, Modular architecture, List & Install abilities |
| [ai-employee](./ai-employee/) | Comprehensive AI Employee for professional environments | Workspace/Space, Shell execution, Gmail integration, Notion integration, Dynamic skillsets |
| [simple-self-improving-agent](./simple-self-improving-agent/) | Self-improving agent that learns from interactions | Backstory file management, Read/Write abilities, Continuous learning |
| [mcp-factory](./mcp-factory/) | Factory-style architecture with multiple MCP server integrations | 4 service skillsets, Modular MCP servers, Clear separation of concerns |
| [dynamic-mcp-search-and-install](./dynamic-mcp-search-and-install/) | Agent that dynamically discovers and installs MCP servers | MCP registry search, On-demand installation, Self-extending capabilities |

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

### Agent Framework Example
Defines an autonomous agent as a project of files (`instructions.md`, `skills/*/SKILL.md`, `heartbeat.md`) and provisions the whole thing — agent, tools, workspace, skills, channels, and schedules — with Terraform. A reference architecture showing the ChatBotKit backend can build anything, entirely as infrastructure-as-code. The agent (named Atlas) is authored as a real project under `agent/` that Terraform uploads to the backend and wires up.

**What you'll learn:**
- Authoring an agent as a project of files, then provisioning it with Terraform
- Driving a bot's backstory from `instructions.md` via `file()`
- Giving an agent a rich toolset with ability packs (`pack/shell` for shell/file/import tools, `pack/cbk/space/skills` for skill discovery)
- File-based skills: `skills/*/SKILL.md` (+ bundled scripts) uploaded into a workspace under `.skills/` with `fileset` + `for_each`, then discovered and read by the agent's space-skill tools
- A Slack channel (with Discord/Teams/WhatsApp ready to uncomment)
- Two kinds of schedules: timed triggers and a frequent heartbeat whose instructions come from `heartbeat.md` (the trigger description)

**Use when:** You want an agent-framework developer experience managed as infrastructure-as-code, with the agent, its project files, and all its surfaces created and torn down reproducibly.

### Multi-Tenant Agents Examples
A pair of examples deploying a separate, isolated agent for each customer in that
customer's own ChatBotKit sub-account (a "partner user"). Both use **one** master
token plus the provider's `run_as` attribute (the `X-RunAs-UserId` header) to
target each customer's sub-account — no per-customer tokens, no `for_each`.

- **[multi-tenant-agents-shared](./multi-tenant-agents-shared/)** — the same agent for every customer, from one module, via provider aliases.
- **[multi-tenant-agents-per-customer](./multi-tenant-agents-per-customer/)** — a bespoke agent per customer, each in its own folder.

**What you'll learn:**
- The platform's sub-account (partner user) model for multi-tenant SaaS
- Using one master token + `run_as` to operate on many sub-accounts (like the AWS provider's `assume_role`)
- Provider aliases composing one reused module (shared) vs. a distinct module per customer (per-customer)
- Per-customer isolation via separate sub-accounts (`run_as`), all from a single shared state and one apply — no `for_each`

**Use when:** You are building a multi-tenant product where each customer needs their own isolated agent and resources, not a shared account.

### Dual-Agent Programmable Workflows Example
Two-agent architecture where a Workflow Architect programs custom scripts and a Task Runner executes them.

**What you'll learn:**
- Multi-agent collaboration patterns
- Shared file and space resources
- Asymmetric access control (architect writes, runner reads)
- Scheduled workflow execution
- Separation of concerns in agent design

**Use when:** Building systems where design and execution should be separate, or when you need maintainable automation frameworks.

### System Diagnostics Agent Example
Self-monitoring agent that introspects its own capabilities and produces diagnostic reports.

**What you'll learn:**
- Self-introspection using blueprint resource discovery
- Scheduled automated diagnostics
- Persistent diagnostic log storage
- Health monitoring patterns
- Agent self-documentation

**Use when:** Building production AI systems that need continuous monitoring and health checks.

### Second Brain Example
Personal knowledge management system integrating Notion, Google Calendar, and Telegram.

**What you'll learn:**
- Personal knowledge management architecture
- Notion integration for notes and databases
- Google Calendar integration for time awareness
- Telegram bot for mobile access
- Template-based OAuth2 secrets
- Multi-skillset organization

**Use when:** Building personal productivity assistants or knowledge management systems.

### Workflow Orchestrator Example
Multi-step workflow execution engine with dynamic skillset loading and comprehensive tracing.

**What you'll learn:**
- Complex multi-step workflow patterns
- Dynamic skillset loading for modular capabilities
- Workflow state persistence
- Comprehensive execution tracing
- Multiple specialized skillsets (Data, Control, Reporting)
- Audit trail generation

**Use when:** Building sophisticated automation workflows that require state management and full traceability.

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

### Simple Self-improving Agent Example
Demonstrates a self-improving AI agent that learns from interactions.

**What you'll learn:**
- File resource management for backstory storage
- Read and write abilities for file manipulation
- Self-improvement patterns
- Continuous learning architectures

**Use when:** Building agents that need to adapt and improve based on real-world usage.

### MCP Factory Example
Factory-style architecture with multiple independent MCP server integrations.

**What you'll learn:**
- Modular MCP server organization
- Multiple skillset patterns
- Service isolation and boundaries
- MCP server integration setup

**Use when:** Providing a suite of distinct AI functionalities to multiple clients or users.

### Dynamic MCP Search and Install Example
Agent that discovers and installs MCP servers on-demand from a registry.

**What you'll learn:**
- MCP registry search capabilities
- Dynamic MCP installation patterns
- Self-extending agent architectures
- On-demand capability loading

**Use when:** Building adaptable agents that discover and integrate tools as needed.

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
