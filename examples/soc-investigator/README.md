# SOC Investigator

An autonomous **security-operations agent** built on ChatBotKit and provisioned
with Terraform. It models the agentic half of an open-source SOC platform (in the
shape of projects like [agentic-soc-platform](https://github.com/FunnyWolf/agentic-soc-platform)):
an analyst that runs on a cycle to **pull** SIEM alerts, **correlate** them into
cases, **triage**, **investigate** the ones that matter, **enrich** indicators,
and **accumulate knowledge** so it gets better over time.

## The idea

The example exists to make one boundary concrete: **what should be deterministic
code, and what should be agent judgment.**

- **Deterministic → scripts.** Pulling alerts, normalizing them, correlating them
  into cases, and looking up threat intel are high-volume, must-be-reliable work.
  They are stdlib Python scripts the agent runs with the shell tools. No model
  touches them.
- **Judgment → skills.** Which cases deserve attention, what actually happened,
  whether it is a real threat, what to recommend — these are `SKILL.md` playbooks
  the agent reads and follows. This is where it reasons.

A real SOC ingests a _stream_ of alerts; here that becomes a **scheduled pull** on
a cycle (which fits triggers cleanly). The correlation UID is what makes polling
safe — re-pulling the same window never duplicates a case or an alert.

## Architecture

```
        trigger (every 15 min)
                │
                ▼
        ┌───────────────┐      reads .skills/*  ┌──────────────────────┐
        │  SOC Analyst  │◀────────────────────▶ │  workspace (space)   │
        │     (bot)     │   runs scripts, r/w   │  data/  cases/       │
        └───────────────┘                       │  knowledge/ .skills/ |
                │                               └──────────────────────┘
   ┌────────────┼─────────────┬──────────────┬───────────────┐
   ▼            ▼             ▼              ▼               ▼
 pull-alerts  triage     investigate      enrich      extract-knowledge
 (script)    (judgment)  (judgment+gate)  (script)    (judgment → knowledge/)
```

### The cycle

1. **Pull** (`pull-alerts`, deterministic) — `pull_alerts.py` polls the mock SIEM,
   normalizes, correlates by a stable UID, and writes open cases to `cases/`.
2. **Triage** (judgment) — pick the cases worth working; dismiss clear known
   false positives (checked against `knowledge/`).
3. **Investigate** (judgment + gate) — for a case: enrich its IOCs, pull knowledge
   and web context, and write a structured verdict report (verdict, attack chain,
   IOCs, severity, remediation) back into the case file.
4. **Enrich** (`enrich`, script-backed) — `ti_lookup.py` returns reputation for an
   IOC; the agent decides how it changes the read.
5. **Learn** (`extract-knowledge`, judgment) — distill resolved cases into
   `knowledge/` so future investigations are faster.

The behavior lives in the files under [`agent/`](./agent) — the instructions, the
skills, and the scripts are the real asset. Edit them and re-apply.

## The approval gate

The analyst investigates and **recommends**; it does not close cases or execute
remediation (isolating a host, disabling a user, blocking an IOC). Investigated
cases are left at `status: pending_review` for a human to approve. This is the one
deterministic guarantee around the agent's autonomy — exactly the kind of invariant
that belongs in structure rather than a prompt.

## Files

```
main.tf                              bot + skillset + workspace + uploads + cycle triggers
agent/instructions.md                the analyst backstory
agent/heartbeat.md                   the cycle tick (trigger description)
agent/data/sample-alerts.json        mock SIEM source (swap for a real SIEM)
agent/skills/pull-alerts/            SKILL.md + scripts/pull_alerts.py   (deterministic)
agent/skills/triage/                 SKILL.md                            (judgment)
agent/skills/investigate/            SKILL.md                            (judgment + gate)
agent/skills/enrich/                 SKILL.md + scripts/ti_lookup.py     (script-backed)
agent/skills/extract-knowledge/      SKILL.md                            (judgment)
```

## Usage

```bash
export CHATBOTKIT_API_KEY="..."
terraform init
terraform apply
```

The cycle trigger runs every 15 minutes. To watch it work, talk to the analyst bot
directly and ask it to "run one SOC cycle now" — it will pull the sample alerts,
correlate them (the two vssadmin alerts on `WIN-FIN-07` collapse into one critical
case), triage, and investigate.

## Production seams

Each seam is marked in the code; the case shape never changes when you swap one:

- **Real SIEM** — replace `read_alerts()` in `pull_alerts.py` with an ELK/Splunk
  query. Everything downstream depends only on the normalized alert shape.
- **A SIRP case DB** — here cases are JSON files in the space. Point reads/writes
  at a real incident-tracking system (via `fetch` abilities or an MCP server) for
  a shared, durable system of record.
- **Threat intel** — set `TI_API_TOKEN` and wire a provider (AlienVault OTX,
  VirusTotal, ...) into `ti_lookup.py`.
- **Vector knowledge** — back `knowledge/` with a dataset (semantic search) instead
  of flat files at scale; the read/write contract is the same.

## Upgrade path

This is a single analyst on a cycle. For high case volume, promote it to the
**orchestrator-worker** shape (see [`deep-researcher`](../deep-researcher)): an
orchestrator triages and dispatches one investigation **task** per hot case, and
the workers run in parallel. Same skills, same scripts — the orchestration becomes
runtime task fan-out.

## What it borrows from agentic-soc-platform

- The pipeline: ingest → correlate → triage → investigate → enrich → learn.
- Correlation UID + time-bucket dedup (what makes a polling pull idempotent).
- Knowledge accumulation from resolved cases (gets smarter over time).
- Threat-intel enrichment as a discrete step.
- The deterministic-modules / agentic-playbooks split — the lesson this example
  is built to show.

## Related examples

- [`agent-framework`](../agent-framework) — the single agent as a project of files.
- [`deep-researcher`](../deep-researcher) — the orchestrator-worker pattern this
  upgrades into.
