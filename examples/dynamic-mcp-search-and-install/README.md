# Dynamic MCP Search and Install Architecture

This example demonstrates a reference architecture for an AI agent that can search for and install Model Context Protocol (MCP) servers dynamically.

## Overview

This blueprint demonstrates a dynamic architecture for an MCP server that allows the AI agent to search for and install Model Context Protocol (MCP) servers from a known registry. By combining search and installation capabilities, this blueprint enables agents to discover and integrate new MCP servers on-demand, creating a self-extending system that grows its functionality as needed.

## Architecture

```
┌──────────────────────────────────────────────────────────┐
│              Dynamic MCP Agent                            │
│                                                           │
│  ┌─────────────────────────────────────────────────────┐ │
│  │           MCP Discovery Skillset                    │ │
│  │  ┌──────────────────┐  ┌──────────────────┐        │ │
│  │  │  Search MCP      │  │  Install MCP     │        │ │
│  │  │   Servers        │  │  (by URL)        │        │ │
│  │  └──────────────────┘  └──────────────────┘        │ │
│  └─────────────────────────────────────────────────────┘ │
│                                                           │
│              ┌───────────────┐                           │
│              │ MCP Registry  │                           │
│              └───────┬───────┘                           │
│       ┌──────────────┼──────────────┬──────────────┐    │
│       │              │              │              │    │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐   │
│  │  File   │  │ GitHub  │  │ Database│  │  Slack  │   │
│  │ System  │  │   MCP   │  │   MCP   │  │   MCP   │   │
│  │  MCP    │  │         │  │         │  │         │   │
│  └─────────┘  └─────────┘  └─────────┘  └─────────┘   │
│                                                           │
│              (Installed on-demand as needed)              │
└──────────────────────────────────────────────────────────┘
```

## Key Features

1. **Dynamic Discovery**: Search for MCP servers by name or description in a registry
2. **On-Demand Installation**: Install MCP servers only when needed for specific tasks
3. **Self-Extending**: Agent gains new capabilities without manual reconfiguration
4. **Flexible Integration**: Connect to any MCP server available in the registry

## Core Abilities

### Search MCP Servers
Uses the `registry/server/search` template to find available MCP servers. The agent can:
- Search by server name
- Search by description or keywords
- Discover servers that provide needed functionality
- Review server metadata before installation

### Install MCP
Uses the `conversation/mcp/install[url]` template to activate an MCP server by its URL. Once installed:
- The server's tools become available to the agent
- Functions can be called directly in the conversation
- Capabilities persist for the duration of the conversation

## How It Works

### Discovery Flow

1. **User Request**: User asks for functionality the agent doesn't have
2. **Search**: Agent searches the MCP registry for relevant servers
3. **Selection**: Agent identifies the most appropriate MCP server
4. **Installation**: Agent installs the MCP server using its URL
5. **Utilization**: Agent uses the newly available tools to complete the task

### Example Interaction

**User**: "Can you help me analyze the files in my project directory?"

**Agent**:
1. Searches for MCP servers related to file system access
2. Finds a "filesystem" MCP server in the registry
3. Installs it using the server's URL
4. Uses the filesystem tools to list and analyze files
5. Provides the requested analysis to the user

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

5. Test the agent by:
   - Asking it to search for specific MCP servers
   - Requesting functionality that requires MCP installation
   - Observing how it discovers and integrates new capabilities

## Testing the Agent

### Search for MCP Servers
```
User: "What MCP servers are available for working with GitHub?"
Agent: [Searches registry and lists GitHub-related MCP servers]
```

### Install and Use an MCP Server
```
User: "I need to check my GitHub repositories"
Agent: [Searches for GitHub MCP, installs it, then uses its tools to list repositories]
```

## MCP Registry

The agent searches a registry of available MCP servers. Common MCP servers include:

- **Filesystem**: Access and manipulate local files
- **GitHub**: Interact with GitHub repositories and issues
- **PostgreSQL**: Query and manage PostgreSQL databases
- **Slack**: Send messages and interact with Slack
- **Google Drive**: Access and manage Google Drive files
- **Memory**: Persistent key-value storage
- **Fetch**: Make HTTP requests to external APIs

## Customization

### Adding Search Filters

Enhance the search ability with filters:

```hcl
resource "chatbotkit_skillset_ability" "search_mcp_advanced" {
  skillset_id = chatbotkit_skillset.mcp_discovery.id
  name        = "Advanced MCP Search"
  description = "Search MCP servers with filters"
  instruction = <<-EOT
    template: registry/server/search
    parameters:
      query: $[query! ys|search terms]
      category: $[category ys|server category (e.g., database, filesystem)]
  EOT
}
```

### Installing Multiple MCP Servers

The agent can install multiple MCP servers in a single conversation:

```hcl
# The agent can call the install ability multiple times
# Each installation adds to the available tools
```

## When to Use This Pattern

This pattern is ideal when:
- Agent needs diverse, specialized capabilities that vary by use case
- You want to avoid pre-installing all possible tools
- Requirements evolve and new MCP servers become available
- Users need access to different sets of integrations
- You want a self-service approach to capability extension

## Benefits

1. **Reduced Initial Complexity**: Start simple, add capabilities as needed
2. **Efficient Resource Usage**: Only load tools when required
3. **Future-Proof**: Automatically gain access to new MCP servers as they're added to the registry
4. **User-Driven**: Agent adapts to specific user needs dynamically
5. **Maintainable**: No need to reconfigure for new integrations

## Important Considerations

1. **Registry Quality**: Ensure the MCP registry contains trustworthy servers
2. **Installation Limits**: Consider rate limits or caps on installations per conversation
3. **Security**: Validate MCP servers before installation to prevent malicious code
4. **Performance**: Installing multiple large MCP servers may impact response times

## Comparison with Static MCP Configuration

| Aspect | Dynamic (This Pattern) | Static Configuration |
|--------|----------------------|---------------------|
| Setup Complexity | Low (just search/install abilities) | High (configure each MCP server) |
| Runtime Flexibility | High (discover new servers) | Low (fixed set of servers) |
| Resource Efficiency | High (load on-demand) | Low (all loaded upfront) |
| User Experience | Adaptive (fits needs) | Consistent (same tools always) |

## Cleanup

To destroy all created resources:
```bash
terraform destroy
```

## Learn More

- [Model Context Protocol Specification](https://modelcontextprotocol.io)
- [ChatBotKit MCP Documentation](https://chatbotkit.com/docs/integrations/mcp)
- [ChatBotKit Registry Documentation](https://chatbotkit.com/docs/registry)
- [Blueprint Reference Architecture Examples](https://chatbotkit.com/blueprints)
