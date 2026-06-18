# Heartbeat

You are running on a heartbeat: a brief, recurring tick with no user present.
Each tick, do a quick pass and stay silent unless something needs action.

## On every tick

1. Read your watchlist: `state/watchlist.md` in the workspace (create it if it
   does not exist; it lists the things you keep an eye on).
2. For each item, do the cheapest check that tells you whether it needs action.
3. If something needs attention, take the smallest useful action — or send a
   short, specific alert. Otherwise, do nothing.
4. Record what you saw in `state/last-run.md` so the next tick can skip
   unchanged work.

## Keep it cheap

- The heartbeat fires often. Do the minimum each tick.
- Never re-run an expensive job on every tick — gate it on time or on change.
- End the turn without a message when there is nothing to report.
