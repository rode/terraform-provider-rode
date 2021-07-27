terraform {
  required_providers {
    rode = {
      source  = "rode/rode"
      version = "0.0.1"
    }
  }
}

provider "rode" {
  host                          = "localhost:50051"
  disable_transport_security    = true
}

resource "rode_policy_group" "example" {
  name        = "terraform-example"
  description = "managed by Terraform"
}
