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

resource "rode_policy" "example" {
  name = "example"
  description = "policy managed by Terraform"
  policy {
    message = "Terraform"
    rego_content =<<EOF
package tf_example

pass {
    true
}

violations[result] {
	result = {
		"pass": true,
		"id": "valid",
		"name": "name",
		"description": "description",
		"message": "message",
	}
}
EOF
  }
}