terraform {
  required_providers {
    bonsai = {
      source = "omc/bonsai"
    }
  }
}

// Bonsai Spaces
provider "bonsai" {}

data "bonsai_space" "get_by_path" {
  path = "omc/websolr/us-east-1/common"
}

data "bonsai_spaces" "list" {}

output "bonsai_space" {
  value = data.bonsai_space.get_by_path
}

output "bonsai_spaces" {
  value = data.bonsai_spaces.list
}

// Bonsai Releases
data "bonsai_release" "get_by_slug" {
  slug = "elasticsearch-6.4.2"
}

data "bonsai_releases" "list" {}

output "bonsai_release" {
  value = data.bonsai_release.get_by_slug
}

output "bonsai_releases" {
  value = data.bonsai_releases.list
}

// Bonsai Plans
data "bonsai_plan" "get_by_slug" {
  slug = "sandbox"
}

data "bonsai_plans" "list" {}

output "bonsai_plan" {
  value = data.bonsai_plan.get_by_slug
}

output "bonsai_plans" {
  value = data.bonsai_plans.list
}

// Bonsai Clusters
data "bonsai_cluster" "get_by_slug" {
  slug = "dcek-group-llc-5240651189"
}

data "bonsai_clusters" "list" {}

output "bonsai_cluster" {
  sensitive = true
  value     = data.bonsai_cluster.get_by_slug
}

output "bonsai_clusters" {
  sensitive = true
  value     = data.bonsai_clusters.list
}

# Create a cluster
resource "bonsai_cluster" "test" {
  name = "complete example"

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