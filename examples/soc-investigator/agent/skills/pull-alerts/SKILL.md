---
description: Pull new SIEM alerts and correlate them into open cases (deterministic).
---

# Pull Alerts

The first step of every cycle. This is deterministic work — a script does it, you
do not reason about individual alerts here. Run it, read the summary, move on.

## Steps

1. From the workspace root, run the puller with the shell tools:

   ```bash
   python .skills/pull-alerts/scripts/pull_alerts.py
   ```

2. Read the JSON summary it prints: how many alerts were pulled, how many new
   cases were opened, how many were correlated into existing cases, and how many
   duplicates were skipped.
3. The script writes open cases to `cases/<id>.json`. That is the case store for
   the rest of the cycle.

## Notes

- Correlation: alerts of the same rule, on the same asset, within the same hour
  collapse into one case (a stable correlation UID). This is the noise-reduction
  step.
- Idempotent: re-running the same window does not duplicate cases or alerts, so
  it is safe to run on a schedule.
- Stdlib only. The script reads `data/sample-alerts.json` (a mock SIEM). In
  production, replace its `read_alerts()` with a real SIEM query — nothing else
  changes.
