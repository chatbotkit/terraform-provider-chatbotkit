# Matillion Data Pipeline Operations Center

This example demonstrates an AI-powered operations center for Matillion Data Productivity Cloud that automates pipeline management, monitoring, and incident response.

## Overview

The Matillion Intelligent Data Pipeline Operations Center is a comprehensive agentic system designed to transform how data engineering teams manage their Matillion infrastructure. This system can save hundreds of hours per month by automating routine operations, incident response, and proactive monitoring.

## Architecture

```
┌──────────────────────────────────────────────────────────────┐
│           Matillion Operations Agent                         │
│                (Claude 4.5 Sonnet)                           │
│                                                              │
│  ┌────────────────────────────────────────────────────────┐ │
│  │            Main Operations Skillset                    │ │
│  │    ┌──────────────────┐  ┌──────────────────┐         │ │
│  │    │  List Skillsets  │  │ Install Skillset │         │ │
│  │    └──────────────────┘  └──────────────────┘         │ │
│  └────────────────────────────────────────────────────────┘ │
│                            │                                │
│      ┌─────────────────────┼─────────────────────┐         │
│      │                     │                     │         │
│  ┌───▼────────┐    ┌──────▼───────┐    ┌────────▼────┐   │
│  │  Pipeline  │    │   Schedule   │    │ Monitoring  │   │
│  │ Operations │    │  Management  │    │     &       │   │
│  │            │    │              │    │ Compliance  │   │
│  │  • List    │    │  • List      │    │  • Get      │   │
│  │  • Execute │    │  • Create    │    │    Consumption │
│  │  • Cancel  │    │  • Update    │    │  • Audit    │   │
│  │            │    │  • Delete    │    │    Events   │   │
│  └────────────┘    └──────────────┘    └─────────────┘   │
│                                                            │
│                    ┌──────────────┐                        │
│                    │Infrastructure│                        │
│                    │  Management  │                        │
│                    │              │                        │
│                    │  • List      │                        │
│                    │  • Commands  │                        │
│                    └──────────────┘                        │
│                            │                               │
│                    ┌───────▼────────┐                      │
│                    │ Matillion API  │                      │
│                    │     Secret     │                      │
│                    └────────────────┘                      │
│                                                            │
│  ┌─────────────────────────────────────────────────────┐  │
│  │         Daily Health Check Trigger                   │  │
│  │       (Scheduled automated monitoring)               │  │
│  └─────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────┘
```

## Key Features

### 1. Four Specialized Skillsets

#### Pipeline Operations
- Execute pipelines on-demand
- Monitor running executions and report status
- Cancel problematic or stuck runs
- Track execution history and identify patterns

#### Schedule Management
- Create pipeline schedules using cron expressions
- Modify existing schedules (timing, enabled/disabled)
- Delete obsolete schedules
- Review and optimize scheduling patterns

#### Monitoring & Compliance
- Track pipeline execution status across projects
- Monitor credit consumption and alert on anomalies
- Review audit events for compliance
- Track data lineage for governance

#### Infrastructure Management
- Monitor Matillion agent health and status
- Restart, pause, or resume agents
- Review configurations and recommend optimizations

### 2. Dynamic Skill Loading
- Main skillset enables discovery of available capabilities
- Install specific skillsets on-demand based on task needs
- Keeps context focused and efficient

### 3. Automated Health Checks
- Daily trigger integration for scheduled monitoring
- Generates operational summaries
- Proactive identification of issues

### 4. Professional Operations Agent
- Technical and precise communication style
- Follows operational guidelines and safety protocols
- Confirms before destructive operations
- Logs significant actions for audit

## Use Cases

This architecture is particularly valuable for:

- **Complex Pipeline Environments**: Multiple projects and pipelines across different environments
- **Operational Burden Reduction**: Automating routine monitoring and management tasks
- **Incident Response**: Quick diagnosis and resolution of pipeline issues
- **Cost Management**: Tracking and alerting on credit consumption
- **Compliance**: Audit trail maintenance and data lineage tracking
- **Team Collaboration**: Slack integration for team-wide notifications

