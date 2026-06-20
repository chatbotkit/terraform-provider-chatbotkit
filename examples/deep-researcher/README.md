# Deep Researcher

An autonomous **deep-research agent** built on ChatBotKit and provisioned with
Terraform, in the shape of the modern **orchestrator-worker** pattern — the same
architecture as systems like [gpt-researcher](https://github.com/assafelovic/gpt-researcher)
and Anthropic's multi-agent research system.

A "big brain" orchestrator takes a question, breaks it into focused
sub-questions, dispatches them to worker agents that run **in parallel**,
supervises and adapts as findings come back, then **synthesizes** everything into
one comprehensive, cited report.

## The idea

The interesting claim this example makes concrete: **deep research does not need
a hard-coded workflow graph.** The orchestration — fan out, wait, join, adapt,
synthesize — is the orchestrator _agent_ doing its job at runtime, not a static
pipeline declared in advance.

So Terraform declares only the **stable topology**: the agents, their tools, the
shared workspace, the entry surface. The **dynamic part** — the actual research
tasks — are runtime artifacts the agents create on demand. There is deliberately
no "task" resource in this file; tasks are spawned by agents, like conversations
are.

## Architecture

```
                 user
                  │  (chat widget)
                  ▼
            ┌───────────┐
            │  intake   │  small/fast: writes a clean brief,
            └─────┬─────┘  commissions a task for the orchestrator
                  │  task/create + task/run
                  ▼
          ┌────────────────┐
          │  orchestrator  │  big brain: plans sub-questions, fans out,
          └───┬───┬───┬────┘  supervises, adapts, synthesizes the report
              │   │   │  task/create + task/run  (parallel)
       ┌──────┘   │   └──────┐
       ▼          ▼          ▼
   ┌──────────┐ ┌──────────┐ ┌──────────┐
   │researcher│ │researcher│ │researcher│  workers: one sub-question each;
   └────┬─────┘ └──┬───────┘ └────┬─────┘  search + fetch + cite
        │          │              │
        └──────────┼──────────────┘
                   ▼
            shared workspace        findings/*.md  ← workers write
            (chatbotkit_space)      report.md      ← orchestrator writes
```

### How the orchestration actually works

- **Fan out:** the orchestrator calls `create_subtask` + `run_subtask`
  (`task/create` / `task/run`) for each sub-question. `run` executes a task
  immediately, so launching several without waiting _is_ the parallel fan-out.
- **Join / wait:** the orchestrator decides its own rhythm — it calls
  `list_subtasks` / `fetch_subtask` to introspect progress and waits between
  checks as it sees fit. A task runs as a durable workflow, so "waiting" is
  intrinsic, not a held process.
- **Adapt:** as findings arrive it can go deeper on rich threads, re-dispatch
  sharper questions, or prune dead ends — the thing a fixed graph can't do.
- **Results channel:** workers write cited findings to the shared workspace and
  return a short summary as their task result; the orchestrator reads both.
- **Synthesize:** the orchestrator reads all `findings/*.md` and writes a single
  `report.md`.

The behavior lives in the instruction files under [`agent/`](./agent) — those
are the real asset and the source of truth. Edit them and re-apply.

## Files

```
main.tf                                  topology: bots, skillsets, workspace, widget
agent/intake/instructions.md.tftpl       front-door agent (orchestrator id injected)
agent/orchestrator/instructions.md.tftpl big-brain agent (researcher id injected)
agent/researcher/instructions.md         worker agent
```

Bot ids are wired between agents with `templatefile()`, so the orchestrator knows
which bot to dispatch to and the intake knows which bot to commission — Terraform
resolves the dependency order automatically.

## Model tiering

Each bot carries its own `model`, so the team uses the right brain for each job:
`claude-4.5-opus` to orchestrate and synthesize, a lighter `claude-4.5-sonnet`
for the workers and the intake. Adjust per bot in `main.tf`. (Model names follow
the platform catalogue.)

## Usage

```bash
export CHATBOTKIT_API_KEY="..."
terraform init
terraform apply
```

Then open the widget (or talk to the intake bot) and ask for something that
benefits from real research, e.g. _"Compare the leading open-source web-research
agents and their trade-offs."_ The intake hands off, the orchestrator fans out
workers, and a cited report lands in the workspace as `report.md`.

## Notes and limits

- **Bounds.** "Research until satisfied" is bounded by per-task `maxTime` /
  `maxIterations` limits as a backstop. The intake commissions the orchestrator
  with a generous `maxTime` (the default is only 15 minutes — too short to
  supervise workers and write a report), and the orchestrator is instructed to
  stop at "good and well-sourced," not "exhaustive." Tune limits to taste.
- **Cost & variance.** The orchestrator-worker pattern is more capable on
  open-ended questions but burns more tokens and varies more run-to-run than a
  single bounded agent. For simpler questions, a single researcher bot with the
  same tools is cheaper.
- **What's emergent vs declared.** The wiring is declared; the _plan_, the
  fan-out width, the waiting rhythm, and the synthesis are emergent agent
  behavior. If you need _enforced_ invariants (e.g. a mandatory review-and-
  approve gate before a report is published), that guarantee belongs in a
  deterministic step, not in a prompt — out of scope for this example.

## Related examples

- [`agent-framework`](../agent-framework) — a single agent as a project of files.
- [`dual-agent-programmable-workflows`](../dual-agent-programmable-workflows) —
  two agents with asymmetric roles and a shared workspace.
- [`workflow-orchestrator`](../workflow-orchestrator) — dynamic skillset loading
  and workflow state.
