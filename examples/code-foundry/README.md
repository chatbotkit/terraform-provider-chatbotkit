# Code Foundry (multi-account)

An autonomous code foundry built on a **multi-account architecture**: one
**shared "tools" account** holds the expensive, sensitive machinery once, and each
**user gets a thin, isolated sub-account** with a Coding Agent that borrows the
shared tools cross-account and works on that user's own repository.

This is the pattern for productising agents to many users without duplicating the
toolbox per user and without putting credentials (a GitHub App key) in each
sub-account.

## Architecture

```
  shared account (partner-user alias: "shared")        per-user sub-accounts
  ┌──────────────────────────────────────────┐         ┌───────────────────────────┐
  │ GitHub bot — mints repo-scoped tokens    │ install │ alice/                    │
  │   (github/repository/token/create + JWT) │◀────────│   Coding Agent            │
  │ Coding Tools skillset                    │ @shared@│   + install shared tools  │
  │   = global-coding-tools (protected)      │ global- │   + heartbeat             │
  │   bash · rw · import · mint-token ·      │ coding- │   + context (repo/vercel) │
  │   design list/read · skills list/read    │ tools   ├───────────────────────────┤
  │ Design space  ·  Coding space            │         │ bob/  (same module)       │
  │ Designs Manager bot + Sync trigger       │         └───────────────────────────┘
  └──────────────────────────────────────────┘
```

One partner/master token (`CHATBOTKIT_API_KEY`) operates on every account; each is
selected with a provider alias + `run_as` (the `X-RunAs-UserId` header) — the same
multi-account mechanism as the [`multi-tenant-agents-shared`](../multi-tenant-agents-shared)
example.

## The two halves

- **`modules/shared`** — the shared account. A GitHub bot that mints
  repository-scoped GitHub App tokens (JWT secret), a **Coding Tools** skillset
  exported account-wide as `global-coding-tools` (`visibility = protected` + a
  stable `alias`), shared Design and Coding spaces, and a Designs Manager bot with
  a Sync trigger. The shared account's partner user **must have the alias `shared`**.
- **`modules/coder`** — one per user. A Coding Agent whose only built-in ability is
  to install the shared toolset cross-account (`conversation/skillset/install` with
  `@shared@global-coding-tools`). Plus a **heartbeat** and the user's **context**.

## How a token gets minted (and why context matters)

The coding agent never holds the GitHub App key. To touch a repo it calls
`mint_github_repo_token`, which `bot/apply`s the **shared GitHub bot**. The GitHub
bot signs an App JWT (from the shared secret) and returns a short-lived,
repository-scoped token.

But _which_ repository? It can't be hard-coded — each agent belongs to a different
user, and they may interact with the agent, so a hard-coded repo would be a
security risk. The repo comes from the sub-user's **context** (`githubOwner` /
`githubRepo` / `vercelProjectId`). That's why setting context per sub-user is part
of the setup, and why the context API needed a programmatic surface (see below).

## The two gaps this example closes

1. **Context via GraphQL → a native Terraform resource.** The partner-user context
   API (`/api/v1/partner/user/{userId}/context/...`) was REST-only. It is now also
   exposed via GraphQL (`contexts` query; `createContext` / `updateContext` /
   `deleteContext` mutations, scoped to the run_as'd sub-user), and the Terraform
   provider has been regenerated from the updated schema — so context is a native
   **`chatbotkit_context`** resource. The coder module uses it directly; created in
   the user's sub-account (the module's provider is run_as'd to it), so each agent is
   scoped to its own repo with no hard-coding.
2. **A coding-agent heartbeat.** Coding tasks span many steps. Each coder
   sub-account has a recurring heartbeat trigger that nudges the agent to make the
   next concrete step on its active task, reusing one conversation within the
   session window so it keeps its place across ticks.

## Files

```
main.tf                       provider aliases (shared + per-user) + module calls
modules/shared/main.tf        the shared "tools" account (GitHub, Coding Tools, spaces, Designs Manager)
modules/coder/main.tf         one user's coding agent + install-shared-tools + heartbeat + context
terraform.tfvars.example      account IDs, GitHub App id/key, git email
```

## Usage

```bash
export CHATBOTKIT_API_KEY="<partner/master token>"
export TF_VAR_github_app_private_key="$(cat github-app.pem)"
cp terraform.tfvars.example terraform.tfvars   # fill in account IDs + app id
terraform init
terraform apply
```

Accounts (partner users) are created out of band; the shared one must be aliased
`shared`. To add a user: add an account-id variable, a provider alias, and a
`./modules/coder` call (mirroring `alice`/`bob`).

## Notes and seams

- **GitHub App JWT secret.** The shared secret is `type = "jwt"` with the App key;
  its structured claims are JSON-encoded into the `config` map (`config` is
  `map(string)` in the provider). Adjust to your provider version if the JWT config
  shape differs.
- **Context resource.** `chatbotkit_context` is generated from the GraphQL schema;
  its `payload` is `map(string)`, so values (repo owner/name, repo URL, Vercel
  project id) are flat strings. Regenerate the provider (the GraphQL → stubs
  pipeline) if you change the schema again.
- **Designs Manager** Sync trigger is `schedule = "never"` (run on demand) and syncs
  an upstream design repo into the shared Design space.

## Related examples

- [`multi-tenant-agents-shared`](../multi-tenant-agents-shared) — the provider-alias + `run_as` multi-account pattern this builds on.
- [`agent-framework`](../agent-framework) — the single agent as a project of files.
