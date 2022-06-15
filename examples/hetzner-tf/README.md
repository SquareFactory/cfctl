# Bootstrapping k0s cluster on Hetzner

This directory provides an example flow with `cfctl` tool together with Terraform using Hetzner as the cloud provider.

## Prerequisites

- You need an account and API token for Hetzner
- Terraform installed
- cfctl installed

## Steps

Create terraform.tfvars file with needed details. You can use the provided terraform.tfvars.example as a baseline.

- `terraform init`
- `terraform apply`
- `terraform output -raw k0s_cluster | cfctl apply --config -`
