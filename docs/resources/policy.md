---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "rode_policy Resource - terraform-provider-rode"
subcategory: ""
description: |-
  An Open Policy Agent Rego policy.
---

# rode_policy (Resource)

An Open Policy Agent Rego policy.

## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **name** (String) Policy name
- **rego_content** (String) The Rego code

### Optional

- **description** (String) A brief summary of the policy
- **id** (String) The ID of this resource.
- **message** (String) A summary of changes since the last version

### Read-Only

- **created** (String) Creation timestamp
- **current_version** (Number) Current version of the policy
- **deleted** (Boolean)
- **policy_version_id** (String) Policy version id
- **updated** (String) Last updated timestamp

