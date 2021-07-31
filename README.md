# terraform-provider-rode

[![test badge](https://github.com/alexashley/terraform-provider-rode/actions/workflows/test.yaml/badge.svg?branch=main)](https://github.com/alexashley/terraform-provider-rode/actions/workflows/test.yaml?query=branch%3Amain)


A Terraform provider for [rode](https://github.com/rode/rode).

## Resources

- [x] `rode_policy`*
- [x] `rode_policy_group`*
- [x] `rode_policy_assignment`*

> \* Implemented for create, read, & delete. Missing update and import 

See the [example](./example) directory for resource usage.

## Local Development

To build the provider, run `make build`; or run `make install` to build and move the binary under `~/.terraform.d`.

To run the acceptance tests, use `make testacc`. These require a running instance of Rode.

### Test Environment

If you have access to a Kubernetes cluster, the `services` directory contains Terraform for standing up Elasticsearch,
Grafeas, and Rode. 

```shell
terraform -chdir=services init
terraform -chdir=services apply
```

### Developer Overrides

To run the examples with a local provider, configure `~/.terraformrc` to use [dev overrides](https://www.terraform.io/docs/cli/config/config-file.html#development-overrides-for-provider-developers):

```
provider_installation {
    dev_overrides {
      "registry.terraform.io/rode/rode" = "/my/path/to/terraform-provider-rode"
    }
}
```
