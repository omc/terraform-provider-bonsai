package plan_test

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func (s *PlanTestSuite) TestPlan_DataSource() {
	const planSlug = "standard-micro-aws-us-east-1"
	resource.Test(s.T(), resource.TestCase{
		ProtoV6ProviderFactories: s.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
						data "bonsai_plan" "get_by_slug" {
						  slug = "%s"
						}

						output "bonsai_plan_slug" {
						  value = data.bonsai_plan.get_by_slug.slug
						}
					`, planSlug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bonsai_plan.get_by_slug", "slug"),
					resource.TestCheckResourceAttr("data.bonsai_plan.get_by_slug", "slug", planSlug),
					resource.TestCheckResourceAttr("data.bonsai_plan.get_by_slug", "name", "Standard Micro"),
					resource.TestCheckOutput("bonsai_plan_slug", planSlug),
				),
			},
		},
	})
}
