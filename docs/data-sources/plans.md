---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "bonsai_plans Data Source - terraform-provider-bonsai"
subcategory: ""
description: |-
  A list of all available plans on your account.
---

# bonsai_plans (Data Source)

A list of all **available** plans on your account.

## Example Usage

```terraform
data "bonsai_plans" "list" {}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `plans` (Attributes List) Plan represents a subscription plan. (see [below for nested schema](#nestedatt--plans))

<a id="nestedatt--plans"></a>
### Nested Schema for `plans`

Optional:

- `slug` (String) The machine-readable name for the plan.

Read-Only:

- `available_releases` (Attributes List) A collection of search release slugs available for the plan. (see [below for nested schema](#nestedatt--plans--available_releases))
- `available_spaces` (Attributes List) A collection of Space paths available for the plan. (see [below for nested schema](#nestedatt--plans--available_spaces))
- `billing_interval_months` (Number) The plan billing interval in months.
- `name` (String) The human-readable name of the plan.
- `price_in_cents` (Number) Represents the plan price in cents.
- `private_network` (Boolean) Indicates whether the plan is on a publicly addressable network. Private plans provide environments that cannot be reached by the public Internet. A VPC connection will be needed to communicate with a private cluster.
- `single_tenant` (Boolean) Indicates whether the plan is single-tenant or not. A value of false indicates the Cluster will share hardware with other Clusters. Single tenant environments can be reached via the public Internet.

<a id="nestedatt--plans--available_releases"></a>
### Nested Schema for `plans.available_releases`

Read-Only:

- `slug` (String) A machine-readable name for the release.


<a id="nestedatt--plans--available_spaces"></a>
### Nested Schema for `plans.available_spaces`

Read-Only:

- `path` (String) A machine-readable name for the server group.
