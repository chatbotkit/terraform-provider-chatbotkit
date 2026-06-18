# Multi-Tenant Agents — shared agent

Deploy the **same** agent for every customer, each isolated in that customer's own
ChatBotKit **sub-account** (a "partner user"). The agent is defined once in
[`modules/agent`](./modules/agent) and rolled out to every sub-account, so you
maintain it in one place and improvements ship to all customers.

For a **different** agent per customer, see the sibling example
[`../multi-tenant-agents-per-customer`](../multi-tenant-agents-per-customer).

## How it works: one token + `run_as`

You hold **one** partner/master token. The provider's `run_as` attribute selects
which sub-account to operate on (it sends the `X-RunAs-UserId` header), so the
same token manages every customer — no per-customer tokens. This is the standard
Terraform multi-account approach: **provider aliases**, one per account, just like
the AWS provider's `assume_role`.

```hcl
provider "chatbotkit" {
  alias  = "acme"
  run_as = var.acme_account_id   # api_key comes from CHATBOTKIT_API_KEY
}

module "acme" {
  source    = "./modules/agent"
  providers = { chatbotkit = chatbotkit.acme }
  customer_name = "Acme Corporation"
}
```

No `for_each` — one provider alias + one module call per customer.

> Where do account IDs come from? Create a sub-account per customer in the
> dashboard, or via the partner API (`partner/user/create`) if you automate
> onboarding. That is a one-time step, separate from Terraform. Account IDs are
> not secret; the single master token is.

## Usage

```bash
export CHATBOTKIT_API_KEY="sk-...your-partner-token..."   # one master token

cp terraform.tfvars.example terraform.tfvars              # fill in account IDs
terraform init
terraform apply
```

A single apply deploys the agent into every customer's sub-account; bot IDs are
in the `bots` output.

## Adding a customer

Add three things: an `<slug>_account_id` variable, a `provider "chatbotkit"`
alias with `run_as`, and a `module "<slug>"` call wired to that alias.

## Trade-off

Centralized and uniform — great when every customer should get the same agent.
If customers need genuinely different agents that evolve independently, use
[`../multi-tenant-agents-per-customer`](../multi-tenant-agents-per-customer).

## Cleanup

`terraform destroy` removes the agents; it does **not** delete the sub-accounts —
remove those via the partner API (`partner/user/{id}/delete`).
