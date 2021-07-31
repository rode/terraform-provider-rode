resource "rode_policy" "example" {
  name         = "example"
  description  = "policy managed by Terraform"
  message      = "Terraform"
  rego_content = <<EOF
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
