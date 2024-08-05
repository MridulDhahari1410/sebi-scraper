package s3

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Config - configuration for the S3 connection.
type Config struct {
	Bucket     string
	MaxRetries int
}

type client struct {
	config Config
	*s3.Client
}

var c *client

// Init initializes and returns an S3 client.
func Init(ctx context.Context, conf Config) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	cfg.RetryMode = aws.RetryModeStandard
	cfg.RetryMaxAttempts = conf.MaxRetries
	if err != nil {
		return err
	}
	c = &client{
		config: conf,
		Client: s3.NewFromConfig(cfg),
	}
	return nil
}

// Upload uploads data to the specified S3 bucket and file.
func Upload(ctx context.Context, key string, data io.Reader) error {
	_, err := c.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(c.config.Bucket),
		Key:    aws.String(key),
		Body:   data,
	})
	return err
}
