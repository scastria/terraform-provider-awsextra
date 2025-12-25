package awsextra

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scastria/terraform-provider-awsextra/awsextra/client"
)

func dataSourceECRRepository() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceECRRepositoryRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceECRRepositoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	ecrClient := ecr.NewFromConfig(c.Config)
	resp, err := ecrClient.DescribeRepositories(ctx, &ecr.DescribeRepositoriesInput{RepositoryNames: []string{d.Get("name").(string)}})
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	repo := resp.Repositories[0]
	d.SetId(*repo.RepositoryArn)
	return diags
}
