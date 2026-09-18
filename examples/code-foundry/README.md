# Code Foundry (multi-User)

An autonomous code foundry built on a **multi-User architecture**: one shared
tools User holds the expensive, sensitive machinery once, and each customer
gets a thin, isolated child User with a Coding Agent that borrows those tools
across Users and works on that customer's repository.

This pattern supports productising agents for many customers without duplicating
the toolbox or putting credentials such as a GitHub App key in each child User.

## Architecture

```
  shared User (alias: "shared")               child Users
  ┌──────────────────────────────────────────┐         ┌───────────────────────────┐
  │ GitHub bot - mints repo-scoped tokens    │ install │ alice/                    │
  │   (github/repository/token/create + JWT) │◀────────│   Coding Agent            │
  │ Coding Tools skillset                    │ @shared@│   + install shared tools  │
  │   = global-coding-tools (protected)      │ global- │   + heartbeat             │
  │   bash · rw · import · mint-token ·      │ coding- │   + context (repo/vercel) │
  │   design list/read · skills list/read    │ tools   ├───────────────────────────┤
  │ Design space  ·  Coding space            │         │ bob/  (same module)       │
  │ Designs Manager bot + Sync trigger       │         └───────────────────────────┘
  └──────────────────────────────────────────┘
```

One API token belonging to the parent User (`CHATBOTKIT_API_TOKEN`) operates on every
child User. Each is selected with a provider alias and `run_as` (the
`X-RunAs-UserId` header), using the same multi-User mechanism as the
[`multi-tenant-agents-shared`](../multi-tenant-agents-shared)
example.

## The two halves

- **`modules/shared`** - the shared User. A GitHub bot that mints
  repository-scoped GitHub App tokens (JWT secret), a **Coding Tools** skillset
  exposed to child Users as `global-coding-tools` (`visibility = protected` plus a
  stable `alias`), shared Design and Coding spaces, and a Designs Manager bot with
  a Sync trigger. The shared User **must have the alias `shared`**.
- **`modules/coder`** - one per child User. A Coding Agent whose only built-in ability is
  to install the shared toolset across Users (`conversation/skillset/install` with
  `@shared@global-coding-tools`), plus a **heartbeat** and the child User's
  **context**.

## How a token gets minted (and why context matters)

The coding agent never holds the GitHub App key. To touch a repo it calls
`mint_github_repo_token`, which `bot/apply`s the **shared GitHub bot**. The GitHub
bot signs an App JWT (from the shared secret) and returns a short-lived,
repository-scoped token.

But _which_ repository? It cannot be hard-coded because each agent belongs to a
different child User. The repository comes from the User's **context**
(`githubOwner` / `githubRepo` / `vercelProjectId`). That is why setting context
per child User is part
of the setup, and why the context API needed a programmatic surface (see below).

## The two gaps this example closes

1. **Context via GraphQL to a native Terraform resource.** The User context
   API (`/api/v1/user/{userId}/context/...`) was REST-only. It is now also
   exposed via GraphQL (`contexts` query; `createContext` / `updateContext` /
   `deleteContext` mutations, scoped to the selected child User), and the
   Terraform provider has been regenerated from the updated schema. Context is
   therefore a native **`chatbotkit_context`** resource. The coder module creates
   it for the child User selected by its provider, so each agent is scoped to its
   own repository without hard-coding.
2. **A coding-agent heartbeat.** Coding tasks span many steps. Each coder
   child User has a recurring heartbeat trigger that nudges the agent to make the
   next concrete step on its active task, reusing one conversation within the
   session window so it keeps its place across ticks.

## Files

```
main.tf                       provider aliases (shared + per-User) + module calls
modules/shared/main.tf        the shared tools User (GitHub, Coding Tools, spaces, Designs Manager)
modules/coder/main.tf         one child User's coding agent + shared tools + heartbeat + context
terraform.tfvars.example      User IDs, GitHub App ID/key, git email
```

## Usage

```bash
export CHATBOTKIT_API_TOKEN="<parent User API token>"
export TF_VAR_github_app_private_key="$(cat github-app.pem)"
cp terraform.tfvars.example terraform.tfvars   # fill in User IDs + App ID
terraform init
terraform apply
```

Child Users are created out of band; the shared User must be aliased `shared`.
To add a customer: add a User ID variable, a provider alias, and a
`./modules/coder` call (mirroring `alice`/`bob`).

## Notes and seams

- **GitHub App JWT secret.** The shared secret is `type = "jwt"` with the App key;
  its structured claims are JSON-encoded into the `config` map (`config` is
  `map(string)` in the provider). Adjust to your provider version if the JWT config
  shape differs.
- **Context resource.** `chatbotkit_context` is generated from the GraphQL schema;
  its `payload` is `map(string)`, so values (repo owner/name, repo URL, Vercel
  project ID) are flat strings. Regenerate the provider (the GraphQL-to-stubs
  pipeline) if you change the schema again.
- **Designs Manager** Sync trigger is `schedule = "never"` (run on demand) and syncs
  an upstream design repo into the shared Design space.

## Related examples

- [`multi-tenant-agents-shared`](../multi-tenant-agents-shared) - the provider-alias plus `run_as` multi-User pattern this builds on.
- [`agent-framework`](../agent-framework) - the single agent as a project of files.
