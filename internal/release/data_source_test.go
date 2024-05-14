package release_test

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func (s *ReleaseTestSuite) TestRelease_DataSource() {
	const releaseSlug = "elasticsearch-5.6.16"
	resource.Test(s.T(), resource.TestCase{
		ProtoV6ProviderFactories: s.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
						data "bonsai_release" "get_by_slug" {
						  slug = "%s"
						}

						output "bonsai_release_slug" {
						  value = data.bonsai_release.get_by_slug.slug
						}
					`, releaseSlug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bonsai_release.get_by_slug", "slug"),
					resource.TestCheckResourceAttr("data.bonsai_release.get_by_slug", "slug", releaseSlug),
					resource.TestCheckOutput("bonsai_release_slug", releaseSlug),
				),
			},
		},
	})
}
