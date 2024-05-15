package cluster_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/omc/bonsai-api-go/v2/bonsai"
	"github.com/omc/terraform-provider-bonsai/internal/test"
)

func testClusterExists(resourceName string, client *bonsai.Client) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		tflog.Info(context.TODO(), fmt.Sprintf("testClusterExists: checking to see that the resource was created...%s", resourceName))

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("testClusterExists: not found: %s", resourceName)
		}

		tflog.Info(context.TODO(), fmt.Sprintf("testClusterExists: Working with terraform instance: %+v", rs.Primary))
		if rs.Primary.ID == "" {
			return errors.New("no cluster ID is set")
		}

		tflog.Info(context.TODO(), fmt.Sprintf("testClusterExists: Fetching with primary terraform instance: %s", rs.Primary.ID))
		_, err := client.Cluster.GetBySlug(context.TODO(), rs.Primary.ID)
		return err
	}
}

func testClusterDestroyed(resourceName string, client *bonsai.Client) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		tflog.Info(context.TODO(), fmt.Sprintf("Checking to see that the resource was destroyed...%s", resourceName))
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return errors.New("no cluster ID is set")
		}
		tflog.Info(context.TODO(), fmt.Sprintf("Working with terraform instance: %+v", rs.Primary))

		result, err := client.Cluster.GetBySlug(context.TODO(), rs.Primary.ID)

		if err != nil && !errors.Is(err, bonsai.ErrHTTPStatusNotFound) {
			return fmt.Errorf("unexpected error checking for deleted cluster: %s", err)
		}

		if err == nil && result.State != bonsai.ClusterStateDeprovisioned {
			return errors.New("expected either a deprovisioned cluster or bonsai.ErrHTTPStatusNotFound error, but got neither")
		}

		return nil
	}
}

func TestCluster_Resource(t *testing.T) {
	clusterSuffix := acctest.RandString(16)
	clusterName := fmt.Sprintf("bonsai test %s", clusterSuffix)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: test.ProtoV6ProviderFactories,
		CheckDestroy:             testClusterDestroyed("bonsai_cluster.test", test.NewApiClient()),
		Steps: []resource.TestStep{
			{
				ResourceName: "bonsai_cluster.test",
				Config: `
                    resource "bonsai_cluster" "test" {
                        name = "never-created-test-cluster"

                        plan = { 
							slug = "invalid-ref"
						}

                        space = { 
							path = "omc/bonsai/us-east-1/common"
						}

                        release = { 
							slug = "opensearch-2.6.0-mt"
						}
                    }
                `,
				// Errors have weird line breaks.
				ExpectError: regexp.MustCompile(`.*?Plan\s+'?[-_a-zA-Z]+'?\s+not\s+found.\s+Please\s+use\s+the\s+/plans\s+endpoint\s+for\s+a\s+list\s+of\s+available\s+plans`), //nolint:lll
			},
			{
				ResourceName: "bonsai_cluster.test",
				Config: `
                    resource "bonsai_cluster" "test" {
                        name = "never-created-test-cluster"

                        plan = { 
							slug = "standard-nano-comped"
						}

                        space = { 
							path = "invalid-ref"
						}

                        release = { 
							slug = "opensearch-2.6.0-mt"
						}
                    }
                `,
				ExpectError: regexp.MustCompile(`.*?Space\s+'?[-_a-zA-Z]+'?\s+not\s+found.\s+Please\s+use\s+the\s+/spaces\s+endpoint\s+for\s+a\s+list\s+of\s+available\s+spaces`), //nolint:lll
			},
			{
				ResourceName: "bonsai_cluster.test",
				Config: `
                    resource "bonsai_cluster" "test" {
                        name = "never-created-test-cluster"

                        plan = { 
							slug = "standard-nano-comped"
						}

                        space = { 
							path = "omc/bonsai/us-east-1/common"
						}

                        release = { 
							slug = "invalid-ref"
						}
                    }
                `,
				ExpectError: regexp.MustCompile(`.*?(Release|Version)\s+'?[-_a-zA-Z]+'?\s+not\s+found.\s+Please\s+use\s+the\s+/releases\s+endpoint\s+for\s+a\s+list\s+of\s+available\s+(releases|versions)`),
			},
			// Create and Read testing
			{
				ResourceName: "bonsai_cluster.test",
				Config: fmt.Sprintf(`
                    resource "bonsai_cluster" "test" {
                        name = "%s"

                        plan = { 
							slug = "standard-nano-comped"
						}

                        space = { 
							path = "omc/bonsai/us-east-1/common"
						}

                        release = { 
							slug = "opensearch-2.6.0-mt"
						}
                    }
                `, clusterName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testClusterExists("bonsai_cluster.test", test.NewApiClient()),
					resource.TestCheckResourceAttr("bonsai_cluster.test",
						"name",
						clusterName,
					),
				),
			},
			// Update testing
			{
				ResourceName: "bonsai_cluster.test",
				Config: fmt.Sprintf(`
			        resource "bonsai_cluster" "test" {
			            name = "bonsai test cluster %s"

			            plan = {
							slug = "standard-nano-comped"
						}

			            space = {
							path = "omc/bonsai/us-east-1/common"
						}

			            release = {
							slug = "opensearch-2.6.0-mt"
						}
			        }
			    `, clusterSuffix),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bonsai_cluster.test", "name", fmt.Sprintf("bonsai test cluster %s", clusterSuffix)),
				),
			},
		},
	})
}
