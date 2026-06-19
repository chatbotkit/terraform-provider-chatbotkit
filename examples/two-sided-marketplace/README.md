# Two-Sided Marketplace — one agent per side of a market

A reference architecture for a **two-sided market**, modelled the way a
recruitment agency works: one side has people to place, the other has roles to
fill, and the market clears between them. It deploys **two agents**, one per
side, over a shared external matching engine.

```
two-sided-marketplace/
├── main.tf                      # secrets, two skillsets + fetch abilities, two bots, channels
└── terraform.tfvars.example     # secret values to copy into terraform.tfvars
```

## The two sides

| Agent | Side | Works for | Does |
|-------|------|-----------|------|
| **Talent Partner** | supply | the candidate | browse open roles, submit candidacy, track submissions |
| **Client Partner** | demand | the hiring company | post roles, search talent, review & **advance** submissions |

The two agents have **opposite loyalties** and never talk to each other. They
meet only through the marketplace, the same way a real agency's candidate desk
and client desk meet over the same pipeline.

## The core is an API, not the agents

The matching engine — the role index, the talent pool, submissions, and the
placement/fee transaction — lives in **your own marketplace service**. The
agents reach it with the `fetch` action and hold none of that logic themselves.
ChatBotKit is the conversational + identity tier on each side; your API is the
shared system of record.

```
   candidate ──► Talent Partner ─┐                        ┌─ Client Partner ◄── employer
                  (supply side)   │  fetch                 │   (demand side)
                                  ▼                        ▼
                       ┌──────────────────────────────────────────┐
                       │        YOUR MARKETPLACE API (core)        │
                       │  roles · talent pool · submissions · fee  │
                       └──────────────────────────────────────────┘
```

## The market clears through shared state

There is no agent-to-agent message. The loop runs through the core's data:

1. Talent Partner calls `submit_candidacy` → a submission is created in the core.
2. It surfaces in Client Partner's `review_submissions` for that role.
3. Client Partner calls `advance_submission` → the core records the stage change
   (interview / offer / placement) and any fee, **server-side**.

`advance_submission` is the market-clearing transaction. The agent *requests*
it; the API *decides and records* it. A placement is never something a model
does on its own — that keeps consent, auditability, and billing in the core
where they belong.

## Identity: one shared gateway token + per-side on-behalf OAuth

| Secret | `kind` | `type` | Represents | Used by |
|--------|--------|--------|------------|---------|
| `Marketplace Service Token` | `shared` | `bearer` | the platform, as a service | both desks (gateway auth) |
| `Talent Identity` | `personal` | `oauth` | the signed-in candidate | Talent Desk actions |
| `Client Identity` | `personal` | `oauth` | the signed-in employer contact | Client Desk actions |

Both desks present the **same** shared token to get past the marketplace
gateway, so the core sees one stable caller. Each desk attributes actions to a
**different** human via its own personal OAuth secret, so the core enforces what
that person may see and do. Same gateway, different sides — that is the whole
identity story of a two-sided market in two secrets.

### How the abilities reference secrets

Every ability reads the same way, on both sides:

| Header | Value | Meaning |
|--------|-------|---------|
| `Authorization` | `!reference SECRET_MARKETPLACE_SERVICE_TOKEN` | always the shared gateway token, by name |
| `X-On-Behalf-Of` | `!reference SECRET_DEFAULT` | this desk's human — present only on actions attributed to a person |

`SECRET_DEFAULT` is the secret linked to the ability via `secret_id`. On the
Talent Desk that is the candidate; on the Client Desk it is the employer. Within
one skillset it always means the same person, so the convention never flips. An
ability links exactly one secret (its `SECRET_DEFAULT`); to use a *second* secret
in the same request you reference it **by name** (`SECRET_<NAME>`, upper-cased
with non-word characters collapsed to `_`) — which is why the action abilities
link the personal identity and name the shared token, while the read-only browse
ability links only the shared token.

> **No `Bearer` prefix.** A `bearer`/`oauth` secret resolves to a complete
> Authorization value *including* its scheme, so the header is just
> `Authorization: !reference SECRET_...`.

## Usage

```bash
export CHATBOTKIT_API_KEY="sk-...your-api-key..."

cp terraform.tfvars.example terraform.tfvars   # fill in the secret values
terraform init
terraform apply
```

The `client_slack_*` variables are optional — leave them empty to skip the Slack
channel and just create the bots, skillsets and secrets.

## Trying it

- **As a candidate** (Talent Partner): *"find me remote senior backend roles"* →
  `search_open_roles` (shared token only, no sign-in). Then *"apply me to role
  R-128"* → first time, you sign in; `submit_candidacy` puts you forward **as
  you**.
- **As an employer** (Client Partner): *"show submissions for role R-128"* →
  `review_submissions` returns the candidate the Talent Partner just submitted.
  Then *"move them to interview"* → `advance_submission` clears that side of the
  market.

## Adapting it

- **Point the URLs at your own service** and the OAuth `config` endpoints at your
  own identity provider. The core API is yours to implement — these abilities are
  just the typed surface the agents call.
- **Make it proactive** with a `chatbotkit_trigger_integration` per bot: the
  Talent Partner sweeps for new matching roles on a schedule; the Client Partner
  sweeps for new matching talent. (Drive per-user fan-out from your backend, not
  one global trigger — a shared agent serves many users.)
- **Add a broker desk-lead** as an optional third bot if you want a human-style
  mediator: a `bot/call` ability (set `bot_id` on the ability) lets one agent
  consult it. Keep the actual placement/fee in the API, not in the broker's
  reasoning.
- **Give the demand side a stronger model** if its judgement work (screening,
  advancing) warrants it — set a different `model` on `client_partner`.
- **Enforce consent in the core**, not the agent: `search_talent` should only
  ever return candidates who opted in, and the on-behalf token is how the core
  knows which employer is asking.
