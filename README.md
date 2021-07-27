# terraform-provider-rode

A Terraform provider for [rode](https://github.com/rode/rode).

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