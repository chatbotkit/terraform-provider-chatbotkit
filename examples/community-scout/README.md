# Community Scout

A **product-led-growth agent** built on ChatBotKit and provisioned with Terraform.
It runs on a cycle to watch public conversations (Reddit today) for threads where
your product genuinely helps, drafts a useful, disclosed reply, and **suggests it to
your team on Slack** for a human to post. It earns attention by being helpful, not by
being loud.

## The idea

Same monitor-on-a-cycle shape as [`soc-investigator`](../soc-investigator), applied
to social listening — with clean boundaries:

- **Deterministic (script).** Fetching the Reddit firehose and deduping it into the
  mention store — done in a script so the bulk of that data never round-trips through
  the model (cheaper and faster). The agent reads only the summary.
- **Tool calls (read).** Reading a _specific_ thread the agent is judging, with the
  read-only Reddit tools. These are read-only — there is no Reddit-post tool.
- **Judgment (skills).** Whether a thread is worth engaging, and what a genuinely
  helpful reply says. This is where the agent thinks — and holds a high bar.
- **Hand-off (Slack).** It suggests the thread + draft to the team; a human posts.

## The gate is structural

Because this is an automated agent, the human-in-the-loop is built into the
**structure**, not just the prompt: the scout has read-only Reddit tools and a
Slack-suggest tool, and **no Reddit-post tool at all.** It physically cannot reply
on Reddit. The cycle's only outbound action is a Slack suggestion; a human reads it
and posts if they choose.

This matters more here than almost anywhere — auto-posting marketing to Reddit is
how accounts get shadow-banned and brands get burned. Removing the post tool turns
"please don't auto-post" from a hope into a guarantee.

## Architecture

```
        trigger (hourly)
              │
              ▼
     ┌─────────────────┐  reads .skills/ + watchlist.md    ┌────────────────────┐
     │ Community Scout │◀─────────────────────────────────▶│ workspace (space)  │
     │      (bot)      │  read-only reddit · runs scripts  │ mentions/          │
     └─────────────────┘                                   │ knowledge/ .skills/│
              │                                            └────────────────────┘
   ┌──────────┼───────────┬───────────┬───────────────┐
   ▼          ▼           ▼           ▼               ▼
 monitor    score       draft       suggest         learn
 (fetch+    (judgment)  (judgment)  (Slack hand-    (judgment →
  script)                            off to team)    knowledge/)
                                          │
                                          ▼
                                   #growth-leads  →  a human posts on Reddit
```

### The cycle

1. **Monitor** — `fetch_mentions.py` reads `watchlist.md`, searches Reddit directly
   (Atom feed + realistic User-Agent), normalizes the hits, and dedups them into
   `mentions/`. The agent just runs it and reads the summary — the firehose never
   enters its context.
2. **Score** (judgment) — rate genuine relevance; skip everything that does not clearly
   clear the bar (`relevance >= 70`).
3. **Draft** (judgment) — write a helpful, disclosed, non-spammy reply for the ones that
   clear it; `status: drafted`.
4. **Suggest** — send each drafted thread to the team on Slack (link, why, draft,
   suggested owner) with `slack/conversation/start`; `status: suggested`. The cycle
   does not post on Reddit.
5. **Learn** — turn the outcomes of threads the team posted into targeting notes under
   `knowledge/`.

The behavior lives in the files under [`agent/`](./agent) — instructions, the
`watchlist.md` config, the skills, and the script are the real asset.

## Reddit + Slack tooling

- **Monitoring fetch (script)** — `fetch_mentions.py` queries Reddit's public Atom
  search directly with a realistic User-Agent (the JSON endpoint and bot-looking
  User-Agents are blocked) and dedups, so the firehose never enters the model context.
- **Targeted reads (ability)** — `pack/reddit[read-only]` lets the agent read a
  _specific_ thread it is judging. Small and on-demand — here the content must enter
  context, because the model has to reason over it.
- **Slack suggest** — `slack/conversation/start` posts the suggestion to the team
  channel. Provide Slack credentials (variables in `main.tf`).
- **No Reddit write** — `reddit/comment/create` is deliberately not wired, so the
  scout cannot post. A human posts from the Slack suggestion.

## Files

```
main.tf                          bot + read-only reddit + slack suggest + workspace + triggers
agent/instructions.md            the scout backstory
agent/watchlist.md               product, keywords, subreddits, team/Slack routing, rules (edit me)
agent/heartbeat.md               the cycle tick (trigger description)
agent/skills/monitor/            SKILL.md + scripts/fetch_mentions.py    (fetch + dedup)
agent/skills/score/              SKILL.md                               (judgment: the bar)
agent/skills/draft/              SKILL.md                               (judgment: the reply)
agent/skills/suggest/            SKILL.md                               (Slack hand-off)
agent/skills/learn/              SKILL.md                               (judgment → knowledge)
```

## Usage

```bash
export CHATBOTKIT_API_KEY="..."
terraform apply \
  -var="slack_bot_token=xoxb-..." \
  -var="slack_signing_secret=..."
```

Edit `agent/watchlist.md` for your product and team, then talk to the bot and ask it
to "run one scouting cycle now." It searches Reddit live, dedups the hits, scores
them, drafts replies for the ones that clear the bar, and posts a suggestion per
thread into your Slack channel. A human picks it up and posts on Reddit.

## The watchlist is a living config

`watchlist.md` is read by the fetch script *and* the agent, so it lives as a file in
the workspace (not in the prompt). Terraform **seeds it once** and then leaves it
alone: the `chatbotkit_space_storage_file` resource uses `lifecycle { ignore_changes }`,
and the provider's storage `Read` is a no-op (write-mostly), so a later `terraform
apply` will **not** revert edits the team — or the agent — makes in the space. To push
a new baseline, edit it in the space, or drop `ignore_changes` for one apply. The
skills under `.skills/` are code, so they stay fully Terraform-managed and *are*
re-applied on change.

## Other sources (seams)

The pipeline only depends on the mention shape (`mentions/<id>.json`), so more sources
slot in:

- **LinkedIn / news** — first-class `linkedin` and `newsapi` ability catalogues can
  write mentions in the same shape.
- **X / Twitter** — no first-class catalogue; add an X API bearer secret + a
  `fetch/request[with-auth]` ability against the recent-search endpoint (see the
  commented block in `main.tf`), or approximate with `search_web` + `site:x.com`.

## Responsible engagement

The example is built to be welcome, not spammy: a high relevance bar, one tailored
reply per thread, affiliation disclosed every time, subreddit rules respected, and a
human in the loop by construction. That is both the ethical default and the only way
this kind of agent stays effective rather than banned.

## Related examples

- [`soc-investigator`](../soc-investigator) — the monitor-on-a-cycle pattern this mirrors.
- [`agent-framework`](../agent-framework) — the single agent as a project of files.
- [`deep-researcher`](../deep-researcher) — the orchestrator-worker upgrade for volume.
