# Multi-Tenant Agents — a different agent per customer

Deploy a **bespoke** agent for each customer, each isolated in that customer's own
ChatBotKit **sub-account** (a "partner user"). Each customer is its own folder/module
with a hand-written agent, and a single shared root composes them all — so one
`terraform apply` deploys every customer's agent, in one shared state.

For the **same** agent across all customers, see the sibling example
[`../multi-tenant-agents-shared`](../multi-tenant-agents-shared).

```
multi-tenant-agents-per-customer/
├── main.tf           # root: provider alias + run_as + module call per customer
├── acme/main.tf      # Acme's bespoke agent (support)
└── globex/main.tf    # Globex's bespoke agent (research + sandbox)
```

`acme` and `globex` are deliberately different shapes — that is the whole point.
Each is a module the root wires to that customer's sub-account.

## How it works: one token + `run_as`

You hold **one** partner/master token (set via `CHATBOTKIT_API_KEY`). The root
configures one provider alias per customer, each with `run_as = var.<customer>_account_id`
(the `X-RunAs-UserId` header), and passes that alias to the customer's module — so
each bespoke agent is created in the right sub-account. No per-customer tokens, no
`for_each`. This mirrors the AWS provider's `assume_role` multi-account pattern.

> Where do account IDs come from? Create a sub-account per customer in the
> dashboard, or via the partner API (`partner/user/create`). Account IDs are not
> secret; the single master token is.

## Usage

```bash
export CHATBOTKIT_API_KEY="sk-...partner-token..."   # one master token

cp terraform.tfvars.example terraform.tfvars         # fill in account IDs
terraform init
terraform apply
```

One apply deploys every customer's bespoke agent into their own sub-account; bot
IDs are in the `bots` output, and everything lives in one shared state.

## Adding a customer

1. Create a `<slug>/main.tf` module describing that customer's agent (copy a
   neighbor and customize).
2. In `main.tf`, add a `<slug>_account_id` variable, a `provider "chatbotkit"`
   alias with `run_as`, and a `module "<slug>"` call wired to that alias.

## Trade-off

Full per-customer flexibility, but every agent is maintained on its own. If every
customer should get the **same** agent you improve in one place, use
[`../multi-tenant-agents-shared`](../multi-tenant-agents-shared).

## Cleanup

`terraform destroy` removes the agents; it does **not** delete the sub-accounts —
remove those via the partner API (`partner/user/{id}/delete`).
