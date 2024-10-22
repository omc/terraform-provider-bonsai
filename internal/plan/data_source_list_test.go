package plan_test

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func (s *PlanTestSuite) TestPlan_ListDataSource() {
	resource.Test(s.T(), resource.TestCase{
		ProtoV6ProviderFactories: s.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
						data "bonsai_plans" "list" {}

						output "bonsai_plans" {
						  value = [for s in data.bonsai_plans.list.plans : s.slug]
						}
					`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Keys in map
					resource.TestCheckResourceAttr("data.bonsai_plans.list", "plans.1.%", "8"),
					// Confirm at least one attr
					resource.TestCheckResourceAttr("data.bonsai_plans.list", "plans.1.billing_interval_months", "1"),
				),
			},
		},
	})
}
