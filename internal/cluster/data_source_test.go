package cluster_test

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func (s *ClusterTestSuite) TestCluster_DataSource() {
	clusterSuffix := acctest.RandString(16)
	clusterName := fmt.Sprintf("bonsai test %s", clusterSuffix)

	resource.Test(s.T(), resource.TestCase{
		ProtoV6ProviderFactories: s.ProtoV6ProviderFactories,
		CheckDestroy:             testClusterDestroyed("bonsai_cluster.test", s.Client),
		Steps: []resource.TestStep{
			{
				ResourceName: "bonsai_cluster.test",
				Config: fmt.Sprintf(`
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

						data "bonsai_cluster" "get_by_slug" {
						  slug = resource.bonsai_cluster.test.id
						}

						output "bonsai_cluster_slug" {
						  value = data.bonsai_cluster.get_by_slug.slug
						}
					`, clusterName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bonsai_cluster.get_by_slug", "slug"),
				),
			},
		},
	})
}
