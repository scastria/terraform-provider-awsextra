package client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

type Client struct {
	region    string
	profile   string
	accessKey string
	secretKey string
	token     string
	Config    aws.Config
}

func NewClient(ctx context.Context, region string, profile string, accessKey string, secretKey string, token string) (*Client, error) {
	c := &Client{
		region:    region,
		profile:   profile,
		accessKey: accessKey,
		secretKey: secretKey,
		token:     token,
	}
	var cfg aws.Config
	var err error
	if profile != "" {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(region), config.WithSharedConfigProfile(profile))
	} else {
		scp := credentials.NewStaticCredentialsProvider(accessKey, secretKey, token)
		cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(region), config.WithCredentialsProvider(scp))
	}
	if err != nil {
		return nil, err
	}
	c.Config = cfg
	return c, nil
}
