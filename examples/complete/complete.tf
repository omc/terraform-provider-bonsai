terraform {
  required_providers {
    bonsai = {
      source = "omc/bonsai"
    }
  }
}

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
