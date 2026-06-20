---
description: Investigate one case and write a structured, cited verdict report.
---

# Investigate

The core agentic playbook. Take one triaged case and produce a structured report:
what happened, whether it is a real threat, and what to do. This is where your
judgment matters.

## Steps

1. Read the case: `cat cases/<id>.json`. Note the rule, asset, user, processes,
   and IOCs across its alerts.
2. Pull context:
   - search `knowledge/` for prior cases of the same pattern (known-good or
     known-bad signatures, past remediations),
   - enrich the IOCs (see the `enrich` skill) to get reputation,
   - use web search / fetch for unfamiliar tooling, CVEs, or attacker techniques
     referenced by the alert.
3. Form a verdict and write it back into the case file under `report`, using this
   structure:

   - `verdict`: one of `malicious`, `suspicious`, `benign`, `inconclusive`
   - `confidence`: `low` | `medium` | `high`
   - `attack_chain`: ordered steps of what happened (with the alert/IOC each step
     rests on)
   - `iocs`: the indicators, each with its enrichment verdict
   - `severity`: confirmed severity after investigation
   - `remediation`: concrete recommended actions
   - `summary`: a few sentences a human can read first

4. Set the case `status` to `pending_review` (never `closed`). Every claim in the
   report must trace to an alert, an enrichment result, or a cited source.

## The approval gate (do not skip)

You investigate and recommend. You do **not** close cases or execute remediation
(isolating a host, disabling a user, blocking an IOC). Those are gated on human
approval. Leave the case at `pending_review` and surface your recommendation. This
is the one deterministic guarantee of the workflow.
