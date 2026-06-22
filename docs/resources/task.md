---
page_title: "chatbotkit_task Resource - terraform-provider-chatbotkit"
subcategory: ""
description: |-
  Manages a ChatBotKit Task resource.
---

# chatbotkit_task (Resource)

Manages a ChatBotKit Task — an autonomous, bounded run of a bot, optionally on a schedule. A task ties a bot to an objective (its `description`), runs it under hard limits (`max_iterations`, `max_time`), and can recur on a `schedule`. It is the durable unit of agent execution behind background work and heartbeats.

Tasks are commonly created at runtime by agents (via the task abilities); this resource is for the cases where you want to declare a task as infrastructure — e.g. a recurring scheduled job for a bot, managed and versioned in Terraform.

## Example Usage

### A recurring scheduled task

```terraform
resource "chatbotkit_bot" "analyst" {
  name      = "Analyst"
  backstory = "You produce a concise daily briefing."
}

resource "chatbotkit_task" "daily_briefing" {
  name        = "Daily Briefing"
  description = "Produce the weekday morning briefing and post it to the team."
  bot_id      = chatbotkit_bot.analyst.id

  schedule = "0 9 * * 1-5" # 09:00 on weekdays
  timezone = "Europe/London"

  max_time       = 900000 # 15 minutes
  max_iterations = 200
}
```

### A one-shot task that runs immediately

```terraform
resource "chatbotkit_task" "kickoff" {
  name        = "Kickoff"
  description = "Set up the project scaffold."
  bot_id      = chatbotkit_bot.analyst.id
  schedule    = "now"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional) The name of the task.
- `description` - (Optional) The task's objective — what the bot should do when it runs.
- `bot_id` - (Optional) The ID of the bot the task runs.
- `contact_id` - (Optional) The ID of a contact to scope the task to.
- `schedule` - (Optional) When the task runs: `now`, a cron expression, a date-time, or an interval keyword such as `daily`. Omit (or `never`) for a manual task.
- `timezone` - (Optional) The IANA timezone the schedule is evaluated in (e.g. `Europe/London`).
- `max_iterations` - (Optional) Maximum reasoning iterations per execution (clamped to platform limits).
- `max_time` - (Optional) Maximum wall-clock time per execution, in milliseconds (clamped to platform limits).
- `session_duration` - (Optional) Session duration in milliseconds controlling conversation reuse across runs (`0` starts a fresh conversation each run).
- `meta` - (Optional) A map of metadata key-value pairs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier of the task.
- `created_at` - The timestamp when the task was created.
- `updated_at` - The timestamp when the task was last updated.

## Import

Tasks can be imported using their ID:

```bash
terraform import chatbotkit_task.example task_abc123def456
```
