resource "bonsai_cluster" "test" {
  name = "comped example"

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