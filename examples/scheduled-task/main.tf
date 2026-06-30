terraform {
  required_providers {
    chatbotkit = {
      source = "chatbotkit/chatbotkit"
    }
  }
}

provider "chatbotkit" {
  # api_key = "..." # Or set the CHATBOTKIT_API_KEY env var
}

# ----------------------------------------------------------------------------
# A blueprint groups a working setup — a bot plus the tasks that drive it — into
# a single template that can be cloned and exported as one unit.
# ----------------------------------------------------------------------------

resource "chatbotkit_blueprint" "briefing_kit" {
  name        = "Daily Briefing Kit"
  description = "An analyst bot plus the scheduled tasks that drive it."
  visibility  = "private"
}

# The bot the tasks run. Setting blueprint_id makes it part of the blueprint.
resource "chatbotkit_bot" "analyst" {
  name      = "Analyst"
  backstory = "You produce a concise, sourced morning briefing for the team."

  blueprint_id = chatbotkit_blueprint.briefing_kit.id
}

# ----------------------------------------------------------------------------
# A recurring task runs the bot on a schedule, bounded by hard execution caps.
# Setting blueprint_id ships the task as part of the blueprint, so cloning the
# blueprint brings the scheduled automation along with the bot it runs.
# ----------------------------------------------------------------------------

resource "chatbotkit_task" "daily_briefing" {
  name        = "Daily Briefing"
  description = "Produce the weekday morning briefing and post it to the team channel."

  bot_id       = chatbotkit_bot.analyst.id
  blueprint_id = chatbotkit_blueprint.briefing_kit.id

  schedule = "0 9 * * 1-5" # 09:00 on weekdays
  timezone = "Europe/London"

  # Per-run execution bounds.
  max_iterations   = 200
  max_time         = 900000 # 15 minutes, in milliseconds
  max_calls        = 2000
  session_duration = 0 # a fresh conversation each run
}

# A one-shot task in the same blueprint that runs immediately on apply.
resource "chatbotkit_task" "kickoff" {
  name        = "Kickoff"
  description = "Draft the first briefing so the team has a sample to react to."

  bot_id       = chatbotkit_bot.analyst.id
  blueprint_id = chatbotkit_blueprint.briefing_kit.id

  schedule = "now"
}

output "blueprint_id" {
  value = chatbotkit_blueprint.briefing_kit.id
}

output "daily_briefing_task_id" {
  value = chatbotkit_task.daily_briefing.id
}
