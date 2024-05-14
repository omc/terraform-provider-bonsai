package release_test

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func (s *ReleaseTestSuite) TestRelease_ListDataSource() {
	resource.Test(s.T(), resource.TestCase{
		ProtoV6ProviderFactories: s.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
						data "bonsai_releases" "list" {}

						output "bonsai_releases" {
						  value = [for s in data.bonsai_releases.list.releases : s.slug]
						}
					`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.bonsai_releases.list", "releases.#", "12"),
					resource.TestCheckResourceAttr("data.bonsai_releases.list", "releases.0.%", "5"),
					resource.TestCheckResourceAttr("data.bonsai_releases.list", "releases.0.slug", "elasticsearch-2.4.0"),
				),
			},
		},
	})
}
