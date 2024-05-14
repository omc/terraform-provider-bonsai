package space_test

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func (s *SpaceTestSuite) TestSpace_DataSource() {
	const spacePath = "omc/bonsai/eu-west-1/common"
	resource.Test(s.T(), resource.TestCase{
		ProtoV6ProviderFactories: s.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
						data "bonsai_space" "get_by_path" {
						  path = "%s"
						}

						output "bonsai_space_path" {
						  value = data.bonsai_space.get_by_path.path
						}
					`, spacePath),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bonsai_space.get_by_path", "path"),
					resource.TestCheckResourceAttr("data.bonsai_space.get_by_path", "path", spacePath),
					resource.TestCheckOutput("bonsai_space_path", spacePath),
				),
			},
		},
	})
}
