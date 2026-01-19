# MCP Factory Reference Architectures

This example demonstrates a factory-style architecture for MCP servers that expose multiple skillsets and abilities through separate MCP server integrations.

## Overview

This blueprint demonstrates a factory-style architecture for an MCP server that exposes multiple skillsets and abilities through separate MCP server integrations. Each skillset is designed to encapsulate specific functionalities, allowing for modular and organized management of AI capabilities.

## Architecture

```
┌──────────────────────────────────────────────────────────────────┐
│                    MCP Factory Architecture                       │
│                                                                   │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  │
│  │  Analytics MCP  │  │  Content MCP    │  │  Research MCP   │  │
│  │     Server      │  │     Server      │  │     Server      │  │
│  └────────┬────────┘  └────────┬────────┘  └────────┬────────┘  │
│           │                    │                    │            │
│  ┌────────▼────────┐  ┌────────▼────────┐  ┌────────▼────────┐  │
│  │   Analytics     │  │    Content      │  │    Research     │  │
│  │   Skillset      │  │   Skillset      │  │    Skillset     │  │
│  │                 │  │                 │  │                 │  │
│  │ • Analyze Data  │  │ • Create        │  │ • Search Info   │  │
│  │ • Generate      │  │ • Edit          │  │ • Summarize     │  │
│  │   Report        │  │                 │  │                 │  │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘  │
│                                                                   │
│  ┌─────────────────┐                                             │
│  │   Tasks MCP     │                                             │
│  │     Server      │                                             │
│  └────────┬────────┘                                             │
│           │                                                       │
│  ┌────────▼────────┐                                             │
│  │     Tasks       │                                             │
│  │   Skillset      │                                             │
│  │                 │                                             │
│  │ • Plan Tasks    │                                             │
│  │ • Track Progress│                                             │
│  └─────────────────┘                                             │
│                                                                   │
└──────────────────────────────────────────────────────────────────┘
```

## Key Features

1. **Modular Architecture**: Each skillset is independent and encapsulates specific functionality
2. **Separate MCP Servers**: Each skillset has its own MCP server integration for clear boundaries
3. **Scalable Design**: Easy to add new services without affecting existing ones
4. **Clear Separation**: Each service has a well-defined scope and responsibility

## Services

### Analytics Service
**Purpose**: Data analytics and reporting capabilities

**Abilities**:
- **Analyze Data**: Perform comprehensive data analysis
- **Generate Report**: Create detailed analytics reports

**Use Cases**:
- Business intelligence dashboards
- Data-driven decision making
- Performance metrics tracking

### Content Service
**Purpose**: Content creation and management

**Abilities**:
- **Create Content**: Generate new content based on specifications
- **Edit Content**: Modify and refine existing content

**Use Cases**:
- Marketing content generation
- Documentation management
- Content workflow automation

### Research Service
**Purpose**: Information gathering and synthesis

**Abilities**:
- **Search Information**: Find relevant information on topics
- **Summarize Findings**: Create concise summaries

**Use Cases**:
- Competitive analysis
- Market research
- Literature reviews

### Task Management Service
**Purpose**: Task planning and execution tracking

**Abilities**:
- **Plan Tasks**: Create structured task plans
- **Track Progress**: Monitor and report on progress

**Use Cases**:
- Project management
- Workflow coordination
- Progress reporting

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

5. Access the MCP servers:
   - Navigate to the ChatBotKit dashboard
   - Find each MCP server integration
   - Copy the MCP server URLs
   - Use these URLs in compatible MCP clients (Claude Desktop, etc.)

## Connecting to MCP Servers

After applying this configuration, each MCP server will be available through the ChatBotKit platform. To connect:

1. **From Claude Desktop or other MCP clients**:
   - Add the MCP server URL to your client configuration
   - Authenticate using your ChatBotKit credentials
   - The abilities from that skillset will be available to the client

2. **From ChatBotKit Bots**:
   - Assign the skillset to a bot
   - The bot will have access to the skillset's abilities
   - Multiple bots can share the same skillsets

## Customization

### Adding a New Service

To add a new service to the factory:

```hcl
resource "chatbotkit_skillset" "monitoring_service" {
  name        = "Monitoring Service"
  description = "System monitoring and alerting capabilities"
}

resource "chatbotkit_skillset_ability" "monitor_health" {
  skillset_id = chatbotkit_skillset.monitoring_service.id
  name        = "Check Health"
  description = "Monitor system health and status"
  instruction = <<-EOT
    # Monitoring logic here
  EOT
}

resource "chatbotkit_mcp_server_integration" "monitoring_mcp" {
  skillset_id = chatbotkit_skillset.monitoring_service.id
  name        = "Monitoring MCP Server"
  description = "MCP server exposing monitoring capabilities"
}
```

### Implementing Real Abilities

Replace placeholder logic with actual implementations:

```hcl
resource "chatbotkit_skillset_ability" "analytics_analyze" {
  skillset_id = chatbotkit_skillset.analytics_service.id
  name        = "Analyze Data"
  description = "Perform comprehensive data analysis"
  instruction = <<-EOT
    template: search/web
    parameters:
      query: $[query! ys|data analysis query]
  EOT
}
```

## When to Use This Pattern

This pattern is ideal when:
- You need to provide a suite of distinct AI functionalities
- Different teams or users need access to different capabilities
- You want clear boundaries between service domains
- Centralized management of AI resources is required
- Multiple client applications need to consume AI services

## Benefits

1. **Clear Boundaries**: Each service is self-contained with well-defined responsibilities
2. **Independent Scaling**: Services can be scaled independently based on demand
3. **Easy Maintenance**: Changes to one service don't affect others
4. **Flexible Access Control**: Can grant access to specific services per user/team
5. **Simplified Monitoring**: Each MCP server can be monitored separately

## Architecture Patterns

### Factory Pattern
This architecture follows the factory pattern where:
- Each skillset is a "product" in the factory
- Each MCP server is a "distribution channel" for that product
- New products (skillsets) can be added without changing existing ones

### Microservices Inspiration
While not true microservices, this architecture borrows concepts:
- Service isolation
- Independent deployment
- Clear service boundaries
- API-first design (via MCP protocol)

## Cleanup

To destroy all created resources:
```bash
terraform destroy
```

## Learn More

- [ChatBotKit MCP Documentation](https://chatbotkit.com/docs/integrations/mcp)
- [Model Context Protocol](https://modelcontextprotocol.io)
- [ChatBotKit Skillsets Documentation](https://chatbotkit.com/docs/resources/skillsets)
- [Blueprint Reference Architecture Examples](https://chatbotkit.com/blueprints)
