# Agent Framework on ChatBotKit

This example defines an autonomous agent as a small **project of files** and
provisions the whole thing — the agent, its tools, its workspace, its skills, its
channels, and its schedules — end to end with Terraform. It's a reference
architecture: the ChatBotKit backend provides the primitives, and Terraform wires
them together. Build anything on top of it.

The agent project lives in [`agent/`](./agent/). `terraform apply` creates the
backend resources and uploads the project files. The files are the source of
truth: edit them and re-apply to update the live agent.

## The agent project

```
agent/
├── instructions.md              # the always-on system prompt (the backstory)
├── heartbeat.md                 # instructions for each recurring "tick"
└── skills/                      # on-demand skills, uploaded to the workspace
    ├── web-research/
    │   └── SKILL.md
    └── data-cleanup/
        ├── SKILL.md
        └── scripts/
            └── summarize_csv.py # a skill can bundle scripts the agent runs
```

## How it maps to ChatBotKit

| Concept             | What it is here                           | ChatBotKit resource(s)                                            |
| ------------------- | ----------------------------------------- | ----------------------------------------------------------------- |
| `instructions.md`   | always-on system prompt                   | `chatbotkit_bot.backstory` (via `file(...)`)                      |
| abilities           | the agent's toolset (ability packs + web) | `chatbotkit_skillset` + `chatbotkit_skillset_ability`             |
| workspace           | isolated, persistent filesystem           | `chatbotkit_space`                                                |
| `skills/*/SKILL.md` | on-demand skills (+ bundled scripts)      | uploaded via `chatbotkit_space_storage_file` (`.skills/`)         |
| `heartbeat.md`      | the recurring-tick instructions           | the `description` of a heartbeat `chatbotkit_trigger_integration` |
| channels            | the surfaces the agent runs on            | `*_integration` (Slack active; others ready to uncomment)         |
| schedules           | timed runs + a heartbeat                  | `chatbotkit_trigger_integration` (cron)                           |

## Architecture

```
   instructions.md ─────────────▶ Atlas (chatbotkit_bot)
   (backstory)                         │ skillset
                                       ▼
                          Abilities (chatbotkit_skillset)
                ┌──────────────────────┴────────────────────────┐
                │ shell  (pack/shell: exec · rw · import)       │
                │ skills (pack/cbk/space/skills: list · read)   │
                │ search_web · fetch_url                        │
                └──────────────────────┬────────────────────────┘
                          shell / skills │ (scoped to the workspace)
                                         ▼
                            Workspace (chatbotkit_space)
                         ┌─────────────────────────────────┐
                         │ .skills/web-research/SKILL.md   │  ◀── uploaded
                         │ .skills/data-cleanup/SKILL.md   │      by Terraform
                         │ .skills/data-cleanup/scripts/   │
                         │ state/ (the agent's own files)  │
                         └─────────────────────────────────┘

   channels   ── Slack (active) · Discord · Teams · WhatsApp (ready) ──▶ Atlas
   schedules  ── Daily Briefing · Weekly Review · Heartbeat (every 5 min) ──▶ Atlas
```

## How the parts work

### `instructions.md` → backstory

The bot's `backstory` is read directly from `agent/instructions.md` with
Terraform's `file()` function, so the always-on system prompt is just a Markdown
file in the project.

### Abilities (the toolset)

Defined in Terraform and grouped into one skillset attached to the agent. Two of
them are ability **packs** — a single ability that installs a whole set of tools:

- **`pack/shell`** → run commands, read/write files, and import URLs, scoped to
  the workspace.
- **`pack/cbk/space/skills`** → list and read the agent's skills from the
  workspace.

Plus standalone **`search_web`** and **`fetch_url`**. Packs are the easy way to
give an agent a rich, coherent toolset without declaring each tool by hand.

### Skills (loaded on demand)

Each skill is a folder with a `SKILL.md` (YAML frontmatter `description:` + a
Markdown body) and, optionally, bundled `scripts/`. Terraform uploads the whole
`skills/` tree into the workspace under `.skills/`. At runtime the agent uses the
space-skill tools to list the available skills and read the relevant `SKILL.md`,
keeping the base prompt lean. The `data-cleanup` skill shows a skill bundling a
real, runnable script that the agent invokes with the shell tools.

### Channels (the surfaces)

The same bot is exposed through channel integrations. **Slack is active** in this
example; Discord, Teams, and WhatsApp are included as commented blocks — uncomment
the integration (and its variables) and supply credentials to activate them.

### Schedules — two kinds

1. **Normal schedules** run the agent at specific times for a defined job
   (`Daily Briefing`, `Weekly Review`).
2. **Heartbeat** is a frequent tick (every 5 minutes by default) that wakes the
   agent with no user present. Its instructions are the `heartbeat.md` file,
   passed straight through as the trigger's `description` — a cheap, continuous
   "is anything up?" loop.

## Usage

1. Set your ChatBotKit API key:

   ```bash
   export CHATBOTKIT_API_KEY="your-api-key"
   ```

2. (Optional) Provide Slack credentials to activate the channel:

   ```bash
   export TF_VAR_slack_bot_token="xoxb-..."
   export TF_VAR_slack_signing_secret="..."
   ```

3. Initialize, review, and apply:

   ```bash
   terraform init
   terraform plan
   terraform apply
   ```

## Customization

### Add a skill

Create `agent/skills/<your-skill>/SKILL.md` (and any `scripts/`). It is uploaded
to `.skills/<your-skill>/` on the next apply — no Terraform changes needed; the
agent discovers it with the space-skill tools.

### Add an ability

Add another `chatbotkit_skillset_ability` to the `Atlas Abilities` skillset —
either a single tool or another pack.

### Activate another channel

Uncomment the relevant `*_integration` block (and its variables) and supply
credentials.

### Tune the heartbeat

Edit the `heartbeat` trigger's `schedule` (e.g. `"* * * * *"` for every minute)
and the `heartbeat.md` instructions.

## When to use this pattern

- You want an agent-framework developer experience (a project of files: prompt,
  abilities, skills, channels, schedules) managed as infrastructure-as-code.
- You want the agent and all of its surfaces created, versioned, and torn down
  reproducibly across environments.

## Cleanup

```bash
terraform destroy
```

## Learn More

- [ChatBotKit Documentation](https://chatbotkit.com/docs)
- [ChatBotKit Terraform Provider](https://registry.terraform.io/providers/chatbotkit/chatbotkit/latest/docs)
