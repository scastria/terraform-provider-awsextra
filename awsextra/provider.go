package awsextra

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scastria/terraform-provider-awsextra/awsextra/client"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"AWS_REGION", "AWS_DEFAULT_REGION"}, nil),
			},
			"profile": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.MultiEnvDefaultFunc([]string{"AWS_PROFILE", "AWS_DEFAULT_PROFILE"}, nil),
				ConflictsWith: []string{"access_key", "secret_key", "token"},
			},
			"access_key": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc("AWS_ACCESS_KEY_ID", nil),
				ConflictsWith: []string{"profile"},
				RequiredWith:  []string{"secret_key"},
			},
			"secret_key": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc("AWS_SECRET_ACCESS_KEY", nil),
				ConflictsWith: []string{"profile"},
				RequiredWith:  []string{"access_key"},
			},
			"token": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc("AWS_SESSION_TOKEN", nil),
				ConflictsWith: []string{"profile"},
			},
		},
		ResourcesMap: map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{
			"awsextra_ecr_repository": dataSourceECRRepository(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	region := d.Get("region").(string)
	profile := d.Get("profile").(string)
	accessKey := d.Get("access_key").(string)
	secretKey := d.Get("secret_key").(string)
	token := d.Get("token").(string)

	//Check for valid authentication
	if (profile == "") && (accessKey == "") && (secretKey == "") && (token == "") {
		return nil, diag.Errorf("You must specify either profile or access_key/secret_key for authentication")
	}

	var diags diag.Diagnostics
	c, err := client.NewClient(ctx, region, profile, accessKey, secretKey, token)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return c, diags
}
