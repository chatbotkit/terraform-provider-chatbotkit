Run one scouting cycle now.

1. **Monitor** — run the `monitor` skill: it runs `fetch_mentions.py`, which searches
   Reddit for the `watchlist.md` terms, normalizes the hits, and dedups new threads
   into `mentions/`. Read the summary.
2. **Score** — run the `score` skill over `status: new` mentions. Hold the bar high;
   skip everything that does not clearly clear it.
3. **Draft** — for mentions that clear the bar, run the `draft` skill to write a
   helpful, disclosed reply (`status: drafted`).
4. **Suggest** — run the `suggest` skill: send each drafted thread to the team on
   Slack (link, why, draft, suggested owner) and mark it `suggested`. You do not
   post on Reddit — a human does.
5. **Learn** — for threads that were posted and you can see outcomes for, run the
   `learn` skill to capture what landed.

If there is nothing new worth engaging, say so briefly and stop. Quality and
welcome over volume — a quiet cycle is a fine outcome.
