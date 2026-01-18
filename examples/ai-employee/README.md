# AI Employee Reference Architecture

This example demonstrates a comprehensive AI Employee designed to work autonomously within professional environments.

## Overview

The AI Employee is a digital team member engineered to operate within professional settings. It can perform tasks such as project management, research, content creation, email management, and problem-solving, all while maintaining professional standards and protocols.

## Architecture

```
┌──────────────────────────────────────────────────────────────┐
│                      AI Employee Bot                         │
│                  (Claude 4.5 Opus Model)                     │
│                                                              │
│  ┌────────────────────────────────────────────────────────┐  │
│  │                 Core Skillset                          │  │
│  │  ┌──────────────┐  ┌──────────────┐  ┌─────────────┐   │  │
│  │  │ List Skills  │  │Install Skills│  │   Shell     │   │  │
│  │  └──────────────┘  └──────────────┘  │  Execution  │   │  │
│  │                                      └─────────────┘   │  │
│  └────────────────────────────────────────────────────────┘  │
│                            │                                 │
│       ┌────────────────────┴────────────────────┐            │
│       │                                         │            │
│  ┌─────────────────┐                  ┌──────────────────┐   │
│  │ Mail Management │                  │ Notion Management│   │
│  │  Skillset       │                  │  Skillset        │   │
│  │─────────────────│                  │──────────────────│   │
│  │ • List Messages │                  │ • Search Pages   │   │
│  │ • Fetch Message │                  │ • Fetch Page     │   │
│  │ • Search Email  │                  │ • Create Items   │   │
│  │ • Send Email    │                  │                  │   │
│  └─────────────────┘                  └──────────────────┘   │
│         │                                     │              │
│  ┌──────▼──────┐                      ┌───────▼─────────┐    │
│  │ Google Mail │                      │     Notion      │    │
│  │   Secret    │                      │     Secret      │    │
│  │  (Personal) │                      │   (Personal)    │    │
│  └─────────────┘                      └─────────────────┘    │
│                                                              │
│  ┌────────────────────────────────────────────────────────┐  │
│  │              AI Employee Workspace                     │  │
│  │         (Secure Sandbox Environment)                   │  │
│  └────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────┘
```

## Key Features

### 1. Professional AI Employee

- Comprehensive professional backstory defining roles and responsibilities
- Operates as a competent colleague rather than a subservient assistant
- Maintains professional tone and ethical standards

### 2. Secure Workspace

- Dedicated Space for sandboxed operations
- Shell execution capability for command-line tasks
- Isolated environment for file management

### 3. Dynamic Skill Loading

- Can list available skillsets
- Install skills on-demand based on task requirements
- Modular architecture for easy capability expansion

### 4. Email Management

- Gmail integration via OAuth2
- List, search, fetch, and send emails
- Personal-kind secret for end-user authentication

### 5. Knowledge Management

- Notion integration for document access
- Search pages, fetch content, create database items
- Personal-kind secret for user-specific access

## Core Capabilities

### Research and Information Retrieval

- Gather data, market insights, and industry trends
- Perform competitive intelligence gathering
- Verify information from multiple sources

### Document Management

- Create and edit business documents
- Generate reports and presentations
- Maintain documentation standards

### Communication Facilitation

- Draft professional correspondence
- Manage email communications
- Provide project updates

### Data Analysis and Project Support

- Process and analyze datasets
- Track tasks and timelines
- Monitor progress and identify roadblocks

## Usage

1. Set your ChatBotKit API key:

```bash
export CHATBOTKIT_API_KEY="your-api-key"
```

2. Configure OAuth2 secrets (done via ChatBotKit platform):

   - Google Mail: Navigate to Secrets → Create → Template → `platform/google/mail`
   - Notion: Navigate to Secrets → Create → Template → `platform/notion`

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

## Personal Kind Secrets

Both the Google Mail and Notion secrets use `kind = "personal"`, which means:

- End users authenticate with their own accounts
- The AI Employee acts on behalf of the authenticated user
- Each user's data remains isolated and secure
- Users must grant OAuth permissions before the AI Employee can access their data

## Customization

### Adding New Skillsets

To extend the AI Employee's capabilities:

```hcl
resource "chatbotkit_skillset" "calendar_skillset" {
  name        = "Calendar Management"
  description = "Manage calendar events and scheduling"
}

resource "chatbotkit_secret" "google_calendar" {
  name        = "Google Calendar OAuth2 Token"
  description = "OAuth2 token for accessing Google Calendar"
  type        = "template"
  kind        = "personal"

  config = jsonencode({
    template = "platform/google/calendar"
  })
}

resource "chatbotkit_skillset_ability" "list_events" {
  skillset_id = chatbotkit_skillset.calendar_skillset.id
  secret_id   = chatbotkit_secret.google_calendar.id

  name        = "List Calendar Events"
  description = "List upcoming calendar events"
  instruction = <<-EOT
    template: google/calendar/event/list
    parameters: {}
  EOT
}
```

### Customizing the Backstory

The backstory can be customized to fit your organization's:

- Specific workflows and protocols
- Industry terminology and standards
- Compliance requirements
- Communication styles

## Integration Options

This AI Employee can be integrated with various channels:

### Slack Integration

```hcl
resource "chatbotkit_slack_integration" "employee_slack" {
  bot_id           = chatbotkit_bot.ai_employee.id
  name             = "AI Employee Slack Bot"
  signing_secret   = "your-slack-signing-secret"
  bot_token        = "your-slack-bot-token"
  session_duration = 0
}
```

### Widget Integration

```hcl
resource "chatbotkit_widget_integration" "employee_widget" {
  bot_id           = chatbotkit_bot.ai_employee.id
  name             = "AI Employee Widget"
  session_duration = 0
}
```

## Security Considerations

- Shell execution is sandboxed within the Space
- OAuth2 secrets are personal-kind for user-level authentication
- The AI Employee respects confidentiality and privacy
- All external integrations require explicit user consent

## When to Use This Pattern

This pattern is ideal for:

- Organizations needing autonomous digital team members
- Environments requiring multi-modal task automation
- Teams wanting AI assistance with professional workflows
- Use cases requiring secure, sandboxed operations
- Scenarios needing dynamic capability expansion

## Cleanup

To destroy all created resources:

```bash
terraform destroy
```

## Learn More

- [ChatBotKit Bots Documentation](https://chatbotkit.com/docs/resources/bots)
- [ChatBotKit Spaces Documentation](https://chatbotkit.com/docs/resources/spaces)
- [ChatBotKit Secrets Documentation](https://chatbotkit.com/docs/resources/secrets)
- [Blueprint Examples](https://chatbotkit.com/blueprints)
