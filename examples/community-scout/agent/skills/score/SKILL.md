---
description: Score new mentions for genuine relevance; drop the ones not worth engaging.
---

# Score

Decide where contributing is genuinely useful. This is the spam-avoidance bar, and
it matters more than reach. Most threads should not get a reply.

## Steps

1. For each mention with `status: new`, read it (and, if useful, pull the thread
   discussion with `reddit/feed/post/comment/list` to understand the context).
2. Judge fit against `watchlist.md`:
   - Is the person asking a question the product genuinely answers, or describing a
     problem it genuinely solves?
   - Would a reply add real value even if the product did not exist?
   - Is it the right venue (subreddit rules allow it; not a competitor's space; not
     an old/closed thread; not a rant looking for agreement)?
3. Set `relevance` on the mention to a score 0–100 with a one-line reason, and move
   `status` to `scored`. Set `status: skipped` (with the reason) for anything that
   does not clearly clear the bar.

## The bar

Only `relevance >= 70` proceeds to a draft. When in doubt, skip — a missed thread
costs nothing; a tone-deaf or salesy reply costs trust and risks a ban. Never
engage purely because a keyword matched.
