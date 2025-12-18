package s3

import (
	"context"
	"io"
	"murky_api/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client interface {
	Upload(key string, bodyReader io.Reader) error
	ListObjectKeys(path string) ([]string, error)
	DeleteObject(key string) error
}

type defaultClient struct {
	ctx      context.Context
	s3Client *s3.Client
	bucket   string
}

func NewClient(conf *config.Config) (Client, error) {
	ctx := context.TODO()

	// Credentials are loaded from env / shared config
	cfg, err := awsConfig.LoadDefaultConfig(ctx, awsConfig.WithRegion(conf.AwsRegion))
	if err != nil {
		return nil, err
	}

	s3client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(conf.AwsEndpointUrl)
		// many S3-compatible servers require path-style addressing
		o.UsePathStyle = true
	})

	client := &defaultClient{
		ctx:      ctx,
		s3Client: s3client,
		bucket:   conf.S3Bucket,
	}

	return client, nil
}

func (client *defaultClient) Upload(key string, bodyReader io.Reader) error {
	_, err := client.s3Client.PutObject(client.ctx, &s3.PutObjectInput{
		Bucket: &client.bucket,
		Key:    &key,
		Body:   bodyReader,
	})
	return err
}

func (client *defaultClient) ListObjectKeys(path string) ([]string, error) {
	objects, err := client.s3Client.ListObjectsV2(client.ctx, &s3.ListObjectsV2Input{
		Bucket: &client.bucket,
		Prefix: &path,
	})
	if err != nil {
		return nil, err
	}

	keys := make([]string, len(objects.Contents))
	for i, object := range objects.Contents {
		keys[i] = *object.Key
	}

	return keys, nil
}

func (client *defaultClient) DeleteObject(key string) error {
	_, err := client.s3Client.DeleteObject(client.ctx, &s3.DeleteObjectInput{
		Bucket: &client.bucket,
		Key:    &key,
	})
	return err
}
