# Scheduled task in a blueprint

Demonstrates the `chatbotkit_task` resource — the durable unit of agent
execution — and how a task is shipped as part of a `chatbotkit_blueprint`.

The configuration creates:

- a **blueprint** (`Daily Briefing Kit`) that groups the setup into one
  cloneable/exportable template,
- a **bot** (`Analyst`) assigned to the blueprint via `blueprint_id`,
- a **recurring task** (`Daily Briefing`) that runs the bot every weekday at
  09:00 Europe/London, bounded by `max_iterations` / `max_time` / `max_calls`
  and started fresh each run (`session_duration = 0`), and
- a **one-shot task** (`Kickoff`) with `schedule = "now"` that runs immediately
  on apply.

Because the bot and both tasks set `blueprint_id`, they belong to the same
blueprint — cloning or exporting `Daily Briefing Kit` carries the scheduled
automation along with the bot it drives.

## Usage

```bash
export CHATBOTKIT_API_KEY="your-api-key"

terraform init
terraform plan
terraform apply
```

## Scheduling

`schedule` accepts `now`, an interval keyword (`hourly`, `daily`, `weekly`, …),
a cron expression (as above), or a one-off date-time. Omit it (or use `never`)
for a task that only runs on demand. `timezone` is the IANA zone the schedule is
evaluated in.
