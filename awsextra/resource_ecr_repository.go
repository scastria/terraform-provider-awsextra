package awsextra

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scastria/terraform-provider-awsextra/awsextra/client"
)

func resourceECRRepository() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceECRRepositoryCreate,
		ReadContext:   resourceECRRepositoryRead,
		UpdateContext: resourceECRRepositoryUpdate,
		DeleteContext: resourceECRRepositoryDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"image_tag_mutability": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "MUTABLE",
				ValidateFunc: validation.StringInSlice([]string{"MUTABLE", "IMMUTABLE", "IMMUTABLE_WITH_EXCLUSION"}, false),
			},
			"image_scanning_configuration": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scan_on_push": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
				DiffSuppressFunc: SuppressMissingOptionalConfigurationBlock,
			},
			"encryption_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"encryption_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      "AES256",
							ValidateFunc: validation.StringInSlice([]string{"AES256", "KMS"}, false),
						},
					},
				},
				DiffSuppressFunc: SuppressMissingOptionalConfigurationBlock,
				ForceNew:         true,
			},
			"force_delete": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			//"tags": {
			//	Type:     schema.TypeMap,
			//	Optional: true,
			//	Elem: &schema.Schema{
			//		Type: schema.TypeString,
			//	},
			//},
			"use_existing": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          false,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool { return d.Id() != "" },
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"registry_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"repository_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func flattenImageScanningConfiguration(isc *types.ImageScanningConfiguration) []map[string]any {
	if isc == nil {
		return nil
	}
	config := make(map[string]any)
	config["scan_on_push"] = isc.ScanOnPush
	return []map[string]any{
		config,
	}
}

func flattenRepositoryEncryptionConfiguration(ec *types.EncryptionConfiguration) []map[string]any {
	if ec == nil {
		return nil
	}
	config := map[string]any{
		"encryption_type": ec.EncryptionType,
	}
	return []map[string]any{
		config,
	}
}

func expandRepositoryEncryptionConfiguration(data []any) *types.EncryptionConfiguration {
	if len(data) == 0 || data[0] == nil {
		return nil
	}
	ec := data[0].(map[string]any)
	config := &types.EncryptionConfiguration{
		EncryptionType: types.EncryptionType((ec["encryption_type"].(string))),
	}
	return config
}

//func convertTagMapToArray(tagMap map[string]interface{}) []types.Tag {
//	retVal := []types.Tag{}
//	for key, value := range tagMap {
//		retVal = append(retVal, types.Tag{
//			Key:   aws.String(key),
//			Value: aws.String(value.(string)),
//		})
//	}
//	return retVal
//}
//
//func convertArrayToTagMap(tags []types.Tag) map[string]string {
//	retVal := map[string]string{}
//	for _, tag := range tags {
//		retVal[aws.ToString(tag.Key)] = aws.ToString(tag.Value)
//	}
//	return retVal
//}

func fillECRRepository(c *ecr.CreateRepositoryInput, d *schema.ResourceData) {
	c.RepositoryName = aws.String(d.Get("name").(string))
	c.ImageTagMutability = types.ImageTagMutability(d.Get("image_tag_mutability").(string))
	if v, ok := d.GetOk("image_scanning_configuration"); ok && len(v.([]any)) > 0 && v.([]any)[0] != nil {
		tfMap := v.([]any)[0].(map[string]any)
		c.ImageScanningConfiguration = &types.ImageScanningConfiguration{
			ScanOnPush: tfMap["scan_on_push"].(bool),
		}
	}
	c.EncryptionConfiguration = expandRepositoryEncryptionConfiguration(d.Get("encryption_configuration").([]any))
	//tagsMap, ok := d.GetOk("tags")
	//if ok {
	//	c.Tags = convertTagMapToArray(tagsMap.(map[string]interface{}))
	//}
}

