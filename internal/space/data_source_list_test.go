package space_test

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func (s *SpaceTestSuite) TestSpace_ListDataSource() {
	resource.Test(s.T(), resource.TestCase{
		ProtoV6ProviderFactories: s.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
						data "bonsai_spaces" "list" {}

						output "bonsai_spaces" {
						  value = [for s in data.bonsai_spaces.list.spaces : s.path]
						}
					`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.bonsai_spaces.list", "spaces.0.%", "3"),
					resource.TestCheckResourceAttr("data.bonsai_spaces.list", "spaces.0.path", "omc/bonsai/eu-west-1/common"),
				),
			},
		},
	})
}
