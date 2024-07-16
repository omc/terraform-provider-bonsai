package cluster_test

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func (s *ClusterTestSuite) TestCluster_ListDataSource() {
	resource.Test(s.T(), resource.TestCase{
		ProtoV6ProviderFactories: s.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
						resource "bonsai_cluster" "test" {
							name = "%s"

							plan = {
								slug = "sandbox"
							}

							space = {
								path = "omc/bonsai/us-east-1/common"
							}

							release = {
								slug = "opensearch-2.6.0-mt"
							}
						}

						data "bonsai_clusters" "list" {}

						output "bonsai_clusters" {
						  value = [for c in data.bonsai_clusters.list.clusters : c.slug]
						}
					`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.bonsai_clusters.list", "clusters.0.%", "9"),
				),
			},
		},
	})
}
