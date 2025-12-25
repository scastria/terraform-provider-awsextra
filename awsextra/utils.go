package awsextra

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func SuppressMissingOptionalConfigurationBlock(k, old, new string, d *schema.ResourceData) bool {
	return old == "1" && new == "0"
}
