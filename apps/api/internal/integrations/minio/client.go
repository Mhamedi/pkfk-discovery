package minio

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	client *minio.Client
	bucketAdapterBundles string
	bucketArtifacts      string
}

func NewClient(endpoint, accessKey, secretKey string, useSSL bool) (*Client, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	c := &Client{
		client:               client,
		bucketAdapterBundles: "adapter-bundles",
		bucketArtifacts:      "artifacts",
	}

	ctx := context.Background()

	// Ensure buckets exist
	if err := c.ensureBucket(ctx, c.bucketAdapterBundles); err != nil {
		return nil, fmt.Errorf("failed to create adapter-bundles bucket: %w", err)
	}

	if err := c.ensureBucket(ctx, c.bucketArtifacts); err != nil {
		return nil, fmt.Errorf("failed to create artifacts bucket: %w", err)
	}

	return c, nil
}

func (c *Client) ensureBucket(ctx context.Context, bucketName string) error {
	exists, err := c.client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}

	if !exists {
		return c.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	}

	return nil
}

func (c *Client) UploadAdapterBundle(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) error {
	_, err := c.client.PutObject(ctx, c.bucketAdapterBundles, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func (c *Client) DownloadAdapterBundle(ctx context.Context, objectName string) (io.ReadCloser, error) {
	return c.client.GetObject(ctx, c.bucketAdapterBundles, objectName, minio.GetObjectOptions{})
}

func (c *Client) UploadArtifact(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) error {
	_, err := c.client.PutObject(ctx, c.bucketArtifacts, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func (c *Client) GetPresignedURL(ctx context.Context, bucket, objectName string, expiry time.Duration) (string, error) {
	return c.client.PresignedGetObject(ctx, bucket, objectName, expiry, nil)
}

