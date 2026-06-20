Run one SOC cycle now.

1. **Pull** — run the `pull-alerts` skill to bring in new alerts and correlate
   them into cases. Read the summary.
2. **Triage** — run the `triage` skill over the open, un-investigated cases and
   pick the ones worth working this cycle. Dismiss clear known false positives.
3. **Investigate** — for each selected case, follow the `investigate` skill:
   enrich its IOCs (`enrich` skill), pull knowledge and context, and write a
   structured verdict report back into the case. Leave it at `pending_review`.
4. **Learn** — for any case a human has approved since last cycle, run the
   `extract-knowledge` skill to capture the lesson.

Respect the approval gate: never close a case or recommend executing remediation
as if it were done. If there is nothing new to do, say so briefly and stop —
do not re-investigate cases that already have a report.
