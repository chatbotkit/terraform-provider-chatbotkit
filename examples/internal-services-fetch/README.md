# Internal Services Assistant — `fetch` + OAuth & shared secrets

An agent that reaches your **internal** corporate services with the `fetch`
action, authenticated by two different kinds of secret. It shows the full
distinction between **machine-to-machine** access and **acting on behalf of a
user**, in one realistic setup.

```
internal-services-fetch/
├── main.tf                      # secrets, skillset + fetch abilities, bot, Slack channel
└── terraform.tfvars.example     # secret values to copy into terraform.tfvars
```

## The two kinds of secret

| Secret | `kind` | `type` | Who it represents | Used for |
|--------|--------|--------|-------------------|----------|
| `Acme Internal Service Token` | `shared` | `bearer` | The agent itself, as a service account | Read-only platform data everyone may see |
| `Acme SSO (Employee)` | `personal` | `oauth` | The individual employee chatting | Anything scoped to that person |

**Shared secret — agent ↔ internal system.** One account-wide bearer token. The
agent presents it to internal systems that trust it as a service account; there
is no human in the loop and every conversation uses the same identity. Good for
service-catalog/health lookups.

**Personal OAuth secret — the user authenticates.** Resolved per end-user. The
first time an employee triggers an ability that uses it, ChatBotKit asks them to
sign in through your SSO and stores *their* token; the agent then calls
permission-enforcing internal systems **as that employee**, so the backend
applies its normal access control.

## How `fetch` references secrets

The abilities use the **structured** fetch format — a `!fetch` YAML action tag
with typed helpers rather than inline string templates:

| Helper | Meaning |
|--------|---------|
| `!reference SECRET_DEFAULT` | the secret linked to the ability via `secret_id` |
| `!reference SECRET_<NAME>` | any other secret in the account, matched by name (e.g. `!reference SECRET_ACME_INTERNAL_SERVICE_TOKEN` → the secret named *Acme Internal Service Token*). Names match case-insensitively, with non-word characters collapsed to `_`. |
| `!string { name, description, optional, enum }` | a parameter the model fills in at call time |

A single ability can reference `SECRET_DEFAULT` **and** a named secret, which is
how `open_incident` uses both at once.

> **No `Bearer` prefix.** A `bearer` or `oauth` secret already resolves to a
> complete Authorization value *including* its scheme (`Bearer <token>`), so the
> header is just `Authorization: !reference SECRET_...`. Adding your own `Bearer `
> would produce `Bearer Bearer <token>`. (The structured format also means there
> are no `${...}` sequences in the instruction, so no Terraform `$$` escaping is
> needed.)

## The abilities

| Ability | Secret(s) | What it does |
|---------|-----------|--------------|
| `lookup_service_status` | shared (`SECRET_DEFAULT`) | Reads internal platform catalog as a service — no sign-in |
| `list_my_tickets` | personal (`SECRET_DEFAULT`) | Lists the signed-in employee's own helpdesk tickets |
| `get_my_timeoff_balance` | personal (`SECRET_DEFAULT`) | Reads the employee's PTO balance (user inferred from token) |
| `open_incident` | **both** — shared via `SECRET_ACME_INTERNAL_SERVICE_TOKEN` + personal via `SECRET_DEFAULT` | Goes through the internal API gateway (machine token) and attributes the incident to the employee (their forwarded OAuth token) |

`open_incident` is the key one: the gateway authenticates the **caller** (the
agent service) via the standard `Authorization` header (shared token), and the
incident is attributed to a real person via the employee's own token forwarded
in `X-Forwarded-Authorization` — both supplied as bearer values by the secrets.

## Usage

```bash
export CHATBOTKIT_API_KEY="sk-...your-api-key..."

cp terraform.tfvars.example terraform.tfvars   # fill in the secret values
terraform init
terraform apply
```

The `slack_*` variables are optional — leave them empty to skip the Slack
channel and just create the bot, skillset and secrets.

## Trying it

- Ask *"is billing-api healthy?"* → uses the **shared** token, answers
  immediately.
- Ask *"what are my open tickets?"* → the first time, the agent relays an Acme
  SSO sign-in prompt; once you authenticate it answers **as you**.
- Ask *"open a high severity incident for billing-api"* → exercises **both**
  secrets in a single request.

## Adapting it

- Point the URLs at your own internal services and the `config` endpoints at
  your own identity provider (the example uses a Keycloak-style realm).
- Keep secret values out of source control — supply them via `terraform.tfvars`
  or `TF_VAR_*` environment variables (all the variables are marked `sensitive`).
- The OAuth `config` keys are camelCase (`clientId`, `clientSecret`,
  `authorizationUrl`, `tokenUrl`, `grantType`, `scope`).
- For OAuth servers that publish metadata, you can replace the explicit
  `authorizationUrl`/`tokenUrl` with a single `resourceUrl` and let the
  endpoints be discovered (RFC 9728 / RFC 8414).