func fillResourceDataFromECRRepository(c *types.Repository /*, tags []types.Tag*/, d *schema.ResourceData) {
	d.Set("name", c.RepositoryName)
	d.Set("image_tag_mutability", c.ImageTagMutability)
	d.Set("image_scanning_configuration", flattenImageScanningConfiguration(c.ImageScanningConfiguration))
	d.Set("encryption_configuration", flattenRepositoryEncryptionConfiguration(c.EncryptionConfiguration))
	//d.Set("tags", convertArrayToTagMap(tags))
	d.Set("arn", c.RepositoryArn)
	d.Set("registry_id", c.RegistryId)
	d.Set("repository_url", c.RepositoryUri)
}

func resourceECRRepositoryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	ecrClient := ecr.NewFromConfig(c.Config)
	name := d.Get("name").(string)
	useExisting := d.Get("use_existing").(bool)
	var repo *types.Repository = nil
	if useExisting {
		// Try to find an existing repository with the given name and return it if found
		resp, err := ecrClient.DescribeRepositories(ctx, &ecr.DescribeRepositoriesInput{RepositoryNames: []string{name}})
		if err != nil {
			var rnf *types.RepositoryNotFoundException
			if !errors.As(err, &rnf) {
				d.SetId("")
				return diag.FromErr(err)
			}
		} else {
			repo = &(resp.Repositories[0])
		}
	}
	if repo == nil {
		newECRRepository := ecr.CreateRepositoryInput{}
		fillECRRepository(&newECRRepository, d)
		resp, err := ecrClient.CreateRepository(ctx, &newECRRepository)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
		repo = resp.Repository
	}
	//tagResp, err := ecrClient.ListTagsForResource(ctx, &ecr.ListTagsForResourceInput{ResourceArn: repo.RepositoryArn})
	//if err != nil {
	//	d.SetId("")
	//	return diag.FromErr(err)
	//}
	fillResourceDataFromECRRepository(repo /*, tagResp.Tags*/, d)
	d.SetId(aws.ToString(repo.RepositoryName))
	return diags
}

func resourceECRRepositoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	ecrClient := ecr.NewFromConfig(c.Config)
	resp, err := ecrClient.DescribeRepositories(ctx, &ecr.DescribeRepositoriesInput{RepositoryNames: []string{d.Id()}})
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	repo := resp.Repositories[0]
	//tagResp, err := ecrClient.ListTagsForResource(ctx, &ecr.ListTagsForResourceInput{ResourceArn: repo.RepositoryArn})
	//if err != nil {
	//	d.SetId("")
	//	return diag.FromErr(err)
	//}
	fillResourceDataFromECRRepository(&repo /*, tagResp.Tags*/, d)
	return diags
}

func resourceECRRepositoryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	ecrClient := ecr.NewFromConfig(c.Config)
	if d.HasChange("image_tag_mutability") {
		input := &ecr.PutImageTagMutabilityInput{
			ImageTagMutability: types.ImageTagMutability((d.Get("image_tag_mutability").(string))),
			RegistryId:         aws.String(d.Get("registry_id").(string)),
			RepositoryName:     aws.String(d.Id()),
		}
		_, err := ecrClient.PutImageTagMutability(ctx, input)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("image_scanning_configuration") {
		input := &ecr.PutImageScanningConfigurationInput{
			ImageScanningConfiguration: &types.ImageScanningConfiguration{},
			RegistryId:                 aws.String(d.Get("registry_id").(string)),
			RepositoryName:             aws.String(d.Id()),
		}
		if v, ok := d.GetOk("image_scanning_configuration"); ok && len(v.([]any)) > 0 && v.([]any)[0] != nil {
			tfMap := v.([]any)[0].(map[string]any)
			input.ImageScanningConfiguration.ScanOnPush = tfMap["scan_on_push"].(bool)
		}
		_, err := ecrClient.PutImageScanningConfiguration(ctx, input)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return diags
}

func resourceECRRepositoryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	ecrClient := ecr.NewFromConfig(c.Config)
	_, err := ecrClient.DeleteRepository(ctx, &ecr.DeleteRepositoryInput{
		Force:          d.Get("force_delete").(bool),
		RegistryId:     aws.String(d.Get("registry_id").(string)),
		RepositoryName: aws.String(d.Id()),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
