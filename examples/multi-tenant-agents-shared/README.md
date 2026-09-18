# Multi-Tenant Agents - shared agent

Deploy the **same** agent for every customer, each isolated in that customer's own
ChatBotKit **child User**. The agent is defined once in
[`modules/agent`](./modules/agent) and rolled out to every child User, so you
maintain it in one place and improvements ship to all customers.

For a **different** agent per customer, see the sibling example
[`../multi-tenant-agents-per-customer`](../multi-tenant-agents-per-customer).

## How it works: one API token plus `run_as`

You hold **one** API token belonging to the parent User. The provider's `run_as`
attribute selects which child User to operate on by sending the
`X-RunAs-UserId` header, so the same API token manages every customer. This uses
Terraform's standard multi-account approach: **provider aliases**, one per child
User, similar to the AWS provider's `assume_role`.

```hcl
provider "chatbotkit" {
  alias  = "acme"
  run_as = var.acme_account_id   # api_token comes from CHATBOTKIT_API_TOKEN
}

module "acme" {
  source    = "./modules/agent"
  providers = { chatbotkit = chatbotkit.acme }
  customer_name = "Acme Corporation"
}
```

No `for_each` is needed: use one provider alias and one module call per customer.

> Where do User IDs come from? Create a child User per customer in the dashboard
> or through the User API (`user/create`) if you automate onboarding. That is a
> one-time step separate from Terraform. User IDs are not secret; the parent
> User API token is.

## Usage

```bash
export CHATBOTKIT_API_TOKEN="sk-...your-parent-user-api-token..." # one parent User API token

cp terraform.tfvars.example terraform.tfvars              # fill in User IDs
terraform init
terraform apply
```

A single apply deploys the agent into every customer's child User; bot IDs are
in the `bots` output.

## Adding a customer

Add three things: an `<slug>_account_id` variable, a `provider "chatbotkit"`
alias with `run_as`, and a `module "<slug>"` call wired to that alias.

## Trade-off

Centralized and uniform, which is useful when every customer should get the same agent.
If customers need genuinely different agents that evolve independently, use
[`../multi-tenant-agents-per-customer`](../multi-tenant-agents-per-customer).

## Cleanup

`terraform destroy` removes the agents; it does **not** delete the child Users.
Remove those through the User API (`user/{id}/delete`).
