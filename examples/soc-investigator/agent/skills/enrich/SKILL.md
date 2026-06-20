---
description: Enrich a case's IOCs with threat-intelligence reputation (script-backed).
---

# Enrich

Attach threat-intelligence reputation to a case's indicators. The lookup itself
is a deterministic script; you decide which IOCs to enrich and how the result
changes your read of the case.

## Steps

1. For each IOC in the case (ip, domain, url, hash), run the lookup from the
   workspace root:

   ```bash
   python .skills/enrich/scripts/ti_lookup.py 185.220.101.7
   ```

2. The script returns a verdict (`malicious` / `suspicious` / `benign`), a
   reputation score, and any feed pulses.
3. Write the results into the case file under `enrichment` (keyed by IOC), so the
   `investigate` step can cite them.

## Notes

- The lookup is a deterministic offline mock so the example runs without
  credentials. Set `TI_API_TOKEN` and wire a real provider (AlienVault OTX,
  VirusTotal, ...) in `ti_lookup.py` to go live — the case shape does not change.
- A single malicious IOC rarely settles a case on its own; combine reputation
  with the behavior in the alerts before you reach a verdict.
