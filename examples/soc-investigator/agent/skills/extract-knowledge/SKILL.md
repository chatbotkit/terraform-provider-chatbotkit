---
description: Distill a resolved case into reusable knowledge for future investigations.
---

# Extract Knowledge

How the SOC gets smarter over time. When a case has been resolved (a human has
approved its disposition), capture the reusable lesson so the next similar case is
faster to judge.

## Steps

1. For a case that reached a final disposition, decide what is worth keeping: a
   detection pattern, a false-positive signature, a confirmed-malicious
   indicator, or a remediation that worked.
2. Write it as a short markdown record under `knowledge/`, named for the pattern:

   ```bash
   # via the shell write tool
   knowledge/edr-vssadmin-shadow-delete.md
   ```

   Each record should carry: the pattern/rule it applies to, the signals that
   distinguished true from false positive, the verdict reached, and the
   remediation. Keep it tight and searchable.
3. Future investigations read `knowledge/` during triage and investigation, so
   the value compounds.

## Notes

- Only extract from cases that reached a human-approved disposition — do not learn
  from your own un-reviewed verdicts, or errors will reinforce themselves.
- For a production system at scale, back the knowledge base with a dataset
  (vector search) instead of flat files; the read/write contract is the same.
