package storageprovider

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

func r2provider() aws.Config {
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", os.Getenv("CLOUDFLARE_ACCOUNT_ID")),
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(os.Getenv("CLOUDFLARE_R2_ACCESS_KEY_ID"), os.Getenv("CLOUDFLARE_R2_ACCESS_KEY_SECRET"), "")),
	)
	if err != nil {
		panic(err)
	}
	return cfg
}

func s3provider() aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}
	return cfg
}
