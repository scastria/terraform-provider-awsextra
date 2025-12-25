package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/scastria/terraform-provider-awsextra/awsextra"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: awsextra.Provider,
	})
}
