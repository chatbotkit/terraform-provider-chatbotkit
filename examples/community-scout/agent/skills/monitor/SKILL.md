---
description: Fetch new Reddit threads matching the watchlist and persist them, deduped.
---

# Monitor

The first step of every cycle. This is deterministic and done by a script — you do
not fetch threads yourself. Pulling the firehose through your context just to filter
it would be slow and expensive; the script does the fetch and the dedup, and you
read only the small summary.

## Steps

1. Run the fetcher from the workspace root:

   ```bash
   python .skills/monitor/scripts/fetch_mentions.py
   ```

   It reads `watchlist.md` (keywords + subreddits), searches Reddit directly,
   normalizes the hits, and writes new threads to `mentions/<id>.json`, skipping
   ones already seen.
2. Read the JSON summary it prints: how many threads were pulled, how many are new,
   how many were duplicates, and any errors.
3. The new mentions (`status: new`) are what you triage next.

## Notes

- Idempotent: re-running never resurfaces a thread already recorded, so you never
  re-score or re-suggest the same conversation.
- It uses Reddit's public Atom search with a realistic User-Agent (the plain JSON
  endpoint is blocked). It needs network egress from the shell sandbox; if a search
  fails it is logged and skipped, not retried blindly.
- Other sources (LinkedIn, a news API, or X via the X API) can write mentions in the
  same shape; the rest of the pipeline does not change.
