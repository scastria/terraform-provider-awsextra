# Resource: awsextra_ecr_repository
Represents an ECR repository
## Example usage
```hcl
resource "awsextra_ecr_repository" "example" {
  name                 = "bar"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
}
```
## Argument Reference
* `name` - **(Required, ForceNew, String)** Name of the repository.
* `encryption_configuration` - **(Optional, ForceNew, Record)** Encryption configuration for the repository. See [below for schema](#encryption_configuration).
* `force_delete` - **(Optional, Boolean)** If `true`, will delete the repository even if it contains images.
  Defaults to `false`.
* `image_tag_mutability` - **(Optional, String)** The tag mutability setting for the repository. Must be one of: `MUTABLE`, `IMMUTABLE`, `IMMUTABLE_WITH_EXCLUSION`, or `MUTABLE_WITH_EXCLUSION`. Defaults to `MUTABLE`.
* `image_scanning_configuration` - **(Optional, Record)** Configuration block that defines image scanning configuration for the repository. By default, image scanning must be manually triggered. See the [ECR User Guide](https://docs.aws.amazon.com/AmazonECR/latest/userguide/image-scanning.html) for more information about image scanning.
    * `scan_on_push` - **(Required, Boolean)** Indicates whether images are scanned after being pushed to the repository (true) or not scanned (false).
* `use_existing` - **(Optional, Boolean, IgnoreDiffs)** During a CREATE only, look for an existing repository with the same `name`.  Prevents the need for an import. Default: `false`
### encryption_configuration
* `encryption_type` - **(Optional, ForceNew, String)** The encryption type to use for the repository. Valid values are `AES256` or `KMS`. Defaults to `AES256`.
## Attribute Reference
* `id` - **(String)** Same as `name`
* `arn` - **(String)** Full ARN of the repository.
* `registry_id` - **(String)** The registry ID where the repository was created.
* `repository_url` - **(String)** The URL of the repository (in the form `aws_account_id.dkr.ecr.region.amazonaws.com/repositoryName`).
## Import
Repositories can be imported using a proper value of `id` as described above
