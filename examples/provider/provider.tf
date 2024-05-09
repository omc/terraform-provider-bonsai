terraform {
  required_providers {
    bonsai = {
      source = "omc/bonsai"
    }
  }
}

provider "bonsai" {}

data "bonsai_clusters" "example" {}
