# Complete ChatBotKit Terraform Example

This example demonstrates a complete ChatBotKit setup using Terraform, including:

- **Dataset**: A knowledge base that stores information the bot can search
- **Skillset**: A container for abilities that extend bot capabilities
- **Ability**: A specific skill attached to the skillset
- **Bot**: The main AI assistant linked to both the dataset and skillset
- **Trigger Integration**: Allows the bot to be invoked via webhooks or schedules

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         Bot                                  │
│  ┌─────────────────────────────────────────────────────────┐│
│  │ Name: AI Assistant                                      ││
│  │ Backstory: Helpful AI assistant with knowledge access   ││
│  └─────────────────────────────────────────────────────────┘│
│                    │                    │                    │
│         ┌─────────┴─────────┐   ┌──────┴──────┐             │
│         ▼                   ▼   │              │             │
│  ┌──────────────┐    ┌─────────────────┐      │             │
│  │   Dataset    │    │    Skillset     │      │             │
│  │──────────────│    │─────────────────│      │             │
│  │ Knowledge    │    │   Bot Skills    │      │             │
│  │ Base         │    │─────────────────│      │             │
│  └──────────────┘    │ ┌─────────────┐ │      │             │
│                      │ │   Ability   │ │      │             │
│                      │ │─────────────│ │      │             │
│                      │ │ Search      │ │      │             │
│                      │ │ Knowledge   │ │      │             │
│                      │ └─────────────┘ │      │             │
│                      └─────────────────┘      │             │
│                                               │             │
│                                    ┌──────────┴──────────┐  │
│                                    │ Trigger Integration │  │
│                                    │────────────────────│  │
│                                    │ Bot Trigger        │  │
│                                    └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## Prerequisites

1. [Terraform](https://www.terraform.io/downloads.html) >= 1.0
2. A ChatBotKit account and API key

## Usage

1. Set your ChatBotKit API key:

```bash
export CHATBOTKIT_API_KEY="your-api-key"
```

Or configure it directly in `main.tf`:

```hcl
provider "chatbotkit" {
  api_key = "your-api-key"
}
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

5. When done, clean up resources:

```bash
terraform destroy
```

## Resource Dependencies

The resources in this example have the following dependencies:

1. **Dataset** - No dependencies
2. **Skillset** - No dependencies
3. **Ability** - Depends on Skillset
4. **Bot** - Depends on Dataset and Skillset
5. **Trigger Integration** - Depends on Bot

Terraform automatically handles the creation order based on these dependencies.

## Customization

### Adding Records to the Dataset

After creating the dataset, you can add records to populate it with knowledge:

```hcl
resource "chatbotkit_dataset_record" "faq_1" {
  dataset_id  = chatbotkit_dataset.knowledge_base.id
  name        = "FAQ: How to reset password"
  description = "Instructions for password reset"
  text        = "To reset your password, go to the login page and click 'Forgot Password'..."
}
```

### Adding More Abilities

You can add multiple abilities to the skillset:

```hcl
resource "chatbotkit_skillset_ability" "email_ability" {
  skillset_id = chatbotkit_skillset.bot_skills.id
  name        = "Send Email"
  description = "Send an email notification"
  instruction = "Use this ability to send email notifications to users."
}
```

### Scheduling the Trigger

Configure the trigger to run on a schedule:

```hcl
resource "chatbotkit_trigger_integration" "scheduled_trigger" {
  name             = "Daily Report Trigger"
  description      = "Runs daily to generate reports"
  bot_id           = chatbotkit_bot.assistant.id
  trigger_schedule = "daily"
}
```

## Outputs

After applying, you'll receive the IDs of all created resources:

- `dataset_id` - The ID of the knowledge base dataset
- `skillset_id` - The ID of the bot skills skillset
- `ability_id` - The ID of the search ability
- `bot_id` - The ID of the AI assistant bot
- `trigger_integration_id` - The ID of the trigger integration

Use these IDs to reference the resources in other configurations or API calls.
