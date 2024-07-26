# Configure the Bonsai Provider using the required_providers stanza.
terraform {
  required_providers {
    bonsai = {
      source  = "omc/bonsai"
      version = "~> 1.0"
    }
  }
}

provider "bonsai" {
  # Optionally omit this entry to get the value from the BONSAI_API_KEY
  # environment variable.
  api_key = var.bonsai_api_key

  # Optionally omit this entry to get the value from the BONSAI_API_TOKEN
  # environment variable.
  api_token = var.bonsai_api_token
}
