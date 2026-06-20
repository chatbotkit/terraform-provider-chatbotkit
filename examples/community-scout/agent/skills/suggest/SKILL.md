---
description: Suggest a thread + draft to the team on Slack for a human to post.
---

# Suggest

The hand-off. You do not post on Reddit — you have no tool to do so. Instead, you
send the team a Slack suggestion with everything a human needs to decide and act in
seconds.

## Steps

1. For each mention with a draft (`status: drafted`), use the `suggest_to_team`
   tool (`slack/conversation/start`) to message the team channel from
   `watchlist.md` (or DM the suggested owner).
2. Make the Slack message scannable and self-contained:
   - the subreddit + a one-line why-it's-relevant,
   - the **link** to the thread,
   - the suggested **draft reply** (so a human can copy, edit, and post),
   - a suggested owner — @mention the team member who covers that topic per the
     watchlist, with a light "want to take this?",
   - the relevance score and any risk flags (subreddit rules, tone).
3. Set the mention `status` to `suggested` and record where/whom you suggested it
   to, so the next cycle does not re-suggest the same thread.

## Notes

- One Slack message per thread; do not batch unrelated threads into a wall of text.
- You suggest; a human posts. That is the whole design — the gate is structural,
  not a matter of you choosing to hold back.
- If Slack is not configured, leave the mention at `drafted` and say so; do not
  attempt any other way to publish.
