# AWS Extra Provider
The AWS Extra provider extends the official AWS provider with `use_existing` flags to make it easier
to import existing resources without having to run `terraform import`.
## Example Usage
```hcl
terraform {
  required_providers {
    awsextra = {
      source  = "scastria/awsextra"
      version = "~> 0.1.0"
    }
  }
}

# Configure the AWS Extra Provider
provider "awsextra" {
  region = "us-west-2"
}
```
## Argument Reference
* `region` - (Optional) AWS Region where the provider will operate. The Region must be set.
  Can also be set with either the `AWS_REGION` or `AWS_DEFAULT_REGION` environment variables,
  or via a shared config file parameter `region` if `profile` is used.
  If credentials are retrieved from the EC2 Instance Metadata Service, the Region can also be retrieved from the metadata.
  Most Regional resources, data sources and ephemeral resources support an optional top-level `region` argument which can be used to override the provider configuration value. See the individual resource's documentation for details.
* `profile` - (Optional) AWS profile name as set in the shared configuration and credentials files.
  Can also be set using either the environment variables `AWS_PROFILE` or `AWS_DEFAULT_PROFILE`.
* `access_key` - (Optional) AWS access key. Can also be set with the `AWS_ACCESS_KEY_ID` environment variable, or via a shared credentials file if `profile` is specified. See also `secret_key`.
* `secret_key` - (Optional) AWS secret key. Can also be set with the `AWS_SECRET_ACCESS_KEY` environment variable, or via a shared configuration and credentials files if `profile` is used. See also `access_key`.
* `token` - (Optional) Session token for validating temporary credentials. Typically provided after successful identity federation or Multi-Factor Authentication (MFA) login. With MFA login, this is the session token provided afterward, not the 6 digit MFA code used to get temporary credentials.  Can also be set with the `AWS_SESSION_TOKEN` environment variable.