## Usage

1. Set your ChatBotKit API key:
```bash
export CHATBOTKIT_API_KEY="your-api-key"
```

2. Configure Matillion API credentials:
   - Navigate to Secrets in ChatBotKit platform
   - Create a new secret of type `template`
   - Select `matillion` template
   - Enter your Matillion API credentials

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

## Operational Guidelines

### Before Executing Pipelines
1. Confirm project, environment, and pipeline name
2. Check for currently running executions
3. Verify agent status if critical
4. Execute and track with execution ID

### For Schedule Changes
1. List current schedules first
2. Validate cron expressions
3. Consider timezone implications
4. Confirm with requestor

### For Incident Response
1. Gather execution details and errors
2. Check for patterns in recent failures
3. Review agent status if issues persist
4. Escalate infrastructure problems

### For Cost Monitoring
1. Track consumption trends
2. Alert on usage anomalies
3. Recommend optimizations

## Safety Protocols

- Always confirm before deleting schedules or projects
- Warn before cancelling running executions
- Provide clear feedback on all operations
- Log significant actions for audit purposes

## Integration Options

### Slack Integration (Optional)
Uncomment the Slack integration block in `main.tf` and provide credentials:

```hcl
resource "chatbotkit_slack_integration" "matillion_slack" {
  bot_id           = chatbotkit_bot.matillion_ops.id
  name             = "Matillion Operations Bot"
  description      = "Slack bot for pipeline operations"
  signing_secret   = "your-slack-signing-secret"
  bot_token        = "your-slack-bot-token"
  session_duration = 0
}
```

Benefits:
- Team-wide visibility into operations
- Conversational interaction with pipelines
- Notifications for alerts and issues
- Natural language pipeline management

### Scheduled Health Checks
The included daily trigger integration automatically:
- Checks pipeline health
- Monitors execution status
- Reviews credit consumption
- Generates daily operational summaries

Customize the schedule by changing `trigger_schedule`:
- `"hourly"` - Every hour
- `"daily"` - Once per day
- `"weekly"` - Once per week
- `"quarterhourly"` - Every 15 minutes

## Customization

### Adding More Abilities

To extend capabilities, add abilities to existing skillsets:

```hcl
resource "chatbotkit_skillset_ability" "get_execution_status" {
  skillset_id = chatbotkit_skillset.pipeline_ops.id
  secret_id   = chatbotkit_secret.matillion_api.id
  
  name        = "Get Execution Status"
  description = "Retrieve status of a specific pipeline execution"
  instruction = <<-EOT
    template: matillion/pipeline-execution/fetch
    parameters: {}
  EOT
}
```

### Creating Custom Skillsets

Add domain-specific skillsets for specialized operations:

```hcl
resource "chatbotkit_skillset" "data_quality" {
  name        = "Data Quality Monitoring"
  description = "Abilities for monitoring data quality metrics"
}
```

## Best Practices

1. **Start with core skillsets**: Install only what you need for specific tasks
2. **Use scheduled triggers**: Automate routine health checks and reporting
3. **Enable Slack integration**: Improve team visibility and collaboration
4. **Review audit logs regularly**: Ensure compliance and track operations
5. **Monitor credit consumption**: Set up alerts for unusual usage patterns
6. **Document custom modifications**: Keep track of specialized abilities added

## Cleanup

To destroy all created resources:
```bash
terraform destroy
```

## Learn More

- [Matillion Data Productivity Cloud Documentation](https://docs.matillion.com/)
- [ChatBotKit Skillsets Documentation](https://chatbotkit.com/docs/resources/skillsets)
- [ChatBotKit Trigger Integrations](https://chatbotkit.com/docs/integrations/trigger)
- [Blueprint Examples](https://chatbotkit.com/blueprints)
