# terraform-provider-rode

A Terraform provider for [rode](https://github.com/rode/rode).

## Resources

- [x] `rode_policy`*
- [x] `rode_policy_group`*
- [ ] `rode_policy_assignment`

> \* Implemented for create, read, & delete. Missing update and import

## Data Sources

- [ ] `rode_policy`
- [ ] `rode_policy_group`

## Local Development

To run the examples, configure `~/.terraformrc` to use your local provider:

```
provider_installation {
    dev_overrides {
      "registry.terraform.io/rode/rode" = "/my/path/to/terraform-provider-rode"
    }
}
```

Then run `make example`.