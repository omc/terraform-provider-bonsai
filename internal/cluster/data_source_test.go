package cluster_test

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func (s *ClusterTestSuite) TestCluster_DataSource() {
	const clusterSlug = "dcek-group-llc-5240651189"
	resource.Test(s.T(), resource.TestCase{
		ProtoV6ProviderFactories: s.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
						data "bonsai_cluster" "get_by_slug" {
						  slug = "%s"
						}

						output "bonsai_cluster_slug" {
						  value = data.bonsai_cluster.get_by_slug.slug
						}
					`, clusterSlug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bonsai_cluster.get_by_slug", "slug"),
					resource.TestCheckResourceAttr("data.bonsai_cluster.get_by_slug", "slug", clusterSlug),
					resource.TestCheckResourceAttr("data.bonsai_cluster.get_by_slug", "name", "DCEK Group, LLC search"),
					resource.TestCheckOutput("bonsai_cluster_slug", clusterSlug),
				),
			},
		},
	})
}
