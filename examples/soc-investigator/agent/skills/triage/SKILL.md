---
description: Decide which open cases are worth investigating, and in what order.
---

# Triage

This is judgment, not a script. After pulling alerts, look at the open cases and
decide where to spend investigation effort. The goal is to protect analyst time —
not every case deserves a full investigation.

## Steps

1. List the open cases:

   ```bash
   ls cases/ && cat cases/*.json
   ```

2. For each case with `status: open` and no `report` yet, weigh:
   - **severity** and alert count (a correlated burst is more urgent than a
     single low alert),
   - **asset** (a finance or domain controller host outranks a dev box),
   - **novelty** — check `knowledge/` for a known pattern or a known
     false-positive signature that already explains it.
3. Pick the top cases to investigate this cycle (a small number — quality over
   coverage). For a case that a knowledge record clearly marks as a benign known
   pattern, set `status: dismissed` with a one-line `report` saying why, and skip
   the full investigation.

## Notes

- Do not change a case's disposition beyond `open` → `dismissed` for clear known
  false positives. Anything that looks like a real threat goes to investigation,
  never straight to closure.
